package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	api "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	util "github.com/webitel/cases/util"
)

type Priority struct {
	storage *Store
}

type PriorityScan func(priority *api.Priority) any

const (
	prioLeft            = "cp"
	priorityDefaultSort = "name"
)

// Create implements store.PriorityStore.
func (p *Priority) Create(rpc options.CreateOptions, add *api.Priority) (*api.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := p.buildCreatePriorityQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &api.Priority{}
	scanArgs := convertToPriorityScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.execution_error", err)
	}

	return tempAdd, nil
}

func (p *Priority) buildCreatePriorityQuery(
	rpc options.CreateOptions,
	priority *api.Priority,
) (sq.SelectBuilder, []PriorityScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Build the INSERT query with a RETURNING clause
	insertBuilder := sq.Insert("cases.priority").
		Columns("name", "dc", "created_at", "description", "created_by", "updated_at", "updated_by", "color").
		Values(
			priority.Name,                                  // name
			rpc.GetAuthOpts().GetDomainId(),                // dc
			rpc.RequestTime(),                              // created_at
			sq.Expr("NULLIF(?, '')", priority.Description), // NULLIF for empty description
			rpc.GetAuthOpts().GetUserId(),                  // created_by
			rpc.RequestTime(),                              // updated_at
			rpc.GetAuthOpts().GetUserId(),                  // updated_by
			priority.Color,                                 // color
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *") // RETURNING all columns for use in the next SELECT

	// Convert the INSERT query into a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.priority.create.query_build_error", err)
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH cp AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, plan, err := buildPrioritySelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(prioLeft)

	return selectBuilder, plan, nil
}

func (p *Priority) Delete(rpc options.DeleteOptions) error {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.priority.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := p.buildDeletePriorityQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.priority.delete.query_build_error", err)
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("postgres.priority.delete.query_to_sql_error", err)
	}

	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("postgres.priority.delete.execution_error", execErr)
	}

	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.priority.delete.no_rows_affected")
	}

	return nil
}

func (p *Priority) buildDeletePriorityQuery(
	rpc options.DeleteOptions,
) (sq.DeleteBuilder, error) {
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, dberr.NewDBInternalError("postgres.priority.delete.missing_ids", fmt.Errorf("no IDs provided for deletion"))
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.priority").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	return deleteBuilder, nil
}

// List implements store.PriorityStore.
func (p *Priority) List(
	rpc options.SearchOptions,
	notInSla int64,
	inSla int64,
) (*api.PriorityList, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := p.buildListPriorityQuery(rpc, notInSla, inSla)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.query_build_error", err)
	}
	query = util2.CompactSQL(query)

	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.execution_error", err)
	}
	defer rows.Close()

	var priorities []*api.Priority
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		priority := &api.Priority{}
		scanArgs := convertToPriorityScanArgs(plan, priority)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.priority.list.row_scan_error", err)
		}

		priorities = append(priorities, priority)
		lCount++
	}

	return &api.PriorityList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: priorities,
	}, nil
}

func (p *Priority) buildListPriorityQuery(
	rpc options.SearchOptions,
	notInSla int64,
	inSla int64,
) (sq.SelectBuilder, []PriorityScan, error) {

	queryBuilder := sq.Select().
		From("cases.priority AS cp").
		Where(sq.Eq{"cp.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Add ID filter if provided
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cp.id": rpc.GetIDs()})
	}

	// Add name filter if provided
	if name, ok := rpc.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "cp.name")
	}

	// Add NOT IN SLA condition if `notInSla` is not 0
	if notInSla != 0 {
		queryBuilder = queryBuilder.Where(sq.Expr(`
				(NOT EXISTS (
					SELECT 1
					FROM cases.sla_condition sc
					JOIN cases.priority_sla_condition psc ON sc.id = psc.sla_condition_id
					WHERE sc.sla_id = ? AND psc.priority_id = cp.id
	))
			`, notInSla))
	}

	if inSla != 0 {
		queryBuilder = queryBuilder.Where(sq.Expr(`
			(EXISTS (
				SELECT 1
				FROM cases.priority_sla_condition psc
				WHERE psc.sla_condition_id = ?
				AND psc.priority_id = cp.id
			)
			OR NOT EXISTS (
				SELECT 1
				FROM cases.sla_condition sc
				JOIN cases.priority_sla_condition psc ON sc.id = psc.sla_condition_id
				WHERE sc.sla_id = (
					SELECT sla_id FROM cases.sla_condition WHERE id = ?
				)
				AND psc.priority_id = cp.id
	))
		`, inSla, inSla))
	}

	// Filter by IDs
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cp.id": rpc.GetIDs()})
	}

	// -------- Apply sorting ----------
	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, priorityDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, plan, err := buildPrioritySelectColumnsAndPlan(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.priority.search.query_build_error", err)
	}

	return queryBuilder, plan, nil
}

// Update implements store.PriorityStore.
func (p *Priority) Update(rpc options.UpdateOptions, update *api.Priority) (*api.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := p.buildUpdatePriorityQuery(rpc, update)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &api.Priority{}
	scanArgs := convertToPriorityScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.execution_error", err)
	}

	return tempAdd, nil
}

func (p *Priority) buildUpdatePriorityQuery(
	rpc options.UpdateOptions,
	priority *api.Priority,
) (sq.SelectBuilder, []PriorityScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Start the UPDATE query
	updateBuilder := sq.Update("cases.priority").
		PlaceholderFormat(sq.Dollar). // Use PostgreSQL-compatible placeholders
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": priority.Id}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically add fields to the `SET` clause
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			if priority.Name != "" {
				updateBuilder = updateBuilder.Set("name", priority.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", priority.Description))
		case "color":
			if priority.Color != "" {
				updateBuilder = updateBuilder.Set("color", priority.Color)
			}
		}
	}

	// Generate the CTE for the update operation
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.priority.update.query_build_error", err)
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH cp AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using `buildPrioritySelectColumnsAndPlan`
	selectBuilder, plan, err := buildPrioritySelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("cp")

	return selectBuilder, plan, nil
}

// Helper function to convert plan to scan arguments.
func convertToPriorityScanArgs(plan []PriorityScan, priority *api.Priority) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(priority))
	}
	return scanArgs
}

// Helper function to dynamically build select columns and plan.
func buildPrioritySelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []PriorityScan, error) {
	var plan []PriorityScan
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(prioLeft, "id"))
			plan = append(plan, func(priority *api.Priority) any {
				return &priority.Id
			})
		case "name":
			base = base.Column(util2.Ident(prioLeft, "name"))
			plan = append(plan, func(priority *api.Priority) any {
				return &priority.Name
			})
		case "description":
			base = base.Column(util2.Ident(prioLeft, "description"))
			plan = append(plan, func(priority *api.Priority) any {
				return scanner.ScanText(&priority.Description)
			})
		case "created_at":
			base = base.Column(util2.Ident(prioLeft, "created_at"))
			plan = append(plan, func(priority *api.Priority) any {
				return scanner.ScanTimestamp(&priority.CreatedAt)
			})
		case "updated_at":
			base = base.Column(util2.Ident(prioLeft, "updated_at"))
			plan = append(plan, func(priority *api.Priority) any {
				return scanner.ScanTimestamp(&priority.UpdatedAt)
			})
		case "created_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.created_by) created_by", prioLeft))
			plan = append(plan, func(priority *api.Priority) any {
				return scanner.ScanRowLookup(&priority.CreatedBy)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by", prioLeft))
			plan = append(plan, func(priority *api.Priority) any {
				return scanner.ScanRowLookup(&priority.UpdatedBy)
			})
		case "color":
			base = base.Column(util2.Ident(prioLeft, "color"))
			plan = append(plan, func(priority *api.Priority) any {
				return &priority.Color
			})
		default:
			return base, nil, dberr.NewDBInternalError("postgres.priority.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return base, plan, nil
}

func NewPriorityStore(store *Store) (store.PriorityStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_priority.check.bad_arguments",
			"error creating priority interface to the status_condition table, main store is nil")
	}
	return &Priority{storage: store}, nil
}
