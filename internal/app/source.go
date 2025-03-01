package app

import (
	"context"
	"strings"

	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
)

const (
	defaultSourceFields = "id, name, description,type"
)

type SourceService struct {
	app *App
	_go.UnimplementedSourcesServer
	objClassName string
}

func (s SourceService) CreateSource(ctx context.Context, req *_go.CreateSourceRequest) (*_go.Source, error) {
	// Validate required fields
	if req.Name == "" {
		return nil, cerror.NewBadRequestError("source_service.create_source.name.required", ErrLookupNameReq)
	}

	// Validate the Type field
	if req.Type == _go.SourceType_TYPE_UNSPECIFIED {
		return nil, cerror.NewBadRequestError("source_service.create_source.type.required", "Source type is required")
	}

	fields := []string{"id", "name", "description", "type", "created_at", "updated_at", "created_by", "updated_by"}

	// Define create options
	createOpts := model.CreateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Fields:  fields,
	}

	// Define the current user as the creator and updater
	currentU := &_go.Lookup{
		Id: createOpts.GetAuthOpts().GetUserId(),
	}

	// Create a new source user_auth
	source := &_go.Source{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		CreatedBy:   currentU,
		UpdatedBy:   currentU,
	}

	// Create the source in the store
	l, e := s.app.Store.Source().Create(&createOpts, source)
	if e != nil {
		return nil, cerror.NewInternalError("source_service.create_source.store.create.failed", e.Error())
	}

	return l, nil
}

func (s SourceService) ListSources(ctx context.Context, req *_go.ListSourceRequest) (*_go.SourceList, error) {
	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultSourceFields, ", ")
	}

	// Use default page size and page number if not provided
	page := req.Page
	if page == 0 {
		page = 1
	}

	searchOptions := model.SearchOptions{
		IDs: req.Id,
		//UserAuthSession: session,
		Fields:  fields,
		Context: ctx,
		Page:    int(page),
		Sort:    req.Sort,
		Size:    int(req.Size),
		Filter:  make(map[string]interface{}),
		Auth:    model.GetAutherOutOfContext(ctx),
	}

	if req.Q != "" {
		searchOptions.Filter["name"] = req.Q
	}

	if len(req.Type) > 0 {
		searchOptions.Filter["type"] = req.Type
	}

	lookups, e := s.app.Store.Source().List(&searchOptions)
	if e != nil {
		return nil, cerror.NewInternalError("source_service.list_sources.store.list.failed", e.Error())
	}

	// // Publish an event to RabbitMQ
	// event := map[string]interface{}{
	// 	"action": "ListSources",
	// 	"user":   session.GetUserId(),
	// 	"query":  req.Q,
	// 	"type":   req.Type,
	// 	"fields": fields,
	// 	"page":   page,
	// 	"size":   req.Size,
	// }

	// eventData, err := json.Marshal(event)
	// if err != nil {
	// 	return nil, cerror.NewInternalError("source_service.list_sources.event_marshal.failed", err.Error())
	// }

	// err = s.app.rabbit.Publish(
	// 	user_auth.APP_SERVICE_NAME,
	// 	"list_sources_key",
	// 	eventData,
	// 	strconv.Itoa(int(session.GetUserId())),
	// 	time.Now(),
	// )
	// if err != nil {
	// 	return nil, cerror.NewInternalError("source_service.list_sources.event_publish.failed", err.Error())
	// }

	return lookups, nil
}

func (s SourceService) UpdateSource(ctx context.Context, req *_go.UpdateSourceRequest) (*_go.Source, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.update_source.id.required", "Source ID is required")
	}

	// Fields to update
	fields := []string{"id", "updated_at", "updated_by"}

	// Validate fields and add them to the update list
	for _, f := range req.XJsonMask {
		switch f {
		case "name":
			// Validate that name is not empty
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("source_service.update_source.name.required", "Name is required and cannot be empty")
			}
			fields = append(fields, "name")

		case "description":
			fields = append(fields, "description")

		case "type":
			// Validate that type is not unspecified
			if req.Input.Type == _go.SourceType_TYPE_UNSPECIFIED {
				return nil, cerror.NewBadRequestError("source_service.update_source.type.required", "Type is required and cannot be unspecified")
			}
			fields = append(fields, "type")
		}
	}

	// Define update options
	updateOpts := model.UpdateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Fields:  fields,
	}

	// Define the current user as the updater
	currentU := &_go.Lookup{
		Id: updateOpts.GetAuthOpts().GetUserId(),
	}

	// Update source user_auth
	source := &_go.Source{
		Id:          req.Id,
		Name:        req.Input.Name,
		Description: req.Input.Description,
		Type:        req.Input.Type,
		UpdatedBy:   currentU,
	}

	// Update the source in the store
	l, e := s.app.Store.Source().Update(&updateOpts, source)
	if e != nil {
		return nil, cerror.NewInternalError("source_service.update_source.store.update.failed", e.Error())
	}

	return l, nil
}

func (s SourceService) DeleteSource(ctx context.Context, req *_go.DeleteSourceRequest) (*_go.Source, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.delete_source.id.required", "Lookup ID is required")
	}

	// Define delete options
	deleteOpts := model.DeleteOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		IDs:     []int64{req.Id},
	}

	// Delete the source in the store
	e := s.app.Store.Source().Delete(&deleteOpts)
	if e != nil {
		return nil, cerror.NewInternalError("source_service.delete_source.store.delete.failed", e.Error())
	}

	return &(_go.Source{Id: req.Id}), nil
}

func (s SourceService) LocateSource(ctx context.Context, req *_go.LocateSourceRequest) (*_go.LocateSourceResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("source_service.locate_source.id.required", "Lookup ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListSourceRequest{
		Id:     []int64{req.Id},
		Fields: req.Fields,
		Page:   1,
		Size:   1, // We only need one item
	}

	// Call the ListSources method
	listResp, err := s.ListSources(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.locate_source.list_sources.error", err.Error())
	}

	// Check if the source was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("source_service.locate_source.not_found", "Source not found")
	}

	// Return the found source
	return &_go.LocateSourceResponse{Source: listResp.Items[0]}, nil
}

func NewSourceService(app *App) (*SourceService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_source_service.args_check.app_nil", "internal is nil")
	}
	return &SourceService{app: app, objClassName: model.ScopeDictionary}, nil
}
