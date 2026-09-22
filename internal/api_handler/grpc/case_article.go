package grpc

import (
	"context"

	"github.com/webitel/cases/api/cases"
	grpcoptions "github.com/webitel/cases/internal/api_handler/grpc/options"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
)

type CaseArticleHandler interface {
	ListCaseArticles(options.Searcher) ([]*model.CaseArticle, error)
	LinkCaseArticle(options.Creator, *model.CaseArticle) (*model.CaseArticle, error)
	UnlinkCaseArticle(options.Deleter) (*model.CaseArticle, error)
}

type CaseArticleService struct {
	cases.UnimplementedCaseArticlesServer

	app CaseArticleHandler
}

func NewCaseArticleService(app CaseArticleHandler) (*CaseArticleService, error) {
	if app == nil {
		return nil, errors.New("case article handler is nil")
	}
	return &CaseArticleService{app: app}, nil
}

func (s *CaseArticleService) ListCaseArticles(ctx context.Context, req *cases.ListCaseArticlesRequest) (*cases.CaseArticleList, error) {
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}
	return s.list(ctx, req, util.EqualFilter("case_id", caseTid.GetOid()))
}

func (s *CaseArticleService) ListArticleCases(ctx context.Context, req *cases.ListArticleCasesRequest) (*cases.CaseArticleList, error) {
	if req.GetArticleId() <= 0 {
		return nil, errors.InvalidArgument("article id is required")
	}
	return s.list(ctx, req, util.EqualFilter("article_id", req.GetArticleId()))
}

func (s *CaseArticleService) LinkCaseArticle(ctx context.Context, req *cases.LinkCaseArticleRequest) (*cases.CaseArticle, error) {
	articleID := req.GetInput().GetArticle().GetId()
	if articleID <= 0 {
		return nil, errors.InvalidArgument("article id is required")
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}
	createOpts, err := grpcoptions.NewCreateOptions(ctx, grpcoptions.WithCreateParentID(caseTid.GetOid()))
	if err != nil {
		return nil, err
	}
	link, err := s.app.LinkCaseArticle(createOpts, &model.CaseArticle{ArticleID: articleID})
	if err != nil {
		return nil, err
	}
	return s.Marshal(link)
}

func (s *CaseArticleService) UnlinkCaseArticle(ctx context.Context, req *cases.UnlinkCaseArticleRequest) (*cases.CaseArticle, error) {
	if req.GetArticleId() <= 0 {
		return nil, errors.InvalidArgument("article id is required")
	}
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseEtag())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}
	deleteOpts, err := grpcoptions.NewDeleteOptions(ctx,
		grpcoptions.WithDeleteID(req.GetArticleId()),
		grpcoptions.WithDeleteParentID(caseTid.GetOid()),
	)
	if err != nil {
		return nil, err
	}
	link, err := s.app.UnlinkCaseArticle(deleteOpts)
	if err != nil {
		return nil, err
	}
	return s.Marshal(link)
}

func (s *CaseArticleService) list(ctx context.Context, req grpcoptions.Pager, filter string) (*cases.CaseArticleList, error) {
	searchOpts, err := grpcoptions.NewSearchOptions(ctx, grpcoptions.WithPagination(req))
	if err != nil {
		return nil, err
	}

	searchOpts.AddFilter(filter)

	links, err := s.app.ListCaseArticles(searchOpts)
	if err != nil {
		return nil, err
	}
	converted, err := utils.ConvertToOutputBulk(links, s.Marshal)
	if err != nil {
		return nil, err
	}
	var res cases.CaseArticleList
	res.Next, res.Items = utils.GetListResult(searchOpts, converted)
	res.Page = int64(searchOpts.GetPage())
	return &res, nil
}

// Marshal converts a model.CaseArticle to cases.CaseArticle.
func (s *CaseArticleService) Marshal(m *model.CaseArticle) (*cases.CaseArticle, error) {
	out := &cases.CaseArticle{
		Article:   &cases.Lookup{Id: m.ArticleID},
		Source:    cases.CaseArticleSource(m.Source),
		CreatedAt: utils.MarshalTime(m.CreatedAt),
	}
	if m.Author != nil && m.Author.Id != nil {
		out.CreatedBy = utils.MarshalLookup(m.Author)
	}
	if m.CaseID != nil {
		ver := utils.Dereference(m.CaseVer)
		caseEtag, err := etag.EncodeEtag(etag.EtagCase, *m.CaseID, ver)
		if err != nil {
			return nil, err
		}
		out.Case = &cases.RelatedCaseLookup{
			Id:      *m.CaseID,
			Etag:    caseEtag,
			Ver:     ver,
			Name:    utils.Dereference(m.CaseName),
			Subject: utils.Dereference(m.CaseSubject),
		}
	}
	return out, nil
}
