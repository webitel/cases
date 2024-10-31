package app

import (
	"context"

	"github.com/webitel/cases/api/cases"

	cerror "github.com/webitel/cases/internal/error"
)

/*

API layer

- etag formation
- graphQL types declaration and validation
- authorization interceptor
- rabbitMQ events
- case name forming
-------------------------------------------------
- OpenTelemetry interceptors and init
- storage interfaces
- additional auth layer with context attributes, client name and service registry


Database layer

proto filter parsing
storages (singleton)
calendar storage and calculation module
sql scripts structure
*/

type CaseService struct {
	app *App
}

/*
SearchCases
Authorization
Obac
Rbac
Fields validation with graph
Search options construction with filters
Database layer with search options
Result construction by fields requested
*/
func (c *CaseService) SearchCases(ctx context.Context, request *cases.SearchCasesRequest) (*cases.CaseList, error) {
	// TODO implement me
	panic("implement me")
}

/*
LocateCase
Authorization
Obac
Rbac
Etag parsing
Fields validation with graph
Search options construction with filters
Database layer with search options
Result construction with etag
*/
func (c *CaseService) LocateCase(ctx context.Context, request *cases.LocateCaseRequest) (*cases.Case, error) {
	// TODO implement me
	panic("implement me")
}

/*
CreateCase
Authorization
Obac
Fields validation with graph
Database layer with create options
Calendar's logic
Result construction
Rabbit event publishing
*/
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

func NewCaseService(app *App) (*CaseService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewBadRequestError("app.case.new_case_service.check_args.app", "unable to init case service, app is nil")
	}
	return &CaseService{app: app}, nil
}
