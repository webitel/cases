package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

// CreateSLA implements cases.SLAsServer.
func (s *App) CreateSLA(
	creator options.Creator,
	input *model.SLA,
) (*model.SLA, error) {
	if input.Name == nil || *input.Name == "" {
		return nil, errors.InvalidArgument("SLA name is required")
	}
	if input.Calendar == nil || input.Calendar.GetId() == nil || *input.Calendar.GetId() == 0 {
		return nil, errors.InvalidArgument("Calendar ID is required")
	}
	if input.ReactionTime == 0 {
		return nil, errors.InvalidArgument("Reaction time is required")
	}
	if input.ResolutionTime == 0 {
		return nil, errors.InvalidArgument("Resolution time is required")
	}

	res, err := s.Store.SLA().Create(creator, input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteSLA implements cases.SLAsServer.
func (s *App) DeleteSLA(
	deleter options.Deleter,
) (*model.SLA, error) {
	if len(deleter.GetIDs()) == 0 {
		return nil, errors.InvalidArgument("SLA ID is required")
	}
	item, err := s.Store.SLA().Delete(deleter)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// ListSLAs implements cases.SLAsServer.
func (s *App) ListSLAs(
	searcher options.Searcher,
) ([]*model.SLA, error) {
	res, err := s.Store.SLA().List(searcher)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// LocateSLA implements cases.SLAsServer.
func (s *App) LocateSLA(
	searcher options.Searcher,
) (*model.SLA, error) {
	list, err := s.Store.SLA().List(searcher)
	if err != nil {
		return nil, err
	}

	// Check if the SLA was found
	if len(list) == 0 {
		return nil, errors.NotFound("SLA not found")
	}

	// Return the found SLA
	return list[0], nil
}

// UpdateSLA implements cases.SLAsServer.
func (s *App) UpdateSLA(
	updator options.Updator,
	input *model.SLA,
) (*model.SLA, error) {
	res, err := s.Store.SLA().Update(updator, input)
	if err != nil {
		return nil, err
	}
	return res, nil
}
