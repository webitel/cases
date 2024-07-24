package model

import (
	"context"
	session "github.com/webitel/cases/auth/model"
	"time"
)

type UpdateOptions struct {
	session         *session.Session
	context.Context //binding
	Time            time.Time
	ID              int64
	Fields          []string
}

func (rpc *UpdateOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now()
		rpc.Time = ts
	}
	return ts
}
