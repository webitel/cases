package postgres

import (
	"context"
	"fmt"
	"maps"
	"slices"

	common "github.com/webitel/cases/model/options"
	options "github.com/webitel/cases/model/options/grpc"
)

func withSearchOptions(ctx context.Context, opts ...options.SearchOption) common.SearchOptions {
	// base, _ := ctx.(common.SearchOptions)
	search, _ := ctx.(*options.SearchOptions)
	if search == nil {
		search, _ = options.NewSearchOptions(ctx,
			func(search *options.SearchOptions) error {

				if base, is := ctx.(common.SearchOptions); is {
					return fromSearchOptions(base)(search)
				}
				if base, is := ctx.(common.UpdateOptions); is {
					return fromUpdateOptions(base)(search)
				}
				if base, is := ctx.(common.CreateOptions); is {
					return fromCreateOptions(base)(search)
				}
				if base, is := ctx.(common.DeleteOptions); is {
					return fromDeleteOptions(base)(search)
				}

				return nil
			},
		)
	}
	var err error
	for _, setup := range opts {
		err = setup(search)
		if err != nil {
			panic(fmt.Errorf("options: %v", err))
		}
	}
	return search
}

func fromSearchOptions(opts common.SearchOptions) options.SearchOption {
	return func(search *options.SearchOptions) (_ error) {

		self, _ := opts.(*options.SearchOptions)
		if search == self {
			return // SELF
		}

		// search.Context = req.(context.Context)
		// search.createdAt = req.RequestTime()
		search.Auth = opts.GetAuthOpts()

		search.Fields = slices.Clone(opts.GetFields())
		search.UnknownFields = slices.Clone(opts.GetUnknownFields())

		search.IDs = slices.Clone(opts.GetIDs())
		search.Filters = maps.Clone(opts.GetFilters())

		search.Sort = opts.GetSort()
		search.Page = opts.GetPage()
		search.Size = opts.GetSize()

		return nil
	}
}

func fromCreateOptions(req common.CreateOptions) options.SearchOption {
	return func(search *options.SearchOptions) error {

		// search.Context = req.(context.Context)
		// search.createdAt = req.RequestTime()
		search.Auth = req.GetAuthOpts()

		search.Fields = slices.Clone(req.GetFields())
		search.UnknownFields = slices.Clone(req.GetUnknownFields())

		// GetDerivedSearchOpts() map[string]*SearchOptions
		// GetIDs() []int64
		// GetParentID() int64
		// GetChildID() int64

		return nil
	}
}

func fromUpdateOptions(req common.UpdateOptions) options.SearchOption {
	return func(search *options.SearchOptions) error {

		// search.Context = req.(context.Context)
		// search.createdAt = req.RequestTime()
		search.Auth = req.GetAuthOpts()

		search.Fields = slices.Clone(req.GetFields())
		search.UnknownFields = slices.Clone(req.GetUnknownFields())

		// GetMask() []string
		// GetEtags() []*etag.Tid
		// GetParentID() int64
		// GetIDs() []int64

		return nil
	}
}

func fromDeleteOptions(req common.DeleteOptions) options.SearchOption {
	return func(search *options.SearchOptions) error {

		// search.Context = req.(context.Context)
		// search.createdAt = req.RequestTime()
		search.Auth = req.GetAuthOpts()

		search.Filters = maps.Clone(req.GetFilters())

		// RemoveFilter(string)
		// AddFilter(string, any)
		// GetFilter(string) any

		// // If connection to parent object required
		// GetParentID() int64

		// // ID filtering
		// GetIDs() []int64

		return nil
	}
}
