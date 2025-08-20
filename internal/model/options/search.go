package options

import (
	"context"
	"time"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/model"
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
	GetFilters() model.Filterer
	AddFilter(model.ConnectionType, model.Filterer)
	// shortcuts
	GetIDs() []int64
	GetQin() string
}
