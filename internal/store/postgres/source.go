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
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/util"
	"google.golang.org/grpc/codes"
)

const (
	sourceLeft        = "s"
	sourceDefaultSort = "name"
)

type Source struct {
	storage *Store
}

func buildSourceSelectColumnsAndPlan(base sq.SelectBuilder, fields []string) (sq.SelectBuilder, error) {
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(sourceLeft, "id"))
		case "name":
			base = base.Column(util2.Ident(sourceLeft, "name"))
		case "description":
			base = base.Column(util2.Ident(sourceLeft, "description"))
		case "type":
			base = base.Column(util2.Ident(sourceLeft, "type"))
		case "created_at":
			base = base.Column(util2.Ident(sourceLeft, "created_at"))
		case "updated_at":
			base = base.Column(util2.Ident(sourceLeft, "updated_at"))
		case "created_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.created_by) created_by",
				sourceLeft))
		case "updated_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by",
				sourceLeft))
		default:
			return base, errors.New(fmt.Sprintf("unknown field: %s", field), errors.WithCode(codes.InvalidArgument))
		}
	}
	return base, nil
}

func (s *Source) buildCreateSourceQuery(rpc options.Creator, source *model.Source) (sq.SelectBuilder, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	insertBuilder := sq.Insert("cases.source").
		Columns("name", "dc", "created_at", "description", "type", "created_by", "updated_at", "updated_by").
		Values(
			source.Name,
			rpc.GetAuthOpts().GetDomainId(),
			rpc.RequestTime(),
			source.Description,
			source.Type,
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, ParseError(err)
	}

	cte := sq.Expr("WITH s AS ("+insertSQL+")", args...)
	selectBuilder, err := buildSourceSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	return selectBuilder.PrefixExpr(cte).From(sourceLeft), nil
}

func (s *Source) Create(rpc options.Creator, source *model.Source) (*model.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.create.database_connection_error: %v", dbErr))
	}

	selectBuilder, err := s.buildCreateSourceQuery(rpc, source)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.create.query_build_error: %v", err))
	}

	item := model.Source{}
	err = pgxscan.Get(rpc, d, &item, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return &item, nil
}

func (s *Source) buildUpdateSourceQuery(rpc options.Updator, source *model.Source) (sq.SelectBuilder, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	updateBuilder := sq.Update("cases.source").
		PlaceholderFormat(sq.Dollar).
		Set("updated_at", rpc.RequestTime()).
		Set("updated_by", rpc.GetAuthOpts().GetUserId()).
		Where(sq.Eq{"id": source.Id}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()})

	for _, field := range rpc.GetMask() {
		switch field {
		case "name":
			if source.Name != nil && *source.Name != "" {
				updateBuilder = updateBuilder.Set("name", source.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", source.Description))
		case "type":
			if source.Type != nil && *source.Type != _go.SourceType_TYPE_UNSPECIFIED.String() {
				updateBuilder = updateBuilder.Set("type", source.Type)
			}
		}
	}

	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, ParseError(err)
	}

	cte := sq.Expr("WITH s AS ("+updateSQL+")", args...)
	selectBuilder, err := buildSourceSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	return selectBuilder.PrefixExpr(cte).From(sourceLeft), nil
}

func (s *Source) Update(rpc options.Updator, source *model.Source) (*model.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.update.database_connection_error: %v", dbErr))
	}

	selectBuilder, err := s.buildUpdateSourceQuery(rpc, source)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.update.query_build_error: %v", err))
	}

	temp := &model.Source{}
	err = pgxscan.Get(rpc, d, temp, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	return temp, nil
}

func (s *Source) buildListSourceQuery(rpc options.Searcher) (sq.SelectBuilder, error) {
	queryBuilder := sq.Select().
		From("cases.source AS s").
		Where(sq.Eq{"s.dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar)

	if len(rpc.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"s.id": rpc.GetIDs()})
	}

	if name, ok := rpc.GetFilter("name").(string); ok && name != "" {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "s.name")
	}

	if types, ok := rpc.GetFilter("type").([]_go.SourceType); ok && len(types) > 0 {
		var typeStrings []string
		for _, t := range types {
			typeStrings = append(typeStrings, t.String())
		}
		queryBuilder = queryBuilder.Where(sq.Eq{"s.type": typeStrings})
	}

	queryBuilder = util2.ApplyDefaultSorting(rpc, queryBuilder, sourceDefaultSort)
	queryBuilder = util2.ApplyPaging(rpc.GetPage(), rpc.GetSize(), queryBuilder)

	return buildSourceSelectColumnsAndPlan(queryBuilder, rpc.GetFields())
}

func (s *Source) List(rpc options.Searcher) ([]*model.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.list.database_connection_error: %v", dbErr))
	}

	selectBuilder, err := s.buildListSourceQuery(rpc)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.list.query_build_error: %v", err))
	}

	var sources []*model.Source
	err = pgxscan.Select(rpc, d, &sources, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}
	return sources, nil
}

func (s *Source) buildDeleteSourceQuery(rpc options.Deleter) (sq.DeleteBuilder, error) {
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, errors.InvalidArgument("no IDs provided for deletion")
	}

	return sq.Delete("cases.source").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar), nil
}

func (s *Source) Delete(rpc options.Deleter) (*model.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.delete.database_connection_error: %v", dbErr))
	}

	deleteBuilder, err := s.buildDeleteSourceQuery(rpc)
	if err != nil {
		return nil, err
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, errors.Internal(fmt.Sprintf("postgres.source.delete.query_build_error: %v", err))
	}

	res, err := d.Exec(rpc, query, args...)
	if err != nil {
		return nil, ParseError(err)
	}

	if res.RowsAffected() == 0 {
		return nil, errors.NotFound("postgres.source.delete.no_rows_affected")
	}

	return nil, nil
}

func NewSourceStore(store *Store) (store.SourceStore, error) {
	if store == nil {
		return nil, errors.New("error creating source interface, main store is nil")
	}
	return &Source{storage: store}, nil
}
