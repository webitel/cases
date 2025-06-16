package app

import (
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

const (
	ErrLookupNameReq    = "Lookup name is required"
	statusDefaultFields = "id, name, description, created_by"
)

// CreateStatus implements api.StatusesServer.
func (s *App) CreateStatus(opts options.Creator, req *model.Status) (*model.Status, error) {
	res, err := s.Store.Status().Create(opts, req)
	if err != nil {
		return nil, cerror.NewInternalError("status.create_status.store.create.failed", err.Error())
	}

	return res, nil
}

// ListStatuses implements api.StatusesServer.
func (s *App) ListStatus(opts options.Searcher) ([]*model.Status, error) {
	res, err := s.Store.Status().List(opts)
	if err != nil {
		return nil, cerror.NewInternalError("status.list_status.store.list.failed", err.Error())
	}
	return res, nil
}

// UpdateStatus implements api.StatusesServer.
func (s *App) UpdateStatus(opts options.Updator, input *model.Status) (*model.Status, error) {
	// Update the input in the store
	res, err := s.Store.Status().Update(opts, input)
	if err != nil {
		return nil, cerror.NewInternalError("status.update_status.store.update.failed", err.Error())
	}

	return res, nil
}

// DeleteStatus implements api.StatusesServer.
func (s *App) DeleteStatus(opts options.Deleter) (*model.Status, error) {
	// TODO: return deleted status
	_, err := s.Store.Status().Delete(opts)
	if err != nil {
		return nil, cerror.NewInternalError("status.delete_status.store.delete.failed", err.Error())
	}

	return nil, nil
}

// LocateStatus implements api.StatusesServer.
func (s *App) LocateStatus(opts options.Searcher) (*model.Status, error) {
	res, err := s.ListStatus(opts)
	if err != nil {
		return nil, cerror.NewInternalError("status.locate_status.list_status.error", err.Error())
	}

	// Check if the lookup was found
	if len(res) == 0 {
		return nil, cerror.NewNotFoundError("status.locate_status.not_found", "Status lookup not found")
	}

	// Return the found status lookup
	return res[0], nil
}
