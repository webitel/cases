package webitel_manager

import (
	"context"

	"github.com/webitel/cases/auth/user_auth"
	autherror "github.com/webitel/cases/internal/error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ user_auth.AuthManager = &WebitelAppAuthManager{}

type WebitelAppAuthManager struct {
	client *user_auth.AuthorizationClient
}

func NewWebitelAppAuthManager(conn *grpc.ClientConn) (user_auth.AuthManager, error) {
	cli, err := user_auth.NewAuthorizationClient(conn)
	if err != nil {
		return nil, err
	}
	manager := &WebitelAppAuthManager{client: cli}

	return manager, nil
}

func (i *WebitelAppAuthManager) AuthorizeFromContext(ctx context.Context) (*user_auth.UserAuthSession, error) {
	var token []string
	var info metadata.MD
	var ok bool

	v := ctx.Value(user_auth.RequestContextName)
	info, ok = v.(metadata.MD)

	if !ok {
		info, ok = metadata.FromIncomingContext(ctx)
	}

	if !ok {
		return nil, autherror.NewForbiddenError("internal.grpc.get_context", "Not found")
	} else {
		token = info.Get(user_auth.AuthTokenName)
	}
	newContext := metadata.NewOutgoingContext(ctx, info)
	if len(token) < 1 {
		return nil, autherror.NewInternalError("webitel_manager.authorize_from_from_context.search_token.not_found", "token not found")
	}
	return i.Authorize(newContext, token[0])
}

func (i *WebitelAppAuthManager) Authorize(ctx context.Context, token string) (*user_auth.UserAuthSession, error) {
	return i.client.UserInfo(ctx, token)
}
