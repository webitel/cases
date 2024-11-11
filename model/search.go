package model

import (
	"context"
	"time"

	"github.com/webitel/cases/internal/server/interceptor"
	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"

	session "github.com/webitel/cases/auth/model"
)

type SearchOptions struct {
	Time time.Time
	context.Context
	Session *session.Session
	Filter  map[string]interface{}
	Search  string
	IDs     []int64
	Sort    []string
	Fields  []string
	Id      int64
	Page    int32
	Size    int32
}

type Searcher interface {
	Locator
	GetPage() int32
	GetSize() int32
}

type Locator interface {
	GetFields() []string
}

func NewSearchOptions(ctx context.Context, searcher Searcher) *SearchOptions {
	sess := ctx.Value(interceptor.SessionHeader).(*session.Session)

	return &SearchOptions{
		Context: ctx,
		Session: sess,
		Fields: util.FieldsFunc(
			searcher.GetFields(), graph.SplitFieldsQ, // explode: by COMMA(',')
		),
		Time: time.Now(),
		Page: searcher.GetPage(),
		Size: searcher.GetSize(),
	}
}

func NewLocateOptions(ctx context.Context, locator Locator) *SearchOptions {
	sess := ctx.Value(interceptor.SessionHeader).(*session.Session)

	return &SearchOptions{
		Context: ctx,
		Session: sess,
		Fields: util.FieldsFunc(
			locator.GetFields(), graph.SplitFieldsQ,
		),
		Time: time.Now(),
		Page: 1,
		Size: 1,
	}
}

// DeafaultSearchSize is a constant integer == 16.
const (
	DefaultSearchSize = 16
)

func (rpc *SearchOptions) GetSize() int32 {
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

func (rpc *SearchOptions) GetPage() int32 {
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

func (rpc *SearchOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}
