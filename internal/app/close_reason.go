package app

import (
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

// CreateCloseReason creates a new close reason in the store.
func (s *App) CreateCloseReason(
	creator options.Creator,
	input *model.CloseReason,
) (*model.CloseReason, error) {
	res, err := s.Store.CloseReason().Create(creator, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.create.store.create.failed", err.Error())
	}
	return res, nil
}

// ListCloseReasons lists close reasons for a group.
func (s *App) ListCloseReasons(
	searcher options.Searcher,
	closeReasonGroupId int64,
) ([]*model.CloseReason, error) {
	res, err := s.Store.CloseReason().List(searcher, closeReasonGroupId)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.list.store.list.failed", err.Error())
	}
	return res, nil
}

// UpdateCloseReason updates a close reason in the store.
func (s *App) UpdateCloseReason(
	updator options.Updator,
	input *model.CloseReason,
) (*model.CloseReason, error) {
	res, err := s.Store.CloseReason().Update(updator, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.update.store.update.failed", err.Error())
	}
	return res, nil 
}

// DeleteCloseReason deletes a close reason from the store.
func (s *App) DeleteCloseReason(
    deleter options.Deleter,
) (*model.CloseReason, error) {
    item, err := s.Store.CloseReason().Delete(deleter)
    if err != nil {
        return nil, cerror.NewInternalError("close_reason_service.delete.store.delete.failed", err.Error())
    }
    return item, nil
}

// LocateCloseReason finds a close reason by criteria.
func (s *App) LocateCloseReason(
	searcher options.Searcher,
	closeReasonGroupId int64,
) (*model.CloseReason, error) {
	list, err := s.Store.CloseReason().List(searcher, closeReasonGroupId)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_service.locate.store.list.failed", err.Error())
	}
	if len(list) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_service.locate.not_found", "Close reason not found")
	}
	return list[0], nil
}
