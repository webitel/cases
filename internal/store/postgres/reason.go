package postgres

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type Reason struct {
	storage store.Store
}

// Create implements store.ReasonStore.
func (s *Reason) Create(rpc *model.CreateOptions, add *_go.Reason) (*_go.Reason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.reason.create.database_connection_error", dbErr.Error())
	}

	query, args, err := s.buildCreateReasonQuery(rpc, add)
	if err != nil {
		return nil, model.NewInternalError("postgres.reason.create.query_build_error", err.Error())
	}

	var (
		createdByLookup, updatedByLookup _go.Lookup
		createdAt, updatedAt             time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name, &add.CloseReasonId,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.reason.create.execution_error", err.Error())
	}

	t := rpc.Time
	return &_go.Reason{
		Id:            add.Id,
		Name:          add.Name,
		Description:   add.Description,
		CreatedAt:     util.Timestamp(t),
		UpdatedAt:     util.Timestamp(t),
		CloseReasonId: add.CloseReasonId,
		CreatedBy:     &createdByLookup,
		UpdatedBy:     &updatedByLookup,
	}, nil
}

// List implements store.ReasonStore.
func (s *Reason) List(rpc *model.SearchOptions, closeReasonId int64) (*_go.ReasonList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.reason.list.database_connection_error", dbErr.Error())
	}

	query, args, err := s.buildSearchReasonQuery(rpc, closeReasonId)
	if err != nil {
		return nil, model.NewInternalError("postgres.reason.list.query_build_error", err.Error())
	}

	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.reason.list.execution_error", err.Error())
	}
	defer rows.Close()

	var lookupList []*_go.Reason
	lCount := 0
	next := false
	// Check if we want to fetch all records
	//
	// If the size is -1, we want to fetch all records
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		// If not fetching all records, check the size limit
		if !fetchAll && lCount >= rpc.GetSize() {
			next = true
			break
		}

		l := &_go.Reason{}

		var (
			createdBy, updatedBy         _go.Lookup
			tempCreatedAt, tempUpdatedAt time.Time
		)

		scanArgs := s.buildScanArgs(rpc.Fields, l, &createdBy, &updatedBy, &tempCreatedAt, &tempUpdatedAt)
		if err := rows.Scan(scanArgs...); err != nil {
			return nil, model.NewInternalError("postgres.reason.list.row_scan_error", err.Error())
		}

		s.populateReasonFields(rpc.Fields, l, &createdBy, &updatedBy, tempCreatedAt, tempUpdatedAt)
		lookupList = append(lookupList, l)
		lCount++
	}

	return &_go.ReasonList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

// Delete implements store.ReasonStore.
func (s *Reason) Delete(rpc *model.DeleteOptions, closeReasonId int64) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.reason.delete.database_connection_error", dbErr.Error())
	}

	query, args, err := s.buildDeleteReasonQuery(rpc)
	if err != nil {
		return model.NewInternalError("postgres.reason.delete.query_build_error", err.Error())
	}

	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return model.NewInternalError("postgres.reason.delete.execution_error", err.Error())
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return model.NewNotFoundError("postgres.reason.delete.no_rows_affected", "No rows affected for deletion")
	}

	return nil
}

// Update implements store.ReasonStore.
func (s *Reason) Update(rpc *model.UpdateOptions, l *_go.Reason) (*_go.Reason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.reason.update.database_connection_error", dbErr.Error())
	}

	query, args, err := s.buildUpdateReasonQuery(rpc, l)
	if err != nil {
		return nil, model.NewInternalError("postgres.reason.update.query_build_error", err.Error())
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
		return nil, model.NewInternalError("postgres.reason.update.execution_error", err.Error())
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedBy

	return l, nil
}

// buildCreateCloseReasonLookupQuery constructs the SQL insert query and returns the query string and arguments.
func (s Reason) buildCreateReasonQuery(rpc *model.CreateOptions, lookup *_go.Reason) (string, []interface{}, error) {
	query := createReasonQuery
	args := []interface{}{
		lookup.Name,
		rpc.Session.GetDomainId(),
		rpc.Time,
		lookup.Description,
		rpc.Session.GetUserId(),
		lookup.CloseReasonId,
	}
	return query, args, nil
}

// buildSearchCloseReasonLookupQuery constructs the SQL search query and returns the query builder.
func (s Reason) buildSearchReasonQuery(rpc *model.SearchOptions, closeReasonId int64) (string, []interface{}, error) {
	queryBuilder := sq.Select().
		From("cases.reason AS g").
		Where(sq.Eq{"g.dc": rpc.Session.GetDomainId(), "g.close_reason_id": closeReasonId}).
		PlaceholderFormat(sq.Dollar)

	fields := rpc.FieldsUtil.FieldsFunc(rpc.Fields, rpc.FieldsUtil.InlineFields)
	rpc.Fields = append(fields, "id")

	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
		case "description":
			// Use COALESCE to handle null values for description
			queryBuilder = queryBuilder.Column("COALESCE(g.description, '') AS description")
		case "created_by":
			// cbi = created_by_id
			// cbn = created_by_name
			// Use COALESCE to handle null values
			queryBuilder = queryBuilder.
				Column("COALESCE(created_by.id, 0) AS cbi").    // Handle NULL as 0 for created_by_id
				Column("COALESCE(created_by.name, '') AS cbn"). // Handle NULL as '' for created_by_name
				LeftJoin("directory.wbt_auth AS created_by ON p.created_by = created_by.id")
		case "updated_by":
			// ubi = updated_by_id
			// ubn = updated_by_name
			// Use COALESCE to handle null values
			queryBuilder = queryBuilder.
				Column("COALESCE(updated_by.id, 0) AS ubi").    // Handle NULL as 0 for updated_by_id
				Column("COALESCE(updated_by.name, '') AS ubn"). // Handle NULL as '' for updated_by_name
				LeftJoin("directory.wbt_auth AS updated_by ON p.updated_by = updated_by.id")
		}
	}

	convertedIds := rpc.FieldsUtil.Int64SliceToStringSlice(rpc.IDs)
	ids := rpc.FieldsUtil.FieldsFunc(convertedIds, rpc.FieldsUtil.InlineFields)

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substrs := rpc.Match.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": combinedLike})
	}

	parsedFields := rpc.FieldsUtil.FieldsFunc(rpc.Sort, rpc.FieldsUtil.InlineFields)
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
			column = "g." + sortField
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

	queryBuilder = queryBuilder.OrderBy(sortFields...)

	size := rpc.GetSize()
	page := rpc.Page

	if rpc.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	if rpc.GetSize() != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.reason.query_build_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

// buildDeleteCloseReasonLookupQuery constructs the SQL delete query and returns the query string and arguments.
func (s Reason) buildDeleteReasonQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := rpc.FieldsUtil.Int64SliceToStringSlice(rpc.IDs)
	ids := rpc.FieldsUtil.FieldsFunc(convertedIds, rpc.FieldsUtil.InlineFields)

	query := deleteReasonQuery
	args := []interface{}{pq.Array(ids), rpc.Session.GetDomainId()}
	return query, args, nil
}

func (s Reason) buildUpdateReasonQuery(rpc *model.UpdateOptions, l *_go.Reason) (string, []interface{}, error) {
	// Initialize Squirrel builder
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Create a Squirrel update builder
	updateBuilder := psql.Update("cases.reason").
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId())

	// Add the fields to the update query if they are provided
	for _, field := range rpc.Fields {
		switch field {
		case "name":
			if l.Name != "" {
				updateBuilder = updateBuilder.Set("name", l.Name)
			}
		case "description":
			// Use NULLIF to store NULL if description is an empty string
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", l.Description))
		}
	}

	// Add the WHERE clause for id and dc
	updateBuilder = updateBuilder.Where(sq.Eq{"id": l.Id, "dc": rpc.Session.GetDomainId()})

	// Build the SQL string and the arguments slice
	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Construct the final SQL query with joins for created_by and updated_by
	query := fmt.Sprintf(`
WITH upd AS (%s
    RETURNING id, name, created_at, updated_at, description, created_by, updated_by)
SELECT upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       COALESCE(upd.description, '')          AS description, -- Use COALESCE to return '' if description is NULL
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

// buildScanArgs prepares the arguments for scanning SQL rows.
func (s Reason) buildScanArgs(fields []string, r *_go.Reason, createdBy, updatedBy *_go.Lookup, tempCreatedAt, tempUpdatedAt *time.Time) []interface{} {
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
func (s Reason) populateReasonFields(fields []string, r *_go.Reason, createdBy, updatedBy *_go.Lookup, tempCreatedAt, tempUpdatedAt time.Time) {
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
func (s Reason) containsField(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

var (
	createReasonQuery = store.CompactSQL(`
	WITH ins AS (
		INSERT INTO cases.reason (name, dc, created_at, description, created_by, updated_at, updated_by, close_reason_id)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $3, $5, $6)  -- Use NULLIF to set NULL if description is ''
		RETURNING *
	)
	SELECT ins.id,
		   ins.name,
		   ins.created_at,
		   COALESCE(ins.description, '')      AS description,  -- Use COALESCE to return '' if description is NULL
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

	deleteReasonQuery = store.CompactSQL(`
DELETE FROM cases.reason
WHERE id = ANY($1) AND dc = $2
`)
)

func NewReasonStore(store store.Store) (store.ReasonStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_reason.check.bad_arguments",
			"error creating reason interface to the status_condition table, main store is nil")
	}
	return &Reason{storage: store}, nil
}
