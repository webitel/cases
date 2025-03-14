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

type StatusScan func(status *_go.Status) any

const (
	statusLeft        = "s"
	statusDefaultSort = "name"
)

type Status struct {
	storage *Store
}

// Helper function to convert plan to scan arguments.
func convertToStatusScanArgs(plan []StatusScan, status *_go.Status) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(status))
	}
	return scanArgs
}

// Helper function to dynamically build select columns and plan.
func buildStatusSelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []StatusScan, error) {
	var plan []StatusScan
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(statusLeft, "id"))
			plan = append(plan, func(status *_go.Status) any {
				return &status.Id
			})
		case "name":
			base = base.Column(util2.Ident(statusLeft, "name"))
			plan = append(plan, func(status *_go.Status) any {
				return &status.Name
			})
		case "description":
			base = base.Column(util2.Ident(statusLeft, "description"))
			plan = append(plan, func(status *_go.Status) any {
				return scanner.ScanText(&status.Description)
			})
		case "created_at":
			base = base.Column(util2.Ident(statusLeft, "created_at"))
			plan = append(plan, func(status *_go.Status) any {
				return scanner.ScanTimestamp(&status.CreatedAt)
			})
		case "updated_at":
			base = base.Column(util2.Ident(statusLeft, "updated_at"))
			plan = append(plan, func(status *_go.Status) any {
				return scanner.ScanTimestamp(&status.UpdatedAt)
			})
		case "created_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.created_by) created_by", statusLeft))
			plan = append(plan, func(status *_go.Status) any {
				return scanner.ScanRowLookup(&status.CreatedBy)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by", statusLeft))
			plan = append(plan, func(status *_go.Status) any {
				return scanner.ScanRowLookup(&status.UpdatedBy)
			})
		default:
			return base, nil, dberr.NewDBInternalError("postgres.status.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return base, plan, nil
}

func (s *Status) buildCreateStatusQuery(
	rpc options.CreateOptions,
	input *_go.Status,
) (sq.SelectBuilder, []StatusScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Build the INSERT query with a RETURNING clause
	insertBuilder := sq.Insert("cases.status").
		Columns("name", "dc", "created_at", "description", "created_by", "updated_at", "updated_by").
		Values(
			input.Name,                                  // name
			rpc.GetAuthOpts().GetDomainId(),             // dc
			rpc.RequestTime(),                           // created_at
			sq.Expr("NULLIF(?, '')", input.Description), // NULLIF for empty description
			rpc.GetAuthOpts().GetUserId(),               // created_by
			rpc.RequestTime(),                           // updated_at
			rpc.GetAuthOpts().GetUserId(),               // updated_by
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *") // RETURNING all columns for use in the next SELECT

	// Convert the INSERT query into a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.input.create.query_build_error", err)
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH s AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, plan, err := buildStatusSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(statusLeft)

	return selectBuilder, plan, nil
}

func (s *Status) Create(rpc options.CreateOptions, input *_go.Status) (*_go.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.status.create.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildCreateStatusQuery(rpc, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.create.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &_go.Status{}
	scanArgs := convertToStatusScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.create.execution_error", err)
	}

	return tempAdd, nil
}

func (s *Status) buildUpdateStatusQuery(
	rpc options.UpdateOptions,
	input *_go.Status,
) (sq.SelectBuilder, []StatusScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Start the UPDATE query
	updateBuilder := sq.Update("cases.status").
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
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.input.update.query_build_error", err)
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH s AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using buildStatusSelectColumnsAndPlan
	selectBuilder, plan, err := buildStatusSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("s")

	return selectBuilder, plan, nil
}

func (s *Status) Update(rpc options.UpdateOptions, input *_go.Status) (*_go.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.status.input.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildUpdateStatusQuery(rpc, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.input.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.input.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &_go.Status{}
	scanArgs := convertToStatusScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.input.execution_error", err)
	}

	return tempAdd, nil
}

func (s *Status) buildListStatusQuery(
	rpc options.SearchOptions,
) (sq.SelectBuilder, []StatusScan, error) {

	queryBuilder := sq.Select().
		From("cases.status AS s").
		Where(sq.Eq{"s.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Add ID filter if provided
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"s.id": rpc.GetIDs()})
	}

	// Add name filter if provided
	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "s.name")
	}

	// -------- Apply sorting ----------
	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, statusDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, plan, err := buildStatusSelectColumnsAndPlan(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.status.search.query_build_error", err)
	}

	return queryBuilder, plan, nil
}

func (s *Status) List(rpc options.SearchOptions) (*_go.StatusList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.status.list.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildListStatusQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.list.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.list.query_build_error", err)
	}
	query = util2.CompactSQL(query)

	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.status.list.execution_error", err)
	}
	defer rows.Close()

	var statuses []*_go.Status
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		status := &_go.Status{}
		scanArgs := convertToStatusScanArgs(plan, status)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.status.list.row_scan_error", err)
		}

		statuses = append(statuses, status)
		lCount++
	}

	return &_go.StatusList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: statuses,
	}, nil
}

func (s *Status) buildDeleteStatusQuery(
	rpc options.DeleteOptions,
) (sq.DeleteBuilder, error) {
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, dberr.NewDBInternalError("postgres.status.delete.missing_ids", fmt.Errorf("no IDs provided for deletion"))
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.status").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	return deleteBuilder, nil
}

func (s *Status) Delete(rpc options.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.status.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := s.buildDeleteStatusQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.status.delete.query_build_error", err)
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("postgres.status.delete.query_to_sql_error", err)
	}

	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("postgres.status.delete.execution_error", execErr)
	}

	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.status.delete.no_rows_affected")
	}

	return nil
}

func NewStatusStore(store *Store) (store.StatusStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_status.check.bad_arguments",
			"error creating status interface, main store is nil")
	}
	return &Status{storage: store}, nil
}
