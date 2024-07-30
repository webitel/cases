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

type AppealLookupService struct {
	app *app.App
}

func (s AppealLookupService) CreateAppealLookup(ctx context.Context, req *_go.CreateAppealLookupRequest) (*_go.AppealLookup, error) {
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
	scope := session.GetScope(model.ScopeLog)
	if !session.HasAccess(scope, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the creator and updater
	currentU := &_general.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new lookup model
	lookup := &_go.AppealLookup{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "name", "description", "created_at", "updated_at", "created_by", "updated_by"}

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
	}

	// Create the group in the store
	l, e := s.app.Store.AppealLookup().Create(&createOpts, lookup)
	if e != nil {
		return nil, e
	}

	return l, nil
}

func (s AppealLookupService) ListAppealLookups(ctx context.Context, req *_go.ListAppealLookupsRequest) (*_go.AppealLookupList, error) {
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

	lookups, e := s.app.Store.AppealLookup().List(&searchOptions)
	if e != nil {
		return nil, e
	}

	return lookups, nil
}

func (s AppealLookupService) UpdateAppealLookup(ctx context.Context, req *_go.UpdateAppealLookupRequest) (*_go.AppealLookup, error) {
	// Validate required fields
	if req.Id == 0 || req.Name == "" {
		return nil, model.NewBadRequestError("groups.id_name.required", "Lookup ID and name are required")
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

	// Update lookup model
	lookup := &_go.AppealLookup{
		Id:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		UpdatedBy:   currentU,
	}

	fields := []string{"id", "name", "description", "updated_at", "updated_by"}

	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
	}

	// Update the lookup in the store
	l, e := s.app.Store.AppealLookup().Update(&updateOpts, lookup)
	if e != nil {
		return nil, e
	}

	return l, nil
}

func (s AppealLookupService) DeleteAppealLookup(ctx context.Context, req *_go.DeleteAppealLookupRequest) (*_go.AppealLookup, error) {
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

	// Delete the lookup in the store
	e := s.app.Store.CloseReasonLookup().Delete(&deleteOpts)
	if e != nil {
		return nil, e
	}

	return &(_go.AppealLookup{Id: req.Id}), nil
}

func (s AppealLookupService) LocateAppealLookup(ctx context.Context, req *_go.LocateAppealLookupRequest) (*_go.LocateAppealLookupResponse, error) {
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

	l, e := s.app.Store.AppealLookup().List(&searchOpts)
	if e != nil {
		return nil, e
	}

	if len(l.Items) == 0 {
		return nil, model.NewNotFoundError("close_reason_lookup.not_found", "CloseReason lookup not found")
	}

	lookup := l.Items[0]

	return &_go.LocateAppealLookupResponse{Lookup: lookup}, nil
}

func NewAppealLookupService(app *app.App) (*AppealLookupService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_appeal_lookup_service.args_check.app_nil", "pkg is nil")
	}
	return &AppealLookupService{app: app}, nil
}
