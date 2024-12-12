package postgres

import (
	"fmt"
	"log"
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

type Status struct {
	storage store.Store
}

// Create creates a new status in the database. Implements the store.StatusStore interface.
func (s Status) Create(rpc *model.CreateOptions, add *_go.Status) (*_go.Status, error) {
	d, dbErr := s.storage.Database()

	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.create.database_connection_error", dbErr)
	}

	query, args, err := s.buildCreateStatusQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.create.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup _go.Lookup
		createdAt, updatedAt             time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.create.execution_error", err)
	}

	t := rpc.Time

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
func (s Status) List(rpc *model.SearchOptions) (*_go.StatusList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.list.database_connection_error", dbErr)
	}

	query, args, err := s.buildSearchStatusQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.list.query_build_error", err)
	}

	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.list.execution_error", err)
	}
	defer rows.Close()

	var lookupList []*_go.Status

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

		l := &_go.Status{}

		var (
			createdBy, updatedBy         _go.Lookup
			tempUpdatedAt, tempCreatedAt time.Time
			scanArgs                     []interface{}
		)

		// Prepare scan arguments based on requested fields
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
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			log.Printf("postgres.cases.status.list Failed to scan row: %v", err)
			return nil, dberr.NewDBInternalError("postgres.cases.status.list.row_scan_error", err)
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

	return &_go.StatusList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

// Delete removes a status from the database. Implements the store.StatusStore interface.
func (s Status) Delete(rpc *model.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.cases.status.delete.database_connection_error", dbErr)
	}

	query, args, err := s.buildDeleteStatusQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.status.delete.query_build_error", err)
	}

	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		log.Printf("postgres.cases.status.delete Failed to execute SQL query: %v", err)
		return dberr.NewDBInternalError("postgres.cases.status.delete.execution_error", err)
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBNoRowsError("postgres.cases.status.delete.no_rows_affected")
	}

	return nil
}

// Update modifies a status in the database. Implements the store.StatusStore interface.
func (s Status) Update(rpc *model.UpdateOptions, l *_go.Status) (*_go.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.update.database_connection_error", dbErr)
	}

	query, args, queryErr := s.buildUpdateStatusQuery(rpc, l)

	if queryErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.status.update.query_build_error", queryErr)
	}

	var (
		createdBy, updatedByLookup _go.Lookup
		createdAt, updatedAt       time.Time
	)

	err := d.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt, &l.Description,
		&createdBy.Id, &createdBy.Name, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		log.Printf("postgres.cases.status.update Failed to execute SQL query: %v", err)
		return nil, dberr.NewDBInternalError("postgres.cases.status.update.execution_error", err)
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedByLookup

	return l, nil
}

// buildCreateStatusLookupQuery constructs the SQL insert query and returns the query string and arguments.
func (s Status) buildCreateStatusQuery(rpc *model.CreateOptions, lookup *_go.Status) (string, []interface{}, error) {
	args := []interface{}{
		lookup.Name,               // $1 - name
		rpc.Session.GetDomainId(), // $2 - dc
		rpc.Time,                  // $3 - created_at / updated_at
		lookup.Description,        // $4 - description
		rpc.Session.GetUserId(),   // $5 - created_by / updated_by
	}
	return createStatusQuery, args, nil
}

func (s Status) buildSearchStatusQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	queryBuilder := sq.Select().
		From("cases.status AS g").
		Where(sq.Eq{"g.dc": rpc.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	fields := util.FieldsFunc(rpc.Fields, util.InlineFields)
	rpc.Fields = append(fields, "id")

	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
		case "description":
			// Use COALESCE to handle null values for description
			queryBuilder = queryBuilder.Column("COALESCE(g.description, '') AS description")
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

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substrs := util.Substring(name)
		combinedLike := strings.Join(substrs, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": combinedLike})
	}

	if len(rpc.Sort) > 0 {
		parsedFields := util.FieldsFunc(rpc.Sort, util.InlineFields)

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
	} else {
		// -------- Apply [Sorting by Name] --------
		queryBuilder = queryBuilder.OrderBy("g.name ASC")
	}

	size := rpc.GetSize()
	page := rpc.GetPage()

	// Apply offset only if page > 1
	if rpc.Page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * size))
	}

	// Apply limit
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1)) // Request one more record to check if there's a next page
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.cases.status.query_build.sql_generation_error", err)
	}

	return store.CompactSQL(query), args, nil
}

// buildDeleteStatusLookupQuery constructs the SQL delete query and returns the query string and arguments.
func (s Status) buildDeleteStatusQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	args := []interface{}{
		pq.Array(ids),             // $1 - id
		rpc.Session.GetDomainId(), // $2 - dc
	}
	return deleteStatusQuery, args, nil
}

func (s Status) buildUpdateStatusQuery(rpc *model.UpdateOptions, l *_go.Status) (string, []interface{}, error) {
	// Initialize Squirrel builder
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Create a Squirrel update builder
	updateBuilder := psql.Update("cases.status").
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
WITH upd AS (
    %s
    RETURNING id, name, created_at, updated_at, description, created_by, updated_by
)
SELECT upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       COALESCE(upd.description, '')      AS description,  -- Use COALESCE to return '' if description is NULL
       upd.created_by                     AS created_by_id,
       COALESCE(c.name::text, c.username, '') AS created_by_name,
       upd.updated_by                     AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name
FROM upd
LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
LEFT JOIN directory.wbt_user c ON c.id = upd.created_by;
    `, sql)

	return store.CompactSQL(query), args, nil
}

// ---- STATIC SQL QUERIES ----
var (
	createStatusQuery = store.CompactSQL(`
	WITH ins AS (
		INSERT INTO cases.status (name, dc, created_at, description, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $3, $5)  -- Use NULLIF to set NULL if description is ''
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
		   COALESCE(u.name::text, u.username) AS updated_by_name
	FROM ins
	LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
	LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;
	`)

	deleteStatusQuery = store.CompactSQL(`
DELETE FROM cases.status
WHERE id = ANY($1) AND dc = $2
`)
)

func NewStatusStore(store store.Store) (store.StatusStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_status.check.bad_arguments",
			"error creating stuas interface to the status table, main store is nil")
	}
	return &Status{storage: store}, nil
}
