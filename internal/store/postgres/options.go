package postgres

import (
	"context"
	"fmt"
	"slices"

	common "github.com/webitel/cases/internal/model/options"
	options "github.com/webitel/cases/internal/model/options/grpc"
)

func withSearchOptions(ctx context.Context, opts ...options.SearchOption) common.Searcher {
	// base, _ := ctx.(common.Searcher)
	search, _ := ctx.(*options.SearchOptions)
	if search == nil {
		search, _ = options.NewSearchOptions(ctx,
			func(search *options.SearchOptions) error {

				if base, is := ctx.(common.Searcher); is {
					return fromSearchOptions(base)(search)
				}
				if base, is := ctx.(common.Updator); is {
					return fromUpdateOptions(base)(search)
				}
				if base, is := ctx.(common.Creator); is {
					return fromCreateOptions(base)(search)
				}
				if base, is := ctx.(common.Deleter); is {
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

func fromSearchOptions(opts common.Searcher) options.SearchOption {
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
		search.Filters = slices.Clone(opts.GetFilters())

		search.Sort = opts.GetSort()
		search.Page = opts.GetPage()
		search.Size = opts.GetSize()

		return nil
	}
}

func fromCreateOptions(req common.Creator) options.SearchOption {
	return func(search *options.SearchOptions) error {

		// search.Context = req.(context.Context)
		// search.createdAt = req.RequestTime()
		search.Auth = req.GetAuthOpts()

		search.Fields = slices.Clone(req.GetFields())
		search.UnknownFields = slices.Clone(req.GetUnknownFields())

		// GetDerivedSearchOpts() map[string]*Searcher
		// GetIDs() []int64
		// GetParentID() int64
		// GetChildID() int64

		return nil
	}
}

func fromUpdateOptions(req common.Updator) options.SearchOption {
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

func fromDeleteOptions(req common.Deleter) options.SearchOption {
	return func(search *options.SearchOptions) error {

		// search.Context = req.(context.Context)
		// search.createdAt = req.RequestTime()
		search.Auth = req.GetAuthOpts()

		search.Filters = slices.Clone(req.GetFilters())

		// RemoveFilter(string)
		// AddFilter(string, any)
		// GetFilter1(string) any

		// // If connection to parent object required
		// GetParentID() int64

		// // ID filtering
		// GetIDs() []int64

		return nil
	}
}
