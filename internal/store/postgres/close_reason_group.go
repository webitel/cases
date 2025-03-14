package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"github.com/webitel/cases/util"
)

type CloseReasonGroup struct {
	storage *Store
}

type CloseReasonGroupScan func(group *_go.CloseReasonGroup) any

const (
	crgLeft                     = "g"
	closeReasonGroupDefaultSort = "name"
)

// Helper function to convert plan to scan arguments.
func convertToCloseReasonGroupScanArgs(plan []CloseReasonGroupScan, group *_go.CloseReasonGroup) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(group))
	}
	return scanArgs
}

// Helper function to dynamically build select columns and plan.
func buildCloseReasonGroupSelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []CloseReasonGroupScan, error) {
	var plan []CloseReasonGroupScan
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(crgLeft, "id"))
			plan = append(plan, func(group *_go.CloseReasonGroup) any {
				return &group.Id
			})
		case "name":
			base = base.Column(util2.Ident(crgLeft, "name"))
			plan = append(plan, func(group *_go.CloseReasonGroup) any {
				return &group.Name
			})
		case "description":
			base = base.Column(util2.Ident(crgLeft, "description"))
			plan = append(plan, func(group *_go.CloseReasonGroup) any {
				return scanner.ScanText(&group.Description)
			})
		case "created_at":
			base = base.Column(util2.Ident(crgLeft, "created_at"))
			plan = append(plan, func(group *_go.CloseReasonGroup) any {
				return scanner.ScanTimestamp(&group.CreatedAt)
			})
		case "updated_at":
			base = base.Column(util2.Ident(crgLeft, "updated_at"))
			plan = append(plan, func(group *_go.CloseReasonGroup) any {
				return scanner.ScanTimestamp(&group.UpdatedAt)
			})
		case "created_by":
			base = base.Column(
				fmt.Sprintf("(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.created_by) created_by",
					crgLeft,
				))
			plan = append(plan, func(group *_go.CloseReasonGroup) any {
				return scanner.ScanRowLookup(&group.CreatedBy)
			})
		case "updated_by":
			base = base.Column(
				fmt.Sprintf("(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by",
					crgLeft,
				))
			plan = append(plan, func(group *_go.CloseReasonGroup) any {
				return scanner.ScanRowLookup(&group.UpdatedBy)
			})
		default:
			return base, nil,
				dberr.NewDBInternalError(
					"postgres.close_reason_group.unknown_field",
					fmt.Errorf("unknown field: %s", field),
				)
		}
	}
	return base, plan, nil
}

func (s CloseReasonGroup) buildCreateCloseReasonGroupQuery(
	rpc options.CreateOptions,
	group *_go.CloseReasonGroup,
) (sq.SelectBuilder, []CloseReasonGroupScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Build the INSERT query with a RETURNING clause
	insertBuilder := sq.Insert("cases.close_reason_group").
		Columns("name", "dc", "created_at", "description", "created_by", "updated_at", "updated_by").
		Values(
			group.Name,
			rpc.GetAuthOpts().GetDomainId(),
			rpc.RequestTime(),
			sq.Expr("NULLIF(?, '')", group.Description),
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	// Convert the INSERT query into a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil,
			dberr.NewDBInternalError(
				"postgres.close_reason_group.create.query_build_error",
				err,
			)
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH g AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, plan, err := buildCloseReasonGroupSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(crgLeft)

	return selectBuilder, plan, nil
}

func (s CloseReasonGroup) Create(rpc options.CreateOptions, input *_go.CloseReasonGroup) (*_go.CloseReasonGroup, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.create.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildCreateCloseReasonGroupQuery(rpc, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.create.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &_go.CloseReasonGroup{}
	scanArgs := convertToCloseReasonGroupScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.create.execution_error", err)
	}

	return tempAdd, nil
}

func (s CloseReasonGroup) buildUpdateCloseReasonGroupQuery(
	rpc options.UpdateOptions,
	input *_go.CloseReasonGroup,
) (sq.SelectBuilder, []CloseReasonGroupScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Start the UPDATE query
	updateBuilder := sq.Update("cases.close_reason_group").
		PlaceholderFormat(sq.Dollar). // Use PostgreSQL-compatible placeholders
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": input.Id}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically add fields to the SET clause
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			if input.Name != "" {
				updateBuilder = updateBuilder.Set("name", input.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", input.Description))
		}
	}

	// Generate the CTE for the update operation
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil,
			dberr.NewDBInternalError(
				"postgres.close_reason_group.update.query_build_error",
				err,
			)
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH g AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using buildCloseReasonGroupSelectColumnsAndPlan
	selectBuilder, plan, err := buildCloseReasonGroupSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("g")

	return selectBuilder, plan, nil
}

func (s CloseReasonGroup) Update(
	rpc options.UpdateOptions,
	input *_go.CloseReasonGroup,
) (*_go.CloseReasonGroup, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.input.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildUpdateCloseReasonGroupQuery(rpc, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.input.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.input.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &_go.CloseReasonGroup{}
	scanArgs := convertToCloseReasonGroupScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.input.execution_error", err)
	}

	return tempAdd, nil
}

func (s CloseReasonGroup) buildListCloseReasonGroupQuery(
	rpc options.SearchOptions,
) (sq.SelectBuilder, []CloseReasonGroupScan, error) {

	queryBuilder := sq.Select().
		From("cases.close_reason_group AS g").
		Where(sq.Eq{"g.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Add ID filter if provided
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"g.id": rpc.GetIDs()})
	}

	// Add name filter if provided
	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "g.name")
	}

	// -------- Apply sorting ----------
	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, closeReasonGroupDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, plan, err := buildCloseReasonGroupSelectColumnsAndPlan(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, nil,
			dberr.NewDBInternalError(
				"postgres.close_reason_group.search.query_build_error",
				err,
			)
	}

	return queryBuilder, plan, nil
}

func (s CloseReasonGroup) List(rpc options.SearchOptions) (*_go.CloseReasonGroupList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.list.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildListCloseReasonGroupQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.list.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.list.query_build_error", err)
	}
	query = util2.CompactSQL(query)

	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason_group.list.execution_error", err)
	}
	defer rows.Close()

	var groups []*_go.CloseReasonGroup
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= rpc.GetSize() {
			next = true
			break
		}

		group := &_go.CloseReasonGroup{}
		scanArgs := convertToCloseReasonGroupScanArgs(plan, group)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.close_reason_group.list.row_scan_error", err)
		}

		groups = append(groups, group)
		lCount++
	}

	return &_go.CloseReasonGroupList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: groups,
	}, nil
}

func (s CloseReasonGroup) buildDeleteCloseReasonGroupQuery(
	rpc options.DeleteOptions,
) (sq.DeleteBuilder, error) {
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{},
			dberr.NewDBInternalError(
				"postgres.close_reason_group.delete.missing_ids",
				fmt.Errorf("no IDs provided for deletion"))
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.close_reason_group").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	return deleteBuilder, nil
}

func (s CloseReasonGroup) Delete(rpc options.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.close_reason_group.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := s.buildDeleteCloseReasonGroupQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.close_reason_group.delete.query_build_error", err)
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("postgres.close_reason_group.delete.query_to_sql_error", err)
	}

	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("postgres.close_reason_group.delete.execution_error", execErr)
	}

	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.close_reason_group.delete.no_rows_affected")
	}

	return nil
}

func NewCloseReasonGroupStore(store *Store) (store.CloseReasonGroupStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_close_reason_group.check.bad_arguments",
			"error creating close_reason_group interface to the close_reason_group table, main store is nil")
	}
	return &CloseReasonGroup{storage: store}, nil
}
