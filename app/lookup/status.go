package lookup

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	"context"
	"github.com/webitel/cases/app"
	"github.com/webitel/cases/model"
)

type StatusLookupService struct {
	app *app.App
}

func (s StatusLookupService) ListStatusLookups(ctx context.Context, request *_go.ListStatusLookupsRequest) (*_go.StatusLookupList, error) {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookupService) CreateStatusLookup(ctx context.Context, request *_go.CreateStatusLookupRequest) (*_go.StatusLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookupService) UpdateStatusLookup(ctx context.Context, request *_go.UpdateStatusLookupRequest) (*_go.StatusLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookupService) DeleteStatusLookup(ctx context.Context, request *_go.DeleteStatusLookupRequest) (*_go.StatusLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (s StatusLookupService) LocateStatusLookup(ctx context.Context, request *_go.LocateStatusLookupRequest) (*_go.LocateStatusLookupResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewStatusLookupService(app *app.App) (*StatusLookupService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_status_lookup_service.args_check.app_nil",
			"app is nil")
	}
	return &StatusLookupService{app: app}, nil
}
