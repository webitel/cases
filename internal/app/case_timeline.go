package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	errors "github.com/webitel/cases/internal/error"
	"github.com/webitel/cases/model"
	"github.com/webitel/webitel-go-kit/etag"
	"log/slog"
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
	// TODO: RBAC check
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_timeline.get_timeline.check_args.invalid_etag", "Invalid case etag")
	}
	searchOpts := model.NewSearchOptions(ctx, request, CaseTimelineMetadata)
	searchOpts.IDs = []int64{tid.GetOid()}
	res, err := c.app.Store.CaseTimeline().Get(searchOpts)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("case_id", tid.GetOid()))
		return nil, AppDatabaseError
	}
	return res, nil

}

func (c CaseTimelineService) GetTimelineCounter(ctx context.Context, request *cases.GetTimelineCounterRequest) (*cases.GetTimelineCounterResponse, error) {
	// TODO: RBAC check
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseEtag())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_timeline.get_timeline_counter.check_args.invalid_etag", "Invalid case etag")
	}
	searchOpts := &model.SearchOptions{Context: ctx, Fields: CaseTimelineMetadata.GetDefaultFields(), ParentId: tid.GetOid()}
	eventTypeCounters, err := c.app.Store.CaseTimeline().GetCounter(searchOpts)
	if err != nil {
		slog.Warn(err.Error(), slog.Int64("case_id", tid.GetOid()))
		return nil, AppDatabaseError
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

		if eventTypeCounter.DateFrom < dateFrom {
			dateFrom = eventTypeCounter.DateFrom
		}

		if eventTypeCounter.DateTo > dateTo {
			dateTo = eventTypeCounter.DateTo
		}
	}
	res.DateFrom = dateFrom
	res.DateTo = dateTo
	return &res, nil
}
