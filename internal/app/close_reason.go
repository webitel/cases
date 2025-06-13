package app

import (
	_go "github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
)

type CloseReasonService struct {
	app *App
	_go.UnimplementedCloseReasonsServer
	objClassName string
}
type CloseReasonService1 struct {
	store store.CloseReasonStore
}

var CloseReasonMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{"id", true},
	{"created_by", true},
	{"created_at", true},
	{"updated_by", false},
	{"updated_at", false},
	{"name", true},
	{"description", true},
	{"close_reason_id", false},
})

// CreateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService1) Create(creator options.Creator, input *model.CloseReason) (*model.CloseReason, error) {
	if creator.GetAuthOpts().GetDomainId() == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.create.domain_id.required", "Domain ID is required")
	}
	return s.store.Create(creator, input)
}

// ListCloseReasons implements api.CloseReasonsServer.
func (s *CloseReasonService1) List(searcher options.Searcher, closeReasonGroupId int64) (*model.CloseReasonList, error) {
	if searcher.GetAuthOpts().GetDomainId() == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.list.domain_id.required", "Domain ID is required")
	}
	return s.store.List(searcher, closeReasonGroupId)
}

// UpdateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService1) Update(updator options.Updator, input *model.CloseReason) (*model.CloseReason, error) {
	if updator.GetAuthOpts().GetDomainId() == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.update.domain_id.required", "Domain ID is required")
	}
	return s.store.Update(updator, input)
}

// DeleteCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService1) Delete(deleter options.Deleter) (*model.CloseReason, error) {
	if deleter.GetAuthOpts().GetDomainId() == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.delete.domain_id.required", "Domain ID is required")
	}
	return s.store.Delete(deleter)
}

// LocateCloseReason implements api.CloseReasonsServer.
func (s *CloseReasonService1) Locate(searcher options.Searcher, closeReasonGroupId int64) (*model.CloseReason, error) {
	if searcher.GetAuthOpts().GetDomainId() == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.locate.domain_id.required", "Domain ID is required")
	}

	list, err := s.store.List(searcher, closeReasonGroupId)
	if err != nil {
		return nil, err
	}
	if len(list.Items) == 0 {
		return nil, cerror.NewNotFoundError("close_reason_service.locate_close_reason.not_found", "Close reason not found")
	}
	return list.Items[0], nil
}

func NewCloseReasonService(app *App) (*CloseReasonService1, cerror.AppError) {
	if app == nil {
		return nil, cerror.NewInternalError("api.config.new_close_reason_service.args_check.app_nil", "internal is nil")
	}
	return &CloseReasonService1{store: app.Store.CloseReason()}, nil
}
