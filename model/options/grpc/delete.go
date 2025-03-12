package grpc

import (
	"context"
	"errors"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"time"
)

type DeleteOptions struct {
	createdAt time.Time
	context.Context
	IDs      []int64
	ParentID int64
	Auth     auth.Auther
	Filters  map[string]any
}

type DeleteOption func(options *DeleteOptions) error

func WithDeleteIDs(ids []int64) DeleteOption {
	return func(options *DeleteOptions) error {
		options.IDs = ids
		return nil
	}
}
func WithDeleteID(id int64) DeleteOption {
	return func(options *DeleteOptions) error {
		options.IDs = []int64{id}
		return nil
	}
}

func WithDeleteIDsAsEtags(tag etag.EtagType, etags ...string) DeleteOption {
	return func(options *DeleteOptions) error {
		ids, err := util.ParseIds(etags, tag)
		if err != nil {
			return err
		}
		options.IDs = ids
		return nil
	}
}
func WithDeleteParentID(id int64) DeleteOption {
	return func(options *DeleteOptions) error {
		options.IDs = []int64{id}
		return nil
	}
}

func WithDeleteParentIDAsEtag(etagType etag.EtagType, tag string) DeleteOption {
	return func(options *DeleteOptions) error {
		id, err := etag.EtagOrId(etagType, tag)
		if err != nil {
			return err
		}
		options.ParentID = id.GetOid()
		return nil
	}
}

func (s *DeleteOptions) RequestTime() time.Time {
	return s.createdAt
}

func (s *DeleteOptions) RemoveFilter(key string) {
	delete(s.Filters, key)
}

func (s *DeleteOptions) AddFilter(key string, value any) {
	s.Filters[key] = value
}

func (s *DeleteOptions) GetFilter(key string) any {
	return s.Filters[key]
}
func (s *DeleteOptions) GetFilters() map[string]any {
	return s.Filters
}

func (s *DeleteOptions) GetParentID() int64 {
	return s.ParentID
}

func (s *DeleteOptions) GetIDs() []int64 {
	return s.IDs
}

func (s *DeleteOptions) GetAuthOpts() auth.Auther {
	return s.Auth
}

// NewDeleteOptions initializes a DeleteOptions instance with the current session, context, and current time.
func NewDeleteOptions(ctx context.Context, opts ...DeleteOption) (*DeleteOptions, error) {

	deleteOpts := &DeleteOptions{
		Context:   ctx,
		Filters:   map[string]any{},
		createdAt: time.Now().UTC(),
	}
	if sess := model.GetAutherOutOfContext(ctx); sess != nil {
		deleteOpts.Auth = sess
	} else {
		return nil, errors.New("can't authorize user")
	}
	for _, opt := range opts {
		err := opt(deleteOpts)
		if err != nil {
			return nil, err
		}
	}
	if len(deleteOpts.IDs) == 0 {
		return nil, errors.New("minimum one id required to delete")
	}

	return deleteOpts, nil
}
