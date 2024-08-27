package postgres

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	api "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type Priority struct {
	storage store.Store
}

// Create implements store.PriorityStore.
func (p *Priority) Create(ctx *model.CreateOptions, add *api.Priority) (*api.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.cases.priority.create.database_connection_error", dbErr.Error())
	}

	query, args, err := p.buildCreatePriorityQuery(ctx, add)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.priority.create.query_build_error", err.Error())
	}

	var createdByLookup, updatedByLookup api.Lookup
	var createdAt, updatedAt time.Time

	err = d.QueryRow(ctx.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.close_reason.create.execution_error", err.Error())
	}

	return &api.Priority{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		CreatedAt:   util.Timestamp(createdAt),
		UpdatedAt:   util.Timestamp(updatedAt),
		CreatedBy:   &createdByLookup,
		UpdatedBy:   &updatedByLookup,
	}, nil
}

// Delete implements store.PriorityStore.
func (p *Priority) Delete(ctx *model.DeleteOptions) error {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return model.NewInternalError("postgres.cases.priority.delete.database_connection_error", dbErr.Error())
	}

	query, args, err := p.buildDeletePriorityQuery(ctx)
	if err != nil {
		return model.NewInternalError("postgres.cases.priority.delete.query_build_error", err.Error())
	}

	res, err := d.Exec(ctx.Context, query, args...)
	if err != nil {
		return model.NewInternalError("postgres.cases.priority.delete.execution_error", err.Error())
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return model.NewNotFoundError("postgres.cases.priority.delete.no_rows_affected", "No rows affected for deletion")
	}

	return nil
}

// List implements store.PriorityStore.
func (p *Priority) List(ctx *model.SearchOptions) (*api.PriorityList, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.cases.priority.list.database_connection_error", dbErr.Error())
	}

	query, args, err := p.buildSearchPriorityQuery(ctx)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.priority.list.query_build_error", err.Error())
	}

	rows, err := d.Query(ctx.Context, query, args...)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.priority.list.execution_error", err.Error())
	}
	defer rows.Close()

	var lookupList []*api.Priority
	lCount := 0
	next := false
	for rows.Next() {
		if lCount >= ctx.GetSize() {
			next = true
			break
		}

		l := &api.Priority{}
		var createdBy, updatedBy api.Lookup
		var tempUpdatedAt, tempCreatedAt time.Time
		var scanArgs []interface{}

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
			case "color":
				scanArgs = append(scanArgs, &l.Color)
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, model.NewInternalError("postgres.cases.close_reason.list.row_scan_error", err.Error())
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

	return &api.PriorityList{
		Page:  int32(ctx.Page), // TODO page should be correctly passes even if user pass negative value
		Next:  next,
		Items: lookupList,
	}, nil
}

// Update implements store.PriorityStore.
func (p *Priority) Update(ctx *model.UpdateOptions, l *api.Priority) (*api.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, model.NewInternalError("postgres.cases.priority.update.database_connection_error", dbErr.Error())
	}

	query, args, queryErr := p.buildUpdatePriorityQuery(ctx, l)
	if queryErr != nil {
		return nil, model.NewInternalError("postgres.cases.priority.update.query_build_error", queryErr.Error())
	}

	var createdBy, updatedByLookup api.Lookup
	var createdAt, updatedAt time.Time

	err := d.QueryRow(ctx.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt,
		&l.Description, &createdBy.Id, &createdBy.Name, &updatedByLookup.Id,
		&updatedByLookup.Name, &l.Color,
	)
	if err != nil {
		return nil, model.NewInternalError("postgres.cases.priority.update.execution_error", err.Error())
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedByLookup

	return l, nil
}

func (s Priority) buildSearchPriorityQuery(ctx *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	queryBuilder := sq.Select().
		From("cases.priority AS p").
		Where(sq.Eq{"p.dc": ctx.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	fields := ctx.FieldsUtil.FieldsFunc(ctx.Fields, ctx.FieldsUtil.InlineFields)
	ctx.Fields = append(fields, "id")

	for _, field := range ctx.Fields {
		switch field {
		case "id", "name", "description", "created_at", "updated_at", "color":
			queryBuilder = queryBuilder.Column("p." + field)
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
		queryBuilder = queryBuilder.Where(sq.Eq{"p.id": ids})
	}

	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substr := ctx.Match.Substring(name)
		queryBuilder = queryBuilder.Where(sq.ILike{"p.name": substr})
	}

	parsedFields := ctx.FieldsUtil.FieldsFunc(ctx.Sort, ctx.FieldsUtil.InlineFields)
	var sortFields []string

	for _, sortField := range parsedFields {
		desc := false
		if strings.HasPrefix(sortField, "!") {
			desc = true
			sortField = strings.TrimPrefix(sortField, "!")
		}

		column := "g." + sortField
		if desc {
			column += " DESC"
		} else {
			column += " ASC"
		}
		sortFields = append(sortFields, column)
	}

	queryBuilder = queryBuilder.OrderBy(sortFields...)

	size := ctx.GetSize()
	page := ctx.Page

	if ctx.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	if ctx.GetSize() != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, model.NewInternalError("postgres.cases.priority.query_build.sql_generation_error", err.Error())
	}

	return store.CompactSQL(query), args, nil
}

func (s Priority) buildUpdatePriorityQuery(ctx *model.UpdateOptions, l *api.Priority) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Initialize the SQL builder
	builder := psql.Update("cases.priority").
		Set("updated_at", ctx.Time).
		Set("updated_by", ctx.Session.GetUserId()).
		Where(sq.Eq{"id": l.Id}).
		Where(sq.Eq{"dc": ctx.Session.GetDomainId()})

	// Fields that could be updated
	updateFields := []string{"name", "description"} // TODO make it empty  |  add XJsonMask to proto

	// Add the fields to the update statement
	for _, field := range updateFields {
		switch field {
		case "name":
			if l.Name != "" {
				builder = builder.Set("name", l.Name)
			}
		case "description":
			if l.Description != "" {
				builder = builder.Set("description", l.Description)
			}
		case "color":
			if l.Color != "" {
				builder = builder.Set("color", l.Color)
			}
		}
	}

	// Generate SQL and arguments from the builder
	sql, args, err := builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Construct the final SQL query with joins for created_by and updated_by
	q := fmt.Sprintf(`
WITH upd AS (
	%s
	RETURNING id, name, created_at, updated_at, description, created_by, updated_by, color
)
SELECT upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       upd.description,
       upd.created_by AS created_by_id,
       COALESCE(c.name::text, c.username) AS created_by_name,
       upd.updated_by AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name
FROM upd
LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
LEFT JOIN directory.wbt_user c ON c.id = upd.created_by;
`, sql)

	return store.CompactSQL(q), args, nil
}

// buildDeleteCloseReasonLookupQuery constructs the SQL delete query and returns the query string and arguments.
func (s Priority) buildDeletePriorityQuery(ctx *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	query := deletePriorityQuery
	args := []interface{}{pq.Array(ids), ctx.Session.GetDomainId()}
	return query, args, nil
}

// buildCreatePriorityQuery constructs the SQL insert query and returns the query string and arguments.
func (s Priority) buildCreatePriorityQuery(ctx *model.CreateOptions, lookup *api.Priority) (string, []interface{}, error) {
	query := createPriorityQuery
	args := []interface{}{
		lookup.Name,
		ctx.Session.GetDomainId(),
		ctx.Time,
		lookup.Description,
		ctx.Session.GetUserId(),
		lookup.Color,
	}
	return query, args, nil
}

var (
	createPriorityQuery = store.CompactSQL(`
	with ins as (
		INSERT INTO cases.priority (name, dc, created_at, description, created_by, updated_at, updated_by, color)
		VALUES ($1, $2, $3, $4, $5, $3, $5, $6)
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

	deletePriorityQuery = store.CompactSQL(`
	DELETE FROM cases.priority
	WHERE id = ANY($1) AND dc = $2
`)
)

func NewPriorityStore(store store.Store) (store.PriorityStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.new_priority.check.bad_arguments",
			"error creating priority interface to the status_condition table, main store is nil")
	}
	return &Priority{storage: store}, nil
}
