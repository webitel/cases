package store

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/cases/model"
)

type Store interface {
	StatusLookup() StatusLookupStore
	LookupStatus() LookupStatusStore
	CloseReasonLookup() CloseReasonLookupStore
	AppealLookup() AppealLookupStore

	// Database connection
	Database() (*sqlx.DB, model.AppError)
	// Open database connection
	Open() model.AppError
	// Close database connection
	Close() model.AppError
}

type StatusLookupStore interface {
	// Create a new status lookup
	Create(rpc *model.CreateOptions, add *_go.StatusLookup) (*_go.StatusLookup, error)
	// List status lookup
	List(rpc *model.SearchOptions) (*_go.StatusLookupList, error)
	// Delete status lookup
	Delete(rpc *model.DeleteOptions) error
	// Update status lookup
	Update(rpc *model.UpdateOptions, lookup *_go.StatusLookup) (*_go.StatusLookup, error)
}

type LookupStatusStore interface {
	// Create a new lookup status
	Attach(ctx *model.CreateOptions, add *_go.LookupStatus) (*_go.LookupStatus, error)
	// List lookup statuses
	List(ctx *model.SearchOptions) (*_go.LookupStatusList, error)
	// Delete lookup status
	Delete(ctx *model.DeleteOptions) error
	// Update lookup status
	Update(ctx *model.UpdateOptions, status *_go.LookupStatus) (*_go.LookupStatus, error)
}

type CloseReasonLookupStore interface {
	// Create a new close reason lookup
	Create(rpc *model.CreateOptions, add *_go.CloseReasonLookup) (*_go.CloseReasonLookup, error)
	// List close reason lookup
	List(rpc *model.SearchOptions) (*_go.CloseReasonLookupList, error)
	// Delete close reason lookup
	Delete(rpc *model.DeleteOptions) error
	// Update close reason lookup
	Update(rpc *model.UpdateOptions, lookup *_go.CloseReasonLookup) (*_go.CloseReasonLookup, error)
}

type AppealLookupStore interface {
	// Create a new appeal lookup
	Create(rpc *model.CreateOptions, add *_go.AppealLookup) (*_go.AppealLookup, error)
	// List appeal lookup
	List(rpc *model.SearchOptions) (*_go.AppealLookupList, error)
	// Delete appeal lookup
	Delete(rpc *model.DeleteOptions) error
	// Update appeal lookup
	Update(rpc *model.UpdateOptions, lookup *_go.AppealLookup) (*_go.AppealLookup, error)
}
