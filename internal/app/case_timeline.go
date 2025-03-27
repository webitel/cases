package app

import (
	"context"
	"github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/errors"
	deferr "github.com/webitel/cases/internal/errors/defaults"
	"github.com/webitel/cases/model"
	grpcopts "github.com/webitel/cases/model/options/grpc"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/etag"
	"log/slog"
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

func NewCaseTimelineService(app *App) (*CaseTimelineService, errors.AppError) {
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
		return nil, NewBadRequestError(err)
	}
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseId())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_timeline.get_timeline.check_args.invalid_etag", "Invalid case etag")
	}
	searchOpts.AddFilter("case_id", tid.GetOid())
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", tid.GetOid()),
	)
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseTimelineMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, tid.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	res, err := c.app.Store.CaseTimeline().Get(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
	}
	return res, nil

}

func (c CaseTimelineService) GetTimelineCounter(ctx context.Context, request *cases.GetTimelineCounterRequest) (*cases.GetTimelineCounterResponse, error) {
	tid, err := etag.EtagOrId(etag.EtagCase, request.GetCaseId())
	if err != nil {
		return nil, errors.NewBadRequestError("app.case_timeline.get_timeline_counter.check_args.invalid_etag", "Invalid case etag")
	}
	searchOpts := &grpcopts.SearchOptions{Context: ctx, Fields: CaseTimelineMetadata.GetDefaultFields(), IDs: []int64{tid.GetOid()}, Auth: model.GetAutherOutOfContext(ctx)}
	logAttributes := slog.Group(
		"context",
		slog.Int64("user_id", searchOpts.GetAuthOpts().GetUserId()),
		slog.Int64("domain_id", searchOpts.GetAuthOpts().GetDomainId()),
		slog.Int64("case_id", tid.GetOid()),
	)
	accessMode := auth.Read
	if searchOpts.GetAuthOpts().IsRbacCheckRequired(CaseTimelineMetadata.GetParentScopeName(), accessMode) {
		access, err := c.app.Store.Case().CheckRbacAccess(searchOpts, searchOpts.GetAuthOpts(), accessMode, tid.GetOid())
		if err != nil {
			slog.ErrorContext(ctx, err.Error(), logAttributes)
			return nil, deferr.ForbiddenError
		}
		if !access {
			slog.ErrorContext(ctx, "user doesn't have required (READ) access to the case", logAttributes)
			return nil, deferr.ForbiddenError
		}
	}
	eventTypeCounters, err := c.app.Store.CaseTimeline().GetCounter(searchOpts)
	if err != nil {
		slog.ErrorContext(ctx, err.Error(), logAttributes)
		return nil, deferr.DatabaseError
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
