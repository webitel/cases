package options

import (
	"context"
	"github.com/webitel/cases/auth"
	"time"
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
	GetFilters() map[string]any
	RemoveFilter(string)
	AddFilter(string, any)
	GetFilter(string) any
	// shortcuts
	GetIDs() []int64
}
