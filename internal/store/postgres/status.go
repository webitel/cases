package postgres

import (
	"fmt"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type Status struct {
	storage store.Store
}

// Create creates a new status in the database. Implements the store.StatusStore interface.
func (s Status) Create(ctx *model.CreateOptions, add *_go.Status) (*_go.Status, error) {
	d, dbErr := s.storage.Database()

	if dbErr != nil {
		return nil, model.NewInternalError("postgres.cases.status.create.database_connection_error", dbErr.Error())
	}

	query, args, err := s.buildCreateStatusQuery(ctx, add)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.status.create.query_build_error", err.Error())
	}

	var createdByLookup, updatedByLookup _go.Lookup
	var createdAt, updatedAt time.Time

	err = d.QueryRow(ctx.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.status.create.execution_error", err.Error())
	}

	t := ctx.Time

	return &_go.Status{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		CreatedAt:   util.Timestamp(t),
		UpdatedAt:   util.Timestamp(t),
		CreatedBy:   &createdByLookup,
		UpdatedBy:   &updatedByLookup,
	}, nil
}

// List retrieves a list of statuses from the database. Implements the store.StatusStore interface.
func (s Status) List(ctx *model.SearchOptions) (*_go.StatusList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.cases.status.list.database_connection_error", dbErr.Error())
	}

	query, args, err := s.buildSearchStatusQuery(ctx)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.status.list.query_build_error", err.Error())
	}

	rows, err := d.Query(ctx.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.status.list.execution_error", err.Error())
	}
	defer rows.Close()

	var lookupList []*_go.Status

	lCount := 0
	next := false
	for rows.Next() {
		if lCount >= ctx.GetSize() {
			next = true
			break
		}

		l := &_go.Status{}
		var createdBy, updatedBy _go.Lookup
		var tempUpdatedAt, tempCreatedAt time.Time
		var scanArgs []interface{}

		// Prepare scan arguments based on requested fields
		for _, field := range ctx.Fields {
			switch field {
			case "id":
				scanArgs = append(scanArgs, &l.Id)
			case "name":
				scanArgs = append(scanArgs, &l.Name)
			case "description":
				scanArgs = append(scanArgs, &l.Description)
			case "created_at":
				scanArgs = append(scanArgs, &tempCreatedAt)
			case "updated_at":
				scanArgs = append(scanArgs, &tempUpdatedAt)
			case "created_by":
				scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
			case "updated_by":
				scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			log.Printf("postgres.cases.status.list Failed to scan row: %v", err)
			return nil, model.NewInternalError("postgres.cases.status.list.row_scan_error", err.Error())
		}

		if ctx.FieldsUtil.ContainsField(ctx.Fields, "created_by") {
			l.CreatedBy = &createdBy
		}
		if ctx.FieldsUtil.ContainsField(ctx.Fields, "updated_by") {
			l.UpdatedBy = &updatedBy
		}

		if ctx.FieldsUtil.ContainsField(ctx.Fields, "created_at") {
			l.CreatedAt = util.Timestamp(tempCreatedAt)
		}
		if ctx.FieldsUtil.ContainsField(ctx.Fields, "updated_at") {
			l.UpdatedAt = util.Timestamp(tempUpdatedAt)
		}

		lookupList = append(lookupList, l)
		lCount++
	}

	return &_go.StatusList{
		Page:  int32(ctx.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

// Delete removes a status from the database. Implements the store.StatusStore interface.
func (s Status) Delete(ctx *model.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.cases.status.delete.database_connection_error", dbErr.Error())
	}

	query, args, err := s.buildDeleteStatusQuery(ctx)
	if err != nil {
		return model.NewInternalError("postgres.cases.status.delete.query_build_error", err.Error())
	}

	res, err := d.Exec(ctx.Context, query, args...)
	if err != nil {
		log.Printf("postgres.cases.status.delete Failed to execute SQL query: %v", err)
		return model.NewInternalError("postgres.cases.status.delete.execution_error", err.Error())
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return model.NewNotFoundError("postgres.cases.status.delete.no_rows_affected", "No rows affected for deletion")
	}

	return nil
}

// Update modifies a status in the database. Implements the store.StatusStore interface.
func (s Status) Update(ctx *model.UpdateOptions, l *_go.Status) (*_go.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.cases.status.update.database_connection_error", dbErr.Error())
	}

	query, args, queryErr := s.buildUpdateStatusQuery(ctx, l)

	if queryErr != nil {
		return nil, model.NewInternalError("postgres.cases.status.update.query_build_error", queryErr.Error())
	}

	var createdBy, updatedByLookup _go.Lookup
	var createdAt, updatedAt time.Time

	err := d.QueryRow(ctx.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt, &l.Description,
		&createdBy.Id, &createdBy.Name, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		log.Printf("postgres.cases.status.update Failed to execute SQL query: %v", err)
		return nil, model.NewInternalError("postgres.cases.status.update.execution_error", err.Error())
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedByLookup

	return l, nil
}

// buildCreateStatusLookupQuery constructs the SQL insert query and returns the query string and arguments.
func (s Status) buildCreateStatusQuery(ctx *model.CreateOptions, lookup *_go.Status) (string, []interface{}, error) {
	args := []interface{}{
		lookup.Name,
		ctx.Session.GetDomainId(),
		ctx.Time,
		lookup.Description,
		ctx.Session.GetUserId(),
	}
	return createStatusQuery, args, nil
}

func (s Status) buildSearchStatusQuery(ctx *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	queryBuilder := sq.Select().
		From("cases.status AS g").
		Where(sq.Eq{"g.dc": ctx.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	fields := ctx.FieldsUtil.FieldsFunc(ctx.Fields, ctx.FieldsUtil.InlineFields)
	ctx.Fields = append(fields, "id")

	for _, field := range ctx.Fields {
		switch field {
		case "id", "name", "description", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
		case "created_by":
			// cbi = created_by_id
			// cbn = created_by_name
			queryBuilder = queryBuilder.Column("created_by.id AS cbi, created_by.name AS cbn").
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
		case "updated_by":
			// ubi = updated_by_id
			// ubn = updated_by_name
			queryBuilder = queryBuilder.Column("updated_by.id AS ubi, updated_by.name AS ubn").
				LeftJoin("directory.wbt_auth AS updated_by ON g.updated_by = updated_by.id")
		}
	}

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substrs := ctx.Match.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": "%" + combinedLike + "%"})
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

	// Apply sorting
	queryBuilder = queryBuilder.OrderBy(sortFields...)

	size := ctx.GetSize()
	page := ctx.GetPage()

	// Apply offset only if page > 1
	if ctx.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	// Apply limit
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1)) // Request one more record to check if there's a next page
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.cases.status.query_build.sql_generation_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

// buildDeleteStatusLookupQuery constructs the SQL delete query and returns the query string and arguments.
func (s Status) buildDeleteStatusQuery(ctx *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	args := []interface{}{pq.Array(ids), ctx.Session.GetDomainId()}
	return deleteStatusQuery, args, nil
}

func (s Status) buildUpdateStatusQuery(ctx *model.UpdateOptions, l *_go.Status) (string, []interface{}, error) {
	// Initialize Squirrel builder
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Create a Squirrel update builder
	updateBuilder := psql.Update("cases.status").
		Set("updated_at", ctx.Time).
		Set("updated_by", ctx.Session.GetUserId())

	// Fields that could be updated
	updateFields := map[string]string{
		"name":        l.Name,
		"description": l.Description,
	}

	// Add the fields to the update query if they are provided
	for field, value := range updateFields {
		if value != "" {
			updateBuilder = updateBuilder.Set(field, value)
		}
	}

	// Add the WHERE clause for id and dc
	updateBuilder = updateBuilder.Where(sq.Eq{"id": l.Id, "dc": ctx.Session.GetDomainId()})

	// Build the SQL string and the arguments slice
	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return "", nil, err
	}

	// Construct the final SQL query with joins for created_by and updated_by
	query := fmt.Sprintf(`
with upd as (%s
    RETURNING id, name, created_at, updated_at, description, created_by, updated_by
)
select upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       upd.description,
       upd.created_by as created_by_id,
       coalesce(c.name::text, c.username) as created_by_name,
       upd.updated_by as updated_by_id,
       coalesce(u.name::text, u.username) as updated_by_name
from upd
  left join directory.wbt_user u on u.id = upd.updated_by
  left join directory.wbt_user c on c.id = upd.created_by;
    `, sql)

	return store.CompactSQL(query), args, nil
}

// ---- STATIC SQL QUERIES ----
var (
	createStatusQuery = store.CompactSQL(`
with ins as (
    INSERT INTO cases.status (name, dc, created_at, description, created_by, updated_at,
updated_by)
    VALUES ($1, $2, $3, $4, $5, $3, $5)
    returning *
)
select ins.id,
    ins.name,
    ins.created_at,
    ins.description,
    ins.created_by created_by_id,
    coalesce(c.name::text, c.username) created_by_name,
    ins.updated_at,
    ins.updated_by updated_by_id,
    coalesce(u.name::text, u.username) updated_by_name
from ins
  left join directory.wbt_user u on u.id = ins.updated_by
  left join directory.wbt_user c on c.id = ins.created_by;
`)

	deleteStatusQuery = store.CompactSQL(`
DELETE FROM cases.status
WHERE id = ANY($1) AND dc = $2
`)
)

func NewStatusStore(store store.Store) (store.StatusStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_status.check.bad_arguments",
			"error creating stuas interface to the status table, main store is nil")
	}
	return &Status{storage: store}, nil
}
