package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/model"
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

func (c CaseTimelineService) GetTimeline(ctx context.Context, request *cases.GetTimelineRequest) (*cases.GetTimelineResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c CaseTimelineService) GetTimelineCounter(ctx context.Context, request *cases.GetTimelineCounterRequest) (*cases.GetTimelineCounterResponse, error) {
	//TODO implement me
	panic("implement me")
}
