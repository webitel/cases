package interceptor

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/model"
)

type obacSession struct {
	auth.Auther

	obac map[string]bool
}

func (s *obacSession) CheckObacAccess(scope string, _ auth.AccessMode) bool {
	return s.obac[scope]
}

func TestValidateSessionPermission(t *testing.T) {
	tests := []struct {
		name     string
		objClass string
		obac     map[string]bool
		want     bool
	}{
		{name: "external comments granted", objClass: model.ScopeCaseComments, obac: map[string]bool{model.ScopeCaseComments: true}, want: true},
		{name: "internal comments only", objClass: model.ScopeCaseComments, obac: map[string]bool{model.ScopeInternalComments: true}, want: true},
		{name: "no comment access", objClass: model.ScopeCaseComments, obac: map[string]bool{}, want: false},
		{name: "no alternative for cases", objClass: model.ScopeCases, obac: map[string]bool{model.ScopeInternalComments: true}, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.want, validateSessionPermission(&obacSession{obac: tc.obac}, tc.objClass, auth.Read))
		})
	}
}
