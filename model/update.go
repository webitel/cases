package model

import (
	"context"
	"time"

	session "github.com/webitel/cases/auth/model"
)

type UpdateOptions struct {
	Time time.Time
	context.Context
	Session *session.Session
	Fields  []string
	IDs     []int64
	ID      int64
}

func (rpc *UpdateOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}
