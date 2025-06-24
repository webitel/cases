package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
)

// CreateSLACondition implements cases.SLAConditionsServer.
func (s *App) CreateSLACondition(opts options.Creator, req *model.SLACondition) (*model.SLACondition, error) {
	// Validate required fields
	if req.Name == nil || *req.Name == "" {
		return nil, errors.New("SLA Condition name is required", errors.WithCode(codes.InvalidArgument))
	}
	if len(req.Priorities) == 0 {
		return nil, errors.New("at least one priority is required", errors.WithCode(codes.InvalidArgument))
	}
	if req.ReactionTime == nil || *req.ReactionTime == 0 {
		return nil, errors.New("reaction time is required", errors.WithCode(codes.InvalidArgument))
	}
	if req.ResolutionTime == nil || *req.ResolutionTime == 0 {
		return nil, errors.New("resolution time is required", errors.WithCode(codes.InvalidArgument))
	}
	if req.SlaId == nil || *req.SlaId == 0 {
		return nil, errors.New("SLA ID is required", errors.WithCode(codes.InvalidArgument))
	}

	// Create the SLACondition in the store
	r, err := s.Store.SLACondition().Create(opts, req)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// DeleteSLACondition implements cases.SLAConditionsServer.
func (s *App) DeleteSLACondition(opts options.Deleter) (*model.SLACondition, error) {
	// Validate required fields
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("SLA Condition ID is required", errors.WithCode(codes.InvalidArgument))
	}

	// Delete the SLACondition in the store
	_, err := s.Store.SLACondition().Delete(opts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ListSLAConditions implements cases.SLAConditionsServer.
func (s *App) ListSLAConditions(opts options.Searcher) ([]*model.SLACondition, error) {
	slaConditions, err := s.Store.SLACondition().List(opts)
	if err != nil {
		return nil, err
	}

	return slaConditions, nil
}

// UpdateSLACondition implements cases.SLAConditionsServer.
func (s *App) UpdateSLACondition(opts options.Updator, req *model.SLACondition) (*model.SLACondition, error) {
	// Validate required fields
	if req.Id == 0 {
		return nil, errors.New("SLA Condition ID is required", errors.WithCode(codes.InvalidArgument))
	}

	// Update the SLACondition in the store
	item, err := s.Store.SLACondition().Update(opts, req)
	if err != nil {
		return nil, err
	}

	return item, nil
}
