package grpc

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"google.golang.org/grpc/codes"
)

// SLAConditionHandler defines the interface for managing SLA conditions.
type SLAConditionHandler interface {
	CreateSLACondition(options.Creator, *model.SLACondition) (*model.SLACondition, error)
	UpdateSLACondition(options.Updator, *model.SLACondition) (*model.SLACondition, error)
	DeleteSLACondition(options.Deleter) (*model.SLACondition, error)
	ListSLAConditions(options.Searcher) ([]*model.SLACondition, error)
}

// SLAConditionService implements the gRPC server for SLA conditions.
type SLAConditionService struct {
	app SLAConditionHandler
	cases.UnimplementedSLAConditionsServer
	objClassName string
}

// SLAConditionMetadata defines the fields available for SLA condition objects.
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

// CreateSLACondition handles the gRPC request to create a new SLA condition.
// It validates the request, creates a new SLACondition model, and calls the
// handler's CreateSLACondition method.
func (s *SLAConditionService) CreateSLACondition(ctx context.Context, req *cases.CreateSLAConditionRequest) (*cases.SLACondition, error) {
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, SLAConditionMetadata),
	)
	if err != nil {
		return nil, err
	}

	// Create a new SLACondition user_session
	reactionTime := int(req.Input.ReactionTime)
	resolutionTime := int(req.Input.ResolutionTime)
	slaId := int64(req.SlaId)
	slaCondition := &model.SLACondition{
		Name:           &req.Input.Name,
		ReactionTime:   &reactionTime,
		ResolutionTime: &resolutionTime,
		SlaId:          &slaId,
	}
	for _, priority := range req.Input.Priorities {
		if priority != nil { // Check for nil to avoid runtime panic
			slaCondition.Priorities = append(slaCondition.Priorities, &model.Priority{Id: priority.GetId(), Name: priority.GetName()})
		}
	}

	// Create the SLACondition in the store
	r, err := s.app.CreateSLACondition(createOpts, slaCondition)
	if err != nil {
		return nil, err
	}

	return s.Marshal(r)
}

// DeleteSLACondition handles the gRPC request to delete an SLA condition.
// It validates the request and calls the handler's DeleteSLACondition method.
func (s *SLAConditionService) DeleteSLACondition(ctx context.Context, req *cases.DeleteSLAConditionRequest) (*cases.SLACondition, error) {
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	deleteOpts.Fields = SLAConditionMetadata.GetAllFields()
	if err != nil {
		return nil, err
	}

	// Delete the SLACondition in the store
	res, err := s.app.DeleteSLACondition(deleteOpts)
	if err != nil {
		return nil, err
	}

	return s.Marshal(res)
}

// ListSLAConditions handles the gRPC request to list SLA conditions with filters and pagination.
// It constructs search options based on the request and calls the handler's ListSLAConditions method.
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
		return nil, err
	}
	if req.SlaId != 0 {
		searchOptions.AddFilter(util.EqualFilter("sla_id", int(req.SlaId)))

	}
	if req.PriorityId != 0 {
		searchOptions.AddFilter(util.EqualFilter("priority_id", req.PriorityId))
	}

	if req.Q != "" {
		searchOptions.AddFilter(util.EqualFilter("name", req.Q))
	}

	slaConditions, err := s.app.ListSLAConditions(searchOptions)
	if err != nil {
		return nil, err
	}
	var res cases.SLAConditionList
	res.Items, err = utils.ConvertToOutputBulk(slaConditions, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOptions, res.Items)
	res.Page = int32(searchOptions.GetPage())
	return &res, nil
}

// LocateSLACondition finds an SLA condition by ID and returns it, or an error if not found.
// It uses the ListSLAConditions method to retrieve the SLA condition.
func (s *SLAConditionService) LocateSLACondition(ctx context.Context, req *cases.LocateSLAConditionRequest) (*cases.LocateSLAConditionResponse, error) {
	// Prepare a list request with necessary parameters
	listReq := &cases.ListSLAConditionRequest{
		SlaId:  req.SlaId,
		Id:     []int64{req.Id},
		Fields: req.Fields,
	}

	// Call the ListSLAConditions method
	listResp, err := s.ListSLAConditions(ctx, listReq)
	if err != nil {
		return nil, err
	}

	// Check if the SLA Condition was found
	if len(listResp.Items) == 0 {
		return nil, errors.New("SLA Condition not found", errors.WithCode(codes.NotFound))
	}

	// Return the found SLA Condition
	return &cases.LocateSLAConditionResponse{SlaCondition: listResp.Items[0]}, nil
}

// UpdateSLACondition handles the gRPC request to update an existing SLA condition.
// It validates the request, constructs an SLACondition model, and calls the
// handler's UpdateSLACondition method.
func (s *SLAConditionService) UpdateSLACondition(ctx context.Context, req *cases.UpdateSLAConditionRequest) (*cases.SLACondition, error) {
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
		return nil, err
	}
	// Update SLACondition user_session
	reactionTime := int(req.Input.ReactionTime)
	resolutionTime := int(req.Input.ResolutionTime)
	slaId := int64(req.SlaId)
	slaCondition := &model.SLACondition{
		Id:             int64(req.Id),
		Name:           &req.Input.Name,
		ReactionTime:   &reactionTime,
		ResolutionTime: &resolutionTime,
		SlaId:          &slaId,
	}
	for _, priority := range req.Input.Priorities {
		if priority != nil { // Check for nil to avoid runtime panic
			slaCondition.Priorities = append(slaCondition.Priorities, &model.Priority{Id: priority.GetId(), Name: priority.GetName()})
		}
	}

	// Update the SLACondition in the store
	item, err := s.app.UpdateSLACondition(updateOpts, slaCondition)
	if err != nil {
		return nil, err
	}

	return s.Marshal(item)
}

// Marshal converts a model.SLACondition to its gRPC representation.
func (s *SLAConditionService) Marshal(in *model.SLACondition) (*cases.SLACondition, error) {
	if in == nil {
		return nil, nil
	}
	res := &cases.SLACondition{
		Id:             int64(in.Id),
		Name:           utils.Dereference(in.Name),
		ReactionTime:   int64(utils.Dereference(in.ReactionTime)),
		ResolutionTime: int64(utils.Dereference(in.ResolutionTime)),
		SlaId:          int64(utils.Dereference(in.SlaId)),
		CreatedAt:      utils.MarshalTime(in.CreatedAt),
		UpdatedAt:      utils.MarshalTime(in.UpdatedAt),
		CreatedBy:      utils.MarshalLookup(in.Author),
		UpdatedBy:      utils.MarshalLookup(in.Editor),
	}
	for _, v := range in.Priorities {
		res.Priorities = append(res.Priorities, &cases.Lookup{
			Id:   v.Id,
			Name: v.Name,
		})
	}
	return res, nil
}

// NewSLAConditionService constructs a new SLAConditionService.
func NewSLAConditionService(app SLAConditionHandler) (*SLAConditionService, error) {
	if app == nil {
		return nil, errors.New("sla condition handler is nil")
	}
	return &SLAConditionService{app: app, objClassName: model.ScopeDictionary}, nil
}
