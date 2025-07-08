package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/model/options/defaults"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	storeUtil "github.com/webitel/cases/internal/store/util"

	_go "github.com/webitel/cases/api/cases"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"
	"github.com/webitel/cases/internal/model"
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

// Publish implements store.CaseCommentStore
func (c *CaseCommentStore) Publish(rpc options.Creator, input *model.CaseComment) (*model.CaseComment, error) {
	d, err := c.storage.Database()
	if err != nil {
		return nil, err
	}

	// Build the insert and select query with RETURNING clause
	selectBuilder, err := c.buildPublishCaseCommentQuery(rpc, input)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}

	var result model.CaseComment
	err = pgxscan.Get(rpc, d, &result, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	if util.ContainsField(rpc.GetFields(), "role_ids") {
		roles, err := c.GetRolesById(rpc, result.Id, auth.Read)
		if err != nil {
			return nil, err
		}
		result.RoleIds = roles
	}

	return &result, nil
}

func (c *CaseCommentStore) buildPublishCaseCommentQuery(
	rpc options.Creator,
	input *model.CaseComment,
) (sq.SelectBuilder, error) {
	// Ensure "id" and "ver" are in the fields list
	fields := util.EnsureIdAndVerField(rpc.GetFields())
	userID := rpc.GetAuthOpts().GetUserId()

	if input.Author != nil && input.Author.GetId() != nil && *input.Author.GetId() != 0 {
		userID = int64(*input.Author.GetId())
	}

	// Build the insert query with a RETURNING clause
	insertBuilder := sq.Insert("cases.case_comment").
		Columns("dc", "case_id", "created_at", "created_by", "updated_at", "updated_by", "comment").
		Values(
			rpc.GetAuthOpts().GetDomainId(), //dc
			input.CaseId,                    //case_id
			rpc.RequestTime(),               //created_at (and updated_at)
			userID,                          //created_by (and updated)by)
			rpc.RequestTime(),               //updated_at
			userID,                          //updated_by
			input.Text,                      //comment text
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

		// Convert insertBuilder to SQL to use it within a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, ParseError(err)
	}

	// Use the insert SQL as a CTE prefix for the main select query
	cte := sq.Expr("WITH cc AS ("+insertSQL+")", args...)

	selectBuilder := sq.Select()

	// Add columns to the select builder
	selectBuilder, err = buildCaseCommentSelectColumns(
		selectBuilder,
		fields,
		rpc.GetAuthOpts(),
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the select query
	selectBuilder = selectBuilder.
		From(caseCommentLeft).
		PrefixExpr(cte)

	return selectBuilder, nil
}

// Delete implements store.CaseCommentStore
func (c *CaseCommentStore) Delete(rpc options.Deleter) (*model.CaseComment, error) {
	// Establish database connection
	d, err := c.storage.Database()
	if err != nil {
		return nil, errors.NewDBInternalError("store.case_comment.delete.database_connection_error", err)
	}

	// Build the delete query
	selectBuilder, err := c.buildDeleteCaseCommentQuery(rpc)
	if err != nil {
		return nil, err
	}

	// Execute the query
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var result model.CaseComment
	err = pgxscan.Get(rpc, d, &result, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return &result, nil
}

func (c *CaseCommentStore) buildDeleteCaseCommentQuery(rpc options.Deleter) (sq.SelectBuilder, error) {

	if len(rpc.GetIDs()) == 0 {
		return sq.SelectBuilder{}, errors.InvalidArgument("no IDs provided for deletion")
	}

	convertedIds := util.Int64SliceToStringSlice(rpc.GetIDs())
	ids := util.FieldsFunc(convertedIds, util.InlineFields)

	deleteBuilder := sq.Delete("cases.case_comment").
		Where("id = ANY(?)", pq.Array(ids)).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	deleteBuilder, err := addCaseCommentRbacConditionForDelete(rpc.GetAuthOpts(), auth.Delete, deleteBuilder, "id")
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	deleteSQL, args, err := deleteBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, errors.Internal("case_comment.delete.query_to_sql_error", errors.WithCause(err))
	}

	cte := sq.Expr("WITH deleted AS ("+deleteSQL+")", args...)

	// First create a select builder
	selectBuilder := sq.Select()

	// Add columns to the select builder
	selectBuilder, err = buildCaseCommentSelectColumns(
		selectBuilder,
		rpc.GetFields(),
		rpc.GetAuthOpts(),
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the select query
	selectBuilder = selectBuilder.
		From("deleted cc").
		PrefixExpr(cte)

	return selectBuilder, nil
}

// List implements store.CaseCommentStore
func (c *CaseCommentStore) List(rpc options.Searcher) ([]*model.CaseComment, error) {
	d, err := c.storage.Database()
	if err != nil {
		return nil, err
	}

	// Build the query and plan builder using BuildListCaseCommentQuery
	selectBuilder, err := c.buildListCaseCommentQuery(rpc)
	if err != nil {
		return nil, err
	}

	// Convert the query to SQL
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	query = storeUtil.CompactSQL(query)

	var comments []*model.CaseComment
	err = pgxscan.Select(rpc, d, &comments, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return comments, nil
}

func (c *CaseCommentStore) buildListCaseCommentQuery(rpc options.Searcher) (sq.SelectBuilder, error) {

	// Begin building the base query
	queryBuilder := sq.Select().
		From("cases.case_comment AS cc").
		Where(sq.Eq{"cc.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	filters := rpc.GetFilter("case_id")
	if len(filters) > 0 {
		if parentId, err := strconv.ParseInt(filters[0].Value, 10, 64); err == nil && parentId != 0 {
			queryBuilder = queryBuilder.Where(sq.Eq{"cc.case_id": parentId})
		}
	}

	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where("cc.id = ANY(?)", rpc.GetIDs())
	}

	queryBuilder, err := addCaseCommentRbacCondition(rpc.GetAuthOpts(), auth.Read, queryBuilder, "cc.id")
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Build select columns and scan plan using buildCommentSelectColumns
	queryBuilder, err = buildCaseCommentSelectColumns(queryBuilder, rpc.GetFields(), rpc.GetAuthOpts())
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// ----------Apply search by text -----------------
	if rpc.GetSearch() != "" {
		// Use "text" which is the alias for the "comment" column
		queryBuilder = storeUtil.AddSearchTerm(queryBuilder, rpc.GetSearch(), "text")
	}

	// -------- Apply sorting by creation date ----------
	queryBuilder = queryBuilder.OrderBy("cc.created_at ASC")
	queryBuilder = storeUtil.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	return queryBuilder, nil
}

func (c *CaseCommentStore) Update(rpc options.Updator, input *model.CaseComment) (*model.CaseComment, error) {
	d, err := c.storage.Database()
	if err != nil {
		return nil, err
	}

	selectBuilder, err := c.buildUpdateCaseCommentQuery(rpc, input)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}

	var result model.CaseComment
	err = pgxscan.Get(rpc, d, &result, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	if util.ContainsField(rpc.GetFields(), "role_ids") {
		roles, err := c.GetRolesById(rpc, result.Id, auth.Read)
		if err != nil {
			return nil, err
		}
		result.RoleIds = roles
	}
	return &result, nil
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

func (c *CaseCommentStore) buildUpdateCaseCommentQuery(
	rpc options.Updator,
	input *model.CaseComment,
) (sq.SelectBuilder, error) {
	// Ensure "id" and "ver" are in the fields list
	fields := util.EnsureIdAndVerField(rpc.GetFields())
	userID := rpc.GetAuthOpts().GetUserId()

	// Check if Editor is provided and has a valid ID
	if input.Editor != nil && input.Editor.GetId() != nil && *input.Editor.GetId() != 0 {
		userID = int64(*input.Editor.GetId())
	}

	// Begin the update statement for `cases.case_comment`
	updateBuilder := sq.Update("cases.case_comment").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", userID).
		Set("ver", sq.Expr("ver + 1")). // Increment version
		// input.Etag == input.ID
		Where(sq.Eq{
			"id":         input.Id,
			"ver":        input.Ver,
			"dc":         rpc.GetAuthOpts().GetDomainId(),
			"created_by": rpc.GetAuthOpts().GetUserId(), // Ensure only the creator can edit
		})

	updateBuilder, err := addCaseCommentRbacConditionForUpdate(rpc.GetAuthOpts(), auth.Edit, updateBuilder, "case_comment.id")
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Update the `comment` field if provided
	if input.Text != "" {
		updateBuilder = updateBuilder.Set("comment", input.Text)
	}

	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, ParseError(err)
	}

	// Generate the CTE for the update operation
	cte := sq.Expr("WITH cc AS ("+updateSQL+")", args...)

	// First create a select builder
	selectBuilder := sq.Select()

	// Add columns to the select builder
	selectBuilder, err = buildCaseCommentSelectColumns(
		selectBuilder,
		fields,
		rpc.GetAuthOpts(),
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the select query
	selectBuilder = selectBuilder.
		From(caseCommentLeft).
		PrefixExpr(cte)

	return selectBuilder, nil
}

// Helper function to build the select columns and scan plan based on the fields requested.
// UserAuthSession required to get some columns
func buildCaseCommentSelectColumns(
	base sq.SelectBuilder,
	fields []string,
	session auth.Auther,
) (sq.SelectBuilder, error) {
	var (
		createdByAlias string
		joinCreatedBy  = func(alias string) string {
			if createdByAlias != "" {
				return createdByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.created_by = %s.id", alias, caseCommentLeft, alias))
			createdByAlias = alias
			return alias
		}
		updatedByAlias string
		joinUpdatedBy  = func(alias string) string {
			if updatedByAlias != "" {
				return updatedByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.updated_by = %s.id", alias, caseCommentLeft, alias))
			updatedByAlias = alias
			return alias
		}
		authorAlias string
		joinAuthor  = func(alias string) string {
			if authorAlias != "" {
				return authorAlias
			}
			createdByAlias = joinCreatedBy("ccb")
			authorAlias = alias
			base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %[1]s.id = %s.contact_id", alias, createdByAlias))
			return alias
		}
	)

	base = base.Column(storeUtil.Ident(caseCommentLeft, "id"))
	base = base.Column(storeUtil.Ident(caseCommentLeft, "ver"))

	for _, field := range fields {
		switch field {
		case "id", "ver":
			// already set
		case "text":
			base = base.Column(fmt.Sprintf("%s.comment text", caseCommentLeft))
		case "created_at":
			base = base.Column(storeUtil.Ident(caseCommentLeft, "created_at"))
		case "updated_at":
			base = base.Column(storeUtil.Ident(caseCommentLeft, "updated_at"))
		case "created_by":
			alias := joinCreatedBy("ccb")
			base = base.Column(fmt.Sprintf("%s.id created_by_id", alias))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) created_by_name", alias, alias))
		case "updated_by":
			alias := joinUpdatedBy("cub")
			base = base.Column(fmt.Sprintf("%s.id updated_by_id", alias))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) updated_by_name", alias, alias))
		case "edited":
			base = base.Column(fmt.Sprintf("(%s.created_at < %[1]s.updated_at) edited", caseCommentLeft))
		case "can_edit":
			if session != nil {
				base = base.Column(fmt.Sprintf("(%s.created_by = %d) can_edit", caseCommentLeft, session.GetUserId()))
			}
		case "author":
			alias := joinAuthor("au")
			base = base.Column(fmt.Sprintf("%s.id AS contact_id", alias))
			base = base.Column(fmt.Sprintf("%s.common_name AS contact_name", alias))
		case "case_id":
			base = base.Column(storeUtil.Ident(caseCommentLeft, "case_id"))
		case "role_ids":
			// skip
		default:
			return base, errors.New(fmt.Sprintf("unknown field: %s", field))
		}
	}

	return base, nil
}

// deprecated function. use until case service is refactored
func buildCommentSelectColumnsAndPlan(
	base sq.SelectBuilder,
	left string,
	fields []string,
	session auth.Auther,
) (sq.SelectBuilder, []func(comment *_go.CaseComment) any, error) {
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
			base = base.Column(storeUtil.Ident(left, "id"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return &comment.Id
			})
		case "ver":
			base = base.Column(storeUtil.Ident(left, "ver"))
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
			base = base.Column(storeUtil.Ident(left, "created_at"))
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
			base = base.Column(storeUtil.Ident(left, "updated_at"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanTimestamp(&comment.UpdatedAt)
			})
		case "text":
			base = base.Column(storeUtil.Ident(left, "comment"))
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
			base = base.Column(storeUtil.Ident(left, "case_id"))
			plan = append(plan, func(comment *_go.CaseComment) any {
				return scanner.ScanInt64(&comment.CaseId)
			})
		default:
			return base, nil, errors.NewDBError("postgres.case_comment.build_comment_select.cycle_fields.unknown", fmt.Sprintf("%s field is unknown", field))
		}
	}

	if len(plan) == 0 {
		return base, nil, errors.NewDBError("postgres.case_comment.build_comment_select.final_check.unknown", "no resulting columns")
	}

	return base, plan, nil
}

// deprecated function. use until case service is refactored
func buildCommentsSelectAsSubquery(auther auth.Auther, fields []string, caseAlias string) (sq.SelectBuilder, []func(link *_go.CaseComment) any, error) {
	alias := "comments"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.
		Select().
		From("cases.case_comment " + alias).
		Where(fmt.Sprintf("%s = %s", storeUtil.Ident(alias, "case_id"), storeUtil.Ident(caseAlias, "id")))
	base, err := addCaseCommentRbacCondition(auther, auth.Read, base, storeUtil.Ident(alias, "id"))
	if err != nil {
		return base, nil, err
	}
	base, plan, dbErr := buildCommentSelectColumnsAndPlan(base, alias, fields, auther)
	if dbErr != nil {
		return base, nil, dbErr
	}
	base = storeUtil.ApplyPaging(1, defaults.DefaultSearchSize, base)
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
		return nil, errors.New("error creating comment case interface to the case_comment table, main store is nil")
	}
	return &CaseCommentStore{storage: store}, nil
}
