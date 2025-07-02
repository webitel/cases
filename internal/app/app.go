package app

import (
	"context"
	"fmt"
	webitelgo "github.com/webitel/cases/api/webitel-go/contacts"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/auth/manager/webitel_app"
	conf "github.com/webitel/cases/config"
	ftsadapter "github.com/webitel/cases/internal/adapters/fts"
	loggeradapter "github.com/webitel/cases/internal/adapters/logger"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/server"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres"
	ftsclient "github.com/webitel/webitel-go-kit/infra/fts_client"
	rabbit "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq"
	brokeradapter "github.com/webitel/webitel-go-kit/infra/pubsub/rabbitmq/pkg/adapter/slog"
	"github.com/webitel/webitel-go-kit/pkg/watcher"
	"log/slog"

	"github.com/webitel/cases/api/engine"
	wlogger "github.com/webitel/webitel-go-kit/infra/logger_client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	AnonymousName = "Anonymous"
)

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
	rabbitConn          *rabbit.Connection
	rabbitPublisher     rabbit.Publisher
	webitelgoClient     webitelgo.GroupsClient
	engineConn          *grpc.ClientConn
	engineAgentClient   engine.AgentServiceClient
	wtelLogger          *wlogger.Logger
	ftsClient           *ftsclient.Client
	watcherManager      watcher.Manager
	caseResolutionTimer *TimerTask[*App]
}

func StartBroker(config *conf.AppConfig) (*rabbit.Connection, error) {
	if config.Rabbit == nil {
		return nil, errors.New("error creating broker, config is nil")
	}

	conf, err := rabbit.NewConfig(config.Rabbit.Url)
	if err != nil {
		return nil, errors.New("error creating broker config", errors.WithCause(err))
	}

	r, err := rabbit.NewConnection(conf, brokeradapter.NewSlogLogger(slog.Default()))
	if err != nil {
		return nil, errors.New("error creating rabbit connection", errors.WithCause(err))
	}
	exchangeConf, err := rabbit.NewExchangeConfig("cases", "topic")
	if err != nil {
		return nil, errors.New("error creating exchange config", errors.WithCause(err))
	}

	err = r.DeclareExchange(context.Background(), exchangeConf)
	if err != nil {
		return nil, errors.New("error declaring exchange", errors.WithCause(err))
	}

	return r, nil
}

func New(config *conf.AppConfig, shutdown func(ctx context.Context) error) (*App, error) {
	// --------- App Initialization ---------
	app := &App{config: config, shutdown: shutdown}
	var err error

	// --------- DB Initialization ---------
	if config.Database == nil {
		return nil, errors.New("error creating store, config is nil")
	}
	app.Store = BuildDatabase(config.Database)

	// --------- Message Broker ( Rabbit ) Initialization ---------
	app.rabbitConn, err = StartBroker(config)
	if err != nil {
		return nil, err
	}
	publisherConf, err := rabbit.NewPublisherConfig()
	if err != nil {
		return nil, errors.New("error creating publisher config", errors.WithCause(err))
	}
	app.rabbitPublisher, err = rabbit.NewPublisher(app.rabbitConn, publisherConf, brokeradapter.NewSlogLogger(slog.Default()))
	if err != nil {
		return nil, err
	}

	// register watchers
	watcherManager := watcher.NewDefaultWatcherManager(config.WatchersEnabled)
	app.watcherManager = watcherManager
	//

	// --------- Webitel App gRPC Connection ---------
	app.webitelAppConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	app.webitelgoClient = webitelgo.NewGroupsClient(app.webitelAppConn)

	if err != nil {
		return nil, errors.New("unable to create contact group client", errors.WithCause(err))
	}

	// --------- Webitel Engine gRPC Connection ---------
	app.engineConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/engine?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	app.engineAgentClient = engine.NewAgentServiceClient(app.engineConn)

	if err != nil {
		return nil, errors.New("unable to create agent client", errors.WithCause(err))
	}

	// --------- Webitel Logger gRPC Connection ---------
	loggerAdapter, err := loggeradapter.New(app.rabbitPublisher)
	if err != nil {
		return nil, err
	}
	app.wtelLogger, err = wlogger.New(loggerAdapter)
	if err != nil {
		return nil, errors.New("unable to create logger client", errors.WithCause(err))
	}

	// --------- Session Manager Initialization ---------
	app.sessionManager, err = webitel_app.New(app.webitelAppConn)
	if err != nil {
		return nil, err
	}

	// --------- Full Text Search Client ---------
	ftsAdapter, err := ftsadapter.NewDefaultClient(app.rabbitPublisher)
	if err != nil {
		return nil, err
	}
	app.ftsClient = ftsclient.New(ftsAdapter)

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
		return nil, errors.New("unable to create storage client", errors.WithCause(err))
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
