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
	db "github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

const (
	sourceDefaultSort = "name"
)

type Source struct {
	storage db.Store
}

func (s Source) Create(rpc *model.CreateOptions, add *_go.Source) (*_go.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.create.database_connection_error", dbErr)
	}

	query, args, err := s.buildCreateSourceQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.create.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup _go.Lookup
		createdAt, updatedAt             time.Time
		tempType                         string
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description, &tempType,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.create.execution_error", err)
	}

	// Convert tempType (string) to the enum Type
	add.Type, err = stringToType(tempType)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.create.type_conversion_error", err)
	}

	return &_go.Source{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		Type:        add.Type,
		CreatedAt:   util.Timestamp(createdAt),
		UpdatedAt:   util.Timestamp(updatedAt),
		CreatedBy:   &createdByLookup,
		UpdatedBy:   &updatedByLookup,
	}, nil
}

func (s Source) List(rpc *model.SearchOptions) (*_go.SourceList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.list.database_connection_error", dbErr)
	}

	query, args, err := s.buildSearchSourceQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.list.query_build_error", err)
	}

	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.list.execution_error", err)
	}
	defer rows.Close()

	var lookupList []*_go.Source
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

		l := &_go.Source{}
		var (
			createdBy, updatedBy         _go.Lookup
			tempUpdatedAt, tempCreatedAt time.Time
			tempType                     string
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
			case "type":
				scanArgs = append(scanArgs, &tempType)
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

		if scanErr := rows.Scan(scanArgs...); scanErr != nil {
			return nil, dberr.NewDBInternalError("postgres.cases.source.list.row_scan_error", err)
		}

		// Convert tempType (string) to the enum Type if "type" is in the requested fields
		if util.ContainsField(rpc.Fields, "type") {
			l.Type, err = stringToType(tempType)
			if err != nil {
				return nil, dberr.NewDBInternalError("postgres.cases.source.list.type_conversion_error", err)
			}
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

	// If fetching all records, set `next` to false as there's no pagination
	if fetchAll {
		next = false
	}

	return &_go.SourceList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

func (s Source) Delete(rpc *model.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.cases.source.delete.database_connection_error", dbErr)
	}

	query, args, err := s.buildDeleteSourceQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.source.delete.query_build_error", err)
	}

	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.source.delete.execution_error", err)
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBError("postgres.cases.source.delete.no_rows_affected", "No rows affected for deletion")
	}

	return nil
}

func (s Source) Update(rpc *model.UpdateOptions, l *_go.Source) (*_go.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.update.database_connection_error", dbErr)
	}

	query, args, queryErr := s.buildUpdateSourceQuery(rpc, l)
	if queryErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.update.query_build_error", queryErr)
	}

	var (
		createdBy, updatedByLookup _go.Lookup
		createdAt, updatedAt       time.Time
		tempType                   string
	)

	err := d.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt, &l.Description, &tempType,
		&createdBy.Id, &createdBy.Name, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.update.execution_error", err)
	}

	// Convert tempType (string) to the enum Type
	l.Type, err = stringToType(tempType)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.source.update.type_conversion_error", err)
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedByLookup

	return l, nil
}

func (s Source) buildCreateSourceQuery(rpc *model.CreateOptions, lookup *_go.Source) (string, []interface{}, error) {
	query := createSourceQuery
	args := []interface{}{
		lookup.Name, rpc.GetAuthOpts().GetDomainId(), rpc.CurrentTime(), lookup.Description, lookup.Type,
		rpc.GetAuthOpts().GetUserId(),
	}
	return query, args, nil
}

func (s Source) buildSearchSourceQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	queryBuilder := sq.Select().
		From("cases.source AS g").
		Where(sq.Eq{"g.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	fields := util.FieldsFunc(rpc.Fields, util.InlineFields)
	rpc.Fields = append(fields, "id")

	// Adding columns based on fields
	for _, field := range rpc.Fields {
		switch field {
		case "id", "name", "type", "created_at", "updated_at", "source":
			queryBuilder = queryBuilder.Column("g." + field)
		case "description":
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

	// Applying filters
	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": ids})
	}

	if name, ok := rpc.Filter["name"].(string); ok && len(name) > 0 {
		substr := util.Substring(name)
		combinedLike := strings.Join(substr, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": combinedLike})
	}

	if types, ok := rpc.Filter["type"].([]_go.Type); ok && len(types) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.type": types})
	}

	// -------- Apply sorting ----------
	queryBuilder = store.ApplyDefaultSorting(rpc, queryBuilder, sourceDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = store.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Generate SQL and arguments
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return "", nil, dberr.NewDBInternalError("postgres.cases.source.query_build.sql_generation_error", err)
	}

	return db.CompactSQL(query), args, nil
}

func (s Source) buildDeleteSourceQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deleteSourceQuery
	args := []interface{}{pq.Array(ids), rpc.GetAuthOpts().GetDomainId()}
	return query, args, nil
}

func (s Source) buildUpdateSourceQuery(rpc *model.UpdateOptions, l *_go.Source) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Update("cases.source").
		Set("updated_at", rpc.CurrentTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": l.Id}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()})

	for _, field := range rpc.Fields {
		switch field {
		case "name":
			if l.Name != "" {
				builder = builder.Set("name", l.Name)
			}
		case "description":
			// Use NULLIF to store NULL when the description is an empty string
			builder = builder.Set("description", sq.Expr("NULLIF(?, '')", l.Description))
		case "type":
			if l.Type != _go.Type_TYPE_UNSPECIFIED {
				builder = builder.Set("type", l.Type)
			}
		}
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	q := fmt.Sprintf(`
WITH upd AS (%s
	RETURNING id, name, created_at, updated_at, description, type, created_by, updated_by)
SELECT upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       COALESCE(upd.description, '')          AS description, -- Return '' if description is NULL
       upd.type,
       upd.created_by                         AS created_by_id,
       COALESCE(c.name::text, c.username, '') AS created_by_name,
       upd.updated_by                         AS updated_by_id,
       COALESCE(u.name::text, u.username)     AS updated_by_name
FROM upd
         LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = upd.created_by;
`, sql)

	return db.CompactSQL(q), args, nil
}

// StringToType converts a string into the corresponding Type enum value.
//
// Types are specified ONLY for Source dictionary and are ENUMS in API.
func stringToType(typeStr string) (_go.Type, error) {
	switch strings.ToUpper(typeStr) {
	case "CALL":
		return _go.Type_CALL, nil
	case "CHAT":
		return _go.Type_CHAT, nil
	case "SOCIAL_MEDIA":
		return _go.Type_SOCIAL_MEDIA, nil
	case "EMAIL":
		return _go.Type_EMAIL, nil
	case "API":
		return _go.Type_API, nil
	case "MANUAL":
		return _go.Type_MANUAL, nil
	default:
		return _go.Type_TYPE_UNSPECIFIED, fmt.Errorf("invalid type value: %s", typeStr)
	}
}

var (
	createSourceQuery = db.CompactSQL(`WITH ins AS (
    INSERT INTO cases.source (name, dc, created_at, description, type, created_by, updated_at, updated_by)
        VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6, $3, $6) -- Use NULLIF to set NULL if description is ''
        RETURNING id, name, created_at, description, type, created_by, updated_at, updated_by)
SELECT ins.id,
       ins.name,
       ins.created_at,
       COALESCE(ins.description, '')      AS description,
       ins.type::text,
       ins.created_by                     AS created_by_id,
       COALESCE(c.name::text, c.username) AS created_by_name,
       ins.updated_at,
       ins.updated_by                     AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name
FROM ins
         LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
         LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;`)

	deleteSourceQuery = db.CompactSQL(
		`DELETE FROM cases.source
    WHERE id = ANY($1) AND dc = $2 `,
	)
)

func NewSourceStore(store db.Store) (db.SourceStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_source.check.bad_arguments",
			"error creating source interface to the source table, main store is nil")
	}
	return &Source{storage: store}, nil
}
