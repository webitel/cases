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

type CloseReasonGroupService struct {
	app *App
}

func (s CloseReasonGroupService) CreateCloseReasonGroup(ctx context.Context, req *_go.CreateCloseReasonGroupRequest) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("close_reason_group_service.create_close_reason_group.name.required", "Lookup name is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("close_reason_group_service.create_close_reason_group.authorization.failed", err.Error())
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
	lookup := &_go.CloseReasonGroup{
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

	// Create the close reason group in the store
	l, e := s.app.Store.CloseReasonGroup().Create(&createOpts, lookup)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.create_close_reason_group.store.create.failed", e.Error())
	}

	return l, nil
}

func (s CloseReasonGroupService) ListCloseReasonGroups(ctx context.Context, req *_go.ListCloseReasonGroupsRequest) (*_go.CloseReasonGroupList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("close_reason_group_service.list_close_reason_groups.authorization.failed", err.Error())
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
		IDs:     req.Id,
		Session: session,
		Fields:  fields,
		Context: ctx,
		Page:    int(page),
		Sort:    req.Sort,
		Size:    int(req.Size),
		Time:    t,
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	lookups, e := s.app.Store.CloseReasonGroup().List(&searchOptions)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.list_close_reason_groups.store.list.failed", e.Error())
	}

	return lookups, nil
}

func (s CloseReasonGroupService) UpdateCloseReasonGroup(ctx context.Context, req *_go.UpdateCloseReasonGroupRequest) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_group_service.update_close_reason_group.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("close_reason_group_service.update_close_reason_group.authorization.failed", err.Error())
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
	lookup := &_go.CloseReasonGroup{
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
				return nil, cerror.NewBadRequestError("close_reason_group_service.update_close_reason_group.name.required", "Lookup name is required and cannot be empty")
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
	l, e := s.app.Store.CloseReasonGroup().Update(&updateOpts, lookup)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.update_close_reason_group.store.update.failed", e.Error())
	}

	return l, nil
}

func (s CloseReasonGroupService) DeleteCloseReasonGroup(ctx context.Context, req *_go.DeleteCloseReasonGroupRequest) (*_go.CloseReasonGroup, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_group_service.delete_close_reason_group.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("close_reason_group_service.delete_close_reason_group.authorization.failed", err.Error())
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
	e := s.app.Store.CloseReasonGroup().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.delete_close_reason_group.store.delete.failed", e.Error())
	}

	return &(_go.CloseReasonGroup{Id: req.Id}), nil
}

func (s CloseReasonGroupService) LocateCloseReasonGroup(ctx context.Context, req *_go.LocateCloseReasonGroupRequest) (*_go.LocateCloseReasonGroupResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("close_reason_group_service.locate_close_reason_group.id.required", "Lookup ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListCloseReasonGroupsRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1, // We only need one item
	}

	// Call the ListCloseReasonGroups method
	listResp, err := s.ListCloseReasonGroups(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.locate_close_reason_group.list_close_reason_groups.error", err.Error())
	}

	// Check if the close reason group was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_group_service.locate_close_reason_group.not_found", "Close reason group not found")
	}

	// Return the found close reason group
	return &_go.LocateCloseReasonGroupResponse{CloseReasonGroup: listResp.Items[0]}, nil
}

func NewCloseReasonGroupsService(app *App) (*CloseReasonGroupService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_close_reason_group_service.args_check.app_nil", "internal is nil")
	}

	return &CloseReasonGroupService{app: app}, nil
}
