package lookup

import (
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/storage"
)

type AppealLookup struct {
	storage storage.Storage
}

func (a AppealLookup) Create(rpc *model.CreateOptions, domainId int64, createdBy int64) error {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookup) Search(rpc *model.SearchOptions, ids []string) error {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookup) Delete(rpc *model.DeleteOptions) error {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookup) Update(rpc *model.UpdateOptions) error {
	//TODO implement me
	panic("implement me")
}

func NewAppealLookupStore(store storage.Storage) (storage.AppealLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_appeal_lookup.check.bad_arguments",
			"error creating config interface to the appeal table, main store is nil")
	}
	return &AppealLookup{storage: store}, nil
}
