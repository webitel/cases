package model

import (
	"context"
	"time"
)

type DeleteOptions struct {
	Time time.Time
	context.Context
	//Session *session.Session
	IDs      []int64
	ID       int64
	ParentID int64
	Auth     Auther
}

func (s *DeleteOptions) SetAuthOpts(a Auther) *DeleteOptions {
	s.Auth = a
	return s
}

func (s *DeleteOptions) GetAuthOpts() Auther {
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
func NewDeleteOptions(ctx context.Context, metadatter ObjectMetadatter) *DeleteOptions {

	deleteOpts := &DeleteOptions{
		Context: ctx,
	}
	deleteOpts.CurrentTime() // Set Time using CurrentTime

	if sess := GetSessionOutOfContext(ctx); sess != nil {
		deleteOpts.Auth = NewSessionAuthOptions(sess, metadatter.GetAllScopeNames()...)
	}

	return deleteOpts
}
