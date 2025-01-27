package user_auth

import (
	"github.com/webitel/cases/auth"
	"strings"
	"time"

	authmodel "buf.build/gen/go/webitel/webitel-go/protocolbuffers/go"
)

type UserAuthSession struct {
	user             *User
	permissions      []string
	scopes           map[string]*Scope
	license          map[string]bool
	roles            []*Role
	domainId         int64
	expiresAt        int64
	superCreate      bool
	superEdit        bool
	superDelete      bool
	superSelect      bool
	mainAccess       auth.AccessMode
	mainObjClassName string
}

// region Auther interface implementation

func (s *UserAuthSession) GetUserId() int64 {
	if s.user == nil || s.user.Id <= 0 {
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
	if sc == "" {
		return nil
	}
	scope, found := s.scopes[sc]
	if !found {
		return nil
	}
	return scope
}

func (s *UserAuthSession) GetAllObjectScopes() []auth.ObjectScoper {
	var res []auth.ObjectScoper
	for _, scope := range s.scopes {
		res = append(res, scope)
	}
	return res
}

func (s *UserAuthSession) GetPermissions() []string {
	return s.permissions
}

func (s *UserAuthSession) CheckLicenseAccess(name string) bool {
	if legit, found := s.license[name]; found {
		return legit
	}
	return false
}

func (s *UserAuthSession) GetMainAccessMode() auth.AccessMode {
	return s.mainAccess
}

func (s *UserAuthSession) GetMainObjClassName() string {
	return s.mainObjClassName
}

func (s *UserAuthSession) CheckObacAccess(scopeName string, accessType auth.AccessMode) bool {
	scope := s.GetObjectScope(scopeName)
	if scope == nil {
		return false
	}

	if scope.IsObacUsed() {
		var (
			bypass  bool
			require string
		)

		switch accessType {
		case auth.Delete, auth.Read | auth.Delete:
			require, bypass = "d", s.superDelete
		case auth.Edit, auth.Read | auth.Edit:
			require, bypass = "w", s.superEdit
		case auth.Read, auth.NONE:
			require, bypass = "r", s.superSelect
		case auth.Add, auth.Read | auth.Add:
			require, bypass = "x", s.superCreate
		}
		if bypass {
			return true
		}
		for i := len(require) - 1; i >= 0; i-- {
			mode := require[i]
			if strings.IndexByte(scope.GetAccess(), mode) < 0 {
				return false
			}
		}
	}

	return true
}

func (s *UserAuthSession) IsRbacCheckRequired(scopeName string, accessType auth.AccessMode) bool {
	scope := s.GetObjectScope(scopeName)
	if scope == nil {
		return false
	}
	rbacEnabled := scope.IsRbacUsed()
	if rbacEnabled {
		var bypass bool

		switch accessType {
		case auth.Delete, auth.Read | auth.Delete:
			bypass = s.superDelete
		case auth.Edit, auth.Read | auth.Edit:
			bypass = s.superEdit
		case auth.Read, auth.NONE:
			bypass = s.superSelect
		case auth.Add, auth.Read | auth.Add:
			bypass = s.superCreate
		}
		if bypass {
			return false
		}
	}
	return rbacEnabled

}

// endregion

func (s *UserAuthSession) IsExpired() bool {
	return time.Now().Unix() > s.expiresAt
}

func ConstructSessionFromUserInfo(userinfo *authmodel.Userinfo, mainObjClass string, mainAccess auth.AccessMode) *UserAuthSession {
	session := &UserAuthSession{
		user: &User{
			Id:        userinfo.UserId,
			Name:      userinfo.Name,
			Username:  userinfo.Username,
			Extension: userinfo.Extension,
		},
		expiresAt:        userinfo.ExpiresAt,
		domainId:         userinfo.Dc,
		permissions:      make([]string, 0),
		license:          map[string]bool{},
		scopes:           map[string]*Scope{},
		mainAccess:       mainAccess,
		mainObjClassName: mainObjClass,
	}
	for _, lic := range userinfo.License {
		session.license[lic.Id] = lic.ExpiresAt > time.Now().UnixMilli()
	}
	for _, permission := range userinfo.Permissions {
		switch permission.GetId() {
		case auth.SuperCreatePermission:
			session.superCreate = true
		case auth.SuperDeletePermission:
			session.superDelete = true
		case auth.SuperEditPermission:
			session.superEdit = true
		case auth.SuperSelectPermission:
			session.superSelect = true
		}
		session.permissions = append(session.permissions, permission.GetId())
	}
	for _, scope := range userinfo.Scope {
		session.scopes[scope.Class] = &Scope{
			Id:     scope.GetId(),
			Name:   scope.GetName(),
			Abac:   scope.Abac,
			Obac:   scope.Obac,
			Rbac:   scope.Rbac,
			Class:  scope.Class,
			Access: scope.Access,
		}
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
