package user_auth

import (
	"github.com/webitel/cases/auth"
	"strings"
	"time"

	authmodel "buf.build/gen/go/webitel/webitel-go/protocolbuffers/go"
)

type UserAuthSession struct {
	user        *User
	permissions []*Permission
	scope       []*Scope
	ifScopes    map[string]auth.ObjectScoper
	roles       []*Role
	domainId    int64
	expiresAt   int64
}

// region Auther interface implementation

func (s *UserAuthSession) GetUserId() int64 {
	if s.user == nil {
		return 0
	}
	return s.user.Id
}

func (s *UserAuthSession) GetDomainId() int64 {
	return s.domainId
}

func (s *UserAuthSession) GetRoles() []int64 {
	roles := []int64{s.GetUserId()}
	for _, role := range s.roles {
		roles = append(
			roles,
			role.Id,
		)
	}
	return roles
}

func (s *UserAuthSession) GetObjectScope(sc string) auth.ObjectScoper {
	return s.ifScopes[sc]
}

func (s *UserAuthSession) GetPermissions() []string {
	//TODO implement me
	panic("implement me")
}

// endregion

func (s *UserAuthSession) HasScope(scopeName string) bool {
	for _, scope := range s.scope {
		if scope.Name == scopeName {
			return true
		}
	}
	return false
}

func (s *UserAuthSession) GetScope(scopeName string) *Scope {
	for _, scope := range s.scope {
		if scope.Class == scopeName {
			return scope
		}
	}
	return nil
}

func (s *UserAuthSession) GetUserName() string {
	if s.user == nil {
		return ""
	}
	return s.user.Name
}

func (s *UserAuthSession) GetUser() *User {
	if s.user == nil {
		return nil
	}
	clone := *s.user
	return &clone
}

func (s *UserAuthSession) IsExpired() bool {
	return time.Now().Unix() > s.expiresAt
}

func (s *UserAuthSession) HasPermission(permissionName string) bool {
	for _, permission := range s.permissions {
		if permission.Id == permissionName {
			return true
		}
	}
	return false
}

func (s *UserAuthSession) HasObacAccess(scopeName string, accessType auth.AccessMode) bool {
	scope := s.GetScope(scopeName)
	if scope == nil {
		return false
	}

	var bypass, require string

	switch accessType {
	case auth.Delete, auth.Read | auth.Delete:
		require, bypass = "d", "delete"
	case auth.Edit, auth.Read | auth.Edit:
		require, bypass = "w", "write"
	case auth.Read, auth.NONE:
		require, bypass = "r", "read"
	case auth.Add, auth.Read | auth.Add:
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

func ConstructSessionFromUserInfo(userinfo *authmodel.Userinfo) *UserAuthSession {
	session := &UserAuthSession{
		user: &User{
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
