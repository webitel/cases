package postgres

import (
	"context"
	"fmt"

	"github.com/lib/pq"
	db "github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
)

type AccessContol struct {
	storage db.Store
}

// Check RBAC rights for the user to access the resource
func (s AccessContol) RbacAccess(ctx context.Context, domainId int64, id int64, groups []int, access uint8, table string) (bool, model.AppError) {
	// Get the database connection from the storage layer
	db, appErr := s.storage.Database()
	if appErr != nil {
		return false, appErr
	}

	// Append "_acl" to the base table name to create the full table name
	tableName := table + "_acl"

	// Format the SQL query string
	sql := fmt.Sprintf(`
		SELECT 1
		FROM %s acl
		WHERE acl.dc = $1
		  AND acl.object = $2
		  AND acl.subject = ANY($3::int[])
		  AND acl.access & $4 = $4
		LIMIT 1`, tableName)

	// Execute the query
	var ac bool
	err := db.QueryRow(ctx, sql, domainId, id, pq.Array(groups), access).Scan(&ac)
	if err != nil {
		return false, model.NewInternalError("postgres.access_control.check_access.scan.error", err.Error())
	}

	return ac, nil
}

func NewAccessControlStore(store db.Store) (db.AccessControlStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_access_control.check.bad_arguments",
			"error creating access control interface, main store is nil")
	}
	return &AccessContol{storage: store}, nil
}
