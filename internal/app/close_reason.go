package app

import (
	"context"
	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
	"log/slog"
	"strings"
)

type CloseReasonService struct {
	app *App
	_go.UnimplementedCloseReasonsServer
	objClassName string
}

var CloseReasonMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
	{"close_reason_id", false},
})

const (
	defaultFieldsCloseReason = "id, name, description, created_by"
)

// CreateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) CreateCloseReason(
	ctx context.Context,
	req *_go.CreateCloseReasonRequest,
) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("close_reason_service.create_close_reason.name.required", "Close reason name is required")
	}

	createOpts, err := model.NewCreateOptions(ctx, req, CloseReasonMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}

	input := &_go.CloseReason{
		Name:               req.Input.Name,
		Description:        req.Input.Description,
		CloseReasonGroupId: req.CloseReasonGroupId,
	}

	// Create the close reason in the store
	res, err := s.app.Store.CloseReason().Create(createOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.create_close_reason.store.create.failed", err.Error())
	}

	return res, nil
}

// ListCloseReasons implements api.CloseReasonsServer.
func (s *CloseReasonService) ListCloseReasons(
	ctx context.Context,
	req *_go.ListCloseReasonRequest,
) (*_go.CloseReasonList, error) {

	searchOpts, err := model.NewSearchOptions(ctx, req, CloseReasonMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}
	searchOpts.IDs = req.Id
	searchOpts.Filter = make(map[string]any)

	if req.Q != "" {
		searchOpts.Filter["name"] = req.Q
	}

	res, err := s.app.Store.CloseReason().List(searchOpts, req.CloseReasonGroupId)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.list_close_reasons.store.list.failed", err.Error())
	}

	return res, nil
}

// UpdateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) UpdateCloseReason(
	ctx context.Context,
	req *_go.UpdateCloseReasonRequest,
) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.update_close_reason.id.required", "Close reason ID is required")
	}

	updateOpts, err := model.NewUpdateOptions(ctx, req, CloseReasonMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}

	// Update close reason user_auth
	input := &_go.CloseReason{
		Id:                 req.Id,
		CloseReasonGroupId: req.CloseReasonGroupId,
		Name:               req.Input.Name,
		Description:        req.Input.Description,
	}

	// Update the close reason in the store
	res, err := s.app.Store.CloseReason().Update(updateOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.update_close_reason.store.update.failed", err.Error())
	}

	return res, nil
}

// DeleteCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) DeleteCloseReason(
	ctx context.Context,
	req *_go.DeleteCloseReasonRequest,
) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.delete_close_reason.id.required", "Close reason ID is required")
	}

	deleteOpts, err := model.NewDeleteOptions(ctx, CloseReasonMetadata)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return nil, InternalError
	}

	deleteOpts.IDs = []int64{req.Id}

	// Delete the close reason in the store
	err = s.app.Store.CloseReason().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.delete_close_reason.store.delete.failed", err.Error())
	}

	return &(_go.CloseReason{Id: req.Id}), nil
}

// LocateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) LocateCloseReason(
	ctx context.Context,
	req *_go.LocateCloseReasonRequest,
) (*_go.LocateCloseReasonResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.locate_close_reason.id.required", "Close reason ID is required")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)
	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsCloseReason, ", ")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListCloseReasonRequest{
		Id:                 []int64{req.Id},
		Fields:             req.Fields,
		Page:               1,
		Size:               1,
		CloseReasonGroupId: req.GetCloseReasonGroupId(),
	}

	// Call the ListCloseReasons method
	res, err := s.ListCloseReasons(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.locate_close_reason.list_close_reasons.error", err.Error())
	}

	// Check if the close reason was found
	if len(res.Items) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_service.locate_close_reason.not_found", "Close reason not found")
	}

	// Return the found close reason
	return &_go.LocateCloseReasonResponse{CloseReason: res.Items[0]}, nil
}

func NewCloseReasonService(app *App) (*CloseReasonService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_close_reason_service.args_check.app_nil", "internal is nil")
	}
	return &CloseReasonService{app: app, objClassName: model.ScopeDictionary}, nil
}
