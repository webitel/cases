package postgres

import (
	"context"

	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
)

type LinkCase struct {
	storage store.Store
}

// Create implements store.LinkCaseStore.
func (l *LinkCase) Create(ctx context.Context, rpc *model.CreateOptions, add *_go.CaseLink) (*_go.CaseLink, error) {
	panic("unimplemented")
}

// Delete implements store.LinkCaseStore.
func (l *LinkCase) Delete(ctx context.Context, req *model.DeleteOptions) (*_go.CaseLink, error) {
	panic("unimplemented")
}

// List implements store.LinkCaseStore.
func (l *LinkCase) List(ctx context.Context, rpc *model.SearchOptions) (*_go.CaseLinkList, error) {
	panic("unimplemented")
}

// Merge implements store.LinkCaseStore.
func (l *LinkCase) Merge(ctx context.Context, req *model.UpdateOptions) (*_go.CaseLinkList, error) {
	panic("unimplemented")
}

// Update implements store.LinkCaseStore.
func (l *LinkCase) Update(ctx context.Context, req *model.UpdateOptions) (*_go.CaseLink, error) {
	panic("unimplemented")
}

func NewLinkCaseStore(store store.Store) (store.LinkCaseStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_link_case.check.bad_arguments",
			"error creating link case interface to the comment_case table, main store is nil")
	}
	return &LinkCase{storage: store}, nil
}
