package lookup

import (
	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	_general "buf.build/gen/go/webitel/general/protocolbuffers/go"
	"context"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/internal/app"
	"github.com/webitel/cases/model"
	"strings"
)

type LookupStatusService struct {
	app *app.App
}

const (
	ErrStatusNameReq          = "Status name is required"
	defaultFieldsLookupStatus = "id, name, description,is_initial,is_final"
)

func (s LookupStatusService) CreateLookupStatus(ctx context.Context, req *_go.CreateLookupStatusRequest) (*_go.
	LookupStatus, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("status.name.required", ErrStatusNameReq)
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the creator and updater
	currentU := &_general.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new status model
	status := &_go.LookupStatus{
		LookupId:    req.LookupId,
		Name:        req.Name,
		Description: req.Description,
		IsInitial:   req.IsInitial,
		IsFinal:     req.IsFinal,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "lookup_id", "name", "description", "is_initial", "is_final", "created_at", "updated_at", "created_by", "updated_by"}

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
	}

	// Create the status in the store
	st, e := s.app.Store.LookupStatus().Attach(&createOpts, status)
	if e != nil {
		return nil, e
	}

	return st, nil
}

func (s LookupStatusService) ListLookupStatuses(ctx context.Context, req *_go.ListLookupStatusesRequest) (*_go.LookupStatusList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsLookupStatus, ", ")
	}

	// Use default page size and page number if not provided
	page := req.Page
	if page == 0 {
		page = 1
	}

	searchOptions := model.SearchOptions{
		IDs:     req.Id,
		Session: session,
		Fields:  req.Fields,
		Context: ctx,
		Page:    int(page),
		Size:    int(req.Size),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	} else if req.Name != "" {
		searchOptions.Filter["name"] = req.Name
	}

	statuses, e := s.app.Store.LookupStatus().List(&searchOptions)
	if e != nil {
		return nil, e
	}

	return statuses, nil
}

func (s LookupStatusService) UpdateLookupStatus(ctx context.Context, req *_go.UpdateLookupStatusRequest) (*_go.LookupStatus, error) {
	// Validate required fields
	if req.Id == 0 || req.Name == "" {
		return nil, model.NewBadRequestError("status.id_name.required", "Status ID and name are required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the updater
	currentU := &_general.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update status model
	status := &_go.LookupStatus{
		Id:          req.Id,
		LookupId:    req.LookupId,
		Name:        req.Name,
		Description: req.Description,
		IsInitial:   req.IsInitial,
		IsFinal:     req.IsFinal,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "lookup_id", "name", "description", "is_initial", "is_final", "updated_at", "updated_by"}

	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
	}

	// Update the status in the store
	st, e := s.app.Store.LookupStatus().Update(&updateOpts, status)
	if e != nil {
		return nil, e
	}

	return st, nil
}

func (s LookupStatusService) DeleteLookupStatus(ctx context.Context, req *_go.DeleteLookupStatusRequest) (*_go.LookupStatus, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("status.id.required", "Status ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define delete options
	deleteOpts := model.DeleteOptions{
		Session: session,
		Context: ctx,
		IDs:     []int64{req.Id},
	}

	// Delete the status in the store
	e := s.app.Store.LookupStatus().Delete(&deleteOpts)
	if e != nil {
		return nil, e
	}

	return &(_go.LookupStatus{Id: req.Id}), nil
}

func (s LookupStatusService) LocateLookupStatus(ctx context.Context, req *_go.LocateLookupStatusRequest) (*_go.LocateLookupStatusResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("status.id.required", "Status ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, err
	}
	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFields, ", ")
	}

	searchOpts := model.SearchOptions{
		IDs:     []int64{req.Id},
		Session: session,
		Context: ctx,
		Fields:  req.Fields,
		Page:    1,
		Size:    1,
	}

	l, e := s.app.Store.LookupStatus().List(&searchOpts)
	if e != nil {
		return nil, e
	}

	if len(l.Items) == 0 {
		return nil, model.NewNotFoundError("status.not_found", "Status not found")
	}

	status := l.Items[0]

	return &_go.LocateLookupStatusResponse{Status: status}, nil
}

func NewLookupStatusService(app *app.App) (*LookupStatusService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_lookup_status_service.args_check.app_nil", "pkg is nil")
	}
	return &LookupStatusService{app: app}, nil
}
