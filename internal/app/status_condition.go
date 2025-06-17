package app

import (
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

// CreateStatusCondition implements api.StatusConditionsServer.
func (s *App) CreateStatusCondition(rpc options.Creator, req *model.StatusCondition) (*model.StatusCondition, error) {
	// Create the status in the store
	st, e := s.Store.StatusCondition().Create(rpc, req)
	if e != nil {
		return nil, cerror.NewInternalError("status_condition.create_status_condition.store.create.failed", e.Error())
	}

	return st, nil
}

// ListStatusConditions implements api.StatusConditionsServer.
func (s *App) ListStatusConditions(opts options.Searcher) ([]*model.StatusCondition, error) {
	statuses, e := s.Store.StatusCondition().List(opts)
	if e != nil {
		return nil, cerror.NewInternalError("status_condition.list_status_conditions.store.list.failed", e.Error())
	}

	return statuses, nil
}

// UpdateStatusCondition implements api.StatusConditionsServer.
func (s *App) UpdateStatusCondition(opts options.Updator, input *model.StatusCondition) (*model.StatusCondition, error) {
	// Update the input in the store
	st, err := s.Store.StatusCondition().Update(opts, input)
	if err != nil {
		return nil, err
	}

	return st, nil
}

// DeleteStatusCondition implements api.StatusConditionsServer.
func (s *App) DeleteStatusCondition(opts options.Deleter) (*model.StatusCondition, error) {
	// Delete the status in the store
	_, err := s.Store.StatusCondition().Delete(opts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// LocateStatusCondition implements api.StatusConditionsServer.
func (s *App) LocateStatusCondition(searcher options.Searcher) (*model.StatusCondition, error) {

	// Call the ListStatusConditions method
	listResp, err := s.ListStatusConditions(searcher)
	if err != nil {
		return nil, cerror.NewInternalError("status_condition.locate_status_condition.list_status_condition.error", err.Error())
	}

	// Check if the status condition was found
	if len(listResp) == 0 {
		return nil, cerror.NewNotFoundError("status_condition.locate_status_condition.not_found", "Status condition not found")
	}

	// Return the found status condition
	return listResp[0], nil
}
