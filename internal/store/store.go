package store

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
)

// Store is an interface that defines all the methods and properties that a store should implement in Cases service
type Store interface {
	Status() StatusStore
	StatusCondition() StatusConditionStore
	CloseReason() CloseReasonStore
	Reason() ReasonStore
	Appeal() AppealStore
	AccessControl() AccessControlStore

	// Database connection
	Database() (*pgxpool.Pool, model.AppError)
	// Open database connection
	Open() model.AppError
	// Close database connection
	Close() model.AppError
}

type AccessControlStore interface {
	// Check if user has Rbac access
	RbacAccess(ctx context.Context, domainId, id int64, groups []int, access uint8, table string) (bool, model.AppError)
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
	// Create a new status сondition
	Create(ctx *model.CreateOptions, add *_go.StatusCondition) (*_go.StatusCondition, error)
	// List status сondition
	List(ctx *model.SearchOptions, statusId int64) (*_go.StatusConditionList, error)
	// Delete status сondition
	Delete(ctx *model.DeleteOptions, statusId int64) error
	// Update status сondition
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

type ReasonStore interface {
	// Create a new reason
	Create(ctx *model.CreateOptions, add *_go.Reason) (*_go.Reason, error)
	// List reasons
	List(ctx *model.SearchOptions, closeReasonId int64) (*_go.ReasonList, error)
	// Delete reason
	Delete(ctx *model.DeleteOptions, closeReasonId int64) error
	// Update reason
	Update(ctx *model.UpdateOptions, lookup *_go.Reason) (*_go.Reason, error)
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
