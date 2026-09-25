package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/webitel/cases/api/cases"
	grpcopts "github.com/webitel/cases/internal/api_handler/grpc/options"
	"github.com/webitel/cases/internal/api_handler/grpc/utils"
	"github.com/webitel/cases/internal/errors"
	"github.com/webitel/cases/internal/model"
	"github.com/webitel/cases/internal/model/options"
	"github.com/webitel/cases/util"
	"github.com/webitel/webitel-go-kit/pkg/etag"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

type CaseTimelineHandler interface {
	GetTimeline(options.Searcher) (*model.CaseTimeline, error)
	GetTimelineCounter(options.Searcher) (*model.TimelineCounterResponse, error)
	GetTimelineItemInfo(searcher options.Searcher, itemType model.CaseTimelineEventType, itemID string) (*model.CaseTimelineItemInfo, error)
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

// GetTimelineItemInfo handles the gRPC request to get variables + postprocessing
// results for a single timeline communication (call | chat | email) of a case.
func (s *CaseTimelineService) GetTimelineItemInfo(
	ctx context.Context,
	req *cases.GetTimelineItemInfoRequest,
) (*cases.GetTimelineItemInfoResponse, error) {
	if req.GetCaseId() == "" {
		return nil, errors.InvalidArgument("case id is required")
	}
	if req.GetId() == "" {
		return nil, errors.InvalidArgument("item id is required")
	}

	// Decode case etag to numeric ID
	caseTid, err := etag.EtagOrId(etag.EtagCase, req.GetCaseId())
	if err != nil {
		return nil, errors.InvalidArgument("invalid case etag", errors.WithCause(err))
	}

	itemType, err := unmarshalCaseTimelineEventType(req.GetType())
	if err != nil {
		return nil, err
	}

	searchOpts, err := grpcopts.NewSearchOptions(ctx)
	if err != nil {
		return nil, err
	}

	searchOpts.AddFilter(util.EqualFilter("case_id", caseTid.GetOid()))

	info, err := s.app.GetTimelineItemInfo(searchOpts, itemType, req.GetId())
	if err != nil {
		return nil, err
	}

	return s.MarshalTimelineItemInfo(info)
}

func unmarshalCaseTimelineEventType(t cases.CaseTimelineEventType) (model.CaseTimelineEventType, error) {
	switch t {
	case cases.CaseTimelineEventType_call:
		return model.TimelineEventTypeCall, nil
	case cases.CaseTimelineEventType_chat:
		return model.TimelineEventTypeChat, nil
	case cases.CaseTimelineEventType_email:
		return model.TimelineEventTypeEmail, nil
	default:
		return "", errors.InvalidArgument(fmt.Sprintf("unknown timeline event type %q", t))
	}
}

// MarshalTimelineItemInfo converts a model.CaseTimelineItemInfo to its gRPC representation.
func (s *CaseTimelineService) MarshalTimelineItemInfo(info *model.CaseTimelineItemInfo) (*cases.GetTimelineItemInfoResponse, error) {
	result := &cases.GetTimelineItemInfoResponse{}
	if info == nil {
		return result, nil
	}

	for _, v := range info.Variables {
		result.Variables = append(result.Variables, &cases.CaseTimelineVariable{
			Key:   v.Key,
			Value: v.Value,
		})
	}

	for _, p := range info.Postprocessing {
		item := &cases.CaseTimelinePostprocessingResult{
			ReportingAt: p.ReportingAt,
		}
		if p.Agent != nil {
			item.Agent = utils.MarshalLookup(p.Agent)
		}
		if p.Form != nil {
			formValue, err := structpb.NewValue(p.Form)
			if err != nil {
				return nil, errors.Internal("failed to marshal postprocessing form", errors.WithCause(err))
			}
			item.Form = formValue
		}
		result.Postprocessing = append(result.Postprocessing, item)
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
			Channel:  file.Channel,
			Type:     callFileType(file.Channel, file.MimeType),
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

func callFileType(channel, mime string) cases.CallFileType {
	const fileChannelScreenRecordings = "screenrecording" // engine/model.FileChannelScreenRecordings
	switch {
	case channel == fileChannelScreenRecordings:
		return cases.CallFileType_file_type_screensharing
	case strings.HasPrefix(mime, "audio/"):
		return cases.CallFileType_file_type_audio
	case strings.HasPrefix(mime, "video/"):
		return cases.CallFileType_file_type_video
	case strings.HasPrefix(mime, "image/"):
		return cases.CallFileType_file_type_screenshot
	case strings.HasPrefix(mime, "application/pdf"):
		return cases.CallFileType_file_type_pdf
	default:
		return cases.CallFileType_file_type_empty
	}
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
