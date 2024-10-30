package app

import (
	"context"

	"github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
)

type CaseLinkService struct {
	app *App
}

func (c *CaseLinkService) LocateLink(ctx context.Context, request *cases.LocateLinkRequest) (*cases.CaseLink, error) {
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
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) MergeLinks(ctx context.Context, request *cases.MergeLinksRequest) (*cases.CaseLinkList, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseLinkService) ResetLinks(ctx context.Context, request *cases.ResetLinksRequest) (*cases.CaseLinkList, error) {
	// TODO implement me
	panic("implement me")
}

func NewCaseLinkService(app *App) (*CaseLinkService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_comment_service.check_args.app", "unable to init service, app is nil")
	}
	return &CaseLinkService{app: app}, nil
}
