package util

import (
	"fmt"
	"testing"
)

func TestSearchOptions_ProcessEtag(t *testing.T) {

	tests := []struct {
		name   string
		fields []string
		want   []string
	}{
		{
			name:   "ETAG present",
			fields: []string{"etag"},
			want:   []string{"id", "ver"},
		},
		{
			name:   "ETAG not present",
			fields: []string{"id"},
			want:   []string{"id"},
		},
		{
			name:   "ID present",
			fields: []string{"etag", "id"},
			want:   []string{"id", "ver"},
		},
		{
			name:   "VER present",
			fields: []string{"etag", "ver"},
			want:   []string{"id", "ver"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := ParseFieldsForEtag(tt.fields)
			if !areSlicesEqualUnordered(tt.want, res) {
				t.Fail()
				t.Log(fmt.Sprintf("want %v, have %v", tt.want, tt.fields))
			}
		})
	}
}

func areSlicesEqualUnordered(slice1, slice2 []string) bool {
	// If the lengths are not equal, the slices can't be equal
	if len(slice1) != len(slice2) {
		return false
	}

	// Create maps to count occurrences of each element in the slices
	countMap1 := make(map[string]int)
	countMap2 := make(map[string]int)

	// Count occurrences in slice1
	for _, str := range slice1 {
		countMap1[str]++
	}

	// Count occurrences in slice2
	for _, str := range slice2 {
		countMap2[str]++
	}

	// Compare the two maps
	for key, count1 := range countMap1 {
		if count2, found := countMap2[key]; !found || count1 != count2 {
			return false
		}
	}

	// If all elements have the same counts, the slices are equal (unordered)
	return true
}

func TestSplitKnownAndUnknownFields(t *testing.T) {
	type args struct {
		requestedFields []string
		modelFields     []string
	}
	tests := []struct {
		name        string
		args        args
		wantKnown   []string
		wantUnknown []string
	}{
		{
			name: "1",
			args: args{
				requestedFields: []string{"foo", "bar", "id", "etag", "known", "super", "yehor"},
				modelFields:     []string{"id", "known", "super"},
			},
			wantKnown:   []string{"id", "known", "super"},
			wantUnknown: []string{"foo", "bar", "etag", "yehor"},
		},
		{
			name: "2",
			args: args{
				requestedFields: []string{"1", "2", "id", "3", "known", "super", "4"},
				modelFields:     []string{"id", "known", "super"},
			},
			wantKnown:   []string{"id", "known", "super"},
			wantUnknown: []string{"1", "2", "3", "4"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotKnown, gotUnknown := SplitKnownAndUnknownFields(tt.args.requestedFields, tt.args.modelFields)
			if !areSlicesEqualUnordered(gotKnown, tt.wantKnown) {
				t.Errorf("SplitKnownAndUnknownFields() gotKnown = %v, want %v", gotKnown, tt.wantKnown)
			}
			if !areSlicesEqualUnordered(gotUnknown, tt.wantUnknown) {
				t.Errorf("SplitKnownAndUnknownFields() gotUnknown = %v, want %v", gotUnknown, tt.wantUnknown)
			}
		})
	}
}
