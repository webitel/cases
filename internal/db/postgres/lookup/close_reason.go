package lookup

import (
	db "github.com/webitel/cases/internal/db"
	"github.com/webitel/cases/model"
)

type CloseReasonLookup struct {
	storage db.DB
}

func (c CloseReasonLookup) Create(rpc *model.CreateOptions) error {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookup) List(rpc *model.SearchOptions, ids []string) error {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookup) Delete(rpc *model.DeleteOptions) error {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookup) Update(rpc *model.UpdateOptions) error {
	//TODO implement me
	panic("implement me")
}

func NewCloseReasonLookupStore(store db.DB) (db.CloseReasonLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_close_reason_lookup.check.bad_arguments",
			"error creating config interface to the close_reason table, main store is nil")
	}
	return &CloseReasonLookup{storage: store}, nil
}
