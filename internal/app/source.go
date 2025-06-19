package app

import (
	"github.com/webitel/cases/internal/model/options"

	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
)

// CreateSource implements api.SourcesServer.
func (s *App) CreateSource(
	opts options.Creator,
	input *model.Source,
) (*model.Source, error) {
	// Create the source in the store
	res, err := s.Store.Source().Create(opts, input)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.create_source.store.create.failed", err.Error())
	}

	return res, nil
}

// ListSources implements api.SourcesServer.
func (s *App) ListSources(
	opts options.Searcher,
) ([]*model.Source, error) {
	res, err := s.Store.Source().List(opts)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.list_sources.store.list.failed", err.Error())
	}

	return res, nil
}

// UpdateSource implements api.SourcesServer.
func (s *App) UpdateSource(
	opts options.Updator,
	req *model.Source,
) (*model.Source, error) {
	// Update the source in the store
	res, err := s.Store.Source().Update(opts, req)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.update_source.store.update.failed", err.Error())
	}

	return res, nil
}

// DeleteSource implements api.SourcesServer.
func (s *App) DeleteSource(
	opts options.Deleter,
) (*model.Source, error) {
	// Delete the source in the store
	_, err := s.Store.Source().Delete(opts)
	if err != nil {
		return nil, cerror.NewInternalError("source_service.delete_source.store.delete.failed", err.Error())
	}

	return nil, nil
}
