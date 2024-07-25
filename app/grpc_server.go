package app

import (
	grpcservice "buf.build/gen/go/webitel/cases/grpc/go/_gogrpc"
	"context"
	"errors"
	"fmt"
	"github.com/webitel/cases/app/lookup"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/registry"
	"github.com/webitel/cases/registry/consul"
	"github.com/webitel/wlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	RequestContextName       = "grpc_ctx"
	AppServiceTtl            = time.Second * 30
	AppDeregesterCriticalTtl = time.Minute * 2
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
	config   *model.ConsulConfig
	exitChan chan model.AppError
	registry registry.ServiceRegistrator
}

func BuildServer(app *App, config *model.ConsulConfig, exitChan chan model.AppError) (*Server, model.AppError) {
	// * Build grpc server
	server, appErr := buildGrpc(app)
	if appErr != nil {
		return nil, appErr
	}
	//  * Open tcp connection
	listener, err := net.Listen("tcp", config.PublicAddress)
	if err != nil {
		return nil, model.NewInternalError("api.grpc_server.serve_requests.listen.error", err.Error())
	}
	reg, appErr := consul.NewConsulRegistry(config)
	if appErr != nil {
		return nil, appErr
	}

	return &Server{
		server:   server,
		listener: listener,
		exitChan: exitChan,
		config:   config,
		registry: reg,
	}, nil
}

func (a *Server) Start() {
	appErr := a.registry.Register()
	if appErr != nil {
		a.exitChan <- appErr
		return
	}
	err := a.server.Serve(a.listener)
	if err != nil {
		a.exitChan <- model.NewInternalError("api.grpc_server.serve_requests.serve.error", err.Error())
		return
	}
}

func (a *Server) Stop() {
	appErr := a.registry.Deregister()
	if appErr != nil {
		a.exitChan <- appErr
		return
	}
	a.server.Stop()
}

func buildGrpc(app *App) (*grpc.Server, model.AppError) {

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(unaryInterceptor))
	// * Creating services
	l, appErr := lookup.NewAppealLookupService(app)
	if appErr != nil {
		return nil, appErr
	}
	c, appErr := lookup.NewStatusLookupService(app)
	if appErr != nil {
		return nil, appErr
	}

	n, appErr := lookup.NewCloseReasonLookupService(app)
	if appErr != nil {
		return nil, appErr
	}

	// * register appeal service
	grpcservice.RegisterAppealLookupsServer(grpcServer, l)
	// * register status service
	grpcservice.RegisterStatusLookupsServer(grpcServer, c)
	// * register close reason service
	grpcservice.RegisterCloseReasonLookupsServer(grpcServer, n)

	return grpcServer, nil

}

func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	var reqCtx context.Context
	var ip string

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqCtx = context.WithValue(ctx, RequestContextName, md)
		ip = getClientIp(md)
	} else {
		ip = "<not found>"
		reqCtx = context.WithValue(ctx, RequestContextName, nil)
	}

	h, err := handler(reqCtx, req)

	if err != nil {
		wlog.Error(fmt.Sprintf("[%s] method %s duration %s, error: %v", ip, info.FullMethod, time.Since(start), err.Error()))

		var appError model.AppError
		switch {
		case errors.As(err, &appError):
			var e model.AppError
			errors.As(err, &e)
			return h, status.Error(httpCodeToGrpc(e.GetStatusCode()), e.ToJson())
		default:
			return h, err
		}
	} else {
		wlog.Debug(fmt.Sprintf("[%s] method %s duration %s", ip, info.FullMethod, time.Since(start)))
	}

	return h, err
}

func httpCodeToGrpc(c int) codes.Code {
	switch c {
	case http.StatusBadRequest:
		return codes.InvalidArgument
	case http.StatusAccepted:
		return codes.ResourceExhausted
	case http.StatusUnauthorized:
		return codes.Unauthenticated
	case http.StatusForbidden:
		return codes.PermissionDenied
	default:
		return codes.Internal
	}
}

func getClientIp(info metadata.MD) string {
	ip := strings.Join(info.Get("x-real-ip"), ",")
	if ip == "" {
		ip = strings.Join(info.Get("x-forwarded-for"), ",")
	}

	return ip
}
