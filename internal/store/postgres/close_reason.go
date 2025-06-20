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
	"github.com/webitel/cases/util"
)

type CloseReason struct {
	storage *Store
}

const (
	crLeft                 = "cr"
	closeReasonDefaultSort = "name"
)

func buildCloseReasonSelectColumns(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, error) {
	const crLeft = "cr"
	var (
		createdByAlias string
		joinCreatedBy  = func(alias string) string {
			if createdByAlias != "" {
				return createdByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.created_by = %s.id", alias, crLeft, alias))
			createdByAlias = alias
			return alias
		}
		updatedByAlias string
		joinUpdatedBy  = func(alias string) string {
			if updatedByAlias != "" {
				return updatedByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.updated_by = %s.id", alias, crLeft, alias))
			updatedByAlias = alias
			return alias
		}
	)
	base = base.Column(fmt.Sprintf("%s.id", crLeft))
	for _, field := range fields {
		switch field {
		case "id":
			// already set
		case "name":
			base = base.Column(fmt.Sprintf("%s.name", crLeft))
		case "description":
			base = base.Column(fmt.Sprintf("%s.description", crLeft))
		case "created_at":
			base = base.Column(fmt.Sprintf("%s.created_at", crLeft))
		case "updated_at":
			base = base.Column(fmt.Sprintf("%s.updated_at", crLeft))
		case "close_reason_id":
			base = base.Column(fmt.Sprintf("%s.close_reason_id", crLeft))
		case "dc":
			base = base.Column(fmt.Sprintf("%s.dc", crLeft))
		case "created_by":
			crb := "crb"
			joinCreatedBy(crb)
			base = base.Column(fmt.Sprintf("%s.id created_by_id", crb))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) created_by_name", crb, crb))
		case "updated_by":
			upb := "upb"
			joinUpdatedBy(upb)
			base = base.Column(fmt.Sprintf("%s.id updated_by_id", upb))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) updated_by_name", upb, upb))
		default:
			return base, dberr.NewDBInternalError(
				"postgres.close_reason.unknown_field",
				fmt.Errorf("unknown field: %s", field),
			)
		}
	}
	return base, nil
}

func (s *CloseReason) buildCreateCloseReasonQuery(
	creator options.Creator,
	input *model.CloseReason,
) (sq.SelectBuilder, []interface{}, error) {
	fields := creator.GetFields()
	fields = util.EnsureIdField(fields)
	if len(fields) == 0 {
		fields = []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}
	}

	insertBuilder := sq.Insert("cases.close_reason").
		Columns("name", "description", "close_reason_id", "created_at", "created_by", "updated_at", "updated_by", "dc").
		Values(
			input.Name,
			sq.Expr("NULLIF(?, '')", input.Description), // NULLIF for empty description
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
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.close_reason.create.query_build_error", err)
	}

	cte := sq.Expr("WITH cr AS ("+insertSQL+")", args...)

	selectBuilder, err := buildCloseReasonSelectColumns(
		sq.Select().PrefixExpr(cte).From("cr"),
		fields,
	)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}
	selectBuilder = selectBuilder.PlaceholderFormat(sq.Dollar)

	return selectBuilder, nil, nil
}

func (s *CloseReason) buildUpdateCloseReasonQuery(
	updator options.Updator,
	input *model.CloseReason,
) (sq.SelectBuilder, []interface{}, error) {
	fields := updator.GetFields()
	fields = util.EnsureIdField(fields)
	if len(fields) == 0 {
		fields = []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}
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
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.close_reason.update.query_build_error", err)
	}

	cte := sq.Expr("WITH updated AS ("+updateSQL+")", args...)

	selectBuilder, err := buildCloseReasonSelectColumns(
		sq.Select().PrefixExpr(cte).From("updated cr"),
		fields,
	)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}
	selectBuilder = selectBuilder.PlaceholderFormat(sq.Dollar)

	return selectBuilder, nil, nil
}

func (s *CloseReason) buildListCloseReasonQuery(
	searcher options.Searcher,
	closeReasonId int64,
) (sq.SelectBuilder, error) {
	fields := searcher.GetFields()
	if len(fields) == 0 {
		fields = []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}
	}

	queryBuilder, err := buildCloseReasonSelectColumns(
		sq.Select().From("cases.close_reason AS cr"),
		fields,
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	queryBuilder = queryBuilder.Where(sq.Eq{"cr.dc": searcher.GetAuthOpts().GetDomainId()})

	if len(searcher.GetIDs()) > 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cr.id": searcher.GetIDs()})
	}
	if name, ok := searcher.GetFilter("name").(string); ok && len(name) > 0 {
		queryBuilder = util2.AddSearchTerm(queryBuilder, name, "cr.name")
	}
	if closeReasonId != 0 {
		queryBuilder = queryBuilder.Where(sq.Eq{"cr.close_reason_id": closeReasonId})
	}

	queryBuilder = util2.ApplyDefaultSorting(searcher, queryBuilder, closeReasonDefaultSort)
	queryBuilder = util2.ApplyPaging(searcher.GetPage(), searcher.GetSize(), queryBuilder)
	queryBuilder = queryBuilder.PlaceholderFormat(sq.Dollar)

	return queryBuilder, nil
}

func (s *CloseReason) buildDeleteCloseReasonQuery(
	deleter options.Deleter) (sq.SelectBuilder, error) {
	if len(deleter.GetIDs()) == 0 {
		return sq.SelectBuilder{}, dberr.NewDBInternalError("postgres.close_reason.delete.missing_ids", fmt.Errorf("no IDs provided for deletion"))
	}
	fields := []string{"id", "name", "description", "close_reason_id", "created_at", "updated_at", "dc", "created_by", "updated_by"}
	deleteBuilder := sq.Delete("cases.close_reason").
		Where(sq.Eq{"id": deleter.GetIDs()}).
		Where(sq.Eq{"dc": deleter.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	deleteSQL, args, err := deleteBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, dberr.NewDBInternalError("postgres.close_reason.delete.query_to_sql_error", err)
	}

	cte := sq.Expr("WITH deleted AS ("+deleteSQL+")", args...)

	selectBuilder, err := buildCloseReasonSelectColumns(
		sq.Select().PrefixExpr(cte).From("deleted cr"),
		fields,
	)
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	selectBuilder = selectBuilder.PlaceholderFormat(sq.Dollar)

	return selectBuilder, nil
}

// --- CRUD Methods ---

func (s *CloseReason) Create(creator options.Creator, input *model.CloseReason) (*model.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.database_connection_error", dbErr)
	}

	selectBuilder, _, err := s.buildCreateCloseReasonQuery(creator, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.query_build_error", err)
	}

	var result model.CloseReason
	err = pgxscan.Get(creator, d, &result, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.create.execution_error", err)
	}

	return &result, nil
}

func (s *CloseReason) Update(updator options.Updator, input *model.CloseReason) (*model.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.database_connection_error", dbErr)
	}

	selectBuilder, _, err := s.buildUpdateCloseReasonQuery(updator, input)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.build_query_error", err)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.query_build_error", err)
	}

	var result model.CloseReason
	err = pgxscan.Get(updator, d, &result, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.update.execution_error", err)
	}

	return &result, nil
}

func (s *CloseReason) List(searcher options.Searcher, closeReasonId int64) ([]*model.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.database_connection_error", dbErr)
	}

	queryBuilder, err := s.buildListCloseReasonQuery(searcher, closeReasonId)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.build_query_error", err)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.query_build_error", err)
	}

	var items []*model.CloseReason
	if err := pgxscan.Select(searcher, d, &items, query, args...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.list.execution_error", err)
	}
	return items, nil
}

func (s *CloseReason) Delete(deleter options.Deleter) (*model.CloseReason, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := s.buildDeleteCloseReasonQuery(deleter)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.delete.build_query_error", err)
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.close_reason.delete.query_to_sql_error", err)
	}

	var result model.CloseReason
	err = pgxscan.Get(deleter, d, &result, query, args...)
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
