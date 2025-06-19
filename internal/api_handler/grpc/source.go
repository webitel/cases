package grpc

import (
	"context"
	"errors"
	grpcerror "github.com/webitel/cases/internal/api_handler/grpc/errors"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/model/options"
	"strings"

	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
)

type SourceHandler interface {
	CreateSource(options.Creator, *model.Source) (*model.Source, error)
	UpdateSource(options.Updator, *model.Source) (*model.Source, error)
	DeleteSource(options.Deleter) (*model.Source, error)
	ListSources(options.Searcher) ([]*model.Source, error)
}

type SourceService struct {
	app SourceHandler
	_go.UnimplementedSourcesServer
	objClassName string
}

var SourceMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
	{"type", true},
})

// CreateSource implements api.SourcesServer.
func (s *SourceService) CreateSource(
	ctx context.Context,
	req *_go.CreateSourceRequest,
) (*_go.Source, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("source_service.create_source.name.required", "Source name is required")
	}
	if req.Input.Type == _go.SourceType_TYPE_UNSPECIFIED {
		return nil, cerror.NewBadRequestError("source_service.create_source.type.required", "Source type is required")
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, SourceMetadata),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	sourceType := req.Input.Type.String()
	input := &model.Source{
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
		Type:        &sourceType,
	}

	// Create the source in the store
	res, err := s.app.CreateSource(createOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.create_source.store.create.failed", err.Error())
	}

	return s.Marshal(res)
}

// ListSources implements api.SourcesServer.
func (s *SourceService) ListSources(
	ctx context.Context,
	req *_go.ListSourceRequest,
) (*_go.SourceList, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, SourceMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
		grpcopts.WithIDs(req.GetId()),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	if req.Q != "" {
		searchOpts.AddFilter("name", req.Q)
	}
	if len(req.Type) > 0 {
		searchOpts.AddFilter("type", req.Type)
	}

	items, err := s.app.ListSources(searchOpts)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.list_sources.store.list.failed", err.Error())
	}

	var res _go.SourceList
	res.Items, err = utils.ConvertToOutputBulk(items, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOpts, res.Items)
	res.Page = int32(searchOpts.GetPage())
	return &res, nil
}

// UpdateSource implements api.SourcesServer.
func (s *SourceService) UpdateSource(
	ctx context.Context,
	req *_go.UpdateSourceRequest,
) (*_go.Source, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.update_source.id.required", "Source ID is required")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, SourceMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	sourceType := req.Input.Type.String()
	input := &model.Source{
		Id:          int(req.Id),
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
		Type:        &sourceType,
	}

	// Update the source in the store
	res, err := s.app.UpdateSource(updateOpts, input)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.update_source.store.update.failed", err.Error())
	}

	return s.Marshal(res)
}

// DeleteSource implements api.SourcesServer.
func (s *SourceService) DeleteSource(
	ctx context.Context,
	req *_go.DeleteSourceRequest,
) (*_go.Source, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.delete_source.id.required", "Source ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteID(req.Id))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}

	deleteOpts.IDs = []int64{req.Id}

	// Delete the source in the store
	_, err = s.app.DeleteSource(deleteOpts)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.delete_source.store.delete.failed", err.Error())
	}

	return &(_go.Source{Id: req.Id}), nil
}

// LocateSource implements api.SourcesServer.
func (s *SourceService) LocateSource(ctx context.Context, req *_go.LocateSourceRequest) (*_go.LocateSourceResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("source ID is required"))
	}
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithID(req.Id), grpcopts.WithFields(req, StatusConditionMetadata, util.EnsureIdField, util.DeduplicateFields))
	if err != nil {
		return nil, grpcerror.NewBadRequestError(err)
	}
	// Call the ListStatusConditions method
	items, err := s.app.ListSources(opts)
	if err != nil {
		return nil, err
	}

	// Check if the status condition was found
	if len(items) > 1 {
		return nil, grpcerror.NewBadRequestError(errors.New("multiple rows found"))
	}
	if len(items) == 0 {
		return nil, grpcerror.NewBadRequestError(errors.New("not found"))
	}

	// Return the found status condition
	var res _go.LocateSourceResponse
	res.Source, err = s.Marshal(items[0])
	if err != nil {
		return nil, grpcerror.ConversionError
	}
	return &res, nil
}

func (s *SourceService) Marshal(in *model.Source) (*_go.Source, error) {
	return &_go.Source{
		Id:          int64(in.Id),
		Name:        utils.Dereference(in.Name),
		Description: utils.Dereference(in.Description),
		Type:        stringToType(utils.Dereference(in.Type)),
		CreatedAt:   utils.MarshalTime(in.CreatedAt),
		UpdatedAt:   utils.MarshalTime(in.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(in.Author),
		UpdatedBy:   utils.MarshalLookup(in.Author),
	}, nil
}

func NewSourceService(app SourceHandler) (*SourceService, error) {
	if app == nil {
		return nil, errors.New("source handler is nil")
	}
	return &SourceService{app: app, objClassName: model.ScopeDictionary}, nil
}

// StringToType converts a string into the corresponding Type enum value.
//
// Types are specified ONLY for Source dictionary and are ENUMS in API.
func stringToType(typeStr string) _go.SourceType {
	switch strings.ToUpper(typeStr) {
	case "CALL":
		return _go.SourceType_CALL
	case "CHAT":
		return _go.SourceType_CHAT
	case "SOCIAL_MEDIA":
		return _go.SourceType_SOCIAL_MEDIA
	case "EMAIL":
		return _go.SourceType_EMAIL
	case "API":
		return _go.SourceType_API
	case "MANUAL":
		return _go.SourceType_MANUAL
	default:
		return _go.SourceType_TYPE_UNSPECIFIED
	}
}
