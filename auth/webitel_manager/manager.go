package webitel_manager

import (
	"context"

	iface "github.com/webitel/cases/auth"
	"github.com/webitel/cases/auth/model"
	autherror "github.com/webitel/cases/internal/error"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var _ iface.AuthManager = &WebitelAppAuthManager{}

type WebitelAppAuthManager struct {
	client *iface.AuthorizationClient
}

func NewWebitelAppAuthManager(conn *grpc.ClientConn) (iface.AuthManager, error) {
	cli, err := iface.NewAuthorizationClient(conn)
	if err != nil {
		return nil, err
	}
	manager := &WebitelAppAuthManager{client: cli}

	return manager, nil
}

func (i *WebitelAppAuthManager) AuthorizeFromContext(ctx context.Context) (*model.Session, error) {
	var token []string
	var info metadata.MD
	var ok bool

	v := ctx.Value(model.RequestContextName)
	info, ok = v.(metadata.MD)

	if !ok {
		info, ok = metadata.FromIncomingContext(ctx)
	}

	if !ok {
		return nil, autherror.NewForbiddenError("internal.grpc.get_context", "Not found")
	} else {
		token = info.Get(model.AuthTokenName)
	}
	newContext := metadata.NewOutgoingContext(ctx, info)
	if len(token) < 1 {
		return nil, autherror.NewInternalError("webitel_manager.authorize_from_from_context.search_token.not_found", "token not found")
	}
	return i.Authorize(newContext, token[0])
}

func (i *WebitelAppAuthManager) Authorize(ctx context.Context, token string) (*model.Session, error) {
	return i.client.UserInfo(ctx, token)
}
