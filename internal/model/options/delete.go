package options

import (
	"context"
	"time"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/util"
)

type Deleter interface {
	context.Context
	GetAuthOpts() auth.Auther
	RequestTime() time.Time
	GetFields() []string

	// Additional filtering

	GetFilter(field string) []util.FilterExpr
	RemoveFilter(string)
	AddFilter(string)
	GetFilters() []string

	// If connection to parent object required
	GetParentID() int64

	// ID filtering
	GetIDs() []int64
}
