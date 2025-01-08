package app

import (
	"context"
	"strings"
	"time"

	_go "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"

	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
)

type StatusService struct {
	app *App
	_go.UnimplementedStatusesServer
}

const (
	ErrLookupNameReq = "Lookup name is required"
	defaultFields    = "id, name, description"
)

// CreateStatus implements api.StatusesServer.
func (s StatusService) CreateStatus(ctx context.Context, req *_go.CreateStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("status.create_status.lookup.name.required", ErrLookupNameReq)
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("status.create_status.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the creator and updater
	currentU := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new lookup model
	lookup := &_go.Status{
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

	// Create the status in the store
	l, e := s.app.Store.Status().Create(&createOpts, lookup)
	if e != nil {
		return nil, cerror.NewInternalError("status.create_status.store.create.failed", e.Error())
	}

	return l, nil
}

// ListStatuses implements api.StatusesServer.
func (s StatusService) ListStatuses(ctx context.Context, req *_go.ListStatusRequest) (*_go.StatusList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("status.list_status.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
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
		IDs: req.Id,
		//Session: session,
		Fields:  fields,
		Context: ctx,
		Page:    int(page),
		Sort:    req.Sort,
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
		Auth:    model.NewSessionAuthOptions(session, "dictionaries"),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	lookups, e := s.app.Store.Status().List(&searchOptions)
	if e != nil {
		return nil, cerror.NewInternalError("status.list_status.store.list.failed", e.Error())
	}

	return lookups, nil
}

// UpdateStatus implements api.StatusesServer.
func (s StatusService) UpdateStatus(ctx context.Context, req *_go.UpdateStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status.update_status.lookup.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("status.update_status.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the updater
	currentU := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update lookup model
	lookup := &_go.Status{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "updated_at", "updated_by"}

	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("status.update_status.name.required", "Lookup name is required and cannot be empty")
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

	// Update the lookup in the store
	l, e := s.app.Store.Status().Update(&updateOpts, lookup)
	if e != nil {
		return nil, cerror.NewInternalError("status.update_status.store.update.failed", e.Error())
	}

	return l, nil
}

// DeleteStatus implements api.StatusesServer.
func (s StatusService) DeleteStatus(ctx context.Context, req *_go.DeleteStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status.delete_status.lookup.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("status.delete_status.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
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
	e := s.app.Store.Status().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("status.delete_status.store.delete.failed", e.Error())
	}

	return &(_go.Status{Id: req.Id}), nil
}

// LocateStatus implements api.StatusesServer.
func (s StatusService) LocateStatus(ctx context.Context, req *_go.LocateStatusRequest) (*_go.LocateStatusResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status.locate_status.lookup.id.required", "Lookup ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListStatusRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1, // We only need one item
	}

	// Call the ListStatuses method
	listResp, err := s.ListStatuses(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("status.locate_status.list_status.error", err.Error())
	}

	// Check if the lookup was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("status.locate_status.not_found", "Status lookup not found")
	}

	// Return the found status lookup
	return &_go.LocateStatusResponse{Status: listResp.Items[0]}, nil
}

func NewStatusService(app *App) (*StatusService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_status.args_check.app_nil", "internal is nil")
	}
	return &StatusService{app: app}, nil
}
