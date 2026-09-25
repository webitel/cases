package postgres

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webitel/cases/internal/model"
)

func TestDecodeCaseTimelineVariables(t *testing.T) {
	tests := []struct {
		name    string
		raw     []byte
		exclude map[string]bool
		want    []*model.CaseTimelineVariable
	}{
		{
			name:    "sorted keys, excluded key dropped",
			raw:     []byte(`{"b":"2","a":"1","message_id":"skip-me"}`),
			exclude: map[string]bool{"message_id": true},
			want: []*model.CaseTimelineVariable{
				{Key: "a", Value: "1"},
				{Key: "b", Value: "2"},
			},
		},
		{
			name: "non-string value falls back to raw json",
			raw:  []byte(`{"count":3}`),
			want: []*model.CaseTimelineVariable{
				{Key: "count", Value: "3"},
			},
		},
		{
			name: "empty payload",
			raw:  []byte(`{}`),
			want: nil,
		},
		{
			name: "no payload",
			raw:  nil,
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := decodeCaseTimelineVariables(tt.raw, tt.exclude)
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDecodeCaseTimelineVariables_ChatSystemKeysExcluded(t *testing.T) {
	raw := []byte(`{"cid":"2022","chat":"telegram","flow":"352","from":"noname","user":"368924478","broadcast":"true","externalChatID":"368924478"}`)

	got, err := decodeCaseTimelineVariables(raw, chatSystemVariableKeys)
	require.NoError(t, err)
	require.Equal(t, []*model.CaseTimelineVariable{
		{Key: "broadcast", Value: "true"},
	}, got)
}

func TestDecodeCaseTimelinePostprocessing(t *testing.T) {
	raw := []byte(`[{"agent":{"id":7,"name":"Agent Smith"},"form":{"result":"ok"},"reporting_at":1699999999000}]`)

	got, err := decodeCaseTimelinePostprocessing(raw)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.Equal(t, int64(1699999999000), got[0].ReportingAt)
	require.Equal(t, "Agent Smith", *got[0].Agent.Name)
	require.Equal(t, map[string]any{"result": "ok"}, got[0].Form)
}

func TestDecodeCaseTimelinePostprocessing_Empty(t *testing.T) {
	got, err := decodeCaseTimelinePostprocessing(nil)
	require.NoError(t, err)
	require.Nil(t, got)

	got, err = decodeCaseTimelinePostprocessing([]byte(`null`))
	require.NoError(t, err)
	require.Nil(t, got)
}
