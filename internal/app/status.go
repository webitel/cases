package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"
)

// CreateStatus implements api.StatusesServer.
func (s *App) CreateStatus(opts options.Creator, req *model.Status) (*model.Status, error) {
	if req.Name == nil || *req.Name == "" {
		return nil, errors.New("lookup name is required", errors.WithCode(codes.InvalidArgument))
	}
	res, err := s.Store.Status().Create(opts, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ListStatuses implements api.StatusesServer.
func (s *App) ListStatus(opts options.Searcher) ([]*model.Status, error) {
	res, err := s.Store.Status().List(opts)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// UpdateStatus implements api.StatusesServer.
func (s *App) UpdateStatus(opts options.Updator, input *model.Status) (*model.Status, error) {
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("status id is required", errors.WithCode(codes.InvalidArgument))
	}
	// Update the input in the store
	res, err := s.Store.Status().Update(opts, input)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteStatus implements api.StatusesServer.
func (s *App) DeleteStatus(opts options.Deleter) (*model.Status, error) {
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("status id is required", errors.WithCode(codes.InvalidArgument))
	}
	// TODO: return deleted status
	_, err := s.Store.Status().Delete(opts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
