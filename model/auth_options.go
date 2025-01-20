package model

import (
	"context"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/server/interceptor"
)

type Auther interface {
	GetRoles() []int64
	GetUserId() int64
	GetDomainId() int64
	GetObjectScope(string) ObjectScope
}

type ObjectScope interface {
	IsRbacUsed() bool
	IsObacUsed() bool
	GetObjectName() string
}

func NewSessionAuthOptions(session *authmodel.Session, requiredScopes ...string) Auther {
	if len(requiredScopes) == 0 || session == nil {
		return nil
	}
	scopes := make(map[string]ObjectScope)
	for _, s := range requiredScopes {
		scopes[s] = newSessionObjectScope(session.GetScope(s))
	}
	return &SessionAuthOptions{userId: session.GetUserId(), roles: session.GetAclRoles(), domainId: session.GetDomainId(), scopes: scopes}
}

type SessionAuthOptions struct {
	roles    []int64
	userId   int64
	domainId int64
	scopes   map[string]ObjectScope
}

func (a *SessionAuthOptions) GetRoles() []int64 {
	return a.roles
}
func (a *SessionAuthOptions) GetUserId() int64 {
	return a.userId
}
func (a *SessionAuthOptions) GetDomainId() int64 {
	return a.domainId
}
func (a *SessionAuthOptions) GetObjectScope(s string) ObjectScope {
	return a.scopes[s]
}

type DefaultScope struct {
	rbac    bool
	obac    bool
	objName string
}

func newSessionObjectScope(scope *authmodel.Scope) ObjectScope {
	if scope == nil {
		return nil
	}
	return &DefaultScope{
		rbac:    scope.IsRbacUsed(),
		obac:    scope.IsObacUsed(),
		objName: scope.Name,
	}
}

func (d *DefaultScope) IsRbacUsed() bool {
	if d == nil {
		return false
	}
	return d.rbac
}

func (d *DefaultScope) IsObacUsed() bool {
	if d == nil {
		return false
	}
	return d.obac
}

func (d *DefaultScope) GetObjectName() string {
	return d.objName
}

func GetSessionOutOfContext(ctx context.Context) *authmodel.Session {
	return ctx.Value(interceptor.SessionHeader).(*authmodel.Session)
}
