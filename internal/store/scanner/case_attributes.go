package scanner

import (
	"encoding/json"
	"fmt"
)

// ScanJSONToStructList scans a JSON array into a slice of structs.
func ScanJSONToStructList[T any](list *[]*T) any {
	return ScanFunc(func(src any) error {
		if src == nil {
			return nil
		}

		var rawJSON string
		switch v := src.(type) {
		case string:
			rawJSON = v
		case []byte:
			rawJSON = string(v)
		default:
			return fmt.Errorf("invalid type: expected string or []byte, got %T", src)
		}

		// Unmarshal the JSON array into the list of structs.
		var result []*T
		if err := json.Unmarshal([]byte(rawJSON), &result); err != nil {
			return fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		*list = result
		return nil
	})
}
