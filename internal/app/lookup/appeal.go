package lookup

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	"context"
	"github.com/webitel/cases/internal/app"
	"github.com/webitel/cases/model"
)

type AppealLookupService struct {
	app *app.App
}

func (a AppealLookupService) ListAppealLookups(ctx context.Context, request *_go.ListAppealLookupsRequest) (*_go.AppealLookupList, error) {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookupService) CreateAppealLookup(ctx context.Context, request *_go.CreateAppealLookupRequest) (*_go.AppealLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookupService) UpdateAppealLookup(ctx context.Context, request *_go.UpdateAppealLookupRequest) (*_go.AppealLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookupService) DeleteAppealLookup(ctx context.Context, request *_go.DeleteAppealLookupRequest) (*_go.AppealLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (a AppealLookupService) LocateAppealLookup(ctx context.Context, request *_go.LocateAppealLookupRequest) (*_go.LocateAppealLookupResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewAppealLookupService(app *app.App) (*AppealLookupService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_appeal_lookup_service.args_check.app_nil", "pkg is nil")
	}
	return &AppealLookupService{app: app}, nil
}
