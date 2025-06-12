package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/model/options"
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
func buildCloseReasonSelectColumns(fields []string) ([]string, error) {
	var cols []string
	for _, field := range fields {
		switch field {
		case "id":
			cols = append(cols, "cr.id")
		case "name":
			cols = append(cols, "cr.name")
		case "description":
			cols = append(cols, "cr.description")
		case "created_at":
			cols = append(cols, "cr.created_at")
		case "updated_at":
			cols = append(cols, "cr.updated_at")
		case "close_reason_id":
			cols = append(cols, "cr.close_reason_id")
		case "dc":
			cols = append(cols, "cr.dc")
		case "created_by":
			cols = append(cols,
				"cb.id as created_by_id",
				"COALESCE(cb.name, cb.username) as created_by_name",
			)
		case "updated_by":
			cols = append(cols,
				"ub.id as updated_by_id",
				"COALESCE(ub.name, ub.username) as updated_by_name",
			)
		default:
			return nil, dberr.NewDBInternalError("postgres.close_reason.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return cols, nil
}

/*
	 func (s *CloseReason) buildCreateCloseReasonQuery(
		rpc options.Creator,
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
		plan, err := buildCloseReasonSelectColumns(fields)
		if err != nil {
			return sq.SelectBuilder{}, nil, err
		}

		// Combine the CTE with the SELECT query
		selectBuilder = selectBuilder.PrefixExpr(cte).From(crLeft)

		return selectBuilder, plan, nil
	}
*/

func (s *CloseReason) Create(creator options.Creator, input *model.CloseReason) (*model.CloseReason, error) {
    d, dbErr := s.storage.Database()
    if dbErr != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.create.database_connection_error", dbErr)
    }

    // Dynamic fields
    fields := creator.GetFields()
    if len(fields) == 0 {
        fields = []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}
    }
    cols, err := buildCloseReasonSelectColumns(fields)
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.create.columns_error", err)
    }

    // Insert builder
    insertBuilder := sq.Insert("cases.close_reason").
        Columns("name", "description", "close_reason_id", "created_at", "created_by", "updated_at", "updated_by", "dc").
        Values(
            input.Name,
            input.Description,
            input.CloseReasonGroupId,
            creator.RequestTime(),
            creator.GetAuthOpts().GetUserId(),
            creator.RequestTime(),
            creator.GetAuthOpts().GetUserId(),
            creator.GetAuthOpts().GetDomainId(),
        ).
        PlaceholderFormat(sq.Dollar).
        Suffix("RETURNING *")

    insertSQL, args, err := insertBuilder.ToSql()
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.create.query_build_error", err)
    }

    // Select with joins, using the CTE
    queryBuilder := sq.Select(cols...).
        From("cr").
        LeftJoin("directory.wbt_user cb ON cb.id = cr.created_by").
        LeftJoin("directory.wbt_user ub ON ub.id = cr.updated_by").
        PlaceholderFormat(sq.Dollar)

    selectSQL, selectArgs, err := queryBuilder.ToSql()
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.create.select_query_build_error", err)
    }

    finalQuery := fmt.Sprintf("WITH cr AS (%s) %s", insertSQL, selectSQL)
    allArgs := append(args, selectArgs...)

    var result model.CloseReason
    err = pgxscan.Get(creator, d, &result, finalQuery, allArgs...)
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.create.execution_error", err)
    }

    return &result, nil
}

/* func (s *CloseReason) buildUpdateCloseReasonQuery(
	rpc options.Updator,
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
	plan, err := buildCloseReasonSelectColumns(fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("cr")

	return selectBuilder, plan, nil
} */

func (s *CloseReason) Update(updator options.Updator, input *model.CloseReason) (*model.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.database_connection_error", dbErr)
	}

	updateBuilder := sq.Update("cases.close_reason").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", updator.RequestTime()).
		Set("updated_by", updator.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": input.Id}).
		Where(sq.Eq{"dc": updator.GetAuthOpts().GetDomainId()})

	for _, field := range updator.GetMask() {
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

	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.query_build_error", err)
	}

	fields := updator.GetFields()
	if len(fields) == 0 {
		fields = []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}
	}
	cols, err := buildCloseReasonSelectColumns(fields)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.columns_error", err)
	}

	queryBuilder := sq.Select(cols...).
		From("updated cr").
		LeftJoin("directory.wbt_user cb ON cb.id = cr.created_by").
		LeftJoin("directory.wbt_user ub ON ub.id = cr.updated_by").
		PlaceholderFormat(sq.Dollar)

	query, selectArgs, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.select_query_build_error", err)
	}

	finalQuery := fmt.Sprintf("WITH updated AS (%s) %s", updateSQL, query)

	var result model.CloseReason
	// Combine update args and select args
	allArgs := append(args, selectArgs...)

	err = pgxscan.Get(updator, d, &result, finalQuery, allArgs...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.execution_error", err)
	}

	return &result, nil
}

/* func (s *CloseReason) buildListCloseReasonQuery(
	rpc options.Searcher,
	closeReasonId int64,
) (sq.SelectBuilder, []string, error) {
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
	plan, err := buildCloseReasonSelectColumns(rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.close_reason.search.query_build_error", err)
	}

	return queryBuilder, plan, nil
} */

func (s *CloseReason) List(searcher options.Searcher, closeReasonId int64) (*model.CloseReasonList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.database_connection_error", dbErr)
	}

	// Build dynamic columns
	fields := searcher.GetFields()
	if len(fields) == 0 {
		fields = []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}
	}
	cols, err := buildCloseReasonSelectColumns(fields)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.columns_error", err)
	}

	queryBuilder := sq.Select(cols...).
		From("cases.close_reason AS cr").
		LeftJoin("directory.wbt_user cb ON cb.id = cr.created_by").
		LeftJoin("directory.wbt_user ub ON ub.id = cr.updated_by").
		Where(sq.Eq{"cr.dc": searcher.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Filters
	if len(searcher.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cr.id": searcher.GetIDs()})
	}
	if name, ok := searcher.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "cr.name")
	}
	if closeReasonId != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cr.close_reason_id": closeReasonId})
	}

	// Sorting and paging
	queryBuilder = util2.ApplyDefaultSorting(searcher, queryBuilder, closeReasonDefaultSort)
	queryBuilder = util2.ApplyPaging(searcher.GetPage(), searcher.GetSize(), queryBuilder)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.query_build_error", err)
	}

	var items []*model.CloseReason
	if err := pgxscan.Select(searcher, d, &items, query, args...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.execution_error", err)
	}

	// Paging helper
	items, next := util2.ResolvePaging(searcher.GetSize(), items)

	return &model.CloseReasonList{
		Page:  searcher.GetPage(),
		Next:  next,
		Items: items,
	}, nil
}

/* func (s *CloseReason) buildDeleteCloseReasonQuery(
	rpc options.Deleter,
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
} */

func (s *CloseReason) Delete(deleter options.Deleter) (*model.CloseReason, error) {
    d, dbErr := s.storage.Database()
    if dbErr != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.delete.database_connection_error", dbErr)
    }

    // default fields because GetFields is not available in the deleter
    fields := []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}

    cols, err := buildCloseReasonSelectColumns(fields)
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.delete.columns_error", err)
    }

    deleteBuilder := sq.Delete("cases.close_reason").
        Where(sq.Eq{"id": deleter.GetIDs()}).
        Where(sq.Eq{"dc": deleter.GetAuthOpts().GetDomainId()}).
        PlaceholderFormat(sq.Dollar).
        Suffix("RETURNING *")

    deleteSQL, args, err := deleteBuilder.ToSql()
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.delete.query_to_sql_error", err)
    }

    queryBuilder := sq.Select(cols...).
        From("deleted cr").
        LeftJoin("directory.wbt_user cb ON cb.id = cr.created_by").
        LeftJoin("directory.wbt_user ub ON ub.id = cr.updated_by").
        PlaceholderFormat(sq.Dollar)

    selectSQL, selectArgs, err := queryBuilder.ToSql()
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.delete.select_query_build_error", err)
    }

    finalQuery := fmt.Sprintf("WITH deleted AS (%s) %s", deleteSQL, selectSQL)
    allArgs := append(args, selectArgs...)

    var result model.CloseReason
    err = pgxscan.Get(deleter, d, &result, finalQuery, allArgs...)
    if err != nil {
        return nil, dberr.NewDBInternalError("postgres.close_reason.delete.execution_error", err)
    }

    return &result, nil
}

func NewCloseReasonStore(store *Store) (store.CloseReasonStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.new_close_reason.check.bad_arguments",
			"error creating close_reason interface, main store is nil")
	}
	return &CloseReason{storage: store}, nil
}
