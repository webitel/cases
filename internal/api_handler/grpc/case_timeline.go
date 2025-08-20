package grpc

import (
	"context"
	"fmt"
	"github.com/webitel/cases/api/cases"
	grpcopts "github.com/webitel/cases/internal/api_handler/grpc/options"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"google.golang.org/grpc/codes"
)

type CaseTimelineHandler interface {
	GetTimeline(options.Searcher) (*model.CaseTimeline, error)
	GetTimelineCounter(options.Searcher) (*model.TimelineCounterResponse, error)
}

type CaseTimelineService struct {
	app CaseTimelineHandler
	cases.UnimplementedCaseTimelineServer
}

func NewCaseTimelineService(app CaseTimelineHandler) (*CaseTimelineService, error) {
	if app == nil {
		return nil, errors.New("case timeline handler is nil")
	}
	return &CaseTimelineService{app: app}, nil
}

var CaseTimelineMetadata = model.NewObjectMetadata("", "cases", []*model.Field{
	{Name: string(model.TimelineEventTypeCall), Default: true},
	{Name: string(model.TimelineEventTypeChat), Default: true},
	{Name: string(model.TimelineEventTypeEmail), Default: true},
})

// GetTimeline handles the gRPC request to get the timeline for a case.
func (s *CaseTimelineService) GetTimeline(
	ctx context.Context,
	req *cases.GetTimelineRequest,
) (*cases.GetTimelineResponse, error) {
	if req.GetCaseId() == "" {
		return nil, errors.InvalidArgument("case id is required")
	}

	// Decode case etag to numeric ID
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseId())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}

	searchOpts, err := grpcopts.NewSearchOptions(
		ctx,
		grpcopts.WithSearch(req),
		grpcopts.WithPagination(req),
		grpcopts.WithFields(req, CaseTimelineMetadata,
			util.DeduplicateFields,
			func(in []string) []string {
				var requestedType []string
				for _, eventType := range req.Type {
					requestedType = append(requestedType, eventType.String())
				}
				if len(requestedType) != 0 {
					in = requestedType
				}
				return in
			},
		),
		grpcopts.WithSort(req),
	)
	if err != nil {
		return nil, err
	}

	searchOpts.AddFilter(fmt.Sprintf("case_id=%d", caseTid.GetOid()))

	// Call the business logic
	timeline, err := s.app.GetTimeline(searchOpts)
	if err != nil {
		return nil, err
	}

	// Convert from model to proto
	result, err := s.MarshalTimeline(timeline)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetTimelineCounter handles the gRPC request to get the counter for a case timeline.
func (s *CaseTimelineService) GetTimelineCounter(
	ctx context.Context,
	req *cases.GetTimelineCounterRequest,
) (*cases.GetTimelineCounterResponse, error) {
	if req.GetCaseId() == "" {
		return nil, errors.InvalidArgument("case id is required")
	}

	// Decode case etag to numeric ID
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseId())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}

	searchOpts, err := grpcopts.NewSearchOptions(ctx)
	if err != nil {
		return nil, err
	}

	searchOpts.AddFilter(util.EqualFilter("case_id", caseTid.GetOid()))

	// Call the business logic
	counter, err := s.app.GetTimelineCounter(searchOpts)
	if err != nil {
		return nil, err
	}

	if counter == nil {
		return nil, errors.New("no timeline data found", errors.WithCode(codes.NotFound))
	}

	// Convert from model to proto
	result := &cases.GetTimelineCounterResponse{
		DateFrom:    counter.DateFrom,
		DateTo:      counter.DateTo,
		ChatsCount:  counter.ChatsCount,
		CallsCount:  counter.CallsCount,
		EmailsCount: counter.EmailsCount,
	}

	return result, nil
}

// MarshalTimeline converts a model.CaseTimeline to its gRPC representation.
func (s *CaseTimelineService) MarshalTimeline(timeline *model.CaseTimeline) (*cases.GetTimelineResponse, error) {
	if timeline == nil {
		return nil, nil
	}

	result := &cases.GetTimelineResponse{
		Page: timeline.Page,
		Next: timeline.Next,
	}

	// Convert days
	for _, day := range timeline.Days {
		protoDay := &cases.DayTimeline{
			DayTimestamp: day.DayTimestamp,
			ChatsCount:   day.ChatsCount,
			CallsCount:   day.CallsCount,
			EmailsCount:  day.EmailsCount,
		}

		// Convert events
		for _, event := range day.Items {
			protoEvent := &cases.Event{
				CreatedAt: event.CreatedAt,
			}

			// Set event type and data based on the event type
			switch event.Type {
			case model.TimelineEventTypeChat:
				protoEvent.Type = cases.CaseTimelineEventType_chat
				if event.Event != nil {
					protoEvent.Event = &cases.Event_Chat{
						Chat: s.MarshalChatEvent(event.Event.(*model.ChatEvent)),
					}
				}
			case model.TimelineEventTypeCall:
				protoEvent.Type = cases.CaseTimelineEventType_call
				if event.Event != nil {
					protoEvent.Event = &cases.Event_Call{
						Call: s.MarshalCallEvent(event.Event.(*model.CallEvent)),
					}
				}
			case model.TimelineEventTypeEmail:
				protoEvent.Type = cases.CaseTimelineEventType_email
				if event.Event != nil {
					protoEvent.Event = &cases.Event_Email{
						Email: s.MarshalEmailEvent(event.Event.(*model.EmailEvent)),
					}
				}
			}

			protoDay.Items = append(protoDay.Items, protoEvent)
		}

		result.Days = append(result.Days, protoDay)
	}

	return result, nil
}

// MarshalChatEvent converts a model.ChatEvent to its gRPC representation.
func (s *CaseTimelineService) MarshalChatEvent(event *model.ChatEvent) *cases.ChatEvent {
	if event == nil {
		return nil
	}

	result := &cases.ChatEvent{
		Id:         event.Id,
		ClosedAt:   event.ClosedAt,
		Duration:   event.Duration,
		IsInbound:  event.IsInbound,
		IsMissed:   event.IsMissed,
		IsDetailed: event.IsDetailed,
	}

	for _, participant := range event.Participants {
		if participant != nil {
			result.Participants = append(result.Participants, utils.MarshalLookup(participant))
		}
	}

	if event.Gateway != nil {
		result.Gateway = utils.MarshalExtendedLookup(event.Gateway)
	}

	if event.FlowScheme != nil {
		result.FlowScheme = utils.MarshalLookup(event.FlowScheme)
	}

	if event.Queue != nil {
		result.Queue = utils.MarshalLookup(event.Queue)
	}

	return result
}

// MarshalCallEvent converts a model.CallEvent to its gRPC representation.
func (s *CaseTimelineService) MarshalCallEvent(event *model.CallEvent) *cases.CallEvent {
	if event == nil {
		return nil
	}

	result := &cases.CallEvent{
		Id:            event.Id,
		ClosedAt:      event.ClosedAt,
		Duration:      event.Duration,
		IsInbound:     event.IsInbound,
		IsMissed:      event.IsMissed,
		IsDetailed:    event.IsDetailed,
		TotalDuration: event.TotalDuration,
	}

	for _, participant := range event.Participants {
		if participant != nil {
			result.Participants = append(result.Participants, utils.MarshalLookup(participant))
		}
	}

	if event.Gateway != nil {
		result.Gateway = utils.MarshalLookup(event.Gateway)
	}

	if event.FlowScheme != nil {
		result.FlowScheme = utils.MarshalLookup(event.FlowScheme)
	}

	if event.Queue != nil {
		result.Queue = utils.MarshalLookup(event.Queue)
	}

	// Convert files
	for _, file := range event.Files {
		result.Files = append(result.Files, &cases.CallFile{
			Id:       file.Id,
			Name:     file.Name,
			Size:     file.Size,
			MimeType: file.MimeType,
			StartAt:  file.StartAt,
			StopAt:   file.StopAt,
		})
	}

	// Convert transcripts
	for _, transcript := range event.Transcripts {
		protoTranscript := &cases.TranscriptLookup{
			Id:     transcript.Id,
			Locale: transcript.Locale,
		}
		if transcript.File != nil {
			protoTranscript.File = utils.MarshalLookup(transcript.File)
		}
		result.Transcripts = append(result.Transcripts, protoTranscript)
	}

	return result
}

// MarshalEmailEvent converts a model.EmailEvent to its gRPC representation.
func (s *CaseTimelineService) MarshalEmailEvent(event *model.EmailEvent) *cases.EmailEvent {
	if event == nil {
		return nil
	}

	result := &cases.EmailEvent{
		Id:         event.Id,
		From:       event.From,
		To:         event.To,
		Sender:     event.Sender,
		Cc:         event.Cc,
		IsInbound:  event.IsInbound,
		Subject:    event.Subject,
		Body:       event.Body,
		Html:       event.Html,
		IsDetailed: event.IsDetailed,
	}

	if event.Owner != nil {
		result.Owner = utils.MarshalLookup(event.Owner)
	}

	if event.Profile != nil {
		result.Profile = utils.MarshalLookup(event.Profile)
	}

	for _, attachment := range event.Attachments {
		result.Attachments = append(result.Attachments, &cases.Attachment{
			Id:   attachment.Id,
			Url:  attachment.Url,
			Mime: attachment.Mime,
			Name: attachment.Name,
			Size: attachment.Size,
		})
	}

	return result
}
