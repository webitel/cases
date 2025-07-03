package postgres

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
	storeutil "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/util"
	"google.golang.org/grpc/codes"
)

type StatusScan func(status *_go.Status) any

const (
	statusLeft        = "s"
	statusDefaultSort = "name"
)

type Status struct {
	storage *Store
}

// Helper function to dynamically build select columns and plan.
func buildStatusSelectColumnsAndPlan(base sq.SelectBuilder, fields []string) (sq.SelectBuilder, error) {
	var (
		createdByAlias string
		updatedByAlias string
	)
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(storeutil.Ident(statusLeft, "id"))
		case "name":
			base = base.Column(storeutil.Ident(statusLeft, "name"))
		case "description":
			base = base.Column(storeutil.Ident(statusLeft, "description"))
		case "created_at":
			base = base.Column(storeutil.Ident(statusLeft, "created_at"))
		case "updated_at":
			base = base.Column(storeutil.Ident(statusLeft, "updated_at"))
		case "created_by":
			alias := "crb"
			if createdByAlias == "" {
				base = storeutil.SetUserColumn(base, statusLeft, alias, field)
			}
			createdByAlias = alias
		case "updated_by":
			alias := "upb"
			if updatedByAlias == "" {
				base = storeutil.SetUserColumn(base, statusLeft, alias, field)
			}
			updatedByAlias = alias
		default:
			return base, errors.New(fmt.Sprintf("unknown field: %s", field), errors.WithCode(codes.InvalidArgument))
		}
	}
	return base, nil
}

func (s *Status) buildCreateStatusQuery(rpc options.Creator, input *model.Status) (sq.SelectBuilder, error) {
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
		return sq.SelectBuilder{}, errors.Internal(fmt.Sprintf("postgres.input.create.query_build_error: %v", err))
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH s AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, err := buildStatusSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(statusLeft)

	return selectBuilder, nil
}

func (s *Status) Create(rpc options.Creator, input *model.Status) (*model.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.create.database_connection_error: %v", dbErr))
	}

	selectBuilder, err := s.buildCreateStatusQuery(rpc, input)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.create.query_build_error: %v", err))
	}
	res := model.Status{}
	err = pgxscan.Get(rpc, d, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return &res, nil
}

func (s *Status) buildUpdateStatusQuery(rpc options.Updator, input *model.Status) (sq.SelectBuilder, error) {
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
			updateBuilder = updateBuilder.Set("name", input.Name)
		case "description":
			updateBuilder = updateBuilder.Set("description", input.Description)
		}
	}

	// Generate the CTE for the update operation
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, errors.Internal(fmt.Sprintf("postgres.input.update.query_build_error: %v", err))
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH s AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using buildStatusSelectColumnsAndPlan
	selectBuilder, err := buildStatusSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From("s")

	return selectBuilder, nil
}

func (s *Status) Update(rpc options.Updator, input *model.Status) (*model.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.input.database_connection_error: %v", dbErr))
	}

	selectBuilder, err := s.buildUpdateStatusQuery(rpc, input)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.input.query_build_error: %v", err))
	}
	// temporary object for scanning
	res := model.Status{}
	err = pgxscan.Get(rpc, d, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return &res, nil
}

func (s *Status) buildListStatusQuery(rpc options.Searcher) (sq.SelectBuilder, error) {

	queryBuilder := sq.Select().
		From("cases.status AS s").
		Where(sq.Eq{"s.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	// Add ID filter if provided
	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"s.id": rpc.GetIDs()})
	}

	// Add name filter if provided
	nameFilters := rpc.GetFilter("name")
	if len(nameFilters) > 0 {
		f := nameFilters[0]
		if f.Operator == "=" || f.Operator == "" {
			queryBuilder = storeutil.AddSearchTerm(queryBuilder, f.Value, "s.name")
		}
	}

	// -------- Apply sorting ----------
	queryBuilder = storeutil.ApplyDefaultSorting(rpc, queryBuilder, statusDefaultSort)

	// ---------Apply paging based on Search Opts ( page ; size ) -----------------
	queryBuilder = storeutil.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	// Add select columns and scan plan for requested fields
	queryBuilder, err := buildStatusSelectColumnsAndPlan(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	return queryBuilder, nil
}

func (s *Status) List(rpc options.Searcher) ([]*model.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.list.database_connection_error: %v", dbErr))
	}

	selectBuilder, err := s.buildListStatusQuery(rpc)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.list.query_build_error: %v", err))
	}
	query = storeutil.CompactSQL(query)

	var statuses []*model.Status
	err = pgxscan.Select(rpc, d, &statuses, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	return statuses, nil
}

func (s *Status) buildDeleteStatusQuery(
	rpc options.Deleter,
) (sq.DeleteBuilder, error) {
	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, errors.InvalidArgument("no IDs provided for deletion")
	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.status").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	return deleteBuilder, nil
}

func (s *Status) Delete(rpc options.Deleter) (*model.Status, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.delete.database_connection_error: %v", dbErr))
	}

	deleteBuilder, err := s.buildDeleteStatusQuery(rpc)
	if err != nil {
		return nil, err
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.status.delete.query_to_sql_error: %v", err))
	}

	res, execErr := d.Exec(rpc, query, args...)
	if execErr != nil {
		return nil, ParseError(execErr)
	}

	if res.RowsAffected() == 0 {
		return nil, errors.NotFound("postgres.status.delete.no_rows_affected")
	}

	return nil, nil
}

func NewStatusStore(store *Store) (store.StatusStore, error) {
	if store == nil {
		return nil, errors.New("error creating status interface, main store is nil")
	}
	return &Status{storage: store}, nil
}
