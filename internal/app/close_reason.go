package app

import (
	"context"
	"strings"
	"time"

	_go "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type CloseReasonService struct {
	app *App
}

func (s CloseReasonService) CreateCloseReason(ctx context.Context, req *_go.CreateCloseReasonRequest) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("close_reason_service.create_close_reason.name.required", "Lookup name is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("close_reason_service.create_close_reason.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the creator and updater
	currentU := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new lookup model
	lookup := &_go.CloseReason{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "name", "description", "created_at", "updated_at", "created_by", "updated_by"}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Create the close reason in the store
	l, e := s.app.Store.CloseReason().Create(&createOpts, lookup)
	if e != nil {
		return nil, model.NewInternalError("close_reason_service.create_close_reason.store.create.failed", e.Error())
	}

	return l, nil
}

func (s CloseReasonService) ListCloseReasons(ctx context.Context, req *_go.ListCloseReasonRequest) (*_go.CloseReasonList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("close_reason_service.list_close_reasons.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFields, ", ")
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
		Page:    int(page),
		Sort:    req.Sort,
		Size:    int(req.Size),
		Filter:  make(map[string]interface{}),
		Time:    t,
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	} else if req.Name != "" {
		searchOptions.Filter["name"] = req.Name
	}

	lookups, e := s.app.Store.CloseReason().List(&searchOptions)
	if e != nil {
		return nil, model.NewInternalError("close_reason_service.list_close_reasons.store.list.failed", e.Error())
	}

	return lookups, nil
}

func (s CloseReasonService) UpdateCloseReason(ctx context.Context, req *_go.UpdateCloseReasonRequest) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("close_reason_service.update_close_reason.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("close_reason_service.update_close_reason.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the updater
	currentU := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update lookup model
	lookup := &_go.CloseReason{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "updated_at", "updated_by"}

	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
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

	// Update the lookup in the store
	l, e := s.app.Store.CloseReason().Update(&updateOpts, lookup)
	if e != nil {
		return nil, model.NewInternalError("close_reason_service.update_close_reason.store.update.failed", e.Error())
	}

	return l, nil
}

func (s CloseReasonService) DeleteCloseReason(ctx context.Context, req *_go.DeleteCloseReasonRequest) (*_go.CloseReason, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("close_reason_service.delete_close_reason.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("close_reason_service.delete_close_reason.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasAccess(scope, accessMode) {
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

	// Delete the lookup in the store
	e := s.app.Store.CloseReason().Delete(&deleteOpts)
	if e != nil {
		return nil, model.NewInternalError("close_reason_service.delete_close_reason.store.delete.failed", e.Error())
	}

	return &(_go.CloseReason{Id: req.Id}), nil
}

func (s CloseReasonService) LocateCloseReason(ctx context.Context, req *_go.LocateCloseReasonRequest) (*_go.LocateCloseReasonResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("close_reason_service.locate_close_reason.id.required", "Lookup ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListCloseReasonRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1, // We only need one item
	}

	// Call the ListCloseReasons method
	listResp, err := s.ListCloseReasons(ctx, listReq)
	if err != nil {
		return nil, model.NewInternalError("close_reason_service.locate_close_reason.list_close_reasons.error", err.Error())
	}

	// Check if the close reason was found
	if len(listResp.Items) == 0 {
		return nil, model.NewNotFoundError("close_reason_service.locate_close_reason.not_found", "Close reason not found")
	}

	// Return the found close reason
	return &_go.LocateCloseReasonResponse{CloseReason: listResp.Items[0]}, nil
}

func NewCloseReasonService(app *App) (*CloseReasonService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_close_reason_service.args_check.app_nil",
			"internal is nil")
	}
	return &CloseReasonService{app: app}, nil
}
