package util

import (
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/webitel/webitel-go-kit/etag"
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

func ProcessEtag(fields []string) (res []string, hasEtag bool, hasId bool, hasVer bool) {

	// Iterate through the fields and update the flags
	for _, field := range fields {
		if field == "etag" {
			hasEtag = true
			continue
		} else if field == "id" {
			hasId = true
		} else if field == "ver" {
			hasVer = true
		}
		res = append(res, field)
	}
	if hasEtag {
		if !hasId {
			res = append(res, "id")
		}
		if !hasVer {
			res = append(res, "ver")
		}
	}
	return
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

// EnsureFields ensures that all specified fields are present in the list of fields.
// If any field is missing, it will be added to the list.
func EnsureFields(fields []string, requiredFields ...string) []string {
	for _, requiredField := range requiredFields {
		if !ContainsField(fields, requiredField) {
			fields = append(fields, requiredField)
		}
	}
	return fields
}

// ParseQin converts a slice of strings (each possibly containing comma-separated eTags or numeric IDs)
// into a slice of int64. For example, given input ["1", "2,3", "etag4"], it converts each to int64 and returns []int64{1, 2, 3, 4}.
func ParseQin(input []string, etagType etag.EtagType) ([]int64, error) {
	var result []int64

	for _, item := range input {
		// Split the item by comma to handle comma-separated values
		parts := strings.Split(item, ",")

		for _, part := range parts {
			// Trim whitespace
			part = strings.TrimSpace(part)

			// Try to parse as int64
			num, err := strconv.ParseInt(part, 10, 64)
			if err == nil {
				// Successfully parsed as int64, add to result
				result = append(result, num)
				continue
			}

			// If parsing as int64 fails, try parsing as eTag
			tag, etagErr := etag.EtagOrId(etagType, part)
			if etagErr != nil {
				return nil, errors.New("invalid eTag or ID: " + part)
			}

			// Add the eTag converted ID to result
			result = append(result, tag.GetOid())
		}
	}

	return result, nil
}
