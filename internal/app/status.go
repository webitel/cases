package app

import (
	"context"
	"strings"
	"time"

	_go "github.com/webitel/cases/api"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type StatusService struct {
	app *App
}

const (
	ErrLookupNameReq = "Lookup name is required"
	defaultFields    = "id, name, description"
)

// CreateStatus implements api.StatusesServer.
func (s StatusService) CreateStatus(ctx context.Context, req *_go.CreateStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("lookup.name.required", ErrLookupNameReq)
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictinary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
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
		return nil, e
	}

	return l, nil
}

// ListStatuses implements api.StatusesServer.
func (s StatusService) ListStatuses(ctx context.Context, req *_go.ListStatusRequest) (*_go.StatusList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictinary)
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

	lookups, e := s.app.Store.Status().List(&searchOptions)
	if e != nil {
		return nil, e
	}

	return lookups, nil
}

// UpdateStatus implements api.StatusesServer.
func (s StatusService) UpdateStatus(ctx context.Context, req *_go.UpdateStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("lookup.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictinary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the updater
	currentU := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update lookup model
	lookup := &_go.Status{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "name", "description", "updated_at", "updated_by"}

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
		return nil, e
	}

	return l, nil
}

// DeleteStatus implements api.StatusesServer.
func (s StatusService) DeleteStatus(ctx context.Context, req *_go.DeleteStatusRequest) (*_go.Status, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("lookup.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictinary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// // RBAC check
	// if scope.IsRbacUsed() {
	// 	access, err := s.app.Store.Status().RbacAccess(ctx, session.GetDomainId(), int64(req.GetId()), session.GetAclRoles(), accessMode.Value())
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if !access {
	// 		return nil, s.app.MakeScopeError(session, scope, accessMode)
	// 	}
	// }

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
		return nil, e
	}

	return &(_go.Status{Id: req.Id}), nil
}

// LocateStatus implements api.StatusesServer.
func (s StatusService) LocateStatus(ctx context.Context, req *_go.LocateStatusRequest) (*_go.LocateStatusResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("status.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictinary)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFields, ", ")
	}

	t := time.Now()

	searchOpts := model.SearchOptions{
		IDs:     []int64{req.Id},
		Session: session,
		Context: ctx,
		Fields:  fields,
		Page:    1,
		Size:    1,
		Time:    t,
	}

	l, e := s.app.Store.Status().List(&searchOpts)
	if e != nil {
		return nil, e
	}

	if len(l.Items) == 0 {
		return nil, model.NewNotFoundError("status_lookup.not_found", "Status lookup not found")
	}

	lookup := l.Items[0]

	return &_go.LocateStatusResponse{Status: lookup}, nil
}

func NewStatusService(app *App) (*StatusService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_status.args_check.app_nil", "internal is nil")
	}
	return &StatusService{app: app}, nil
}
