package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model/options"
	"google.golang.org/grpc/codes"

	"github.com/webitel/cases/internal/model"
)

// CreateSource implements api.SourcesServer.
func (s *App) CreateSource(
	opts options.Creator,
	input *model.Source,
) (*model.Source, error) {
	// Validate required fields
	if input.Name == nil || *input.Name == "" {
		return nil, errors.New("source name is required", errors.WithCode(codes.InvalidArgument))
	}
	if input.Type == nil || *input.Type == "" {
		return nil, errors.New("source type is required", errors.WithCode(codes.InvalidArgument))
	}
	// Create the source in the store
	res, err := s.Store.Source().Create(opts, input)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// ListSources implements api.SourcesServer.
func (s *App) ListSources(
	opts options.Searcher,
) ([]*model.Source, error) {
	res, err := s.Store.Source().List(opts)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// UpdateSource implements api.SourcesServer.
func (s *App) UpdateSource(
	opts options.Updator,
	req *model.Source,
) (*model.Source, error) {
	// Validate required fields
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("source ID is required", errors.WithCode(codes.InvalidArgument))
	}
	// Update the source in the store
	res, err := s.Store.Source().Update(opts, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// DeleteSource implements api.SourcesServer.
func (s *App) DeleteSource(
	opts options.Deleter,
) (*model.Source, error) {
	// Validate required fields
	if len(opts.GetIDs()) == 0 {
		return nil, errors.New("source ID is required", errors.WithCode(codes.InvalidArgument))
	}
	// Delete the source in the store
	_, err := s.Store.Source().Delete(opts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
