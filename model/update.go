package model

import (
	"context"
	"strings"
	"time"

	"github.com/webitel/cases/model/graph"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"

	session "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/server/interceptor"
)

// UpdateOptions defines options for updating an entity with fields, mask, filter, and pagination
type UpdateOptions struct {
	Time time.Time
	context.Context
	Session *session.Session
	// output
	Fields            []string
	UnknownFields     []string
	DerivedSearchOpts map[string]*SearchOptions
	// update
	Mask []string
	// filters
	IDs   []int64
	Etags []*etag.Tid
	// ID      int64
}

type Updator interface {
	GetFields() []string
	GetXJsonMask() []string
}

// NewUpdateOptions initializes UpdateOptions with values from a context and an Updator-compliant struct
func NewUpdateOptions(ctx context.Context, req Updator, objMetadata *ObjectMetadata) *UpdateOptions {
	opts := &UpdateOptions{
		Context: ctx,
		Session: ctx.Value(interceptor.SessionHeader).(*session.Session),
		Mask:    req.GetXJsonMask(),
		Time:    time.Now(),
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

	return opts
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
