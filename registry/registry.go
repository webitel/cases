package registry

import (
	"time"
)

const (
	DeregisterCriticalServiceAfter = 30 * time.Second
	ServiceName                    = "webitel.cases"
	CheckInterval                  = 1 * time.Minute
)

// ServiceRegistrator interface for managing service registration.
type ServiceRegistrator interface {
	Register() error
	Deregister() error
}
