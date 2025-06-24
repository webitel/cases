package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
)

func (s *App) CreateCloseReasonGroup(
	rpc options.Creator,
	input *model.CloseReasonGroup,
) (*model.CloseReasonGroup, error) {
	// Validate required fields
	if input.Name == nil || *input.Name == "" {
		return nil, errors.New("lookup name is required", errors.WithCode(codes.InvalidArgument))
	}

	// Create the close reason group in the store
	res, err := s.Store.CloseReasonGroup().Create(rpc, input)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *App) ListCloseReasonGroup(
	rpc options.Searcher,
) ([]*model.CloseReasonGroup, error) {
	res, err := s.Store.CloseReasonGroup().List(rpc)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *App) UpdateCloseReasonGroup(
	rpc options.Updator,
	input *model.CloseReasonGroup,
) (*model.CloseReasonGroup, error) {
	// Validate required fields
	if len(rpc.GetIDs()) == 0 {
		return nil, errors.New("lookup ID is required", errors.WithCode(codes.InvalidArgument))
	}

	// Update the lookup in the store
	res, err := s.Store.CloseReasonGroup().Update(rpc, input)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *App) DeleteCloseReasonGroup(
	rpc options.Deleter,
) (*model.CloseReasonGroup, error) {
	// Validate required fields
	if len(rpc.GetIDs()) == 0 {
		return nil, errors.New("lookup ID is required", errors.WithCode(codes.InvalidArgument))
	}
	// Delete the lookup in the store
	err := s.Store.CloseReasonGroup().Delete(rpc)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
