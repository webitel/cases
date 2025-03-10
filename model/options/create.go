package options

import (
	"context"
	"github.com/webitel/cases/auth"

	"time"
)

type CreateOptions interface {
	context.Context
	GetTime() time.Time
	GetFields() []string
	GetDerivedSearchOpts() map[string]*SearchOptions
	GetUnknownFields() []string
	GetIDs() []int64
	GetParentID() int64
	GetChildID() int64
	GetAuth() auth.Auther
}
