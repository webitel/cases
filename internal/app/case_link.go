package app

import (
	"context"

	cases "github.com/webitel/cases/api/cases"
	casegraph "github.com/webitel/cases/internal/app/graph"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/model/graph"
)

// In search options extract from context user
// Remove from search options fields functions

type CaseLinkService struct {
	app *App
	cases.UnimplementedCaseLinksServer
}

func (c *CaseLinkService) LocateLink(ctx context.Context, request *cases.LocateLinkRequest) (*cases.CaseLink, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) CreateLink(ctx context.Context, request *cases.CreateLinkRequest) (*cases.CaseLink, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) UpdateLink(ctx context.Context, request *cases.UpdateLinkRequest) (*cases.CaseLink, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) DeleteLink(ctx context.Context, request *cases.DeleteLinkRequest) (*cases.CaseLink, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) ListLinks(ctx context.Context, request *cases.ListLinksRequest) (*cases.CaseLinkList, error) {
	searchOpts := model.NewSearchOptions(ctx, request)
	// output: validate & normalize & defaults
	graphLinkModel := struct {
		graph.Query
		FieldsParse func(rawFields []string, decode ...graph.FieldEncoding) (fields graph.FieldsQ, err error)
		// Output      func(*cases.CaseLinkList, *graph.Query)
	}{
		Query: graph.Query{
			Name: "listLinks",
		},
		FieldsParse: casegraph.Schema.Case.Link.Output.ParseFields,
	}
	graphParsedFields, err := graphLinkModel.FieldsParse(searchOpts.Fields)
	if err != nil {
		return nil, err
	}
	graphLinkModel.Fields = graphParsedFields
	// output: validate & normalize & defaults

	panic("implement me")
}

func NewCaseLinkService(app *App) (*CaseLinkService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseLinkService{app: app}, nil
}
