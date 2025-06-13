package app

import (
	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

func (s *App) CreateCloseReasonGroup(
	rpc options.Creator,
	input *_go.CloseReasonGroup,
) (*_go.CloseReasonGroup, error) {
	// Create the close reason group in the store
	res, err := s.Store.CloseReasonGroup().Create(rpc, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.create_close_reason_group.store.create.failed", err.Error())
	}

	return res, nil
}

func (s *App) ListCloseReasonGroups(
	rpc options.Searcher,
) (*_go.CloseReasonGroupList, error) {
	res, err := s.Store.CloseReasonGroup().List(rpc)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.list_close_reason_groups.store.list.failed", err.Error())
	}
	return res, nil
}

func (s *App) UpdateCloseReasonGroup(
	rpc options.Searcher,
	input *model.CloseReasonGroup,
) (*_go.CloseReasonGroup, error) {
	// Update the lookup in the store
	res, err := s.Store.CloseReasonGroup().Update(rpc, input)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.update_close_reason_group.store.update.failed", err.Error())
	}

	return res, nil
}

func (s *App) DeleteCloseReasonGroup(
	rpc options.Deleter,
) (*_go.CloseReasonGroup, error) {
	// Delete the lookup in the store
	err := s.Store.CloseReasonGroup().Delete(rpc)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.delete_close_reason_group.store.delete.failed", err.Error())
	}
	return nil, nil
}

func (s *App) LocateCloseReasonGroup(
	ctx options.Searcher,
) (*_go.LocateCloseReasonGroupResponse, error) {

	// Call the ListCloseReasonGroups method
	res, err := s.ListCloseReasonGroups(ctx)
	if err != nil {
		return nil, cerror.NewInternalError("close_reason_group_service.locate_close_reason_group.list_close_reason_groups.error", err.Error())
	}

	// Check if the close reason group was found
	if len(res.Items) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_group_service.locate_close_reason_group.not_found", "Close reason group not found")
	}

	// Return the found close reason group
	return &_go.LocateCloseReasonGroupResponse{CloseReasonGroup: res.Items[0]}, nil
}
