package app

import (
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

// CreatePriority creates a new priority in the store.
func (s *App) CreatePriority(
	creator options.Creator,
	input *model.Priority,
) (*model.Priority, error) {
	res, err := s.Store.Priority().Create(creator, input)
	if err != nil {
		return nil, cerror.NewInternalError("priority_service.create.store.create.failed", err.Error())
	}
	return res, nil
}

// ListPriorities lists priorities with optional SLA filters.
func (s *App) ListPriorities(
	searcher options.Searcher,
	notInSla int64,
	inSla int64,
) ([]*model.Priority, error) {
	res, err := s.Store.Priority().List(searcher, notInSla, inSla)
	if err != nil {
		return nil, cerror.NewInternalError("priority_service.list.store.list.failed", err.Error())
	}
	return res, nil
}

// UpdatePriority updates a priority in the store.
func (s *App) UpdatePriority(
	updator options.Updator,
	input *model.Priority,
) (*model.Priority, error) {
	res, err := s.Store.Priority().Update(updator, input)
	if err != nil {
		return nil, cerror.NewInternalError("priority_service.update.store.update.failed", err.Error())
	}
	return res, nil
}

// DeletePriority deletes a priority from the store.
func (s *App) DeletePriority(
	deleter options.Deleter,
) (*model.Priority, error) {
	item, err := s.Store.Priority().Delete(deleter)
	if err != nil {
		return nil, cerror.NewInternalError("priority_service.delete.store.delete.failed", err.Error())
	}
	return item, nil
}

// LocatePriority finds a priority by criteria.
func (s *App) LocatePriority(
	searcher options.Searcher,
) (*model.Priority, error) {
	list, err := s.Store.Priority().List(searcher, 0, 0)
	if err != nil {
		return nil, cerror.NewInternalError("priority_service.locate.store.list.failed", err.Error())
	}
	if len(list) == 0 {
		return nil, cerror.NewNotFoundError("priority_service.locate.not_found", "Priority not found")
	}
	return list[0], nil
}
