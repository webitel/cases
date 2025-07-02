package postgres

import (
	"fmt"
	"net/url"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/model/options/defaults"
	"github.com/webitel/cases/internal/store"
	dbutil "github.com/webitel/cases/internal/store/util"
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
	parentId, ok := opts.GetFilter("case_id").(int64)
	if !ok || parentId == 0 {
		return nil, ParseError(errors.New("case id required"))
	}
	fields := opts.GetFields()
	if len(fields) == 0 {
		fields = []string{"id", "ver", "created_by", "created_at", "updated_by", "updated_at", "author", "name", "url"}
	}
	selectBuilder, err := buildListCaseLinkQuery(opts, parentId, fields)
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
	fields []string,
) (sq.SelectBuilder, error) {
	base := sq.Select().From("cases.case_link "+caseLinkLeft).
		Where(fmt.Sprintf("%s = ?", dbutil.Ident(caseLinkLeft, "dc")), opts.GetAuthOpts().GetDomainId()).
		Where(fmt.Sprintf("%s = ?", dbutil.Ident(caseLinkLeft, "case_id")), parentId).
		PlaceholderFormat(sq.Dollar)
	if len(opts.GetIDs()) != 0 {
		base = base.Where(fmt.Sprintf("%s = any(?)", dbutil.Ident(caseLinkLeft, "id")), opts.GetIDs())
	}
	base = dbutil.ApplyPaging(opts.GetPage(), opts.GetSize(), base)
	base = dbutil.ApplyDefaultSorting(opts, base, linkDefaultSort)
	return buildLinkSelectColumns(base, caseLinkLeft, fields)
}

func buildLinkSelectAsSubquery(fields []string, caseAlias string) (updatedBase sq.SelectBuilder, dbErr error) {//for cases service (need to refactor cases first)
	alias := "links"
	if caseAlias == alias {
		alias = "sub_" + alias
	}
	base := sq.Select().From("cases.case_link " + alias).
		Where(fmt.Sprintf("%s = %s", dbutil.Ident(alias, "case_id"), dbutil.Ident(caseAlias, "id")))

	base, dbErr = buildLinkSelectColumns(base, alias, fields)
	if dbErr != nil {
		return base, dbErr
	}
	base = dbutil.ApplyPaging(1, defaults.DefaultSearchSize, base)

	return base, nil
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
