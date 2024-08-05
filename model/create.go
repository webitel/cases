package model

import (
	"context"
	"time"

	session "github.com/webitel/cases/auth/model"
)

type CreateOptions struct {
	Session         *session.Session
	context.Context // binding
	Time            time.Time
	Fields          []string
}

func (rpc *CreateOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now()
		rpc.Time = ts
	}
	return ts
}
