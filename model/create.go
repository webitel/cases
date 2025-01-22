package model

import (
	"context"
	"errors"
	"github.com/webitel/cases/auth"
	"time"

	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"
)

type CreateOptions struct {
	//Session         *session.Session
	context.Context // binding
	Time            time.Time
	// output
	Fields            []string
	DerivedSearchOpts map[string]*SearchOptions
	UnknownFields     []string
	// filters
	Ids []int64
	// ParentID is the attribute to represent parent object, that creation process connected to
	ParentID int64
	// ChildID is the attribute to represent child object, that creation process connect
	ChildID int64
	Auth    auth.Auther
}

func (s *CreateOptions) SetAuthOpts(a auth.Auther) *CreateOptions {
	s.Auth = a
	return s
}

func (s *CreateOptions) GetAuthOpts() auth.Auther {
	return s.Auth
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

func NewCreateOptions(ctx context.Context, creator Creator, objMetadata ObjectMetadatter) (*CreateOptions, error) {
	createOpts := &CreateOptions{
		Context: ctx,
	}

	// set current time
	createOpts.CurrentTime()
	if sess := GetAutherOutOfContext(ctx); sess != nil {
		createOpts.Auth = sess
	} else {
		return nil, errors.New("can't authorize user")
	}

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
	return createOpts, nil
}
