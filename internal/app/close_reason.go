package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
)

// CreateCloseReason creates a new close reason in the store.
func (s *App) CreateCloseReason(
	creator options.Creator,
	input *model.CloseReason,
) (*model.CloseReason, error) {
	if input.Name == "" {
		return nil, errors.New("close reason name is required", errors.WithCode(codes.InvalidArgument))
	}
	res, err := s.Store.CloseReason().Create(creator, input)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return res, nil
}

// UpdateCloseReason updates a close reason in the store.
func (s *App) UpdateCloseReason(
	updator options.Updator,
	input *model.CloseReason,
) (*model.CloseReason, error) {
	// Validate required fields
	if len(updator.GetIDs()) == 0 {
		return nil, errors.New("reason ID is required", errors.WithCode(codes.InvalidArgument))
	}
	res, err := s.Store.CloseReason().Update(updator, input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeleteCloseReason deletes a close reason from the store.
func (s *App) DeleteCloseReason(
	deleter options.Deleter,
) (*model.CloseReason, error) {
	if len(deleter.GetIDs()) == 0 {
		return nil, errors.New("reason ID is required", errors.WithCode(codes.InvalidArgument))
	}
	item, err := s.Store.CloseReason().Delete(deleter)
	if err != nil {
		return nil, err
	}
	return item, nil
}
