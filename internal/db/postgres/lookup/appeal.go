package lookup

import (
	db "github.com/webitel/cases/internal/db"
	"github.com/webitel/cases/model"
)

type AppealLookup struct {
	storage db.DB
}

func (a AppealLookup) Create(rpc *model.CreateOptions) error {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookup) List(rpc *model.SearchOptions, ids []string) error {
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

func NewAppealLookupStore(store db.DB) (db.AppealLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_appeal_lookup.check.bad_arguments",
			"error creating config interface to the appeal table, main store is nil")
	}
	return &AppealLookup{storage: store}, nil
}
