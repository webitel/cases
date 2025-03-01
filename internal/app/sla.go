package app

import (
	"context"
	"strings"
	"time"

	cases "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/util"

	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
)

type SLAService struct {
	app *App
	cases.UnimplementedSLAsServer
	objClassName string
}

const (
	defaultFieldsSLA = "id, name, description, calendar"
)

// CreateSLA implements cases.SLAsServer.
func (s *SLAService) CreateSLA(ctx context.Context, req *cases.CreateSLARequest) (*cases.SLA, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.name.required", "SLA name is required")
	}
	if req.Calendar.GetId() == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.calendar_id.required", "Calendar ID is required")
	}
	if req.ReactionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.reaction_time.required", "Reaction time is required")
	}
	if req.ResolutionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.resolution_time.required", "Resolution time is required")
	}

	fields := []string{
		"id", "lookup_id", "name", "description", "valid_from",
		"valid_to", "calendar_id", "reaction_time", "resolution_time", "created_at", "updated_at",
		"created_by", "updated_by",
	}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Context: ctx,
		Fields:  fields,
		Time:    t,
		Auth:    model.GetAutherOutOfContext(ctx),
	}

	// Define the current user as the creator and updater
	currentU := &cases.Lookup{
		Id: createOpts.GetAuthOpts().GetUserId(),
	}

	// Create a new SLA user_auth
	sla := &cases.SLA{
		Name:           req.Name,
		Description:    req.Description,
		ValidFrom:      req.ValidFrom,
		ValidTo:        req.ValidTo,
		Calendar:       req.Calendar,
		ReactionTime:   req.ReactionTime,
		ResolutionTime: req.ResolutionTime,
		CreatedBy:      currentU,
		UpdatedBy:      currentU,
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

	t := time.Now()
	// Define delete options
	deleteOpts := model.DeleteOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
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
		IDs: req.Id,
		//UserAuthSession: session,
		Fields:  fields,
		Context: ctx,
		Sort:    req.Sort,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
		Auth:    model.GetAutherOutOfContext(ctx),
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

	fields := []string{"id"}

	// Map XJsonMask fields to the corresponding SLA fields
	for _, f := range req.XJsonMask {
		if strings.HasPrefix(f, "calendar") {
			// Handle fields with "calendar." prefix
			if !util.ContainsField(fields, "calendar_id") {
				fields = append(fields, "calendar_id")
			}
			if req.Input.Calendar.GetId() == 0 {
				return nil, cerror.NewBadRequestError("sla_service.update_sla.calendar_id.required", "Calendar ID is required")
			}
			continue
		}
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
		case "reaction_time":
			fields = append(fields, "reaction_time")
		case "resolution_time":
			fields = append(fields, "resolution_time")
		}
	}

	t := time.Now()

	// Define update options
	updateOpts := model.UpdateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Define the current user as the updater
	u := &cases.Lookup{
		Id: updateOpts.GetAuthOpts().GetUserId(),
	}

	// Update SLA user_auth
	sla := &cases.SLA{
		Id:             req.Id,
		Name:           req.Input.Name,
		Description:    req.Input.Description,
		ValidFrom:      req.Input.ValidFrom,
		ValidTo:        req.Input.ValidTo,
		Calendar:       req.Input.Calendar,
		ReactionTime:   req.Input.ReactionTime,
		ResolutionTime: req.Input.ResolutionTime,
		UpdatedBy:      u,
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
	return &SLAService{app: app, objClassName: model.ScopeDictionary}, nil
}
