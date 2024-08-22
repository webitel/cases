package model

import (
	"context"
	"time"

	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/util"
)

type SearchOptions struct {
	FieldsUtil util.FieldsUtils
	Match      util.Match
	Time       time.Time
	context.Context
	Session *session.Session
	Filter  map[string]interface{}
	Search  string
	IDs     []int64
	Sort    []string
	Fields  []string
	Page    int
	Size    int
}

// DeafaultSearchSize is a constant integer == 16.
const (
	DefaultSearchSize = 16
)

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

func (rpc *SearchOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}
