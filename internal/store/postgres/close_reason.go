package postgres

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	_go "github.com/webitel/cases/api"
	db "github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
)

type CloseReason struct {
	storage db.Store
}

func (s CloseReason) Create(ctx *model.CreateOptions, add *_go.CloseReason) (*_go.CloseReason, error) {
	d, dbErr := s.storage.Database()

	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	query, args, err := s.buildCreateCloseReasonQuery(ctx, add)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	var createdByLookup, updatedByLookup _go.Lookup
	var createdAt, updatedAt time.Time

	err = d.QueryRow(ctx.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &add.Description,
		&createdByLookup.Id, &createdByLookup.Name,
		&updatedAt, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	t := ctx.Time

	return &_go.CloseReason{
		Id:          add.Id,
		Name:        add.Name,
		Description: add.Description,
		// When we create a new lookup - CREATED/UPDATED_AT are the same
		CreatedAt: t.Unix(),
		UpdatedAt: t.Unix(),
		CreatedBy: &createdByLookup,
		UpdatedBy: &updatedByLookup,
	}, nil
}

func (s CloseReason) List(ctx *model.SearchOptions) (*_go.CloseReasonList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	cte, err := s.buildSearchCloseReasonQuery(ctx)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	// Request one more record to check if there's a next page
	query, args, err := cte.Limit(uint64(ctx.GetSize() + 1)).ToSql()
	if err != nil {
		log.Printf("Failed to generate SQL query: %v", err)
		return nil, err
	}

	rows, err := d.Query(ctx.Context, query, args...)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var lookupList []*_go.CloseReason

	lCount := 0
	next := false
	for rows.Next() {
		if lCount >= ctx.GetSize() {
			// We've retrieved more records than the page size, so there is a next page
			next = true
			break
		}

		l := &_go.CloseReason{}
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
			log.Printf("Failed to scan row: %v", err)
			return nil, err
		}

		if ctx.FieldsUtil.ContainsField(ctx.Fields, "created_by") {
			l.CreatedBy = &createdBy
		}
		if ctx.FieldsUtil.ContainsField(ctx.Fields, "updated_by") {
			l.UpdatedBy = &updatedBy
		}

		if ctx.FieldsUtil.ContainsField(ctx.Fields, "created_at") {
			l.CreatedAt = tempCreatedAt.Unix()
		}
		if ctx.FieldsUtil.ContainsField(ctx.Fields, "updated_at") {
			l.UpdatedAt = tempUpdatedAt.Unix()
		}

		lookupList = append(lookupList, l)
		lCount++
	}

	return &_go.CloseReasonList{
		Page:  int32(ctx.Page),
		Next:  next,
		Items: lookupList,
	}, nil
}

func (s CloseReason) Delete(ctx *model.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return dbErr
	}

	query, args, err := s.buildDeleteCloseReasonQuery(ctx)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return err
	}

	res, err := d.Exec(ctx.Context, query, args...)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return err
	}

	affected := res.RowsAffected()
	if affected == 0 {
		return errors.New("no rows affected for deletion")
	}

	return nil
}

func (s CloseReason) Update(ctx *model.UpdateOptions, l *_go.CloseReason) (*_go.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}
	// Build the query and args using the helper function
	query, args, queryErr := s.buildUpdateCloseReasonQuery(ctx, l)

	if queryErr != nil {
		log.Printf("Failed to build SQL query: %v", queryErr)
		return nil, queryErr
	}

	var createdBy, updatedByLookup _go.Lookup
	var createdAt, updatedAt time.Time

	err := d.QueryRow(ctx.Context, query, args...).Scan(
		&l.Id, &l.Name, &createdAt, &updatedAt, &l.Description,
		&createdBy.Id, &createdBy.Name, &updatedByLookup.Id, &updatedByLookup.Name,
	)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	// Assigning the fields to the lookup
	l.CreatedAt = createdAt.Unix()
	l.UpdatedAt = updatedAt.Unix()
	l.CreatedBy = &createdBy
	l.UpdatedBy = &updatedByLookup

	return l, nil
}

// buildCreateCloseReasonLookupQuery constructs the SQL insert query and returns the query string and arguments.
func (s CloseReason) buildCreateCloseReasonQuery(ctx *model.CreateOptions, lookup *_go.CloseReason) (string, []interface{}, error) {
	query := `
with ins as (
    INSERT INTO cases.close_reason (name, dc, created_at, description, created_by, updated_at,
updated_by)
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
	args := []interface{}{
		lookup.Name, ctx.Session.GetDomainId(), ctx.Time, lookup.Description, ctx.Session.GetUserId(),
		ctx.Time, ctx.Session.GetUserId(),
	}
	return query, args, nil
}

// buildSearchCloseReasonLookupQuery constructs the SQL search query and returns the query builder.
func (s CloseReason) buildSearchCloseReasonQuery(ctx *model.SearchOptions) (sq.SelectBuilder, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)

	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	queryBuilder := sq.Select().
		From("cases.close_reason AS g").
		Where(sq.Eq{"g.dc": ctx.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

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

	size := ctx.GetSize()
	queryBuilder = queryBuilder.OrderBy(sortFields...).Offset(uint64((ctx.Page - 1) * size))
	if size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size))
	}

	return queryBuilder, nil
}

// buildDeleteCloseReasonLookupQuery constructs the SQL delete query and returns the query string and arguments.
func (s CloseReason) buildDeleteCloseReasonQuery(ctx *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	query := `
        DELETE FROM cases.close_reason
        WHERE id = ANY($1) AND dc = $2
    `
	args := []interface{}{pq.Array(ids), ctx.Session.GetDomainId()}
	return query, args, nil
}

func (s CloseReason) buildUpdateCloseReasonQuery(ctx *model.UpdateOptions, l *_go.CloseReason) (string, []interface{}, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// Initialize the SQL builder
	builder := psql.Update("cases.close_reason").
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
		}
	}

	// Generate SQL and arguments from the builder
	sql, args, err := builder.ToSql()
	if err != nil {
		return "", nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Construct the final SQL query with joins for created_by and updated_by
	finalQuery := fmt.Sprintf(`
WITH upd AS (
	%s
	RETURNING id, name, created_at, updated_at, description, created_by, updated_by
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

	return finalQuery, args, nil
}

// // buildUpdateCloseReasonLookupQuery constructs the SQL update query and returns the query string and arguments.
// func (s CloseReason) buildUpdateCloseReasonQuery(ctx *model.UpdateOptions, l *_go.CloseReason) (string,
// 	[]interface{},
// ) {
// 	var setClauses []string
// 	var args []interface{}

// 	args = append(args, ctx.Time, ctx.Session.GetUserId())

// 	// Add the updated_at and updated_by fields to the set clauses
// 	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", len(args)-1))
// 	setClauses = append(setClauses, fmt.Sprintf("updated_by = $%d", len(args)))

// 	// Fields that could be updated
// 	updateFields := []string{"name", "description"}

// 	// Add the fields to the set clauses
// 	for _, field := range updateFields {
// 		switch field {
// 		case "name":
// 			if l.Name != "" {
// 				args = append(args, l.Name)
// 				setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)))
// 			}
// 		case "description":
// 			if l.Description != "" {
// 				args = append(args, l.Description)
// 				setClauses = append(setClauses, fmt.Sprintf("description = $%d", len(args)))
// 			}
// 		}
// 	}

// 	// Construct the SQL query with joins for created_by and updated_by
// 	query := fmt.Sprintf(`
// with upd as (
//     UPDATE cases.close_reason
//     SET %s
//     WHERE id = $%d AND dc = $%d
//     RETURNING id, name, created_at, updated_at, description, created_by, updated_by
// )
// select upd.id,
//        upd.name,
//        upd.created_at,
//        upd.updated_at,
//        upd.description,
//        upd.created_by as created_by_id,
//        coalesce(c.name::text, c.username) as created_by_name,
//        upd.updated_by as updated_by_id,
//        coalesce(u.name::text, u.username) as updated_by_name
// from upd
//   left join directory.wbt_user u on u.id = upd.updated_by
//   left join directory.wbt_user c on c.id = upd.created_by;
//     `, strings.Join(setClauses, ", "), len(args)+1, len(args)+2)

// 	args = append(args, l.Id, ctx.Session.GetDomainId())
// 	return query, args
// }

func NewCloseReasonStore(store db.Store) (db.CloseReasonStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_close_reason.check.bad_arguments",
			"error creating config interface to the close_reason table, main store is nil")
	}
	return &CloseReason{storage: store}, nil
}
