package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/webitel/cases/auth"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"github.com/webitel/cases/model/options/defaults"

	"github.com/jackc/pgx"
	_go "github.com/webitel/cases/api/cases"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type CaseCommentStore struct {
	storage *Store
}

type CommentScan func(comment *_go.CaseComment) any

const (
	caseCommentLeft              = "cc"
	caseCommentAuthorAlias       = "au"
	caseCommentCreatedByAlias    = "cb"
	caseCommentUpdatedByAlias    = "ub"
	caseCommentObjClassScopeName = model.ScopeCaseComments
)

// Publish implements store.CommentCaseStore for publishing a single comment.
func (c *CaseCommentStore) Publish(
	rpc options.CreateOptions,
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
	if err = d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("store.case_comment.publish.scan_error", err)
	}

	for _, field := range rpc.GetFields() {
		if field == "role_ids" {
			roles, defErr := c.GetRolesById(rpc, add.GetId(), auth.Read)
			if defErr != nil {
				return nil, defErr
			}
			add.RoleIds = roles
			break
		}
	}

	return add, nil
}

func (c *CaseCommentStore) buildPublishCommentsSqlizer(
	rpc options.CreateOptions,
	input *_go.InputCaseComment,
) (sq.Sqlizer, []func(comment *_go.CaseComment) any, error) {
	// Ensure "id" and "ver" are in the fields list
	fields := rpc.GetFields()
	fields = util.EnsureIdAndVerField(rpc.GetFields())
	var err error

	userID := rpc.GetAuthOpts().GetUserId()
	if createdBy := input.GetUserID(); createdBy != nil && createdBy.Id != 0 {
		userID = createdBy.Id
	}

	// Build the insert query with a RETURNING clause
	insertBuilder := sq.
		Insert("cases.case_comment").
		Columns("dc", "case_id", "created_at", "created_by", "updated_at", "updated_by", "comment").
		Values(
			rpc.GetAuthOpts().GetDomainId(), // dc
			rpc.GetParentID(),               // case_id
			rpc.RequestTime(),               // created_at (and updated_at)
			userID,                          // created_by (and updated_by)
			rpc.RequestTime(),               // updated_at
			userID,                          // updated_by
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
		fields,
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
	rpc options.DeleteOptions,
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
	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("store.case_comment.delete.exec_error", execErr)
	}

	// Check if any rows were affected
	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("store.case_comment.delete.not_found")
	}

	return nil
}

func (c CaseCommentStore) buildDeleteCaseCommentQuery(rpc options.DeleteOptions) (sq.DeleteBuilder, error) {
	var err error
	convertedIds := util.Int64SliceToStringSlice(rpc.GetIDs())
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

var deleteCaseCommentQuery = util2.CompactSQL(`
	DELETE FROM cases.case_comment
	WHERE id = ANY($1) AND dc = $2
`)

func (c *CaseCommentStore) List(rpc options.SearchOptions) (*_go.CaseCommentList, error) {
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
	rows, err := d.Query(rpc, query, args...)
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
		Page:  int64(rpc.GetPage()),
		Next:  next,
		Items: commentList,
	}, nil
}

func (c *CaseCommentStore) BuildListCaseCommentsSqlizer(
	rpc options.SearchOptions,
) (sq.Sqlizer, func(*_go.CaseComment) []any, error) {
	var defErr error

	parentId, ok := rpc.GetFilter("case_id").(int64)
	if !ok || parentId == 0 {
		return nil, nil, errors.New("case id required")
	}
	// Begin building the base query
	queryBuilder := sq.Select().
		From("cases.case_comment AS cc").
		Where(sq.Eq{"cc.dc": rpc.GetAuthOpts().GetDomainId()}).
		Where(sq.Eq{"cc.case_id": parentId}).
		PlaceholderFormat(sq.Dollar)

	queryBuilder, defErr = addCaseCommentRbacCondition(rpc.GetAuthOpts(), auth.Read, queryBuilder, "cc.id")
	if defErr != nil {
		return nil, nil, defErr
	}

	// Build select columns and scan plan using buildCommentSelectColumnsAndPlan
	queryBuilder, plan, err := buildCommentSelectColumnsAndPlan(
		queryBuilder,
		caseCommentLeft,
		rpc.GetFields(),
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
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.id": rpc.GetIDs()})
	}

	if caseID, ok := rpc.GetFilter("case_id").(string); ok && caseID != "" {
		queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": caseID})
	}

	// ----------Apply search by text -----------------
	if rpc.GetSearch() != "" {
		queryBuilder = util2.AddSearchTerm(queryBuilder, util2.Ident(caseLeft, "text"))
	}

	// -------- Apply sorting by creation date ----------
	queryBuilder = queryBuilder.OrderBy("created_at ASC")

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	return queryBuilder, planBuilder, nil
}

func (c *CaseCommentStore) Update(
	rpc options.UpdateOptions,
	input *_go.CaseComment,
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
			Text   string
			Id     int64
			UserID int64
		}{
			Text:   input.Text,
			Id:     input.Id,
			UserID: input.UpdatedBy.GetId(),
		})
	if err != nil {
		return nil, err
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.query_build_error", err)
	}

	// Convert plan to scanArgs
	scanArgs := convertToScanArgs(plan, input)

	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Explicitly indicate that the user is not the creator
			return nil, dberr.NewDBNotFoundError("postgres.case_comment.update.scan_ver.not_found", "Comment not found")
		}
		return nil, dberr.NewDBInternalError("postgres.cases.case_comment.update.execution_error", err)
	}
	for _, field := range rpc.GetFields() {
		if field == "role_ids" {
			roles, defErr := c.GetRolesById(rpc, input.GetId(), auth.Read)
			if defErr != nil {
				return nil, defErr
			}
			input.RoleIds = roles
			break
		}
	}

	return input, nil
}

func (c *CaseCommentStore) GetRolesById(
	ctx context.Context,
	commentId int64,
	access auth.AccessMode,
) ([]int64, error) {

	db, err := c.storage.Database()
	if err != nil {
		return nil, err
	}
	//// Establish database connection
	//query := "(SELECT ARRAY_AGG(DISTINCT subject) rbac_r FROM cases.case_acl WHERE object = ? AND access & ? = ?)"
	query := sq.Select("ARRAY_AGG(DISTINCT subject)").From("cases.case_comment_acl").Where("object = ?", commentId).Where("access & ? = ?", uint8(access), uint8(access)).PlaceholderFormat(sq.Dollar)
	sql, args, _ := query.ToSql()
	row := db.QueryRow(ctx, sql, args...)

	var res []int64
	defErr := row.Scan(&res)
	if defErr != nil {
		return nil, defErr
	}

	return res, nil
}

func (c *CaseCommentStore) BuildUpdateCaseCommentSqlizer(
	rpc options.UpdateOptions,
	input struct {
		Text   string
		Id     int64
		UserID int64
	},
) (sq.Sqlizer, []func(comment *_go.CaseComment) any, error) {
	var defErr error
	// Ensure "id" and "ver" are in the fields list
	fields := rpc.GetFields()
	fields = util.EnsureIdAndVerField(rpc.GetFields())

	userID := rpc.GetAuthOpts().GetUserId()
	if util.ContainsField(rpc.GetMask(), "userID") {
		if updatedBy := input.UserID; updatedBy != 0 {
			userID = updatedBy
		}
	}

	// Begin the update statement for `cases.case_comment`
	updateBuilder := sq.Update("cases.case_comment").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", userID).
		Set("ver", sq.Expr("ver + 1")). // Increment version
		// input.Etag == input.ID
		Where(sq.Eq{
			"id":         rpc.GetEtags()[0].GetOid(),
			"ver":        rpc.GetEtags()[0].GetVer(),
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
		fields,
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
			base = base.Column(util2.Ident(left, "id"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Id
			})
		case "ver":
			base = base.Column(util2.Ident(left, "ver"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Ver
			})
		case "created_by":
			joinCreatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, coalesce(%[1]s.name, %[1]s.username))::text created_by", caseCommentCreatedByAlias))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanRowLookup(&comment.CreatedBy)
			})
		case "created_at":
			base = base.Column(util2.Ident(left, "created_at"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanTimestamp(&comment.CreatedAt)
			})
		case "updated_by":
			joinUpdatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, coalesce(%[1]s.name, %[1]s.username))::text updated_by", caseCommentUpdatedByAlias))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanRowLookup(&comment.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(util2.Ident(left, "updated_at"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanTimestamp(&comment.UpdatedAt)
			})
		case "text":
			base = base.Column(util2.Ident(left, "comment"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Text
			})
		case "author":
			joinAuthor()
			base = base.Column(fmt.Sprintf(`ROW(%[1]s.id, coalesce(%[1]s.common_name, %[1]s.given_name))::text author`, caseCommentAuthorAlias))
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
		case "role_ids":
			// skip
		case "case_id":
			base = base.Column(util2.Ident(left, "case_id"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanInt64(&comment.CaseId)
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

func buildCommentsSelectAsSubquery(auther auth.Auther, fields []string, caseAlias string) (sq.SelectBuilder, []func(link *_go.CaseComment) any, *dberr.DBError) {
	alias := "comments"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.
		Select().
		From("cases.case_comment " + alias).
		Where(fmt.Sprintf("%s = %s", util2.Ident(alias, "case_id"), util2.Ident(caseAlias, "id")))
	base, err := addCaseCommentRbacCondition(auther, auth.Read, base, util2.Ident(alias, "id"))
	if err != nil {
		return base, nil, dberr.NewDBError("store.case_comment.build_comments_subquery.rbac_err", err.Error())
	}
	base, plan, dbErr := buildCommentSelectColumnsAndPlan(base, alias, fields, auther)
	if dbErr != nil {
		return base, nil, dbErr
	}
	base = util2.ApplyPaging(1, defaults.DefaultSearchSize, base)
	return base, plan, nil
}

func addCaseCommentRbacCondition(auth auth.Auther, access auth.AccessMode, query sq.SelectBuilder, dependencyColumn string) (sq.SelectBuilder, error) {
	if auth != nil && auth.IsRbacCheckRequired(caseCommentObjClassScopeName, access) {
		return query.Where(sq.Expr(fmt.Sprintf("EXISTS(SELECT acl.object FROM cases.case_comment_acl acl WHERE acl.dc = ? AND acl.object = %s AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)", dependencyColumn),
			auth.GetDomainId(), pq.Array(auth.GetRoles()), int64(access), int64(access))), nil

	}
	return query, nil
}

func addCaseCommentRbacConditionForDelete(auth auth.Auther, access auth.AccessMode, query sq.DeleteBuilder, dependencyColumn string) (sq.DeleteBuilder, error) {
	if auth != nil && auth.IsRbacCheckRequired(caseCommentObjClassScopeName, access) {
		return query.Where(sq.Expr(fmt.Sprintf("EXISTS(SELECT acl.object FROM cases.case_comment_acl acl WHERE acl.dc = ? AND acl.object = %s AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)", dependencyColumn),
			auth.GetDomainId(), pq.Array(auth.GetRoles()), int64(access), int64(access))), nil

	}
	return query, nil
}

func addCaseCommentRbacConditionForUpdate(auth auth.Auther, access auth.AccessMode, query sq.UpdateBuilder, dependencyColumn string) (sq.UpdateBuilder, error) {
	if auth != nil && auth.IsRbacCheckRequired(caseCommentObjClassScopeName, access) {
		return query.Where(sq.Expr(fmt.Sprintf("EXISTS(SELECT acl.object FROM cases.case_comment_acl acl WHERE acl.dc = ? AND acl.object = %s AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)", dependencyColumn),
			auth.GetDomainId(), pq.Array(auth.GetRoles()), int64(access), int64(access))), nil

	}
	return query, nil
}

func NewCaseCommentStore(store *Store) (store.CaseCommentStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_case_comment.check.bad_arguments",
			"error creating comment case interface to the case_comment table, main store is nil")
	}
	return &CaseCommentStore{storage: store}, nil
}
