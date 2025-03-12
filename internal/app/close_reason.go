package app

import (
	"context"
	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
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

// CreateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) CreateCloseReason(
	ctx context.Context,
	req *_go.CreateCloseReasonRequest,
) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("close_reason_service.create_close_reason.name.required", "Close reason name is required")
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CloseReasonMetadata),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
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

	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CloseReasonMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	searchOpts.AddFilter("name", req.Q)

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

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CloseReasonMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
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

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

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
