package model

import (
	"context"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/server/interceptor"
)

const (
	APP_SERVICE_NAME = "cases"
	NAMESPACE_NAME   = "webitel"
)

func GetAutherOutOfContext(ctx context.Context) auth.Auther {
	return ctx.Value(interceptor.SessionHeader).(auth.Auther)
}
