package store

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	"github.com/jmoiron/sqlx"
	"github.com/webitel/cases/model"
)

type Store interface {
	Status() StatusStore
	StatusCondition() StatusConditionStore
	CloseReason() CloseReasonStore
	Appeal() AppealStore

	// Database connection
	Database() (*sqlx.DB, model.AppError)
	// Open database connection
	Open() model.AppError
	// Close database connection
	Close() model.AppError
}

type StatusStore interface {
	// Create a new status lookup
	Create(rpc *model.CreateOptions, add *_go.Status) (*_go.Status, error)
	// List status lookup
	List(rpc *model.SearchOptions) (*_go.StatusList, error)
	// Delete status lookup
	Delete(rpc *model.DeleteOptions) error
	// Update status lookup
	Update(rpc *model.UpdateOptions, lookup *_go.Status) (*_go.Status, error)
}

type StatusConditionStore interface {
	// Create a new status to a lookup
	Create(ctx *model.CreateOptions, add *_go.StatusCondition) (*_go.StatusCondition, error)
	// List lookup statuses
	List(ctx *model.SearchOptions) (*_go.StatusConditionList, error)
	// Delete lookup status
	Delete(ctx *model.DeleteOptions) error
	// Update lookup status
	Update(ctx *model.UpdateOptions, status *_go.StatusCondition) (*_go.StatusCondition, error)
}

type CloseReasonStore interface {
	// Create a new close reason lookup
	Create(rpc *model.CreateOptions, add *_go.CloseReason) (*_go.CloseReason, error)
	// List close reason lookup
	List(rpc *model.SearchOptions) (*_go.CloseReasonList, error)
	// Delete close reason lookup
	Delete(rpc *model.DeleteOptions) error
	// Update close reason lookup
	Update(rpc *model.UpdateOptions, lookup *_go.CloseReason) (*_go.CloseReason, error)
}

type AppealStore interface {
	// Create a new appeal lookup
	Create(rpc *model.CreateOptions, add *_go.Appeal) (*_go.Appeal, error)
	// List appeal lookup
	List(rpc *model.SearchOptions) (*_go.AppealList, error)
	// Delete appeal lookup
	Delete(rpc *model.DeleteOptions) error
	// Update appeal lookup
	Update(rpc *model.UpdateOptions, lookup *_go.Appeal) (*_go.Appeal, error)
}
