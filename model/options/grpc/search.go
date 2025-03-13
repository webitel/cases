package grpc

import (
	"context"
	"errors"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/model/options/grpc/shared"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"strings"
	"time"
)

type SearchOption func(options *SearchOptions) error

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

func WithFields(fielder shared.Fielder, md model.ObjectMetadatter, fieldModifiers ...func(in []string) []string) SearchOption {
	return func(options *SearchOptions) error {
		if requestedFields := fielder.GetFields(); len(requestedFields) == 0 {
			options.Fields = md.GetDefaultFields()

		} else {
			options.Fields = requestedFields
		}
		for _, modifier := range fieldModifiers {
			options.Fields = modifier(options.Fields)
		}
		options.Fields, options.UnknownFields = util.SplitKnownAndUnknownFields(options.Fields, md.GetAllFields())

		return nil
	}
}

func WithPagination(pager Pager) SearchOption {
	return func(options *SearchOptions) error {
		options.Page = int(pager.GetPage())
		options.Size = int(pager.GetSize())
		if options.Page == 0 {
			options.Page = 1
		}
		if options.Size < 0 {
			options.Size = -1
		}
		return nil
	}
}

func WithFilters(filterer Filterer) SearchOption {
	return func(options *SearchOptions) error {
		for _, s := range filterer.GetFilters() {
			str := strings.Split(s, "=")
			if len(str) != 2 {
				continue
			}
			column := str[0]
			value := strings.TrimSpace(str[1])
			options.AddFilter(column, value)

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

type SearchOptions struct {
	createdAt time.Time
	context.Context
	// filters
	IDs     []int64
	Filters map[string]any
	// search
	Search string
	// output
	Fields        []string
	UnknownFields []string
	// paging
	Page int
	Size int
	Sort string
	// Auth opts
	Auth auth.Auther
}

func (s *SearchOptions) GetAuthOpts() auth.Auther {
	return s.Auth
}

func (s *SearchOptions) RequestTime() time.Time {
	return s.createdAt
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
		return DefaultSearchSize
	}
	switch {
	case s.Size < 0:
		return -1
	case s.Size > 0:
		// CHECK for too big values !
		return s.Size
	case s.Size == 0:
		return DefaultSearchSize
	}
	panic("unreachable code")
}

func (s *SearchOptions) GetSort() string {
	return s.Sort
}

func (s *SearchOptions) GetFilters() map[string]any {
	return s.Filters
}

func (s *SearchOptions) RemoveFilter(key string) {
	delete(s.Filters, key)
}

func (s *SearchOptions) AddFilter(key string, value any) {
	s.Filters[key] = value
}

func (s *SearchOptions) GetFilter(key string) any {
	return s.Filters[key]
}

func (s *SearchOptions) GetIDs() []int64 {
	return s.IDs
}

func NewSearchOptions(ctx context.Context, opts ...SearchOption) (*SearchOptions, error) {
	search := &SearchOptions{
		createdAt: time.Now().UTC(),
		Context:   ctx,
		Filters:   make(map[string]any),
	}
	if sess := model.GetAutherOutOfContext(ctx); sess != nil {
		search.Auth = sess
	} else {
		return nil, errors.New("can't authorize user")
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
		return nil, errors.New("locate options require id to locate")
	}
	if len(locate.IDs) > 1 {
		return nil, errors.New("locate options require only one id")
	}

	return locate, nil
}

// DeafaultSearchSize is a constant integer == 10.
const (
	DefaultSearchSize = 10
)
