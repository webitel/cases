package registry

import (
	"time"

	"github.com/webitel/cases/model"
)

const (
	DeregisterCriticalServiceAfter = 30 * time.Second
	ServiceName                    = "webitel.cases"
	CheckInterval                  = 1 * time.Minute
)

type ServiceRegistrator interface {
	Register() model.AppError
	Deregister() model.AppError
}
