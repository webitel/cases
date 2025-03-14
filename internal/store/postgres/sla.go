package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"github.com/webitel/cases/util"
)

type SLAScan func(sla *cases.SLA) any

const (
	slaLeft        = "s"
	slaDefaultSort = "name"
)

type SLAStore struct {
	storage *Store
}

// Helper function to convert plan to scan arguments.
func convertToSLAScanArgs(plan []SLAScan, sla *cases.SLA) []any {
	var scanArgs []any
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(sla))
	}
	return scanArgs
}

// Helper function to dynamically build select columns and plan.
func buildSLASelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []SLAScan, error) {
	var plan []SLAScan
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(slaLeft, "id"))
			plan = append(plan, func(sla *cases.SLA) any {
				return &sla.Id
			})
		case "name":
			base = base.Column(util2.Ident(slaLeft, "name"))
			plan = append(plan, func(sla *cases.SLA) any {
				return &sla.Name
			})
		case "description":
			base = base.Column(util2.Ident(slaLeft, "description"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanText(&sla.Description)
			})
		case "valid_from":
			base = base.Column(util2.Ident(slaLeft, "valid_from"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.ValidFrom)
			})
		case "valid_to":
			base = base.Column(util2.Ident(slaLeft, "valid_to"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.ValidTo)
			})
		case "calendar":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM flow.calendar WHERE id = %s.calendar_id) calendar", slaLeft))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanRowLookup(&sla.Calendar)
			})
		case "reaction_time":
			base = base.Column(util2.Ident(slaLeft, "reaction_time"))
			plan = append(plan, func(sla *cases.SLA) any {
				return &sla.ReactionTime
			})
		case "resolution_time":
			base = base.Column(util2.Ident(slaLeft, "resolution_time"))
			plan = append(plan, func(sla *cases.SLA) any {
				return &sla.ResolutionTime
			})
		case "created_at":
			base = base.Column(util2.Ident(slaLeft, "created_at"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.CreatedAt)
			})
		case "updated_at":
			base = base.Column(util2.Ident(slaLeft, "updated_at"))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanTimestamp(&sla.UpdatedAt)
			})
		case "created_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.created_by) created_by", slaLeft))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanRowLookup(&sla.CreatedBy)
			})
		case "updated_by":
			base = base.Column(fmt.Sprintf("(SELECT ROW(id, name)::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by", slaLeft))
			plan = append(plan, func(sla *cases.SLA) any {
				return scanner.ScanRowLookup(&sla.UpdatedBy)
			})
		default:
			return base, nil, dberr.NewDBInternalError("postgres.sla.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return base, plan, nil
}

func (s *SLAStore) buildCreateSLAQuery(
	rpc options.CreateOptions,
	sla *cases.SLA,
) (sq.SelectBuilder, []SLAScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Build the INSERT query with a RETURNING clause
	insertBuilder := sq.Insert("cases.sla").
		Columns(
			"name", "dc", "created_at",
			"description", "created_by", "updated_at",
			"updated_by", "valid_from", "valid_to",
			"calendar_id", "reaction_time", "resolution_time",
		).
		Values(
			sla.Name,
			rpc.GetAuthOpts().GetDomainId(),
			rpc.RequestTime(),
			sq.Expr("NULLIF(?, '')", sla.Description),
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
			util.LocalTime(sla.ValidFrom),
			util.LocalTime(sla.ValidTo),
			sla.Calendar.Id,
			sla.ReactionTime,
			sla.ResolutionTime,
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *") // RETURNING all columns for use in the next SELECT

	// Convert the INSERT query into a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.sla.create.query_build_error", err)
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH s AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, plan, err := buildSLASelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(slaLeft)

	return selectBuilder, plan, nil
}

func (s *SLAStore) Create(rpc options.CreateOptions, add *cases.SLA) (*cases.SLA, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildCreateSLAQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &cases.SLA{}
	scanArgs := convertToSLAScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.create.execution_error", err)
	}

	return tempAdd, nil
}

func (s *SLAStore) buildUpdateSLAQuery(
	rpc options.UpdateOptions,
	sla *cases.SLA,
) (sq.SelectBuilder, []SLAScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	// Start the UPDATE query
	updateBuilder := sq.Update("cases.sla").
		PlaceholderFormat(sq.Dollar). // Use PostgreSQL-compatible placeholders
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": sla.Id}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically add fields to the SET clause
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			if sla.Name != "" {
				updateBuilder = updateBuilder.Set("name", sla.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", sla.Description))
		case "valid_from":
			updateBuilder = updateBuilder.Set("valid_from", util.LocalTime(sla.ValidFrom))
		case "valid_to":
			updateBuilder = updateBuilder.Set("valid_to", util.LocalTime(sla.ValidTo))
		case "calendar_id":
			if sla.Calendar.Id != 0 {
				updateBuilder = updateBuilder.Set("calendar_id", sla.Calendar.Id)
			}
		case "reaction_time":
			if sla.ReactionTime != 0 {
				updateBuilder = updateBuilder.Set("reaction_time", sla.ReactionTime)
			}
		case "resolution_time":
			if sla.ResolutionTime != 0 {
				updateBuilder = updateBuilder.Set("resolution_time", sla.ResolutionTime)
			}
		}
	}

	// Generate the CTE for the update operation
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.sla.update.query_build_error", err)
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH s AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using buildSLASelectColumnsAndPlan
	selectBuilder, plan, err := buildSLASelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("s")

	return selectBuilder, plan, nil
}

func (s *SLAStore) Update(rpc options.UpdateOptions, update *cases.SLA) (*cases.SLA, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.update.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildUpdateSLAQuery(rpc, update)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.update.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.update.query_build_error", err)
	}
	// temporary object for scanning
	tempAdd := &cases.SLA{}
	scanArgs := convertToSLAScanArgs(plan, tempAdd)
	if err := d.QueryRow(rpc, query, args...).Scan(scanArgs...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.update.execution_error", err)
	}

	return tempAdd, nil
}

func (s *SLAStore) buildListSLAQuery(
	rpc options.SearchOptions,
) (sq.SelectBuilder, []SLAScan, error) {

	queryBuilder := sq.Select().
		From("cases.sla AS s").
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
	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, slaDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, plan, err := buildSLASelectColumnsAndPlan(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.sla.search.query_build_error", err)
	}

	return queryBuilder, plan, nil
}

func (s *SLAStore) List(rpc options.SearchOptions) (*cases.SLAList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildListSLAQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.query_build_error", err)
	}
	query = util2.CompactSQL(query)

	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.sla.list.execution_error", err)
	}
	defer rows.Close()

	var slas []*cases.SLA
	lCount := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && lCount >= int(rpc.GetSize()) {
			next = true
			break
		}

		sla := &cases.SLA{}
		scanArgs := convertToSLAScanArgs(plan, sla)

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.sla.list.row_scan_error", err)
		}

		slas = append(slas, sla)
		lCount++
	}

	return &cases.SLAList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: slas,
	}, nil
}

func (s *SLAStore) buildDeleteSLAQuery(
	rpc options.DeleteOptions,
) (sq.DeleteBuilder, error) {
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, dberr.NewDBInternalError("postgres.sla.delete.missing_ids", fmt.Errorf("no IDs provided for deletion"))
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.sla").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	return deleteBuilder, nil
}

func (s *SLAStore) Delete(rpc options.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := s.buildDeleteSLAQuery(rpc)
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.query_build_error", err)
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.query_to_sql_error", err)
	}

	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return dberr.NewDBInternalError("postgres.sla.delete.execution_error", execErr)
	}

	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.sla.delete.no_rows_affected")
	}

	return nil
}

func NewSLAStore(store *Store) (store.SLAStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_sla.check.bad_arguments",
			"error creating SLA interface, main store is nil")
	}
	return &SLAStore{storage: store}, nil
}
