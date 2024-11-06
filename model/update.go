package model

import (
	"context"
	"time"

	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/server/interceptor"
)

// UpdateOptions defines options for updating an entity with fields, mask, filter, and pagination
type UpdateOptions struct {
	Time time.Time
	context.Context
	Session *session.Session
	Fields  []string
	Mask    []string
	IDs     []int64
	// ID      int64
}

type Updator interface {
	GetFields() []string
	GetXJsonMask() []string
}

// NewUpdateOptions initializes UpdateOptions with values from a context and an Updator-compliant struct
func NewUpdateOptions(ctx context.Context, req Updator) *UpdateOptions {
	sess := ctx.Value(interceptor.SessionHeader).(*session.Session)

	return &UpdateOptions{
		Context: ctx,
		Session: sess,
		Fields:  req.GetFields(),
		Mask:    req.GetXJsonMask(),
		Time:    time.Now(),
	}
}

// CurrentTime ensures Time is set to the current time if not already set, and returns it
func (opts *UpdateOptions) CurrentTime() time.Time {
	ts := opts.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		opts.Time = ts
	}
	return ts
}
