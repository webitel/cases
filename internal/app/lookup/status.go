package lookup

import (
	"context"
	"github.com/webitel/cases/internal/app"
	"strings"
	"time"

	_go "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	_general "buf.build/gen/go/webitel/general/protocolbuffers/go"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type StatusLookupService struct {
	app *app.App
}

const (
	ErrLookupNameReq = "Lookup name is required"
	defaultFields    = "id, name, description"
)

func (s StatusLookupService) CreateStatusLookup(ctx context.Context, req *_go.CreateStatusLookupRequest) (*_go.StatusLookup, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("groups.name.required", ErrLookupNameReq)
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

	now := time.Now().UTC()

	// Define the current user as the creator and updater
	currentU := &_general.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new lookup model
	lookup := &_go.StatusLookup{
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
		Time:    now,
	}

	// Create the group in the db
	l, e := s.app.DB.Status().Create(&createOpts, lookup)
	if e != nil {
		return nil, e
	}

	return l, nil
}

func (s StatusLookupService) SearchStatusLookups(ctx context.Context, req *_go.ListStatusLookupsRequest) (*_go.StatusLookupList, error) {
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

	lookups, e := s.app.DB.Status().Search(&searchOptions)
	if e != nil {
		return nil, e
	}

	return lookups, nil
}

func (s StatusLookupService) UpdateStatusLookup(ctx context.Context, req *_go.UpdateStatusLookupRequest) (*_go.StatusLookup, error) {
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

	now := time.Now().UTC()

	// Define the current user as the updater
	currentU := &_general.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update lookup model
	lookup := &_go.StatusLookup{
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
		Time:    now,
	}

	// Update the lookup in the db
	l, e := s.app.DB.Status().Update(&updateOpts, lookup)
	if e != nil {
		return nil, e
	}

	return l, nil
}

func (s StatusLookupService) DeleteStatusLookup(ctx context.Context, req *_go.DeleteStatusLookupRequest) (*_go.StatusLookup, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, model.NewBadRequestError("groups.id.required", "Lookup ID is required")
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

	// Delete the lookup in the db
	e := s.app.DB.Status().Delete(&deleteOpts)
	if e != nil {
		return nil, e
	}

	return &(_go.StatusLookup{Id: req.Id}), nil
}

func (s StatusLookupService) LocateStatusLookup(ctx context.Context, req *_go.LocateStatusLookupRequest) (*_go.LocateStatusLookupResponse, error) {
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

	l, e := s.app.DB.Status().Search(&searchOpts)
	if e != nil {
		return nil, e
	}

	if len(l.Items) == 0 {
		return nil, model.NewNotFoundError("status_lookup.not_found", "Status lookup not found")
	}

	lookup := l.Items[0]

	return &_go.LocateStatusLookupResponse{Lookup: lookup}, nil
}

func NewStatusLookupService(app *app.App) (*StatusLookupService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_status_lookup_service.args_check.app_nil", "pkg is nil")
	}
	return &StatusLookupService{app: app}, nil
}
