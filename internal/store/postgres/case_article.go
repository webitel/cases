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
	"github.com/webitel/cases/internal/store/postgres/transaction"
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

// selectCaseArticles reads the links of the k source with their authors.
func selectCaseArticles() sq.SelectBuilder {
	return sq.Select(
		"k.article_id",
		"k.source",
		"k.created_at",
		"cau.id AS created_by_id",
		"COALESCE(cau.name, cau.username) AS created_by_name",
	).
		LeftJoin("directory.wbt_user cau ON cau.id = k.created_by").
		PlaceholderFormat(sq.Dollar)
}

func buildListCaseArticleQuery(opts options.Searcher) (sq.SelectBuilder, error) {
	sess := opts.GetAuthOpts()

	caseID, err := caseArticleFilter(opts, "case_id")
	if err != nil {
		return sq.SelectBuilder{}, err
	}
	articleID, err := caseArticleFilter(opts, "article_id")
	if err != nil {
		return sq.SelectBuilder{}, err
	}

	query := selectCaseArticles().
		From("cases.case_article k").
		Where("k.dc = ?", sess.GetDomainId())

	switch {
	case caseID > 0:
		query = query.
			Where("k.case_id = ?", caseID).
			OrderBy("k.source DESC", "k.created_at DESC", "k.article_id DESC")
	case articleID > 0:
		rbac, err := getCaseRbacCondition(sess, auth.Read, "k.case_id")
		if err != nil {
			return sq.SelectBuilder{}, err
		}
		query = query.
			Columns("lc.id AS case_id", "lc.ver AS case_ver", "lc.name AS case_name", "lc.subject AS case_subject").
			Join(`cases."case" lc ON lc.id = k.case_id`).
			Where("k.article_id = ?", articleID).
			Where(rbac).
			OrderBy("k.case_id DESC")
	default:
		return sq.SelectBuilder{}, errors.InvalidArgument("case id or article id required")
	}
	return storeUtil.ApplyPaging(opts.GetPage(), opts.GetSize(), query), nil
}

// buildLinkCaseArticleQuery adds a manual link; an existing link, the close article included, is returned as is.
func buildLinkCaseArticleQuery(rpc options.Creator, add *model.CaseArticle) sq.SelectBuilder {
	sess := rpc.GetAuthOpts()
	insert := sq.Insert("cases.case_article").
		Columns("dc", "case_id", "article_id", "created_by").
		Select(sq.Select().
			Columns("c.dc", "c.id").
			Column("?::bigint", add.ArticleID).
			Column("?::bigint", sess.GetUserId()).
			From(`cases."case" c`).
			Where("c.id = ?", rpc.GetParentID()).
			Where("c.dc = ?", sess.GetDomainId())).
		Suffix("ON CONFLICT (case_id, article_id) DO UPDATE SET article_id = EXCLUDED.article_id " +
			"RETURNING article_id, source, created_by, created_at")
	return selectCaseArticles().PrefixExpr(sq.Expr("WITH k AS (?)", insert)).From("k")
}

// buildUnlinkCaseArticleQuery removes a manual link; the close article is returned untouched.
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
		Where("source = ?", model.CaseArticleSourceManual).
		Suffix("RETURNING article_id, source, created_by, created_at")
	kept := sq.Select("article_id", "source", "created_by", "created_at").
		From("cases.case_article").
		Where("dc = ?", domainID).
		Where("case_id = ?", caseID).
		Where("article_id = ?", articleID).
		Where("source = ?", model.CaseArticleSourceResolution)
	return selectCaseArticles().PrefixExpr(sq.Expr("WITH d AS (?), k AS (SELECT * FROM d UNION ALL ?)", del, kept)).From("k")
}

// setCloseArticle makes the article the only close article link of the case; nil clears it.
func setCloseArticle(ctx context.Context, tx *transaction.TxManager, domainID, caseID, userID int64, articleID *int64) error {
	_, err := tx.Exec(ctx,
		`DELETE FROM cases.case_article WHERE case_id = $1 AND source = $2 AND article_id IS DISTINCT FROM $3`,
		caseID, model.CaseArticleSourceResolution, articleID)
	if err != nil || articleID == nil {
		return err
	}
	_, err = tx.Exec(ctx,
		`INSERT INTO cases.case_article (dc, case_id, article_id, source, created_by)
		VALUES ($1, $2, $3, $4, NULLIF($5, 0))
		ON CONFLICT (case_id, article_id) DO UPDATE
		SET source = EXCLUDED.source, created_by = EXCLUDED.created_by, created_at = EXCLUDED.created_at
		WHERE case_article.source <> EXCLUDED.source`,
		domainID, caseID, *articleID, model.CaseArticleSourceResolution, userID)
	return err
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
