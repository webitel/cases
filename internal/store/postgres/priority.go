package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
	util2 "github.com/webitel/cases/internal/store/util"
	util "github.com/webitel/cases/util"
)

type Priority struct {
	storage *Store
}

const (
	prioLeft            = "cp"
	priorityDefaultSort = "name"
)

func buildPrioritySelectColumns(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, error) {

	var (
		createdByAlias string
		joinCreatedBy  = func(alias string) string {
			if createdByAlias != "" {
				return createdByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.created_by = %s.id", alias, prioLeft, alias))
			createdByAlias = alias
			return alias
		}
		updatedByAlias string
		joinUpdatedBy  = func(alias string) string {
			if updatedByAlias != "" {
				return updatedByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.updated_by = %s.id", alias, prioLeft, alias))
			updatedByAlias = alias
			return alias
		}
	)
	base = base.Column(util2.Ident(prioLeft, "id"))
	for _, field := range fields {
		switch field {
		case "id":
			// already set
		case "name":
			base = base.Column(util2.Ident(prioLeft, "name"))
		case "description":
			base = base.Column(util2.Ident(prioLeft, "description"))
		case "created_at":
			base = base.Column(util2.Ident(prioLeft, "created_at"))
		case "updated_at":
			base = base.Column(util2.Ident(prioLeft, "updated_at"))
		case "created_by":
			alias := "prcb"
			joinCreatedBy(alias)
			base = base.Column(fmt.Sprintf("%s.id created_by_id", alias))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) created_by_name", alias, alias))
		case "updated_by":
			alias := "prub"
			joinUpdatedBy(alias)
			base = base.Column(fmt.Sprintf("%s.id updated_by_id", alias))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) updated_by_name", alias, alias))
		case "color":
			base = base.Column(util2.Ident(prioLeft, "color"))
		default:
			return base, dberr.NewDBInternalError("postgres.priority.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return base, nil
}

// Create implements store.PriorityStore.
func (p *Priority) Create(rpc options.Creator, add *model.Priority) (*model.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.database_connection_error", dbErr)
	}

	selectBuilder, err := p.buildCreatePriorityQuery(rpc, add)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.query_build_error", err)
	}
	// temporary object for scanning
	var res model.Priority
	err = pgxscan.Get(rpc, d, &res, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.create.execution_error", err)
	}

	return &res, nil
}

func (p *Priority) buildCreatePriorityQuery(
	rpc options.Creator,
	priority *model.Priority,
) (sq.SelectBuilder, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(fields)
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
		return sq.SelectBuilder{}, dberr.NewDBInternalError("postgres.priority.create.query_build_error", err)
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH cp AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, err := buildPrioritySelectColumns(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(prioLeft)

	return selectBuilder, nil
}

func (p *Priority) Delete(rpc options.Deleter) (*model.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.delete.database_connection_error", dbErr)
	}

	selectBuilder, err := p.buildDeletePriorityQuery(rpc)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.delete.query_build_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.delete.query_to_sql_error", err)
	}

	var result model.Priority
	err = pgxscan.Get(rpc, d, &result, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.delete.execution_error", err)
	}

	return &result, nil
}

func (p *Priority) buildDeletePriorityQuery(
	rpc options.Deleter,
) (sq.SelectBuilder, error) {
	fields := []string{"id", "name", "description", "created_at", "updated_at", "created_by", "updated_by", "color"}
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.SelectBuilder{}, dberr.NewDBInternalError("postgres.priority.delete.missing_ids", fmt.Errorf("no IDs provided for deletion"))
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.priority").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *") // RETURNING all columns for use in the next SELECT

	deleteSQL, args, err := deleteBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, dberr.NewDBInternalError("postgres.priority.delete.query_to_sql_error", err)
	}

	cte := sq.Expr("WITH deleted AS ("+deleteSQL+")", args...)

	selectBuilder, err := buildPrioritySelectColumns(
		sq.Select().PrefixExpr(cte).From("deleted cp"),
		fields,
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	selectBuilder = selectBuilder.PlaceholderFormat(sq.Dollar)

	return selectBuilder, nil
}

// List implements store.PriorityStore.
func (p *Priority) List(
	rpc options.Searcher,
	notInSla int64,
	inSla int64,
) ([]*model.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.database_connection_error", dbErr)
	}

	selectBuilder, err := p.buildListPriorityQuery(rpc, notInSla, inSla)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.query_build_error", err)
	}
	query = util2.CompactSQL(query)

	var priorities []*model.Priority

	err = pgxscan.Select(rpc, d, &priorities, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.list.execution_error", err)
	}
	return priorities, nil
}

func (p *Priority) buildListPriorityQuery(
	rpc options.Searcher,
	notInSla int64,
	inSla int64,
) (sq.SelectBuilder, error) {

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

	// -------- Apply sorting ----------
	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, priorityDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, err := buildPrioritySelectColumns(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, dberr.NewDBInternalError("postgres.priority.search.query_build_error", err)
	}

	return queryBuilder, nil
}

// Update implements store.PriorityStore.
func (p *Priority) Update(rpc options.Updator, update *model.Priority) (*model.Priority, error) {
	d, dbErr := p.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.database_connection_error", dbErr)
	}

	selectBuilder, err := p.buildUpdatePriorityQuery(rpc, update)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.query_build_error", err)
	}
	// temporary object for scanning
	var res model.Priority
	err = pgxscan.Get(rpc, d, &res, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.priority.update.execution_error", err)
	}

	return &res, nil
}

func (p *Priority) buildUpdatePriorityQuery(
	rpc options.Updator,
	priority *model.Priority,
) (sq.SelectBuilder, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(fields)
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
		return sq.SelectBuilder{}, dberr.NewDBInternalError("postgres.priority.update.query_build_error", err)
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH cp AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using `buildPrioritySelectColumnsAndPlan`
	selectBuilder, err := buildPrioritySelectColumns(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(prioLeft)

	return selectBuilder, nil
}

func NewPriorityStore(store *Store) (store.PriorityStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_priority.check.bad_arguments",
			"error creating priority interface to the status_condition table, main store is nil")
	}
	return &Priority{storage: store}, nil
}
