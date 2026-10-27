package app

import (
	"context"
	"testing"
)

func TestArchiveGuardSkipped(t *testing.T) {
	ctx := context.Background()
	if archiveGuardSkipped(ctx) {
		t.Fatal("plain context must not skip the archive guard")
	}

	if !archiveGuardSkipped(withoutArchiveGuard(ctx)) {
		t.Fatal("withoutArchiveGuard context must skip the archive guard")
	}
}

func TestArchivedCaseUpdateAllowed(t *testing.T) {
	tests := []struct {
		name      string
		mask      []string
		canManage bool
		want      error
	}{
		{
			name:      "status change without permission",
			mask:      []string{"etag", "status_condition"},
			canManage: false,
			want:      ErrCaseArchived,
		},
		{
			name:      "field change without permission",
			mask:      []string{"etag", "subject"},
			canManage: false,
			want:      ErrCaseArchived,
		},
		{
			name:      "status change with permission",
			mask:      []string{"etag", "status_condition"},
			canManage: true,
		},
		{
			name:      "status change with id and ver and permission",
			mask:      []string{"id", "ver", "etag", "status_condition"},
			canManage: true,
		},
		{
			name:      "status change with close info and permission",
			mask:      []string{"etag", "userID", "status", "status_condition", "close_reason", "close_result", "close_article"},
			canManage: true,
		},
		{
			name:      "field change with permission",
			mask:      []string{"etag", "status_condition", "assignee"},
			canManage: true,
			want:      ErrArchivedCaseStatusOnly,
		},
		{
			name:      "custom fields with permission",
			mask:      []string{"custom"},
			canManage: true,
			want:      ErrArchivedCaseStatusOnly,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := archivedCaseUpdateAllowed(tc.mask, tc.canManage)
			if got != tc.want {
				t.Fatalf("archivedCaseUpdateAllowed(%v, %v) = %v, want %v", tc.mask, tc.canManage, got, tc.want)
			}
		})
	}
}
