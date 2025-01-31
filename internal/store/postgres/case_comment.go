package postgres

import (
	"errors"
	"fmt"
	"github.com/webitel/cases/auth"

	"github.com/jackc/pgx"
	_go "github.com/webitel/cases/api/cases"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	dberr "github.com/webitel/cases/internal/errors"
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
	caseCommentLeft              = "cc"
	caseCommentAuthorAlias       = "au"
	caseCommentCreatedByAlias    = "cb"
	caseCommentUpdatedByAlias    = "cb"
	caseCommentObjClassScopeName = "case_comments"
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

	return add, nil
}

func (c *CaseCommentStore) buildPublishCommentsSqlizer(
	rpc *model.CreateOptions,
	input *_go.InputCaseComment,
) (sq.Sqlizer, []func(comment *_go.CaseComment) any, error) {
	// Ensure "id" and "ver" are in the fields list
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)
	var (
		err error
	)
	// Build the insert query with a RETURNING clause
	insertBuilder := sq.
		Insert("cases.case_comment").
		Columns("dc", "case_id", "created_at", "created_by", "updated_at", "updated_by", "comment").
		Values(
			rpc.GetAuthOpts().GetDomainId(), // dc
			rpc.ParentID,                    // case_id
			rpc.CurrentTime(),               // created_at (and updated_at)
			rpc.GetAuthOpts().GetUserId(),   // created_by (and updated_by)
			rpc.CurrentTime(),               // updated_at
			rpc.GetAuthOpts().GetUserId(),   // updated_by
			input.Text,                      // comment text
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	// Convert insertBuilder to SQL to use it within a CTE
	insertSQL, insertArgs, err := insertBuilder.ToSql()
	if err != nil {
		return nil, nil, dberr.NewDBError("store.case_comment.build_publish_comments_sqlizer.insert_query_error", err.Error())
	}

	// Use the insert SQL as a CTE prefix for the main select query
	ctePrefix := sq.Expr("WITH cc AS ("+insertSQL+")", insertArgs...)

	// Build select clause and scan plan dynamically using buildCommentSelectColumnsAndPlan
	selectBuilder := sq.Select()
	selectBuilder, plan, dbErr := buildCommentSelectColumnsAndPlan(
		selectBuilder,
		caseCommentLeft,
		rpc.Fields,
		rpc.GetAuthOpts(),
	)
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
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("store.case_comment.delete.database_connection_error", dbErr)
	}

	// Build the delete query
	base, err := c.buildDeleteCaseCommentQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("store.case_comment.delete.query_build_error", dbErr)
	}
	// Execute the query
	query, args, err := base.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("store.case_comment.delete.to_sql.err", err)
	}
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

func (c CaseCommentStore) buildDeleteCaseCommentQuery(rpc *model.DeleteOptions) (sq.DeleteBuilder, error) {
	var err error
	convertedIds := util.Int64SliceToStringSlice(rpc.IDs)
	ids := util.FieldsFunc(convertedIds, util.InlineFields)
	base := sq.
		Delete("cases.case_comment c").
		Where("id = ANY(?)", pq.Array(ids)).
		Where("dc = ?", rpc.GetAuthOpts().GetDomainId()).
		PlaceholderFormat(sq.Dollar)
	base, err = addCaseCommentRbacConditionForDelete(rpc.GetAuthOpts(), auth.Delete, base, "c.id")
	if err != nil {
		return base, err
	}
	return base, nil
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
	var defErr error
	// Begin building the base query
	queryBuilder := sq.Select().
		From("cases.case_comment AS cc").
		Where(sq.Eq{"cc.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	if rpc.ParentId != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": rpc.ParentId})
	}
	queryBuilder, defErr = addCaseCommentRbacCondition(rpc.GetAuthOpts(), auth.Read, queryBuilder, "cc.id")
	if defErr != nil {
		return nil, nil, defErr
	}

	// Ensure necessary fields are included
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)
	if util.ContainsField(rpc.Fields, "edited") {
		rpc.Fields = util.EnsureFields(rpc.Fields, "updated_at", "created_at")
	}

	// Build select columns and scan plan using buildCommentSelectColumnsAndPlan
	queryBuilder, plan, err := buildCommentSelectColumnsAndPlan(
		queryBuilder,
		caseCommentLeft,
		rpc.Fields,
		rpc.GetAuthOpts(),
	)
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

	// ----------Apply search by text -----------------
	if rpc.Search != "" {
		queryBuilder = store.AddSearchTerm(queryBuilder, store.Ident(caseLeft, "text"))
	}

	// -------- Apply sorting by creation date ----------
	queryBuilder = queryBuilder.OrderBy("created_at ASC")

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = store.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	return queryBuilder, planBuilder, nil
}

func (c *CaseCommentStore) Update(
	rpc *model.UpdateOptions,
	upd *_go.CaseComment,
) (*_go.CaseComment, error) {
	// Get the database connection
	d, dbErr := c.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.database_connection_error", dbErr)
	}

	// Build the update query
	queryBuilder, plan, err := c.BuildUpdateCaseCommentSqlizer(
		rpc,
		struct {
			Text string
			Id   int64
		}{
			Text: upd.Text,
			Id:   upd.Id,
		})
	if err != nil {
		return nil, err
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.query_build_error", err)
	}

	// Convert plan to scanArgs
	scanArgs := convertToScanArgs(plan, upd)

	if err := d.QueryRow(rpc.Context, query, args...).Scan(scanArgs...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Explicitly indicate that the user is not the creator
			return nil, dberr.NewDBNotFoundError("postgres.case_comment.update.scan_ver.not_found", "Comment not found")
		}
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.execution_error", err)
	}

	return upd, nil
}

func (c *CaseCommentStore) BuildUpdateCaseCommentSqlizer(
	rpc *model.UpdateOptions,
	input struct {
		Text string
		Id   int64
	},
) (sq.Sqlizer, []func(comment *_go.CaseComment) any, error) {
	var defErr error
	// Ensure "id" and "ver" are in the fields list
	rpc.Fields = util.EnsureIdAndVerField(rpc.Fields)

	// Begin the update statement for `cases.case_comment`
	updateBuilder := sq.Update("cases.case_comment").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.CurrentTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Set("ver", sq.Expr("ver + 1")). // Increment version
		// input.Etag == input.ID
		Where(sq.Eq{
			"id":         rpc.Etags[0].GetOid(),
			"ver":        rpc.Etags[0].GetVer(),
			"dc":         rpc.GetAuthOpts().GetDomainId(),
			"created_by": rpc.GetAuthOpts().GetUserId(), // Ensure only the creator can edit
		})
	updateBuilder, defErr = addCaseCommentRbacConditionForUpdate(rpc.GetAuthOpts(), auth.Edit, updateBuilder, "case_comment.id")
	if defErr != nil {
		return nil, nil, defErr
	}
	// Update the `comment` field if provided
	if input.Text != "" {
		updateBuilder = updateBuilder.Set("comment", input.Text)
	} else {
		return nil, nil, dberr.NewDBInternalError("store.case_comment.update.text_required", nil)
	}

	// Generate the CTE for the update operation
	selectBuilder := sq.Select().PrefixExpr(sq.Expr("WITH cc AS (?)", updateBuilder.Suffix("RETURNING *"))).From("cc")

	// Use `buildCommentSelectColumnsAndPlan` to build select columns and plan based on `rpc.Fields`
	selectBuilder, plan, err := buildCommentSelectColumnsAndPlan(
		selectBuilder,
		caseCommentLeft,
		rpc.Fields,
		rpc.GetAuthOpts(),
	)
	if err != nil {
		return nil, nil, err
	}

	return selectBuilder, plan, nil
}

// Helper function to convert a slice of CommentScan functions to a slice of empty interface{} suitable for scanning.
func convertToScanArgs(plan []func(comment *_go.CaseComment) any, comment *_go.CaseComment) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(comment))
	}
	return scanArgs
}

// Helper function to build the select columns and scan plan based on the fields requested.
// UserAuthSession required to get some columns
func buildCommentSelectColumnsAndPlan(
	base sq.SelectBuilder,
	left string,
	fields []string,
	session auth.Auther,
) (sq.SelectBuilder, []func(comment *_go.CaseComment) any, *dberr.DBError) {
	var (
		plan           []func(link *_go.CaseComment) any
		createdByAlias string
		joinCreatedBy  = func() {
			if createdByAlias != "" {
				return
			}
			createdByAlias = caseLinkCreatedByAlias
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.created_by", caseCommentCreatedByAlias, left))
		}
		updatedByAlias string
		joinUpdatedBy  = func() {
			if updatedByAlias != "" {
				return
			}
			updatedByAlias = caseLinkUpdatedByAlias
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.updated_by", caseCommentUpdatedByAlias, left))
		}
		authorAlias string
		joinAuthor  = func() {
			if authorAlias != "" {
				return
			}
			joinCreatedBy()
			authorAlias = caseLinkAuthorAlias
			base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %[1]s.id = %s.contact_id", authorAlias, createdByAlias))
		}
	)

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
			joinCreatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, %[1]s.name)::text created_by", caseCommentCreatedByAlias))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanRowLookup(&comment.CreatedBy)
			})
		case "created_at":
			base = base.Column(store.Ident(left, "created_at"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanTimestamp(&comment.CreatedAt)
			})
		case "updated_by":
			joinUpdatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, %[1]s.name)::text updated_by", caseCommentUpdatedByAlias))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanRowLookup(&comment.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(store.Ident(left, "updated_at"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanTimestamp(&comment.UpdatedAt)
			})
		case "text":
			base = base.Column(store.Ident(left, "comment"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Text
			})
		case "author":
			joinAuthor()
			base = base.Column(fmt.Sprintf(`ROW(%[1]s.id, %[1]s.common_name)::text author`, caseCommentAuthorAlias))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanRowLookup(&comment.Author)
			})
		case "edited":
			base = base.Column(fmt.Sprintf(`(%s.created_at < %[1]s.updated_at) edited`, left))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Edited
			})
		case "can_edit":
			if session != nil {
				base = base.Column(fmt.Sprintf(`(%s.created_by = %d) can_edit`, left, session.GetUserId()))
				plan = append(plan, func(comment *_go.CaseComment) any {
					return &comment.CanEdit
				})
			}

		default:
			return base, nil, dberr.NewDBError("postgres.case_comment.build_comment_select.cycle_fields.unknown", fmt.Sprintf("%s field is unknown", field))
		}
	}

	if len(plan) == 0 {
		return base, nil, dberr.NewDBError("postgres.case_comment.build_comment_select.final_check.unknown", "no resulting columns")
	}

	return base, plan, nil
}

func buildCommentsSelectAsSubquery(opts *model.SearchOptions, caseAlias string) (sq.SelectBuilder, []func(link *_go.CaseComment) any, int, *dberr.DBError) {
	alias := "comments"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.
		Select().
		From("cases.case_comment " + alias).
		Where(fmt.Sprintf("%s = %s", store.Ident(alias, "case_id"), store.Ident(caseAlias, "id")))
	base, err := addCaseCommentRbacCondition(opts.Auth, auth.Read, base, store.Ident(alias, "id"))
	if err != nil {
		return base, nil, 0, dberr.NewDBError("store.case_comment.build_comments_subquery.rbac_err", err.Error())
	}
	base, plan, dbErr := buildCommentSelectColumnsAndPlan(base, alias, opts.Fields, opts.GetAuthOpts())
	if dbErr != nil {
		return base, nil, 0, dbErr
	}
	base, applied, dbErr := applyCaseCommentFilters(opts, base, alias)
	if dbErr != nil {
		return base, nil, 0, dbErr
	}
	base = store.ApplyPaging(opts.GetPage(), opts.GetSize(), base)
	return base, plan, applied, nil
}

func applyCaseCommentFilters(
	opts *model.SearchOptions,
	base sq.SelectBuilder,
	alias string,
) (updatedBase sq.SelectBuilder, filtersApplied int, err *dberr.DBError) {
	if opts == nil || len(opts.Filter) == 0 {
		return base, 0, nil
	}

	for column, value := range opts.Filter {
		if !util.ContainsStringIgnoreCase(opts.Fields, column) {
			continue
		}
		switch column {
		case "created_by":
			switch v := value.(type) {
			case int64, int, int32, *int64, *int, *int32:
				base = base.Where(fmt.Sprintf("%s = ?", store.Ident(caseCommentCreatedByAlias, "id")), v)
			case string, *string:
				// apply search
				// base = store.AddSearchTerm(base, )
			}
		case "author":
			switch v := value.(type) {
			case int64, int, int32, *int64, *int, *int32:
				//
				base = base.Where(fmt.Sprintf("%s = ?", store.Ident(caseCommentAuthorAlias, "id")), v)
			case string, *string:
				// apply search
				// base = store.AddSearchTerm(base, )
			}
		case "updated_by":
			switch v := value.(type) {
			case int64, int, int32, *int64, *int, *int32:
				//
				base = base.Where(fmt.Sprintf("%s = ?", store.Ident(caseCommentUpdatedByAlias, "id")), v)
			case string, *string:
				// apply search
				// base = store.AddSearchTerm(base, )
			}
			filtersApplied++
		}

	}
	return
}

func addCaseCommentRbacCondition(auth auth.Auther, access auth.AccessMode, query sq.SelectBuilder, dependencyColumn string) (sq.SelectBuilder, error) {
	if auth != nil && auth.GetObjectScope(caseCommentObjClassScopeName).IsRbacUsed() {
		subquery := sq.Select("acl.object").From("cases.case_comment_acl acl").
			Where("acl.dc = ?", auth.GetDomainId()).
			Where(fmt.Sprintf("acl.object = %s", dependencyColumn)).
			Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
			Where("acl.access & ? = ?", int64(access), int64(access)).
			Limit(1)
		return query.Where("exists(?)", subquery), nil

	}
	return query, nil
}

func addCaseCommentRbacConditionForDelete(auth auth.Auther, access auth.AccessMode, query sq.DeleteBuilder, dependencyColumn string) (sq.DeleteBuilder, error) {
	if auth != nil && auth.GetObjectScope(caseCommentObjClassScopeName).IsRbacUsed() {
		subquery := sq.Select("acl.object").From("cases.case_comment_acl acl").
			Where("acl.dc = ?", auth.GetDomainId()).
			Where(fmt.Sprintf("acl.object = %s", dependencyColumn)).
			Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
			Where("acl.access & ? = ?", int64(access), int64(access)).
			Limit(1)
		return query.Where("exists(?)", subquery), nil

	}
	return query, nil
}

func addCaseCommentRbacConditionForUpdate(auth auth.Auther, access auth.AccessMode, query sq.UpdateBuilder, dependencyColumn string) (sq.UpdateBuilder, error) {
	if auth != nil && auth.GetObjectScope(caseCommentObjClassScopeName).IsRbacUsed() {
		subquery := sq.Select("acl.object").From("cases.case_comment_acl acl").
			Where("acl.dc = ?", auth.GetDomainId()).
			Where(fmt.Sprintf("acl.object = %s", dependencyColumn)).
			Where("acl.subject = any( ?::int[])", pq.Array(auth.GetRoles())).
			Where("acl.access & ? = ?", int64(access), int64(access)).
			Limit(1)
		return query.Where("exists(?)", subquery), nil

	}
	return query, nil
}

func NewCaseCommentStore(store store.Store) (store.CaseCommentStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_comment.check.bad_arguments",
			"error creating comment case interface to the case_comment table, main store is nil")
	}
	return &CaseCommentStore{storage: store}, nil
}
