package options

import (
	"context"
	"github.com/webitel/cases/auth"
	"time"
)

type Deleter interface {
	context.Context
	GetAuthOpts() auth.Auther
	RequestTime() time.Time

	// Additional filtering

	GetFilters() map[string]any
	RemoveFilter(string)
	AddFilter(string, any)
	GetFilter(string) any

	// If connection to parent object required
	GetParentID() int64

	// ID filtering
	GetIDs() []int64
}
