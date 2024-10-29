package app

import (
	"context"

	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
)

// parallel processes
//
// etag formation
// authorization interceptor
// additional auth layer with context attributes, client name and service registry
// OpenTelemetry
// graphQL declaration
// proto filter parsing

type CaseService struct {
	app *App
}

func (c *CaseService) SearchCases(ctx context.Context, request *cases.SearchCasesRequest) (*cases.CaseList, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseService) LocateCase(ctx context.Context, request *cases.LocateCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseService) CreateCase(ctx context.Context, request *cases.CreateCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseService) UpdateCase(ctx context.Context, request *cases.UpdateCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

func (c *CaseService) DeleteCase(ctx context.Context, request *cases.DeleteCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

func NewCaseService(app *App) (*CaseService, model.AppError) {
	if app == nil {
		return nil, model.NewBadRequestError("app.case.new_case_service.check_args.app", "unable to init case service, app is nil")
	}
	return &CaseService{app: app}, nil
}
