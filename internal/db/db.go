package db

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/cases/model"
)

type DB interface {
	Status() StatusLookupStore
	CloseReason() CloseReasonLookupStore
	Appeal() AppealLookupStore

	// Database connection
	Database() (*sqlx.DB, model.AppError)
	Open() model.AppError
	Close() model.AppError
}

type StatusLookupStore interface {
	// Create a new status lookup
	Create(rpc *model.CreateOptions, add *_go.StatusLookup) (*_go.StatusLookup, error)
	// List status lookup
	List(rpc *model.SearchOptions) (*_go.StatusLookupList, error)
	// Delete status lookup
	Delete(rpc *model.DeleteOptions) error
	// Update status lookupßß
	Update(rpc *model.UpdateOptions, lookup *_go.StatusLookup) (*_go.StatusLookup, error)
}

type CloseReasonLookupStore interface {
	// Create a new close reason lookup
	Create(rpc *model.CreateOptions) error
	// List close reason lookup
	List(rpc *model.SearchOptions, ids []string) error
	// Delete close reason lookup
	Delete(rpc *model.DeleteOptions) error
	// Update close reason lookup
	Update(rpc *model.UpdateOptions) error
}

type AppealLookupStore interface {
	// Create a new appeal lookup
	Create(rpc *model.CreateOptions) error
	// List appeal lookup
	List(rpc *model.SearchOptions, ids []string) error
	// Delete appeal lookup
	Delete(rpc *model.DeleteOptions) error
	// Update appeal lookup
	Update(rpc *model.UpdateOptions) error
}
