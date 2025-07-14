package model

import (
	"encoding/json"
)

// Type of timeline event
type CaseTimelineEventType string

const (
	TimelineEventTypeChat  CaseTimelineEventType = "chat"
	TimelineEventTypeCall  CaseTimelineEventType = "call"
	TimelineEventTypeEmail CaseTimelineEventType = "email"
)

// CaseTimeline represents the timeline data structure
type CaseTimeline struct {
	Days []*DayTimeline
	Page int32
	Next bool
}

// DayTimeline represents a group of events for a specific day
type DayTimeline struct {
	Items        []*TimelineEvent `db:"-"`
	DayTimestamp int64            `db:"day_timestamp"`
	ChatsCount   int64            `db:"chats_count"`
	CallsCount   int64            `db:"calls_count"`
	EmailsCount  int64            `db:"emails_count"`
	ItemsJSON    json.RawMessage  `db:"items"`
}

// UnmarshalItems for processing the JSONB
func (d *DayTimeline) UnmarshalItems() error {
	if len(d.ItemsJSON) == 0 || string(d.ItemsJSON) == "null" || string(d.ItemsJSON) == "[]" {
		d.Items = []*TimelineEvent{}
		return nil
	}

	return json.Unmarshal(d.ItemsJSON, &d.Items)
}

// TimelineEvent represents a single event in the timeline
type TimelineEvent struct {
	Type      CaseTimelineEventType `json:"type"`
	CreatedAt int64                 `json:"created_at"`
	EventData json.RawMessage       `json:"event_data"`
	Event     any                   `json:"-"` // One of ChatEvent, CallEvent, EmailEvent - populated after scan
}

// UnmarshalEventData populates the Event field based on the Type and EventData
func (e *TimelineEvent) UnmarshalEventData() error {
	if len(e.EventData) == 0 {
		return nil
	}

	switch e.Type {
	case TimelineEventTypeChat:
		var chatEvent ChatEvent
		if err := json.Unmarshal(e.EventData, &chatEvent); err != nil {
			return err
		}
		e.Event = &chatEvent
	case TimelineEventTypeCall:
		var callEvent CallEvent
		if err := json.Unmarshal(e.EventData, &callEvent); err != nil {
			return err
		}
		e.Event = &callEvent
	case TimelineEventTypeEmail:
		var emailEvent EmailEvent
		if err := json.Unmarshal(e.EventData, &emailEvent); err != nil {
			return err
		}
		e.Event = &emailEvent
	}
	return nil
}

type ChatEvent struct {
	Id           string                 `json:"id"`
	ClosedAt     int64                  `json:"closed_at"`
	Duration     int64                  `json:"duration"`
	IsInbound    bool                   `json:"is_inbound"`
	IsMissed     bool                   `json:"is_missed"`
	Participants []*GeneralLookup       `json:"participants"`
	Gateway      *GeneralExtendedLookup `json:"gateway"`
	FlowScheme   *GeneralLookup         `json:"flow_scheme"`
	Queue        *GeneralLookup         `json:"queue"`
	IsDetailed   bool                   `json:"is_detailed"`
}

type CallEvent struct {
	Id            string              `json:"id"`
	ClosedAt      int64               `json:"closed_at"`
	Duration      int64               `json:"duration"`
	IsInbound     bool                `json:"is_inbound"`
	IsMissed      bool                `json:"is_missed"`
	Participants  []*GeneralLookup    `json:"participants"`
	Gateway       *GeneralLookup      `json:"gateway"`
	FlowScheme    *GeneralLookup      `json:"flow_scheme"`
	Queue         *GeneralLookup      `json:"queue"`
	IsDetailed    bool                `json:"is_detailed"`
	Files         []*CallFile         `json:"files"`
	Transcripts   []*TranscriptLookup `json:"transcripts"`
	TotalDuration int64               `json:"total_duration"`
}

type EmailEvent struct {
	Id          string         `json:"id"`
	From        []string       `json:"from"`
	To          []string       `json:"to"`
	Sender      []string       `json:"sender"`
	Cc          []string       `json:"cc"`
	IsInbound   bool           `json:"is_inbound"`
	Subject     string         `json:"subject"`
	Body        string         `json:"body"`
	Html        string         `json:"html"`
	Owner       *GeneralLookup `json:"owner"`
	Attachments []*Attachment  `json:"attachments"`
	Profile     *GeneralLookup `json:"profile"`
	IsDetailed  bool           `json:"is_detailed"`
}

// CallFile represents a file associated with a call
type CallFile struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	MimeType string `json:"mime_type"`
	StartAt  int64  `json:"start_at"`
	StopAt   int64  `json:"stop_at"`
}

// Attachment represents a file attached to an email
type Attachment struct {
	Id   int64  `json:"id"`
	Url  string `json:"url"`
	Mime string `json:"mime"`
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// TranscriptLookup represents a transcript associated with a call
type TranscriptLookup struct {
	Id     int64          `json:"id"`
	Locale string         `json:"locale"`
	File   *GeneralLookup `json:"file"`
}

// TimelineCounter represents a counter for a specific event type
type TimelineCounter struct {
	EventType string `db:"event_type"`
	Count     int64  `db:"count"`
	DateFrom  int64  `db:"date_from"`
	DateTo    int64  `db:"date_to"`
}

// TimelineCounterResponse represents the response for timeline counter
type TimelineCounterResponse struct {
	DateFrom    int64
	DateTo      int64
	ChatsCount  int64
	CallsCount  int64
	EmailsCount int64
}
