package lookup

import (
	"github.com/webitel/cases/db"
	db "github.com/webitel/cases/internal/db"
	"github.com/webitel/cases/model"
)

type AppealLookup struct {
	storage db.Storage
}

func (a AppealLookup) Create(rpc *model.CreateOptions) error {
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

func NewAppealLookupStore(store db.Storage) (db2.AppealLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_appeal_lookup.check.bad_arguments",
			"error creating config interface to the appeal table, main store is nil")
	}
	return &AppealLookup{storage: store}, nil
}
