package model

import (
	"context"
	"time"

	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/server/interceptor"
)

type CreateOptions struct {
	Session         *session.Session
	context.Context // binding
	Time            time.Time
	Fields          []string
	Ids             []int64
	ID              int64
}

func (rpc *CreateOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}

func NewCreateOptions(ctx context.Context) *CreateOptions {
	sess := ctx.Value(interceptor.SessionHeader).(*session.Session)

	createOpts := &CreateOptions{
		Context: ctx,
		Session: sess,
	}
	createOpts.Time = createOpts.CurrentTime()
	return createOpts
}
