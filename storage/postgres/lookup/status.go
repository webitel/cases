package lookup

import (
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/storage"
)

type StatusLookup struct {
	storage storage.Storage
}

func (s StatusLookup) Create(rpc *model.CreateOptions, domainId int64, createdBy int64) error {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) Search(rpc *model.SearchOptions, ids []string) error {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) Delete(rpc *model.DeleteOptions) error {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) Update(rpc *model.UpdateOptions) error {
	//TODO implement me
	panic("implement me")
}

func NewStatusLookupStore(store storage.Storage) (storage.StatusLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_status_lookup.check.bad_arguments",
			"error creating config interface to the status_lookup table, main store is nil")
	}
	return &StatusLookup{storage: store}, nil
}
