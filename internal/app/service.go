package app

import (
	"context"
	api "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/util"
	"strings"
)

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
	app *App
	api.UnimplementedServicesServer
	objClassName string
}

// CreateService implements cases.ServicesServer.
func (s *ServiceService) CreateService(ctx context.Context, req *api.CreateServiceRequest) (*api.Service, error) {
	if req.Input.Name == "" {
		return nil, errors.InvalidArgument("Service name is required")
	}
	if req.Input.RootId == 0 {
		return nil, errors.InvalidArgument("Root ID is required")
	}

	createOpts, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, ServiceMetadata),
	)
	if err != nil {
		return nil, err
	}

	// Create a new Service user_session
	service := &api.Service{
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Code:        req.Input.Code,
		Sla:         req.Input.Sla,
		Group:       req.Input.Group,
		Assignee:    req.Input.Assignee,
		State:       req.Input.State,
		RootId:      req.Input.RootId,
		CatalogId:   req.Input.CatalogId,
	}

	// Create the Service in the store
	r, e := s.app.Store.Service().Create(createOpts, service)
	if e != nil {
		return nil, e
	}
	return r, nil
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

	e := s.app.Store.Service().Delete(deleteOpts)
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
		searchOptions.AddFilter("name", req.Q)
	}

	if req.RootId != 0 {
		searchOptions.AddFilter("root_id", req.RootId)
	}

	if req.State {
		searchOptions.AddFilter("state", req.State)
	}

	services, e := s.app.Store.Service().List(searchOptions)
	if e != nil {
		return nil, e
	}
	return services, nil
}

// LocateService implements cases.ServicesServer.
func (s *ServiceService) LocateService(ctx context.Context, req *api.LocateServiceRequest) (*api.LocateServiceResponse, error) {
	if req.Id == 0 {
		return nil, errors.InvalidArgument("Service ID is required")
	}

	if len(req.Fields) == 0 {
		req.Fields = strings.Split(defaultSubfields, ", ")
	} else {
		req.Fields = util.FieldsFunc(req.Fields, util.InlineFields)
	}

	if !util.ContainsField(req.Fields, "id") {
		req.Fields = append(req.Fields, "id")
	}

	listReq := &api.ListServiceRequest{
		Fields: req.Fields,
		Id:     []int64{req.Id},
		Page:   1,
		Size:   1,
	}

	listResp, err := s.ListServices(ctx, listReq)
	if err != nil {
		return nil, err
	}

	if len(listResp.Items) == 0 {
		return nil, errors.NotFound("Service not found")
	}

	return &api.LocateServiceResponse{Service: listResp.Items[0]}, nil
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

	service := &api.Service{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Code:        req.Input.Code,
		Sla:         req.Input.Sla,
		Group:       req.Input.Group,
		Assignee:    req.Input.Assignee,
		State:       req.Input.State,
		RootId:      req.Input.RootId,
	}

	r, e := s.app.Store.Service().Update(updateOpts, service)
	if e != nil {
		return nil, e
	}
	return r, nil
}

// NewServiceService creates a new ServiceService.
func NewServiceService(app *App) (*ServiceService, error) {
	if app == nil {
		return nil, errors.Internal("internal is nil")
	}
	return &ServiceService{app: app, objClassName: model.ScopeDictionary}, nil
}
