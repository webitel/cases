package model

import (
	authmodel "buf.build/gen/go/webitel/webitel-go/protocolbuffers/go"
	"strings"
	"time"
)

type Session struct {
	user        *AuthorizedUser
	permissions []*Permission
	scope       []*Scope
	roles       []*Role
	domainId    int64
	expiresAt   int64
}

func (s *Session) HasScope(scopeName string) bool {
	for _, scope := range s.scope {
		if scope.Name == scopeName {
			return true
		}
	}
	return false
}

func (s *Session) GetScope(scopeName string) *Scope {
	for _, scope := range s.scope {
		if scope.Class == scopeName {
			return scope
		}
	}
	return nil
}

func (s *Session) GetUserId() int64 {
	if s.user == nil {
		return 0
	}
	return s.user.Id
}

func (s *Session) GetUser() *AuthorizedUser {
	if s.user == nil {
		return nil
	}
	clone := *s.user
	return &clone
}

func (s *Session) GetDomainId() int64 {
	return s.domainId
}

func (s *Session) GetAclRoles() []int {
	var roles []int
	for _, role := range s.roles {
		roles = append(
			roles,
			int(role.Id),
		)
	}
	return roles
}

func (s *Session) IsExpired() bool {
	return time.Now().Unix() > s.expiresAt
}

func (s *Session) HasPermission(permissionName string) bool {
	for _, permission := range s.permissions {
		if permission.Name == permissionName {
			return true
		}
	}
	return false
}

func (s *Session) HasAccess(scope *Scope, accessType AccessMode) bool {
	if scope == nil {
		return false
	}

	var (
		bypass, require string
	)

	switch accessType {
	case Delete, Read | Delete:
		require, bypass = "d", "delete"
	case Edit, Read | Edit:
		require, bypass = "w", "write"
	case Read, NONE:
		require, bypass = "r", "read"
	case Add, Read | Add:
		require, bypass = "x", "add"
	}
	if bypass != "" && s.HasPermission(bypass) {
		return true
	}
	for i := len(require) - 1; i >= 0; i-- {
		mode := require[i]
		if strings.IndexByte(scope.Access, mode) < 0 {
			return false
		}
	}

	return true
}

func ConstructSessionFromUserInfo(userinfo *authmodel.Userinfo) *Session {
	session := &Session{
		user: &AuthorizedUser{
			Id:        userinfo.UserId,
			Name:      userinfo.Name,
			Username:  userinfo.Username,
			Extension: userinfo.Extension,
		},
		expiresAt: userinfo.ExpiresAt,
		domainId:  userinfo.Dc,
	}
	for i, permission := range userinfo.Permissions {
		if i == 0 {
			session.permissions = make([]*Permission, 0)
		}
		session.permissions = append(session.permissions, &Permission{
			Id:   permission.GetId(),
			Name: permission.GetName(),
		})
	}
	for i, scope := range userinfo.Scope {
		if i == 0 {
			session.scope = make([]*Scope, 0)
		}
		session.scope = append(session.scope, &Scope{
			Id:     scope.GetId(),
			Name:   scope.GetName(),
			Abac:   scope.Abac,
			Obac:   scope.Obac,
			Rbac:   scope.Rbac,
			Class:  scope.Class,
			Access: scope.Access,
		})
	}

	for i, role := range userinfo.Roles {
		if i == 0 {
			session.roles = make([]*Role, 0)
		}
		session.roles = append(session.roles, &Role{
			Id:   role.GetId(),
			Name: role.GetName(),
		})
	}
	return session
}
