package app

import (
	"context"
	"strings"
	"time"

	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
)

type CloseReasonService struct {
	app *App
	_go.UnimplementedCloseReasonsServer
	objClassName string
}

// CreateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) CreateCloseReason(ctx context.Context, req *_go.CreateCloseReasonRequest) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("close_reason_service.create_close_reason.name.required", "Close reason name is required")
	}

	fields := []string{"id", "lookup_id", "name", "description", "created_at", "updated_at", "created_by", "updated_by"}

	t := time.Now()

	// Define create options
	createOpts := &model.CreateOptions{
		Context: ctx,
		Fields:  fields,
		Time:    t,
		Auth:    model.GetAutherOutOfContext(ctx),
	}
	// Define the current user as the creator and updater
	currentU := &_go.Lookup{
		Id: createOpts.GetAuthOpts().GetUserId(),
	}
	// Create a new close reason user_auth
	closeReason := &_go.CloseReason{
		Name:               req.Name,
		Description:        req.Description,
		CreatedBy:          currentU,
		UpdatedBy:          currentU,
		CloseReasonGroupId: req.CloseReasonGroupId,
	}

	// Create the close reason in the store
	r, e := s.app.Store.CloseReason().Create(createOpts, closeReason)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_service.create_close_reason.store.create.failed", e.Error())
	}

	return r, nil
}

// ListCloseReasons implements api.CloseReasonsServer.
func (s *CloseReasonService) ListCloseReasons(ctx context.Context, req *_go.ListCloseReasonRequest) (*_go.CloseReasonList, error) {

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsStatus, ", ")
	}

	// Use default page size and page number if not provided
	page := req.Page
	if page == 0 {
		page = 1
	}

	t := time.Now()
	searchOptions := &model.SearchOptions{
		IDs: req.Id,
		//UserAuthSession: session,
		Fields:  fields,
		Context: ctx,
		Sort:    req.Sort,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
		Auth:    model.GetAutherOutOfContext(ctx),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	closeReasons, e := s.app.Store.CloseReason().List(searchOptions, req.CloseReasonGroupId)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_service.list_close_reasons.store.list.failed", e.Error())
	}

	return closeReasons, nil
}

// UpdateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) UpdateCloseReason(ctx context.Context, req *_go.UpdateCloseReasonRequest) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.update_close_reason.id.required", "Close reason ID is required")
	}

	fields := []string{"id", "lookup_id"}

	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("close_reason_service.update_close_reason.name.required", "Close reason name is required and cannot be empty")
			}
		case "description":
			fields = append(fields, "description")
		}
	}

	t := time.Now()

	// Define update options
	updateOpts := &model.UpdateOptions{
		Context: ctx,
		Fields:  fields,
		Time:    t,
		Auth:    model.GetAutherOutOfContext(ctx),
	}
	// Define the current user as the updater
	u := &_go.Lookup{
		Id: updateOpts.GetAuthOpts().GetUserId(),
	}
	// Update close reason user_auth
	closeReason := &_go.CloseReason{
		Id:                 req.Id,
		CloseReasonGroupId: req.CloseReasonGroupId,
		Name:               req.Input.Name,
		Description:        req.Input.Description,
		UpdatedBy:          u,
	}

	// Update the close reason in the store
	r, e := s.app.Store.CloseReason().Update(updateOpts, closeReason)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_service.update_close_reason.store.update.failed", e.Error())
	}

	return r, nil
}

// DeleteCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) DeleteCloseReason(ctx context.Context, req *_go.DeleteCloseReasonRequest) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.delete_close_reason.id.required", "Close reason ID is required")
	}

	t := time.Now()
	// Define delete options
	deleteOpts := &model.DeleteOptions{
		Context: ctx,
		IDs:     []int64{req.Id},
		Time:    t,
		Auth:    model.GetAutherOutOfContext(ctx),
	}

	// Delete the close reason in the store
	e := s.app.Store.CloseReason().Delete(deleteOpts, req.CloseReasonGroupId)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_service.delete_close_reason.store.delete.failed", e.Error())
	}

	return &(_go.CloseReason{Id: req.Id}), nil
}

// LocateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService) LocateCloseReason(ctx context.Context, req *_go.LocateCloseReasonRequest) (*_go.LocateCloseReasonResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.locate_close_reason.id.required", "Close reason ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListCloseReasonRequest{
		Id:                 []int64{req.Id},
		Fields:             req.Fields,
		Page:               1,
		Size:               1, // We only need one item
		CloseReasonGroupId: req.GetCloseReasonGroupId(),
	}

	// Call the ListCloseReasons method
	listResp, err := s.ListCloseReasons(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.locate_close_reason.list_close_reasons.error", err.Error())
	}

	// Check if the close reason was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_service.locate_close_reason.not_found", "Close reason not found")
	}

	// Return the found close reason
	return &_go.LocateCloseReasonResponse{CloseReason: listResp.Items[0]}, nil
}

func NewCloseReasonService(app *App) (*CloseReasonService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_close_reason_service.args_check.app_nil", "internal is nil")
	}
	return &CloseReasonService{app: app, objClassName: model.ScopeDictionary}, nil
}
