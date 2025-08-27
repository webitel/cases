package options

import (
	"context"
	"strings"
	"time"

	"github.com/google/cel-go/cel"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/api_handler/grpc/options/shared"
	optsutil "github.com/webitel/cases/internal/api_handler/grpc/options/util"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/model/options/defaults"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"github.com/webitel/webitel-go-kit/pkg/filters"
	"google.golang.org/grpc/codes"
)

type SearchOption func(options *SearchOptions) error

var _ options.Searcher = (*SearchOptions)(nil)

type SearchOptions struct {
	createdAt time.Time
	context.Context
	// filters
	IDs []int64
	// Deprecated: use FiltersV1
	Filters   []string
	FiltersV1 filters.Filterer
	// search
	Search string
	// output
	Fields        []string
	UnknownFields []string
	CustomContext map[string]any

	// paging
	Page int
	Size int
	Sort string
	// Auth opts
	Auth auth.Auther
	Qin  string
}

type Pager interface {
	GetPage() int32
	GetSize() int32
}

type Sorter interface {
	GetSort() string
}

type Searcher interface {
	GetQ() string
}

type Filterer interface {
	GetFilters() []string
}

func WithFilters(filters []string) SearchOption {
	return func(options *SearchOptions) error {
		options.Filters = filters
		return nil
	}
}

func WithFiltersV1(env *cel.Env, query string) SearchOption {
	return func(options *SearchOptions) error {
		filters, err := filters.ParseFilters(env, query)
		if err != nil {
			return err
		}
		options.FiltersV1 = filters
		return nil
	}
}

func WithFields(fielder shared.Fielder, md model.ObjectMetadatter, fieldModifiers ...func(in []string) []string) SearchOption {
	return func(options *SearchOptions) error {
		if requestedFields := fielder.GetFields(); len(requestedFields) == 0 {
			options.Fields = md.GetDefaultFields()

		} else {
			options.Fields = util.FieldsFunc(requestedFields, util.InlineFields)
		}
		for _, modifier := range fieldModifiers {
			options.Fields = modifier(options.Fields)
		}
		options.Fields, options.UnknownFields = util.SplitKnownAndUnknownFields(options.Fields, md.GetAllFields())

		return nil
	}
}

func WithPagination(pager Pager) SearchOption {
	return func(ops *SearchOptions) error {
		ops.Page = int(pager.GetPage())
		ops.Size = int(pager.GetSize())
		if ops.Page == 0 {
			ops.Page = 1
		}
		if ops.Size < 0 {
			ops.Size = -1
		}
		if ops.Size == 0 {
			ops.Size = options.DefaultSearchSize
		}
		return nil
	}
}

func WithSearch(searcher Searcher) SearchOption {
	return func(options *SearchOptions) error {
		if s := searcher.GetQ(); s != "" {
			options.Search = searcher.GetQ()
		}
		return nil
	}
}

func WithSearchAsParam(query string) SearchOption {
	return func(options *SearchOptions) error {
		options.Search = query
		return nil
	}
}

func WithIDs(ids []int64) SearchOption {
	return func(options *SearchOptions) error {
		options.IDs = ids
		return nil
	}
}

func WithID(id int64) SearchOption {
	return func(options *SearchOptions) error {
		options.IDs = append(options.IDs, id)
		return nil
	}
}

func WithSort(sorter Sorter) SearchOption {
	return func(options *SearchOptions) error {
		options.Sort = sorter.GetSort()
		return nil
	}
}

func WithQin(qin string) SearchOption {
	return func(o *SearchOptions) error {
		o.Qin = strings.ToLower(qin)

		return nil
	}
}

func WithIDsAsEtags(tag etag.EtagType, etags ...string) SearchOption {
	return func(options *SearchOptions) error {
		ids, err := util.ParseIds(etags, tag)
		if err != nil {
			return err
		}
		options.IDs = ids
		return nil
	}
}

func (s *SearchOptions) GetAuthOpts() auth.Auther {
	return s.Auth
}

func (s *SearchOptions) RequestTime() time.Time {
	return s.createdAt
}

func (s *SearchOptions) GetQin() string {
	if s == nil {
		return ""
	}
	return s.Qin
}

func (s *SearchOptions) GetFields() []string {
	return s.Fields
}

func (s *SearchOptions) GetUnknownFields() []string {
	return s.UnknownFields
}

func (s *SearchOptions) GetSearch() string {
	return s.Search
}

func (s *SearchOptions) GetPage() int {
	return s.Page
}

func (s *SearchOptions) GetSize() int {
	if s == nil {
		return defaults.DefaultSearchSize
	}
	switch {
	case s.Size < 0:
		return -1
	case s.Size > 0:
		// CHECK for too big values !
		return s.Size
	case s.Size == 0:
		return defaults.DefaultSearchSize
	}
	panic("unreachable code")
}

func (s *SearchOptions) GetSort() string {
	return s.Sort
}

func (s *SearchOptions) GetFiltersV1() filters.Filterer {
	return s.FiltersV1
}

// Deprecated: use GetFiltersV1
func (s *SearchOptions) GetFilters() []string {
	return s.Filters
}

func (s *SearchOptions) AddFilter(f string) {
	s.Filters = append(s.Filters, f)
}

// Deprecated
// GetFilter returns all filters for a given field, with operator and value
func (s *SearchOptions) GetFilter(f string) []util.FilterExpr {
	return util.GetFilter(s.Filters, f)
}

func (s *SearchOptions) GetCustomContext() map[string]any {
	return s.CustomContext
}

func (s *SearchOptions) GetIDs() []int64 {
	return s.IDs
}

func (s *SearchOptions) AddCustomContext(key string, value any) {
	if s.CustomContext == nil {
		s.CustomContext = make(map[string]any)
	}
	s.CustomContext[key] = value
}

func NewSearchOptions(ctx context.Context, opts ...SearchOption) (*SearchOptions, error) {
	search := &SearchOptions{
		createdAt:     time.Now().UTC(),
		Context:       ctx,
		CustomContext: make(map[string]any),
	}
	if sess := optsutil.GetAutherOutOfContext(ctx); sess != nil {
		search.Auth = sess
	} else {
		return nil, errors.New("can't authorize user", errors.WithCode(codes.Unauthenticated))
	}
	for _, opt := range opts {
		err := opt(search)
		if err != nil {
			return nil, err
		}
	}
	return search, nil
}

func NewLocateOptions(ctx context.Context, opts ...SearchOption) (*SearchOptions, error) {
	locate, err := NewSearchOptions(ctx, opts...)
	if err != nil {
		return nil, err
	}
	if len(locate.IDs) == 0 {
		return nil, errors.New("locate options require id to locate", errors.WithCode(codes.InvalidArgument))
	}
	if len(locate.IDs) > 1 {
		return nil, errors.New("locate options require only one id", errors.WithCode(codes.InvalidArgument))
	}

	return locate, nil
}
