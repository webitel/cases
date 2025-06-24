package app

import (
	"github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

type SLAService struct {
	app *App
	cases.UnimplementedSLAsServer
	objClassName string
}

// CreateSLA implements cases.SLAsServer.
func (s *App) CreateSLA(
	creator options.Creator,
	input *model.SLA,
) (*model.SLA, error) {
	if input.Name == nil || *input.Name == "" {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.name.required", "SLA name is required")
	}
	if input.Calendar == nil || input.Calendar.GetId() == nil || *input.Calendar.GetId() == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.calendar_id.required", "Calendar ID is required")
	}
	if input.ReactionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.reaction_time.required", "Reaction time is required")
	}
	if input.ResolutionTime == 0 {
		return nil, cerror.NewBadRequestError("sla_service.create_sla.resolution_time.required", "Resolution time is required")
	}

	res, err := s.Store.SLA().Create(creator, input)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.create_sla.store.create.failed", err.Error())
	}
	return res, nil
}

// DeleteSLA implements cases.SLAsServer.
func (s *App) DeleteSLA(
	deleter options.Deleter,
) (*model.SLA, error) {
	item, err := s.Store.SLA().Delete(deleter)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.delete_sla.store.delete.failed", err.Error())
	}
	return item, nil
}

// ListSLAs implements cases.SLAsServer.
func (s *App) ListSLAs(
	searcher options.Searcher,
) ([]*model.SLA, error) {
	res, err := s.Store.SLA().List(searcher)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.list_sla.store.list.failed", err.Error())
	}
	return res, nil
}

// LocateSLA implements cases.SLAsServer.
func (s *App) LocateSLA(
	searcher options.Searcher,
) (*model.SLA, error) {
	list, err := s.Store.SLA().List(searcher)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.locate_sla.list_slas.error", err.Error())
	}

	// Check if the SLA was found
	if len(list) == 0 {
		return nil, cerror.NewNotFoundError("sla_service.locate_sla.not_found", "SLA not found")
	}

	// Return the found SLA
	return list[0], nil
}

// UpdateSLA implements cases.SLAsServer.
func (s *App) UpdateSLA(
	updator options.Updator,
	input *model.SLA,
) (*model.SLA, error) { // Validate required fields
	// Update the SLA in the store
	res, err := s.Store.SLA().Update(updator, input)
	if err != nil {
		return nil, cerror.NewInternalError("sla_service.update_sla.store.update.failed", err.Error())
	}
	return res, nil
}