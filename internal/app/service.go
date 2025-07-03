package app

import (
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
)

// CreateService implements cases.ServicesServer.
func (s *App) CreateService(opts options.Creator, item *model.Service) (*model.Service, error) {
	if item.Name == nil || *item.Name == "" {
		return nil, errors.InvalidArgument("Service name is required")
	}
	if item.RootId == nil || *item.RootId == 0 {
		return nil, errors.InvalidArgument("Root ID is required")
	}
	// Create the Service in the store
	r, e := s.Store.Service().Create(opts, item)
	if e != nil {
		return nil, e
	}
	return r, nil
}

// DeleteService implements cases.ServicesServer.
func (s *App) DeleteService(opts options.Deleter) (*model.Service, error) {
	if len(opts.GetIDs()) == 0 {
		return nil, errors.InvalidArgument("Service ID is required")
	}
	_, err := s.Store.Service().Delete(opts)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ListServices implements cases.ServicesServer.
func (s *App) ListServices(opts options.Searcher) ([]*model.Service, error) {
	services, e := s.Store.Service().List(opts)
	if e != nil {
		return nil, e
	}
	return services, nil
}

// UpdateService implements cases.ServicesServer.
func (s *App) UpdateService(opts options.Updator, req *model.Service) (*model.Service, error) {
	if req.Id == 0 {
		return nil, errors.InvalidArgument("Service ID is required")
	}
	r, e := s.Store.Service().Update(opts, req)
	if e != nil {
		return nil, e
	}
	return r, nil
}
