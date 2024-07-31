package lookup

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	_gen "buf.build/gen/go/webitel/general/protocolbuffers/go"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	"log"
	"strings"
	"time"
)

type LookupStatusStore struct {
	storage store.Store
}

func (s LookupStatusStore) Attach(ctx *model.CreateOptions, add *_go.LookupStatus) (*_go.LookupStatus, error) {
	query, args, err := s.buildAttachLookupStatusQuery(ctx, add)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	d, dbErr := s.storage.Database()

	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	var createdByLookup, updatedByLookup _gen.Lookup

	t := ctx.CurrentTime()

	err = d.QueryRowContext(ctx.Context, query, args...).Scan(
		&add.Id, &add.LookupId, &add.Name, t, &add.Description, &add.IsInitial, &add.IsFinal,
		&createdByLookup.Id, &createdByLookup.Name,
		t, &updatedByLookup.Id, &updatedByLookup.Name,
	)

	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	return &_go.LookupStatus{
		Id:          add.Id,
		LookupId:    add.LookupId,
		Name:        add.Name,
		Description: add.Description,
		IsInitial:   add.IsInitial,
		IsFinal:     add.IsFinal,
		CreatedAt:   t.Unix(),
		UpdatedAt:   t.Unix(),
		CreatedBy:   &createdByLookup,
		UpdatedBy:   &updatedByLookup,
	}, nil
}

func (s LookupStatusStore) List(ctx *model.SearchOptions) (*_go.LookupStatusList, error) {
	queryBuilder, err := s.buildListLookupStatusQuery(ctx)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		log.Printf("Failed to generate SQL query: %v", err)
		return nil, err
	}

	d, dbErr := s.storage.Database()

	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	rows, err := d.QueryContext(ctx.Context, query, args...)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	var statusList []*_go.LookupStatus
	lCount := 0
	next := false
	for rows.Next() {
		if lCount >= ctx.GetSize() {
			next = true
			break
		}

		st := &_go.LookupStatus{}
		var createdBy, updatedBy _gen.Lookup
		var tempUpdatedAt, tempCreatedAt time.Time
		var scanArgs []interface{}

		for _, field := range ctx.Fields {
			switch field {
			case "id":
				scanArgs = append(scanArgs, &st.Id)
			case "lookup_id":
				scanArgs = append(scanArgs, &st.LookupId)
			case "name":
				scanArgs = append(scanArgs, &st.Name)
			case "description":
				scanArgs = append(scanArgs, &st.Description)
			case "is_initial":
				scanArgs = append(scanArgs, &st.IsInitial)
			case "is_final":
				scanArgs = append(scanArgs, &st.IsFinal)
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
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		if ctx.FieldsUtil.ContainsField(ctx.Fields, "created_by") {
			st.CreatedBy = &createdBy
		}
		if ctx.FieldsUtil.ContainsField(ctx.Fields, "updated_by") {
			st.UpdatedBy = &updatedBy
		}

		if ctx.FieldsUtil.ContainsField(ctx.Fields, "created_at") {
			st.CreatedAt = tempCreatedAt.Unix()
		}
		if ctx.FieldsUtil.ContainsField(ctx.Fields, "updated_at") {
			st.UpdatedAt = tempUpdatedAt.Unix()
		}

		statusList = append(statusList, st)
		lCount++
	}

	return &_go.LookupStatusList{
		Page:  int32(ctx.Page),
		Next:  next,
		Items: statusList,
	}, nil
}

func (s LookupStatusStore) Delete(ctx *model.DeleteOptions) error {
	query, args, err := s.buildDeleteLookupStatusQuery(ctx)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return err
	}

	d, dbErr := s.storage.Database()

	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return dbErr
	}

	res, err := d.ExecContext(ctx.Context, query, args...)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Failed to get affected rows: %v", err)
		return err
	}
	if affected == 0 {
		return errors.New("no rows affected for deletion")
	}

	return nil
}

func (s LookupStatusStore) Update(ctx *model.UpdateOptions, st *_go.LookupStatus) (*_go.LookupStatus, error) {
	query, args := s.buildUpdateLookupStatusQuery(ctx, st)

	d, dbErr := s.storage.Database()

	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	var createdBy, updatedByLookup _gen.Lookup
	var createdAt, updatedAt time.Time

	err := d.QueryRowContext(ctx.Context, query, args...).Scan(
		&st.Id, &st.LookupId, &st.Name, &createdAt, &updatedAt, &st.Description, &st.IsInitial, &st.IsFinal,
		&createdBy.Id, &createdBy.Name, &updatedByLookup.Id, &updatedByLookup.Name,
	)

	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	st.CreatedAt = createdAt.Unix()
	st.UpdatedAt = updatedAt.Unix()
	st.CreatedBy = &createdBy
	st.UpdatedBy = &updatedByLookup

	return st, nil
}

func (s LookupStatusStore) buildAttachLookupStatusQuery(ctx *model.CreateOptions, status *_go.LookupStatus) (string, []interface{}, error) {
	query := `
with ins as (
    INSERT INTO cases.lookup_status (lookup_id, name, created_at, description, is_initial, is_final, created_by,updated_at, updated_by, dc)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    returning *
)
select ins.id,
    ins.lookup_id,
    ins.name,
    ins.created_at,
    ins.description,
    ins.is_initial,
    ins.is_final,
    ins.created_by created_by_id,
    coalesce(c.name::text, c.username) created_by_name,
    ins.updated_at,
    ins.updated_by updated_by_id,
    coalesce(u.name::text, u.username) updated_by_name
from ins
  left join directory.wbt_user u on u.id = ins.updated_by
  left join directory.wbt_user c on c.id = ins.created_by;
`
	args := []interface{}{status.LookupId, status.Name, ctx.CurrentTime(), status.Description, status.IsInitial, status.IsFinal, ctx.Session.GetUserId(), ctx.CurrentTime(), ctx.Session.GetUserId(), ctx.Session.GetDomainId()}
	return query, args, nil
}

func (s LookupStatusStore) buildListLookupStatusQuery(ctx *model.SearchOptions) (squirrel.SelectBuilder, error) {
	queryBuilder := squirrel.Select().
		From("cases.lookup_status AS s").
		Where(squirrel.Eq{"s.dc": ctx.Session.GetDomainId()}).
		PlaceholderFormat(squirrel.Dollar)

	fields := ctx.FieldsUtil.FieldsFunc(ctx.Fields, ctx.FieldsUtil.InlineFields)

	ctx.Fields = append(fields, "id")

	for _, field := range ctx.Fields {
		switch field {
		case "id", "lookup_id", "name", "description", "is_initial", "is_final", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("s." + field)
		case "created_by":
			queryBuilder = queryBuilder.Column("created_by.id AS created_by_id, created_by.name AS created_by_name").
				LeftJoin("directory.wbt_auth AS created_by ON s.created_by = created_by.id")
		case "updated_by":
			queryBuilder = queryBuilder.Column("updated_by.id AS updated_by_id, updated_by.name AS updated_by_name").
				LeftJoin("directory.wbt_auth AS updated_by ON s.updated_by = updated_by.id")
		}
	}

	if len(ctx.IDs) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"s.id": ctx.IDs})
	}

	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substr := ctx.Match.Substring(name)
		queryBuilder = queryBuilder.Where(squirrel.ILike{"s.name": substr})
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

	size := ctx.GetSize()
	queryBuilder = queryBuilder.OrderBy(sortFields...).Offset(uint64((ctx.Page - 1) * size))
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size))
	}

	return queryBuilder, nil
}

func (s LookupStatusStore) buildDeleteLookupStatusQuery(ctx *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	query := fmt.Sprintf(`
        DELETE FROM cases.lookup_status
        WHERE id = ANY($1) AND dc = $2
    `)
	args := []interface{}{pq.Array(ids), ctx.Session.GetDomainId()}
	return query, args, nil
}

func (s LookupStatusStore) buildUpdateLookupStatusQuery(ctx *model.UpdateOptions, st *_go.LookupStatus) (string, []interface{}) {
	var setClauses []string
	var args []interface{}

	args = append(args, ctx.CurrentTime(), ctx.Session.GetUserId())

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", len(args)-1))
	setClauses = append(setClauses, fmt.Sprintf("updated_by = $%d", len(args)))

	updateFields := []string{"name", "description", "is_initial", "is_final"}

	for _, field := range updateFields {
		switch field {
		case "name":
			if st.Name != "" {
				args = append(args, st.Name)
				setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)))
			}
		case "description":
			if st.Description != "" {
				args = append(args, st.Description)
				setClauses = append(setClauses, fmt.Sprintf("description = $%d", len(args)))
			}
		case "is_initial":
			args = append(args, st.IsInitial)
			setClauses = append(setClauses, fmt.Sprintf("is_initial = $%d", len(args)))
		case "is_final":
			args = append(args, st.IsFinal)
			setClauses = append(setClauses, fmt.Sprintf("is_final = $%d", len(args)))
		}
	}

	query := fmt.Sprintf(`
with upd as (
    UPDATE cases.lookup_status
    SET %s
    WHERE id = $%d AND dc = $%d
    RETURNING id, lookup_id, name, created_at, updated_at, description, is_initial, is_final, created_by, updated_by
)
select upd.id,
       upd.lookup_id,
       upd.name,
       upd.created_at,
       upd.updated_at,
       upd.description,
       upd.is_initial,
       upd.is_final,
       upd.created_by as created_by_id,
       coalesce(c.name::text, c.username) as created_by_name,
       upd.updated_by as updated_by_id,
       coalesce(u.name::text, u.username) as updated_by_name
from upd
  left join directory.wbt_user u on u.id = upd.updated_by
  left join directory.wbt_user c on c.id = upd.created_by;
    `, strings.Join(setClauses, ", "), len(args)+1, len(args)+2)

	args = append(args, ctx.ID, ctx.Session.GetDomainId())
	return query, args
}

func NewLookupStatusStore(store store.Store) (store.LookupStatusStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_lookup_status.check.bad_arguments",
			"error creating config interface to the status_lookup table, main store is nil")
	}
	return &LookupStatusStore{storage: store}, nil
}
