package options

import (
	"context"
	"time"

	"github.com/webitel/cases/auth"
)

type Searcher interface {
	context.Context
	GetAuthOpts() auth.Auther
	RequestTime() time.Time
	GetFields() []string
	GetUnknownFields() []string
	GetSearch() string
	// Paging
	GetPage() int
	GetSize() int
	// Sorting
	GetSort() string
	// Filtering
	AddFilter(string)
	GetFilters() []string
	GetFilter(string) (string, bool)
	RemoveFilter(string)
	// shortcuts
	GetIDs() []int64
	AddCustomContext(key string, value any)
	GetCustomContext() map[string]any
}
