package auth

import (
	"context"

	authclient "buf.build/gen/go/webitel/webitel-go/grpc/go/_gogrpc"
	authmodel "buf.build/gen/go/webitel/webitel-go/protocolbuffers/go"
	"github.com/golang/groupcache/singleflight"

	// neccessary import for client setup ( if not imported - add [443:] to the end of the address )
	// if not impoerted cause such error:
	// ! failed to exit idle mode: invalid target address consul://10.9.8.111:8500/go.webitel.internal, error info: address consul://10.9.8.111:8500/go.webitel.internal:443: too many colons in address
	_ "github.com/mbobakov/grpc-consul-resolver"
	model "github.com/webitel/cases/auth/model"
	autherror "github.com/webitel/cases/internal/error"
	"google.golang.org/grpc"
)

type AuthManager interface {
	Authorize(ctx context.Context, token string) (*model.Session, error)
	AuthorizeFromContext(ctx context.Context) (*model.Session, error)
}

type AuthorizationClient struct {
	Client     authclient.AuthClient
	Group      singleflight.Group
	Connection *grpc.ClientConn
}

func NewAuthorizationClient(conn *grpc.ClientConn) (*AuthorizationClient, error) {
	if conn == nil {
		return nil, autherror.NewInternalError("auth.manager.new_auth_client.validate_params.connection", "invalid GRPC connection")
	}
	return &AuthorizationClient{
		Client:     authclient.NewAuthClient(conn),
		Group:      singleflight.Group{},
		Connection: conn,
	}, nil
}

func (c *AuthorizationClient) UserInfo(ctx context.Context, token string) (*model.Session, error) {
	interfacedSession, err := c.Group.Do(token, func() (interface{}, error) {
		info, err := c.Client.UserInfo(ctx, &authmodel.UserinfoRequest{AccessToken: token})
		if err != nil {
			return nil, err
		}
		return model.ConstructSessionFromUserInfo(info), nil
	})
	if err != nil {
		return nil, autherror.NewUnauthorizedError("auth.manager.user_info.do_request.error", err.Error())
	}
	return interfacedSession.(*model.Session), nil
}
