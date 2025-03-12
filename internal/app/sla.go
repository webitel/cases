package app

import (
	"context"

	"github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
)

type SLAService struct {
	app *App
	cases.UnimplementedSLAsServer
	objClassName string
}

var SLAMetadata = model.NewObjectMetadata(
	model.ScopeDictionary,
	"",
	[]*model.Field{
		{"id", true},
		{"created_by", true},
		{"created_at", true},
		{"updated_by", false},
		{"updated_at", false},
		{"name", true},
		{"description", true},
		{"valid_from", true},
		{"valid_to", true},
		{"calendar", true},
		{"reaction_time", true},
		{"resolution_time", true},
	})

// CreateSLA implements cases.SLAsServer.
func (s *SLAService) CreateSLA(ctx context.Context, req *cases.CreateSLARequest) (*cases.SLA, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.name.required", "SLA name is required")
	}
	if req.Input.Calendar.GetId() == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.calendar_id.required", "Calendar ID is required")
	}
	if req.Input.ReactionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.reaction_time.required", "Reaction time is required")
	}
	if req.Input.ResolutionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.resolution_time.required", "Resolution time is required")
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, SLAMetadata),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Create a new SLA user_auth
	input := &cases.SLA{
		Name:           req.Input.Name,
		Description:    req.Input.Description,
		ValidFrom:      req.Input.ValidFrom,
		ValidTo:        req.Input.ValidTo,
		Calendar:       req.Input.Calendar,
		ReactionTime:   req.Input.ReactionTime,
		ResolutionTime: req.Input.ResolutionTime,
	}

	// Create the SLA in the store
	res, err := s.app.Store.SLA().Create(createOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.create_sla.store.create.failed", err.Error())
	}

	return res, nil
}

// DeleteSLA implements cases.SLAsServer.
func (s *SLAService) DeleteSLA(ctx context.Context, req *cases.DeleteSLARequest) (*cases.SLA, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_service.delete_sla.id.required", "SLA ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Delete the SLA in the store
	err = s.app.Store.SLA().Delete(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.delete_sla.store.delete.failed", err.Error())
	}

	return &cases.SLA{Id: req.Id}, nil
}

// ListSLAs implements cases.SLAsServer.
func (s *SLAService) ListSLAs(ctx context.Context, req *cases.ListSLARequest) (*cases.SLAList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, SLAMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDs(req.GetId()),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	searchOpts.AddFilter("name", req.GetQ())

	res, err := s.app.Store.SLA().List(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.list_slas.store.list.failed", err.Error())
	}

	return res, nil
}

// LocateSLA implements cases.SLAsServer.
func (s *SLAService) LocateSLA(ctx context.Context, req *cases.LocateSLARequest) (*cases.LocateSLAResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_service.locate_sla.id.required", "SLA ID is required")
	}

	fields := util.FieldsFunc(req.Fields, util.InlineFields)

	// Prepare a list request with necessary parameters
	listReq := &cases.ListSLARequest{
		Id:     []int64{req.Id},
		Fields: fields,
		Page:   1,
		Size:   1,
	}

	// Call the ListSLAs method
	res, err := s.ListSLAs(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.locate_sla.list_slas.error", err.Error())
	}

	// Check if the SLA was found
	if len(res.Items) == 0 {
		return nil, cerror.NewNotFoundError("sla_service.locate_sla.not_found", "SLA not found")
	}

	// Return the found SLA
	return &cases.LocateSLAResponse{Sla: res.Items[0]}, nil
}

// UpdateSLA implements cases.SLAsServer.
func (s *SLAService) UpdateSLA(ctx context.Context, req *cases.UpdateSLARequest) (*cases.SLA, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_service.update_sla.id.required", "SLA ID is required")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, SLAMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Update SLA user_auth
	input := &cases.SLA{
		Id:             req.Id,
		Name:           req.Input.Name,
		Description:    req.Input.Description,
		ValidFrom:      req.Input.ValidFrom,
		ValidTo:        req.Input.ValidTo,
		Calendar:       req.Input.Calendar,
		ReactionTime:   req.Input.ReactionTime,
		ResolutionTime: req.Input.ResolutionTime,
	}

	// Update the SLA in the store
	res, err := s.app.Store.SLA().Update(updateOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.update_sla.store.update.failed", err.Error())
	}

	return res, nil
}

func NewSLAService(app *App) (*SLAService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_sla_service.args_check.app_nil", "internal is nil")
	}
	return &SLAService{app: app, objClassName: model.ScopeDictionary}, nil
}
