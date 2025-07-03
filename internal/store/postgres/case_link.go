package postgres

import (
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/webitel/cases/util"
	"net/url"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/model/options/defaults"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/internal/store/postgres/scanner"
	storeUtil "github.com/webitel/cases/internal/store/util"
)

const (
	caseLinkLeft           = "cl"
	caseLinkCreatedByAlias = "cb"
	caseLinkUpdatedByAlias = "ub"
	caseLinkAuthorAlias    = "au"
	linkDefaultSort        = "created_at"
)

type CaseLinkStore struct {
	storage   *Store
	mainTable string
}

var CaseLinkFields = []string{
	"created_by", "created_at", "updated_by", "updated_at", "id", "ver", "author", "name", "url",
}

// Create implements store.CaseLinkStore.
func (l *CaseLinkStore) Create(rpc options.Creator, add *model.CaseLink) (*model.CaseLink, error) {
	if rpc == nil {
		return nil, errors.InvalidArgument("create options required")
	}
	if err := ValidateLinkCreate(rpc.GetParentID(), add); err != nil {
		return nil, ParseError(err)
	}
	fields := rpc.GetFields()
	if len(fields) == 0 {
		fields = []string{"id", "ver", "created_by", "created_at", "updated_by", "updated_at", "author", "name", "url"}
	}
	selectBuilder, err := buildCreateCaseLinkQuery(rpc, add, fields)
	if err != nil {
		return nil, ParseError(err)
	}
	db, err := l.storage.Database()
	if err != nil {
		return nil, err
	}
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	var result model.CaseLink
	if err := pgxscan.Get(rpc, db, &result, query, args...); err != nil {
		return nil, ParseError(err)
	}
	return &result, nil
}

// Delete implements store.CaseLinkStore.
func (l *CaseLinkStore) Delete(opts options.Deleter) (*model.CaseLink, error) {
	if opts == nil {
		return nil, errors.InvalidArgument("delete options required")
	}
	if len(opts.GetIDs()) == 0 {
		return nil, errors.InvalidArgument("id required")
	}
	if opts.GetParentID() == 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	fields := []string{"id", "ver", "created_by", "created_at", "updated_by", "updated_at", "author", "name", "url"}
	selectBuilder, err := buildDeleteCaseLinkQuery(opts, fields)
	if err != nil {
		return nil, ParseError(err)
	}
	db, dbErr := l.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	var result model.CaseLink
	if err := pgxscan.Get(opts, db, &result, query, args...); err != nil {
		return nil, ParseError(err)
	}
	return &result, nil
}

// List implements store.CaseLinkStore.
func (l *CaseLinkStore) List(opts options.Searcher) ([]*model.CaseLink, error) {
	if opts == nil {
		return nil, ParseError(errors.New("search options required"))
	}

	filters := opts.GetFilter("case_id")
	if len(filters) == 0 {
		return nil, errors.NewDBError("postgres.case_link.list.check_args.parent_id", "case id required")
	}
	parentId, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil {
		return nil, errors.New("case id is not valid: %s", errors.WithCause(err))
	}
	selectBuilder, err := buildListCaseLinkQuery(opts, parentId)
	if err != nil {
		return nil, ParseError(err)
	}
	db, dbErr := l.storage.Database()
	if dbErr != nil {
		return nil, dbErr
	}
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	var items []*model.CaseLink
	if err := pgxscan.Select(opts, db, &items, query, args...); err != nil {
		return nil, ParseError(err)
	}
	return items, nil
}

// Update implements store.CaseLinkStore.
func (l *CaseLinkStore) Update(opts options.Updator, upd *model.CaseLink) (*model.CaseLink, error) {
	if opts == nil {
		return nil, errors.InvalidArgument("update options required")
	}
	fields := opts.GetFields()
	if len(fields) == 0 {
		fields = []string{"id", "ver", "created_by", "created_at", "updated_by", "updated_at", "author", "name", "url"}
	}
	selectBuilder, err := buildUpdateCaseLinkQuery(opts, upd, fields)
	if err != nil {
		return nil, ParseError(err)
	}
	db, err := l.storage.Database()
	if err != nil {
		return nil, err
	}
	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	var result model.CaseLink
	if err := pgxscan.Get(opts, db, &result, query, args...); err != nil {
		return nil, ParseError(err)
	}
	return &result, nil
}

func NewCaseLinkStore(store *Store) (store.CaseLinkStore, error) {
	if store == nil {
		return nil, errors.New("error creating link case interface to the comment_case table, main store is nil")
	}
	return &CaseLinkStore{storage: store, mainTable: "cases.case_link"}, nil
}

func buildLinkSelectColumns(
	base sq.SelectBuilder,
	tableAlias string,
	fields []string,
) (sq.SelectBuilder, error) {
	var (
		createdByAlias string
		joinCreatedBy  = func(alias string) string {
			if createdByAlias != "" {
				return createdByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.created_by = %s.id", alias, tableAlias, alias))
			createdByAlias = alias
			return alias
		}
		updatedByAlias string
		joinUpdatedBy  = func(alias string) string {
			if updatedByAlias != "" {
				return updatedByAlias
			}
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %s.updated_by = %s.id", alias, tableAlias, alias))
			updatedByAlias = alias
			return alias
		}
		authorAlias string
		joinAuthor  = func(alias string) string {
			if authorAlias != "" {
				return authorAlias
			}
			cb := createdByAlias
			if cb == "" {
				cb = joinCreatedBy("clcb")
			}
			authorAlias = alias
			base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %s.contact_id = %s.id", alias, cb, alias))
			return alias
		}
	)
	base = base.Column(fmt.Sprintf("%s.id", tableAlias))
	for _, field := range fields {
		switch field {
		case "id":
			// already set
		case "name":
			base = base.Column(fmt.Sprintf("%s.name", tableAlias))
		case "url":
			base = base.Column(fmt.Sprintf("%s.url", tableAlias))
		case "ver":
			base = base.Column(fmt.Sprintf("%s.ver", tableAlias))
		case "created_at":
			base = base.Column(fmt.Sprintf("%s.created_at", tableAlias))
		case "updated_at":
			base = base.Column(fmt.Sprintf("%s.updated_at", tableAlias))
		case "created_by":
			cb := "clcb"
			joinCreatedBy(cb)
			base = base.Column(fmt.Sprintf("%s.id AS created_by_id", cb))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) AS created_by_name", cb, cb))
		case "updated_by":
			ub := "club"
			joinUpdatedBy(ub)
			base = base.Column(fmt.Sprintf("%s.id AS updated_by_id", ub))
			base = base.Column(fmt.Sprintf("COALESCE(%s.name, %s.username) AS updated_by_name", ub, ub))
		case "author":
			au := "clau"
			joinAuthor(au)
			base = base.Column(fmt.Sprintf("%s.id AS contact_id", au))
			base = base.Column(fmt.Sprintf("%s.common_name AS contact_name", au))
		default:
			return base, errors.InvalidArgument("unknown field: " + field)
		}
	}
	return base, nil
}

func buildCreateCaseLinkQuery(
	rpc options.Creator,
	input *model.CaseLink,
	fields []string,
) (sq.SelectBuilder, error) {
	userID := rpc.GetAuthOpts().GetUserId()
	if input != nil && input.Author != nil && input.Author.Id != nil {
		userID = int64(*input.Author.Id)
	}
	insertBuilder := sq.Insert("cases.case_link").
		Columns("created_by", "updated_by", "name", "url", "case_id", "dc").
		Values(userID, userID, input.Name, input.Url, rpc.GetParentID(), rpc.GetAuthOpts().GetDomainId()).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")
	insertSQL, args, err := insertBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, ParseError(err)
	}
	cte := sq.Expr("WITH cl AS ("+insertSQL+")", args...)
	selectBuilder, err := buildLinkSelectColumns(sq.Select(), caseLinkLeft, fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	selectBuilder = selectBuilder.PrefixExpr(cte).From(caseLinkLeft)
	return selectBuilder, nil
}

func buildUpdateCaseLinkQuery(
	opts options.Updator,
	input *model.CaseLink,
	fields []string,
) (sq.SelectBuilder, error) {
	if len(opts.GetEtags()) == 0 {
		return sq.SelectBuilder{}, errors.InvalidArgument("link etag required")
	}
	if len(opts.GetMask()) == 0 {
		return sq.SelectBuilder{}, errors.InvalidArgument("link update mask required")
	}
	tid := opts.GetEtags()[0]
	userID := opts.GetAuthOpts().GetUserId()
	if input != nil && input.Author != nil && input.Author.Id != nil {
		userID = int64(*input.Author.Id)
	}
	updateBuilder := sq.Update("cases.case_link").
		Set("updated_by", userID).
		Set("updated_at", opts.RequestTime()).
		Set("ver", sq.Expr("ver+1")).
		Where("id = ?", tid.GetOid()).
		Where("ver = ?", tid.GetVer()).
		Where("dc = ?", opts.GetAuthOpts().GetDomainId()).
		Where("case_id = ?", opts.GetParentID()).
		PlaceholderFormat(sq.Dollar)
	for _, field := range opts.GetMask() {
		switch field {
		case "url":
			_, err := url.Parse(input.Url)
			if err != nil {
				return sq.SelectBuilder{}, ParseError(err)
			}
			updateBuilder = updateBuilder.Set("url", input.Url)
		case "name":
			updateBuilder = updateBuilder.Set("name", input.Name)
		}
	}
	updateSQL, args, err := updateBuilder.Suffix("RETURNING *").ToSql()
	if err != nil {
		return sq.SelectBuilder{}, ParseError(err)
	}
	cte := sq.Expr("WITH cl AS ("+updateSQL+")", args...)
	selectBuilder, err := buildLinkSelectColumns(sq.Select(), caseLinkLeft, fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	selectBuilder = selectBuilder.PrefixExpr(cte).From(caseLinkLeft)
	return selectBuilder, nil
}

func buildDeleteCaseLinkQuery(
	opts options.Deleter,
	fields []string,
) (sq.SelectBuilder, error) {
	deleteBuilder := sq.Delete("cases.case_link").
		Where("id = ANY(?)", opts.GetIDs()).
		Where("dc = ?", opts.GetAuthOpts().GetDomainId()).
		Where("case_id = ?", opts.GetParentID()).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")
	deleteSQL, args, err := deleteBuilder.ToSql()
	if err != nil {
		return sq.SelectBuilder{}, ParseError(err)
	}
	cte := sq.Expr("WITH cl AS ("+deleteSQL+")", args...)
	selectBuilder, err := buildLinkSelectColumns(sq.Select(), caseLinkLeft, fields)
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	selectBuilder = selectBuilder.PrefixExpr(cte).From(caseLinkLeft)
	return selectBuilder, nil
}

func buildListCaseLinkQuery(
	opts options.Searcher,
	parentId int64,
) (sq.SelectBuilder, error) {
	base := sq.Select().From("cases.case_link "+caseLinkLeft).
		Where(fmt.Sprintf("%s = ?", storeUtil.Ident(caseLinkLeft, "dc")), opts.GetAuthOpts().GetDomainId()).
		Where(fmt.Sprintf("%s = ?", storeUtil.Ident(caseLinkLeft, "case_id")), parentId).
		PlaceholderFormat(sq.Dollar)
	if len(opts.GetIDs()) != 0 {
		base = base.Where(fmt.Sprintf("%s = any(?)", storeUtil.Ident(caseLinkLeft, "id")), opts.GetIDs())
	}
	base = storeUtil.ApplyPaging(opts.GetPage(), opts.GetSize(), base)
	base = storeUtil.ApplyDefaultSorting(opts, base, linkDefaultSort)
	return buildLinkSelectColumns(base, caseLinkLeft, opts.GetFields())
}

// Deprecated. Use until cases service is refactored.
func buildLinkSelectAsSubquery(fields []string, caseAlias string) (updatedBase sq.SelectBuilder, scanPlan []func(link *_go.CaseLink) any, dbErr *errors.DBError) {
	alias := "links"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.
		Select().
		From("cases.case_link " + alias).
		Where(fmt.Sprintf("%s = %s", storeUtil.Ident(alias, "case_id"), storeUtil.Ident(caseAlias, "id")))

	base, plan, dbErr := buildLinkSelectColumnsAndPlan(base, alias, fields)
	if dbErr != nil {
		return base, nil, dbErr
	}
	base = storeUtil.ApplyPaging(1, defaults.DefaultSearchSize, base)

	return base, plan, nil
}

// Deprecated. Use until cases service is refactored.
func buildLinkSelectColumnsAndPlan(base sq.SelectBuilder, left string, fields []string) (sq.SelectBuilder, []func(link *_go.CaseLink) any, *errors.DBError) {
	var (
		plan           []func(link *_go.CaseLink) any
		createdByAlias string
		joinCreatedBy  = func() {
			if createdByAlias != "" {
				return
			}
			createdByAlias = caseLinkCreatedByAlias
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.created_by", createdByAlias, left))
			return
		}
		updatedByAlias string
		joinUpdatedBy  = func() {
			if updatedByAlias != "" {
				return
			}
			updatedByAlias = caseLinkUpdatedByAlias
			base = base.LeftJoin(fmt.Sprintf("directory.wbt_user %s ON %[1]s.id = %s.updated_by", updatedByAlias, left))
			return
		}
		authorAlias string
		joinAuthor  = func() {
			if authorAlias != "" {
				return
			}
			joinCreatedBy()
			authorAlias = caseLinkAuthorAlias
			base = base.LeftJoin(fmt.Sprintf("contacts.contact %s ON %[1]s.id = %s.contact_id", authorAlias, createdByAlias))
			return
		}
	)
	if len(fields) == 0 {
		fields = CaseLinkFields
	}

	for _, field := range fields {
		switch field {
		case "id":
			base = base.Column(storeUtil.Ident(left, "id"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Id
			})
		case "ver":
			base = base.Column(storeUtil.Ident(left, "ver"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Ver
			})
		case "created_by":
			joinCreatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, %[1]s.name)::text created_by", createdByAlias))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.CreatedBy)
			})
		case "created_at":
			base = base.Column(storeUtil.Ident(left, "created_at"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanTimestamp(&link.CreatedAt)
			})
		case "updated_by":
			joinUpdatedBy()
			base = base.Column(fmt.Sprintf("ROW(%[1]s.id, %[1]s.name)::text updated_by", updatedByAlias))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.UpdatedBy)
			})
		case "updated_at":
			base = base.Column(storeUtil.Ident(left, "updated_at"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanTimestamp(&link.UpdatedAt)
			})
		case "name":
			base = base.Column(storeUtil.Ident(left, "name"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanText(&link.Name)
			})
		case "url":
			base = base.Column(storeUtil.Ident(left, "url"))
			plan = append(plan, func(link *_go.CaseLink) any {
				return &link.Url
			})
		case "author":
			joinAuthor()
			base = base.Column(fmt.Sprintf(`ROW(%[1]s.id, %[1]s.common_name)::text author`, authorAlias))
			plan = append(plan, func(link *_go.CaseLink) any {
				return scanner.ScanRowLookup(&link.Author)
			})
		default:
			return base, nil, errors.NewDBError("postgres.case_link.build_link_select.cycle_fields.unknown", fmt.Sprintf("%s field is unknown", field))
		}
	}
	if len(plan) == 0 {
		return base, nil, errors.NewDBError("postgres.case_link.build_link_select.final_check.unknown", "no resulting columns")
	}
	return base, plan, nil
}

func buildCreateLinkQuery(
	rpc options.Creator,
	fields []string,
	input *_go.InputCaseLink,
) (
	sq.Sqlizer,
	[]func(link *_go.CaseLink) any,
	error,
) {
	// Default user from token
	userID := rpc.GetAuthOpts().GetUserId()

	// Override if input.CreatedBy is explicitly provided
	if createdBy := input.GetUserID(); createdBy != nil && createdBy.Id != 0 {
		userID = createdBy.Id
	}
	insertAlias := "i"
	insert := sq.
		Insert("cases.case_link").
		Columns("created_by", "updated_by", "name", "url", "case_id", "dc").
		Values(userID, userID, input.GetName(), input.GetUrl(), rpc.GetParentID(), rpc.GetAuthOpts().GetDomainId()).
		Suffix("RETURNING *")
	// select
	query, args, _ := storeUtil.FormAsCTE(insert, insertAlias)
	base := sq.Select().From(insertAlias).Prefix(query, args...).PlaceholderFormat(sq.Dollar)
	// build plan from columns
	return buildLinkSelectColumnsAndPlan(base, insertAlias, fields)
}

func buildUpdateLinkQuery(opts options.Updator, input *_go.InputCaseLink) (sq.Sqlizer, []func(link *_go.CaseLink) any, error) {
	base := sq.Select().From("cases.case_link cl")
	if len(opts.GetEtags()) == 0 {
		return base, nil, errors.InvalidArgument("link etag required", errors.WithID("postgres.case_link.update.etag.empty"))
	}
	if len(opts.GetMask()) == 0 {
		return nil, nil, errors.InvalidArgument("link update mask required", errors.WithID("postgres.case_link.update.mask.empty"))
	}
	tid := opts.GetEtags()[0]

	userID := opts.GetAuthOpts().GetUserId()
	if util.ContainsField(opts.GetMask(), "userID") {
		if updatedBy := input.GetUserID(); updatedBy != nil && updatedBy.Id != 0 {
			userID = updatedBy.Id
		}
	}

	// insert
	update := sq.
		Update("cases.case_link").
		Set("updated_by", userID).
		Set("updated_at", opts.RequestTime()).
		Set("ver", sq.Expr("ver+1")).
		Where("id = ?", tid.GetOid()).
		Where("ver = ?", tid.GetVer()).
		Where("dc = ?", opts.GetAuthOpts().GetDomainId()).
		Where("case_id = ?", opts.GetParentID()).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)
	for _, field := range opts.GetMask() {
		switch field {
		case "url":
			_, err := url.Parse(input.Url)
			if err != nil {
				return nil, nil, errors.InvalidArgument(err.Error(), errors.WithID("postgres.case_link.build_update_query.url.validate.error"))
			}
			update = update.Set("url", input.Url)
		case "name":
			update = update.Set("name", input.Name)
		}
	}
	prefixAlias := "upd"
	prefix, args, err := storeUtil.FormAsCTE(update, prefixAlias)
	if err != nil {
		return nil, nil, errors.New(err.Error(), errors.WithID("postgres.case_link.build_update_query.form_cte.error"))
	}
	slct := sq.Select().Prefix(prefix, args...).From(prefixAlias)
	// build plan from columns
	return buildLinkSelectColumnsAndPlan(slct, prefixAlias, opts.GetFields())
}

func (l *CaseLinkStore) scanLinks(rows pgx.Rows, plan []func(link *_go.CaseLink) any) ([]*_go.CaseLink, error) {
	var res []*_go.CaseLink

	for rows.Next() {
		link, err := l.scanLink(rows, plan)
		if err != nil {
			return nil, err
		}
		res = append(res, link)
	}
	return res, nil
}

func (l *CaseLinkStore) scanLink(row pgx.Row, plan []func(link *_go.CaseLink) any) (*_go.CaseLink, error) {
	var link _go.CaseLink
	var scanPlan []any
	for _, scan := range plan {
		scanPlan = append(scanPlan, scan(&link))
	}
	err := row.Scan(scanPlan...)
	if err != nil {
		return nil, err
	}

	return &link, nil
}
func ValidateLinkCreate(caseId int64, input *model.CaseLink) error {
	if caseId <= 0 {
		return errors.InvalidArgument("case id required")
	}
	if input == nil || input.Url == "" {
		return errors.InvalidArgument("input for link required")
	}
	_, err := url.Parse(input.Url)
	if err != nil {
		return ParseError(err)
	}
	return nil
}
