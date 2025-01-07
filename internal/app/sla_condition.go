package app

import (
	"context"
	"strings"
	"time"

	cases "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"

	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
)

type SLAConditionService struct {
	app *App
	cases.UnimplementedSLAConditionsServer
}

const (
	defaultFieldsSLACondition = "id, name, priority"
)

// CreateSLACondition implements cases.SLAConditionsServer.
func (s *SLAConditionService) CreateSLACondition(ctx context.Context, req *cases.CreateSLAConditionRequest) (*cases.SLACondition, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.name.required", "SLA Condition name is required")
	}
	if len(req.Priorities) == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.priorities.required", "At least one priority is required")
	}
	if req.ReactionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.reaction_time.required", "Reaction time is required")
	}
	if req.ResolutionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.resolution_time.required", "Resolution time is required")
	}
	if req.SlaId == 0 {
		return nil, cerror.NewBadRequestError("sla_condition_service.create_sla_condition.sla_id.required", "SLA ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_condition_service.create_sla_condition.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the creator and updater
	currentU := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new SLACondition model
	slaCondition := &cases.SLACondition{
		Name:           req.Name,
		ReactionTime:   req.ReactionTime,
		ResolutionTime: req.ResolutionTime,
		SlaId:          req.SlaId,
		CreatedBy:      currentU,
		UpdatedBy:      currentU,
	}

	fields := []string{
		"id", "name", "reaction_time", "resolution_time", "sla_id",
		"created_at", "updated_at", "created_by", "updated_by",
	}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
		Ids:     req.Priorities,
	}

	// Create the SLACondition in the store
	r, e := s.app.Store.SLACondition().Create(&createOpts, slaCondition, req.Priorities)
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

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_condition_service.delete_sla_condition.authorization.failed", err.Error())
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

	// Delete the SLACondition in the store
	e := s.app.Store.SLACondition().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("sla_condition_service.delete_sla_condition.store.delete.failed", e.Error())
	}

	return &cases.SLACondition{Id: req.Id}, nil
}

// ListSLAConditions implements cases.SLAConditionsServer.
func (s *SLAConditionService) ListSLAConditions(ctx context.Context, req *cases.ListSLAConditionRequest) (*cases.SLAConditionList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_condition_service.list_sla_conditions.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsSLACondition, ", ")
	}

	// Use default page size and page number if not provided
	page := req.Page
	if page == 0 {
		page = 1
	}

	t := time.Now()
	searchOptions := model.SearchOptions{
		ParentId: req.SlaId,
		IDs:      req.Id,
		//Session:  session,
		Fields:  fields,
		Context: ctx,
		Sort:    req.Sort,
		Page:    int(page),
		Size:    int(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
		ID:      req.PriorityId,
		Auth:    model.NewDefaultAuthOptions(session, "dictionaries"),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	slaConditions, e := s.app.Store.SLACondition().List(&searchOptions)
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
		Page:   1,
		Size:   1, // We only need one item
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

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("sla_condition_service.update_sla_condition.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the updater
	u := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Update SLACondition model
	slaCondition := &cases.SLACondition{
		Id:             req.Id,
		Name:           req.Input.Name,
		ReactionTime:   req.Input.ReactionTime,
		ResolutionTime: req.Input.ResolutionTime,
		SlaId:          req.Input.SlaId,
		UpdatedBy:      u,
	}

	fields := []string{"id"}

	// Map XJsonMask fields to the corresponding SLACondition fields
	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("sla_condition_service.update_sla_condition.name.required", "SLA Condition name is required and cannot be empty")
			}
		case "reaction_time":
			fields = append(fields, "reaction_time")
		case "resolution_time":
			fields = append(fields, "resolution_time")
		case "sla_id":
			fields = append(fields, "sla_id")
		case "priorities":
			fields = append(fields, "priorities")
		}
	}

	t := time.Now()

	// Define update options
	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
		IDs:     req.Input.Priorities,
	}

	// Update the SLACondition in the store
	r, e := s.app.Store.SLACondition().Update(&updateOpts, slaCondition)
	if e != nil {
		return nil, cerror.NewInternalError("sla_condition_service.update_sla_condition.store.update.failed", e.Error())
	}

	return r, nil
}

func NewSLAConditionService(app *App) (*SLAConditionService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_sla_condition_service.args_check.app_nil", "internal is nil")
	}
	return &SLAConditionService{app: app}, nil
}
