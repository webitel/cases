package auth

import (
	authclient "buf.build/gen/go/webitel/webitel-go/grpc/go/_gogrpc"
	authmodel "buf.build/gen/go/webitel/webitel-go/protocolbuffers/go"
	"context"
	"github.com/golang/groupcache/singleflight"
	"github.com/webitel/logger/auth/model"
	errors "github.com/webitel/logger/model"
	"google.golang.org/grpc"
)

type AuthManager interface {
	Authorize(ctx context.Context, token string) (*model.Session, errors.AppError)
	AuthorizeFromContext(ctx context.Context) (*model.Session, errors.AppError)
}

type AuthorizationClient struct {
	Client     authclient.AuthClient
	Group      singleflight.Group
	Connection *grpc.ClientConn
}

func NewAuthorizationClient(conn *grpc.ClientConn) (*AuthorizationClient, errors.AppError) {
	if conn == nil {
		return nil, errors.NewInternalError("auth.manager.new_auth_client.validate_params.connection", "invalid GRPC connection")
	}
	return &AuthorizationClient{
		Client:     authclient.NewAuthClient(conn),
		Group:      singleflight.Group{},
		Connection: conn,
	}, nil
}

func (c *AuthorizationClient) UserInfo(ctx context.Context, token string) (*model.Session, errors.AppError) {
	interfacedSession, err := c.Group.Do(token, func() (interface{}, error) {
		info, err := c.Client.UserInfo(ctx, &authmodel.UserinfoRequest{AccessToken: token})
		if err != nil {
			return nil, err
		}
		return model.ConstructSessionFromUserInfo(info), nil
	})
	if err != nil {
		return nil, errors.NewUnauthorizedError("auth.manager.user_info.do_request.error", err.Error())
	}
	return interfacedSession.(*model.Session), nil
}
