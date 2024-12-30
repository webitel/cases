package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	errors "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/webitel-go-kit/etag"
	"time"
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
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, err
	}
	searchOpts := &model.SearchOptions{Context: ctx, Fields: CaseTimelineMetadata.GetDefaultFields(), ParentId: tid.GetOid()}
	eventTypeCounters, err := c.app.Store.CaseTimeline().GetCounter(searchOpts)
	if err != nil {
		return nil, err
	}
	if len(eventTypeCounters) == 0 {
		return nil, nil
	}
	var (
		dateFrom = time.Now().UnixMilli()
		dateTo   = int64(0)
		res      cases.GetTimelineCounterResponse
	)

	for _, eventTypeCounter := range eventTypeCounters {
		// find max and min date
		switch cases.CaseCommunicationsTypes(eventTypeCounter.EventType) {
		case cases.CaseCommunicationsTypes_COMMUNICATION_CHAT:
			res.ChatsCount = eventTypeCounter.Count
		case cases.CaseCommunicationsTypes_COMMUNICATION_CALL:
			res.CallsCount = eventTypeCounter.Count
		case cases.CaseCommunicationsTypes_COMMUNICATION_EMAIL:
			res.EmailsCount = eventTypeCounter.Count
		}

		if v := eventTypeCounter.DateFrom; v < dateFrom {
			dateFrom = v
		}

		if v := eventTypeCounter.DateTo; v > dateTo {
			dateTo = v
		}
	}
	res.DateFrom = dateFrom
	res.DateTo = dateTo
	return &res, nil
}
