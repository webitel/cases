package app

import (
	"net"
	"time"

	grpcservice "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/registry"
	"github.com/webitel/cases/registry/consul"
	otelgrpc "github.com/webitel/webitel-go-kit/tracing/grpc"
	"google.golang.org/grpc"
)

var (
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
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithMessageEvents(otelgrpc.SentEvents, otelgrpc.ReceivedEvents),
		)),
		grpc.UnaryInterceptor(unaryInterceptor),
	)

	// * Creating services
	// Source Lookup service
	l, appErr := NewSourceService(app)
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

	// Priority service
	p, appErr := NewPriorityService(app)
	if appErr != nil {
		return nil, appErr
	}

	// SLA service
	sla, appErr := NewSLAService(app)
	if appErr != nil {
		return nil, appErr
	}

	// SLA condition service
	slaC, appErr := NewSLAConditionService(app)
	if appErr != nil {
		return nil, appErr
	}

	// Catalog service
	catalog, appErr := NewCatalogService(app)
	if appErr != nil {
		return nil, appErr
	}

	// Service service
	service, appErr := NewServiceService(app)
	if appErr != nil {
		return nil, appErr
	}

	// * Register the services
	grpcservice.RegisterSourcesServer(grpcServer, l)
	grpcservice.RegisterStatusesServer(grpcServer, c)
	grpcservice.RegisterStatusConditionsServer(grpcServer, s)
	grpcservice.RegisterCloseReasonsServer(grpcServer, n)
	grpcservice.RegisterReasonsServer(grpcServer, r)
	grpcservice.RegisterPrioritiesServer(grpcServer, p)
	grpcservice.RegisterSLAsServer(grpcServer, sla)
	grpcservice.RegisterSLAConditionsServer(grpcServer, slaC)
	grpcservice.RegisterCatalogsServer(grpcServer, catalog)
	grpcservice.RegisterServicesServer(grpcServer, service)

	return grpcServer, nil
}
