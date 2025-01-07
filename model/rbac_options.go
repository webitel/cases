package model

import authmodel "github.com/webitel/cases/auth/model"

type Auther interface {
	GetRoles() []int64
	IsRbacEnabled() bool
	GetObjectName() string
	GetUserId() int64
	GetDomainId() int64
}

type DefaultAuthOptions struct {
	roles    []int64
	userId   int64
	domainId int64
	scope    *authmodel.Scope
}

func NewDefaultAuthOptions(session *authmodel.Session, objClassName string) Auther {
	if objClassName == "" {
		return nil
	}
	if session == nil {
		return nil
	}
	return &DefaultAuthOptions{userId: session.GetUserId(), roles: session.GetAclRoles(), domainId: session.GetDomainId(), scope: session.GetScope(objClassName)}
}

func (a *DefaultAuthOptions) GetRoles() []int64 {
	return a.roles
}

func (a *DefaultAuthOptions) IsRbacEnabled() bool {
	return a.scope.IsRbacUsed()
}

func (a *DefaultAuthOptions) GetObjectName() string {
	return a.scope.Name
}
func (a *DefaultAuthOptions) GetUserId() int64 {
	return a.userId
}
func (a *DefaultAuthOptions) GetDomainId() int64 {
	return a.domainId
}
