package registry

import (
	"github.com/webitel/cases/model"
	"time"
)

const (
	DeregisterCriticalServiceAfter = 30 * time.Second
	ServiceName                    = "cases"
	CheckInterval                  = 1 * time.Minute
)

type ServiceRegistrator interface {
	Register() model.AppError
	Deregister() model.AppError
}
