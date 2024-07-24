package storage

import (
	"github.com/jmoiron/sqlx"
	model "github.com/webitel/cases/model"
)

type Storage interface {
	//Status lookup
	Status() StatusLookupStore
	// CloseReason  lookup
	CloseReason() CloseReasonLookupStore
	// Appeal lookup
	Appeal() AppealLookupStore

	// Database connection
	Database() (*sqlx.DB, model.AppError)

	Open() model.AppError
	Close() model.AppError
}

type StatusLookupStore interface {
	Create(rpc *model.CreateOptions, domainId int64, createdBy int64) error
	Search(rpc *model.SearchOptions, ids []string) error
	Delete(rpc *model.DeleteOptions) error
	Update(rpc *model.UpdateOptions) error
}

type CloseReasonLookupStore interface {
	Create(rpc *model.CreateOptions, domainId int64, createdBy int64) error
	Search(rpc *model.SearchOptions, ids []string) error
	Delete(rpc *model.DeleteOptions) error
	Update(rpc *model.UpdateOptions) error
}

type AppealLookupStore interface {
	Create(rpc *model.CreateOptions, domainId int64, createdBy int64) error
	Search(rpc *model.SearchOptions, ids []string) error
	Delete(rpc *model.DeleteOptions) error
	Update(rpc *model.UpdateOptions) error
}
