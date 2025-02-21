package model

import (
	"context"
	"errors"
	"time"

	"github.com/webitel/cases/auth"

	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"
)

func NewSearchOptions(ctx context.Context, searcher Lister, objMetadata ObjectMetadatter) (*SearchOptions, error) {
	opts := &SearchOptions{
		Context: ctx,
		Page:    int(searcher.GetPage()),
		Size:    int(searcher.GetSize()),
		Sort:    searcher.GetSort(),
		Search:  searcher.GetQ(),
		Filter:  make(map[string]any),
	}
	if sess := GetAutherOutOfContext(ctx); sess != nil {
		opts.Auth = sess
	} else {
		return nil, errors.New("can't authorize user")
	}
	// set current time
	opts.CurrentTime()
	// normalize fields and deduplicate fields
	var resultingFields []string
	if requestedFields := searcher.GetFields(); len(requestedFields) == 0 {
		resultingFields = objMetadata.GetDefaultFields()
	} else {
		resultingFields = util.DeduplicateFields(util.FieldsFunc(
			requestedFields, graph.SplitFieldsQ,
		))

	}

	resultingFields, opts.UnknownFields = util.SplitKnownAndUnknownFields(resultingFields, objMetadata.GetAllFields())
	opts.Fields = util.ParseFieldsForEtag(resultingFields)
	return opts, nil
}

type SearchOptions struct {
	Time time.Time
	context.Context
	//Session *session.Session
	// filters
	Filter    map[string]any
	Search    string
	IDs       []int64
	ParentId  int64
	ContactId int64
	// output
	Fields            []string
	UnknownFields     []string
	DerivedSearchOpts map[string]*SearchOptions
	// paging
	Page int
	Size int
	Sort string
	// filtering by single id
	ID int64
	// Auth opts
	Auth auth.Auther
}

func (s *SearchOptions) SearchDerivedOptionByField(field string) *SearchOptions {
	for s2, options := range s.DerivedSearchOpts {
		if s2 == field {
			return options
		}
	}
	return nil
}

func (s *SearchOptions) SetAuthOpts(a auth.Auther) *SearchOptions {
	s.Auth = a
	return s
}

func (s *SearchOptions) GetAuthOpts() auth.Auther {
	return s.Auth
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

func (s *SearchOptions) GetPage() int {
	if s != nil {
		// Limited ? either: manual -or- default !
		if s.GetSize() > 0 {
			// Valid ?page= specified ?
			if s.Page > 0 {
				return s.Page
			}
			// default: always the first one !
			return 1
		}
	}
	// <nop> -or- <nolimit>
	return 0
}

func (s *SearchOptions) GetSort() string {
	return s.Sort
}

func (s *SearchOptions) CurrentTime() time.Time {
	ts := s.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		s.Time = ts
	}
	return ts
}

type Lister interface {
	Fielder
	Pager
	Searcher
	Sorter
}

type Sorter interface {
	GetSort() string
}

type Pager interface {
	GetPage() int32
	GetSize() int32
}

type Searcher interface {
	GetQ() string
}

type Fielder interface {
	GetFields() []string
}

func NewLocateOptions(ctx context.Context, locator Fielder, objMetadata ObjectMetadatter) (*SearchOptions, error) {
	opts := &SearchOptions{
		Context: ctx,
		//Session: ctx.Value(interceptor.SessionHeader).(*session.Session),
		Time: time.Now(),
		Page: 1,
		Size: 1,
	}
	// set current time
	opts.CurrentTime()
	if sess := GetAutherOutOfContext(ctx); sess != nil {
		opts.Auth = sess
	} else {
		return nil, errors.New("can't authorize user")
	}

	// normalize fields
	var resultingFields []string
	if requestedFields := locator.GetFields(); len(requestedFields) == 0 {
		resultingFields = make([]string, len(objMetadata.GetDefaultFields()))
		copy(resultingFields, objMetadata.GetDefaultFields())
	} else {
		resultingFields = util.DeduplicateFields(util.FieldsFunc(
			requestedFields, graph.SplitFieldsQ,
		))
	}
	resultingFields, opts.UnknownFields = util.SplitKnownAndUnknownFields(resultingFields, objMetadata.GetAllFields())
	opts.Fields = util.ParseFieldsForEtag(resultingFields)
	return opts, nil
}

// DeafaultSearchSize is a constant integer == 10.
const (
	DefaultSearchSize = 10
)
