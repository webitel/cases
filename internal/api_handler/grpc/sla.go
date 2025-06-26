package grpc

import (
	"context"
	"google.golang.org/grpc/codes"

	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
)

type SLAHandler interface {
	ListSLAs(options.Searcher) ([]*model.SLA, error)
	CreateSLA(options.Creator, *model.SLA) (*model.SLA, error)
	UpdateSLA(options.Updator, *model.SLA) (*model.SLA, error)
	DeleteSLA(options.Deleter) (*model.SLA, error)
}

type SLAService struct {
	app SLAHandler
	cases.UnimplementedSLAsServer
	objClassName string
}

func NewSLAService(app SLAHandler) (*SLAService, error) {
	if app == nil {
		return nil, errors.New("sla handler is nil")
	}
	return &SLAService{app: app, objClassName: model.ScopeDictionary}, nil
}

var SlaMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: true},
	{Name: "updated_at", Default: true},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "valid_from", Default: true},
	{Name: "valid_to", Default: true},
	{Name: "calendar", Default: true},
	{Name: "reaction_time", Default: true},
	{Name: "resolution_time", Default: true},
})

func (s *SLAService) CreateSLA(
	ctx context.Context,
	req *cases.CreateSLARequest,
) (*cases.SLA, error) {
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, SlaMetadata),
	)
	if err != nil {
		return nil, err
	}

	input := &model.SLA{
		Name:           &req.Input.Name,
		Description:    &req.Input.Description,
		ValidFrom:      utils.TimePtr(req.Input.ValidFrom),
		ValidTo:        utils.TimePtr(req.Input.ValidTo),
		Calendar:       utils.UnmarshalLookup(req.Input.Calendar, &model.Calendar{}),
		ReactionTime:   int(req.Input.ReactionTime),
		ResolutionTime: int(req.Input.ResolutionTime),
	}

	m, err := s.app.CreateSLA(createOpts, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(m)
}

func (s *SLAService) ListSLAs(
	ctx context.Context,
	req *cases.ListSLARequest,
) (*cases.SLAList, error) {
	searcher, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, SlaMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDs(req.GetId()),
	)
	if err != nil {
		return nil, err
	}
	searcher.AddFilter("name", req.GetQ())

	items, err := s.app.ListSLAs(searcher)
	if err != nil {
		return nil, err
	}

	var res cases.SLAList
	converted, err := utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searcher, converted)
	res.Page = req.GetPage()

	return &res, nil
}

func (s *SLAService) UpdateSLA(
	ctx context.Context,
	req *cases.UpdateSLARequest,
) (*cases.SLA, error) {
	if req.GetId() == 0 {
		return nil, errors.New("SLA ID is required", errors.WithCode(codes.InvalidArgument))
	}

	updator, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, SlaMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, err
	}

	input := &model.SLA{
		Id:             int(req.Id),
		Name:           &req.Input.Name,
		Description:    &req.Input.Description,
		ValidFrom:      utils.TimePtr(req.Input.ValidFrom),
		ValidTo:        utils.TimePtr(req.Input.ValidTo),
		Calendar:       utils.UnmarshalLookup(req.Input.Calendar, &model.Calendar{}),
		ReactionTime:   int(req.Input.ReactionTime),
		ResolutionTime: int(req.Input.ResolutionTime),
	}

	updated, err := s.app.UpdateSLA(updator, input)
	if err != nil {
		return nil, err
	}
	return s.Marshal(updated)
}

func (s *SLAService) DeleteSLA(
	ctx context.Context,
	req *cases.DeleteSLARequest,
) (*cases.SLA, error) {
	if req.GetId() == 0 {
		return nil, errors.New("SLA ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(
		ctx,
		grpcopts.WithDeleteID(req.Id),
	)
	if err != nil {
		return nil, err
	}

	item, err := s.app.DeleteSLA(deleteOpts)
	if err != nil {
		return nil, err
	}
	return s.Marshal(item)
}

func (s *SLAService) LocateSLA(
	ctx context.Context,
	req *cases.LocateSLARequest,
) (*cases.LocateSLAResponse, error) {
	if req.Id == 0 {
		return nil, errors.New("SLA ID is required", errors.WithCode(codes.InvalidArgument))
	}

	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithFields(req, SlaMetadata,
		util.DeduplicateFields,
		util.EnsureIdField,
	), grpcopts.WithID(req.Id))
	if err != nil {
		return nil, err
	}

	items, err := s.app.ListSLAs(opts)
	if err != nil {
		return nil, err
	}

	res, err := s.Marshal(items[0])
	if err != nil {
		return nil, err
	}
	return &cases.LocateSLAResponse{Sla: res}, nil
}

func (s *SLAService) Marshal(in *model.SLA) (*cases.SLA, error) {
	if in == nil {
		return nil, nil
	}
	return &cases.SLA{
		Id:             int64(in.Id),
		Name:           utils.Dereference(in.Name),
		Description:    utils.Dereference(in.Description),
		ValidFrom:      utils.MarshalTime(in.ValidFrom),
		ValidTo:        utils.MarshalTime(in.ValidTo),
		Calendar:       utils.MarshalLookup(in.Calendar),
		ReactionTime:   int64(in.ReactionTime),
		ResolutionTime: int64(in.ResolutionTime),
		CreatedAt:      utils.MarshalTime(in.CreatedAt),
		UpdatedAt:      utils.MarshalTime(in.UpdatedAt),
		CreatedBy:      utils.MarshalLookup(in.Author),
		UpdatedBy:      utils.MarshalLookup(in.Editor),
	}, nil
}
