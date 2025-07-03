package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
)

// CreateStatusCondition implements api.StatusConditionsServer.
func (s *App) CreateStatusCondition(opts options.Creator, req *model.StatusCondition) (*model.StatusCondition, error) {
	// Validate required fields
	if req.Name == nil || *req.Name == "" {
		return nil, errors.New("status name is required", errors.WithCode(codes.InvalidArgument))
	}
	// Create the status in the store
	st, err := s.Store.StatusCondition().Create(opts, req)
	if err != nil {
		return nil, err
	}

	return st, nil
}

// ListStatusConditions implements api.StatusConditionsServer.
func (s *App) ListStatusConditions(opts options.Searcher) ([]*model.StatusCondition, error) {
	statuses, err := s.Store.StatusCondition().List(opts)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

// UpdateStatusCondition implements api.StatusConditionsServer.
func (s *App) UpdateStatusCondition(opts options.Updator, input *model.StatusCondition) (*model.StatusCondition, error) {
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("status condition id is required", errors.WithCode(codes.InvalidArgument))
	}
	// Update the input in the store
	st, err := s.Store.StatusCondition().Update(opts, input)
	if err != nil {
		return nil, err
	}

	return st, nil
}

// DeleteStatusCondition implements api.StatusConditionsServer.
func (s *App) DeleteStatusCondition(opts options.Deleter) (*model.StatusCondition, error) {
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("id for delete required", errors.WithCode(codes.InvalidArgument))
	}
	// Delete the status in the store
	_, err := s.Store.StatusCondition().Delete(opts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// LocateStatusCondition implements api.StatusConditionsServer.
func (s *App) LocateStatusCondition(opts options.Searcher) (*model.StatusCondition, error) {
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("id for locate required")
	}
	// Call the ListStatusConditions method
	listResp, err := s.ListStatusConditions(opts)
	if err != nil {
		return nil, err
	}

	// Check if the status condition was found
	if len(listResp) == 0 {
		return nil, errors.New("status condition not found", errors.WithCode(codes.InvalidArgument))
	}

	// Return the found status condition
	return listResp[0], nil
}
