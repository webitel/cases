package options

import (
	"context"
	"github.com/webitel/cases/auth"
	"time"
)

type Creator interface {
	context.Context
	GetAuthOpts() auth.Auther
	RequestTime() time.Time
	GetFields() []string
	GetDerivedSearchOpts() map[string]*Searcher
	GetUnknownFields() []string
	GetIDs() []int64
	GetParentID() int64
	GetChildID() int64
}
