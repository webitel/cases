package model

import (
	"context"
	"time"

	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/util"
)

type DeleteOptions struct {
	FieldsUtil util.FieldsUtils
	Time       time.Time
	context.Context
	Session *session.Session
	IDs     []int64
}

func (rpc *DeleteOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}
