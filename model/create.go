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
	hasEtag  bool
	hasId    bool
	hasVer   bool
}

func (s *CreateOptions) HasEtag() bool {
	return s.hasEtag
}
func (s *CreateOptions) HasId() bool {
	return s.hasId
}
func (s *CreateOptions) HasVer() bool {
	return s.hasVer
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

func NewCreateOptions(ctx context.Context, creator Creator, defaultFields []string) *CreateOptions {
	createOpts := &CreateOptions{
		Context: ctx,
		Session: ctx.Value(interceptor.SessionHeader).(*session.Session),
	}

	// set current time
	createOpts.CurrentTime()

	// normalize fields
	fields := util.FieldsFunc(
		creator.GetFields(), graph.SplitFieldsQ,
	)
	if len(fields) == 0 {
		fields = defaultFields
	}
	createOpts.Fields, createOpts.hasEtag, createOpts.hasId, createOpts.hasVer = util.ProcessEtag(fields)
	return createOpts
}
