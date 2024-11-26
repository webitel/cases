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
	Session  *session.Session
	Filter   map[string]interface{}
	Search   string
	IDs      []int64
	Sort     []string
	Fields   []string
	ParentId int64
	Page     int32
	Size     int32
}

type Searcher interface {
	Locator
	GetPage() int32
	GetSize() int32
	GetQ() string
}

type Locator interface {
	GetFields() []string
}

func NewSearchOptions(ctx context.Context, searcher Searcher, defaultFields []string) *SearchOptions {
	opts := &SearchOptions{
		Context: ctx,
		Session: ctx.Value(interceptor.SessionHeader).(*session.Session),
		Page:    searcher.GetPage(),
		Size:    searcher.GetSize(),
		Search:  searcher.GetQ(),
	}
	// set current time
	opts.CurrentTime()

	// normalize fields
	fields := util.FieldsFunc(
		searcher.GetFields(), graph.SplitFieldsQ,
	)
	if len(fields) == 0 {
		fields = make([]string, len(defaultFields))
		copy(fields, defaultFields)
	}
	opts.Fields = util.ParseFieldsForEtag(fields)
	return opts
}

func NewLocateOptions(ctx context.Context, locator Locator, defaultFields []string) *SearchOptions {
	opts := &SearchOptions{
		Context: ctx,
		Session: ctx.Value(interceptor.SessionHeader).(*session.Session),
		Time:    time.Now(),
		Page:    1,
		Size:    1,
	}
	// set current time
	opts.CurrentTime()

	// normalize fields
	fields := util.FieldsFunc(
		locator.GetFields(), graph.SplitFieldsQ,
	)
	if len(fields) == 0 {
		copy(fields, defaultFields)
	}
	// find and process etag field
	opts.Fields = util.ParseFieldsForEtag(fields)
	return opts
}

// DeafaultSearchSize is a constant integer == 16.
const (
	DefaultSearchSize = 10
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
