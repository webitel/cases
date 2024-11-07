package util

import (
	"strconv"
	"strings"
	"unicode"
)

// InlineFields explodes an inline 'attr,attr2 attr3' selector into ['attr','attr2','attr3'].
func InlineFields(selector string) []string {
	// split func to explode inline userattrs selector
	split := func(r rune) bool {
		return r == ',' || unicode.IsSpace(r)
	}
	selector = strings.ToLower(selector)
	return strings.FieldsFunc(selector, split)
}

// FieldsFunc normalizes a selection list src of the attributes to be returned.
//
//  1. An empty list with no attributes requests the return of all user attributes.
//  2. A list containing "*" (with zero or more attribute descriptions)
//     requests the return of all user attributes in addition to other listed (operational) attributes.
//
// e.g.: ['id,name','display'] returns ['id','name','display']
func FieldsFunc(src []string, fn func(string) []string) []string {
	if len(src) == 0 {
		return fn("")
	}

	var dst []string
	for i := 0; i < len(src); i++ {
		// explode single selection attr
		switch set := fn(src[i]); len(set) {
		case 0: // none
			src = append(src[:i], src[i+1:]...)
			i-- // process this i again
		case 1: // one
			if len(set[0]) == 0 {
				src = append(src[:i], src[i+1:]...)
				i--
			} else if dst == nil {
				src[i] = set[0]
			} else {
				dst = MergeFields(dst, set)
			}
		default: // many
			// NOTE: should rebuild output
			if dst == nil && i > 0 {
				// copy processed entries
				dst = make([]string, i, len(src)-1+len(set))
				copy(dst, src[:i])
			}
			dst = MergeFields(dst, set)
		}
	}
	if dst == nil {
		return src
	}
	return dst
}

// MergeFields appends a unique set from src to dst.
func MergeFields(dst, src []string) []string {
	if len(src) == 0 {
		return dst
	}
	//
	if cap(dst)-len(dst) < len(src) {
		ext := make([]string, len(dst), len(dst)+len(src))
		copy(ext, dst)
		dst = ext
	}

next: // append unique set of src to dst
	for _, attr := range src {
		if len(attr) == 0 {
			continue
		}
		// look backwards for duplicates
		for j := len(dst) - 1; j >= 0; j-- {
			if strings.EqualFold(dst[j], attr) {
				continue next // duplicate found
			}
		}
		// append unique attr
		dst = append(dst, attr)
	}
	return dst
}

func ContainsField(fields []string, field string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

func Int64SliceToStringSlice(ids []int64) []string {
	strIds := make([]string, len(ids))
	for i, id := range ids {
		strIds[i] = strconv.FormatInt(id, 10)
	}
	return strIds
}

// Helper function to check if a field exists in the update options
// ---------------------------------------------------------------------//
// ---- Example Usage ----
// if !util.FieldExists("name", rpc.Fields) {
func FieldExists(field string, fields []string) bool {
	for _, f := range fields {
		if f == field {
			return true
		}
	}
	return false
}

// EnsureIdAndVerFields ensures that "id" and "ver" are present in the rpc.Fields.
// Need it for etag encoding as ver + id is required.
func EnsureIdAndVerField(fields []string) []string {
	hasId := false
	hasVer := false

	// Check for "id" and "ver" in the fields
	for _, field := range fields {
		if field == "id" {
			hasId = true
		}
		if field == "ver" {
			hasVer = true
		}
	}

	// Add "id" if not found
	if !hasId {
		fields = append(fields, "id")
	}
	// Add "ver" if not found
	// Necessary for etag encoding as ver is required
	if !hasVer {
		fields = append(fields, "ver")
	}

	return fields
}
