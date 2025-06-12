package grpc_api

import (
	"context"
	"time"

	"github.com/webitel/cases/api/cases"
	cerror "github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
)

type CloseReasonAPI struct {
	closeReasonService store.CloseReasonStore
	cases.UnimplementedCloseReasonsServer
}

func NewCloseReasonAPI(service store.CloseReasonStore) *CloseReasonAPI {
	return &CloseReasonAPI{closeReasonService: service}
}

var CloseReasonMetadata = model.NewObjectMetadata(model.ScopeDictionary, "", []*model.Field{
	{Name: "id", Default: true},
	{Name: "created_by", Default: true},
	{Name: "created_at", Default: true},
	{Name: "updated_by", Default: false},
	{Name: "updated_at", Default: false},
	{Name: "name", Default: true},
	{Name: "description", Default: true},
	{Name: "close_reason_id", Default: false},
})

func (api *CloseReasonAPI) LocateCloseReason(ctx context.Context, req *cases.LocateCloseReasonRequest) (*cases.LocateCloseReasonResponse, error) {
    if req.GetId() == 0 {
        return nil, cerror.NewBadRequestError("close_reason_service.locate_close_reason.id.required", "Close reason ID is required")
    }

    listReq := &cases.ListCloseReasonRequest{
        Id:                 []int64{req.GetId()},
        Fields:             req.GetFields(),
        Page:               1,
        Size:               1,
        CloseReasonGroupId: req.GetCloseReasonGroupId(),
    }
    searcher, err := grpcopts.NewSearchOptions(
        ctx,
        grpcopts.WithSearch(listReq),
        grpcopts.WithFields(listReq, CloseReasonMetadata,
            util.DeduplicateFields,
            util.EnsureIdField,
        ),
    )
    if err != nil {
        return nil, cerror.NewBadRequestError("close_reason_service.locate.options.invalid", err.Error())
    }

    result, err := api.closeReasonService.List(searcher, req.GetCloseReasonGroupId())
    if err != nil {
        return nil, err
    }
    if len(result.Items) == 0 {
        return nil, cerror.NewNotFoundError("close_reason_service.locate_close_reason.not_found", "Close reason not found")
    }

    return &cases.LocateCloseReasonResponse{CloseReason: ModelToProtoCloseReason(result.Items[0])}, nil
}

func (api *CloseReasonAPI) CreateCloseReason(ctx context.Context, req *cases.CreateCloseReasonRequest) (*cases.CloseReason, error) {
	creator, err := grpcopts.NewCreateOptions(
		ctx,
		grpcopts.WithCreateFields(req, CloseReasonMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
	)
	if err != nil {
		return nil, cerror.NewBadRequestError("close_reason_service.create.options.invalid", err.Error())
	}

	input := &model.CloseReason{
		Name:               req.GetInput().GetName(),
		Description:        strPtr(req.GetInput().GetDescription()),
		CloseReasonGroupId: req.GetCloseReasonGroupId(),
	}

	created, err := api.closeReasonService.Create(creator, input)
	if err != nil {
		return nil, err
	}
	return ModelToProtoCloseReason(created), nil
}

func (api *CloseReasonAPI) UpdateCloseReason(ctx context.Context, req *cases.UpdateCloseReasonRequest) (*cases.CloseReason, error) {
	if req.GetId() == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.update_close_reason.id.required", "Close reason ID is required")
	}

	updator, err := grpcopts.NewUpdateOptions(
		ctx,
		grpcopts.WithUpdateFields(req, CloseReasonMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithUpdateMasker(req),
	)
	if err != nil {
		return nil, cerror.NewBadRequestError("close_reason_service.update.options.invalid", err.Error())
	}

	input := &model.CloseReason{
		Id:                 req.GetId(),
		Name:               req.GetInput().GetName(),
		Description:        strPtr(req.GetInput().GetDescription()),
		CloseReasonGroupId: req.GetCloseReasonGroupId(),
	}

	updated, err := api.closeReasonService.Update(updator, input)
	if err != nil {
		return nil, err
	}
	return ModelToProtoCloseReason(updated), nil
}

func strPtr(s string) *string {
	return &s
}

func (api *CloseReasonAPI) ListCloseReasons(ctx context.Context, req *cases.ListCloseReasonRequest) (*cases.CloseReasonList, error) {
	searcher, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CloseReasonMetadata,
			util.DeduplicateFields,
			util.EnsureIdField,
		),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, cerror.NewBadRequestError("close_reason_service.list.options.invalid", err.Error())
	}
	searcher.AddFilter("name", req.GetQ())

	result, err := api.closeReasonService.List(searcher, req.GetCloseReasonGroupId())
	if err != nil {
		return nil, err
	}

	items := make([]*cases.CloseReason, len(result.Items))
	for i, m := range result.Items {
		items[i] = ModelToProtoCloseReason(m)
	}
	return &cases.CloseReasonList{
		Page:  int32(result.Page),
		Next:  result.Next,
		Items: items,
	}, nil
}

func (api *CloseReasonAPI) DeleteCloseReason(ctx context.Context, req *cases.DeleteCloseReasonRequest) (*cases.CloseReason, error) {
	if req.GetId() == 0 {
		return nil, cerror.NewBadRequestError("close_reason_service.delete_close_reason.id.required", "Close reason ID is required")
	}

	deleteOpts, err := grpcopts.NewDeleteOptions(
		ctx,
		grpcopts.WithDeleteID(req.GetId()),
	)
	if err != nil {
		return nil, cerror.NewBadRequestError("close_reason_service.delete.options.invalid", err.Error())
	}

	deleted, err := api.closeReasonService.Delete(deleteOpts)
	if err != nil {
		return nil, err
	}
	return ModelToProtoCloseReason(deleted), nil
}

func ModelToProtoCloseReason(m *model.CloseReason) *cases.CloseReason {
	if m == nil {
		return nil
	}
	return &cases.CloseReason{
		Id:                 m.Id,
		Name:               m.Name,
		Description:        dereferString(m.Description),
		CloseReasonGroupId: m.CloseReasonGroupId,
		CreatedAt:          m.CreatedAt.Unix(),
		UpdatedAt:          m.UpdatedAt.Unix(),
		CreatedBy:          ModelToProtoAuthor(m.Author),
		UpdatedBy:          ModelToProtoEditor(m.Editor),
	}
}

func ProtoToModelCloseReason(p *cases.CloseReason) *model.CloseReason {
	if p == nil {
		return nil
	}
	return &model.CloseReason{
		Id:                 p.Id,
		Name:               p.Name,
		Description:        &p.Description,
		CloseReasonGroupId: p.CloseReasonGroupId,
		CreatedAt:          time.Unix(p.CreatedAt, 0),
		UpdatedAt:          time.Unix(p.UpdatedAt, 0),
		Author:             ProtoToModelAuthor(p.CreatedBy),
		Editor:             ProtoToModelEditor(p.UpdatedBy),
	}
}

func ModelToProtoAuthor(a *model.Author) *cases.Lookup {
	if a == nil || a.Id == nil {
		return nil
	}
	return &cases.Lookup{
		Id:   *a.Id,
		Name: dereferString(a.Name),
	}
}

func ModelToProtoEditor(e *model.Editor) *cases.Lookup {
	if e == nil || e.Id == nil {
		return nil
	}
	return &cases.Lookup{
		Id:   *e.Id,
		Name: dereferString(e.Name),
	}
}

func dereferString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ProtoToModelAuthor(p *cases.Lookup) *model.Author {
	if p == nil {
		return nil
	}
	return &model.Author{
		Id:   &p.Id,
		Name: &p.Name,
	}
}

func ProtoToModelEditor(p *cases.Lookup) *model.Editor {
	if p == nil {
		return nil
	}
	return &model.Editor{
		Id:   &p.Id,
		Name: &p.Name,
	}
}
