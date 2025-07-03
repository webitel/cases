package grpc

import (
	"context"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model/options"
	optsutil "github.com/webitel/cases/internal/model/options/grpc/util"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"google.golang.org/grpc/codes"
	"time"
)

type DeleteOption func(options *DeleteOptions) error

var _ options.Deleter = (*DeleteOptions)(nil)

type DeleteOptions struct {
	createdAt time.Time
	context.Context
	IDs      []int64
	ParentID int64
	Auth     auth.Auther
	Filters  map[string]any
}

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
		options.ParentID = id
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
	if sess := optsutil.GetAutherOutOfContext(ctx); sess != nil {
		deleteOpts.Auth = sess
	} else {
		return nil, errors.New("can't authorize user", errors.WithCode(codes.Unauthenticated))
	}
	for _, opt := range opts {
		err := opt(deleteOpts)
		if err != nil {
			return nil, err
		}
	}
	if len(deleteOpts.IDs) == 0 {
		return nil, errors.New("minimum one id required to delete", errors.WithCode(codes.InvalidArgument))
	}

	return deleteOpts, nil
}
