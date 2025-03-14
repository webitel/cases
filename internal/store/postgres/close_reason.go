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

type CloseReasonScan func(reason *_go.CloseReason) any

const (
	crLeft                 = "cr"
	closeReasonDefaultSort = "name"
)

type CloseReason struct {
	storage *Store
}

// Helper function to convert plan to scan arguments.
func convertToCloseReasonScanArgs(plan []CloseReasonScan, reason *_go.CloseReason) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(reason))
	}
	return scanArgs
}

// Helper function to dynamically build select columns and plan.
func buildCloseReasonSelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []CloseReasonScan, error) {
	var plan []CloseReasonScan
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(crLeft, "id"))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return &reason.Id
			})
		case "name":
			base = base.Column(util2.Ident(crLeft, "name"))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return &reason.Name
			})
		case "description":
			base = base.Column(util2.Ident(crLeft, "description"))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return scanner.ScanText(&reason.Description)
			})
		case "created_at":
			base = base.Column(util2.Ident(crLeft, "created_at"))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return scanner.ScanTimestamp(&reason.CreatedAt)
			})
		case "updated_at":
			base = base.Column(util2.Ident(crLeft, "updated_at"))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return scanner.ScanTimestamp(&reason.UpdatedAt)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by",
				crLeft))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return scanner.ScanRowLookup(&reason.UpdatedBy)
			})
		case "created_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.created_by) created_by",
				crLeft))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return scanner.ScanRowLookup(&reason.CreatedBy)
			})
		case "close_reason_id":
			base = base.Column(util2.Ident(crLeft, "close_reason_id"))
			plan = append(plan, func(reason *_go.CloseReason) any {
				return &reason.CloseReasonGroupId
			})
		default:
			return base, nil, dberr.NewDBInternalError("postgres.close_reason.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return base, plan, nil
}

func (s *CloseReason) buildCreateCloseReasonQuery(
	rpc options.CreateOptions,
	reason *_go.CloseReason,
) (sq.SelectBuilder, []CloseReasonScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	insertBuilder := sq.Insert("cases.close_reason").
		Columns("name", "dc", "created_at", "description", "created_by", "updated_at", "updated_by", "close_reason_id").
		Values(
			reason.Name,
			rpc.GetAuthOpts().GetDomainId(),
			rpc.RequestTime(),
			sq.Expr("NULLIF(?, '')", reason.Description),
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
			reason.CloseReasonGroupId,
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	// Convert the INSERT query into a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.close_reason.create.query_build_error", err)
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH cr AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, plan, err := buildCloseReasonSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(crLeft)

	return selectBuilder, plan, nil
}

func (s *CloseReason) Create(rpc options.CreateOptions, input *_go.CloseReason) (*_go.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildCreateCloseReasonQuery(rpc, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &_go.CloseReason{}
	scanArgs := convertToCloseReasonScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.execution_error", err)
	}

	return tempAdd, nil
}

func (s *CloseReason) buildUpdateCloseReasonQuery(
	rpc options.UpdateOptions,
	input *_go.CloseReason,
) (sq.SelectBuilder, []CloseReasonScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Start the UPDATE query
	updateBuilder := sq.Update("cases.close_reason").
		PlaceholderFormat(sq.Dollar).
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
		case "close_reason_id":
			if input.CloseReasonGroupId != 0 {
				updateBuilder = updateBuilder.Set("close_reason_id", input.CloseReasonGroupId)
			}
		}
	}

	// Generate the CTE for the update operation
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.close_reason.update.query_build_error", err)
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH cr AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using buildCloseReasonSelectColumnsAndPlan
	selectBuilder, plan, err := buildCloseReasonSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("cr")

	return selectBuilder, plan, nil
}

func (s *CloseReason) Update(rpc options.UpdateOptions, input *_go.CloseReason) (*_go.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.input.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildUpdateCloseReasonQuery(rpc, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.input.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.input.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &_go.CloseReason{}
	scanArgs := convertToCloseReasonScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.input.execution_error", err)
	}

	return tempAdd, nil
}

func (s *CloseReason) buildListCloseReasonQuery(
	rpc options.SearchOptions,
	closeReasonId int64,
) (sq.SelectBuilder, []CloseReasonScan, error) {
	queryBuilder := sq.Select().
		From("cases.close_reason AS cr").
		Where(sq.Eq{"cr.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Add ID filter if provided
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cr.id": rpc.GetIDs()})
	}

	// Add name filter if provided
	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "cr.name")
	}

	// Add close reason group filter if provided
	if closeReasonId != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cr.close_reason_id": closeReasonId})
	}

	// -------- Apply sorting ----------
	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, closeReasonDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, plan, err := buildCloseReasonSelectColumnsAndPlan(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.close_reason.search.query_build_error", err)
	}

	return queryBuilder, plan, nil
}

func (s *CloseReason) List(rpc options.SearchOptions, closeReasonId int64) (*_go.CloseReasonList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildListCloseReasonQuery(rpc, closeReasonId)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.query_build_error", err)
	}
	query = util2.CompactSQL(query)

	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.execution_error", err)
	}
	defer rows.Close()

	var reasons []*_go.CloseReason
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		reason := &_go.CloseReason{}
		scanArgs := convertToCloseReasonScanArgs(plan, reason)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.close_reason.list.row_scan_error", err)
		}

		reasons = append(reasons, reason)
		lCount++
	}

	return &_go.CloseReasonList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: reasons,
	}, nil
}

func (s *CloseReason) buildDeleteCloseReasonQuery(
	rpc options.DeleteOptions,
) (sq.DeleteBuilder, error) {
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, dberr.NewDBInternalError("postgres.close_reason.delete.missing_ids", fmt.Errorf("no IDs provided for deletion"))
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.close_reason").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	return deleteBuilder, nil
}

func (s *CloseReason) Delete(rpc options.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.close_reason.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := s.buildDeleteCloseReasonQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.close_reason.delete.query_build_error", err)
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("postgres.close_reason.delete.query_to_sql_error", err)
	}

	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("postgres.close_reason.delete.execution_error", execErr)
	}

	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.close_reason.delete.no_rows_affected")
	}

	return nil
}

func NewCloseReasonStore(store *Store) (store.CloseReasonStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_close_reason.check.bad_arguments",
			"error creating close_reason interface, main store is nil")
	}
	return &CloseReason{storage: store}, nil
}
