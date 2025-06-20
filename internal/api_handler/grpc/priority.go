package grpc

import (
	"context"
	"errors"

	api "github.com/webitel/cases/api/cases"
	grpcerror "github.com/webitel/cases/internal/api_handler/grpc/errors"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
)

type PriorityHandler interface {
	CreatePriority(options.Creator, *model.Priority) (*model.Priority, error)
	UpdatePriority(options.Updator, *model.Priority) (*model.Priority, error)
	DeletePriority(options.Deleter) (*model.Priority, error)
	ListPriorities(options.Searcher, int64, int64) ([]*model.Priority, error)
	LocatePriority(options.Searcher) (*model.Priority, error)
}

type PriorityService struct {
	app PriorityHandler
	api.UnimplementedPrioritiesServer
}

func NewPriorityService(app PriorityHandler) *PriorityService {
	return &PriorityService{app: app}
}

var PriorityMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "color", Default: true},
})

func (s *PriorityService) CreatePriority(ctx context.Context, req *api.CreatePriorityRequest) (*api.Priority, error) {
	if req.GetInput().GetName() == "" {
		return nil, grpcerror.NewBadRequestError(errors.New("priority name is required"))
	}
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, PriorityMetadata),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	input := &model.Priority{
		Name:        req.Input.Name,
		Description: &req.Input.Description,
		Color:       req.Input.Color,
	}
	m, err := s.app.CreatePriority(createOpts, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(m)
}

func (s *PriorityService) ListPriorities(ctx context.Context, req *api.ListPriorityRequest) (*api.PriorityList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, PriorityMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDs(req.GetId()),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	searchOpts.AddFilter("name", req.Q)

	items, err := s.app.ListPriorities(searchOpts, req.NotInSla, req.InSlaCond)
	if err != nil {
		return nil, err
	}
	var res api.PriorityList
	converted, err := utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}

	res.Next, res.Items = utils.GetListResult(searchOpts, converted)
	res.Page = req.GetPage()

	return &res, nil
}

func (s *PriorityService) UpdatePriority(ctx context.Context, req *api.UpdatePriorityRequest) (*api.Priority, error) {
	if req.GetId() == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("priority ID is required"))
	}
	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, PriorityMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	input := &model.Priority{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: &req.Input.Description,
		Color:       req.Input.Color,
	}
	m, err := s.app.UpdatePriority(updateOpts, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(m)
}

func (s *PriorityService) DeletePriority(ctx context.Context, req *api.DeletePriorityRequest) (*api.Priority, error) {
	if req.GetId() == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("priority ID is required"))
	}
	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	item, err := s.app.DeletePriority(deleteOpts)
	if err != nil {
		return nil, err
	}
	return s.Marshal(item)
}

func (s *PriorityService) LocatePriority(ctx context.Context, req *api.LocatePriorityRequest) (*api.LocatePriorityResponse, error) {
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("priority ID is required"))
	}
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithFields(req, PriorityMetadata,
		util.DeduplicateFields,
		util.EnsureIdField,
	), grpcopts.WithID(req.Id))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	item, err := s.app.LocatePriority(opts)
	if err != nil {
		return nil, err
	}
	res, err := s.Marshal(item)
	if err != nil {
		return nil, err
	}
	return &api.LocatePriorityResponse{Priority: res}, nil
}

func (s *PriorityService) Marshal(model *model.Priority) (*api.Priority, error) {
	if model == nil {
		return nil, nil
	}
	return &api.Priority{
		Id:          model.Id,
		Name:        model.Name,
		Description: utils.Dereference(model.Description),
		Color:       model.Color,
		CreatedAt:   utils.MarshalTime(model.CreatedAt),
		UpdatedAt:   utils.MarshalTime(model.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(model.Author),
		UpdatedBy:   utils.MarshalLookup(model.Editor),
	}, nil
}
