package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	grpcservice "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/registry"
	"github.com/webitel/cases/registry/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	// Appeal Lookup service
	l, appErr := NewAppealService(app)
	if appErr != nil {
		return nil, appErr
	}
	// Status lookup service
	c, appErr := NewStatusService(app)
	if appErr != nil {
		return nil, appErr
	}

	// Lookup status condition
	s, appErr := NewStatusConditionService(app)
	if appErr != nil {
		return nil, appErr
	}

	// Close reason lookup service
	n, appErr := NewCloseReasonService(app)
	if appErr != nil {
		return nil, appErr
	}

	// Reason service
	r, appErr := NewReasonService(app)
	if appErr != nil {
		return nil, appErr
	}

	// * register appeal service
	grpcservice.RegisterAppealsServer(grpcServer, l)
	// * register status service
	grpcservice.RegisterStatusesServer(grpcServer, c)
	// * register lookup status service
	grpcservice.RegisterStatusConditionsServer(grpcServer, s)
	// * register close reason service
	grpcservice.RegisterCloseReasonsServer(grpcServer, n)
	// * register reason service
	grpcservice.RegisterReasonsServer(grpcServer, r)

	return grpcServer, nil
}

func unaryInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	var reqCtx context.Context
	var ip string

	// Extract metadata from incoming context
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		reqCtx = context.WithValue(ctx, RequestContextName, md)
		ip = getClientIp(md)
	} else {
		ip = "<not found>"
		reqCtx = context.WithValue(ctx, RequestContextName, nil)
	}

	// Log the start of the request for tracing
	slog.Info("cases.grpc_server.request_started",
		slog.String("method", info.FullMethod),
		slog.Time("start_time", start),
	)

	// Handle the request
	h, err := handler(reqCtx, req)

	// Log the result
	if err != nil {
		slog.Error("cases.grpc_server.request_error",
			slog.String("ip", ip),
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)),
			slog.String("error", err.Error()))
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
		slog.Info("cases.grpc_server.request_success",
			slog.String("method", info.FullMethod),
			slog.Duration("duration", time.Since(start)))
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
