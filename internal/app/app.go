package app

import (
	"context"
	"fmt"
	webitelgo "github.com/webitel/cases/api/webitel-go/contacts"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/auth/manager/webitel_app"
	conf "github.com/webitel/cases/config"
	ftspublisher "github.com/webitel/cases/fts_client"
	"github.com/webitel/cases/internal/errors"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/server"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres"
	broker "github.com/webitel/cases/rabbit"
	ftsclient "github.com/webitel/webitel-go-kit/fts_client"
	"log/slog"

	engine "github.com/webitel/cases/api/engine"
	wlogger "github.com/webitel/logger/pkg/client/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AnonymousName = "Anonymous"
)

func NewBadRequestError(err error) errors.AppError {
	return errors.NewBadRequestError("app.process_api.validation.error", err.Error())
}

type App struct {
	config              *conf.AppConfig
	Store               store.Store
	server              *server.Server
	exitChan            chan error
	storageConn         *grpc.ClientConn
	sessionManager      auth.Manager
	webitelAppConn      *grpc.ClientConn
	shutdown            func(ctx context.Context) error
	log                 *slog.Logger
	rabbit              *broker.RabbitBroker
	rabbitExitChan      chan cerror.AppError
	webitelgoClient     webitelgo.GroupsClient
	engineConn          *grpc.ClientConn
	engineAgentClient   engine.AgentServiceClient
	wtelLogger          *wlogger.LoggerClient
	ftsClient           *ftsclient.Client
	watcherManager      WatcherManager
	caseResolutionTimer *TimerTask[*App]
}

func New(config *conf.AppConfig, shutdown func(ctx context.Context) error) (*App, error) {
	// --------- App Initialization ---------
	app := &App{config: config, shutdown: shutdown}
	var err error

	// --------- DB Initialization ---------
	if config.Database == nil {
		return nil, cerror.NewInternalError("internal.internal.new.database_config.bad_arguments", "error creating store, config is nil")
	}
	app.Store = BuildDatabase(config.Database)

	// --------- Message Broker ( Rabbit ) Initialization ---------

	r, appErr := broker.BuildRabbit(app.config.Rabbit, app.rabbitExitChan)
	if appErr != nil {
		return nil, appErr
	}
	app.rabbit = r

	// Start the Rabbit connection and consumers
	appErr = app.rabbit.Start()
	if appErr != nil {
		return nil, cerror.NewInternalError("internal.internal.new_app.rabbit.start.error", appErr.Error())
	}

	// register watchers
	watcherManager := NewDefaultWatcherManager(config.WatchersEnabled)
	app.watcherManager = watcherManager
	//

	// --------- Webitel App gRPC Connection ---------
	app.webitelAppConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	app.webitelgoClient = webitelgo.NewGroupsClient(app.webitelAppConn)

	if err != nil {
		return nil, cerror.NewInternalError("internal.internal.new_app.grpc_conn.error", err.Error())
	}

	// --------- Webitel Engine gRPC Connection ---------
	app.engineConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/engine?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	app.engineAgentClient = engine.NewAgentServiceClient(app.engineConn)

	if err != nil {
		return nil, cerror.NewInternalError("internal.internal.new_engine.grpc_conn.error", err.Error())
	}

	// --------- Webitel Logger gRPC Connection ---------
	app.wtelLogger, err = wlogger.NewLoggerClient(wlogger.WithAmqpConnectionString(app.config.Rabbit.Url), wlogger.WithGrpcConsulAddress(config.Consul.Address))
	if err != nil {
		return nil, cerror.NewInternalError("internal.internal.new_app.grpc_conn.error", err.Error())
	}

	// --------- Session Manager Initialization ---------
	app.sessionManager, err = webitel_app.New(app.webitelAppConn)
	if err != nil {
		return nil, err
	}

	// --------- Full Text Search Client ---------
	app.ftsClient, err = ftspublisher.NewDefaultClient(app.rabbit)
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

	a.initCustom()

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
	err := a.storageConn.Close()
	if err != nil {
		return err
	}
	err = a.webitelAppConn.Close()
	if err != nil {
		return err
	}

	if a.caseResolutionTimer != nil {
		a.caseResolutionTimer.Stop()
	}

	// ----- Call the shutdown function for OTel ----- //
	if a.shutdown != nil {
		err := a.shutdown(context.Background())
		if err != nil {
			return err
		}
	}

	return nil
}
