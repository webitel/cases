package postgres

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/webitel/webitel-go-kit/etag"

	"github.com/jackc/pgx"
	_go "github.com/webitel/cases/api/cases"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	dberr "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/scanner"
	"github.com/webitel/cases/model"
	util "github.com/webitel/cases/util"
)

type CaseCommentStore struct {
	storage store.Store
}

type CommentScan func(comment *_go.CaseComment) any

const (
	caseCommentLeft = "cc"
)

// Publish implements store.CommentCaseStore for publishing a single comment.
func (c *CaseCommentStore) Publish(
	rpc *model.CreateOptions,
	add *_go.CaseComment,
) (*_go.CaseComment, error) {
	// Establish database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.database_connection_error", dbErr)
	}

	// Build the insert and select query with RETURNING clause
	sq, plan, err := c.buildPublishCommentsSqlizer(rpc, &_go.InputCaseComment{Text: add.Text})
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.build_sqlizer_error", err)
	}

	query, args, err := sq.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.query_to_sql_error", err)
	}

	// Convert plan to scanArgs
	scanArgs := convertToScanArgs(plan, add)

	// Execute the query and scan the result directly into `add`
	if err = d.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.scan_error", err)
	}

	// Convert the returned ID to integer and handle any error
	commId, err := strconv.Atoi(add.Id)
	if err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.convert_id_error", err)
	}

	// Encode etag from the comment ID and version
	e := etag.EncodeEtag(etag.EtagCaseComment, int64(commId), add.Ver)
	add.Id = e

	return add, nil
}

func (c *CaseCommentStore) buildPublishCommentsSqlizer(
	rpc *model.CreateOptions,
	input *_go.InputCaseComment,
) (sq.Sqlizer, []CommentScan, error) {
	// Ensure "id" and "ver" are in the fields list
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)

	// Build the insert query with a RETURNING clause
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

	// Build select clause and scan plan dynamically using buildCommentSelectColumnsAndPlan
	selectBuilder := sq.Select()
	selectBuilder, plan, dbErr := buildCommentSelectColumnsAndPlan(selectBuilder, caseCommentLeft, rpc.Fields)
	if dbErr != nil {
		return nil, nil, dbErr
	}

	// Combine the CTE with the select query
	sqBuilder := selectBuilder.
		From(caseCommentLeft).
		PrefixExpr(ctePrefix)

	return sqBuilder, plan, nil
}

// Delete implements store.CommentCaseStore.
func (c *CaseCommentStore) Delete(
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

func (c CaseCommentStore) buildDeleteCaseCommentQuery(rpc *model.DeleteOptions) (string, []interface{}, error) {
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

func (c *CaseCommentStore) List(rpc *model.SearchOptions) (*_go.CaseCommentList, error) {
	// Connect to the database
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.list.database_connection_error", dbErr)
	}

	// Build the query and plan builder using BuildListCaseCommentsSqlizer
	queryBuilder, planBuilder, err := c.BuildListCaseCommentsSqlizer(rpc)
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

		// Create a new comment object
		comment := &_go.CaseComment{}
		// Build the scan plan using the planBuilder function
		plan := planBuilder(comment)

		// Scan row into the comment fields using the plan
		if err := rows.Scan(plan...); err != nil {
			return nil, dberr.NewDBInternalError("store.case_comment.list.row_scan_error", err)
		}

		commentList = append(commentList, comment)
		lCount++
	}

	return &_go.CaseCommentList{
		Page:  int64(rpc.Page),
		Next:  next,
		Items: commentList,
	}, nil
}

func (c *CaseCommentStore) BuildListCaseCommentsSqlizer(
	rpc *model.SearchOptions,
) (sq.Sqlizer, func(*_go.CaseComment) []any, error) {
	// Begin building the base query
	queryBuilder := sq.Select().
		From("cases.case_comment AS cc").
		Where(sq.Eq{"cc.dc": rpc.Session.GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	if rpc.ParentId != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": rpc.ParentId})
	}

	// Ensure necessary fields are included
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)
	if util.ContainsField(rpc.Fields, "edited") {
		rpc.Fields = util.EnsureFields(rpc.Fields, "updated_at", "created_at")
	}

	// Build select columns and scan plan using buildCommentSelectColumnsAndPlan
	queryBuilder, plan, err := buildCommentSelectColumnsAndPlan(queryBuilder, caseCommentLeft, rpc.Fields)
	if err != nil {
		return nil, nil, err
	}

	// Define the plan builder function
	planBuilder := func(output *_go.CaseComment) []any {
		var scanPlan []any
		for _, scanFunc := range plan {
			scanPlan = append(scanPlan, scanFunc(output))
		}
		return scanPlan
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

		column := caseCommentLeft + sortField
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

func (c *CaseCommentStore) Update(
	rpc *model.UpdateOptions,
	upd *_go.CaseComment,
) (*_go.CaseComment, error) {
	commId, err := strconv.Atoi(upd.Id)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.id_error", err)
	}

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

	// Scan the current version of the comment
	ver, err := c.ScanVer(rpc.Context, int64(commId), txManager)
	if err != nil {
		return nil, err
	}

	if upd.Ver != int32(ver) {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.conflict_error", fmt.Errorf("version mismatch"))
	}

	// Build the update query
	queryBuilder, plan, err := c.BuildUpdateCaseCommentSqlizer(rpc, &_go.InputCaseComment{Text: upd.Text, Etag: upd.Id})
	if err != nil {
		return nil, err
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.query_build_error", err)
	}

	// Convert plan to scanArgs
	scanArgs := convertToScanArgs(plan, upd)

	// Execute the query and scan the result
	if err := txManager.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.execution_error", err)
	}

	// Commit the transaction
	if err := tx.Commit(rpc.Context); err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.commit_error", err)
	}

	return upd, nil
}

func (c *CaseCommentStore) ScanVer(
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

func (c *CaseCommentStore) BuildUpdateCaseCommentSqlizer(
	rpc *model.UpdateOptions,
	input *_go.InputCaseComment,
) (sq.Sqlizer, []CommentScan, error) {
	// Ensure "id" and "ver" are in the fields list
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)

	// Begin the update statement for `cases.case_comment`
	updateBuilder := sq.Update("cases.case_comment").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.CurrentTime()).
		Set("updated_by", rpc.Session.GetUserId()).
		Set("ver", sq.Expr("ver + 1")). // Increment version
		// input.Etag == input.ID
		Where(sq.Eq{"id": input.Etag, "dc": rpc.Session.GetDomainId()})

	// Update the `comment` field if provided
	if input.Text != "" {
		updateBuilder = updateBuilder.Set("comment", input.Text)
	} else {
		return nil, nil, dberr.NewDBInternalError("store.case_comment.update.text_required", nil)
	}

	// Generate the CTE for the update operation
	selectBuilder := sq.Select().PrefixExpr(sq.Expr("WITH cc AS (?)", updateBuilder.Suffix("RETURNING *"))).From("cc")

	// Use `buildCommentSelectColumnsAndPlan` to build select columns and plan based on `rpc.Fields`
	selectBuilder, plan, err := buildCommentSelectColumnsAndPlan(selectBuilder, "cc", rpc.Fields)
	if err != nil {
		return nil, nil, err
	}

	return selectBuilder, plan, nil
}

// Helper function to convert a slice of CommentScan functions to a slice of empty interface{} suitable for scanning.
func convertToScanArgs(plan []CommentScan, comment *_go.CaseComment) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(comment))
	}
	return scanArgs
}

// Helper function to build the select columns and scan plan based on the fields requested.
func buildCommentSelectColumnsAndPlan(
	base sq.SelectBuilder,
	left string,
	fields []string,
) (sq.SelectBuilder, []CommentScan, *dberr.DBError) {
	var plan []CommentScan

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(store.Ident(left, "id"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Id
			})
		case "ver":
			base = base.Column(store.Ident(left, "ver"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Ver
			})
		case "created_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.created_by) created_by", left))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanRowLookup(&comment.CreatedBy)
			})
		case "created_at":
			base = base.Column(store.Ident(left, "created_at"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanTimestamp(&comment.CreatedAt)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by", left))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanRowLookup(&comment.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(store.Ident(left, "updated_at"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanTimestamp(&comment.UpdatedAt)
			})
		case "comment":
			base = base.Column(store.Ident(left, "comment"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Text
			})
		default:
			return base, nil, dberr.NewDBError("postgres.case_comment.build_comment_select.cycle_fields.unknown", fmt.Sprintf("%s field is unknown", field))
		}
	}

	if len(plan) == 0 {
		return base, nil, dberr.NewDBError("postgres.case_comment.build_comment_select.final_check.unknown", "no resulting columns")
	}

	return base, plan, nil
}

func NewCaseCommentStore(store store.Store) (store.CaseCommentStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_comment.check.bad_arguments",
			"error creating comment case interface to the case_comment table, main store is nil")
	}
	return &CaseCommentStore{storage: store}, nil
}
