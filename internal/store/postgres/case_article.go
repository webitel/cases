package postgres

import (
	"context"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
	storeUtil "github.com/webitel/cases/internal/store/util"
)

type CaseArticleStore struct {
	storage *Store
}

func NewCaseArticleStore(store *Store) (store.CaseArticleStore, error) {
	if store == nil {
		return nil, errors.New("error creating case article interface, main store is nil")
	}
	return &CaseArticleStore{storage: store}, nil
}

// Link implements store.CaseArticleStore.
func (s *CaseArticleStore) Link(rpc options.Creator, add *model.CaseArticle) (*model.CaseArticle, error) {
	if rpc == nil {
		return nil, errors.InvalidArgument("create options required")
	}
	if rpc.GetParentID() <= 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	if add == nil || add.ArticleID <= 0 {
		return nil, errors.InvalidArgument("article id required")
	}
	return s.get(rpc, buildLinkCaseArticleQuery(rpc, add))
}

// List implements store.CaseArticleStore.
func (s *CaseArticleStore) List(opts options.Searcher) ([]*model.CaseArticle, error) {
	if opts == nil {
		return nil, errors.InvalidArgument("search options required")
	}
	query, err := buildListCaseArticleQuery(opts)
	if err != nil {
		return nil, err
	}
	db, err := s.storage.Database()
	if err != nil {
		return nil, err
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	var items []*model.CaseArticle
	if err := pgxscan.Select(opts, db, &items, sql, args...); err != nil {
		return nil, ParseError(err)
	}
	return items, nil
}

// Unlink implements store.CaseArticleStore.
func (s *CaseArticleStore) Unlink(opts options.Deleter) (*model.CaseArticle, error) {
	if opts == nil {
		return nil, errors.InvalidArgument("delete options required")
	}
	if opts.GetParentID() <= 0 {
		return nil, errors.InvalidArgument("case id required")
	}
	if len(opts.GetIDs()) != 1 || opts.GetIDs()[0] <= 0 {
		return nil, errors.InvalidArgument("article id required")
	}
	return s.get(opts, buildUnlinkCaseArticleQuery(opts))
}

func (s *CaseArticleStore) get(ctx context.Context, query sq.SelectBuilder) (*model.CaseArticle, error) {
	db, err := s.storage.Database()
	if err != nil {
		return nil, err
	}
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, ParseError(err)
	}
	var result model.CaseArticle
	if err := pgxscan.Get(ctx, db, &result, sql, args...); err != nil {
		return nil, ParseError(err)
	}
	return &result, nil
}

// buildCaseArticleLinks selects the manual links and the close article links as one set.
func buildCaseArticleLinks(domainID int64) (manual, resolution sq.SelectBuilder) {
	manual = sq.Select(
		"ca.case_id",
		"ca.article_id",
		fmt.Sprintf("%d AS source", model.CaseArticleSourceManual),
		"ca.created_by",
		"ca.created_at",
	).
		From("cases.case_article ca").
		Join(`cases."case" mc ON mc.id = ca.case_id`).
		Where("ca.dc = ?", domainID).
		Where("mc.close_article_id IS DISTINCT FROM ca.article_id")
	resolution = sq.Select(
		"rc.id AS case_id",
		"rc.close_article_id AS article_id",
		fmt.Sprintf("%d AS source", model.CaseArticleSourceResolution),
		"NULL::bigint AS created_by",
		"NULL::timestamp AS created_at",
	).
		From(`cases."case" rc`).
		Where("rc.dc = ?", domainID).
		Where("rc.close_article_id IS NOT NULL")
	return manual, resolution
}

// selectCaseArticles reads the links of the prefixed CTE with their authors.
func selectCaseArticles(prefix sq.Sqlizer) sq.SelectBuilder {
	return sq.Select(
		"k.article_id",
		"k.source",
		"k.created_at",
		"cau.id AS created_by_id",
		"COALESCE(cau.name, cau.username) AS created_by_name",
	).
		PrefixExpr(prefix).
		From("k").
		LeftJoin("directory.wbt_user cau ON cau.id = k.created_by").
		PlaceholderFormat(sq.Dollar)
}

func buildListCaseArticleQuery(opts options.Searcher) (sq.SelectBuilder, error) {
	sess := opts.GetAuthOpts()
	manual, resolution := buildCaseArticleLinks(sess.GetDomainId())

	caseID, err := caseArticleFilter(opts, "case_id")
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	articleID, err := caseArticleFilter(opts, "article_id")
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	var query sq.SelectBuilder

	switch {
	case caseID > 0:
		manual = manual.Where("ca.case_id = ?", caseID)
		resolution = resolution.Where("rc.id = ?", caseID)
		query = selectCaseArticles(sq.Expr("WITH k AS (? UNION ALL ?)", manual, resolution)).
			OrderBy("k.source DESC", "k.created_at DESC", "k.article_id DESC")
	case articleID > 0:
		manual = manual.Where("ca.article_id = ?", articleID)
		resolution = resolution.Where("rc.close_article_id = ?", articleID)
		rbac, err := getCaseRbacCondition(sess, auth.Read, "k.case_id")
		if err != nil {
			return sq.SelectBuilder{}, err
		}
		query = selectCaseArticles(sq.Expr("WITH k AS (? UNION ALL ?)", manual, resolution)).
			Columns("lc.id AS case_id", "lc.ver AS case_ver", "lc.name AS case_name", "lc.subject AS case_subject").
			Join(`cases."case" lc ON lc.id = k.case_id`).
			Where(rbac).
			OrderBy("k.case_id DESC")
	default:
		return sq.SelectBuilder{}, errors.InvalidArgument("case id or article id required")
	}
	return storeUtil.ApplyPaging(opts.GetPage(), opts.GetSize(), query), nil
}

func buildLinkCaseArticleQuery(rpc options.Creator, add *model.CaseArticle) sq.SelectBuilder {
	var (
		sess      = rpc.GetAuthOpts()
		caseID    = rpc.GetParentID()
		articleID = add.ArticleID
	)
	insert := sq.Insert("cases.case_article").
		Columns("dc", "case_id", "article_id", "created_by").
		Select(sq.Select().
			Columns("c.dc", "c.id").
			Column("?::bigint", articleID).
			Column("?::bigint", sess.GetUserId()).
			From(`cases."case" c`).
			Where("c.id = ?", caseID).
			Where("c.dc = ?", sess.GetDomainId()).
			Where("c.close_article_id IS DISTINCT FROM ?", articleID)).
		Suffix("ON CONFLICT (case_id, article_id) DO UPDATE SET article_id = EXCLUDED.article_id " +
			"RETURNING article_id, created_by, created_at")
	inserted := sq.Select(
		"i.article_id",
		fmt.Sprintf("%d AS source", model.CaseArticleSourceManual),
		"i.created_by",
		"i.created_at",
	).From("i")
	return selectCaseArticles(sq.Expr("WITH i AS (?), k AS (? UNION ALL ?)",
		insert, inserted, buildCloseArticleLink(sess.GetDomainId(), caseID, articleID)))
}

func buildUnlinkCaseArticleQuery(opts options.Deleter) sq.SelectBuilder {
	var (
		domainID  = opts.GetAuthOpts().GetDomainId()
		caseID    = opts.GetParentID()
		articleID = opts.GetIDs()[0]
	)
	del := sq.Delete("cases.case_article").
		Where("dc = ?", domainID).
		Where("case_id = ?", caseID).
		Where("article_id = ?", articleID).
		Where(`NOT EXISTS (SELECT 1 FROM cases."case" uc WHERE uc.id = ? AND uc.close_article_id = ?)`, caseID, articleID).
		Suffix("RETURNING article_id, created_by, created_at")
	deleted := sq.Select(
		"d.article_id",
		fmt.Sprintf("%d AS source", model.CaseArticleSourceManual),
		"d.created_by",
		"d.created_at",
	).From("d")
	return selectCaseArticles(sq.Expr("WITH d AS (?), k AS (? UNION ALL ?)",
		del, deleted, buildCloseArticleLink(domainID, caseID, articleID)))
}

// buildCloseArticleLink selects the article when it is the close article of the case; writes leave it untouched.
func buildCloseArticleLink(domainID, caseID, articleID int64) sq.SelectBuilder {
	return sq.Select(
		"rc.close_article_id",
		fmt.Sprintf("%d", model.CaseArticleSourceResolution),
		"NULL::bigint",
		"NULL::timestamp",
	).
		From(`cases."case" rc`).
		Where("rc.id = ?", caseID).
		Where("rc.dc = ?", domainID).
		Where("rc.close_article_id = ?", articleID)
}

func caseArticleFilter(opts options.Searcher, field string) (int64, error) {
	filters := opts.GetFilter(field)
	if len(filters) == 0 {
		return 0, nil
	}
	id, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.InvalidArgument(fmt.Sprintf("invalid %s", field))
	}
	return id, nil
}
