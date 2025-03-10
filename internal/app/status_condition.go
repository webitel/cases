package app

import (
	"context"
	"strings"
	"time"

	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/model"
)

type StatusConditionService struct {
	app *App
	_go.UnimplementedStatusConditionsServer
	objClassName string
}

const (
	ErrStatusNameReq    = "Status name is required"
	defaultFieldsStatus = "id, name, description, initial, final"
)

// CreateStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) CreateStatusCondition(ctx context.Context, req *_go.CreateStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Input.Name == "" {
		return nil, cerror.NewBadRequestError("status_condition.create_status_condition.name.required", ErrStatusNameReq)
	}

	fields := []string{"id", "lookup_id", "name", "description", "initial", "final", "created_at", "updated_at", "created_by", "updated_by"}

	t := time.Now()

	// Define create options
	createOpts := model.CreateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Create a new status user_auth
	status := &_go.StatusCondition{
		Name:        req.Input.Name,
		Description: req.Input.Description,
		StatusId:    req.StatusId,
	}

	// Create the status in the store
	st, e := s.app.Store.StatusCondition().Create(&createOpts, status)
	if e != nil {
		return nil, cerror.NewInternalError("status_condition.create_status_condition.store.create.failed", e.Error())
	}

	return st, nil
}

// ListStatusConditions implements api.StatusConditionsServer.
func (s StatusConditionService) ListStatusConditions(ctx context.Context, req *_go.ListStatusConditionRequest) (*_go.StatusConditionList, error) {
	fields := req.Fields
	if len(fields) == 0 {
		fields = strings.Split(defaultFieldsStatus, ", ")
	}

	// Use default page size and page number if not provided
	page := req.Page
	if page == 0 {
		page = 1
	}

	t := time.Now()
	searchOptions := model.SearchOptions{
		IDs: req.Id,
		// UserAuthSession: session,
		Fields:  fields,
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

	statuses, e := s.app.Store.StatusCondition().List(&searchOptions, req.StatusId)
	if e != nil {
		return nil, cerror.NewInternalError("status_condition.list_status_conditions.store.list.failed", e.Error())
	}

	return statuses, nil
}

// UpdateStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) UpdateStatusCondition(ctx context.Context, req *_go.UpdateStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status_condition.update_status_condition.id.required", "Status ID is required")
	}

	// Update status user_auth
	status := &_go.StatusCondition{
		Id:          req.Id,
		StatusId:    req.StatusId,
		Name:        req.Input.Name,
		Description: req.Input.Description,
	}

	fields := []string{"id", "lookup_id"}

	for _, f := range req.XJsonMask {
		if f == "name" {
			fields = append(fields, "name")
			if req.Input.Name == "" {
				return nil, cerror.NewBadRequestError("status_condition.update_status_condition.name.required", "Status name is required and cannot be empty")
			}
		}
		if f == "description" {
			fields = append(fields, "description")
		}
		if f == "initial" {
			fields = append(fields, "initial")
			status.Initial = req.Input.Initial.Value
		}
		if f == "final" {
			fields = append(fields, "final")
			status.Final = req.Input.Final.Value
		}
	}

	t := time.Now()

	// Define update options
	updateOpts := model.UpdateOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		Fields:  fields,
		Time:    t,
	}

	// Define the current user as the updater
	u := &_go.Lookup{
		Id: updateOpts.GetAuthOpts().GetUserId(),
	}
	status.UpdatedBy = u

	// Update the status in the store
	st, err := s.app.Store.StatusCondition().Update(&updateOpts, status)
	if err != nil {
		switch err.(type) {
		case *cerror.DBCheckViolationError:
			return nil, cerror.NewBadRequestError(
				"app.status_condition.update.initial_false_not_allowed",
				"update not allowed: there must be at least one initial = TRUE for the given dc and status_id",
			)
		case *cerror.DBInternalError:
			return nil, cerror.NewBadRequestError(
				"app.status_condition.update.error",
				err.Error(),
			)
		}
		return nil, cerror.NewInternalError(
			"app.status_condition.update.error",
			err.Error(),
		)
	}

	return st, nil
}

// DeleteStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) DeleteStatusCondition(ctx context.Context, req *_go.DeleteStatusConditionRequest) (*_go.StatusCondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status_condition.delete_status_condition.id.required", "Status ID is required")
	}

	t := time.Now()
	// Define delete options
	deleteOpts := model.DeleteOptions{
		Auth:    model.GetAutherOutOfContext(ctx),
		Context: ctx,
		ID:      req.Id,
		Time:    t,
	}

	// Delete the status in the store
	err := s.app.Store.StatusCondition().Delete(&deleteOpts, req.StatusId)
	if err != nil {
		switch err.(type) {
		case *cerror.DBNoRowsError:
			return nil, cerror.NewBadRequestError(
				"status_condition.delete_status_condition.not_found",
				"delete not allowed",
			)
		}
		return nil, cerror.NewInternalError(
			"status_condition.delete_status_condition.error",
			err.Error(),
		)
	}

	return &(_go.StatusCondition{Id: req.Id}), nil
}

// LocateStatusCondition implements api.StatusConditionsServer.
func (s StatusConditionService) LocateStatusCondition(ctx context.Context, req *_go.LocateStatusConditionRequest) (*_go.LocateStatusConditionResponse, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, cerror.NewBadRequestError("status_condition.locate_status_condition.id.required", "Status ID is required")
	}

	// Prepare a list request with necessary parameters
	listReq := &_go.ListStatusConditionRequest{
		Id:       []int64{req.Id},
		Fields:   req.Fields,
		Page:     1,
		Size:     1, // We only need one item
		StatusId: req.StatusId,
	}

	// Call the ListStatusConditions method
	listResp, err := s.ListStatusConditions(ctx, listReq)
	if err != nil {
		return nil, cerror.NewInternalError("status_condition.locate_status_condition.list_status_condition.error", err.Error())
	}

	// Check if the status condition was found
	if len(listResp.Items) == 0 {
		return nil, cerror.NewNotFoundError("status_condition.locate_status_condition.not_found", "Status condition not found")
	}

	// Return the found status condition
	return &_go.LocateStatusConditionResponse{Status: listResp.Items[0]}, nil
}

func NewStatusConditionService(app *App) (*StatusConditionService, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_status_condition_service.args_check.app_nil", "internal is nil")
	}
	return &StatusConditionService{app: app, objClassName: model.ScopeDictionary}, nil
}
