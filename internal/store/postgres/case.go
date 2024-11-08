package postgres

import (
	"context"

	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
)

type Case struct {
	storage store.Store
}

// Create implements store.CaseStore.
func (c *Case) Create(ctx context.Context, rpc *model.CreateOptions, add *_go.Case) (*_go.Case, error) {
	panic("unimplemented")
}

// Delete implements store.CaseStore.
func (c *Case) Delete(ctx context.Context, req *model.DeleteOptions) (*_go.Case, error) {
	panic("unimplemented")
}

// List implements store.CaseStore.
func (c *Case) List(ctx context.Context, rpc *model.SearchOptions) (*_go.CaseList, error) {
	panic("unimplemented")
}

// Merge implements store.CaseStore.
func (c *Case) Merge(ctx context.Context, req *model.CreateOptions) (*_go.CaseList, error) {
	panic("unimplemented")
}

// Update implements store.CaseStore.
func (c *Case) Update(ctx context.Context, req *model.UpdateOptions) (*_go.Case, error) {
	panic("unimplemented")
}

func NewCaseStore(store store.Store) (store.CaseStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case.check.bad_arguments",
			"error creating case interface to the case table, main store is nil")
	}
	return &Case{storage: store}, nil
}
