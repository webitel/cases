package app

import (
	"context"
	"fmt"
	"github.com/webitel/cases/auth"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/auth/webitel_manager"
	"github.com/webitel/cases/internal/db"
	"github.com/webitel/cases/internal/db/postgres"
	server2 "github.com/webitel/cases/internal/server"
	"github.com/webitel/cases/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	config         *model.AppConfig
	DB             db.DB
	server         *server2.Server
	exitChan       chan model.AppError
	storageConn    *grpc.ClientConn
	sessionManager auth.AuthManager
	webitelAppConn *grpc.ClientConn
}

func New(config *model.AppConfig) (*App, model.AppError) {
	app := &App{config: config}
	var err error

	// init of a database
	if config.Database == nil {
		model.NewInternalError("pkg.pkg.new.database_config.bad_arguments", "error creating db, config is nil")
	}
	app.DB = BuildDatabase(config.Database)

	// init of grpc server
	s, appErr := server2.BuildServer(app, app.config.Consul, app.exitChan)
	if appErr != nil {
		return nil, appErr
	}
	app.server = s

	// init service connections
	app.storageConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/db?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, model.NewInternalError("pkg.pkg.new_app.grpc_conn.error", err.Error())
	}

	app.webitelAppConn, err = grpc.NewClient(fmt.Sprintf("consul://%s/go.webitel.pkg?wait=14s", config.Consul.Address),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, model.NewInternalError("pkg.pkg.new_app.grpc_conn.error", err.Error())
	}

	app.sessionManager, appErr = webitel_manager.NewWebitelAppAuthManager(app.webitelAppConn)
	if appErr != nil {
		return nil, appErr
	}

	return app, nil
}

func BuildDatabase(config *model.DatabaseConfig) db.DB {
	return postgres.New(config)
}

func (a *App) Start() model.AppError {

	err := a.DB.Open()
	if err != nil {
		return err
	}

	// * run grpc server
	go a.server.Start()
	//go ServeRequests(a, a.config.Consul, a.exitChan)
	return <-a.exitChan
}

func (a *App) Stop() model.AppError {
	// close massive modules
	a.server.Stop()
	// close db connection
	a.DB.Close()
	// close grpc connections
	a.storageConn.Close()
	// close webitel pkg connection
	a.webitelAppConn.Close()

	return nil
}

func (a *App) AuthorizeFromContext(ctx context.Context) (*authmodel.Session, model.AppError) {
	session, err := a.sessionManager.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	if session.IsExpired() {
		return nil, model.NewUnauthorizedError("pkg.pkg.authorize_from_context.validate_session.expired", "session expired")
	}
	return session, nil
}

func (a *App) MakePermissionError(session *authmodel.Session) model.AppError {
	if session == nil {
		return model.NewForbiddenError("pkg.permissions.check_access.denied", "access denied")
	}
	return model.NewForbiddenError("pkg.permissions.check_access.denied", fmt.Sprintf("userId=%d, access denied", session.GetUserId()))
}

func (a *App) MakeScopeError(session *authmodel.Session, scope *authmodel.Scope, access authmodel.AccessMode) model.AppError {
	if session == nil || session.GetUser() == nil || scope == nil {
		return model.NewForbiddenError("pkg.scope.check_access.denied", fmt.Sprintf("access denied"))
	}
	return model.NewForbiddenError("pkg.scope.check_access.denied", fmt.Sprintf("access denied scope=%s access=%d for user %d", scope.Name, access, session.GetUserId()))
}
