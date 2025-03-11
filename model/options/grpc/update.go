package grpc

import (
	"context"
	"errors"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/model/options"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"time"
)

type UpdateOption func(*UpdateOptions) error

type UpdateMasker interface {
	GetXJsonMask() []string
}

func WithUpdateFields(
	fielder Fielder,
	md model.ObjectMetadatter,
	fieldsModifiers ...func(fields []string) []string,
) UpdateOption {
	return func(o *UpdateOptions) error {
		if requestedFields := fielder.GetFields(); len(requestedFields) == 0 {
			o.Fields = md.GetDefaultFields()
		} else {
			o.Fields = util.DeduplicateFields(util.FieldsFunc(
				requestedFields, util.InlineFields,
			))
		}
		o.Fields, o.UnknownFields = util.SplitKnownAndUnknownFields(o.Fields, md.GetAllFields())
		o.Fields = util.ParseFieldsForEtag(o.Fields)
		return nil
	}
}

func WithUpdateMasker(m UpdateMasker) UpdateOption {
	return func(o *UpdateOptions) error {
		o.Mask = append(o.Mask, m.GetXJsonMask()...)
		return nil
	}
}

// WithUpdateEtag adds an etag to the UpdateOptions
func WithUpdateEtag(etags ...*etag.Tid) UpdateOption {
	return func(o *UpdateOptions) error {
		o.Etags = append(o.Etags, etags...)
		return nil
	}
}

func WithUpdateParentID(parentID int64) UpdateOption {
	return func(o *UpdateOptions) error {
		o.ParentID = parentID
		return nil
	}
}

func WithUpdateIDs(ids []int64) UpdateOption {
	return func(o *UpdateOptions) error {
		o.IDs = ids
		return nil
	}
}

type UpdateOptions struct {
	context.Context
	Time              time.Time
	Fields            []string
	UnknownFields     []string
	DerivedSearchOpts map[string]*options.SearchOptions
	Mask              []string
	Etags             []*etag.Tid
	Auth              auth.Auther
	ParentID          int64
	IDs               []int64
}

func (s *UpdateOptions) GetAuthOpts() auth.Auther {
	return s.Auth
}
func (s *UpdateOptions) SetAuthOpts(auth auth.Auther) *UpdateOptions {
	s.Auth = auth
	return s
}
func (s *UpdateOptions) GetIDs() []int64            { return s.IDs }
func (s *UpdateOptions) GetParentID() int64         { return s.ParentID }
func (s *UpdateOptions) GetFields() []string        { return s.Fields }
func (s *UpdateOptions) GetUnknownFields() []string { return s.UnknownFields }
func (s *UpdateOptions) GetDerivedSearchOpts() map[string]*options.SearchOptions {
	return s.DerivedSearchOpts
}
func (s *UpdateOptions) SetDerivedSearchOpts(opts map[string]*options.SearchOptions) *UpdateOptions {
	s.DerivedSearchOpts = opts
	return s
}
func (s *UpdateOptions) GetMask() []string     { return s.Mask }
func (s *UpdateOptions) GetEtags() []*etag.Tid { return s.Etags }
func (s *UpdateOptions) SetEtags(etags ...*etag.Tid) *UpdateOptions {
	s.Etags = append(s.Etags, etags...)
	return s
}
func (s *UpdateOptions) GetTime() time.Time { return s.Time }

func NewUpdateOptions(ctx context.Context, opts ...UpdateOption) (*UpdateOptions, error) {
	updateOpts := &UpdateOptions{
		Context: ctx,
		Time:    time.Now().UTC(),
	}

	// Apply functional updateOpts
	for _, opt := range opts {
		if err := opt(updateOpts); err != nil {
			return nil, err
		}
	}

	// Set authentication
	if err := setUpdateAuthOptions(ctx, updateOpts); err != nil {
		return nil, err
	}

	// Deduplicate and trim mask prefixes
	updateOpts.Mask = DeduplicateMaskPrefixes(updateOpts.Mask)

	return updateOpts, nil
}

// setUpdateAuthOptions extracts authentication from context and sets it in options
func setUpdateAuthOptions(ctx context.Context, options *UpdateOptions) error {
	if sess := model.GetAutherOutOfContext(ctx); sess != nil {
		options.Auth = sess
		return nil
	}
	return errors.New("can't authorize user")
}
