package lookup

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	"context"
	"github.com/webitel/cases/app"
	"github.com/webitel/cases/model"
)

type CloseReasonLookupService struct {
	app *app.App
}

func (c CloseReasonLookupService) ListCloseReasonLookups(ctx context.Context, request *_go.ListCloseReasonLookupsRequest) (*_go.CloseReasonLookupList, error) {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookupService) CreateCloseReasonLookup(ctx context.Context, request *_go.CreateCloseReasonLookupRequest) (*_go.CloseReasonLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookupService) UpdateCloseReasonLookup(ctx context.Context, request *_go.UpdateCloseReasonLookupRequest) (*_go.CloseReasonLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookupService) DeleteCloseReasonLookup(ctx context.Context, request *_go.DeleteCloseReasonLookupRequest) (*_go.CloseReasonLookup, error) {
	//TODO implement me
	panic("implement me")
}

func (c CloseReasonLookupService) LocateCloseReasonLookup(ctx context.Context, request *_go.LocateCloseReasonLookupRequest) (*_go.LocateCloseReasonLookupResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewCloseReasonLookupService(app *app.App) (*CloseReasonLookupService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_close_reason_lookup_service.args_check.app_nil",
			"app is nil")
	}
	return &CloseReasonLookupService{app: app}, nil
}
