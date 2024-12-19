package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	errors "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/webitel-go-kit/etag"
)

var CaseTimelineMetadata = model.NewObjectMetadata(
	[]*model.Field{
		{"calls", true},
		{"chats", true},
		{"emails", true},
	})

type CaseTimelineService struct {
	app *App
	cases.UnimplementedCaseTimelineServer
}

func NewCaseTimelineService(app *App) (*CaseTimelineService, errors.AppError) {
	return &CaseTimelineService{app: app}, nil
}

func (c CaseTimelineService) GetTimeline(ctx context.Context, request *cases.GetTimelineRequest) (*cases.GetTimelineResponse, error) {
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, err
	}
	searchOpts := model.NewSearchOptions(ctx, request, CaseTimelineMetadata)
	searchOpts.IDs = []int64{tid.GetOid()}
	res, err := c.app.Store.CaseTimeline().Get(searchOpts)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (c CaseTimelineService) GetTimelineCounter(ctx context.Context, request *cases.GetTimelineCounterRequest) (*cases.GetTimelineCounterResponse, error) {
	//TODO implement me
	panic("implement me")
}
