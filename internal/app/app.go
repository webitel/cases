package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/webitel/cases/auth"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/auth/webitel_manager"
	conf "github.com/webitel/cases/config"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/server"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres"
	broker "github.com/webitel/cases/rabbit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	config         *conf.AppConfig
	Store          store.Store
	server         *server.Server
	exitChan       chan error
	storageConn    *grpc.ClientConn
	sessionManager auth.AuthManager
	webitelAppConn *grpc.ClientConn
	shutdown       func(ctx context.Context) error
	log            *slog.Logger
	rabbit         *broker.RabbitBroker
	rabbitExitChan chan cerror.AppError
}

func New(config *conf.AppConfig, shutdown func(ctx context.Context) error) (*App, error) {
	// --------- App Initialization ---------
	app := &App{config: config, shutdown: shutdown}
	var err error

	// --------- DB Initialization ---------
	if config.Database == nil {
		return nil, cerror.NewBadRequestError("internal.internal.new.database_config.bad_arguments", "error creating store, config is nil")
	}
	app.Store = BuildDatabase(config.Database)

	// --------- Message Broker ( Rabbit ) Initialization ---------

	r, appErr := broker.BuildRabbit(app.config.Rabbit, app.rabbitExitChan)
	if appErr != nil {
		return nil, appErr
	}
	app.rabbit = r

	// --------- Webitel App gRPC Connection ( Consul )---------
	app.webitelAppConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, cerror.NewInternalError("internal.internal.new_app.grpc_conn.error", err.Error())
	}

	// --------- Session Manager Initialization ---------
	app.sessionManager, err = webitel_manager.NewWebitelAppAuthManager(app.webitelAppConn)
	if err != nil {
		return nil, err
	}

	// --------- gRPC Server Initialization ---------
	s, err := server.BuildServer(app.config.Consul, app.sessionManager, app.exitChan)
	if err != nil {
		return nil, err
	}
	app.server = s

	// --------- Service Registration ---------
	RegisterServices(app.server.Server, app)

	// --------- Storage gRPC Connection ---------
	app.storageConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/store?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, cerror.NewInternalError("internal.internal.new_app.grpc_conn.error", err.Error())
	}

	return app, nil
}

func BuildDatabase(config *conf.DatabaseConfig) store.Store {
	return postgres.New(config)
}

func (a *App) Start() error { // Change return type to standard error
	err := a.Store.Open()
	if err != nil {
		return err
	}

	// * run grpc server
	go a.server.Start()
	return <-a.exitChan
}

func (a *App) Stop() error { // Change return type to standard error
	// close massive modules
	a.server.Stop()
	// close store connection
	a.Store.Close()
	// close grpc connections
	a.storageConn.Close()
	a.webitelAppConn.Close()

	// ----- Call the shutdown function for OTel ----- //
	if a.shutdown != nil {
		a.shutdown(context.Background())
	}

	return nil
}

func (a *App) AuthorizeFromContext(ctx context.Context) (*authmodel.Session, error) { // Change return type to standard error
	session, err := a.sessionManager.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, cerror.NewUnauthorizedError("internal.internal.authorize_from_context.validate_session.expired", "session expired")
	}
	return session, nil
}
