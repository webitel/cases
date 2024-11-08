package app

import (
	"context"

	cases "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
)

type RelatedCaseService struct {
	app *App
	cases.UnimplementedRelatedCasesServer
}

func (r *RelatedCaseService) LocateRelatedCase(ctx context.Context, request *cases.LocateRelatedCaseRequest) (*cases.RelatedCase, error) {
	// TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) CreateRelatedCase(ctx context.Context, request *cases.CreateRelatedCaseRequest) (*cases.RelatedCase, error) {
	// TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) UpdateRelatedCase(ctx context.Context, request *cases.UpdateRelatedCaseRequest) (*cases.RelatedCase, error) {
	// TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) DeleteRelatedCase(ctx context.Context, request *cases.DeleteRelatedCaseRequest) (*cases.RelatedCase, error) {
	// TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) ListRelatedCases(ctx context.Context, request *cases.ListRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	// TODO implement me
	panic("implement me")
}

func NewCaseRelatedService(app *App) (*RelatedCaseService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_related_service.check_args.app", "unable to init service, app is nil")
	}
	return &RelatedCaseService{app: app}, nil
}
