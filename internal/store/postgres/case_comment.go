package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
	_go "github.com/webitel/cases/api/cases"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
)

type CaseComment struct {
	storage store.Store
}

// Publish implements store.CommentCaseStore for publishing a single comment.
func (c *CaseComment) Publish(
	rpc *model.CreateOptions,
	add *_go.CaseComment,
) (*_go.CaseComment, error) {
	// Establish database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.database_connection_error", dbErr)
	}

	// Build the insert and select query with RETURNING clause
	sq, plan, err := c.buildPublishCommentsSqlizer(rpc, &_go.InputCaseComment{Text: add.Text}, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.build_sqlizer_error", err)
	}

	query, args, err := sq.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.query_to_sql_error", err)
	}

	// Temporary variables for `created_at` and `updated_at` timestamps
	var createdAt, updatedAt time.Time

	// Replace `created_at` and `updated_at` in `plan` with time.Time
	for i, field := range rpc.Fields {
		switch field {
		case "created_at":
			plan[i] = &createdAt
		case "updated_at":
			plan[i] = &updatedAt
		}
	}

	// Execute the query and scan the result directly into `add`
	if err = d.QueryRow(rpc.Context, query, args...).Scan(plan...); err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.scan_error", err)
	}

	// Convert the returned ID to integer and handle any error
	commId, err := strconv.Atoi(add.Id)
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.convert_id_error", err)
	}

	// Convert `created_at` and `updated_at` to Unix timestamps for protobuf fields
	add.CreatedAt = util.Timestamp(createdAt)
	add.UpdatedAt = util.Timestamp(updatedAt)

	// Encode etag from the comment ID and version
	e := etag.EncodeEtag(etag.EtagCaseComment, int64(commId), add.Ver)
	add.Id = e

	return add, nil
}

// buildPublishCommentsSqlizer builds a single query that inserts one comment
// and returns only the specified fields from the inserted row using a CTE.
func (c *CaseComment) buildPublishCommentsSqlizer(
	rpc *model.CreateOptions,
	input *_go.InputCaseComment,
	output *_go.CaseComment,
) (sq.Sqlizer, []any, error) {
	// Ensure "id" and "ver" are in the fields list
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)

	// Prepare the scan plan for the query
	var plan []any

	// Start building the insert part of the query using Squirrel
	insertBuilder := sq.
		Insert("cases.case_comment").
		Columns("dc", "case_id", "created_at", "created_by", "updated_at", "updated_by", "comment").
		Values(
			rpc.Session.GetDomainId(), // dc
			rpc.ParentID,              // case_id
			rpc.CurrentTime(),         // created_at (and updated_at)
			rpc.Session.GetUserId(),   // created_by (and updated_by)
			rpc.CurrentTime(),         // updated_at
			rpc.Session.GetUserId(),   // updated_by
			input.Text,                // comment text
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	// Convert insertBuilder to SQL to use it within a CTE
	insertSQL, insertArgs, err := insertBuilder.ToSql()
	if err != nil {
		return nil, nil, dberr.NewDBInternalError("store.case_comment.build_publish_comments_sqlizer.insert_query_error", err)
	}

	// Use the insert SQL as a CTE prefix for the main select query
	ctePrefix := sq.Expr("WITH cc AS ("+insertSQL+")", insertArgs...)

	// Dynamically build the SELECT query for retrieving only the specified fields
	selectBuilder := sq.Select()

	// Add only the fields specified in rpc.Fields to the SELECT clause
	for _, field := range rpc.Fields {
		switch field {
		case "id":
			selectBuilder = selectBuilder.Column("cc.id")
			plan = append(plan, &output.Id)
		case "case_id":
			selectBuilder = selectBuilder.Column("cc.case_id")
			plan = append(plan, &output.CaseId)
		case "created_at":
			selectBuilder = selectBuilder.Column("cc.created_at")
			plan = append(plan, &output.CreatedAt)
		case "comment":
			selectBuilder = selectBuilder.Column("cc.comment")
			plan = append(plan, &output.Text)
		case "ver":
			selectBuilder = selectBuilder.Column("cc.ver")
			plan = append(plan, &output.Ver)
		case "created_by":
			if output.CreatedBy == nil {
				output.CreatedBy = &_go.Lookup{}
			}
			selectBuilder = selectBuilder.
				Column("(SELECT ROW (id, name)::text FROM directory.wbt_user WHERE id = cc.created_by) AS created_by")
			plan = append(plan, scanner.ScanRowLookup(&output.CreatedBy))
		case "updated_by":
			if output.UpdatedBy == nil {
				output.UpdatedBy = &_go.Lookup{}
			}
			selectBuilder = selectBuilder.
				Column("(SELECT ROW (id, name)::text FROM directory.wbt_user WHERE id = cc.updated_by) AS updated_by")
			plan = append(plan, scanner.ScanRowLookup(&output.UpdatedBy))
		case "updated_at":
			selectBuilder = selectBuilder.Column("cc.updated_at")
			plan = append(plan, &output.UpdatedAt)
		}
	}

	// Combine the CTE with the select query
	sqBuilder := selectBuilder.
		From("cc").
		PrefixExpr(ctePrefix)

	return sqBuilder, plan, nil
}

// Delete implements store.CommentCaseStore.
func (c *CaseComment) Delete(
	rpc *model.DeleteOptions,
) error {
	// Establish database connection
	d, err := c.storage.Database()
	if err != nil {
		return dberr.NewDBInternalError("store.case_comment.delete.database_connection_error", err)
	}

	// Build the delete query
	query, args, dbErr := c.buildDeleteCaseCommentQuery(rpc)
	if dbErr != nil {
		return dberr.NewDBInternalError("store.case_comment.delete.query_build_error", dbErr)
	}

	// Execute the query
	res, execErr := d.Exec(rpc.Context, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("store.case_comment.delete.exec_error", execErr)
	}

	// Check if any rows were affected
	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("store.case_comment.delete.not_found")
	}

	return nil
}

func (c CaseComment) buildDeleteCaseCommentQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	query := deleteCaseCommentQuery
	args := []interface{}{pq.Array(ids), rpc.Session.GetDomainId()}
	return query, args, nil
}

var deleteCaseCommentQuery = store.CompactSQL(`
	DELETE FROM cases.case_comment
	WHERE id = ANY($1) AND dc = $2
`)

func (c *CaseComment) List(rpc *model.SearchOptions) (*_go.CaseCommentList, error) {
	// Connect to the database
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.list.database_connection_error", dbErr)
	}

	// Build the query and plan builder using BuildListCaseCommentsSqlizer
	queryBuilder, planBuilder, err := c.BuildListCaseCommentsSqlizer(rpc, &_go.CaseComment{})
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.list.query_build_error", err)
	}

	// Convert the query to SQL
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.list.query_to_sql_error", err)
	}

	// Execute the query
	rows, err := d.Query(rpc.Context, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.list.execution_error", err)
	}
	defer rows.Close()

	var commentList []*_go.CaseComment
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		comment := &_go.CaseComment{}

		// Temporary variables to hold timestamp values
		var createdAt, updatedAt time.Time

		// Plan with temporary time.Time variables for created_at and updated_at
		plan := planBuilder(comment)
		for i, field := range rpc.Fields {
			switch field {
			case "created_at":
				plan[i] = &createdAt
			case "updated_at":
				plan[i] = &updatedAt
			}
		}

		// Scan row into the comment fields using the plan
		if err := rows.Scan(plan...); err != nil {
			return nil, dberr.NewDBInternalError("store.case_comment.list.row_scan_error", err)
		}

		// Convert the time.Time values to Unix timestamps for protobuf fields
		comment.CreatedAt = util.Timestamp(createdAt)
		comment.UpdatedAt = util.Timestamp(updatedAt)

		// Encode the `id` and `ver` fields into an etag
		commId, err := strconv.Atoi(comment.Id)
		if err != nil {
			return nil, dberr.NewDBInternalError("store.case_comment.list.id_conversion_error", err)
		}
		comment.Id = etag.EncodeEtag(etag.EtagCaseComment, int64(commId), comment.Ver)

		commentList = append(commentList, comment)
		lCount++
	}

	return &_go.CaseCommentList{
		Page:  int64(rpc.Page),
		Next:  next,
		Items: commentList,
	}, nil
}

func (c *CaseComment) BuildListCaseCommentsSqlizer(
	rpc *model.SearchOptions,
	_ *_go.CaseComment,
) (sq.Sqlizer, func(*_go.CaseComment) []any, error) {
	// Begin building the query
	queryBuilder := sq.Select().
		From("cases.case_comment AS cc").
		Where(sq.Eq{"cc.dc": rpc.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	if rpc.Id != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": rpc.Id})
	}

	// Ensure necessary fields are included
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)
	if util.ContainsField(rpc.Fields, "edited") {
		rpc.Fields = util.EnsureFields(rpc.Fields, "updated_at", "created_at")
	}

	// Define field mappings
	fieldMappings := map[string]struct {
		addToPlan func(output *_go.CaseComment) any
		column    string
	}{
		"id": {
			column:    "cc.id",
			addToPlan: func(output *_go.CaseComment) any { return &output.Id },
		},
		"comment": {
			column:    "cc.comment",
			addToPlan: func(output *_go.CaseComment) any { return &output.Text },
		},
		"case_id": {
			column:    "cc.case_id",
			addToPlan: func(output *_go.CaseComment) any { return &output.CaseId },
		},
		"ver": {
			column:    "cc.ver",
			addToPlan: func(output *_go.CaseComment) any { return &output.Ver },
		},
		"created_at": {
			column:    "cc.created_at",
			addToPlan: func(output *_go.CaseComment) any { return &output.CreatedAt },
		},
		"updated_at": {
			column:    "cc.updated_at",
			addToPlan: func(output *_go.CaseComment) any { return &output.UpdatedAt },
		},
		"created_by": {
			column: "(SELECT ROW (id, name)::text FROM directory.wbt_user WHERE id = cc.created_by) AS created_by",
			addToPlan: func(output *_go.CaseComment) any {
				if output.CreatedBy == nil {
					output.CreatedBy = &_go.Lookup{}
				}
				return scanner.ScanRowLookup(&output.CreatedBy)
			},
		},
		"updated_by": {
			column: "(SELECT ROW (id, name)::text FROM directory.wbt_user WHERE id = cc.updated_by) AS updated_by",
			addToPlan: func(output *_go.CaseComment) any {
				if output.UpdatedBy == nil {
					output.UpdatedBy = &_go.Lookup{}
				}
				return scanner.ScanRowLookup(&output.UpdatedBy)
			},
		},
	}

	// Loop over fields to add columns and prepare plan builder
	for _, field := range rpc.Fields {
		if mapping, ok := fieldMappings[field]; ok {
			queryBuilder = queryBuilder.Column(mapping.column)
		}
	}

	// Define the plan builder function
	planBuilder := func(output *_go.CaseComment) []any {
		var plan []any
		for _, field := range rpc.Fields {
			if mapping, ok := fieldMappings[field]; ok {
				plan = append(plan, mapping.addToPlan(output))
			}
		}
		return plan
	}

	// Apply additional filters, sorting, and pagination as needed
	if len(rpc.IDs) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.id": rpc.IDs})
	}

	if caseID, ok := rpc.Filter["case_id"].(string); ok && caseID != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": caseID})
	}

	if text, ok := rpc.Filter["text"].(string); ok && len(text) > 0 {
		substr := util.Substring(text)
		combinedLike := strings.Join(substr, "%")
		queryBuilder = queryBuilder.Where(sq.ILike{"cc.text": combinedLike})
	}

	var sortFields []string
	for _, sortField := range util.FieldsFunc(rpc.Sort, util.InlineFields) {
		desc := strings.HasPrefix(sortField, "!")
		if desc {
			sortField = strings.TrimPrefix(sortField, "!")
		}

		column := "cc." + sortField
		if desc {
			column += " DESC"
		} else {
			column += " ASC"
		}
		sortFields = append(sortFields, column)
	}

	queryBuilder = queryBuilder.OrderBy(sortFields...)

	// Pagination
	if size := rpc.GetSize(); size != -1 {
		queryBuilder = queryBuilder.Limit(uint64(size + 1))
	}
	if page := rpc.Page; page > 1 {
		queryBuilder = queryBuilder.Offset(uint64((page - 1) * rpc.GetSize()))
	}

	return queryBuilder, planBuilder, nil
}

func (c *CaseComment) Update(
	rpc *model.UpdateOptions,
	upd *_go.CaseComment,
) (*_go.CaseComment, error) {
	// Get the database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.database_connection_error", dbErr)
	}

	// Begin a transaction
	tx, err := d.Begin(rpc.Context)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.transaction_error", err)
	}
	defer tx.Rollback(rpc.Context)
	txManager := store.NewTxManager(tx)

	commId, err := strconv.Atoi(upd.Id)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.id_error", err)
	}

	// Scan the current version of the comment
	ver, err := c.ScanVer(rpc.Context, int64(commId), txManager)
	if err != nil {
		return nil, err
	}

	if upd.Ver != int32(ver) {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.conflict_error", fmt.Errorf("version mismatch"))
	}

	// Build the update query
	queryBuilder, plan, err := c.BuildUpdateCaseCommentSqlizer(rpc, upd)
	if err != nil {
		return nil, err
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.query_build_error", err)
	}

	// Temporary variables for `created_at` and `updated_at` timestamps
	var createdAt, updatedAt time.Time

	// Replace `created_at` and `updated_at` in `plan` with time.Time
	for i, field := range rpc.Fields {
		switch field {
		case "created_at":
			plan[i] = &createdAt
		case "updated_at":
			plan[i] = &updatedAt
		}
	}

	// Execute the query and scan the result
	if err := txManager.QueryRow(rpc.Context, query, args...).Scan(plan...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.execution_error", err)
	}

	// Convert `created_at` and `updated_at` to Unix timestamps
	upd.CreatedAt = util.Timestamp(createdAt)
	upd.UpdatedAt = util.Timestamp(updatedAt)

	// Encode etag from the comment ID and version
	e := etag.EncodeEtag(etag.EtagCaseComment, int64(commId), upd.Ver)
	upd.Id = e

	// Commit the transaction
	if err := tx.Commit(rpc.Context); err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.commit_error", err)
	}

	return upd, nil
}

func (c *CaseComment) ScanVer(
	ctx context.Context,
	commentID int64,
	txManager *store.TxManager,
) (int64, error) {
	// Retrieve the current version (`ver`) of the comment
	var ver int64
	err := txManager.QueryRow(ctx, "SELECT ver FROM cases.case_comment WHERE id = $1", commentID).Scan(&ver)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return a specific error if no comment with the given ID is found
			return 0, dberr.NewDBNotFoundError("postgres.cases.case_comment.scan_ver.not_found", "Comment not found")
		}
		return 0, dberr.NewDBInternalError("postgres.cases.case_comment.scan_ver.query_error", err)
	}
	return ver, nil
}

func (c *CaseComment) BuildUpdateCaseCommentSqlizer(
	rpc *model.UpdateOptions,
	input *_go.CaseComment,
) (sq.Sqlizer, []any, error) {
	// Ensure "id" and "ver" are in the fields list
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)

	// Create the `UPDATE` statement inside a CTE named `cc`
	updateBuilder := sq.Update("cases.case_comment").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.CurrentTime()).
		Set("updated_by", rpc.Session.GetUserId()).
		Set("ver", sq.Expr("ver + 1")). // Increment version
		Where(sq.Eq{"id": input.Id, "dc": rpc.Session.GetDomainId()})

	// Set the comment text if provided
	if input.Text != "" {
		updateBuilder = updateBuilder.Set("comment", input.Text)
	} else {
		return nil, nil, dberr.NewDBInternalError("store.case_comment.update.text_required", nil)
	}

	// Initialize the main selectBuilder with the `WITH` clause containing the `updateBuilder`
	selectBuilder := sq.Select().PrefixExpr(sq.Expr("WITH cc AS (?)", updateBuilder.Suffix("RETURNING *"))).From("cc")

	// Prepare the scan plan
	var plan []any
	for _, field := range rpc.Fields {
		switch field {
		case "id":
			selectBuilder = selectBuilder.Column("cc.id")
			plan = append(plan, &input.Id)
		case "comment":
			selectBuilder = selectBuilder.Column("cc.comment")
			plan = append(plan, &input.Text)
		case "created_at":
			selectBuilder = selectBuilder.Column("cc.created_at")
			plan = append(plan, &input.CreatedAt)
		case "updated_at":
			selectBuilder = selectBuilder.Column("cc.updated_at")
			plan = append(plan, &input.UpdatedAt)
		case "ver":
			selectBuilder = selectBuilder.Column("cc.ver")
			plan = append(plan, &input.Ver)
		case "created_by":
			selectBuilder = selectBuilder.
				Column("(SELECT ROW (id, name)::text FROM directory.wbt_user WHERE id = cc.created_by) AS created_by")
			if input.CreatedBy == nil {
				input.CreatedBy = &_go.Lookup{}
			}
			plan = append(plan, scanner.ScanRowLookup(&input.CreatedBy))
		case "updated_by":
			selectBuilder = selectBuilder.
				Column("(SELECT ROW (id, name)::text FROM directory.wbt_user WHERE id = cc.updated_by) AS updated_by")
			if input.UpdatedBy == nil {
				input.UpdatedBy = &_go.Lookup{}
			}
			plan = append(plan, scanner.ScanRowLookup(&input.UpdatedBy))
		}
	}

	return selectBuilder, plan, nil
}

func NewCaseCommentStore(store store.Store) (store.CaseCommentStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_comment.check.bad_arguments",
			"error creating comment case interface to the case_comment table, main store is nil")
	}
	return &CaseComment{storage: store}, nil
}
