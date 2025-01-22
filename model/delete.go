package model

import (
	"context"
	"errors"
	"github.com/webitel/cases/auth"
	"time"
)

type DeleteOptions struct {
	Time time.Time
	context.Context
	//Session *session.Session
	IDs      []int64
	ID       int64
	ParentID int64
	Auth     auth.Auther
}

func (s *DeleteOptions) SetAuthOpts(a auth.Auther) *DeleteOptions {
	s.Auth = a
	return s
}

func (s *DeleteOptions) GetAuthOpts() auth.Auther {
	return s.Auth
}

// CurrentTime sets and returns the current time if Time is zero.
func (rpc *DeleteOptions) CurrentTime() time.Time {
	ts := rpc.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		rpc.Time = ts
	}
	return ts
}

// NewDeleteOptions initializes a DeleteOptions instance with the current session, context, and current time.
func NewDeleteOptions(ctx context.Context, metadatter ObjectMetadatter) (*DeleteOptions, error) {

	deleteOpts := &DeleteOptions{
		Context: ctx,
	}
	deleteOpts.CurrentTime() // Set Time using CurrentTime

	if sess := GetAutherOutOfContext(ctx); sess != nil {
		deleteOpts.Auth = sess
	} else {
		return nil, errors.New("can't authorize user")
	}

	return deleteOpts, nil
}
