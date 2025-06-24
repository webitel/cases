package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
)

// CreatePriority creates a new priority in the store.
func (s *App) CreatePriority(
	creator options.Creator,
	input *model.Priority,
) (*model.Priority, error) {
	if input.Name == "" {
		return nil, errors.New("priority name is required", errors.WithCode(codes.InvalidArgument))
	}
	res, err := s.Store.Priority().Create(creator, input)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return res, nil
}

// UpdatePriority updates a priority in the store.
func (s *App) UpdatePriority(
	updator options.Updator,
	input *model.Priority,
) (*model.Priority, error) {
	if len(updator.GetIDs()) == 0 {
		return nil, errors.New("priority ID is required", errors.WithCode(codes.InvalidArgument))
	}
	res, err := s.Store.Priority().Update(updator, input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeletePriority deletes a priority from the store.
func (s *App) DeletePriority(
	deleter options.Deleter,
) (*model.Priority, error) {
	if len(deleter.GetIDs()) == 0 {
		return nil, errors.New("priority ID is required", errors.WithCode(codes.InvalidArgument))
	}
	item, err := s.Store.Priority().Delete(deleter)
	if err != nil {
		return nil, err
	}
	return item, nil
}
