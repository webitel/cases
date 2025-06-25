package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	grpcopts "github.com/webitel/cases/internal/model/options/grpc"
	"github.com/webitel/cases/internal/store"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"time"
)

var CaseTimelineMetadata = model.NewObjectMetadata("", caseObjScope, []*model.Field{
	{cases.CaseTimelineEventType_call.String(), true},
	{cases.CaseTimelineEventType_chat.String(), true},
	{cases.CaseTimelineEventType_email.String(), true},
})

type CaseTimelineService struct {
	app *App
	cases.UnimplementedCaseTimelineServer
}

func NewCaseTimelineService(app *App) (*CaseTimelineService, error) {
	return &CaseTimelineService{app: app}, nil
}

func (c CaseTimelineService) GetTimeline(ctx context.Context, request *cases.GetTimelineRequest) (*cases.GetTimelineResponse, error) {
	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(request),
		grpcopts.WithPagination(request),
		grpcopts.WithFields(request, CaseTimelineMetadata,
			util.DeduplicateFields,
			func(in []string) []string {
				var requestedType []string
				for _, eventType := range request.Type {
					requestedType = append(requestedType, eventType.String())
				}
				if len(requestedType) != 0 {
					in = requestedType
				}
				return in
			},
		),
		grpcopts.WithSort(request),
	)
	if err != nil {
		return nil, err
	}
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseId())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid case etag", errors.WithCause(err))
	}
	searchOpts.AddFilter("case_id", tid.GetOid())

	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseTimelineMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, tid.GetOid())
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case", errors.WithCause(err))
		}
	}
	res, err := c.app.Store.CaseTimeline().Get(searchOpts)
	if err != nil {
		return nil, err
	}
	return res, nil

}

func (c CaseTimelineService) GetTimelineCounter(ctx context.Context, request *cases.GetTimelineCounterRequest) (*cases.GetTimelineCounterResponse, error) {
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseId())
	if err != nil {
		return nil, errors.InvalidArgument("Invalid case etag", errors.WithCause(err))
	}
	//opts, err := grpcopts.NewSearchOptions(ctx, grpcopts.WithIDsAsEtags(etag.EtagCase, request.GetCaseId()))
	//if err != nil {
	//	return nil, err
	//}
	searchOpts := &grpcopts.SearchOptions{Context: ctx, Fields: CaseTimelineMetadata.GetDefaultFields(), IDs: []int64{tid.GetOid()} /*Auth: model.GetAutherOutOfContext(ctx)*/}
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseTimelineMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, tid.GetOid())
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case", errors.WithCause(err))
		}
	}
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
		switch eventTypeCounter.EventType {
		case store.CommunicationChat:
			res.ChatsCount = eventTypeCounter.Count
		case store.CommunicationCall:
			res.CallsCount = eventTypeCounter.Count
		case store.CommunicationEmail:
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
