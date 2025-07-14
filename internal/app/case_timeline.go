package app

import (
	"strconv"
	"time"

	"github.com/webitel/cases/auth"
	"github.com/webitel/cases/internal/api_handler/grpc"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/internal/store"
)

// GetTimeline retrieves the timeline for a case
func (s *App) GetTimeline(searcher options.Searcher) (*model.CaseTimeline, error) {
	filters := searcher.GetFilter("case_id")
	if len(filters) == 0 {
		return nil, errors.InvalidArgument("case id required")
	}

	caseID, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil {
		return nil, errors.InvalidArgument("invalid case id", errors.WithCause(err))
	}

	accessMode := auth.Read
	if searcher.GetAuthOpts().IsRbacCheckRequired(grpc.CaseTimelineMetadata.GetParentScopeName(), accessMode) {
		access, err := s.Store.Case().CheckRbacAccess(searcher, searcher.GetAuthOpts(), accessMode, caseID)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case", errors.WithCause(err))
		}
	}

	timeline, err := s.Store.CaseTimeline().Get(searcher)
	if err != nil {
		return nil, err
	}

	return timeline, nil
}

// GetTimelineCounter retrieves the counter for a case timeline
func (s *App) GetTimelineCounter(searcher options.Searcher) (*model.TimelineCounterResponse, error) {
	filters := searcher.GetFilter("case_id")
	if len(filters) == 0 {
		return nil, errors.InvalidArgument("case id required")
	}

	caseID, err := strconv.ParseInt(filters[0].Value, 10, 64)
	if err != nil {
		return nil, errors.InvalidArgument("invalid case id", errors.WithCause(err))
	}

	accessMode := auth.Read
	if searcher.GetAuthOpts().IsRbacCheckRequired(grpc.CaseTimelineMetadata.GetParentScopeName(), accessMode) {
		access, err := s.Store.Case().CheckRbacAccess(searcher, searcher.GetAuthOpts(), accessMode, caseID)
		if err != nil {
			return nil, err
		}
		if !access {
			return nil, errors.Forbidden("user doesn't have required (READ) access to the case", errors.WithCause(err))
		}
	}

	eventTypeCounters, err := s.Store.CaseTimeline().GetCounter(searcher)
	if err != nil {
		return nil, err
	}

	if len(eventTypeCounters) == 0 {
		return &model.TimelineCounterResponse{}, nil
	}

	var (
		dateFrom = time.Now().UnixMilli()
		dateTo   = int64(0)
		response = &model.TimelineCounterResponse{}
	)

	for _, eventTypeCounter := range eventTypeCounters {
		// find max and min date
		switch eventTypeCounter.EventType {
		case store.CommunicationChat:
			response.ChatsCount = eventTypeCounter.Count
		case store.CommunicationCall:
			response.CallsCount = eventTypeCounter.Count
		case store.CommunicationEmail:
			response.EmailsCount = eventTypeCounter.Count
		}

		if eventTypeCounter.DateFrom < dateFrom && eventTypeCounter.DateFrom != 0 {
			dateFrom = eventTypeCounter.DateFrom
		}

		if eventTypeCounter.DateTo > dateTo {
			dateTo = eventTypeCounter.DateTo
		}
	}

	response.DateFrom = dateFrom
	response.DateTo = dateTo

	return response, nil
}
