package app

import (
	"context"
	"fmt"

	"github.com/webitel/cases/auth"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/auth/webitel_manager"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres"
	"github.com/webitel/cases/model"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	config         *model.AppConfig
	Store          store.Store
	server         *Server
	exitChan       chan model.AppError
	storageConn    *grpc.ClientConn
	sessionManager auth.AuthManager
	webitelAppConn *grpc.ClientConn
	shutdown       func(ctx context.Context) error
}

func New(config *model.AppConfig, shutdown func(ctx context.Context) error) (*App, model.AppError) {
	app := &App{config: config, shutdown: shutdown}
	var err error

	// init of a database
	if config.Database == nil {
		model.NewInternalError("internal.internal.new.database_config.bad_arguments", "error creating store, config is nil")
	}
	app.Store = BuildDatabase(config.Database)

	// init of grpc server
	s, appErr := BuildServer(app, app.config.Consul, app.exitChan)
	if appErr != nil {
		return nil, appErr
	}
	app.server = s

	// init service connections
	app.storageConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/store?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, model.NewInternalError("internal.internal.new_app.grpc_conn.error", err.Error())
	}

	app.webitelAppConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/go.webitel.app?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, model.NewInternalError("internal.internal.new_app.grpc_conn.error", err.Error())
	}

	app.sessionManager, appErr = webitel_manager.NewWebitelAppAuthManager(app.webitelAppConn)
	if appErr != nil {
		return nil, appErr
	}

	return app, nil
}

func BuildDatabase(config *model.DatabaseConfig) store.Store {
	return postgres.New(config)
}

func (a *App) Start() model.AppError {
	err := a.Store.Open()
	if err != nil {
		return err
	}

	// * run grpc server
	go a.server.Start()
	return <-a.exitChan
}

func (a *App) Stop() model.AppError {
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

func (a *App) AuthorizeFromContext(ctx context.Context) (*authmodel.Session, model.AppError) {
	session, err := a.sessionManager.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, model.NewUnauthorizedError("internal.internal.authorize_from_context.validate_session.expired", "session expired")
	}
	return session, nil
}

func (a *App) MakePermissionError(session *authmodel.Session) model.AppError {
	if session == nil {
		return model.NewForbiddenError("internal.permissions.check_access.denied", "access denied")
	}
	return model.NewForbiddenError("internal.permissions.check_access.denied", fmt.Sprintf("userId=%d, access denied", session.GetUserId()))
}

func (a *App) MakeScopeError(session *authmodel.Session, scope *authmodel.Scope, access authmodel.AccessMode) model.AppError {
	if session == nil || session.GetUser() == nil || scope == nil {
		return model.NewForbiddenError("internal.scope.check_access.denied", "access denied")
	}
	return model.NewForbiddenError("internal.scope.check_access.denied", fmt.Sprintf("access denied scope=%s access=%d for user %d", scope.Name, access, session.GetUserId()))
}
