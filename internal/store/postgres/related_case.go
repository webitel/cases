package postgres

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/model/options/defaults"
	storeutil "github.com/webitel/cases/internal/store/util"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/util"
)

type RelatedCaseStore struct {
	storage *Store
}

const (
	relatedCaseLeft              = "rc"
	relatedCaseAlias             = "rca"
	relatedCasePriorityAlias     = "rcpa"
	primaryCaseAlias             = "pca"
	primaryCasePriorityAlias     = "pcpa"
	relatedCaseCreatedByAlias    = "cb"
	relatedCaseUpdatedByAlias    = "ub"
	relatedCaseObjClassScopeName = model.ScopeCases
)

// Create implements store.RelatedCaseStore for creating a new related case.
func (r *RelatedCaseStore) Create(
	rpc options.Creator,
	relation *cases.RelationType,
	userID int64,
) (*cases.RelatedCase, error) {
	// Establish database connection
	d, err := r.storage.Database()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.related_case.create.database_connection_error: %v", err))
	}

	// Build SQLizer
	queryBuilder, plan, err := r.buildCreateRelatedCaseSqlizer(rpc, relation, userID)
	if err != nil {
		return nil, err
	}

	// Convert queryBuilder to SQL
	query, args, sqlErr := queryBuilder.ToSql()
	query = storeutil.CompactSQL(query)

	if sqlErr != nil {
		return nil, errors.New("error building create query", errors.WithCause(sqlErr))
	}

	// Execute the query and scan the results
	relatedCase := &cases.RelatedCase{}
	scanArgs := convertToRelatedCaseScanArgs(plan, relatedCase)

	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, ParseError(err)
	}

	return relatedCase, nil
}

// buildCreateRelatedCaseSqlizer builds the insert and select SQLizer for creating related cases.
func (r *RelatedCaseStore) buildCreateRelatedCaseSqlizer(
	rpc options.Creator,
	relation *cases.RelationType,
	inputUserID int64,
) (sq.Sqlizer, []func(*cases.RelatedCase) any, error) {

	userID := rpc.GetAuthOpts().GetUserId()
	if createdBy := inputUserID; createdBy != 0 {
		userID = createdBy
	}
	// Insert query
	insertBuilder := sq.
		Insert("cases.related_case").
		Columns("dc", "primary_case_id", "related_case_id", "relation_type", "created_at", "created_by", "updated_at", "updated_by").
		Values(
			rpc.GetAuthOpts().GetDomainId(), // dc
			rpc.GetParentID(),               // primary_case_id
			rpc.GetChildID(),                // related_case_id
			relation,                        // relation_type
			rpc.RequestTime(),               // created_at
			userID,                          // created_by
			rpc.RequestTime(),               // updated_at
			userID,                          // updated_by
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	// Convert insertBuilder to SQL to use it within a CTE
	insertSQL, insertArgs, err := insertBuilder.ToSql()
	if err != nil {
		return nil, nil, errors.New("error building insert query", errors.WithCause(err))
	}

	// Use the insert SQL as a CTE prefix for the main select query
	ctePrefix := sq.Expr("WITH rc AS ("+insertSQL+")", insertArgs...)
	// Build select clause and scan plan dynamically using buildCommentSelectColumnsAndPlan
	selectBuilder := sq.Select()

	// Use helper to build select columns and scan plan
	selectBuilder, plan, dbErr := buildRelatedCasesSelectColumnsAndPlan(selectBuilder, relatedCaseLeft, rpc.GetFields())
	if dbErr != nil {
		return nil, nil, dbErr
	}
	// Combine the CTE with the select query
	sqBuilder := selectBuilder.
		From(relatedCaseLeft).
		PrefixExpr(ctePrefix)

	return sqBuilder, plan, nil
}

// Delete implements store.RelatedCaseStore for deleting a related case.
func (r *RelatedCaseStore) Delete(
	rpc options.Deleter,
) error {
	// Get database connection
	d, err := r.storage.Database()
	if err != nil {
		return errors.Internal(fmt.Sprintf("postgres.related_case.delete.database_connection_error: %v", err))
	}

	// Build the delete query
	query, args, err := r.buildDeleteRelatedCaseQuery(rpc)
	if err != nil {
		return err
	}

	// Execute the query
	res, err := d.Exec(rpc, query, args...)
	if err != nil {
		return ParseError(err)
	}

	// Check if any rows were affected
	if res.RowsAffected() == 0 {
		return errors.NotFound("not found related case to delete, no rows affected")
	}

	return nil
}

func (c RelatedCaseStore) buildDeleteRelatedCaseQuery(rpc options.Deleter) (string, []interface{}, error) {
	query := deleteRelatedCaseQuery
	args := []interface{}{rpc.GetIDs(), rpc.GetAuthOpts().GetDomainId(), rpc.GetParentID()}
	return query, args, nil
}

var deleteRelatedCaseQuery = storeutil.CompactSQL(`
	DELETE FROM cases.related_case
	WHERE id = ANY($1) AND dc = $2 AND (primary_case_id = $3 OR related_case_id = $3)
`)

// List implements store.RelatedCaseStore for fetching related cases.
func (r *RelatedCaseStore) List(
	rpc options.Searcher,
) (*cases.RelatedCaseList, error) {
	// Get database connection
	d, err := r.storage.Database()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.related_case.list.database_connection_error: %v", err))
	}

	// Build the query and scan plan
	queryBuilder, planBuilder, err := r.buildListRelatedCaseSqlizer(rpc)
	if err != nil {
		return nil, err
	}

	// Convert queryBuilder to SQL
	query, args, sqlErr := queryBuilder.ToSql()
	if sqlErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.related_case.list.query_build_error: %v", sqlErr))
	}

	// Execute the query
	rows, execErr := d.Query(rpc, query, args...)
	if execErr != nil {
		return nil, ParseError(err)
	}
	defer rows.Close()

	// Prepare to scan rows
	var relatedCases []*cases.RelatedCase
	count := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		// Stop fetching more rows if size limit is reached
		if !fetchAll && count >= int(rpc.GetSize()) {
			next = true
			break
		}

		relatedCase := &cases.RelatedCase{}
		scanArgs := planBuilder(relatedCase)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, errors.Internal(fmt.Sprintf("postgres.related_case.list.row_scan_error: %v", err))
		}

		// Parse and reverse relation type
		parsedRelationType, parseErr := r.ParseRelationTypeWithReversion(relatedCase.RelationType.String())
		if parseErr != nil {
			return nil, errors.Internal(fmt.Sprintf("postgres.related_case.list.relation_parse_error: %v", parseErr))
		}
		relatedCase.RelationType = parsedRelationType

		relatedCases = append(relatedCases, relatedCase)
		count++
	}

	// Return the paginated list
	return &cases.RelatedCaseList{
		Page: int64(rpc.GetPage()),
		Next: next,
		Data: relatedCases,
	}, nil
}

// ParseRelationTypeWithReversion determines the relation type based on parent-case matching.
func (r *RelatedCaseStore) ParseRelationTypeWithReversion(
	rawType string,
) (cases.RelationType, error) {
	switch rawType {
	case "RELATION_TYPE_UNSPECIFIED":
		return cases.RelationType_RELATION_TYPE_UNSPECIFIED, nil
	case "DUPLICATES":
		return cases.RelationType_DUPLICATES, nil
	case "IS_DUPLICATED_BY":
		return cases.RelationType_IS_DUPLICATED_BY, nil
	case "BLOCKS":
		return cases.RelationType_BLOCKS, nil
	case "IS_BLOCKED_BY":
		return cases.RelationType_IS_BLOCKED_BY, nil
	case "CAUSES":
		return cases.RelationType_CAUSES, nil
	case "IS_CAUSED_BY":
		return cases.RelationType_IS_CAUSED_BY, nil
	case "IS_CHILD_OF":
		return cases.RelationType_IS_CHILD_OF, nil
	case "IS_PARENT_OF":
		return cases.RelationType_IS_PARENT_OF, nil
	case "RELATES_TO":
		return cases.RelationType_RELATES_TO, nil
	default:
		return cases.RelationType_RELATION_TYPE_UNSPECIFIED, fmt.Errorf("invalid relation type: %s", rawType)
	}
}

// buildListRelatedCaseSqlizer dynamically builds the SELECT query for related cases.
func (r *RelatedCaseStore) buildListRelatedCaseSqlizer(
	rpc options.Searcher,
) (sq.SelectBuilder, func(*cases.RelatedCase) []any, error) {

	// Start building the base query
	queryBuilder := sq.Select().
		From("cases.related_case AS rc").
		Where(sq.Eq{"rc.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Filter by parent case if provided
	caseIDFilters := rpc.GetFilter("case_id")
	if len(caseIDFilters) == 0 {
		return queryBuilder, nil, errors.InvalidArgument("case id required")
	}
	f := caseIDFilters[0]
	if (f.Operator == "=" || f.Operator == "") && f.Value != "" {
		parentId, perr := strconv.ParseInt(f.Value, 10, 64)
		if perr != nil || parentId == 0 {
			return queryBuilder, nil, errors.InvalidArgument(
				"case id required", errors.WithID(
					"postgres.case_timeline.build_case_timeline_sqlizer.check_args.case_id"),
			)
		}
		queryBuilder = queryBuilder.Where(sq.Or{
			sq.Eq{"rc.primary_case_id": parentId},
			sq.Eq{"rc.related_case_id": parentId},
		})
	}

	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"rc.id": rpc.GetIDs()})
	}

	// -------- Apply sorting by creation date ----------
	queryBuilder = queryBuilder.OrderBy("created_at ASC")

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = storeutil.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Build columns dynamically using helper
	queryBuilder, plan, err := buildRelatedCasesSelectColumnsAndPlan(queryBuilder, relatedCaseLeft, rpc.GetFields())
	if err != nil {
		return queryBuilder, nil, err
	}

	// Define scan plan function
	planBuilder := func(rc *cases.RelatedCase) []any {
		var scanArgs []any
		for _, scan := range plan {
			scanArgs = append(scanArgs, scan(rc))
		}
		return scanArgs
	}

	return queryBuilder, planBuilder, nil
}

func (r *RelatedCaseStore) Update(
	rpc options.Updator,
	input *cases.InputRelatedCase,
	userID int64,
) (*cases.RelatedCase, error) {
	// Get database connection
	d, err := r.storage.Database()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.related_case.update.database_connection_error: %v", err))
	}

	// Build update SQLizer
	queryBuilder, plan, err := r.buildUpdateRelatedCaseSqlizer(rpc, input, userID)
	if err != nil {
		return nil, err
	}

	// Convert queryBuilder to SQL
	query, args, sqlErr := queryBuilder.ToSql()
	if sqlErr != nil {
		return nil, errors.New("error building update query", errors.WithCause(sqlErr))
	}

	// Prepare object for result scanning
	updatedCase := &cases.RelatedCase{}
	scanArgs := convertToRelatedCaseScanArgs(plan, updatedCase)

	// Execute query and scan the result
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, ParseError(err)
	}

	return updatedCase, nil
}

// buildUpdateRelatedCaseSqlizer dynamically builds the update query for related cases.
func (r *RelatedCaseStore) buildUpdateRelatedCaseSqlizer(
	rpc options.Updator,
	input *cases.InputRelatedCase,
	inputUserID int64,
) (sq.Sqlizer, []func(*cases.RelatedCase) any, error) {
	// Ensure "id" and "ver" are included
	fields := rpc.GetFields()
	fields = util.EnsureIdAndVerField(rpc.GetFields())

	userID := rpc.GetAuthOpts().GetUserId()
	if util.ContainsField(rpc.GetMask(), "userID") {
		if updatedBy := inputUserID; updatedBy != 0 {
			userID = updatedBy
		}
	}

	// Start building the update query
	updateBuilder := sq.Update("cases.related_case").
		PlaceholderFormat(sq.Dollar).
		Set("relation_type", input.RelationType).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", userID).
		Set("ver", sq.Expr("ver + 1")).
		Where(sq.Eq{
			"id":  rpc.GetEtags()[0].GetOid(),
			"ver": rpc.GetEtags()[0].GetVer(),
			"dc":  rpc.GetAuthOpts().GetDomainId(),
		})

	for _, mask := range rpc.GetMask() {
		switch mask {
		case "primaryCaseId":
			updateBuilder = updateBuilder.Set("primary_case_id", input.GetPrimaryCase().GetId())
		case "relatedCaseId":
			updateBuilder = updateBuilder.Set("related_case_id", input.GetRelatedCase().GetId())
		}
	}

	// Prepare a SELECT query with the updated values
	selectBuilder := sq.Select().
		PrefixExpr(sq.Expr(fmt.Sprintf(`WITH %s AS (?)`, relatedCaseLeft), updateBuilder.Suffix("RETURNING *"))).
		From(relatedCaseLeft)

	// Use helper function to dynamically build columns and scan plan
	selectBuilder, plan, err := buildRelatedCasesSelectColumnsAndPlan(selectBuilder, relatedCaseLeft, fields)
	if err != nil {
		return nil, nil, err
	}

	return selectBuilder, plan, nil
}

func buildRelatedCasesSelectColumnsAndPlan(
	base sq.SelectBuilder,
	left string,
	fields []string,
) (sq.SelectBuilder, []func(relatedCase *cases.RelatedCase) any, error) {
	var (
		plan []func(relatedCase *cases.RelatedCase) any

		joinCreatedBy = func() {
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.created_by", relatedCaseCreatedByAlias, left))
		}
		joinUpdatedBy = func() {
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.updated_by", relatedCaseUpdatedByAlias, left))
		}
		joinRelatedCase = func() {
			base = base.
				LeftJoin(fmt.Sprintf("cases.case %s ON %[1]s.id = %s.related_case_id", relatedCaseAlias, left)).
				LeftJoin(fmt.Sprintf("cases.priority %s ON %[1]s.id = %s.priority", relatedCasePriorityAlias, relatedCaseAlias))
		}
		joinPrimaryCase = func() {
			base = base.
				LeftJoin(fmt.Sprintf("cases.case %s ON %[1]s.id = %s.primary_case_id", primaryCaseAlias, left)).
				LeftJoin(fmt.Sprintf("cases.priority %s ON %[1]s.id = %s.priority", primaryCasePriorityAlias, primaryCaseAlias))
		}
	)

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(storeutil.Ident(left, "id"))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return &rc.Id
			})
		case "ver":
			base = base.Column(storeutil.Ident(left, "ver"))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return &rc.Ver
			})
		case "created_by":
			joinCreatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, coalesce(%[1]s.name, %[1]s.username))::text created_by", relatedCaseCreatedByAlias))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return scanner.ScanRowLookup(&rc.CreatedBy)
			})
		case "created_at":
			base = base.Column(storeutil.Ident(left, "created_at"))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return scanner.ScanTimestamp(&rc.CreatedAt)
			})
		case "updated_by":
			joinUpdatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, coalesce(%[1]s.name, %[1]s.username))::text updated_by", relatedCaseUpdatedByAlias))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return scanner.ScanRowLookup(&rc.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(storeutil.Ident(left, "updated_at"))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return scanner.ScanTimestamp(&rc.UpdatedAt)
			})
		//case "relation":
		//	base = base.Column(util2.Ident(left, "relation_type"))
		//	plan = append(plan, func(rc *cases.RelatedCase) any {
		//		return scanner.TextDecoder(func(src []byte) error {
		//			if len(src) == 0 {
		//				rc.RelationType = 0
		//
		//				return nil
		//			}
		//			var relType int32
		//			s := string(src)
		//			_, err := fmt.Sscanf(s, "%d", &relType)
		//			if err != nil {
		//				return err
		//			}
		//			rc.RelationType = cases.RelationType(relType)
		//			return nil
		//		})
		//	})
		case "relation":
			base = base.Column(storeutil.Ident(left, "relation_type"))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				// Scan directly into int32 (handles NULL as 0 automatically)
				return (*int32)(&rc.RelationType)
			})
		case "related_case":
			joinRelatedCase()
			base = base.Column(fmt.Sprintf(
				"ROW(%[1]s.id, %[1]s.name, %[1]s.subject, %[1]s.ver, %[2]s.color)::text related_case",
				relatedCaseAlias, relatedCasePriorityAlias))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return scanner.ScanRelatedCaseLookup(&rc.RelatedCase)
			})
		case "primary_case":
			joinPrimaryCase()
			base = base.Column(fmt.Sprintf(
				"ROW(%[1]s.id, %[1]s.name, %[1]s.subject, %[1]s.ver, %[2]s.color)::text primary_case",
				primaryCaseAlias, primaryCasePriorityAlias))
			plan = append(plan, func(rc *cases.RelatedCase) any {
				return scanner.ScanRelatedCaseLookup(&rc.PrimaryCase)
			})
		default:
			return base, nil, errors.Internal(fmt.Sprintf("store.related_case.build_columns.unknown_field: Unknown field: %s", field))
		}
	}

	if len(plan) == 0 {
		return base, nil, errors.Internal("store.related_case.build_columns.no_fields: No columns specified")
	}

	return base, plan, nil
}

// convertToRelatedCaseScanArgs converts scan plan to arguments for Scan function.
func convertToRelatedCaseScanArgs(plan []func(*cases.RelatedCase) any, rc *cases.RelatedCase) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(rc))
	}
	return scanArgs
}

func NewRelatedCaseStore(store *Store) (store.RelatedCaseStore, error) {
	if store == nil {
		return nil, errors.Internal("postgres.new_related_case.check.bad_arguments: error creating related case interface, main store is nil")
	}
	return &RelatedCaseStore{storage: store}, nil
}

func buildRelatedCasesSelectAsSubquery(auther auth.Auther, fields []string, caseAlias string) (sq.SelectBuilder, []func(*cases.RelatedCase) any, error) {
	alias := "related"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.
		Select().
		From("cases.related_case " + alias).
		Where(fmt.Sprintf("%s = %s", storeutil.Ident(alias, "primary_case_id"), storeutil.Ident(caseAlias, "id")))

	base, err := addRelatedCaseRbacCondition(auther, auth.Read, base, storeutil.Ident(alias, "id"))
	if err != nil {
		return base, nil, err
	}
	base, plan, dbErr := buildRelatedCasesSelectColumnsAndPlan(base, alias, fields)
	if dbErr != nil {
		return base, nil, dbErr
	}
	base = storeutil.ApplyPaging(1, defaults.DefaultSearchSize, base)
	return base, plan, nil
}

func addRelatedCaseRbacCondition(auth auth.Auther, access auth.AccessMode, query sq.SelectBuilder, dependencyColumn string) (sq.SelectBuilder, error) {
	if auth != nil && auth.IsRbacCheckRequired(relatedCaseObjClassScopeName, access) {
		return query.Where(sq.Expr(fmt.Sprintf(
			"EXISTS(SELECT acl.object FROM cases.case_acl acl WHERE acl.dc = ? AND acl.object = %s AND acl.subject = any( ?::int[]) AND acl.access & ? = ? LIMIT 1)",
			dependencyColumn),
			auth.GetDomainId(), pq.Array(auth.GetRoles()), int64(access), int64(access))), nil
	}
	return query, nil
}
