package model

import (
	"context"
	"time"

	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"

	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/server/interceptor"
)

type CreateOptions struct {
	Session         *session.Session
	context.Context // binding
	Time            time.Time
	Fields          []string
	UnknownFields   []string
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

func NewCreateOptions(ctx context.Context, creator Creator, objMetadata ObjectMetadatter) *CreateOptions {
	createOpts := &CreateOptions{
		Context: ctx,
		Session: ctx.Value(interceptor.SessionHeader).(*session.Session),
	}

	// set current time
	createOpts.CurrentTime()

	// normalize fields
	var resultingFields []string
	if requestedFields := creator.GetFields(); len(requestedFields) == 0 {
		resultingFields = make([]string, len(objMetadata.GetDefaultFields()))
		copy(resultingFields, objMetadata.GetDefaultFields())
	} else {
		resultingFields = util.FieldsFunc(
			requestedFields, graph.SplitFieldsQ,
		)
	}
	resultingFields, createOpts.UnknownFields = util.SplitKnownAndUnknownFields(resultingFields, objMetadata.GetAllFields())
	createOpts.Fields = util.ParseFieldsForEtag(resultingFields)
	return createOpts
}
