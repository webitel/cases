package options

import (
	"context"
	"github.com/webitel/cases/auth"
	"github.com/webitel/webitel-go-kit/etag"
	"time"
)

type UpdateOptions interface {
	context.Context
	GetAuthOpts() auth.Auther
	GetFields() []string
	GetUnknownFields() []string
	GetDerivedSearchOpts() map[string]*SearchOptions
	RequestTime() time.Time
	GetMask() []string
	GetEtags() []*etag.Tid
	GetParentID() int64
	GetIDs() []int64
}
