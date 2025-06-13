package app

import (
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

func (s *App) CreateCloseReasonGroup(
	rpc options.Creator,
	input *model.CloseReasonGroup,
) (*model.CloseReasonGroup, error) {
	// Create the close reason group in the store
	res, err := s.Store.CloseReasonGroup().Create(rpc, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.create_close_reason_group.store.create.failed", err.Error())
	}

	return res, nil
}

func (s *App) ListCloseReasonGroup(
	rpc options.Searcher,
) ([]*model.CloseReasonGroup, error) {
	res, err := s.Store.CloseReasonGroup().List(rpc)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.list_close_reason_groups.store.list.failed", err.Error())
	}
	return res, nil
}

func (s *App) UpdateCloseReasonGroup(
	rpc options.Updator,
	input *model.CloseReasonGroup,
) (*model.CloseReasonGroup, error) {
	// Update the lookup in the store
	res, err := s.Store.CloseReasonGroup().Update(rpc, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.update_close_reason_group.store.update.failed", err.Error())
	}

	return res, nil
}

func (s *App) DeleteCloseReasonGroup(
	rpc options.Deleter,
) (*model.CloseReasonGroup, error) {
	// Delete the lookup in the store
	err := s.Store.CloseReasonGroup().Delete(rpc)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.delete_close_reason_group.store.delete.failed", err.Error())
	}
	return nil, nil
}

func (s *App) LocateCloseReasonGroup(
	ctx options.Searcher,
) (*model.CloseReasonGroup, error) {

	// Call the ListCloseReasonGroups method
	res, err := s.ListCloseReasonGroup(ctx)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.locate_close_reason_group.list_close_reason_groups.error", err.Error())
	}

	// Check if the close reason group was found
	if len(res) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_group_service.locate_close_reason_group.not_found", "close reason group not found")
	}

	// Return the found close reason group
	return res[0], nil
}
