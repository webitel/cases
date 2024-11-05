package app

import (
	"context"
	"strings"
	"time"

	cases "buf.build/gen/go/webitel/cases/protocolbuffers/go"
	general "buf.build/gen/go/webitel/general/protocolbuffers/go"
	authmodel "github.com/webitel/cases/auth/model"

	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
)

type SLAService struct {
	app *App
}

const (
	defaultFieldsSLA = "id, name, description"
)

// CreateSLA implements cases.SLAsServer.
func (s *SLAService) CreateSLA(ctx context.Context, req *cases.CreateSLARequest) (*cases.SLA, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.name.required", "SLA name is required")
	}
	if req.CalendarId == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.calendar_id.required", "Calendar ID is required")
	}
	if req.ReactionTime.Hours == 0 && req.ReactionTime.Minutes == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.reaction_time.required", "Reaction time is required")
	}
	if req.ResolutionTime.Hours == 0 && req.ResolutionTime.Minutes == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.resolution_time.required", "Resolution time is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_service.create_sla.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the creator and updater
	currentU := &general.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new SLA model
	sla := &cases.SLA{
		Name:        req.Name,
		Description: req.Description,
		ValidFrom:   req.ValidFrom.AsTime().Unix(),
		ValidTo:     req.ValidTo.AsTime().Unix(),
		Calendar:    &general.Lookup{Id: req.CalendarId},
		ReactionTime: &cases.ReactionTime{
			Hours:   req.ReactionTime.Hours,
			Minutes: req.ReactionTime.Minutes,
		},
		ResolutionTime: &cases.ResolutionTime{
			Hours:   req.ResolutionTime.Hours,
			Minutes: req.ResolutionTime.Minutes,
		},
		CreatedBy: currentU,
		UpdatedBy: currentU,
	}

	fields := []string{
		"id", "lookup_id", "name", "description", "valid_from",
		"valid_to", "calendar_id", "reaction_time_hours", "reaction_time_minutes",
		"resolution_time_hours", "resolution_time_minutes", "created_at", "updated_at",
		"created_by", "updated_by",
	}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Create the SLA in the store
	r, e := s.app.Store.SLA().Create(&createOpts, sla)
	if e != nil {
		return nil, cerror.NewInternalError("sla_service.create_sla.store.create.failed", e.Error())
	}

	return r, nil
}

// DeleteSLA implements cases.SLAsServer.
func (s *SLAService) DeleteSLA(ctx context.Context, req *cases.DeleteSLARequest) (*cases.SLA, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_service.delete_sla.id.required", "SLA ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_service.delete_sla.authorization.failed", err.Error())
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

	// Delete the SLA in the store
	e := s.app.Store.SLA().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("sla_service.delete_sla.store.delete.failed", e.Error())
	}

	return &cases.SLA{Id: req.Id}, nil
}

// ListSLAs implements cases.SLAsServer.
func (s *SLAService) ListSLAs(ctx context.Context, req *cases.ListSLARequest) (*cases.SLAList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_service.list_slas.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsSLA, ", ")
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
		Page:    int32(page),
		Size:    int32(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	slas, e := s.app.Store.SLA().List(&searchOptions)
	if e != nil {
		return nil, cerror.NewInternalError("sla_service.list_slas.store.list.failed", e.Error())
	}

	return slas, nil
}

// LocateSLA implements cases.SLAsServer.
func (s *SLAService) LocateSLA(ctx context.Context, req *cases.LocateSLARequest) (*cases.LocateSLAResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_service.locate_sla.id.required", "SLA ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &cases.ListSLARequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1, // We only need one item
	}

	// Call the ListSLAs method
	listResp, err := s.ListSLAs(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.locate_sla.list_slas.error", err.Error())
	}

	// Check if the SLA was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("sla_service.locate_sla.not_found", "SLA not found")
	}

	// Return the found SLA
	return &cases.LocateSLAResponse{Sla: listResp.Items[0]}, nil
}

// UpdateSLA implements cases.SLAsServer.
func (s *SLAService) UpdateSLA(ctx context.Context, req *cases.UpdateSLARequest) (*cases.SLA, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_service.update_sla.id.required", "SLA ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_service.update_sla.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the updater
	u := &general.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Initialize ReactionTime and ResolutionTime if nil
	if req.Input.ReactionTime == nil {
		req.Input.ReactionTime = &cases.ReactionTime{}
	}
	if req.Input.ResolutionTime == nil {
		req.Input.ResolutionTime = &cases.ResolutionTime{}
	}

	// Update SLA model
	sla := &cases.SLA{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		ValidFrom:   req.Input.ValidFrom.AsTime().Unix(),
		ValidTo:     req.Input.ValidTo.AsTime().Unix(),
		Calendar:    &general.Lookup{Id: req.Input.CalendarId},
		ReactionTime: &cases.ReactionTime{
			Hours:   req.Input.ReactionTime.Hours,
			Minutes: req.Input.ReactionTime.Minutes,
		},
		ResolutionTime: &cases.ResolutionTime{
			Hours:   req.Input.ResolutionTime.Hours,
			Minutes: req.Input.ResolutionTime.Minutes,
		},
		UpdatedBy: u,
	}

	fields := []string{"id"}

	// Map XJsonMask fields to the corresponding SLA fields
	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("sla_service.update_sla.name.required", "SLA name is required and cannot be empty")
			}
		case "description":
			fields = append(fields, "description")
		case "valid_from":
			fields = append(fields, "valid_from")
		case "valid_to":
			fields = append(fields, "valid_to")
		case "calendar_id":
			fields = append(fields, "calendar_id")
			if req.Input.CalendarId == 0 {
				return nil, cerror.NewBadRequestError("sla_service.update_sla.calendar_id.required", "Calendar ID is required")
			}
		case "reaction_time_hours":
			fields = append(fields, "reaction_time_hours")
		case "reaction_time_minutes":
			fields = append(fields, "reaction_time_minutes")
		case "resolution_time_hours":
			fields = append(fields, "resolution_time_hours")
		case "resolution_time_minutes":
			fields = append(fields, "resolution_time_minutes")
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

	// Update the SLA in the store
	r, e := s.app.Store.SLA().Update(&updateOpts, sla)
	if e != nil {
		return nil, cerror.NewInternalError("sla_service.update_sla.store.update.failed", e.Error())
	}

	return r, nil
}

func NewSLAService(app *App) (*SLAService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_sla_service.args_check.app_nil", "internal is nil")
	}
	return &SLAService{app: app}, nil
}
