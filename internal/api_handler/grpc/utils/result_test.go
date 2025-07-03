package utils

import (
	"errors"
	"testing"
)

// mockLister implements Lister for testing.
type mockLister struct {
	size int
}

func (m mockLister) GetSize() int {
	return m.size
}

func TestGetListResult(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		items    []int
		wantNext bool
		wantList []int
	}{
		{
			name:     "next page exists",
			size:     2,
			items:    []int{1, 2, 3},
			wantNext: true,
			wantList: []int{1, 2},
		},
		{
			name:     "no next page",
			size:     3,
			items:    []int{1, 2, 3},
			wantNext: false,
			wantList: []int{1, 2, 3},
		},
		{
			name:     "empty items",
			size:     2,
			items:    []int{},
			wantNext: false,
			wantList: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			next, got := GetListResult(mockLister{tt.size}, tt.items)
			if next != tt.wantNext {
				t.Errorf("GetListResult() next = %v, want %v", next, tt.wantNext)
			}
			if len(got) != len(tt.wantList) {
				t.Errorf("GetListResult() got = %v, want %v", got, tt.wantList)
			}
			for i := range got {
				if got[i] != tt.wantList[i] {
					t.Errorf("GetListResult() got[%d] = %v, want %v", i, got[i], tt.wantList[i])
				}
			}
		})
	}
}

func TestConvertToOutputBulk(t *testing.T) {
	convert := func(i int) (string, error) {
		if i < 0 {
			return "", errors.New("negative value")
		}
		return string(rune('A' + i)), nil
	}

	t.Run("all ok", func(t *testing.T) {
		in := []int{0, 1, 2}
		want := []string{"A", "B", "C"}
		got, err := ConvertToOutputBulk(in, convert)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(got) != len(want) {
			t.Fatalf("got len %d, want %d", len(got), len(want))
		}
		for i := range got {
			if got[i] != want[i] {
				t.Errorf("got[%d]=%v, want %v", i, got[i], want[i])
			}
		}
	})

	t.Run("conversion error", func(t *testing.T) {
		in := []int{0, -1, 2}
		_, err := ConvertToOutputBulk(in, convert)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestResolvePaging(t *testing.T) {
	tests := []struct {
		name     string
		size     int
		items    []int
		wantList []int
		wantNext bool
	}{
		{
			name:     "paging with next",
			size:     2,
			items:    []int{1, 2, 3},
			wantList: []int{1, 2},
			wantNext: true,
		},
		{
			name:     "paging without next",
			size:     5,
			items:    []int{1, 2, 3},
			wantList: []int{1, 2, 3},
			wantNext: false,
		},
		{
			name:     "size zero returns all",
			size:     0,
			items:    []int{1, 2, 3},
			wantList: []int{1, 2, 3},
			wantNext: false,
		},
		{
			name:     "empty items",
			size:     2,
			items:    []int{},
			wantList: []int{},
			wantNext: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, next := ResolvePaging(tt.size, tt.items)
			if next != tt.wantNext {
				t.Errorf("ResolvePaging() next = %v, want %v", next, tt.wantNext)
			}
			if len(got) != len(tt.wantList) {
				t.Errorf("ResolvePaging() got = %v, want %v", got, tt.wantList)
			}
			for i := range got {
				if got[i] != tt.wantList[i] {
					t.Errorf("ResolvePaging() got[%d] = %v, want %v", i, got[i], tt.wantList[i])
				}
			}
		})
	}
}
