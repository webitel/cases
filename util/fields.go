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

func DeduplicateFields(in []string) []string {
	seen := make(map[string]struct{}) // Using struct{} to save memory
	var result []string

	for _, str := range in {
		if _, exists := seen[str]; !exists {
			seen[str] = struct{}{}
			result = append(result, str)
		}
	}
	return result
}

// AddVersionAndIdByEtag searches for etag, id, ver fields and determines what fields should be added
// to provide full functionality of etag
func AddVersionAndIdByEtag(fields []string) {
	var hasEtag, hasId, hasVer bool
	hasEtag, hasId, hasVer = FindEtagFields(fields)
	if hasEtag {
		if !hasId {
			fields = append(fields, "id")
		}
		if !hasVer {
			fields = append(fields, "ver")
		}
	}
}

func RemoveElements(arr []string, elementsToRemove ...string) []string {
	elementSet := make(map[string]bool)
	for _, elem := range elementsToRemove {
		elementSet[elem] = true
	}

	filteredArr := []string{}
	for _, item := range arr {
		if !elementSet[item] {
			filteredArr = append(filteredArr, item)
		}
	}

	return filteredArr
}

// ParseFieldsForEtag searches for id, ver fields and adds missing
// to provide full functionality of etag (do not change fields, returns fully new slice)
func ParseFieldsForEtag(fields []string) []string {
	var (
		res                    []string
		hasEtag, hasId, hasVer bool
	)
	for _, field := range fields {
		switch field {
		case "etag":
			hasEtag = true
		case "id":
			res = append(res, field)
			hasId = true
		case "ver":
			res = append(res, field)
			hasVer = true
		default:
			res = append(res, field)
		}
	}

	if hasEtag {
		if !hasId {
			res = append(res, "id")
		}
		if !hasVer {
			res = append(res, "ver")
		}
	}
	return res
}

func FindEtagFields(fields []string) (hasEtag bool, hasId bool, hasVer bool) {
	// Iterate through the fields and update the flags
	for _, field := range fields {
		switch field {
		case "etag":
			hasEtag = true
		case "id":
			hasId = true
		case "ver":
			hasVer = true
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

// EnsureIdField ensures that "id" is present in the rpc.Fields.
// Necessary when the "id" field is required for specific operations.
func EnsureIdField(fields []string) []string {
	hasId := false

	// Check for "id" in the fields
	for _, field := range fields {
		if field == "id" {
			hasId = true
			break
		}
	}

	// Add "id" if not found
	if !hasId {
		fields = append(fields, "id")
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

// NormalizeEtag normalizes etag, id, ver fields visibility for the response depending on what fields were requested.
func NormalizeEtag(etagType etag.EtagType, fields []string, etg *string, id *int64, ver *int32) error {
	var err error
	hasEtag, hasId, hasVer := FindEtagFields(fields)
	if hasEtag {
		*etg, err = etag.EncodeEtag(etagType, *id, *ver)
		if err != nil {
			return err
		}
		// hide
		if !hasId {
			*id = 0
		}
		if !hasVer {
			*ver = 0
		}
	}
	return nil
}

// NormalizeEtags normalizes etag, id, ver fields visibility for the response depending on what fields were requested.
// Function is an optimized-for-slice copy of NormalizeEtag function.
//
// Before usage you should call FindEtagFields on requested fields to find hasEtag, hasId, hasVer args.
// Optimization consists in reduced number of cycles running through requested fields. If you are going to use
// NormalizeEtag for each item in slice of your elements, FindEtagFields will be executed for each slice element.
func NormalizeEtags(t etag.EtagType, hasEtag bool, hasId bool, hasVer bool, etg *string, id *int64, ver *int32) (err error) {
	if hasEtag {
		*etg, err = etag.EncodeEtag(t, *id, *ver)
		if err != nil {
			return err
		}
		// hide
		if !hasId {
			*id = 0
		}
		if !hasVer {
			*ver = 0
		}
	}
	return err
}

// ParseIds converts a slice of strings (each possibly containing comma-separated eTags or numeric IDs)
// into a slice of int64. For example, given input ["1", "2,3", "etag4"], it converts each to int64 and returns []int64{1, 2, 3, 4}.
func ParseIds(input []string, etagType etag.EtagType) ([]int64, error) {
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

func SplitKnownAndUnknownFields(requestedFields []string, modelFields []string) (known []string, unknown []string) {
	for _, field := range requestedFields {
		var found bool
		for _, modelField := range modelFields {
			if field == modelField {
				known = append(known, field)
				found = true
				break
			}
		}
		if !found {
			unknown = append(unknown, field)
		}
	}
	return
}

func ContainsStringIgnoreCase(slice []string, target string) bool {
	targetLower := strings.ToLower(target)
	for _, str := range slice {
		if strings.ToLower(str) == targetLower {
			return true
		}
	}
	return false
}
