package app

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"

	watcherkit "github.com/webitel/webitel-go-kit/pkg/watcher"

	cfg "github.com/webitel/cases/config"
	"github.com/webitel/cases/internal/model"
)

func TestCaseCommentFtsEvent(t *testing.T) {
	const (
		internal = model.CaseCommentTypeInternal
		external = model.CaseCommentTypeExternal
	)

	tests := []struct {
		name      string
		event     watcherkit.EventType
		typ, prev int32
		want      watcherkit.EventType
		wantOK    bool
	}{
		{name: "create external", event: watcherkit.EventTypeCreate, typ: external, want: watcherkit.EventTypeCreate, wantOK: true},
		{name: "create internal skipped", event: watcherkit.EventTypeCreate, typ: internal, wantOK: false},
		{name: "update external", event: watcherkit.EventTypeUpdate, typ: external, prev: external, want: watcherkit.EventTypeUpdate, wantOK: true},
		{name: "update external unknown prev", event: watcherkit.EventTypeUpdate, typ: external, want: watcherkit.EventTypeUpdate, wantOK: true},
		{name: "internal to external indexes", event: watcherkit.EventTypeUpdate, typ: external, prev: internal, want: watcherkit.EventTypeCreate, wantOK: true},
		{name: "external to internal removes", event: watcherkit.EventTypeUpdate, typ: internal, prev: external, want: watcherkit.EventTypeDelete, wantOK: true},
		{name: "update internal removes", event: watcherkit.EventTypeUpdate, typ: internal, prev: internal, want: watcherkit.EventTypeDelete, wantOK: true},
		{name: "delete external", event: watcherkit.EventTypeDelete, typ: external, want: watcherkit.EventTypeDelete, wantOK: true},
		{name: "delete internal", event: watcherkit.EventTypeDelete, typ: internal, want: watcherkit.EventTypeDelete, wantOK: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := caseCommentFtsEvent(tc.event, tc.typ, tc.prev)
			require.Equal(t, tc.wantOK, ok)

			if tc.wantOK {
				require.Equal(t, tc.want, got)
			}
		})
	}
}

type capturePublisher struct {
	keys []string
}

func (p *capturePublisher) Publish(_ context.Context, _, routingKey string, _ []byte, _ amqp091.Table) error {
	p.keys = append(p.keys, routingKey)

	return nil
}

func TestCaseCommentTriggerRoutingKey(t *testing.T) {
	tests := []struct {
		name  string
		typ   int32
		event watcherkit.EventType
		want  string
	}{
		{name: "external create", typ: model.CaseCommentTypeExternal, event: watcherkit.EventTypeCreate, want: "cases.case_comments.create.7"},
		{name: "internal create", typ: model.CaseCommentTypeInternal, event: watcherkit.EventTypeCreate, want: "cases.internal_comments.create.7"},
		{name: "changed to external", typ: model.CaseCommentTypeExternal, event: watcherkit.EventTypeUpdate, want: "cases.case_comments.update.7"},
		{name: "changed to internal", typ: model.CaseCommentTypeInternal, event: watcherkit.EventTypeUpdate, want: "cases.internal_comments.update.7"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pub := &capturePublisher{}
			obs, err := NewTriggerObserver(pub, &cfg.TriggerWatcherConfig{ExchangeName: "cases", TopicName: "*"}, formCaseCommentTriggerModel, slog.Default())
			require.NoError(t, err)

			now := time.Now()
			comment := &model.CaseComment{Id: 1, CaseId: 2, Type: tc.typ, CreatedAt: &now, UpdatedAt: &now}
			err = obs.Update(tc.event, map[string]any{"obj": comment, "domain_id": int64(7)})
			require.NoError(t, err)
			require.Equal(t, []string{tc.want}, pub.keys)
		})
	}
}

func TestCaseCommentScope(t *testing.T) {
	require.Equal(t, model.ScopeCaseComments, model.CaseCommentScope(model.CaseCommentTypeExternal))
	require.Equal(t, model.ScopeInternalComments, model.CaseCommentScope(model.CaseCommentTypeInternal))
	require.Equal(t, model.ScopeInternalComments, model.CaseCommentScope(0))
	require.Equal(t, model.CaseCommentTypeInternal, model.CaseCommentTypeOrDefault(0))
	require.Equal(t, model.CaseCommentTypeExternal, model.CaseCommentTypeOrDefault(model.CaseCommentTypeExternal))
}
