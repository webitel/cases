package postgres

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

const (
	closeReasonDefaultSort = "name"
)

type CloseReason struct {
	storage store.Store
}

// Create implements store.CloseReasonStore.
func (s *CloseReason) Create(rpc *model.CreateOptions, add *_go.CloseReason) (*_go.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.database_connection_error", dbErr)
	}

	query, args, err := s.buildCreateCloseReasonQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup _go.Lookup
		createdAt, updatedAt             time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name, &add.CloseReasonGroupId,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.execution_error", err)
	}

	t := rpc.Time
	return &_go.CloseReason{
		Id:                 add.Id,
		Name:               add.Name,
		Description:        add.Description,
		CreatedAt:          util.Timestamp(t),
		UpdatedAt:          util.Timestamp(t),
		CloseReasonGroupId: add.CloseReasonGroupId,
		CreatedBy:          &createdByLookup,
		UpdatedBy:          &updatedByLookup,
	}, nil
}

// List implements store.CloseReasonStore.
func (s *CloseReason) List(rpc *model.SearchOptions, closeReasonId int64) (*_go.CloseReasonList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.database_connection_error", dbErr)
	}

	query, args, err := s.buildSearchCloseReasonQuery(rpc, closeReasonId)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.query_build_error", err)
	}

	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.execution_error", err)
	}
	defer rows.Close()

	var lookupList []*_go.CloseReason
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		l := &_go.CloseReason{}

		var (
			createdBy, updatedBy         _go.Lookup
			tempCreatedAt, tempUpdatedAt time.Time
		)

		scanArgs := s.buildScanArgs(rpc.Fields, l, &createdBy, &updatedBy, &tempCreatedAt, &tempUpdatedAt)
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.close_reason.list.row_scan_error", err)
		}

		s.populateCloseReasonFields(rpc.Fields, l, &createdBy, &updatedBy, tempCreatedAt, tempUpdatedAt)
		lookupList = append(lookupList, l)
		lCount++
	}

	return &_go.CloseReasonList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

// Delete implements store.CloseReasonStore.
func (s *CloseReason) Delete(rpc *model.DeleteOptions, closeReasonId int64) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.close_reason.delete.database_connection_error", dbErr)
	}

	query, args, err := s.buildDeleteCloseReasonQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.close_reason.delete.query_build_error", err)
	}

	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.close_reason.delete.execution_error", err)
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBNoRowsError("postgres.close_reason.delete.no_rows_affected")
	}

	return nil
}

// Update implements store.CloseReasonStore.
func (s *CloseReason) Update(rpc *model.UpdateOptions, l *_go.CloseReason) (*_go.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.database_connection_error", dbErr)
	}

	query, args, err := s.buildUpdateCloseReasonQuery(rpc, l)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.query_build_error", err)
	}

	var (
		createdBy, updatedBy _go.Lookup
		createdAt, updatedAt time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt, &l.Description,
		&createdBy.Id, &createdBy.Name, &updatedBy.Id, &updatedBy.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.execution_error", err)
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedBy

	return l, nil
}

// buildCreateCloseReasonQuery constructs the SQL insert query and returns the query string and arguments.
func (s CloseReason) buildCreateCloseReasonQuery(rpc *model.CreateOptions, lookup *_go.CloseReason) (string, []interface{}, error) {
	query := createCloseReasonQuery
	args := []interface{}{
		lookup.Name,
		rpc.Session.GetDomainId(),
		rpc.Time,
		lookup.Description,
		rpc.Session.GetUserId(),
		lookup.CloseReasonGroupId,
	}
	return query, args, nil
}

// buildSearchCloseReasonQuery constructs the SQL search query and returns the query builder.
func (s CloseReason) buildSearchCloseReasonQuery(rpc *model.SearchOptions, closeReasonId int64) (string, []interface{}, error) {
	queryBuilder := sq.Select().
		From("cases.close_reason AS g").
		Where(sq.Eq{"g.dc": rpc.Session.GetDomainId(), "g.close_reason_id": closeReasonId}).
		PlaceholderFormat(sq.Dollar)

	fields := util.FieldsFunc(rpc.Fields, util.InlineFields)
	rpc.Fields = append(fields, "id")

	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
		case "description":
			queryBuilder = queryBuilder.Column("COALESCE(g.description, '') AS description")
		case "created_by":
			queryBuilder = queryBuilder.
				Column("COALESCE(created_by.id, 0) AS cbi").
				Column("COALESCE(created_by.name, '') AS cbn").
				LeftJoin("directory.wbt_auth AS created_by ON p.created_by = created_by.id")
		case "updated_by":
			queryBuilder = queryBuilder.
				Column("COALESCE(updated_by.id, 0) AS ubi").
				Column("COALESCE(updated_by.name, '') AS ubn").
				LeftJoin("directory.wbt_auth AS updated_by ON p.updated_by = updated_by.id")
		}
	}

	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substrs := util.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": combinedLike})
	}

	// -------- Apply sorting ----------
	queryBuilder = store.ApplyDefaultSorting(rpc, queryBuilder, closeReasonDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = store.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.close_reason.query_build_error", err)
	}

	return store.CompactSQL(query), args, nil
}

// buildDeleteCloseReasonQuery constructs the SQL delete query and returns the query string and arguments.
func (s CloseReason) buildDeleteCloseReasonQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deleteCloseReasonQuery
	args := []interface{}{pq.Array(ids), rpc.Session.GetDomainId()}
	return query, args, nil
}

func (s CloseReason) buildUpdateCloseReasonQuery(rpc *model.UpdateOptions, l *_go.CloseReason) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	updateBuilder := psql.Update("cases.close_reason").
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId())

	for _, field := range rpc.Fields {
		switch field {
		case "name":
			if l.Name != "" {
				updateBuilder = updateBuilder.Set("name", l.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", l.Description))
		}
	}

	updateBuilder = updateBuilder.Where(sq.Eq{"id": l.Id, "dc": rpc.Session.GetDomainId()})

	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	query := fmt.Sprintf(`
WITH upd AS (%s
    RETURNING id, name, created_at, updated_at, description, created_by, updated_by)
SELECT upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       COALESCE(upd.description, '')          AS description,
       upd.created_by                         AS created_by_id,
       COALESCE(c.name::text, c.username, '') AS created_by_name,
       upd.updated_by                         AS updated_by_id,
       COALESCE(u.name::text, u.username)     AS updated_by_name
FROM upd
         LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = upd.created_by;
`, sql)

	return store.CompactSQL(query), args, nil
}

var (
	createCloseReasonQuery = store.CompactSQL(`
	WITH ins AS (
		INSERT INTO cases.close_reason (name, dc, created_at, description, created_by, updated_at, updated_by, close_reason_id)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $3, $5, $6)
		RETURNING *
	)
	SELECT ins.id,
		   ins.name,
		   ins.created_at,
		   COALESCE(ins.description, '')      AS description,
		   ins.created_by                     AS created_by_id,
		   COALESCE(c.name::text, c.username) AS created_by_name,
		   ins.updated_at,
		   ins.updated_by                     AS updated_by_id,
		   COALESCE(u.name::text, u.username) AS updated_by_name,
		   ins.close_reason_id
	FROM ins
	LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
	LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;
	`)

	deleteCloseReasonQuery = store.CompactSQL(`
DELETE FROM cases.close_reason
WHERE id = ANY($1) AND dc = $2
`)
)

// buildScanArgs prepares the arguments for scanning SQL rows.
func (s CloseReason) buildScanArgs(fields []string, r *_go.CloseReason, createdBy, updatedBy *_go.Lookup, tempCreatedAt, tempUpdatedAt *time.Time) []interface{} {
	var scanArgs []interface{}

	for _, field := range fields {
		switch field {
		case "id":
			scanArgs = append(scanArgs, &r.Id)
		case "name":
			scanArgs = append(scanArgs, &r.Name)
		case "description":
			scanArgs = append(scanArgs, &r.Description)
		case "created_at":
			scanArgs = append(scanArgs, tempCreatedAt)
		case "updated_at":
			scanArgs = append(scanArgs, tempUpdatedAt)
		case "created_by":
			scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
		case "updated_by":
			scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
		}
	}
	return scanArgs
}

// populateReasonFields populates the Reason struct with the scanned values.
func (s CloseReason) populateCloseReasonFields(fields []string, r *_go.CloseReason, createdBy, updatedBy *_go.Lookup, tempCreatedAt, tempUpdatedAt time.Time) {
	if s.containsField(fields, "created_by") {
		r.CreatedBy = createdBy
	}
	if s.containsField(fields, "updated_by") {
		r.UpdatedBy = updatedBy
	}
	if s.containsField(fields, "created_at") {
		r.CreatedAt = util.Timestamp(tempCreatedAt)
	}
	if s.containsField(fields, "updated_at") {
		r.UpdatedAt = util.Timestamp(tempUpdatedAt)
	}
}

// containsField checks if a field is in the list of fields.
func (s CloseReason) containsField(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func NewCloseReasonStore(store store.Store) (store.CloseReasonStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_close_reason.check.bad_arguments",
			"error creating close_reason interface, main store is nil")
	}
	return &CloseReason{storage: store}, nil
}
