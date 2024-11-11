package model

import (
	"context"
	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"
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
	// ParentID is the attribute to represent parent object, that creation process connected to
	ParentID int64
}

type Creator interface {
	GetFields() []string
}

func (rpc *CreateOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}

func NewCreateOptions(ctx context.Context, creator Creator) *CreateOptions {
	sess := ctx.Value(interceptor.SessionHeader).(*session.Session)

	createOpts := &CreateOptions{
		Context: ctx,
		Session: sess,
		Fields: util.FieldsFunc(
			creator.GetFields(), graph.SplitFieldsQ,
		),
	}
	createOpts.Time = createOpts.CurrentTime()
	return createOpts
}
