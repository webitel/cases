package model

import (
	"context"
	"time"

	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/server/interceptor"
)

type DeleteOptions struct {
	Time time.Time
	context.Context
	Session *session.Session
	IDs     []int64
	ID      int64
}

// CurrentTime sets and returns the current time if Time is zero.
func (rpc *DeleteOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}

// NewDeleteOptions initializes a DeleteOptions instance with the current session, context, and current time.
func NewDeleteOptions(ctx context.Context) *DeleteOptions {
	sess := ctx.Value(interceptor.SessionHeader).(*session.Session)

	deleteOpts := &DeleteOptions{
		Context: ctx,
		Session: sess,
	}
	deleteOpts.Time = deleteOpts.CurrentTime() // Set Time using CurrentTime

	return deleteOpts
}
