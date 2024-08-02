package postgres

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

type StatusConditionStore struct {
	storage store.Store
}

func (s StatusConditionStore) Create(ctx *model.CreateOptions, add *_go.StatusCondition) (*_go.StatusCondition, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	tx, err := d.BeginTx(ctx.Context, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
			panic(p)
		} else if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
		} else {
			err = tx.Commit()
		}
	}()

	query, args, err := s.buildCreateStatusConditionQuery(ctx, add)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return nil, err
	}

	var createdByLookup, updatedByLookup _gen.Lookup
	var createdAt, updatedAt time.Time

	err = tx.QueryRowContext(ctx.Context, query, args...).Scan(
		&add.Id, &add.Name, &createdAt, &updatedAt, &add.Description, &add.Initial, &add.Final,
		&createdByLookup.Id, &createdByLookup.Name, &updatedByLookup.Id, &updatedByLookup.Name, &add.Id,
	)

	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	add.CreatedAt = createdAt.Unix()
	add.UpdatedAt = updatedAt.Unix()
	add.CreatedBy = &createdByLookup
	add.UpdatedBy = &updatedByLookup

	return add, nil
}

func (s StatusConditionStore) List(ctx *model.SearchOptions) (*_go.StatusConditionList, error) {
	queryBuilder, err := s.buildListStatusConditionQuery(ctx)
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

	var statusList []*_go.StatusCondition
	lCount := 0
	next := false
	for rows.Next() {
		if lCount >= ctx.GetSize() {
			next = true
			break
		}

		st := &_go.StatusCondition{}
		var createdBy, updatedBy _gen.Lookup
		var tempUpdatedAt, tempCreatedAt time.Time
		var scanArgs []interface{}

		for _, field := range ctx.Fields {
			switch field {
			case "id":
				scanArgs = append(scanArgs, &st.Id)
			case "name":
				scanArgs = append(scanArgs, &st.Name)
			case "description":
				scanArgs = append(scanArgs, &st.Description)
			case "initial":
				scanArgs = append(scanArgs, &st.Initial)
			case "final":
				scanArgs = append(scanArgs, &st.Final)
			case "created_at":
				scanArgs = append(scanArgs, &tempCreatedAt)
			case "updated_at":
				scanArgs = append(scanArgs, &tempUpdatedAt)
			case "created_by":
				scanArgs = append(scanArgs, &createdBy.Id, &createdBy.Name)
			case "updated_by":
				scanArgs = append(scanArgs, &updatedBy.Id, &updatedBy.Name)
			case "status_id":
				scanArgs = append(scanArgs, &st.Id)
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

	return &_go.StatusConditionList{
		Page:  int32(ctx.Page),
		Next:  next,
		Items: statusList,
	}, nil
}

func (s StatusConditionStore) Delete(ctx *model.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return dbErr
	}

	tx, err := d.BeginTx(ctx.Context, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
			panic(p)
		} else if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
		} else {
			err = tx.Commit()
		}
	}()

	query, args, err := s.buildDeleteStatusConditionQuery(ctx)
	if err != nil {
		log.Printf("Failed to build SQL query: %v", err)
		return err
	}

	rows, err := tx.QueryContext(ctx.Context, query, args...)
	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Failed to close rows: %v", err)
		}
	}(rows)

	affected := 0
	for rows.Next() {
		affected++
	}

	if affected == 0 {
		var initialCount, finalCount int
		checkQuery := `
			SELECT 
				(SELECT COUNT(*) FROM cases.status_condition WHERE initial = TRUE AND dc = $1) AS initial_count,
				(SELECT COUNT(*) FROM cases.status_condition WHERE final = TRUE AND dc = $1) AS final_count
		`
		err = tx.QueryRowContext(ctx.Context, checkQuery, ctx.Session.GetDomainId()).Scan(&initialCount, &finalCount)
		if err != nil {
			log.Printf("Failed to execute check query: %v", err)
			return err
		}

		if initialCount == 1 {
			return errors.New("cannot delete the last initial status condition")
		}
		if finalCount == 1 {
			return errors.New("cannot delete the last final status condition")
		}

		return errors.New("no rows affected for deletion or constraints violated")
	}

	return nil
}

func (s StatusConditionStore) Update(ctx *model.UpdateOptions, st *_go.StatusCondition) (*_go.StatusCondition, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		log.Printf("Failed to get database connection: %v", dbErr)
		return nil, dbErr
	}

	tx, err := d.BeginTx(ctx.Context, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
			panic(p)
		} else if err != nil {
			err := tx.Rollback()
			if err != nil {
				log.Printf("Failed to rollback transaction: %v", err)
			}
		} else {
			err = tx.Commit()
		}
	}()

	query, args := s.buildUpdateStatusConditionQuery(ctx, st)

	var createdBy, updatedByLookup _gen.Lookup
	var createdAt, updatedAt time.Time
	var remainingInitial, remainingFinal int

	err = tx.QueryRowContext(ctx.Context, query, args...).Scan(
		&st.Id, &st.Name, &createdAt, &updatedAt, &st.Description, &st.Initial, &st.Final,
		&createdBy.Id, &createdBy.Name, &updatedByLookup.Id, &updatedByLookup.Name, &st.Id,
		&remainingInitial, &remainingFinal,
	)

	if err != nil {
		log.Printf("Failed to execute SQL query: %v", err)
		return nil, err
	}

	// Check constraints
	if remainingInitial == 1 {
		return nil, errors.New("cannot update to remove the last initial status condition")
	}
	if remainingFinal == 1 {
		return nil, errors.New("cannot update to remove the last final status condition")
	}

	st.CreatedAt = createdAt.Unix()
	st.UpdatedAt = updatedAt.Unix()
	st.CreatedBy = &createdBy
	st.UpdatedBy = &updatedByLookup

	return st, nil
}

func (s StatusConditionStore) buildCreateStatusConditionQuery(ctx *model.CreateOptions, status *_go.StatusCondition) (string, []interface{}, error) {
	query := `
WITH existing_status AS (
    SELECT COUNT(*) AS count
    FROM cases.status_condition
    WHERE dc = $9 AND status_id = $10
),
default_values AS (
    SELECT 
        CASE WHEN (SELECT count FROM existing_status) = 0 THEN TRUE ELSE FALSE END AS initial_default,
        CASE WHEN (SELECT count FROM existing_status) = 0 THEN TRUE ELSE FALSE END AS final_default
),
ins AS (
    INSERT INTO cases.status_condition (name, created_at, description, initial, final, created_by, updated_at, updated_by, dc, status_id)
    VALUES ($1, $2, $3, COALESCE($4, (SELECT initial_default FROM default_values)), COALESCE($5, (SELECT final_default FROM default_values)), $6, $7, $8, $9, $10)
    RETURNING id, name, created_at, description, initial, final, created_by, updated_at, updated_by, status_id

SELECT 
    ins.id,
    ins.name,
    ins.created_at,
    ins.description,
    ins.initial,
    ins.final,
    ins.created_by AS created_by_id,
    COALESCE(c.name::text, c.username) AS created_by_name,
    ins.updated_at,
    ins.updated_by AS updated_by_id,
    COALESCE(u.name::text, u.username) AS updated_by_name,
    ins.status_id
FROM ins
LEFT JOIN directory.wbt_user u ON u.id = ins.updated_by
LEFT JOIN directory.wbt_user c ON c.id = ins.created_by;
`
	args := []interface{}{
		status.Name, ctx.CurrentTime(), status.Description, status.Initial, status.Final,
		ctx.Session.GetUserId(), ctx.CurrentTime(), ctx.Session.GetUserId(), ctx.Session.GetDomainId(), status.Id,
	}
	return query, args, nil
}

func (s StatusConditionStore) buildListStatusConditionQuery(ctx *model.SearchOptions) (squirrel.SelectBuilder, error) {
	queryBuilder := squirrel.Select().
		From("cases.status_condition AS s").
		Where(squirrel.Eq{"s.dc": ctx.Session.GetDomainId()}).
		PlaceholderFormat(squirrel.Dollar)

	fields := ctx.FieldsUtil.FieldsFunc(ctx.Fields, ctx.FieldsUtil.InlineFields)

	ctx.Fields = append(fields, "id")

	for _, field := range ctx.Fields {
		switch field {
		case "id", "name", "description", "initial", "final", "created_at", "updated_at":
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

func (s StatusConditionStore) buildDeleteStatusConditionQuery(ctx *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := ctx.FieldsUtil.Int64SliceToStringSlice(ctx.IDs)
	ids := ctx.FieldsUtil.FieldsFunc(convertedIds, ctx.FieldsUtil.InlineFields)

	query := `
WITH
    to_delete AS (
        SELECT id, initial, final
        FROM cases.status_condition
        WHERE id = ANY($1) AND dc = $2
    ),
    initial_count AS (
        SELECT COUNT(*) AS count
        FROM cases.status_condition
        WHERE initial = TRUE AND id NOT IN (SELECT id FROM to_delete) AND dc = $2
    ),
    final_count AS (
        SELECT COUNT(*) AS count
        FROM cases.status_condition
        WHERE final = TRUE AND id NOT IN (SELECT id FROM to_delete) AND dc = $2
    ),
    checks AS (
        SELECT
            (SELECT COUNT(*) FROM to_delete WHERE initial = TRUE) AS deleting_initial,
            (SELECT COUNT(*) FROM to_delete WHERE final = TRUE) AS deleting_final,
            (SELECT count FROM initial_count) AS remaining_initial,
            (SELECT count FROM final_count) AS remaining_final
    ),
    perform_delete AS (
        DELETE FROM cases.status_condition
        WHERE id IN (SELECT id FROM to_delete)
        RETURNING id
    )
SELECT * FROM perform_delete
WHERE
    (SELECT deleting_initial FROM checks) = 0 OR (SELECT remaining_initial FROM checks) > 0
    AND
    (SELECT deleting_final FROM checks) = 0 OR (SELECT remaining_final FROM checks) > 0;
`

	args := []interface{}{pq.Array(ids), ctx.Session.GetDomainId()}
	return query, args, nil
}

func (s StatusConditionStore) buildUpdateStatusConditionQuery(ctx *model.UpdateOptions, st *_go.StatusCondition) (string, []interface{}) {
	var setClauses []string
	var args []interface{}

	args = append(args, ctx.CurrentTime(), ctx.Session.GetUserId())

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", len(args)-1))
	setClauses = append(setClauses, fmt.Sprintf("updated_by = $%d", len(args)))

	updateFields := []string{"name", "description", "initial", "final"}

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
		case "initial":
			args = append(args, st.Initial)
			setClauses = append(setClauses, fmt.Sprintf("initial = $%d", len(args)))
		case "final":
			args = append(args, st.Final)
			setClauses = append(setClauses, fmt.Sprintf("final = $%d", len(args)))
		}
	}

	query := fmt.Sprintf(`
WITH 
    reset_initial AS (
        UPDATE cases.status_condition
        SET initial = FALSE
        WHERE initial = TRUE
        AND $%d IS TRUE
        RETURNING id
    ),
    ensure_initial AS (
        SELECT COUNT(*) AS count
        FROM cases.status_condition
        WHERE initial = TRUE AND id <> $%d AND dc = $%d
    ),
    ensure_final AS (
        SELECT COUNT(*) AS count
        FROM cases.status_condition
        WHERE final = TRUE AND id <> $%d AND dc = $%d
    ),
    upd AS (
        UPDATE cases.status_condition
        SET %s
        WHERE id = $%d AND dc = $%d
        RETURNING id, name, created_at, updated_at, description, initial, final, created_by, updated_by, status_id
    )
SELECT 
    upd.id,
    upd.name,
    upd.created_at,
    upd.updated_at,
    upd.description,
    upd.initial,
    upd.final,
    upd.created_by AS created_by_id,
    COALESCE(c.name::text, c.username) AS created_by_name,
    upd.updated_by AS updated_by_id,
    COALESCE(u.name::text, u.username) AS updated_by_name,
    upd.status_id,
    (SELECT count FROM ensure_initial) AS remaining_initial,
    (SELECT count FROM ensure_final) AS remaining_final
FROM upd
LEFT JOIN directory.wbt_user u ON u.id = upd.updated_by
LEFT JOIN directory.wbt_user c ON c.id = upd.created_by;
`, len(args)+1, st.Id, ctx.Session.GetDomainId(), st.Id, ctx.Session.GetDomainId(), strings.Join(setClauses, ", "), len(args)+2, len(args)+3)

	args = append(args, st.Initial, st.Final)
	return query, args
}

func NewStatusConditionStore(store store.Store) (store.StatusConditionStore, model.AppError) {
	if store == nil {
		return nil, model.NewInternalError("postgres.config.new_status_condition.check.bad_arguments",
			"error creating config interface to the status_condition table, main store is nil")
	}
	return &StatusConditionStore{storage: store}, nil
}
