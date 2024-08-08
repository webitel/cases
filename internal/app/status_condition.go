package app

import (
	"context"
	"strings"
	"time"

	_go "github.com/webitel/cases/api"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type StatusConditionService struct {
	app *App
}

const (
	ErrStatusNameReq    = "Status name is required"
	defaultFieldsStatus = "id, name, description,is_initial,is_final"
)

// CreateStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) CreateStatusCondition(ctx context.Context, req *_go.CreateStatusConditionRequest) (*_go.StatusCondition, error) {
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
	currentU := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new status model
	status := &_go.StatusCondition{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
		StatusId:    req.StatusId,
	}

	fields := []string{"id", "lookup_id", "name", "description", "initial", "final", "created_at", "updated_at", "created_by", "updated_by"}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Create the status in the store
	st, e := s.app.Store.StatusCondition().Create(&createOpts, status)
	if e != nil {
		return nil, e
	}

	return st, nil
}

// ListStatusConditions implements api.StatusConditionsServer.
func (s StatusConditionService) ListStatusConditions(ctx context.Context, req *_go.ListStatusConditionRequest) (*_go.StatusConditionList, error) {
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
		fields = strings.Split(defaultFieldsStatus, ", ")
	}

	// Use default page size and page number if not provided
	page := req.Page
	if page == 0 {
		page = 1
	}

	t := time.Now()
	searchOptions := model.SearchOptions{
		Session: session,
		Fields:  fields,
		Context: ctx,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	} else if req.Name != "" {
		searchOptions.Filter["name"] = req.Name
	}

	statuses, e := s.app.Store.StatusCondition().List(&searchOptions, req.StatusId)
	if e != nil {
		return nil, e
	}

	return statuses, nil
}

// UpdateStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) UpdateStatusCondition(ctx context.Context, req *_go.UpdateStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("status.id.required", "Status ID is required")
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
	u := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update status model
	status := &_go.StatusCondition{
		Id:          req.Id,
		StatusId:    req.StatusId,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Initial:     req.Input.Initial.Value,
		Final:       req.Input.Final.Value,
		UpdatedBy:   u,
	}

	fields := []string{"id", "lookup_id", "name", "description", "initial", "final", "updated_at", "updated_by"}

	t := time.Now()

	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Update the status in the store
	st, e := s.app.Store.StatusCondition().Update(&updateOpts, status)
	if e != nil {
		return nil, e
	}

	return st, nil
}

// PatchStatusCondition implements api.StatusConditionsServer.
func (s *StatusConditionService) PatchStatusCondition(ctx context.Context, req *_go.PatchStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("status.id.required", "Status ID is required")
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
	u := &_go.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Initialize the status update object
	status := &_go.StatusCondition{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		UpdatedBy:   u,
		StatusId:    req.StatusId,
	}

	t := time.Now()

	// Collect fields to be updated
	var fields []string
	if req.Input.Initial != nil {
		fields = append(fields, "initial")
		status.Initial = req.Input.Initial.Value
	}
	if req.Input.Final != nil {
		fields = append(fields, "final")
		status.Final = req.Input.Final.Value
	}
	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Time:    t,
		Fields:  fields,
	}

	// Update the status in the store
	st, e := s.app.Store.StatusCondition().Update(&updateOpts, status)
	if e != nil {
		return nil, e
	}

	return st, nil
}

// DeleteStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) DeleteStatusCondition(ctx context.Context, req *_go.DeleteStatusConditionRequest) (*_go.StatusCondition, error) {
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

	t := time.Now()
	// Define delete options
	deleteOpts := model.DeleteOptions{
		Session: session,
		Context: ctx,
		IDs:     []int64{req.Id},
		Time:    t,
	}

	// Delete the status in the store
	e := s.app.Store.StatusCondition().Delete(&deleteOpts, req.StatusId)
	if e != nil {
		return nil, e
	}

	return &(_go.StatusCondition{Id: req.Id}), nil
}

// LocateStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) LocateStatusCondition(ctx context.Context, req *_go.LocateStatusConditionRequest) (*_go.LocateStatusConditionResponse,
	error,
) {
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

	l, e := s.app.Store.StatusCondition().List(&searchOpts, req.StatusId)
	if e != nil {
		return nil, e
	}

	if len(l.Items) == 0 {
		return nil, model.NewNotFoundError("status.not_found", "Status not found")
	}

	status := l.Items[0]

	return &_go.LocateStatusConditionResponse{Status: status}, nil
}

func NewStatusConditionService(app *App) (*StatusConditionService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_status_condition_service.args_check.app_nil", "internal is nil")
	}
	return &StatusConditionService{app: app}, nil
}
