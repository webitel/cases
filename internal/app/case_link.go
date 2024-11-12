package app

import (
	"context"
	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/webitel-go-kit/etag"
)

// In search options extract from context user
// Remove from search options fields functions

type CaseLinkService struct {
	app *App
	cases.UnimplementedCaseLinksServer
}

var DefaultCaseLinkFields = []string{
	"etag", "created_by", "created_at", "author", "name", "url",
}

func (c *CaseLinkService) LocateLink(ctx context.Context, req *cases.LocateLinkRequest) (*cases.CaseLink, error) {
	// Validate required fields
	//if req.Etag == "" {
	//	return nil, cerror.NewBadRequestError("app.case_link.locate.case_etag.check_args.etag", "Etag is required")
	//}
	//
	//// Convert the etag to an internal identifier (Tid) for filtering by ID and Ver
	////etag, err := etag.EtagOrId(etag.EtagCaseLink, req.Etag)
	////if err != nil {
	////	return nil, cerror.NewBadRequestError("app.case_link.locate.case_etag.check_args.etag", err.Error())
	////}
	//
	//searchOpts := model.NewLocateOptions(ctx, req, DefaultCaseLinkFields)
	//
	//// Use ListComments to retrieve the specific comment
	//commentList, err := c.app.Store.CaseComment().List(searchOpts)
	//if err != nil {
	//	return nil, cerror.NewInternalError("case_comment_service.locate_comment.fetch_error", err.Error())
	//}
	//
	//// Ensure we found exactly one comment
	//if len(commentList.Items) == 0 {
	//	return nil, cerror.NewNotFoundError("case_comment_service.locate_comment.not_found", "Comment not found")
	//} else if len(commentList.Items) > 1 {
	//	return nil, cerror.NewInternalError("case_comment_service.locate_comment.multiple_found", "Multiple comments found")
	//}

	// Return the located comment
	//return commentList.Items[0], nil
	return nil, nil
}

func (c *CaseLinkService) CreateLink(ctx context.Context, req *cases.CreateLinkRequest) (*cases.CaseLink, error) {

	// Validate request
	if req.CaseEtag == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.case_etag", "Case etag is required")
	} else if req.Input.GetUrl() == "" {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.check_args.url", "Url is required for each link")
	}
	caseTID, err := etag.EtagOrId(etag.EtagCaseLink, req.CaseEtag)
	if err != nil {
		return nil, cerror.NewBadRequestError("app.case_link.create.case_etag.parse.error", err.Error())
	}

	createOpts := model.NewCreateOptions(ctx, req, DefaultCaseLinkFields)
	createOpts.ParentID = caseTID.GetOid()
	res, dbErr := c.app.Store.CaseLink().Create(createOpts, req.Input)
	if dbErr != nil {
		return nil, dbErr
	}

	if createOpts.HasEtag() {
		res.Etag = etag.EncodeEtag(etag.EtagCaseLink, res.Id, res.Ver)
		// hide
		if !createOpts.HasId() {
			res.Id = 0
		}
		if !createOpts.HasVer() {
			res.Ver = 0
		}
	}

	return res, nil
}

func (c *CaseLinkService) UpdateLink(ctx context.Context, req *cases.UpdateLinkRequest) (*cases.CaseLink, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) DeleteLink(ctx context.Context, req *cases.DeleteLinkRequest) (*cases.CaseLink, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) ListLinks(ctx context.Context, req *cases.ListLinksRequest) (*cases.CaseLinkList, error) {
	//searchOpts := model.NewSearchOptions(ctx, req)
	//// output: validate & normalize & defaults
	//graphLinkModel := struct {
	//	graph.Query
	//	FieldsParse func(rawFields []string, decode ...graph.FieldEncoding) (fields graph.FieldsQ, err error)
	//	// Output      func(*cases.CaseLinkList, *graph.Query)
	//}{
	//	Query: graph.Query{
	//		Name: "listLinks",
	//	},
	//	FieldsParse: casegraph.Schema.Case.Link.Output.ParseFields,
	//}
	//graphParsedFields, err := graphLinkModel.FieldsParse(searchOpts.Fields)
	//if err != nil {
	//	return nil, err
	//}
	//graphLinkModel.Fields = graphParsedFields
	// output: validate & normalize & defaults

	panic("implement me")
}

func NewCaseLinkService(app *App) (*CaseLinkService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseLinkService{app: app}, nil
}
