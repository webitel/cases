package options

import (
	"context"
	"time"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/filters"
)

const DefaultSearchSize = 10

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
	// Deprecated: use GetFiltersV1
	GetFilters() []string
	// Deprecated: use FiltersV1
	AddFilter(string)
	// Deprecated: use FiltersV1
	GetFilter(f string) []util.FilterExpr

	GetFiltersV1() *filters.FilterExpr
	// shortcuts
	GetIDs() []int64
	GetQin() string

	AddCustomContext(key string, value any)
	GetCustomContext() map[string]any
}
