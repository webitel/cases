package lookup

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	_gen "buf.build/gen/go/webitel/general/protocolbuffers/go"
	"database/sql"
	"github.com/Masterminds/squirrel"
	db "github.com/webitel/cases/internal/db"
	"github.com/webitel/cases/model"
	"log"
	"strings"
	"time"
)

type StatusLookup struct {
	storage db.DB
}

func (s StatusLookup) Create(rpc *model.CreateOptions, add *_go.StatusLookup) (*_go.StatusLookup, error) {
	query, args, err := s.buildCreateStatusLookupQuery(rpc.Session.GetDomainId(), rpc.Session.GetUserId(), rpc.Time, add)
	d, dbErr := s.storage.Database()

	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	var createdByLookup, updatedByLookup _gen.Lookup

	err = d.QueryRowContext(rpc.Context, query, args...).Scan(
		&add.Id, &add.Name, &rpc.Time, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&rpc.Time, &updatedByLookup.Id, &updatedByLookup.Name,
	)

	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	//When we create a new lookup - CREATED/UPDATED_AT are the same
	t := rpc.Time.Unix()

	return &_go.StatusLookup{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		CreatedAt:   t,
		UpdatedAt:   t,
		CreatedBy:   &createdByLookup,
		UpdatedBy:   &updatedByLookup,
	}, nil
}

func (s StatusLookup) Search(rpc *model.SearchOptions, ids []string) ([]*_go.StatusLookup, error) {
	cte, err := s.buildSearchStatusLookupQuery(rpc, rpc.Session.GetDomainId(), ids)

	d, dbErr := s.storage.Database()
	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	// Request one more record to check if there's a next page
	query, args, err := cte.Limit(uint64(rpc.GetSize() + 1)).ToSql()
	if err != nil {
		log.Printf("Failed to generate SQL query: %v", err)
		return nil, err
	}

	rows, err := d.DB.QueryContext(rpc.Context, query, args...)

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

	var lookupList []*_go.StatusLookup

	lCount := 0
	for rows.Next() {
		if lCount >= rpc.GetSize() {
			// We've retrieved more records than the page size, so there is a next page
			break
		}

		l := &_go.StatusLookup{}
		var createdBy, updatedBy _gen.Lookup
		var tempUpdatedAt, tempCreatedAt time.Time
		var scanArgs []interface{}

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
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		// Assign the lookup fields to the lookup
		if rpc.FieldsUtil.ContainsField(rpc.Fields, "created_by") {
			l.CreatedBy = &createdBy
		}
		if rpc.FieldsUtil.ContainsField(rpc.Fields, "updated_by") {
			l.UpdatedBy = &updatedBy
		}

		if rpc.FieldsUtil.ContainsField(rpc.Fields, "created_at") {
			l.CreatedAt = tempCreatedAt.Unix()
		}
		if rpc.FieldsUtil.ContainsField(rpc.Fields, "updated_at") {
			l.UpdatedAt = tempUpdatedAt.Unix()
		}

		lookupList = append(lookupList, l)
		lCount++
	}

	return lookupList, nil
}

func (s StatusLookup) Delete(rpc *model.DeleteOptions, id string) error {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) Update(rpc *model.UpdateOptions, lookup *_go.StatusLookup) (*_go.StatusLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookup) buildCreateStatusLookupQuery(domainID int64, createdBy int64, t time.Time, lookup *_go.StatusLookup) (string,
	[]interface{}, error) {
	query := `
with ins as (
    INSERT INTO cases.status_lookup (name, dc, created_at, description, created_by, updated_at, 
updated_by) //TODO CREATE TABLE FOR CASES
    VALUES ($1, $2, $3, $4, $5, $6, $7)
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
`
	args := []interface{}{lookup.Name, domainID, t, lookup.Description, createdBy, t, createdBy}
	return query, args, nil
}

func (s StatusLookup) buildSearchStatusLookupQuery(ctx *model.SearchOptions, domainID int64, ids []string) (squirrel.SelectBuilder, error) {
	queryBuilder := squirrel.Select().
		From("cases.status_lookup AS g").
		Where(squirrel.Eq{"g.dc": domainID}).
		PlaceholderFormat(squirrel.Dollar)

	fields := ctx.FieldsUtil.FieldsFunc(ctx.Fields, ctx.FieldsUtil.InlineFields)

	ctx.Fields = append(fields, "id")

	for _, field := range ctx.Fields {
		switch field {
		case "id", "name", "description", "created_at", "updated_at":
			queryBuilder = queryBuilder.Column("g." + field)
		case "created_by":
			queryBuilder = queryBuilder.Column("created_by.id AS created_by_id, created_by.name AS created_by_name").
				LeftJoin("directory.wbt_auth AS created_by ON g.created_by = created_by.id")
		case "updated_by":
			queryBuilder = queryBuilder.Column("updated_by.id AS updated_by_id, updated_by.name AS updated_by_name").
				LeftJoin("directory.wbt_auth AS updated_by ON g.updated_by = updated_by.id")
		}
	}

	if len(ids) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"g.id": ids})
	}

	if name, ok := ctx.Filter["name"].(string); ok && len(name) > 0 {
		substr := ctx.Match.Substring(name)
		queryBuilder = queryBuilder.Where(squirrel.ILike{"g.name": substr})
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

	size := ctx.GetSize()
	queryBuilder = queryBuilder.OrderBy(sortFields...).Offset(uint64((ctx.Page - 1) * size))
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size))
	}

	return queryBuilder, nil
}

func NewStatusLookupStore(store db.DB) (db.StatusLookupStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_status_lookup.check.bad_arguments",
			"error creating config interface to the status_lookup table, main store is nil")
	}
	return &StatusLookup{storage: store}, nil
}
