package model

import (
	"context"
	"time"

	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"
)

func NewSearchOptions(ctx context.Context, searcher Lister, objMetadata ObjectMetadatter) *SearchOptions {
	opts := &SearchOptions{
		Context: ctx,
		Page:    int(searcher.GetPage()),
		Size:    int(searcher.GetSize()),
		Search:  searcher.GetQ(),
		Filter:  make(map[string]any),
	}
	if sess := GetSessionOutOfContext(ctx); sess != nil {
		opts.Auth = NewSessionAuthOptions(sess, objMetadata.GetAllScopeNames()...)
	}
	// set current time
	opts.CurrentTime()
	// normalize fields
	var resultingFields []string
	if requestedFields := searcher.GetFields(); len(requestedFields) == 0 {
		resultingFields = objMetadata.GetDefaultFields()
	} else {
		resultingFields = util.FieldsFunc(
			requestedFields, graph.SplitFieldsQ,
		)
	}

	resultingFields, opts.UnknownFields = util.SplitKnownAndUnknownFields(resultingFields, objMetadata.GetAllFields())
	opts.Fields = util.ParseFieldsForEtag(resultingFields)
	return opts
}

type SearchOptions struct {
	Time time.Time
	context.Context
	//Session *session.Session
	// filters
	Filter   map[string]any
	Search   string
	IDs      []int64
	ParentId int64
	// output
	Fields            []string
	UnknownFields     []string
	DerivedSearchOpts map[string]*SearchOptions
	// paging
	Page int
	Size int
	Sort []string
	// filtering by single id
	ID int64
	// Auth opts
	Auth Auther
}

func (s *SearchOptions) SearchDerivedOptionByField(field string) *SearchOptions {
	for s2, options := range s.DerivedSearchOpts {
		if s2 == field {
			return options
		}
	}
	return nil
}

func (s *SearchOptions) SetAuthOpts(a Auther) *SearchOptions {
	s.Auth = a
	return s
}

func (s *SearchOptions) GetAuthOpts() Auther {
	return s.Auth
}

func (rpc *SearchOptions) GetSize() int {
	if rpc == nil {
		return DefaultSearchSize
	}
	switch {
	case rpc.Size < 0:
		return -1
	case rpc.Size > 0:
		// CHECK for too big values !
		return rpc.Size
	case rpc.Size == 0:
		return DefaultSearchSize
	}
	panic("unreachable code")
}

func (rpc *SearchOptions) GetPage() int {
	if rpc != nil {
		// Limited ? either: manual -or- default !
		if rpc.GetSize() > 0 {
			// Valid ?page= specified ?
			if rpc.Page > 0 {
				return rpc.Page
			}
			// default: always the first one !
			return 1
		}
	}
	// <nop> -or- <nolimit>
	return 0
}

func (rpc *SearchOptions) GetSort() []string {
	return rpc.Sort
}

func (rpc *SearchOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}

type Lister interface {
	Fielder
	Pager
	Searcher
}

type Sorter interface {
	GetSort() []string
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

func NewLocateOptions(ctx context.Context, locator Fielder, objMetadata ObjectMetadatter) *SearchOptions {
	opts := &SearchOptions{
		Context: ctx,
		//Session: ctx.Value(interceptor.SessionHeader).(*session.Session),
		Time: time.Now(),
		Page: 1,
		Size: 1,
	}
	// set current time
	opts.CurrentTime()
	if sess := GetSessionOutOfContext(ctx); sess != nil {
		opts.Auth = NewSessionAuthOptions(sess, objMetadata.GetAllScopeNames()...)
	}

	// normalize fields
	var resultingFields []string
	if requestedFields := locator.GetFields(); len(requestedFields) == 0 {
		resultingFields = make([]string, len(objMetadata.GetDefaultFields()))
		copy(resultingFields, objMetadata.GetDefaultFields())
	} else {
		resultingFields = util.FieldsFunc(
			requestedFields, graph.SplitFieldsQ,
		)
	}
	resultingFields, opts.UnknownFields = util.SplitKnownAndUnknownFields(resultingFields, objMetadata.GetAllFields())
	opts.Fields = util.ParseFieldsForEtag(resultingFields)
	return opts
}

// DeafaultSearchSize is a constant integer == 10.
const (
	DefaultSearchSize = 10
)
