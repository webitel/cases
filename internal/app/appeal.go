package app

import (
	"context"
	"strings"

	_go "github.com/webitel/cases/api"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type AppealService struct {
	app *App
}

func (s AppealService) CreateAppeal(ctx context.Context, req *_go.CreateAppealRequest) (*_go.Appeal, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("lookup.name.required", ErrLookupNameReq)
	}

	// Validate the Type field
	if req.Type == _go.Type_TYPE_UNSPECIFIED {
		return nil, model.NewBadRequestError("lookup.type.required", "Appeal type is required")
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

	// Create a new appeal model
	appeal := &_go.Appeal{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "name", "description", "type", "created_at", "updated_at", "created_by", "updated_by"}

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
	}

	// Create the appeal in the store
	l, e := s.app.Store.Appeal().Create(&createOpts, appeal)
	if e != nil {
		return nil, e
	}

	return l, nil
}

func (s AppealService) ListAppeals(ctx context.Context, req *_go.ListAppealRequest) (*_go.AppealList, error) {
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

	searchOptions := model.SearchOptions{
		IDs:     req.Id,
		Session: session,
		Fields:  fields,
		Context: ctx,
		Page:    int(page),
		Size:    int(req.Size),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	} else if req.Name != "" {
		searchOptions.Filter["name"] = req.Name
	}

	if len(req.Type) > 0 {
		searchOptions.Filter["type"] = req.Type
	}

	lookups, e := s.app.Store.Appeal().List(&searchOptions)
	if e != nil {
		return nil, e
	}

	return lookups, nil
}

func (s AppealService) UpdateAppeal(ctx context.Context, req *_go.UpdateAppealRequest) (*_go.Appeal, error) {
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

	// Update appeal model
	appeal := &_go.Appeal{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "name", "description", "type", "updated_at", "updated_by"}

	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
	}

	// Update the appeal in the store
	l, e := s.app.Store.Appeal().Update(&updateOpts, appeal)
	if e != nil {
		return nil, e
	}

	return l, nil
}

func (s AppealService) DeleteAppeal(ctx context.Context, req *_go.DeleteAppealRequest) (*_go.Appeal, error) {
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

	// Define delete options
	deleteOpts := model.DeleteOptions{
		Session: session,
		Context: ctx,
		IDs:     []int64{req.Id},
	}

	// Delete the appeal in the store
	e := s.app.Store.Appeal().Delete(&deleteOpts)
	if e != nil {
		return nil, e
	}

	return &(_go.Appeal{Id: req.Id}), nil
}

func (s AppealService) LocateAppeal(ctx context.Context, req *_go.LocateAppealRequest) (*_go.LocateAppealResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("groups.id.required", "Lookup ID is required")
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

	searchOpts := model.SearchOptions{
		IDs:     []int64{req.Id},
		Session: session,
		Context: ctx,
		Fields:  fields,
		Page:    1,
		Size:    1,
	}

	l, e := s.app.Store.Appeal().List(&searchOpts)
	if e != nil {
		return nil, e
	}

	if len(l.Items) == 0 {
		return nil, model.NewNotFoundError("appeal.not_found", "Appeal not found")
	}

	lookup := l.Items[0]

	return &_go.LocateAppealResponse{Appeal: lookup}, nil
}

func NewAppealService(app *App) (*AppealService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_appeal_service.args_check.app_nil", "internal is nil")
	}
	return &AppealService{app: app}, nil
}
