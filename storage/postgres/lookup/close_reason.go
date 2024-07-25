package lookup

import (
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/storage"
)

type CloseReasonLookup struct {
	storage storage.Storage
}

func (c CloseReasonLookup) Create(rpc *model.CreateOptions, domainId int64, createdBy int64) error {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookup) Search(rpc *model.SearchOptions, ids []string) error {
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

func NewCloseReasonLookupStore(store storage.Storage) (storage.CloseReasonLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_close_reason_lookup.check.bad_arguments",
			"error creating config interface to the close_reason table, main store is nil")
	}
	return &CloseReasonLookup{storage: store}, nil
}
