package postgres

import (
	"fmt"
	"strings"
	"time"

	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	general "buf.build/gen/go/webitel/general/protocolbuffers/go"
	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	dberr "github.com/webitel/cases/internal/error"
	db "github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"

	util "github.com/webitel/cases/util"
)

type CloseReasonGroup struct {
	storage db.Store
}

func (s CloseReasonGroup) Create(rpc *model.CreateOptions, add *_go.CloseReasonGroup) (*_go.CloseReasonGroup, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.create.database_connection_error", dbErr)
	}

	query, args, err := s.buildCreateCloseReasonGroupQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.create.query_build_error", err)
	}

	var (
		createdByLookup, updatedByLookup general.Lookup
		createdAt, updatedAt             time.Time
	)

	err = d.QueryRow(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.create.execution_error", err)
	}

	return &_go.CloseReasonGroup{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		CreatedAt:   util.Timestamp(createdAt),
		UpdatedAt:   util.Timestamp(updatedAt),
		CreatedBy:   &createdByLookup,
		UpdatedBy:   &updatedByLookup,
	}, nil
}

func (s CloseReasonGroup) List(rpc *model.SearchOptions) (*_go.CloseReasonGroupList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.list.database_connection_error", dbErr)
	}

	query, args, err := s.buildSearchCloseReasonGroupQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.list.query_build_error", err)
	}

	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.list.execution_error", err)
	}
	defer rows.Close()

	var lookupList []*_go.CloseReasonGroup
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		l := &_go.CloseReasonGroup{}

		var (
			createdBy, updatedBy         general.Lookup
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
			}
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.list.row_scan_error", err)
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

	return &_go.CloseReasonGroupList{
		Page:  int32(rpc.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

func (s CloseReasonGroup) Delete(rpc *model.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.cases.close_reason_group.delete.database_connection_error", dbErr)
	}

	query, args, err := s.buildDeleteCloseReasonGroupQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.close_reason_group.delete.query_build_error", err)
	}

	res, err := d.Exec(rpc.Context, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.cases.close_reason_group.delete.execution_error", err)
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return dberr.NewDBNoRowsError("postgres.cases.close_reason_group.delete.no_rows_affected")
	}

	return nil
}

func (s CloseReasonGroup) Update(rpc *model.UpdateOptions, l *_go.CloseReasonGroup) (*_go.CloseReasonGroup, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.update.database_connection_error", dbErr)
	}

	query, args, queryErr := s.buildUpdateCloseReasonGroupQuery(rpc, l)
	if queryErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.update.query_build_error", queryErr)
	}

	var (
		createdBy, updatedByLookup general.Lookup
		createdAt, updatedAt       time.Time
	)

	err := d.QueryRow(rpc.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt, &l.Description,
		&createdBy.Id, &createdBy.Name, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.update.execution_error", err)
	}

	l.CreatedAt = util.Timestamp(createdAt)
	l.UpdatedAt = util.Timestamp(updatedAt)
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedByLookup

	return l, nil
}

// buildCreateCloseReasonGroupQuery constructs the SQL insert query and returns the query string and arguments.
func (s CloseReasonGroup) buildCreateCloseReasonGroupQuery(rpc *model.CreateOptions, lookup *_go.CloseReasonGroup) (string, []interface{}, error) {
	query := createCloseReasonGroupQuery
	args := []interface{}{
		lookup.Name, rpc.Session.GetDomainId(), rpc.Time, lookup.Description, rpc.Session.GetUserId(),
	}
	return query, args, nil
}

func (s CloseReasonGroup) buildSearchCloseReasonGroupQuery(rpc *model.SearchOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	queryBuilder := sq.Select().
		From("cases.close_reason_group AS g").
		Where(sq.Eq{"g.dc": rpc.Session.GetDomainId()}).
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
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
		case "updated_by":
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
		substr := util.Substring(name)
		combinedLike := strings.Join(substr, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"g.name": combinedLike})
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
		return "", nil, dberr.NewDBInternalError("postgres.cases.close_reason_group.query_build.sql_generation_error", err)
	}

	return db.CompactSQL(query), args, nil
}

func (s CloseReasonGroup) buildDeleteCloseReasonGroupQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deleteCloseReasonGroupQuery
	args := []interface{}{pq.Array(ids), rpc.Session.GetDomainId()}
	return query, args, nil
}

func (s CloseReasonGroup) buildUpdateCloseReasonGroupQuery(rpc *model.UpdateOptions, l *_go.CloseReasonGroup) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	builder := psql.Update("cases.close_reason_group").
		Set("updated_at", rpc.Time).
		Set("updated_by", rpc.Session.GetUserId()).
		Where(sq.Eq{"id": l.Id}).
		Where(sq.Eq{"dc": rpc.Session.GetDomainId()})

	for _, field := range rpc.Fields {
		switch field {
		case "name":
			if l.Name != "" {
				builder = builder.Set("name", l.Name)
			}
		case "description":
			builder = builder.Set("description", sq.Expr("NULLIF(?, '')", l.Description))
		}
	}

	sql, args, err := builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	q := fmt.Sprintf(`
WITH upd AS (
	%s
	RETURNING id, name, created_at, updated_at, description, created_by, updated_by
)
SELECT upd.id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       COALESCE(upd.description, '') AS description,
       upd.created_by AS created_by_id,
       COALESCE(c.name::text, c.username, '') AS created_by_name,
       upd.updated_by AS updated_by_id,
       COALESCE(u.name::text, u.username) AS updated_by_name
FROM upd
LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
LEFT JOIN directory.wbt_user c ON c.id = upd.created_by;
`, sql)

	return db.CompactSQL(q), args, nil
}

var (
	createCloseReasonGroupQuery = db.CompactSQL(`
	WITH ins AS (
		INSERT INTO cases.close_reason_group (name, dc, created_at, description, created_by, updated_at, updated_by)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $3, $5)
		RETURNING *
	)
	SELECT ins.id,
		   ins.name,
		   ins.created_at,
		   COALESCE(ins.description, '') AS description,
		   ins.created_by AS created_by_id,
		   COALESCE(c.name::text, c.username) AS created_by_name,
		   ins.updated_at,
		   ins.updated_by AS updated_by_id,
		   COALESCE(u.name::text, u.username) AS updated_by_name
	FROM ins
	LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
	LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;
	`)

	deleteCloseReasonGroupQuery = db.CompactSQL(`
	DELETE FROM cases.close_reason_group
	WHERE id = ANY($1) AND dc = $2
`)
)

func NewCloseReasonGroupStore(store db.Store) (db.CloseReasonGroupStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.config.new_close_reason_group.check.bad_arguments",
			"error creating config interface to the close_reason_group table, main store is nil")
	}
	return &CloseReasonGroup{storage: store}, nil
}
