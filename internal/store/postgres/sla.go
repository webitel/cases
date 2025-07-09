package postgres

import (
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"

	errors "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
	storeutil "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/util"
)

const (
	slaLeft        = "s"
	slaDefaultSort = "name"
)

type SLAStore struct {
	storage *Store
}

// Helper function to dynamically build select columns.
func buildSLASelectColumns(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, error) {
	var (
		createdByAlias string
		joinCreatedBy  = func(alias string) string {
			if createdByAlias != "" {
				return createdByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.created_by = %s.id", alias, slaLeft, alias))
			createdByAlias = alias
			return alias
		}
		updatedByAlias string
		joinUpdatedBy  = func(alias string) string {
			if updatedByAlias != "" {
				return updatedByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.updated_by = %s.id", alias, slaLeft, alias))
			updatedByAlias = alias
			return alias
		}
	)
	base = base.Column(storeutil.Ident(slaLeft, "id"))
	for _, field := range fields {
		switch field {
		case "id":
			// already set
		case "name":
			base = base.Column(storeutil.Ident(slaLeft, "name"))
		case "description":
			base = base.Column(storeutil.Ident(slaLeft, "description"))
		case "valid_from":
			base = base.Column(storeutil.Ident(slaLeft, "valid_from"))
		case "valid_to":
			base = base.Column(storeutil.Ident(slaLeft, "valid_to"))
		case "calendar":
			base = base.
				LeftJoin("flow.calendar cal ON cal.id = s.calendar_id").
				Column("cal.id as calendar_id").
				Column("cal.name as calendar_name")
		case "reaction_time":
			base = base.Column(storeutil.Ident(slaLeft, "reaction_time"))
		case "resolution_time":
			base = base.Column(storeutil.Ident(slaLeft, "resolution_time"))
		case "created_at":
			base = base.Column(storeutil.Ident(slaLeft, "created_at"))
		case "updated_at":
			base = base.Column(storeutil.Ident(slaLeft, "updated_at"))
		case "created_by":
			alias := "slacb"
			joinCreatedBy(alias)
			base = base.Column(fmt.Sprintf("%s.id created_by_id", alias))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) created_by_name", alias, alias))
		case "updated_by":
			alias := "slaub"
			joinUpdatedBy(alias)
			base = base.Column(fmt.Sprintf("%s.id updated_by_id", alias))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) updated_by_name", alias, alias))
		default:
			return base, errors.New(fmt.Sprintf("unknown field: %s", field))
		}
	}
	return base, nil
}

func (s *SLAStore) buildCreateSLAQuery(rpc options.Creator, sla *model.SLA) (sq.SelectBuilder, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(fields)
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
			sla.Description,
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
			sla.ValidFrom,
			sla.ValidTo,
			sla.Calendar.Id,
			sla.ReactionTime,
			sla.ResolutionTime,
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *") // RETURNING all columns for use in the next SELECT

	// Convert the INSERT query into a CTE
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, errors.New("unable to convert to sql", errors.WithCause(err))
	}

	// Use the INSERT query as a CTE (Common Table Expression)
	cte := sq.Expr("WITH s AS ("+insertSQL+")", args...)

	// Dynamically build the SELECT query for the resulting row
	selectBuilder, err := buildSLASelectColumns(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(slaLeft)

	return selectBuilder, nil
}

func (s *SLAStore) Create(rpc options.Creator, input *model.SLA) (*model.SLA, error) {
	db, err := s.storage.Database()
	if err != nil {
		return nil, err
	}

	selectBuilder, err := s.buildCreateSLAQuery(rpc, input)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var res model.SLA
	err = pgxscan.Get(rpc, db, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return &res, nil
}

func (s *SLAStore) buildUpdateSLAQuery(
	rpc options.Updator,
	input *model.SLA,
) (sq.SelectBuilder, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(fields)
	// Start the UPDATE query
	updateBuilder := sq.Update("cases.sla").
		PlaceholderFormat(sq.Dollar). // Use PostgreSQL-compatible placeholders
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": input.Id}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()})

	// Dynamically add fields to the SET clause
	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			if input.Name != nil && *input.Name != "" {
				updateBuilder = updateBuilder.Set("name", input.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", input.Description))
		case "valid_from":
			updateBuilder = updateBuilder.Set("valid_from", input.ValidFrom)
		case "valid_to":
			updateBuilder = updateBuilder.Set("valid_to", input.ValidTo)
		case "calendar":
			if input.Calendar != nil && input.Calendar.Id != nil {
				updateBuilder = updateBuilder.Set("calendar_id", input.Calendar.Id)
			}
		case "reaction_time":
			updateBuilder = updateBuilder.
				Set("reaction_time", input.ReactionTime)
		case "resolution_time":
			updateBuilder = updateBuilder.
				Set("resolution_time", input.ResolutionTime)
		}
	}

	// Generate the CTE for the update operation
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, errors.New("unable to convert to sql", errors.WithCause(err))
	}

	// Use the UPDATE query as a CTE
	cte := sq.Expr("WITH s AS ("+updateSQL+")", args...)

	// Build select clause and scan plan dynamically using buildSLASelectColumnsAndPlan
	selectBuilder, err := buildSLASelectColumns(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	// Combine the CTE with the SELECT query
	selectBuilder = selectBuilder.PrefixExpr(cte).From(slaLeft)

	return selectBuilder, nil
}

func (s *SLAStore) Update(rpc options.Updator, input *model.SLA) (*model.SLA, error) {
	db, err := s.storage.Database()
	if err != nil {
		return nil, err
	}

	selectBuilder, err := s.buildUpdateSLAQuery(
		rpc,
		input,
	)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.New("unable to convert to sql", errors.WithCause(err))
	}

	var res model.SLA
	err = pgxscan.Get(rpc, db, &res, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return &res, nil
}

func (s *SLAStore) buildListSLAQuery(rpc options.Searcher) (sq.SelectBuilder, error) {
	queryBuilder := sq.Select().
		From("cases.sla AS s").
		Where(sq.Eq{"s.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"s.id": rpc.GetIDs()})
	}

	// Add name filter if provided
	nameFilters := rpc.GetFilter("name")
	if len(nameFilters) > 0 {
		f := nameFilters[0]
		if (f.Operator == "=" || f.Operator == "") && len(f.Value) > 0 {
			queryBuilder = storeutil.AddSearchTerm(queryBuilder, f.Value, "s.name")
		}
	}

	switch field, op := storeutil.GetSortingOperator(rpc.GetSort()); field {
	case "calendar":
		queryBuilder = queryBuilder.
			LeftJoin("flow.calendar cal_sort ON cal_sort.id = s.calendar_id")
		queryBuilder = queryBuilder.OrderBy(fmt.Sprintf("%s %s", "cal_sort.name", op))
	default:
		queryBuilder = storeutil.ApplyDefaultSorting(rpc, queryBuilder, slaDefaultSort)
	}
	queryBuilder = storeutil.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	queryBuilder, err := buildSLASelectColumns(queryBuilder, rpc.GetFields())
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	return queryBuilder, nil
}

func (s *SLAStore) List(rpc options.Searcher) ([]*model.SLA, error) {
	d, err := s.storage.Database()
	if err != nil {
		return nil, err
	}

	selectBuilder, err := s.buildListSLAQuery(rpc)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.New("unable to convert to sql", errors.WithCause(err))
	}
	query = storeutil.CompactSQL(query)

	var slas []*model.SLA
	err = pgxscan.Select(rpc, d, &slas, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	return slas, nil
}

func (s *SLAStore) buildDeleteSLAQuery(
	rpc options.Deleter,
) (sq.SelectBuilder, error) {
	fields := []string{"id", "name", "description", "created_at", "updated_at", "created_by", "updated_by", "calendar", "reaction_time", "resolution_time", "valid_from", "valid_to"}

	// Ensure IDs are provided
	if len(rpc.GetIDs()) == 0 {
		return sq.SelectBuilder{}, errors.InvalidArgument("no IDs provided for deletion")

	}

	// Build the delete query
	deleteBuilder := sq.Delete("cases.sla").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	deleteSQL, args, err := deleteBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, errors.New("unable to convert to sql", errors.WithCause(err))
	}

	cte := sq.Expr("WITH deleted AS ("+deleteSQL+")", args...)

	selectBuilder, err := buildSLASelectColumns(
		sq.Select().PrefixExpr(cte).From("deleted s"),
		fields,
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	selectBuilder = selectBuilder.PlaceholderFormat(sq.Dollar)

	return selectBuilder, nil
}

func (s *SLAStore) Delete(rpc options.Deleter) (*model.SLA, error) {
	d, err := s.storage.Database()
	if err != nil {
		return nil, err
	}

	deleteBuilder, err := s.buildDeleteSLAQuery(rpc)
	if err != nil {
		return nil, err
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, errors.New("unable to convert to sql", errors.WithCause(err))
	}

	var result model.SLA

	err = pgxscan.Get(rpc, d, &result, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	return &result, nil
}

func NewSLAStore(store *Store) (store.SLAStore, error) {
	if store == nil {
		return nil, errors.New(
			"error creating SLA interface, main store is nil")
	}
	return &SLAStore{storage: store}, nil
}
