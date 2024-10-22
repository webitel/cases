package app

import (
	"context"
	"strings"
	"time"

	_go "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type ReasonService struct {
	app *App
}

// CreateReason implements api.ReasonsServer.
func (s *ReasonService) CreateReason(ctx context.Context, req *_go.CreateReasonRequest) (*_go.Reason, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("reason_service.create_reason.name.required", "Reason name is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("reason_service.create_reason.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the creator and updater
	currentU := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new reason model
	reason := &_go.Reason{
		Name:          req.Name,
		Description:   req.Description,
		CreatedBy:     currentU,
		UpdatedBy:     currentU,
		CloseReasonId: req.CloseReasonId,
	}

	fields := []string{"id", "lookup_id", "name", "description", "created_at", "updated_at", "created_by", "updated_by"}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Create the reason in the store
	r, e := s.app.Store.Reason().Create(&createOpts, reason)
	if e != nil {
		return nil, model.NewInternalError("reason_service.create_reason.store.create.failed", e.Error())
	}

	return r, nil
}

// ListReasons implements api.ReasonsServer.
func (s *ReasonService) ListReasons(ctx context.Context, req *_go.ListReasonRequest) (*_go.ReasonList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("reason_service.list_reasons.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

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
	searchOptions := model.SearchOptions{
		IDs:     req.Id,
		Session: session,
		Fields:  fields,
		Context: ctx,
		Sort:    req.Sort,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	reasons, e := s.app.Store.Reason().List(&searchOptions, req.CloseReasonId)
	if e != nil {
		return nil, model.NewInternalError("reason_service.list_reasons.store.list.failed", e.Error())
	}

	return reasons, nil
}

// UpdateReason implements api.ReasonsServer.
func (s *ReasonService) UpdateReason(ctx context.Context, req *_go.UpdateReasonRequest) (*_go.Reason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("reason_service.update_reason.id.required", "Reason ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("reason_service.update_reason.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the updater
	u := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update reason model
	reason := &_go.Reason{
		Id:            req.Id,
		CloseReasonId: req.CloseReasonId,
		Name:          req.Input.Name,
		Description:   req.Input.Description,
		UpdatedBy:     u,
	}

	fields := []string{"id", "lookup_id"}

	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, model.NewBadRequestError("reason_service.update_reason.name.required", "Reason name is required and cannot be empty")
			}
		case "description":
			fields = append(fields, "description")
		}
	}

	t := time.Now()

	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Update the reason in the store
	r, e := s.app.Store.Reason().Update(&updateOpts, reason)
	if e != nil {
		return nil, model.NewInternalError("reason_service.update_reason.store.update.failed", e.Error())
	}

	return r, nil
}

// DeleteReason implements api.ReasonsServer.
func (s *ReasonService) DeleteReason(ctx context.Context, req *_go.DeleteReasonRequest) (*_go.Reason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("reason_service.delete_reason.id.required", "Reason ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("reason_service.delete_reason.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	t := time.Now()
	// Define delete options
	deleteOpts := model.DeleteOptions{
		Session: session,
		Context: ctx,
		IDs:     []int64{req.Id},
		Time:    t,
	}

	// Delete the reason in the store
	e := s.app.Store.Reason().Delete(&deleteOpts, req.CloseReasonId)
	if e != nil {
		return nil, model.NewInternalError("reason_service.delete_reason.store.delete.failed", e.Error())
	}

	return &(_go.Reason{Id: req.Id}), nil
}

// LocateReason implements api.ReasonsServer.
func (s *ReasonService) LocateReason(ctx context.Context, req *_go.LocateReasonRequest) (*_go.LocateReasonResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("reason_service.locate_reason.id.required", "Reason ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListReasonRequest{
		Id:            []int64{req.Id},
		Fields:        req.Fields,
		Page:          1,
		Size:          1, // We only need one item
		CloseReasonId: req.CloseReasonId,
	}

	// Call the ListReasons method
	listResp, err := s.ListReasons(ctx, listReq)
	if err != nil {
		return nil, model.NewInternalError("reason_service.locate_reason.list_reasons.error", err.Error())
	}

	// Check if the reason was found
	if len(listResp.Items) == 0 {
		return nil, model.NewNotFoundError("reason_service.locate_reason.not_found", "Reason not found")
	}

	// Return the found reason
	return &_go.LocateReasonResponse{Reason: listResp.Items[0]}, nil
}

func NewReasonService(app *App) (*ReasonService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_reason_service.args_check.app_nil", "internal is nil")
	}
	return &ReasonService{app: app}, nil
}
