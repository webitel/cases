package app

import (
	"context"
	"strings"
	"time"

	api "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/cases/util"
)

type ServiceService struct {
	app *App
	api.UnimplementedServicesServer
	objClassName string
}

// CreateService implements cases.ServicesServer.
func (s *ServiceService) CreateService(ctx context.Context, req *api.CreateServiceRequest) (*api.Service, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("service.create_service.name.required", "Service name is required")
	}

	if req.RootId == 0 {
		return nil, cerror.NewBadRequestError("service.create_service.root_id.required", "Root ID is required")
	}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Time:    t,
	}

	// Define the current user as the creator and updater
	currentU := &api.Lookup{
		Id: createOpts.GetAuthOpts().GetUserId(),
	}

	// Create a new Service user_auth
	service := &api.Service{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Sla:         req.Sla,
		Group:       req.Group,
		Assignee:    req.Assignee,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
		State:       req.State,
		RootId:      req.RootId,
		CatalogId:   req.CatalogId,
	}

	// Create the Service in the store
	r, e := s.app.Store.Service().Create(&createOpts, service)
	if e != nil {
		return nil, cerror.NewInternalError("service.create_service.store.create.failed", e.Error())
	}

	return r, nil
}

// DeleteService implements cases.ServicesServer.
func (s *ServiceService) DeleteService(ctx context.Context, req *api.DeleteServiceRequest) (*api.ServiceList, error) {
	if len(req.Id) == 0 {
		return nil, cerror.NewBadRequestError("service.delete_service.id.required", "Service ID is required")
	}

	t := time.Now()
	deleteOpts := model.DeleteOptions{
		Context: ctx,
		IDs:     req.Id,
		Time:    t,
		Auth:    model.GetAutherOutOfContext(ctx),
	}

	e := s.app.Store.Service().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("service.delete_service.store.delete.failed", e.Error())
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
	page := req.Page
	if page == 0 {
		page = 1
	}

	if len(req.Fields) == 0 {
		req.Fields = strings.Split(defaultSubfields, ", ")
	} else {
		req.Fields = util.FieldsFunc(req.Fields, util.InlineFields)
	}

	if !util.ContainsField(req.Fields, "id") {
		req.Fields = append(req.Fields, "id")
	}

	t := time.Now()
	searchOptions := &model.SearchOptions{
		Fields: req.Fields,
		IDs:    req.Id,
		// UserAuthSession: session,
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

	if req.RootId != 0 {
		searchOptions.Filter["root_id"] = req.RootId
	}

	if req.State {
		searchOptions.Filter["state"] = req.State
	}

	services, e := s.app.Store.Service().List(searchOptions)
	if e != nil {
		return nil, cerror.NewInternalError("service.list_services.store.list.failed", e.Error())
	}

	return services, nil
}

// LocateService implements cases.ServicesServer.
func (s *ServiceService) LocateService(ctx context.Context, req *api.LocateServiceRequest) (*api.LocateServiceResponse, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("service.locate_service.id.required", "Service ID is required")
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
		return nil, cerror.NewInternalError("service.locate_service.list_services.error", err.Error())
	}

	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("service.locate_service.not_found", "Service not found")
	}

	return &api.LocateServiceResponse{Service: listResp.Items[0]}, nil
}

// UpdateService implements cases.ServicesServer.
func (s *ServiceService) UpdateService(ctx context.Context, req *api.UpdateServiceRequest) (*api.Service, error) {
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("service.update_service.id.required", "Service ID is required")
	}

	fields := []string{"id"}

	for _, f := range req.XJsonMask {
		// Handle fields with specific prefixes
		if strings.HasPrefix(f, "sla") {
			if !util.ContainsField(fields, "sla_id") {
				fields = append(fields, "sla_id")
			}
			continue
		}

		if strings.HasPrefix(f, "group") {
			if !util.ContainsField(fields, "group_id") {
				fields = append(fields, "group_id")
			}
			continue
		}

		if strings.HasPrefix(f, "assignee") {
			if !util.ContainsField(fields, "assignee_id") {
				fields = append(fields, "assignee_id")
			}
			continue
		}

		// Handle exact matches
		switch f {
		case "name":
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("service.update_service.name.required", "Service name is required and cannot be empty")
			}
		case "description":
			fields = append(fields, "description")
		case "root_id":
			fields = append(fields, "root_id")
		case "code":
			fields = append(fields, "code")
		case "state":
			fields = append(fields, "state")
		}
	}

	t := time.Now()

	updateOpts := model.UpdateOptions{
		Context: ctx,
		Fields:  fields,
		Time:    t,
		Auth:    model.GetAutherOutOfContext(ctx),
	}

	u := &api.Lookup{
		Id: updateOpts.GetAuthOpts().GetUserId(),
	}

	service := &api.Service{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Code:        req.Input.Code,
		Sla:         req.Input.Sla,
		Group:       req.Input.Group,
		Assignee:    req.Input.Assignee,
		UpdatedBy:   u,
		State:       req.Input.State,
		RootId:      req.Input.RootId,
	}

	r, e := s.app.Store.Service().Update(&updateOpts, service)
	if e != nil {
		return nil, cerror.NewInternalError("service.update_service.store.update.failed", e.Error())
	}

	return r, nil
}

// NewServiceService creates a new ServiceService.
func NewServiceService(app *App) (*ServiceService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_service.args_check.app_nil", "internal is nil")
	}
	return &ServiceService{app: app, objClassName: model.ScopeDictionary}, nil
}
