package grpc

import (
	"context"
	api "github.com/webitel/cases/api/cases"
	grpcopts "github.com/webitel/cases/internal/api_handler/grpc/options"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/util"
)

type ServiceHandler interface {
	CreateService(options.Creator, *model.Service) (*model.Service, error)
	ListServices(options.Searcher) ([]*model.Service, error)
	UpdateService(options.Updator, *model.Service) (*model.Service, error)
	DeleteService(options.Deleter) (*model.Service, error)
}

var ServiceMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "root_id", Default: true},
	{Name: "code", Default: true},
	{Name: "state", Default: true},
	{Name: "sla", Default: true},
	{Name: "group", Default: true},
	{Name: "assignee", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "catalog_id", Default: false},
})

type ServiceService struct {
	app ServiceHandler
	api.UnimplementedServicesServer
	objClassName string
}

// CreateService implements cases.ServicesServer.
func (s *ServiceService) CreateService(ctx context.Context, req *api.CreateServiceRequest) (*api.Service, error) {
	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, ServiceMetadata),
	)
	if err != nil {
		return nil, err
	}
	// Create a new Service user_session
	rootId := int(req.Input.RootId)
	catalogId := int(req.Input.CatalogId)
	service := &model.Service{
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
		Code:        &req.Input.Code,
		Sla:         utils.UnmarshalLookup(req.Input.Sla, &model.GeneralLookup{}),
		Group:       utils.UnmarshalExtendedLookup(req.Input.Group, &model.GeneralExtendedLookup{}),
		Assignee:    utils.UnmarshalLookup(req.Input.Assignee, &model.GeneralLookup{}),
		State:       &req.Input.State,
		RootId:      &rootId,
		CatalogId:   &catalogId,
	}

	// Create the Service in the store
	r, e := s.app.CreateService(createOpts, service)
	if e != nil {
		return nil, e
	}
	return s.Marshal(r)
}

// DeleteService implements cases.ServicesServer.
func (s *ServiceService) DeleteService(ctx context.Context, req *api.DeleteServiceRequest) (*api.ServiceList, error) {
	if len(req.Id) == 0 {
		return nil, errors.InvalidArgument("Service ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(ctx, grpcopts.WithDeleteIDs(req.Id))
	if err != nil {
		return nil, err
	}

	_, e := s.app.DeleteService(deleteOpts)
	if e != nil {
		return nil, e
	}

	deletedServices := make([]*api.Service, len(req.Id))
	for i, id := range req.Id {
		deletedServices[i] = &api.Service{Id: id}
	}

	return &api.ServiceList{
		Items: deletedServices,
	}, nil
}

// ListServices implements cases.ServicesServer.
func (s *ServiceService) ListServices(ctx context.Context, req *api.ListServiceRequest) (*api.ServiceList, error) {
	searchOptions, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, ServiceMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithIDs(req.Id),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, err
	}

	if req.Q != "" {
		searchOptions.AddFilter(util.EqualFilter("name", req.Q))
	}

	if req.RootId != 0 {
		searchOptions.AddFilter(util.EqualFilter("root_id", req.RootId))
	}

	if req.State {
		searchOptions.AddFilter(util.EqualFilter("state", req.State))
	}

	services, e := s.app.ListServices(searchOptions)
	if e != nil {
		return nil, e
	}
	var res api.ServiceList
	res.Items, err = utils.ConvertToOutputBulk(services, s.Marshal)
	if err != nil {
		return nil, err
	}
	res.Next, res.Items = utils.GetListResult(searchOptions, res.Items)
	res.Page = req.GetPage()
	return &res, nil
}

// LocateService implements cases.ServicesServer.
func (s *ServiceService) LocateService(ctx context.Context, req *api.LocateServiceRequest) (*api.LocateServiceResponse, error) {
	if req.Id == 0 {
		return nil, errors.InvalidArgument("Service ID is required")
	}
	opts, err := grpcopts.NewLocateOptions(ctx, grpcopts.WithID(req.Id), grpcopts.WithFields(req, ServiceMetadata, util.EnsureIdField, util.DeduplicateFields))
	if err != nil {
		return nil, err
	}
	listResp, err := s.app.ListServices(opts)
	if err != nil {
		return nil, err
	}

	if len(listResp) == 0 {
		return nil, errors.NotFound("service not found")
	}
	if len(listResp) > 1 {
		return nil, errors.InvalidArgument("multiple services found with the same ID")
	}
	var res api.LocateServiceResponse
	res.Service, err = s.Marshal(listResp[0])
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// UpdateService implements cases.ServicesServer.
func (s *ServiceService) UpdateService(ctx context.Context, req *api.UpdateServiceRequest) (*api.Service, error) {
	if req.Id == 0 {
		return nil, errors.InvalidArgument("Service ID is required")
	}

	updateOpts, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, ServiceMetadata),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, err
	}
	rootId := int(req.Input.RootId)
	service := &model.Service{
		Id:          int(req.Id),
		Name:        &req.Input.Name,
		Description: &req.Input.Description,
		Code:        &req.Input.Code,
		Sla:         utils.UnmarshalLookup(req.Input.Sla, &model.GeneralLookup{}),
		Group:       utils.UnmarshalExtendedLookup(req.Input.Group, &model.GeneralExtendedLookup{}),
		Assignee:    utils.UnmarshalLookup(req.Input.Assignee, &model.GeneralLookup{}),
		State:       &req.Input.State,
		RootId:      &rootId,
	}

	r, e := s.app.UpdateService(updateOpts, service)
	if e != nil {
		return nil, e
	}
	return s.Marshal(r)
}

func (s *ServiceService) Marshal(in *model.Service) (*api.Service, error) {
	if in == nil {
		return nil, nil
	}
	var resServices []*api.Service
	for _, service := range in.Services {
		if service != nil {
			marshaledService, err := s.Marshal(service)
			if err != nil {
				return nil, err
			}
			resServices = append(resServices, marshaledService)
		}
	}

	return &api.Service{
		Id:          int64(in.Id),
		Name:        utils.Dereference(in.Name),
		RootId:      int64(utils.Dereference(in.RootId)),
		Description: utils.Dereference(in.Description),
		Code:        utils.Dereference(in.Code),
		State:       utils.Dereference(in.State),
		Sla:         utils.MarshalLookup(in.Sla),
		Group:       utils.MarshalExtendedLookup(in.Group),
		Assignee:    utils.MarshalLookup(in.Assignee),
		CreatedAt:   utils.MarshalTime(in.CreatedAt),
		UpdatedAt:   utils.MarshalTime(in.UpdatedAt),
		CreatedBy:   utils.MarshalLookup(in.Author),
		UpdatedBy:   utils.MarshalLookup(in.Editor),
		CatalogId:   int64(utils.Dereference(in.CatalogId)),
		// Service and Searched fields can be set as needed
		Service:  resServices,
		Searched: utils.Dereference(in.Searched),
	}, nil
}

// NewServiceService creates a new ServiceService.
func NewServiceService(app ServiceHandler) (*ServiceService, error) {
	if app == nil {
		return nil, errors.Internal("internal is nil")
	}
	return &ServiceService{app: app, objClassName: model.ScopeDictionary}, nil
}
