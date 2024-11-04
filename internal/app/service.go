package app

import (
	"context"
	"time"

	api "github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"

	cerror "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
)

type ServiceService struct {
	app *App
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

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("service.create_service.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	// Define the current user as the creator and updater
	currentU := &api.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new Service model
	service := &api.Service{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Sla:         &api.Lookup{Id: req.SlaId},
		Group:       &api.Lookup{Id: req.GroupId},
		Assignee:    &api.Lookup{Id: req.AssigneeId},
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
		State:       req.State,
		RootId:      req.RootId,
		CatalogId:   req.CatalogId,
	}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Session: session,
		Context: ctx,
		Time:    t,
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

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("service.delete_service.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	t := time.Now()
	deleteOpts := model.DeleteOptions{
		Session: session,
		Context: ctx,
		IDs:     req.Id,
		Time:    t,
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
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("service.list_services.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	page := req.Page
	if page == 0 {
		page = 1
	}

	t := time.Now()
	searchOptions := model.SearchOptions{
		IDs:     req.Id,
		Session: session,
		Context: ctx,
		Sort:    req.Sort,
		Page:    int32(page),
		Size:    int32(req.Size),
		Time:    t,
		Filter:  make(map[string]interface{}),
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

	services, e := s.app.Store.Service().List(&searchOptions)
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

	listReq := &api.ListServiceRequest{
		Id:   []int64{req.Id},
		Page: 1,
		Size: 1,
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

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, cerror.NewUnauthorizedError("service.update_service.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, cerror.MakeScopeError(session.GetUserId(), scope.Class, int(accessMode))
	}

	u := &api.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	service := &api.Service{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Code:        req.Input.Code,
		Sla:         &api.Lookup{Id: req.Input.SlaId},
		Group:       &api.Lookup{Id: req.Input.GroupId},
		Assignee:    &api.Lookup{Id: req.Input.AssigneeId},
		UpdatedBy:   u,
		State:       req.Input.State,
		RootId:      req.Input.RootId,
	}

	fields := []string{"id"}

	for _, f := range req.XJsonMask {
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
		case "sla_id":
			fields = append(fields, "sla_id")
		case "group_id":
			fields = append(fields, "group_id")
		case "assignee_id":
			fields = append(fields, "assignee_id")
		case "state":
			fields = append(fields, "state")
		}
	}

	t := time.Now()

	updateOpts := model.UpdateOptions{
		Session: session,
		Context: ctx,
		Fields:  fields,
		Time:    t,
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
	return &ServiceService{app: app}, nil
}
