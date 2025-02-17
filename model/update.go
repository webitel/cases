package model

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/webitel/cases/auth"

	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

// UpdateOptions defines options for updating an entity with fields, mask, filter, and pagination
type UpdateOptions struct {
	Time time.Time
	context.Context
	// Session *session.Session
	// output
	Fields            []string
	UnknownFields     []string
	DerivedSearchOpts map[string]*SearchOptions
	// update
	Mask []string
	// filters
	ParentID int64
	IDs      []int64
	Etags    []*etag.Tid
	// ID      int64
	Auth auth.Auther
}

func (s *UpdateOptions) SetAuthOpts(a auth.Auther) *UpdateOptions {
	s.Auth = a
	return s
}

func (s *UpdateOptions) GetAuthOpts() auth.Auther {
	return s.Auth
}

type Updator interface {
	GetFields() []string
	GetXJsonMask() []string
}

// NewUpdateOptions initializes UpdateOptions with values from a context and an Updator-compliant struct
func NewUpdateOptions(ctx context.Context, req Updator, objMetadata ObjectMetadatter) (*UpdateOptions, error) {
	opts := &UpdateOptions{
		Context: ctx,
		Mask:    req.GetXJsonMask(),
	}
	opts.CurrentTime()
	if sess := GetAutherOutOfContext(ctx); sess != nil {
		opts.Auth = sess
	} else {
		return nil, errors.New("can't authorize user")
	}
	// Normalize fields
	var resultingFields []string
	if requestedFields := req.GetFields(); len(requestedFields) == 0 {
		resultingFields = make([]string, len(objMetadata.GetDefaultFields()))
		copy(resultingFields, objMetadata.GetDefaultFields())
	} else {
		resultingFields = util.FieldsFunc(
			requestedFields, graph.SplitFieldsQ,
		)
	}

	// Deduplicate and trim prefixes in the mask
	uniquePrefixes := make(map[string]struct{})
	var trimmedMask []string
	for _, field := range opts.Mask {
		prefix := field
		if dotIndex := strings.Index(field, "."); dotIndex > 0 {
			prefix = field[:dotIndex] // Trim after the dot
		}
		if _, exists := uniquePrefixes[prefix]; !exists {
			uniquePrefixes[prefix] = struct{}{}
			trimmedMask = append(trimmedMask, prefix)
		}
	}
	opts.Mask = trimmedMask

	// Split known and unknown fields
	resultingFields, opts.UnknownFields = util.SplitKnownAndUnknownFields(resultingFields, objMetadata.GetAllFields())
	opts.Fields = util.ParseFieldsForEtag(resultingFields)

	return opts, nil
}

// CurrentTime ensures Time is set to the current time if not already set, and returns it
func (opts *UpdateOptions) CurrentTime() time.Time {
	ts := opts.Time
	if ts.IsZero() {
		ts = time.Now().UTC()
		opts.Time = ts
	}
	return ts
}
