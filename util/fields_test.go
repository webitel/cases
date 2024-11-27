package util

import (
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
			// res, _, _, _ := ProcessEtag(tt.fields)
			// if !areSlicesEqualUnordered(tt.want, res) {
			// 	t.Fail()
			// 	t.Log(fmt.Sprintf("want %v, have %v", tt.want, res))
			// }
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
