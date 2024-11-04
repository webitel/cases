package postgres

import (
	"fmt"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	api "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type Priority struct {
	storage store.Store
}

// Create implements store.PriorityStore.
func (p *Priority) Create(rpc *model.CreateOptions, add *api.Priority) (*api.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.create.database_connection_error", dbErr)
	}

	query, args, err := p.buildCreatePriorityQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.create.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup api.Lookup
		createdAt, updatedAt             time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason.create.execution_error", err)
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
func (p *Priority) Delete(rpc *model.DeleteOptions) error {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.cases.priority.delete.database_connection_error", dbErr)
	}

	query, args, err := p.buildDeletePriorityQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.priority.delete.query_build_error", err)
	}

	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.priority.delete.execution_error", err)
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBNoRowsError("postgres.cases.priority.delete.no_rows_affected")
	}

	return nil
}

// List implements store.PriorityStore.
func (p *Priority) List(rpc *model.SearchOptions) (*api.PriorityList, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.list.database_connection_error", dbErr)
	}

	query, args, err := p.buildSearchPriorityQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.list.query_build_error", err)
	}

	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.list.execution_error", err)
	}
	defer rows.Close()

	var lookupList []*api.Priority
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

		l := &api.Priority{}

		var (
			createdBy, updatedBy         api.Lookup
			tempUpdatedAt, tempCreatedAt time.Time
			scanArgs                     []interface{}
		)

		for _, field := range rpc.Fields {
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
			return nil, dberr.NewDBInternalError("postgres.cases.close_reason.list.row_scan_error", err)
		}

		if util.ContainsField(rpc.Fields, "created_by") {
			l.CreatedBy = &createdBy
		}
		if util.ContainsField(rpc.Fields, "updated_by") {
			l.UpdatedBy = &updatedBy
		}
		if util.ContainsField(rpc.Fields, "created_at") {
			l.CreatedAt = util.Timestamp(tempCreatedAt)
		}
		if util.ContainsField(rpc.Fields, "updated_at") {
			l.UpdatedAt = util.Timestamp(tempUpdatedAt)
		}

		lookupList = append(lookupList, l)
		lCount++
	}

	return &api.PriorityList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

// Update implements store.PriorityStore.
func (p *Priority) Update(rpc *model.UpdateOptions, l *api.Priority) (*api.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.update.database_connection_error", dbErr)
	}

	query, args, queryErr := p.buildUpdatePriorityQuery(rpc, l)
	if queryErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.update.query_build_error", queryErr)
	}

	var (
		createdBy, updatedByLookup api.Lookup
		createdAt, updatedAt       time.Time
	)

	err := d.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt,
		&l.Description, &createdBy.Id, &createdBy.Name, &updatedByLookup.Id,
		&updatedByLookup.Name, &l.Color,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.priority.update.execution_error", err)
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedByLookup

	return l, nil
}

func (s Priority) buildSearchPriorityQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	queryBuilder := sq.Select().
		From("cases.priority AS p").
		Where(sq.Eq{"p.dc": rpc.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	fields := util.FieldsFunc(rpc.Fields, util.InlineFields)
	rpc.Fields = append(fields, "id")

	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "created_at", "updated_at", "color":
			queryBuilder = queryBuilder.Column("p." + field)
		case "description":
			// Use COALESCE to handle NULL values for description
			queryBuilder = queryBuilder.Column("COALESCE(p.description, '') AS description")
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

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"p.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substr := util.Substring(name)
		combinedLike := strings.Join(substr, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"p.name": combinedLike})
	}

	parsedFields := util.FieldsFunc(rpc.Sort, util.InlineFields)
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
		return "", nil, dberr.NewDBInternalError("postgres.cases.priority.query_build.sql_generation_error", err)
	}

	return store.CompactSQL(query), args, nil
}

func (s Priority) buildUpdatePriorityQuery(rpc *model.UpdateOptions, l *api.Priority) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Initialize the SQL builder
	builder := psql.Update("cases.priority").
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId()).
		Where(sq.Eq{"id": l.Id}).
		Where(sq.Eq{"dc": rpc.Session.GetDomainId()})

	// Add the fields to the update statement
	for _, field := range rpc.Fields {
		switch field {
		case "name":
			if l.Name != "" {
				builder = builder.Set("name", l.Name)
			}
		case "description":
			// Use NULLIF to store NULL if description is an empty string
			builder = builder.Set("description", sq.Expr("NULLIF(?, '')", l.Description))
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
       COALESCE(c.name::text, c.username, '') AS created_by_name,
       upd.updated_by AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name,
	   upd.color
FROM upd
LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
LEFT JOIN directory.wbt_user c ON c.id = upd.created_by;
`, sql)

	return store.CompactSQL(q), args, nil
}

// buildDeleteCloseReasonLookupQuery constructs the SQL delete query and returns the query string and arguments.
func (s Priority) buildDeletePriorityQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deletePriorityQuery
	args := []interface{}{pq.Array(ids), rpc.Session.GetDomainId()}
	return query, args, nil
}

// buildCreatePriorityQuery constructs the SQL insert query and returns the query string and arguments.
func (s Priority) buildCreatePriorityQuery(rpc *model.CreateOptions, lookup *api.Priority) (string, []interface{}, error) {
	query := createPriorityQuery
	args := []interface{}{
		lookup.Name,
		rpc.Session.GetDomainId(),
		rpc.Time,
		lookup.Description,
		rpc.Session.GetUserId(),
		lookup.Color,
	}
	return query, args, nil
}

var (
	createPriorityQuery = store.CompactSQL(`
	WITH ins AS (
		INSERT INTO cases.priority (name, dc, created_at, description, created_by, updated_at, updated_by, color)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $3, $5, $6) -- Use NULLIF to set NULL if description is ''
		RETURNING *
	)
	SELECT ins.id,
		   ins.name,
		   ins.created_at,
		   COALESCE(ins.description, '')      AS description, -- Use COALESCE to return '' if description is NULL
		   ins.created_by                     AS created_by_id,
		   COALESCE(c.name::text, c.username) AS created_by_name,
		   ins.updated_at,
		   ins.updated_by                     AS updated_by_id,
		   COALESCE(u.name::text, u.username) AS updated_by_name
	FROM ins
	LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
	LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;
	`)

	deletePriorityQuery = store.CompactSQL(`
	DELETE FROM cases.priority
	WHERE id = ANY($1) AND dc = $2
`)
)

func NewPriorityStore(store store.Store) (store.PriorityStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_priority.check.bad_arguments",
			"error creating priority interface to the status_condition table, main store is nil")
	}
	return &Priority{storage: store}, nil
}
