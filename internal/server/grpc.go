package server

import (
	"fmt"
	"net"

	"github.com/webitel/cases/auth/user_auth"

	"github.com/bufbuild/protovalidate-go"
	conf "github.com/webitel/cases/config"
	grpcerr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/server/interceptor"
	"github.com/webitel/cases/registry"
	"github.com/webitel/cases/registry/consul"
	otelgrpc "github.com/webitel/webitel-go-kit/tracing/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Server   *grpc.Server
	listener net.Listener
	config   *conf.ConsulConfig
	exitChan chan error
	registry registry.ServiceRegistrator
}

// BuildServer constructs and configures a new gRPC server with interceptors.
func BuildServer(config *conf.ConsulConfig, authManager user_auth.AuthManager, exitChan chan error) (*Server, error) {
	// Initialize protovalidate validator
	val, err := protovalidate.New(protovalidate.WithFailFast(true))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize protovalidate: %w", err)
	}

	// Create a new gRPC server with interceptors and tracing
	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithMessageEvents(otelgrpc.SentEvents, otelgrpc.ReceivedEvents),
		)),
		grpc.ChainUnaryInterceptor(
			interceptor.OuterInterceptor(),
			interceptor.AuthUnaryServerInterceptor(authManager),
			interceptor.ValidateUnaryServerInterceptor(val),
		),
	)

	// Open a TCP listener on the configured address
	listener, err := net.Listen("tcp", config.PublicAddress)
	if err != nil {
		return nil, grpcerr.NewInternalError("server.build.listen.error", err.Error())
	}

	// Initialize Consul service registry
	reg, err := consul.NewConsulRegistry(config)
	if err != nil {
		return nil, grpcerr.NewInternalError("server.build.consul_registry.error", err.Error())
	}

	// Register gRPC reflection for debugging
	reflection.Register(s)

	return &Server{
		Server:   s,
		listener: listener,
		exitChan: exitChan,
		config:   config,
		registry: reg,
	}, nil
}

// Start registers and starts the gRPC server
func (s *Server) Start() {
	if err := s.registry.Register(); err != nil {
		s.exitChan <- err
		return
	}
	if err := s.Server.Serve(s.listener); err != nil {
		s.exitChan <- grpcerr.NewInternalError("server.start.serve.error", err.Error())
	}
}

// Stop deregisters the service and gracefully stops the gRPC server
func (s *Server) Stop() {
	if err := s.registry.Deregister(); err != nil {
		s.exitChan <- err
		return
	}
	s.Server.Stop()
}
