package app

import (
	"context"
	"time"

	"github.com/webitel/cases/api/cases"
	authmodel "github.com/webitel/cases/auth/model"
	"github.com/webitel/cases/model"
)

type ServiceService struct {
	app *App
}

// CreateService implements cases.ServicesServer.
func (s *ServiceService) CreateService(ctx context.Context, req *cases.CreateServiceRequest) (*cases.Service, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, model.NewBadRequestError("service.create_service.name.required", "Service name is required")
	}

	if req.RootId == 0 {
		return nil, model.NewBadRequestError("service.create_service.root_id.required", "Root ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("service.create_service.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Add
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	// Define the current user as the creator and updater
	currentU := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	// Create a new Service model
	service := &cases.Service{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		Sla:         &cases.Lookup{Id: req.SlaId},
		Group:       &cases.Lookup{Id: req.GroupId},
		Assignee:    &cases.Lookup{Id: req.AssigneeId}, // Added Assignee field
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
		return nil, model.NewInternalError("service.create_service.store.create.failed", e.Error())
	}

	return r, nil
}

// DeleteService implements cases.ServicesServer.
func (s *ServiceService) DeleteService(ctx context.Context, req *cases.DeleteServiceRequest) (*cases.ServiceList, error) {
	if len(req.Id) == 0 {
		return nil, model.NewBadRequestError("service.delete_service.id.required", "Service ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("service.delete_service.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Delete
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
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
		return nil, model.NewInternalError("service.delete_service.store.delete.failed", e.Error())
	}

	deletedServices := make([]*cases.Service, len(req.Id))
	for i, id := range req.Id {
		deletedServices[i] = &cases.Service{Id: id}
	}

	return &cases.ServiceList{
		Items: deletedServices,
	}, nil
}

// ListServices implements cases.ServicesServer.
func (s *ServiceService) ListServices(ctx context.Context, req *cases.ListServiceRequest) (*cases.ServiceList, error) {
	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("service.list_services.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Read
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
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
		Page:    int(page),
		Size:    int(req.Size),
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
		return nil, model.NewInternalError("service.list_services.store.list.failed", e.Error())
	}

	return services, nil
}

// LocateService implements cases.ServicesServer.
func (s *ServiceService) LocateService(ctx context.Context, req *cases.LocateServiceRequest) (*cases.LocateServiceResponse, error) {
	if req.Id == 0 {
		return nil, model.NewBadRequestError("service.locate_service.id.required", "Service ID is required")
	}

	listReq := &cases.ListServiceRequest{
		Id:   []int64{req.Id},
		Page: 1,
		Size: 1,
	}

	listResp, err := s.ListServices(ctx, listReq)
	if err != nil {
		return nil, model.NewInternalError("service.locate_service.list_services.error", err.Error())
	}

	if len(listResp.Items) == 0 {
		return nil, model.NewNotFoundError("service.locate_service.not_found", "Service not found")
	}

	return &cases.LocateServiceResponse{Service: listResp.Items[0]}, nil
}

// UpdateService implements cases.ServicesServer.
func (s *ServiceService) UpdateService(ctx context.Context, req *cases.UpdateServiceRequest) (*cases.Service, error) {
	if req.Id == 0 {
		return nil, model.NewBadRequestError("service.update_service.id.required", "Service ID is required")
	}

	session, err := s.app.AuthorizeFromContext(ctx)
	if err != nil {
		return nil, model.NewUnauthorizedError("service.update_service.authorization.failed", err.Error())
	}

	// OBAC check
	accessMode := authmodel.Edit
	scope := session.GetScope(model.ScopeDictionary)
	if !session.HasObacAccess(scope.Class, accessMode) {
		return nil, s.app.MakeScopeError(session, scope, accessMode)
	}

	u := &cases.Lookup{
		Id:   session.GetUserId(),
		Name: session.GetUserName(),
	}

	service := &cases.Service{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Code:        req.Input.Code,
		Sla:         &cases.Lookup{Id: req.Input.SlaId},
		Group:       &cases.Lookup{Id: req.Input.GroupId},
		Assignee:    &cases.Lookup{Id: req.Input.AssigneeId},
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
				return nil, model.NewBadRequestError("service.update_service.name.required", "Service name is required and cannot be empty")
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
		return nil, model.NewInternalError("service.update_service.store.update.failed", e.Error())
	}

	return r, nil
}

// NewServiceService creates a new ServiceService.
func NewServiceService(app *App) (*ServiceService, model.AppError) {
	if app == nil {
		return nil, model.NewInternalError("api.config.new_service.args_check.app_nil", "internal is nil")
	}
	return &ServiceService{app: app}, nil
}
