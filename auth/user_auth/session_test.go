package user_auth

import (
	"github.com/webitel/cases/auth"
	"reflect"
	"testing"
)

func TestUserAuthSession_CheckLicenseAccess(t *testing.T) {
	type fields struct {
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
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Active license",
			args: args{name: "CALL_CENTER"},
			fields: fields{
				license: map[string]bool{
					"CALL_CENTER":  true,
					"CALL_MANAGER": false,
					"CHATS":        true,
					"LICENSE1":     true,
					"LICENSE2":     true,
					"LICENSE3":     false,
					"LICENSE4":     true,
				},
			},
			want: true,
		},
		{
			name: "Expired license but present in map",
			args: args{name: "CALL_MANAGER"},
			fields: fields{
				license: map[string]bool{
					"CALL_CENTER":  true,
					"CALL_MANAGER": false,
					"CHATS":        true,
					"LICENSE1":     true,
					"LICENSE2":     true,
					"LICENSE3":     false,
					"LICENSE4":     true,
				},
			},
			want: false,
		},
		{
			name: "Not present in license",
			args: args{name: "EXPIRED_LICENSE"},
			fields: fields{
				license: map[string]bool{
					"CALL_CENTER":  true,
					"CALL_MANAGER": false,
					"CHATS":        true,
					"LICENSE1":     true,
					"LICENSE2":     true,
					"LICENSE3":     false,
					"LICENSE4":     true,
				},
			},
			want: false,
		},
		{
			name: "Combined licenses",
			args: args{name: "CALL_CENTER,CALL_MANAGER"},
			fields: fields{
				license: map[string]bool{
					"CALL_CENTER":  true,
					"CALL_MANAGER": false,
					"CHATS":        true,
					"LICENSE1":     true,
					"LICENSE2":     true,
					"LICENSE3":     false,
					"LICENSE4":     true,
				},
			},
		},
		{
			name: "Not present license but user has super rights",
			args: args{name: "UNKNOWN_LICENSE"},
			fields: fields{
				license: map[string]bool{
					"CALL_CENTER":  true,
					"CALL_MANAGER": false,
					"CHATS":        true,
					"LICENSE1":     true,
					"LICENSE2":     true,
					"LICENSE3":     false,
					"LICENSE4":     true,
				},
				superCreate: true,
				superSelect: true,
				superEdit:   true,
				superDelete: true,
			},
			want: false,
		},
		{
			name: "Present expired license and user has super rights",
			args: args{name: "CALL_MANAGER"},
			fields: fields{
				license: map[string]bool{
					"CALL_CENTER":  true,
					"CALL_MANAGER": false,
					"CHATS":        true,
					"LICENSE1":     true,
					"LICENSE2":     true,
					"LICENSE3":     false,
					"LICENSE4":     true,
				},
				superCreate: true,
				superSelect: true,
				superEdit:   true,
				superDelete: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserAuthSession{
				user:             tt.fields.user,
				permissions:      tt.fields.permissions,
				scopes:           tt.fields.scopes,
				license:          tt.fields.license,
				roles:            tt.fields.roles,
				domainId:         tt.fields.domainId,
				expiresAt:        tt.fields.expiresAt,
				superCreate:      tt.fields.superCreate,
				superEdit:        tt.fields.superEdit,
				superDelete:      tt.fields.superDelete,
				superSelect:      tt.fields.superSelect,
				mainAccess:       tt.fields.mainAccess,
				mainObjClassName: tt.fields.mainObjClassName,
			}
			if got := s.CheckLicenseAccess(tt.args.name); got != tt.want {
				t.Errorf("CheckLicenseAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAuthSession_CheckObacAccess(t *testing.T) {
	type fields struct {
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
	type args struct {
		scopeName  string
		accessType auth.AccessMode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Existing scope with enabled obac and name with read access",
			args: args{scopeName: "chats", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "xrwd", Obac: true},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with enabled obac and name with edit access",
			args: args{scopeName: "chats", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "xrwd", Obac: true},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with enabled obac and name with delete access",
			args: args{scopeName: "chats", accessType: auth.Delete},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "xrwd", Obac: true},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with enabled obac and name with create access",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "xrwd", Obac: true},
				},
			},
			want: true,
		},
		// without access
		{
			name: "Existing scope with enabled obac and name without read access",
			args: args{scopeName: "chats", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
			},
			want: false,
		},
		{
			name: "Existing scope with enabled obac and name without edit access",
			args: args{scopeName: "chats", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
			},
			want: false,
		},
		{
			name: "Existing scope with enabled obac and name without delete access",
			args: args{scopeName: "chats", accessType: auth.Delete},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
			},
			want: false,
		},
		{
			name: "Existing scope with enabled obac and name with create access",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
			},
			want: false,
		},
		// disabled obac
		{
			name: "Existing scope with disabled obac and name without read access",
			args: args{scopeName: "chats", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with disabled obac and name without edit access",
			args: args{scopeName: "chats", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with disabled obac and name without delete access",
			args: args{scopeName: "chats", accessType: auth.Delete},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with disabled obac and name with create access",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
			},
			want: true,
		},
		// super permissions
		{
			name: "Existing scope with enabled obac and name without read access but with super read",
			args: args{scopeName: "chats", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
				superSelect: true,
			},
			want: true,
		},
		{
			name: "Existing scope with enabled obac and name without edit access but with super edit",
			args: args{scopeName: "chats", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
				superEdit: true,
			},
			want: true,
		},
		{
			name: "Existing scope with enabled obac and name without delete access but with super delete",
			args: args{scopeName: "chats", accessType: auth.Delete},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
				superDelete: true,
			},
			want: true,
		},
		{
			name: "Existing scope with enabled obac and name without create access but with super add",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: false},
				},
				superCreate: true,
			},
			want: true,
		},
		// super permissions that not match required permission
		{
			name: "Existing scope with enabled obac and name without read access but with super read",
			args: args{scopeName: "chats", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
				superEdit:   true,
				superDelete: true,
				superCreate: true,
			},
			want: false,
		},
		{
			name: "Existing scope with enabled obac and name without edit access but with super edit",
			args: args{scopeName: "chats", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
				superSelect: true,
				superDelete: true,
				superCreate: true,
			},
			want: false,
		},
		{
			name: "Existing scope with enabled obac and name without delete access but with super delete",
			args: args{scopeName: "chats", accessType: auth.Delete},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
				superSelect: true,
				superEdit:   true,
				superCreate: true,
			},
			want: false,
		},
		{
			name: "Existing scope with enabled obac and name without create access but with super add",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Access: "", Obac: true},
				},
				superSelect: true,
				superEdit:   true,
				superDelete: true,
			},
			want: false,
		},
		{
			name: "Non-Existing scope",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: make(map[string]*Scope),
			},
			want: false,
		},
		{
			name: "Non-Existing scope with all super permissions",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes:      make(map[string]*Scope),
				superSelect: true,
				superEdit:   true,
				superDelete: true,
				superCreate: true,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserAuthSession{
				user:             tt.fields.user,
				permissions:      tt.fields.permissions,
				scopes:           tt.fields.scopes,
				license:          tt.fields.license,
				roles:            tt.fields.roles,
				domainId:         tt.fields.domainId,
				expiresAt:        tt.fields.expiresAt,
				superCreate:      tt.fields.superCreate,
				superEdit:        tt.fields.superEdit,
				superDelete:      tt.fields.superDelete,
				superSelect:      tt.fields.superSelect,
				mainAccess:       tt.fields.mainAccess,
				mainObjClassName: tt.fields.mainObjClassName,
			}
			if got := s.CheckObacAccess(tt.args.scopeName, tt.args.accessType); got != tt.want {
				t.Errorf("CheckObacAccess() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAuthSession_GetDomainId(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserAuthSession{
				user:             tt.fields.user,
				permissions:      tt.fields.permissions,
				scopes:           tt.fields.scopes,
				license:          tt.fields.license,
				roles:            tt.fields.roles,
				domainId:         tt.fields.domainId,
				expiresAt:        tt.fields.expiresAt,
				superCreate:      tt.fields.superCreate,
				superEdit:        tt.fields.superEdit,
				superDelete:      tt.fields.superDelete,
				superSelect:      tt.fields.superSelect,
				mainAccess:       tt.fields.mainAccess,
				mainObjClassName: tt.fields.mainObjClassName,
			}
			if got := s.GetDomainId(); got != tt.want {
				t.Errorf("GetDomainId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAuthSession_GetRoles(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   []int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserAuthSession{
				user:             tt.fields.user,
				permissions:      tt.fields.permissions,
				scopes:           tt.fields.scopes,
				license:          tt.fields.license,
				roles:            tt.fields.roles,
				domainId:         tt.fields.domainId,
				expiresAt:        tt.fields.expiresAt,
				superCreate:      tt.fields.superCreate,
				superEdit:        tt.fields.superEdit,
				superDelete:      tt.fields.superDelete,
				superSelect:      tt.fields.superSelect,
				mainAccess:       tt.fields.mainAccess,
				mainObjClassName: tt.fields.mainObjClassName,
			}
			if got := s.GetRoles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAuthSession_GetUserId(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   int64
	}{
		{
			name:   "Nil user",
			fields: fields{},
			want:   0,
		},
		{
			name: "Not nil user with id < 0",
			fields: fields{
				user: &User{Id: -100},
			},
			want: 0,
		},
		{
			name: "Not nil user with id = 0",
			fields: fields{
				user: &User{Id: -100},
			},
			want: 0,
		},
		{
			name: "Not nil user with id > 0",
			fields: fields{
				user: &User{Id: 10},
			},
			want: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserAuthSession{
				user:             tt.fields.user,
				permissions:      tt.fields.permissions,
				scopes:           tt.fields.scopes,
				license:          tt.fields.license,
				roles:            tt.fields.roles,
				domainId:         tt.fields.domainId,
				expiresAt:        tt.fields.expiresAt,
				superCreate:      tt.fields.superCreate,
				superEdit:        tt.fields.superEdit,
				superDelete:      tt.fields.superDelete,
				superSelect:      tt.fields.superSelect,
				mainAccess:       tt.fields.mainAccess,
				mainObjClassName: tt.fields.mainObjClassName,
			}
			if got := s.GetUserId(); got != tt.want {
				t.Errorf("GetUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAuthSession_IsExpired(t *testing.T) {
	type fields struct {
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
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserAuthSession{
				user:             tt.fields.user,
				permissions:      tt.fields.permissions,
				scopes:           tt.fields.scopes,
				license:          tt.fields.license,
				roles:            tt.fields.roles,
				domainId:         tt.fields.domainId,
				expiresAt:        tt.fields.expiresAt,
				superCreate:      tt.fields.superCreate,
				superEdit:        tt.fields.superEdit,
				superDelete:      tt.fields.superDelete,
				superSelect:      tt.fields.superSelect,
				mainAccess:       tt.fields.mainAccess,
				mainObjClassName: tt.fields.mainObjClassName,
			}
			if got := s.IsExpired(); got != tt.want {
				t.Errorf("IsExpired() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserAuthSession_IsRbacCheckRequired(t *testing.T) {
	type fields struct {
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
	type args struct {
		scopeName  string
		accessType auth.AccessMode
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Non-existing scope",
			args: args{scopeName: "non-existent", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: false,
		},
		{
			name: "Non-existing scope",
			args: args{scopeName: "non-existent", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: false,
		},
		{
			name: "Non-existing scope",
			args: args{scopeName: "non-existent", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: false,
		},
		{
			name: "Non-existing scope",
			args: args{scopeName: "non-existent", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: false,
		},

		// existing scope
		{
			name: "Existing scope with edit",
			args: args{scopeName: "chats", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with delete",
			args: args{scopeName: "chats", accessType: auth.Delete},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with add",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: true,
		},
		{
			name: "Existing scope with read",
			args: args{scopeName: "chats", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: true},
				},
			},
			want: true,
		},
		// existing with rbac disabled
		{
			name: "Existing scope with edit",
			args: args{scopeName: "chats", accessType: auth.Edit},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: false},
				},
			},
			want: false,
		},
		{
			name: "Existing scope with delete",
			args: args{scopeName: "chats", accessType: auth.Delete},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: false},
				},
			},
			want: false,
		},
		{
			name: "Existing scope with add",
			args: args{scopeName: "chats", accessType: auth.Add},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: false},
				},
			},
			want: false,
		},
		{
			name: "Existing scope with read",
			args: args{scopeName: "chats", accessType: auth.Read},
			fields: fields{
				scopes: map[string]*Scope{
					"chats": {Rbac: false},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &UserAuthSession{
				user:             tt.fields.user,
				permissions:      tt.fields.permissions,
				scopes:           tt.fields.scopes,
				license:          tt.fields.license,
				roles:            tt.fields.roles,
				domainId:         tt.fields.domainId,
				expiresAt:        tt.fields.expiresAt,
				superCreate:      tt.fields.superCreate,
				superEdit:        tt.fields.superEdit,
				superDelete:      tt.fields.superDelete,
				superSelect:      tt.fields.superSelect,
				mainAccess:       tt.fields.mainAccess,
				mainObjClassName: tt.fields.mainObjClassName,
			}
			if got := s.IsRbacCheckRequired(tt.args.scopeName, tt.args.accessType); got != tt.want {
				t.Errorf("IsRbacCheckRequired() = %v, want %v", got, tt.want)
			}
		})
	}
}
