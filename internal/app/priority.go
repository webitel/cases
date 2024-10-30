package app

import (
	"context"
	"strings"
	"time"

	api "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"

	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
)

type PriorityService struct {
	app *App
}

// CreatePriority implements api.PrioritiesServer.
func (p *PriorityService) CreatePriority(ctx context.Context, req *api.CreatePriorityRequest) (*api.Priority, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("priority_service.create_priority.name.required", "Lookup name is required")
	}

	if req.Color == "" {
		return nil, cerror.NewBadRequestError("priority_service.create_priority.color.required", "Color is required")
	}

	session, err := p.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("priority_service.create_priority.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the creator and updater
	currentU := &api.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new lookup model
	lookup := &api.Priority{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
		Color:       req.Color,
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

	// Create the priority in the store
	l, e := p.app.Store.Priority().Create(&createOpts, lookup)
	if e != nil {
		return nil, cerror.NewInternalError("priority_service.create_priority.store.create.failed", e.Error())
	}

	return l, nil
}

// ListPriorities implements api.PrioritiesServer.
func (p *PriorityService) ListPriorities(ctx context.Context, req *api.ListPriorityRequest) (*api.PriorityList, error) {
	session, err := p.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("priority.list_priorities.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
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

	prios, e := p.app.Store.Priority().List(&searchOptions)
	if e != nil {
		return nil, cerror.NewInternalError("priority.list_priority.store.list.failed", e.Error())
	}

	return prios, nil
}

// UpdatePriority implements api.PrioritiesServer.
func (p *PriorityService) UpdatePriority(ctx context.Context, req *api.UpdatePriorityRequest) (*api.Priority, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("priority_service.update_priority.id.required", "Lookup ID is required")
	}

	session, err := p.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("priority_service.update_priority.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the updater
	currentU := &api.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update lookup model
	lookup := &api.Priority{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		UpdatedBy:   currentU,
		Color:       req.Input.Color,
	}

	fields := []string{"id", "updated_at", "updated_by"}

	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("priority_service.update_priority.name.required", "Lookup name is required and cannot be empty")
			}
		case "description":
			fields = append(fields, "description")
		case "color":
			fields = append(fields, "color")
			if req.Input.Color == "" {
				return nil, cerror.NewBadRequestError("priority_service.update_priority.color.required", "Color is required")
			}
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
	l, e := p.app.Store.Priority().Update(&updateOpts, lookup)
	if e != nil {
		return nil, cerror.NewInternalError("priority_service.update_priority.store.update.failed", e.Error())
	}

	return l, nil
}

// DeletePriority implements api.PrioritiesServer.
func (p *PriorityService) DeletePriority(ctx context.Context, req *api.DeletePriorityRequest) (*api.Priority, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("priority.delete_priority.id.required", "Priority ID is required")
	}

	session, err := p.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("priority.delete_priority.authorization.failed", err.Error())
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

	// Delete the priority in the store
	e := p.app.Store.Priority().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("priority.delete_priority.store.delete.failed", e.Error())
	}

	return &api.Priority{Id: req.Id}, nil
}

// LocatePriority implements api.PrioritiesServer.
func (p *PriorityService) LocatePriority(ctx context.Context, req *api.LocatePriorityRequest) (*api.LocatePriorityResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("priority_service.locate_priority.id.required", "Lookup ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &api.ListPriorityRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1, // We only need one item
	}

	// Call the ListPriorities method
	listResp, err := p.ListPriorities(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("priority_service.locate_priority.list_priorities.error", err.Error())
	}

	// Check if the priority was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("priority_service.locate_priority.not_found", "Priority not found")
	}

	// Return the found priority
	return &api.LocatePriorityResponse{Priority: listResp.Items[0]}, nil
}

func NewPriorityService(app *App) (*PriorityService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_priority_service.args_check.app_nil", "internal is nil")
	}
	return &PriorityService{app: app}, nil
}
