package app

import (
	"context"
	"strings"

	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"log/slog"
)

const (
	defaultFieldsSource = "id, name, description, type, created_by"
)

type SourceService struct {
	app *App
	_go.UnimplementedSourcesServer
	objClassName string
}

var SourceMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
	{"type", true},
})

// CreateSource implements api.SourcesServer.
func (s *SourceService) CreateSource(
	ctx context.Context,
	req *_go.CreateSourceRequest,
) (*_go.Source, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("source_service.create_source.name.required", "Source name is required")
	}
	if req.Input.Type == _go.SourceType_TYPE_UNSPECIFIED {
		return nil, cerror.NewBadRequestError("source_service.create_source.type.required", "Source type is required")
	}

	createOpts, err := model.NewCreateOptions(ctx, req, SourceMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}

	input := &_go.Source{
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Type:        req.Input.Type,
	}

	// Create the source in the store
	res, err := s.app.Store.Source().Create(createOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.create_source.store.create.failed", err.Error())
	}

	return res, nil
}

// ListSources implements api.SourcesServer.
func (s *SourceService) ListSources(
	ctx context.Context,
	req *_go.ListSourceRequest,
) (*_go.SourceList, error) {

	searchOpts, err := model.NewSearchOptions(ctx, req, SourceMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}
	searchOpts.IDs = req.Id
	searchOpts.Filter = make(map[string]any)

	if req.Q != "" {
		searchOpts.Filter["name"] = req.Q
	}

	if len(req.Type) > 0 {
		searchOpts.Filter["type"] = req.Type
	}

	res, err := s.app.Store.Source().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.list_sources.store.list.failed", err.Error())
	}

	return res, nil
}

// UpdateSource implements api.SourcesServer.
func (s *SourceService) UpdateSource(
	ctx context.Context,
	req *_go.UpdateSourceRequest,
) (*_go.Source, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.update_source.id.required", "Source ID is required")
	}

	updateOpts, err := model.NewUpdateOptions(ctx, req, SourceMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}

	input := &_go.Source{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Type:        req.Input.Type,
	}

	// Update the source in the store
	res, err := s.app.Store.Source().Update(updateOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.update_source.store.update.failed", err.Error())
	}

	return res, nil
}

// DeleteSource implements api.SourcesServer.
func (s *SourceService) DeleteSource(
	ctx context.Context,
	req *_go.DeleteSourceRequest,
) (*_go.Source, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.delete_source.id.required", "Source ID is required")
	}

	deleteOpts, err := model.NewDeleteOptions(ctx, SourceMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}

	deleteOpts.IDs = []int64{req.Id}

	// Delete the source in the store
	err = s.app.Store.Source().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.delete_source.store.delete.failed", err.Error())
	}

	return &(_go.Source{Id: req.Id}), nil
}

// LocateSource implements api.SourcesServer.
func (s *SourceService) LocateSource(
	ctx context.Context,
	req *_go.LocateSourceRequest,
) (*_go.LocateSourceResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.locate_source.id.required", "Source ID is required")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsSource, ", ")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListSourceRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1,
	}

	res, err := s.ListSources(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.locate_source.list_sources.error", err.Error())
	}

	if len(res.Items) == 0 {
		return nil, cerror.NewNotFoundError("source_service.locate_source.not_found", "Source not found")
	}

	// Return the found source
	return &_go.LocateSourceResponse{Source: res.Items[0]}, nil
}

func NewSourceService(app *App) (*SourceService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_source_service.args_check.app_nil", "internal is nil")
	}
	return &SourceService{app: app, objClassName: model.ScopeDictionary}, nil
}
