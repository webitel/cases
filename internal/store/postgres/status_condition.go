package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	_go "github.com/webitel/cases/api"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
)

type StatusConditionStore struct {
	storage store.Store
}

func (s StatusConditionStore) Create(ctx *model.CreateOptions, add *_go.StatusCondition) (*_go.StatusCondition, error) {
	db, err := s.getDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return nil, err
	}

	tx, err := db.BeginTx(ctx.Context, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer s.handleTx(ctx.Context, tx, &err)

	query, args, err := s.buildCreateStatusConditionQuery(ctx, add)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	var createdBy, updatedBy _go.Lookup
	var createdAt, updatedAt time.Time

	err = tx.QueryRow(ctx.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &updatedAt, &add.Description, &add.Initial, &add.Final,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, &add.StatusId,
	)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	add.CreatedAt = createdAt.Unix()
	add.UpdatedAt = updatedAt.Unix()
	add.CreatedBy = &createdBy
	add.UpdatedBy = &updatedBy

	return add, nil
}

func (s StatusConditionStore) List(ctx *model.SearchOptions, statusId int64) (*_go.StatusConditionList, error) {
	db, err := s.getDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return nil, err
	}

	queryBuilder, err := s.buildListStatusConditionQuery(ctx, statusId)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		log.Printf("Failed to generate SQL query: %v", err)
		return nil, err
	}

	rows, err := db.Query(ctx.Context, query, args...)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var statusList []*_go.StatusCondition
	lCount := 0
	next := false
	for rows.Next() {
		if lCount >= ctx.GetSize() {
			next = true
			break
		}

		st := &_go.StatusCondition{}
		var createdBy, updatedBy _go.Lookup
		var tempCreatedAt, tempUpdatedAt time.Time

		scanArgs := s.buildScanArgs(ctx.Fields, st, &createdBy, &updatedBy, &tempCreatedAt, &tempUpdatedAt)
		if err := rows.Scan(scanArgs...); err != nil {
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		s.populateStatusConditionFields(ctx.Fields, st, &createdBy, &updatedBy, tempCreatedAt, tempUpdatedAt)
		statusList = append(statusList, st)
		lCount++
	}

	return &_go.StatusConditionList{
		Page:  int32(ctx.Page),
		Next:  next,
		Items: statusList,
	}, nil
}

func (s StatusConditionStore) Delete(ctx *model.DeleteOptions, statusId int64) error {
	// Prepare the arguments for the query
	ids := []int64{statusId}
	domainId := ctx.Session.GetDomainId()

	// Build the delete query with constraints checking
	query, args, err := s.buildDeleteStatusConditionQuery(ids, domainId, statusId)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return err
	}

	// Get the database connection
	db, err := s.getDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return err
	}

	// Execute the query
	rows, err := db.Query(ctx.Context, query, args...)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return err
	}
	defer rows.Close()

	// Process the results
	var deletedIds []int64
	for rows.Next() {
		var deletedId int64
		if err := rows.Scan(&deletedId); err != nil {
			log.Printf("Failed to scan deleted ID: %v", err)
			return err
		}
		deletedIds = append(deletedIds, deletedId)
	}

	// Check if any records were deleted
	if len(deletedIds) == 0 {
		return errors.New("operation would violate constraints: at least one initial and one final record must remain")
	}

	return nil
}

func (s StatusConditionStore) Update(ctx *model.UpdateOptions, st *_go.StatusCondition) (*_go.StatusCondition, error) {
	db, err := s.getDBConnection()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return nil, err
	}

	tx, err := db.BeginTx(ctx.Context, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer s.handleTx(ctx.Context, tx, &err)

	for _, field := range ctx.Fields {
		switch field {
		case "initial":
			if !st.Initial {
				log.Printf("Not allowed to update: initial = FALSE")
				return nil, fmt.Errorf("update not allowed: there must be at least one initial = TRUE for the given dc and status_id")
			}
		}
	}

	// for _, field := range ctx.Fields {
	// 	switch field {
	// 	case "initial":
	// 		if !st.Initial {
	// 			// Check if it's the last initial status condition
	// 			isLast, initialErr := s.isLastStatusCondition(ctx.Context, tx, ctx.Session.GetDomainId(), st.StatusId, st.Id, "initial")
	// 			if initialErr != nil {
	// 				log.Printf("Failed to check initial status condition: %v", initialErr)
	// 				return nil, initialErr
	// 			}
	// 			if isLast {
	// 				return nil, fmt.Errorf("update not allowed: there must be at least one initial = TRUE for the given dc and status_id")
	// 			}
	// 		}
	// 	case "final":
	// 		if !st.Final {
	// 			// Check if it's the last final status condition
	// 			isLast, finalErr := s.isLastStatusCondition(ctx.Context, tx, ctx.Session.GetDomainId(), st.StatusId, st.Id, "final")
	// 			if finalErr != nil {
	// 				log.Printf("Failed to check final status condition: %v", finalErr)
	// 				return nil, finalErr
	// 			}
	// 			if isLast {
	// 				return nil, fmt.Errorf("update not allowed: there must be at least one final = TRUE for the given dc and status_id")
	// 			}
	// 		}
	// 	}
	// }

	// Build the update query
	query, args := s.buildUpdateStatusConditionQuery(ctx, st)

	// Log the final query and arguments
	fmt.Printf("Final query: %s\n", query)
	fmt.Printf("Final args: %v\n", args)

	var createdBy, updatedBy _go.Lookup
	var createdAt, updatedAt time.Time

	err = tx.QueryRow(ctx.Context, query, args...).Scan(
		&st.Id, &st.Name, &createdAt, &updatedAt, &st.Description, &st.Initial, &st.Final,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, &st.StatusId,
	)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, fmt.Errorf("update not allowed: there must be at least one initial or final = TRUE for the given dc and status_id")
	}

	st.CreatedAt = createdAt.Unix()
	st.UpdatedAt = updatedAt.Unix()
	st.CreatedBy = &createdBy
	st.UpdatedBy = &updatedBy

	return st, nil
}

func (s StatusConditionStore) buildCreateStatusConditionQuery(ctx *model.CreateOptions, status *_go.StatusCondition) (string, []interface{}, error) {
	query := `
WITH existing_status AS (
    SELECT COUNT(*) AS count
    FROM cases.status_condition
    WHERE dc = $7 AND status_id = $8
),
default_values AS (
    SELECT
        CASE WHEN (SELECT count FROM existing_status) = 0 THEN TRUE ELSE FALSE END AS initial_default,
        CASE WHEN (SELECT count FROM existing_status) = 0 THEN TRUE ELSE FALSE END AS final_default
),
ins AS (
    INSERT INTO cases.status_condition (name, created_at, description, initial, final, created_by, updated_at, updated_by, dc, status_id)
    VALUES (
        $1, $2, $3,
        (SELECT initial_default FROM default_values),
        (SELECT final_default FROM default_values),
        $4, $5, $6, $7, $8
    )
    RETURNING id, name, created_at, updated_at, description, initial, final, created_by, updated_by, status_id
)
SELECT
    ins.id,
    ins.name,
    ins.created_at,
    ins.updated_at,
    ins.description,
    ins.initial,
    ins.final,
    ins.created_by AS created_by_id,
    COALESCE(c.name::text, c.username) AS created_by_name,
    ins.updated_by AS updated_by_id,
    COALESCE(u.name::text, u.username) AS updated_by_name,
    ins.status_id
FROM ins
LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;
`
	args := []interface{}{
		status.Name, ctx.Time, status.Description,
		ctx.Session.GetUserId(), ctx.Time, ctx.Session.GetUserId(), ctx.Session.GetDomainId(), status.StatusId,
	}
	return query, args, nil
}

func (s StatusConditionStore) buildListStatusConditionQuery(ctx *model.SearchOptions, statusId int64) (squirrel.SelectBuilder, error) {
	queryBuilder := squirrel.Select().
		From("cases.status_condition AS s").
		Where(squirrel.Eq{"s.dc": ctx.Session.GetDomainId(), "s.status_id": statusId}).
		PlaceholderFormat(squirrel.Dollar)

	fields := ctx.FieldsUtil.FieldsFunc(ctx.Fields, ctx.FieldsUtil.InlineFields)

	ctx.Fields = append(fields, "id")

	for _, field := range ctx.Fields {
		switch field {
		case "id", "name", "description", "initial", "final", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("s." + field)
		case "created_by":
			queryBuilder = queryBuilder.Column("created_by.id AS created_by_id, created_by.name AS created_by_name").
				LeftJoin("directory.wbt_auth AS created_by ON s.created_by = created_by.id")
		case "updated_by":
			queryBuilder = queryBuilder.Column("updated_by.id AS updated_by_id, updated_by.name AS updated_by_name").
				LeftJoin("directory.wbt_auth AS updated_by ON s.updated_by = updated_by.id")
		}
	}

	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)

	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	if len(ctx.IDs) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"s.id": ids})
	}

	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substrs := ctx.Match.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(squirrel.ILike{"s.name": "%" + combinedLike + "%"})
	}

	parsedFields := ctx.FieldsUtil.FieldsFunc(ctx.Sort, ctx.FieldsUtil.InlineFields)

	var sortFields []string

	for _, sortField := range parsedFields {
		desc := false
		if strings.HasPrefix(sortField, "!") {
			desc = true
			sortField = strings.TrimPrefix(sortField, "!")
		}

		var column string
		switch sortField {
		case "name", "description":
			column = "s." + sortField
		default:
			continue
		}

		if desc {
			column += " DESC"
		} else {
			column += " ASC"
		}

		sortFields = append(sortFields, column)
	}

	size := ctx.GetSize()
	queryBuilder = queryBuilder.OrderBy(sortFields...).Offset(uint64((ctx.Page - 1) * size))
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size))
	}

	return queryBuilder, nil
}

func (s StatusConditionStore) buildDeleteStatusConditionQuery(ids []int64, domainId, statusId int64) (string, []interface{}, error) {
	query := `
WITH
    to_check AS (
        SELECT id, initial, final
        FROM cases.status_condition
        WHERE id = ANY($1) AND dc = $2 AND status_id = $3
    ),
    initial_remaining AS (
        SELECT COUNT(*) AS count
        FROM cases.status_condition
        WHERE initial = TRUE AND id NOT IN (SELECT id FROM to_check) AND dc = $2 AND status_id = $3
    ),
    final_remaining AS (
        SELECT COUNT(*) AS count
        FROM cases.status_condition
        WHERE final = TRUE AND id NOT IN (SELECT id FROM to_check) AND dc = $2 AND status_id = $3
    ),
    initial_to_check AS (
        SELECT COUNT(*) AS count
        FROM to_check
        WHERE initial = TRUE
    ),
    final_to_check AS (
        SELECT COUNT(*) AS count
        FROM to_check
        WHERE final = TRUE
    ),
    delete_conditions AS (
        SELECT
            (SELECT count FROM initial_remaining) AS remaining_initial,
            (SELECT count FROM final_remaining) AS remaining_final,
            (SELECT count FROM initial_to_check) AS checking_initial,
            (SELECT count FROM final_to_check) AS checking_final
        FROM to_check
        LIMIT 1
    )
DELETE FROM cases.status_condition
WHERE id IN (SELECT id FROM to_check)
AND (
    (SELECT remaining_initial FROM delete_conditions) > 0 OR
    (SELECT checking_initial FROM delete_conditions) = 0
)
AND (
    (SELECT remaining_final FROM delete_conditions) > 0 OR
    (SELECT checking_final FROM delete_conditions) = 0
)
RETURNING id;
`
	args := []interface{}{pq.Array(ids), domainId, statusId}
	return query, args, nil
}

func (s StatusConditionStore) buildUpdateStatusConditionQuery(ctx *model.UpdateOptions, st *_go.StatusCondition) (string, []interface{}) {
	var setClauses []string
	var args []interface{}

	// Start placeholder numbering at 1
	placeholderIndex := 1

	// Add common fields with correct types
	args = append(args, ctx.Time)
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", placeholderIndex))
	placeholderIndex++

	args = append(args, ctx.Session.GetUserId())
	setClauses = append(setClauses, fmt.Sprintf("updated_by = $%d", placeholderIndex))
	placeholderIndex++

	// Track whether "initial" or "final" are being updated
	updateInitial := false
	updateFinal := false

	// Check and add update-specific fields if provided by the user
	for _, field := range ctx.Fields {
		switch field {
		case "name":
			if st.Name != "" {
				args = append(args, st.Name)
				setClauses = append(setClauses, fmt.Sprintf("name = $%d", placeholderIndex))
				placeholderIndex++
			}
		case "description":
			if st.Description != "" {
				args = append(args, st.Description)
				setClauses = append(setClauses, fmt.Sprintf("description = $%d", placeholderIndex))
				placeholderIndex++
			}
		case "initial":
			args = append(args, st.Initial)
			setClauses = append(setClauses, fmt.Sprintf("initial = $%d", placeholderIndex))
			placeholderIndex++
			updateInitial = true
		case "final":
			args = append(args, st.Final)
			setClauses = append(setClauses, fmt.Sprintf("final = $%d", placeholderIndex))
			placeholderIndex++
			updateFinal = true
		}
	}

	dc := ctx.Session.GetDomainId()
	statusId := st.StatusId
	id := st.Id

	query := fmt.Sprintf(`
   WITH	final_remaining AS (
		SELECT COUNT(*) AS count
		FROM cases.status_condition
		WHERE dc = $%d AND status_id = $%d AND final = TRUE
	),
	 set_initial_false AS (
		UPDATE cases.status_condition
		SET initial = FALSE
		WHERE dc = $%d AND status_id = $%d AND id <> $%d AND $%d = TRUE
	),
	upd AS (
		UPDATE cases.status_condition
		SET %s
		WHERE id = $%d AND dc = $%d
		RETURNING id, name, created_at, updated_at, description, initial, final, created_by, updated_by, status_id
	)
	SELECT
		upd.id,
		upd.name,
		upd.created_at,
		upd.updated_at,
		upd.description,
		upd.initial,
		upd.final,
		upd.created_by AS created_by_id,
		COALESCE(c.name::text, c.username) AS created_by_name,
		upd.updated_by AS updated_by_id,
		COALESCE(u.name::text, u.username) AS updated_by_name,
		upd.status_id
	FROM upd
	LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
	LEFT JOIN directory.wbt_user c ON c.id = upd.created_by
WHERE
    -- We DO NOT UPDATE initial, only try to update final to FALSE and checking if it's NOT the last one
    ($%d::boolean = FALSE AND $%d::boolean = FALSE AND (SELECT count FROM final_remaining) >1)

    OR
    -- WE ONLY UPDATE INITIAL [initial always true] and DON'T UPDATE FINAL
    ($%d::boolean = TRUE AND $%d::boolean = FALSE)

    OR
    -- WE UPDATE FINAL + INITIAL but final is FALSE so we checking if it's NOT the last one
    ($%d::boolean = TRUE AND $%d::boolean = TRUE
    AND $%d::boolean = FALSE
    AND (SELECT count FROM final_remaining) > 1)

	OR
    -- WE UPDATE FINAL + INITIAL - both are TRUE
    ($%d::boolean = TRUE AND $%d::boolean = TRUE
    AND $%d::boolean = TRUE AND $%d::boolean = TRUE)

	OR
    -- WE ONLY UPDATE FINAL TO TRUE
    ($%d::boolean = TRUE AND $%d::boolean = TRUE);
	`,
		// Arguments for set_initial_false
		placeholderIndex, placeholderIndex+1, placeholderIndex+2, placeholderIndex+3,
		// Arguments for final_remaining
		placeholderIndex+4, placeholderIndex+5,
		// Set clause for the update
		strings.Join(setClauses, ", "),
		// Arguments for upd subquery
		placeholderIndex+6, placeholderIndex+7,

		// --- Arguments for the WHERE clause checks ---

		// We DO NOT UPDATE initial, only try to update final to FALSE and checking if it's NOT the last one
		placeholderIndex+8, placeholderIndex+9,
		// WE ONLY UPDATE INITIAL [initial always true] and DON'T UPDATE FINAL
		placeholderIndex+10, placeholderIndex+11,
		// WE UPDATE FINAL + INITIAL but final is FALSE so we checking if it's NOT the last one
		placeholderIndex+12, placeholderIndex+13, placeholderIndex+14,
		// WE UPDATE FINAL + INITIAL - both are TRUE
		placeholderIndex+15, placeholderIndex+16, placeholderIndex+17, placeholderIndex+18,
		// WE ONLY UPDATE FINAL TO TRUE
		placeholderIndex+19, placeholderIndex+20,
	)

	// final_remaining check count for dc and status_id
	args = append(args, dc, statusId)
	// Arguments for set_initial_false
	args = append(args, dc, statusId, id, st.Initial)
	// args for upd subquery
	args = append(args, id, dc) // Arguments for upd subquery

	// --- Arguments for WHERE clause checks ---

	args = append(args,
		updateInitial, st.Final,
		st.Initial, updateFinal,
		updateFinal, updateInitial, st.Final,
		updateFinal, updateInitial, st.Final, st.Initial,
		updateFinal, st.Final,
	)

	return query, args
}

// func (s StatusConditionStore) isLastStatusCondition(ctx context.Context, tx pgx.Tx, dc int64, statusID int64, id int64, columnName string) (bool, error) {
// 	var count int
// 	query := fmt.Sprintf(`
// 		SELECT COUNT(*)
// 		FROM cases.status_condition
// 		WHERE dc = $1 AND status_id = $2 AND %s = true AND id <> $3
// 	`, columnName)

// 	err := tx.QueryRow(ctx, query, dc, statusID, id).Scan(&count)
// 	if err != nil {
// 		return false, err
// 	}
// 	return count == 0, nil
// }

func (s StatusConditionStore) getDBConnection() (*pgxpool.Pool, error) {
	db, err := s.storage.Database()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return nil, err
	}
	return db, nil
}

func (s StatusConditionStore) handleTx(ctx context.Context, tx pgx.Tx, err *error) {
	if p := recover(); p != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			log.Printf("Failed to rollback transaction: %v", rbErr)
		}
		panic(p)
	} else if *err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			log.Printf("Failed to rollback transaction: %v", rbErr)
		}
	} else {
		*err = tx.Commit(ctx)
	}
}

func (s StatusConditionStore) buildScanArgs(fields []string, st *_go.StatusCondition, createdBy, updatedBy *_go.Lookup, tempCreatedAt, tempUpdatedAt *time.Time) []interface{} {
	var scanArgs []interface{}
	for _, field := range fields {
		switch field {
		case "id":
			scanArgs = append(scanArgs, &st.Id)
		case "name":
			scanArgs = append(scanArgs, &st.Name)
		case "description":
			scanArgs = append(scanArgs, &st.Description)
		case "initial":
			scanArgs = append(scanArgs, &st.Initial)
		case "final":
			scanArgs = append(scanArgs, &st.Final)
		case "created_at":
			scanArgs = append(scanArgs, tempCreatedAt)
		case "updated_at":
			scanArgs = append(scanArgs, tempUpdatedAt)
		case "created_by":
			scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
		case "updated_by":
			scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
		case "status_id":
			scanArgs = append(scanArgs, &st.Id)
		}
	}
	return scanArgs
}

func (s StatusConditionStore) populateStatusConditionFields(fields []string, st *_go.StatusCondition, createdBy, updatedBy *_go.Lookup, tempCreatedAt, tempUpdatedAt time.Time) {
	if s.containsField(fields, "created_by") {
		st.CreatedBy = createdBy
	}
	if s.containsField(fields, "updated_by") {
		st.UpdatedBy = updatedBy
	}
	if s.containsField(fields, "created_at") {
		st.CreatedAt = tempCreatedAt.Unix()
	}
	if s.containsField(fields, "updated_at") {
		st.UpdatedAt = tempUpdatedAt.Unix()
	}
}

func (s StatusConditionStore) containsField(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func NewStatusConditionStore(store store.Store) (store.StatusConditionStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_status_condition.check.bad_arguments",
			"error creating config interface to the status_condition table, main store is nil")
	}
	return &StatusConditionStore{storage: store}, nil
}
