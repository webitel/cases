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
	GetFields() []string

	// Additional filtering

	GetFilter(string) (string, bool)
	RemoveFilter(string)
	AddFilter(string)
	GetFilters() []string

	// If connection to parent object required
	GetParentID() int64

	// ID filtering
	GetIDs() []int64
}
