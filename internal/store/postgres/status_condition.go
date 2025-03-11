package postgres

import (
	"context"
	"fmt"
	"github.com/webitel/cases/model/options"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

const (
	statusConditionDefaultSort = "name"
)

type StatusConditionStore struct {
	storage *Store
}

func (s StatusConditionStore) Create(rpc options.CreateOptions, add *_go.StatusCondition) (*_go.StatusCondition, error) {
	db, err := s.getDBConnection()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.create.database_connection_error", err)
	}

	tx, err := db.BeginTx(rpc, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.create.transaction_begin_error", err)
	}
	defer s.handleTx(rpc, tx, &err)

	query, args, err := s.buildCreateStatusConditionQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.create.query_build_error", err)
	}

	var (
		createdBy, updatedBy _go.Lookup
		createdAt, updatedAt time.Time
	)

	err = tx.QueryRow(rpc, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &updatedAt, &add.Description, &add.Initial, &add.Final,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, &add.StatusId,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.create.execution_error", err)
	}

	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)
	add.CreatedBy = &createdBy
	add.UpdatedBy = &updatedBy

	return add, nil
}

func (s StatusConditionStore) List(rpc *model.SearchOptions, statusId int64) (*_go.StatusConditionList, error) {
	db, err := s.getDBConnection()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.list.database_connection_error", err)
	}

	query, args, err := s.buildListStatusConditionQuery(rpc, statusId)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.list.query_build_error", err)
	}

	rows, err := db.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.list.execution_error", err)
	}
	defer rows.Close()

	var statusList []*_go.StatusCondition
	lCount := 0
	next := false
	// Check if we want to fetch all records
	//
	// If the size is -1, we want to fetch all records
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		// If not fetching all records, check the size limit
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		st := &_go.StatusCondition{}

		var (
			createdBy, updatedBy         _go.Lookup
			tempCreatedAt, tempUpdatedAt time.Time
		)

		scanArgs := s.buildScanArgs(rpc.Fields, st, &createdBy, &updatedBy, &tempCreatedAt, &tempUpdatedAt)
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.status_condition.list.row_scan_error", err)
		}

		s.populateStatusConditionFields(rpc.Fields, st, &createdBy, &updatedBy, tempCreatedAt, tempUpdatedAt)
		statusList = append(statusList, st)
		lCount++
	}

	return &_go.StatusConditionList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: statusList,
	}, nil
}

func (s StatusConditionStore) Delete(rpc *model.DeleteOptions, statusId int64) error {
	domainId := rpc.GetAuthOpts().GetDomainId()

	query, args, err := s.buildDeleteStatusConditionQuery(rpc.ID, domainId, statusId)
	if err != nil {
		return dberr.NewDBInternalError("postgres.status_condition.delete.query_build_error", err)
	}

	db, err := s.getDBConnection()
	if err != nil {
		return dberr.NewDBInternalError("postgres.status_condition.delete.database_connection_error", err)
	}

	res, err := db.Exec(rpc.Context, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.status_condition.delete.execution_error", err)
	}

	// Check if any rows were affected
	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.status_condition.delete.not_found")
	}

	return nil
}

func (s StatusConditionStore) Update(rpc options.UpdateOptions, st *_go.StatusCondition) (*_go.StatusCondition, error) {
	db, err := s.getDBConnection()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.update.database_connection_error", err)
	}

	tx, err := db.BeginTx(rpc, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.update.transaction_begin_error", err)
	}
	defer s.handleTx(rpc, tx, &err)

	for _, field := range rpc.GetMask() {
		switch field {
		case "initial":
			if !st.Initial {
				return nil, dberr.NewDBCheckViolationError("postgres.status_condition.update.initial_false_not_allowed", "update not allowed: there must be at least one initial = TRUE for the given dc and status_id")
			}
		}
	}

	query, args := s.buildUpdateStatusConditionQuery(rpc, st)

	var (
		createdBy, updatedBy _go.Lookup
		createdAt, updatedAt time.Time
	)

	err = tx.QueryRow(rpc, query, args...).Scan(
		&st.Id, &st.Name, &createdAt, &updatedAt, &st.Description, &st.Initial, &st.Final,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name, &st.StatusId,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status_condition.update.execution_error", err)
	}

	st.CreatedAt = util.Timestamp(createdAt)
	st.UpdatedAt = util.Timestamp(updatedAt)
	st.CreatedBy = &createdBy
	st.UpdatedBy = &updatedBy

	return st, nil
}

func (s StatusConditionStore) buildCreateStatusConditionQuery(rpc options.CreateOptions, status *_go.StatusCondition) (string, []interface{}, error) {
	query := createStatusConditionQuery
	args := []interface{}{
		status.Name,                 // $1 name
		rpc.GetTime(),               // $2 created_at / updated_at
		status.Description,          // $3 description
		rpc.GetAuth().GetUserId(),   // $4 created_by / updated_by
		rpc.GetAuth().GetDomainId(), // $5 dc
		status.StatusId,             // $6 status_id
	}
	return query, args, nil
}

func (s StatusConditionStore) buildListStatusConditionQuery(rpc *model.SearchOptions, statusId int64) (string, []interface{}, error) {
	queryBuilder := sq.Select().
		From("cases.status_condition AS s").
		Where(sq.Eq{"s.dc": rpc.GetAuthOpts().GetDomainId(), "s.status_id": statusId}).
		PlaceholderFormat(sq.Dollar)

	fields := util.FieldsFunc(rpc.Fields, util.InlineFields)
	rpc.Fields = append(fields, "id")

	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "initial", "final", "created_at", "updated_at", "description":
			queryBuilder = queryBuilder.Column("s." + field)
		case "created_by":
			// Handle nulls using COALESCE for created_by
			queryBuilder = queryBuilder.
				Column("COALESCE(created_by.id, 0) AS cbi").
				Column("COALESCE(created_by.name, '') AS cbn").
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
		case "updated_by":
			// Handle nulls using COALESCE for updated_by
			queryBuilder = queryBuilder.
				Column("COALESCE(updated_by.id, 0) AS ubi").
				Column("COALESCE(updated_by.name, '') AS ubn").
				LeftJoin("directory.wbt_auth AS updated_by ON g.updated_by = updated_by.id")
		}
	}

	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"s.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substrs := util.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"s.name": combinedLike})
	}

	// -------- Apply sorting ----------
	queryBuilder = store.ApplyDefaultSorting(rpc, queryBuilder, statusConditionDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = store.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Convert the query to SQL and arguments
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.status_condition.list.query_build_error", err)
	}

	return store.CompactSQL(query), args, nil
}

func (s StatusConditionStore) buildDeleteStatusConditionQuery(id int64, domainId, statusId int64) (string, []interface{}, error) {
	query := deleteStatusConditionQuery

	args := []interface{}{
		id,       // $1 id
		domainId, // $2 dc
		statusId, // $3 status_id
	}
	return query, args, nil
}

func (s StatusConditionStore) buildUpdateStatusConditionQuery(rpc options.UpdateOptions, st *_go.StatusCondition) (string, []interface{}) {
	var args []interface{}

	// 1. Squirrel operations: Building the dynamic part of the "upd" query
	updBuilder := sq.Update("cases.status_condition").
		Set("updated_at", rpc.GetTime()).
		Set("updated_by", rpc.GetAuth().GetUserId())

	// Track whether "initial" or "final" are being updated
	updateInitial := false
	updateFinal := false

	// Add update-specific fields if provided by the user
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			if st.Name != "" {
				updBuilder = updBuilder.Set("name", st.Name)
			}
		case "description":
			// Set description to NULL if it's an empty string
			updBuilder = updBuilder.Set("description", sq.Expr("NULLIF(?, '')", st.Description))
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
		Where(sq.Eq{"dc": rpc.GetAuth().GetDomainId()}).
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
       WITH final_remaining AS (SELECT COUNT(*) AS count
                         FROM cases.status_condition
                         WHERE dc = $1
                           AND status_id = $2
                           AND final = TRUE),
     set_initial_false AS (
         UPDATE cases.status_condition
             SET initial = FALSE
             WHERE dc = $1 AND status_id = $2 AND id <> $3 AND $7 = TRUE),
     upd AS (%s)
SELECT upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       COALESCE(upd.description, '')              AS description,
       upd.initial,
       upd.final,
       upd.created_by                             AS created_by_id,
       COALESCE(c.name::text, c.username,'')      AS created_by_name,
       upd.updated_by                             AS updated_by_id,
       COALESCE(u.name::text, u.username)         AS updated_by_name,
       upd.status_id
FROM upd
         LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = upd.created_by
WHERE CASE
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
		rpc.GetAuth().GetDomainId(), // $1
		st.StatusId,                 // $2
		st.Id,                       // $3
		updateInitial,               // $4
		updateFinal,                 // $5
		st.Final,                    // $6
		st.Initial,                  // $7
	)

	// Append the dynamic query arguments
	args = append(args, updArgs...)
	// fmt.Printf("Executing SQL: %s\nWith args: %v\n", query, args)

	return store.CompactSQL(query), args
}

func (s StatusConditionStore) getDBConnection() (*pgxpool.Pool, error) {
	db, err := s.storage.Database()
	if err != nil {
		log.Printf("Failed to get database connection: %v", err)
		return nil, err
	}
	return db, nil
}

func (s StatusConditionStore) handleTx(rpc context.Context, tx pgx.Tx, err *error) {
	if p := recover(); p != nil {
		if rbErr := tx.Rollback(rpc); rbErr != nil {
			log.Printf("Failed to rollback transaction: %v", rbErr)
		}
		panic(p)
	} else if *err != nil {
		if rbErr := tx.Rollback(rpc); rbErr != nil {
			log.Printf("Failed to rollback transaction: %v", rbErr)
		}
	} else {
		*err = tx.Commit(rpc)
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
			scanArgs = append(scanArgs, scanner.ScanText(&st.Description))
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
WITH existing_status AS (SELECT COUNT(*) AS count FROM cases.status_condition WHERE dc = $5 AND status_id = $6),
     default_values
         AS (SELECT CASE WHEN (SELECT count FROM existing_status) = 0 THEN TRUE ELSE FALSE END AS initial_default,
                    CASE WHEN (SELECT count FROM existing_status) = 0 THEN TRUE ELSE FALSE END AS final_default),
     ins AS (INSERT INTO cases.status_condition (name, created_at, description, initial, final, created_by, updated_at,
                                                 updated_by, dc, status_id)
         VALUES ($1, $2, NULLIF($3, ''), (SELECT initial_default FROM default_values),
                 (SELECT final_default FROM default_values), $4, $2, $4, $5, $6)
         RETURNING id, name, created_at, updated_at, description, initial, final, created_by, updated_by, status_id)
SELECT ins.id,
       ins.name,
       ins.created_at,
       ins.updated_at,
       COALESCE(ins.description, '')      AS description,
       ins.initial,
       ins.final,
       ins.created_by                     AS created_by_id,
       COALESCE(c.name::text, c.username) AS created_by_name,
       ins.updated_by                     AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name,
       ins.status_id
FROM ins
         LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;`)

	deleteStatusConditionQuery = store.CompactSQL(`
		 WITH
			 to_check AS (
				 SELECT id, initial, final
				 FROM cases.status_condition
				 WHERE id = $1 AND dc = $2 AND status_id = $3
			 ),
			 remaining_conditions AS (
				 SELECT
					 COUNT(*) FILTER (WHERE initial = TRUE AND id NOT IN (SELECT id FROM to_check)) AS remaining_initial,
					 COUNT(*) FILTER (WHERE final = TRUE AND id NOT IN (SELECT id FROM to_check)) AS remaining_final
				 FROM cases.status_condition
				 WHERE dc = $2 AND status_id = $3
			 ),
			 to_check_conditions AS (
				 SELECT
					 COUNT(*) FILTER (WHERE initial = TRUE) AS checking_initial,
					 COUNT(*) FILTER (WHERE final = TRUE) AS checking_final
				 FROM to_check
			 ),
			 delete_allowed AS (
				 SELECT
					 (remaining_conditions.remaining_initial > 0 OR to_check_conditions.checking_initial = 0) AS can_delete_initial,
					 (remaining_conditions.remaining_final > 0 OR to_check_conditions.checking_final = 0) AS can_delete_final
				 FROM remaining_conditions, to_check_conditions
			 )
		 DELETE
		 FROM cases.status_condition
		 WHERE id IN (SELECT id FROM to_check)
		   AND (SELECT can_delete_initial FROM delete_allowed)
		   AND (SELECT can_delete_final FROM delete_allowed)
		 RETURNING id;
		 `)
)

func NewStatusConditionStore(store *Store) (store.StatusConditionStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_status_condition.check.bad_arguments",
			"error creating status condition interface to the status_condition table, main store is nil")
	}
	return &StatusConditionStore{storage: store}, nil
}
