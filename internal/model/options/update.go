package options

import (
	"context"
	"github.com/webitel/cases/auth"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"time"
)

type Updator interface {
	context.Context
	GetAuthOpts() auth.Auther
	GetFields() []string
	GetUnknownFields() []string
	GetDerivedSearchOpts() map[string]*Searcher
	RequestTime() time.Time
	GetMask() []string
	GetEtags() []*etag.Tid
	GetParentID() int64
	GetIDs() []int64
}
