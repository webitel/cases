package postgres

import (
	"fmt"
	util2 "github.com/webitel/cases/internal/store/util"
	"github.com/webitel/cases/model/options"
	"strings"

	sq "github.com/Masterminds/squirrel"
	_go "github.com/webitel/cases/api/cases"
	dberr "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	"github.com/webitel/cases/util"
)

const (
	sourceLeft        = "s"
	sourceDefaultSort = "name"
)

type SourceScan func(source *_go.Source) any

type Source struct {
	storage *Store
}

func convertToSourceScanArgs(plan []SourceScan, source *_go.Source) []any {
	scanArgs := make([]any, 0, len(plan))
	for _, scan := range plan {
		scanArgs = append(scanArgs, scan(source))
	}
	return scanArgs
}

func buildSourceSelectColumnsAndPlan(
	base sq.SelectBuilder,
	fields []string,
) (sq.SelectBuilder, []SourceScan, error) {
	var plan []SourceScan
	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(util2.Ident(sourceLeft, "id"))
			plan = append(plan, func(s *_go.Source) any { return scanner.ScanInt64(&s.Id) })
		case "name":
			base = base.Column(util2.Ident(sourceLeft, "name"))
			plan = append(plan, func(s *_go.Source) any { return scanner.ScanText(&s.Name) })
		case "description":
			base = base.Column(util2.Ident(sourceLeft, "description"))
			plan = append(plan, func(s *_go.Source) any { return scanner.ScanText(&s.Description) })
		case "type":
			base = base.Column(util2.Ident(sourceLeft, "type"))
			plan = append(plan, func(s *_go.Source) any {
				return &scanner.SourceTypeScanner{SourceType: &s.Type}
			})
		case "created_at":
			base = base.Column(util2.Ident(sourceLeft, "created_at"))
			plan = append(plan, func(s *_go.Source) any { return scanner.ScanTimestamp(&s.CreatedAt) })
		case "updated_at":
			base = base.Column(util2.Ident(sourceLeft, "updated_at"))
			plan = append(plan, func(s *_go.Source) any { return scanner.ScanTimestamp(&s.UpdatedAt) })
		case "created_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.created_by) created_by",
				sourceLeft))
			plan = append(plan, func(s *_go.Source) any { return scanner.ScanRowLookup(&s.CreatedBy) })
		case "updated_by":
			base = base.Column(fmt.Sprintf(
				"(SELECT ROW(id, COALESCE(name, username))::text FROM directory.wbt_user WHERE id = %s.updated_by) updated_by",
				sourceLeft))
			plan = append(plan, func(s *_go.Source) any { return scanner.ScanRowLookup(&s.UpdatedBy) })
		default:
			return base, nil, dberr.NewDBInternalError("postgres.source.unknown_field", fmt.Errorf("unknown field: %s", field))
		}
	}
	return base, plan, nil
}

// StringToType converts a string into the corresponding Type enum value.
//
// Types are specified ONLY for Source dictionary and are ENUMS in API.
func stringToType(typeStr string) (_go.SourceType, error) {
	switch strings.ToUpper(typeStr) {
	case "CALL":
		return _go.SourceType_CALL, nil
	case "CHAT":
		return _go.SourceType_CHAT, nil
	case "SOCIAL_MEDIA":
		return _go.SourceType_SOCIAL_MEDIA, nil
	case "EMAIL":
		return _go.SourceType_EMAIL, nil
	case "API":
		return _go.SourceType_API, nil
	case "MANUAL":
		return _go.SourceType_MANUAL, nil
	default:
		return _go.SourceType_TYPE_UNSPECIFIED, fmt.Errorf("invalid type value: %s", typeStr)
	}
}

func (s *Source) buildCreateSourceQuery(rpc options.CreateOptions, source *_go.Source) (sq.SelectBuilder, []SourceScan, error) {
	fields := rpc.GetFields()
	fields = util.EnsureIdField(rpc.GetFields())
	insertBuilder := sq.Insert("cases.source").
		Columns("name", "dc", "created_at", "description", "type", "created_by", "updated_at", "updated_by").
		Values(
			source.Name,
			rpc.GetAuthOpts().GetDomainId(),
			rpc.RequestTime(),
			sq.Expr("NULLIF(?, '')", source.Description),
			source.Type.String(),
			rpc.GetAuthOpts().GetUserId(),
			rpc.RequestTime(),
			rpc.GetAuthOpts().GetUserId(),
		).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.source.create.query_build_error", err)
	}

	cte := sq.Expr("WITH s AS ("+insertSQL+")", args...)
	selectBuilder, plan, err := buildSourceSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	return selectBuilder.PrefixExpr(cte).From(sourceLeft), plan, nil
}

func (s *Source) Create(rpc options.CreateOptions, source *_go.Source) (*_go.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.source.create.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildCreateSourceQuery(rpc, source)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.source.create.query_build_error", err)
	}

	temp := &_go.Source{}
	if err := d.QueryRow(rpc, query, args...).Scan(convertToSourceScanArgs(plan, temp)...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.source.create.execution_error", err)
	}

	return temp, nil
}

func (s *Source) buildUpdateSourceQuery(rpc options.UpdateOptions, source *_go.Source) (sq.SelectBuilder, []SourceScan, error) {
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
			if source.Name != "" {
				updateBuilder = updateBuilder.Set("name", source.Name)
			}
		case "description":
			updateBuilder = updateBuilder.Set("description", sq.Expr("NULLIF(?, '')", source.Description))
		case "type":
			if source.Type != _go.SourceType_TYPE_UNSPECIFIED {
				updateBuilder = updateBuilder.Set("type", source.Type.String())
			}
		}
	}

	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, nil, dberr.NewDBInternalError("postgres.source.update.query_build_error", err)
	}

	cte := sq.Expr("WITH s AS ("+updateSQL+")", args...)
	selectBuilder, plan, err := buildSourceSelectColumnsAndPlan(sq.Select(), fields)
	if err != nil {
		return sq.SelectBuilder{}, nil, err
	}

	return selectBuilder.PrefixExpr(cte).From(sourceLeft), plan, nil
}

func (s *Source) Update(rpc options.UpdateOptions, source *_go.Source) (*_go.Source, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.source.update.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildUpdateSourceQuery(rpc, source)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.source.update.query_build_error", err)
	}

	temp := &_go.Source{}
	if err := d.QueryRow(rpc, query, args...).Scan(convertToSourceScanArgs(plan, temp)...); err != nil {
		return nil, dberr.NewDBInternalError("postgres.source.update.execution_error", err)
	}

	return temp, nil
}

func (s *Source) buildListSourceQuery(rpc options.SearchOptions) (sq.SelectBuilder, []SourceScan, error) {
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

func (s *Source) List(rpc options.SearchOptions) (*_go.SourceList, error) {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return nil, dberr.NewDBInternalError("postgres.source.list.database_connection_error", dbErr)
	}

	selectBuilder, plan, err := s.buildListSourceQuery(rpc)
	if err != nil {
		return nil, err
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.source.list.query_build_error", err)
	}

	rows, err := d.Query(rpc, query, args...)
	if err != nil {
		return nil, dberr.NewDBInternalError("postgres.source.list.execution_error", err)
	}
	defer rows.Close()

	var sources []*_go.Source
	count := 0
	next := false
	fetchAll := rpc.GetSize() == -1

	for rows.Next() {
		if !fetchAll && count >= int(rpc.GetSize()) {
			next = true
			break
		}

		src := &_go.Source{}
		if err := rows.Scan(convertToSourceScanArgs(plan, src)...); err != nil {
			return nil, dberr.NewDBInternalError("postgres.source.list.row_scan_error", err)
		}

		sources = append(sources, src)
		count++
	}

	return &_go.SourceList{
		Page:  int32(rpc.GetPage()),
		Next:  next,
		Items: sources,
	}, nil
}

func (s *Source) buildDeleteSourceQuery(rpc options.DeleteOptions) (sq.DeleteBuilder, error) {
	if len(rpc.GetIDs()) == 0 {
		return sq.DeleteBuilder{}, dberr.NewDBInternalError("postgres.source.delete.missing_ids", fmt.Errorf("no IDs provided"))
	}

	return sq.Delete("cases.source").
		Where(sq.Eq{"id": rpc.GetIDs()}).
		Where(sq.Eq{"dc": rpc.GetAuthOpts().GetDomainId()}).
		PlaceholderFormat(sq.Dollar), nil
}

func (s *Source) Delete(rpc options.DeleteOptions) error {
	d, dbErr := s.storage.Database()
	if dbErr != nil {
		return dberr.NewDBInternalError("postgres.source.delete.database_connection_error", dbErr)
	}

	deleteBuilder, err := s.buildDeleteSourceQuery(rpc)
	if err != nil {
		return err
	}

	query, args, err := deleteBuilder.ToSql()
	if err != nil {
		return dberr.NewDBInternalError("postgres.source.delete.query_build_error", err)
	}

	res, err := d.Exec(rpc, query, args...)
	if err != nil {
		return dberr.NewDBInternalError("postgres.source.delete.execution_error", err)
	}

	if res.RowsAffected() == 0 {
		return dberr.NewDBNoRowsError("postgres.source.delete.no_rows_affected")
	}

	return nil
}

func NewSourceStore(store *Store) (store.SourceStore, error) {
	if store == nil {
		return nil, dberr.NewDBError("postgres.source.store.nil", "nil store instance")
	}
	return &Source{storage: store}, nil
}
