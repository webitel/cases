package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
)

type SLAConditionService struct {
	app *App
	cases.UnimplementedSLAConditionsServer
	objClassName string
}

var SLAConditionMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "name", Default: true},
	{Name: "priorities", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "reaction_time", Default: true},
	{Name: "resolution_time", Default: true},
	{Name: "sla_id", Default: true},
})

// CreateSLACondition implements cases.SLAConditionsServer.
func (s *SLAConditionService) CreateSLACondition(ctx context.Context, req *cases.CreateSLAConditionRequest) (*cases.SLACondition, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.name.required", "SLA Condition name is required")
	}
	if len(req.Input.Priorities) == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.priorities.required", "At least one priority is required")
	}
	if req.Input.ReactionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.reaction_time.required", "Reaction time is required")
	}
	if req.Input.ResolutionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.resolution_time.required", "Resolution time is required")
	}
	if req.SlaId == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.sla_id.required", "SLA ID is required")
	}

	// Convert []*cases.Lookup to []int64
	var priorityIDs []int64
	for _, priority := range req.Input.Priorities {
		if priority != nil { // Check for nil to avoid runtime panic
			priorityIDs = append(priorityIDs, priority.GetId()) // Use GetId() to ensure proper handling
		}
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, SLAConditionMetadata),
		grpcopts.WithCreateIDs(priorityIDs),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Create a new SLACondition user_auth
	slaCondition := &cases.SLACondition{
		Name:           req.Input.Name,
		ReactionTime:   req.Input.ReactionTime,
		ResolutionTime: req.Input.ResolutionTime,
		SlaId:          req.SlaId,
	}

	// Create the SLACondition in the store
	r, e := s.app.Store.SLACondition().Create(createOpts, slaCondition, priorityIDs)
	if e != nil {
		return nil, cerror.NewInternalError("sla_condition_service.create_sla_condition.store.create.failed", e.Error())
	}

	return r, nil
}

// DeleteSLACondition implements cases.SLAConditionsServer.
func (s *SLAConditionService) DeleteSLACondition(ctx context.Context, req *cases.DeleteSLAConditionRequest) (*cases.SLACondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.delete_sla_condition.id.required", "SLA Condition ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Delete the SLACondition in the store
	e := s.app.Store.SLACondition().Delete(deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("sla_condition_service.delete_sla_condition.store.delete.failed", e.Error())
	}

	return &cases.SLACondition{Id: req.Id}, nil
}

// ListSLAConditions implements cases.SLAConditionsServer.
func (s *SLAConditionService) ListSLAConditions(ctx context.Context, req *cases.ListSLAConditionRequest) (*cases.SLAConditionList, error) {
	searchOptions, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, SLAConditionMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithIDs(req.Id),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}
	searchOptions.AddFilter("sla_id", req.SlaId)
	if req.PriorityId != 0 {
		searchOptions.AddFilter("priority_id", req.PriorityId)
	}

	if req.Q != "" {
		searchOptions.AddFilter("name", req.Q)
	}

	slaConditions, e := s.app.Store.SLACondition().List(searchOptions)
	if e != nil {
		return nil, cerror.NewInternalError("sla_condition_service.list_sla_conditions.store.list.failed", e.Error())
	}

	return slaConditions, nil
}

// LocateSLACondition implements cases.SLAConditionsServer.
func (s *SLAConditionService) LocateSLACondition(ctx context.Context, req *cases.LocateSLAConditionRequest) (*cases.LocateSLAConditionResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.locate_sla_condition.id.required", "SLA Condition ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &cases.ListSLAConditionRequest{
		SlaId:  req.SlaId,
		Id:     []int64{req.Id},
		Fields: req.Fields,
	}

	// Call the ListSLAConditions method
	listResp, err := s.ListSLAConditions(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("sla_condition_service.locate_sla_condition.list_sla_conditions.error", err.Error())
	}

	// Check if the SLA Condition was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("sla_condition_service.locate_sla_condition.not_found", "SLA Condition not found")
	}

	// Return the found SLA Condition
	return &cases.LocateSLAConditionResponse{SlaCondition: listResp.Items[0]}, nil
}

// UpdateSLACondition implements cases.SLAConditionsServer.
func (s *SLAConditionService) UpdateSLACondition(ctx context.Context, req *cases.UpdateSLAConditionRequest) (*cases.SLACondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.update_sla_condition.id.required", "SLA Condition ID is required")
	}

	// Convert []*cases.Lookup to []int64
	var priorityIDs []int64
	for _, priority := range req.Input.Priorities {
		if priority != nil { // Check for nil to avoid runtime panic
			priorityIDs = append(priorityIDs, priority.GetId()) // Use GetId() to ensure proper handling
		}
	}

	// Define update options
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, SLAConditionMetadata),
		grpcopts.WithUpdateMasker(req),
		grpcopts.WithUpdateIDs(priorityIDs),
	)
	if err != nil {
		return nil, NewBadRequestError(err)
	}

	// Define the current user as the updater
	u := &cases.Lookup{
		Id: updateOpts.GetAuthOpts().GetUserId(),
	}

	// Update SLACondition user_auth
	slaCondition := &cases.SLACondition{
		Id:             req.Id,
		Name:           req.Input.Name,
		ReactionTime:   req.Input.ReactionTime,
		ResolutionTime: req.Input.ResolutionTime,
		SlaId:          req.SlaId,
		UpdatedBy:      u,
	}

	// Update the SLACondition in the store
	r, e := s.app.Store.SLACondition().Update(updateOpts, slaCondition)
	if e != nil {
		return nil, cerror.NewInternalError("sla_condition_service.update_sla_condition.store.update.failed", e.Error())
	}

	return r, nil
}

func NewSLAConditionService(app *App) (*SLAConditionService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_sla_condition_service.args_check.app_nil", "internal is nil")
	}
	return &SLAConditionService{app: app, objClassName: model.ScopeDictionary}, nil
}
