package model

import (
	"context"
	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/util"
	"time"
)

type DeleteOptions struct {
	Session         *session.Session
	context.Context //binding
	Time            time.Time
	IDs             []int64
	FieldsUtil      util.FieldsUtils
}

func (rpc *DeleteOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now()
		rpc.Time = ts
	}
	return ts
}
