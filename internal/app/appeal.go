package app

import (
	"context"
	"strings"

	_go "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type AppealService struct {
	app *App
}

func (s AppealService) CreateAppeal(ctx context.Context, req *_go.CreateAppealRequest) (*_go.Appeal, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("appeal_service.create_appeal.name.required", ErrLookupNameReq)
	}

	// Validate the Type field
	if req.Type == _go.Type_TYPE_UNSPECIFIED {
		return nil, model.NewBadRequestError("appeal_service.create_appeal.type.required", "Appeal type is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("appeal_service.create_appeal.authorization.failed", err.Error())
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
		return nil, model.NewInternalError("appeal_service.create_appeal.store.create.failed", e.Error())
	}

	return l, nil
}

func (s AppealService) ListAppeals(ctx context.Context, req *_go.ListAppealRequest) (*_go.AppealList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("appeal_service.list_appeals.authorization.failed", err.Error())
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

	searchOptions := model.SearchOptions{
		IDs:     req.Id,
		Session: session,
		Fields:  fields,
		Context: ctx,
		Page:    int(page),
		Size:    int(req.Size),
		Filter:  make(map[string]interface{}),
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
		return nil, model.NewInternalError("appeal_service.list_appeals.store.list.failed", e.Error())
	}

	return lookups, nil
}

func (s AppealService) UpdateAppeal(ctx context.Context, req *_go.UpdateAppealRequest) (*_go.Appeal, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("appeal_service.update_appeal.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("appeal_service.update_appeal.authorization.failed", err.Error())
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

	// Update appeal model
	appeal := &_go.Appeal{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "updated_at", "updated_by"}

	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
		case "description":
			fields = append(fields, "description")
		case "type":
			fields = append(fields, "type")
		}
	}

	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
	}

	// Update the appeal in the store
	l, e := s.app.Store.Appeal().Update(&updateOpts, appeal)
	if e != nil {
		return nil, model.NewInternalError("appeal_service.update_appeal.store.update.failed", e.Error())
	}

	return l, nil
}

func (s AppealService) DeleteAppeal(ctx context.Context, req *_go.DeleteAppealRequest) (*_go.Appeal, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("appeal_service.delete_appeal.id.required", "Lookup ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("appeal_service.delete_appeal.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
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
		return nil, model.NewInternalError("appeal_service.delete_appeal.store.delete.failed", e.Error())
	}

	return &(_go.Appeal{Id: req.Id}), nil
}

func (s AppealService) LocateAppeal(ctx context.Context, req *_go.LocateAppealRequest) (*_go.LocateAppealResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("appeal_service.locate_appeal.id.required", "Lookup ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListAppealRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1, // We only need one item
	}

	// Call the ListAppeals method
	listResp, err := s.ListAppeals(ctx, listReq)
	if err != nil {
		return nil, model.NewInternalError("appeal_service.locate_appeal.list_appeals.error", err.Error())
	}

	// Check if the appeal was found
	if len(listResp.Items) == 0 {
		return nil, model.NewNotFoundError("appeal_service.locate_appeal.not_found", "Appeal not found")
	}

	// Return the found appeal
	return &_go.LocateAppealResponse{Appeal: listResp.Items[0]}, nil
}

func NewAppealService(app *App) (*AppealService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_appeal_service.args_check.app_nil", "internal is nil")
	}
	return &AppealService{app: app}, nil
}
