package postgres

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type StatusConditionStore struct {
	storage store.Store
}

func (s StatusConditionStore) Create(ctx *model.CreateOptions, add *_go.StatusCondition) (*_go.StatusCondition, error) {
	db, err := s.getDBConnection()
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.create.database_connection_error", err.Error())
	}

	tx, err := db.BeginTx(ctx.Context, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.create.transaction_begin_error", err.Error())
	}
	defer s.handleTx(ctx.Context, tx, &err)

	query, args, err := s.buildCreateStatusConditionQuery(ctx, add)
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.create.query_build_error", err.Error())
	}

	var createdBy, updatedBy _go.Lookup
	var createdAt, updatedAt time.Time

	err = tx.QueryRow(ctx.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &updatedAt, &add.Description, &add.Initial, &add.Final,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, &add.StatusId,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.create.execution_error", err.Error())
	}

	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)
	add.CreatedBy = &createdBy
	add.UpdatedBy = &updatedBy

	return add, nil
}

func (s StatusConditionStore) List(ctx *model.SearchOptions, statusId int64) (*_go.StatusConditionList, error) {
	db, err := s.getDBConnection()
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.list.database_connection_error", err.Error())
	}

	query, args, err := s.buildListStatusConditionQuery(ctx, statusId)
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.list.query_build_error", err.Error())
	}

	rows, err := db.Query(ctx.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.list.execution_error", err.Error())
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
			return nil, model.NewInternalError("postgres.status_condition.list.row_scan_error", err.Error())
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
	ids := []int64{statusId}
	domainId := ctx.Session.GetDomainId()

	query, args, err := s.buildDeleteStatusConditionQuery(ids, domainId, statusId)
	if err != nil {
		return model.NewInternalError("postgres.status_condition.delete.query_build_error", err.Error())
	}

	db, err := s.getDBConnection()
	if err != nil {
		return model.NewInternalError("postgres.status_condition.delete.database_connection_error", err.Error())
	}

	rows, err := db.Query(ctx.Context, query, args...)
	if err != nil {
		return model.NewInternalError("postgres.status_condition.delete.execution_error", err.Error())
	}
	defer rows.Close()

	var deletedIds []int64
	for rows.Next() {
		var deletedId int64
		if err := rows.Scan(&deletedId); err != nil {
			return model.NewInternalError("postgres.status_condition.delete.scan_error", err.Error())
		}
		deletedIds = append(deletedIds, deletedId)
	}

	if len(deletedIds) == 0 {
		return model.NewInternalError("postgres.status_condition.delete.constraint_violation", "operation would violate constraints: at least one initial and one final record must remain")
	}

	return nil
}

func (s StatusConditionStore) Update(ctx *model.UpdateOptions, st *_go.StatusCondition) (*_go.StatusCondition, error) {
	db, err := s.getDBConnection()
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.update.database_connection_error", err.Error())
	}

	tx, err := db.BeginTx(ctx.Context, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.update.transaction_begin_error", err.Error())
	}
	defer s.handleTx(ctx.Context, tx, &err)

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

	for _, field := range ctx.Fields {
		switch field {
		case "initial":
			if !st.Initial {
				return nil, model.NewInternalError("postgres.status_condition.update.initial_false_not_allowed", "update not allowed: there must be at least one initial = TRUE for the given dc and status_id")
			}
		}
	}

	query, args := s.buildUpdateStatusConditionQuery(ctx, st)

	var createdBy, updatedBy _go.Lookup
	var createdAt, updatedAt time.Time

	err = tx.QueryRow(ctx.Context, query, args...).Scan(
		&st.Id, &st.Name, &createdAt, &updatedAt, &st.Description, &st.Initial, &st.Final,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, &st.StatusId,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.status_condition.update.execution_error", err.Error())
	}

	st.CreatedAt = util.Timestamp(createdAt)
	st.UpdatedAt = util.Timestamp(updatedAt)
	st.CreatedBy = &createdBy
	st.UpdatedBy = &updatedBy

	return st, nil
}

func (s StatusConditionStore) buildCreateStatusConditionQuery(ctx *model.CreateOptions, status *_go.StatusCondition) (string, []interface{}, error) {
	query := createStatusConditionQuery
	args := []interface{}{
		status.Name, ctx.Time, status.Description,
		ctx.Session.GetUserId(), ctx.Session.GetDomainId(), status.StatusId,
	}
	return query, args, nil
}

func (s StatusConditionStore) buildListStatusConditionQuery(ctx *model.SearchOptions, statusId int64) (string, []interface{}, error) {
	queryBuilder := sq.Select().
		From("cases.status_condition AS s").
		Where(sq.Eq{"s.dc": ctx.Session.GetDomainId(), "s.status_id": statusId}).
		PlaceholderFormat(sq.Dollar)

	fields := ctx.FieldsUtil.FieldsFunc(ctx.Fields, ctx.FieldsUtil.InlineFields)
	ctx.Fields = append(fields, "id")

	for _, field := range ctx.Fields {
		switch field {
		case "id", "name", "description", "initial", "final", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("s." + field)
		case "created_by":
			// cbi = created_by_id
			// cbn = created_by_name
			queryBuilder = queryBuilder.Column("created_by.id AS cbi, created_by.name AS cbn").
				LeftJoin("directory.wbt_auth AS created_by ON s.created_by = created_by.id")
		case "updated_by":
			// ubi = updated_by_id
			// ubn = updated_by_name
			queryBuilder = queryBuilder.Column("updated_by.id AS ubi, updated_by.name AS ubn").
				LeftJoin("directory.wbt_auth AS updated_by ON s.updated_by = updated_by.id")
		}
	}

	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"s.id": ids})
	}

	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substrs := ctx.Match.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"s.name": "%" + combinedLike + "%"})
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

	page := ctx.GetPage()
	size := ctx.GetSize()

	// Apply sorting
	queryBuilder = queryBuilder.OrderBy(sortFields...)

	// Apply offset only if page > 1
	if page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	// Apply limit
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1))
	}

	// Convert the query to SQL and arguments
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.status_condition.list.query_build_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

func (s StatusConditionStore) buildDeleteStatusConditionQuery(ids []int64, domainId, statusId int64) (string, []interface{}, error) {
	query := deleteStatusConditionQuery
	args := []interface{}{pq.Array(ids), domainId, statusId}
	return query, args, nil
}

func (s StatusConditionStore) buildUpdateStatusConditionQuery(ctx *model.UpdateOptions, st *_go.StatusCondition) (string, []interface{}) {
	var args []interface{}

	// 1. Squirrel operations: Building the dynamic part of the "upd" query
	updBuilder := sq.Update("cases.status_condition").
		Set("updated_at", ctx.Time).
		Set("updated_by", ctx.Session.GetUserId())

	// Track whether "initial" or "final" are being updated
	updateInitial := false
	updateFinal := false

	// Add update-specific fields if provided by the user
	for _, field := range ctx.Fields {
		switch field {
		case "name":
			if st.Name != "" {
				updBuilder = updBuilder.Set("name", st.Name)
			}
		case "description":
			if st.Description != "" {
				updBuilder = updBuilder.Set("description", st.Description)
			}
		case "initial":
			updBuilder = updBuilder.Set("initial", st.Initial)
			updateInitial = true
		case "final":
			updBuilder = updBuilder.Set("final", st.Final)
			updateFinal = true
		}
	}

	// Build the dynamic part of the "upd" query using squirrel
	updSql, updArgs, err := updBuilder.
		Where(sq.Eq{"id": st.Id}).
		Where(sq.Eq{"dc": ctx.Session.GetDomainId()}).
		Suffix("RETURNING id, name, created_at, updated_at, description, initial, final, created_by, updated_by, status_id").
		ToSql()
	if err != nil {
		return "", nil
	}

	// Manually replace "?" placeholders with "$n" placeholders
	// assuming the starting index of placeholders is $8, following the existing $1 to $7
	for i := range updArgs {
		placeholder := fmt.Sprintf("$%d", i+8)
		updSql = strings.Replace(updSql, "?", placeholder, 1)
	}

	// 2. Main query using fmt.Sprintf without Squirrel placeholder format
	query := fmt.Sprintf(`
        WITH final_remaining AS (
            SELECT COUNT(*) AS count
            FROM cases.status_condition
            WHERE dc = $1 AND status_id = $2  AND final = TRUE
        ),
        set_initial_false AS (
            UPDATE cases.status_condition
            SET initial = FALSE
            WHERE dc = $1 AND status_id = $2 AND id <> $3 AND $7 = TRUE
        ),
        upd AS (
            %s
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
            CASE
                -- Ensure the update only happens if the conditions are met
				-- WE DO NOT UPDATE FINAL & INITIAL
				WHEN $4::boolean = FALSE AND $5::boolean = FALSE THEN TRUE

				-- WE ONLY UPDATE FINAL - so checking if it's NOT the last one
                WHEN $4::boolean = FALSE AND $6::boolean = FALSE THEN (SELECT count FROM final_remaining) > 1

				-- WE ONLY UPDATE INITIAL [initial always true] and DON'T UPDATE FINAL
                WHEN $7::boolean = TRUE AND $5::boolean = FALSE THEN TRUE

				-- WE UPDATE FINAL + INITIAL but final is FALSE so we checking if it's NOT the last one
                WHEN $5::boolean = TRUE AND $4::boolean = TRUE THEN
                    CASE
                        WHEN $6::boolean = FALSE THEN (SELECT count FROM final_remaining) > 1
                        ELSE TRUE
                    END
                ELSE TRUE
            END
    `, updSql)

	// 3. Adding all arguments
	args = append(args,
		ctx.Session.GetDomainId(), // $1
		st.StatusId,               // $2
		st.Id,                     // $3
		updateInitial,             // $4
		updateFinal,               // $5
		st.Final,                  // $6
		st.Initial,                // $7
	)

	// Append the dynamic query arguments
	args = append(args, updArgs...)
	fmt.Printf("Executing SQL: %s\nWith args: %v\n", query, args)

	return store.CompactSQL(query), args
}

// func (s StatusConditionStore) buildUpdateStatusConditionQuery(ctx *model.UpdateOptions, st *_go.StatusCondition) (string, []interface{}) {
// 	var setClauses []string
// 	var args []interface{}

// 	// Start placeholder numbering at 1
// 	placeholderIndex := 1

// 	// Add common fields with correct types
// 	args = append(args, ctx.Time)
// 	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", placeholderIndex))
// 	placeholderIndex++

// 	args = append(args, ctx.Session.GetUserId())
// 	setClauses = append(setClauses, fmt.Sprintf("updated_by = $%d", placeholderIndex))
// 	placeholderIndex++

// 	// Track whether "initial" or "final" are being updated
// 	updateInitial := false
// 	updateFinal := false

// 	// Check and add update-specific fields if provided by the user
// 	for _, field := range ctx.Fields {
// 		switch field {
// 		case "name":
// 			if st.Name != "" {
// 				args = append(args, st.Name)
// 				setClauses = append(setClauses, fmt.Sprintf("name = $%d", placeholderIndex))
// 				placeholderIndex++
// 			}
// 		case "description":
// 			if st.Description != "" {
// 				args = append(args, st.Description)
// 				setClauses = append(setClauses, fmt.Sprintf("description = $%d", placeholderIndex))
// 				placeholderIndex++
// 			}
// 		case "initial":
// 			args = append(args, st.Initial)
// 			setClauses = append(setClauses, fmt.Sprintf("initial = $%d", placeholderIndex))
// 			placeholderIndex++
// 			updateInitial = true
// 		case "final":
// 			args = append(args, st.Final)
// 			setClauses = append(setClauses, fmt.Sprintf("final = $%d", placeholderIndex))
// 			placeholderIndex++
// 			updateFinal = true
// 		}
// 	}

// 	query := fmt.Sprintf(`
// 	    -- Count of final status conditions for the given dc and status_id and final == TRUE
//         WITH final_remaining AS (
//             SELECT COUNT(*) AS count
//             FROM cases.status_condition
//             WHERE dc = $%d AND status_id = $%d AND final = TRUE
//         ),
// 		-- If user set initial to TRUE - set all another initial for this dc and status_id to FALSE
//         set_initial_false AS (
//             UPDATE cases.status_condition
//             SET initial = FALSE
//             WHERE dc = $%d AND status_id = $%d AND id <> $%d AND $%d = TRUE
//         ),
//         upd AS (
//             UPDATE cases.status_condition
//             SET %s
//             WHERE id = $%d AND dc = $%d
//             RETURNING id, name, created_at, updated_at, description, initial, final, created_by, updated_by, status_id
//         )
//         SELECT
//             upd.id,
//             upd.name,
//             upd.created_at,
//             upd.updated_at,
//             upd.description,
//             upd.initial,
//             upd.final,
//             upd.created_by AS created_by_id,
//             COALESCE(c.name::text, c.username) AS created_by_name,
//             upd.updated_by AS updated_by_id,
//             COALESCE(u.name::text, u.username) AS updated_by_name,
//             upd.status_id
//         FROM upd
//         LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
//         LEFT JOIN directory.wbt_user c ON c.id = upd.created_by
//         WHERE
//             CASE
// 			    -- WE DO NOT UPDATE initial, only try to update final to FALSE and checking if it's NOT the last one
//                 WHEN $%d::boolean = FALSE AND $%d::boolean = FALSE THEN (SELECT count FROM final_remaining) > 1

// 				-- WE ONLY UPDATE INITIAL [initial always true] and DON'T UPDATE FINAL
//                 WHEN $%d::boolean = TRUE AND $%d::boolean = FALSE THEN TRUE

// 				-- WE UPDATE FINAL + INITIAL but final is FALSE so we checking if it's NOT the last one
//                 WHEN $%d::boolean = TRUE AND $%d::boolean = TRUE THEN
//                     CASE
//                         WHEN $%d::boolean = FALSE THEN (SELECT count FROM final_remaining) > 1
//                         ELSE TRUE
//                     END
//                 ELSE TRUE
//             END
//         `,

// 		// Arguments for final_remaining
// 		placeholderIndex, placeholderIndex+1,
// 		// Arguments for set_initial_false
// 		placeholderIndex+2, placeholderIndex+3, placeholderIndex+4, placeholderIndex+5,
// 		// Set clause for the update
// 		strings.Join(setClauses, ", "),
// 		// Arguments for upd subquery
// 		placeholderIndex+6, placeholderIndex+7,

// 		// --- Arguments for the WHERE clause checks ---

// 		// We DO NOT UPDATE initial, only try to update final to FALSE and checking if it's NOT the last one
// 		placeholderIndex+8, placeholderIndex+9,
// 		// WE ONLY UPDATE INITIAL [initial always true] and DON'T UPDATE FINAL
// 		placeholderIndex+10, placeholderIndex+11,
// 		// WE UPDATE FINAL + INITIAL but final is FALSE so we checking if it's NOT the last one
// 		placeholderIndex+12, placeholderIndex+13, placeholderIndex+14,
// 	)

// 	// final_remaining check count for dc and status_id
// 	args = append(args, ctx.Session.GetDomainId(), st.StatusId)
// 	// Arguments for set_initial_false
// 	args = append(args, ctx.Session.GetDomainId(), st.StatusId, st.Id, st.Initial)
// 	// args for upd subquery
// 	args = append(args, st.Id, ctx.Session.GetDomainId())

// 	// --- Arguments for WHERE clause checks ---
// 	args = append(args,
// 		updateInitial, st.Final,
// 		st.Initial, updateFinal,
// 		updateFinal, updateInitial, st.Final,
// 	)

// 	return query, args
// }

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
		st.CreatedAt = util.Timestamp(tempCreatedAt)
	}
	if s.containsField(fields, "updated_at") {
		st.UpdatedAt = util.Timestamp(tempUpdatedAt)
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

// ---- STATIC SQL QUERIES ----
var (
	createStatusConditionQuery = store.CompactSQL(`
	WITH existing_status AS (
    SELECT COUNT(*) AS count
    FROM cases.status_condition
    WHERE dc = $5 AND status_id = $6
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
        $4, $2, $4, $5, $6
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
LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;`)

	deleteStatusConditionQuery = store.CompactSQL(`WITH
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
RETURNING id;`)
)

func NewStatusConditionStore(store store.Store) (store.StatusConditionStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_status_condition.check.bad_arguments",
			"error creating status condition interface to the status_condition table, main store is nil")
	}
	return &StatusConditionStore{storage: store}, nil
}
