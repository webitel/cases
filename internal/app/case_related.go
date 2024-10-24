package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
)

type RelatedCaseService struct {
	app *App
}

func (r *RelatedCaseService) LocateRelatedCase(ctx context.Context, request *cases.LocateRelatedCaseRequest) (*cases.RelatedCase, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) UpdateRelatedCase(ctx context.Context, request *cases.UpdateRelatedCaseRequest) (*cases.RelatedCase, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) DeleteRelatedCase(ctx context.Context, request *cases.DeleteRelatedCaseRequest) (*cases.RelatedCase, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) ListRelatedCases(ctx context.Context, request *cases.ListRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) MergeRelatedCases(ctx context.Context, request *cases.MergeRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RelatedCaseService) ResetRelatedCases(ctx context.Context, request *cases.ResetRelatedCasesRequest) (*cases.RelatedCaseList, error) {
	//TODO implement me
	panic("implement me")
}

func NewCaseRelatedService(app *App) (*RelatedCaseService, model.AppError) {
	if app == nil {
		return nil, model.NewBadRequestError("app.case.new_case_related_service.check_args.app", "unable to init service, app is nil")
	}
	return &RelatedCaseService{app: app}, nil
}
