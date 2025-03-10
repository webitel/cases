package scanner

import (
	"fmt"
	"github.com/webitel/cases/api/cases"
	_go "github.com/webitel/cases/api/cases"
	"strings"
)

// ScanSourceType is a custom scanner for SourceType enum.
func ScanSourceType(dest *cases.SourceType) func(src interface{}) error {
	return func(src interface{}) error {
		switch v := src.(type) {
		case nil:
			*dest = cases.SourceType_TYPE_UNSPECIFIED // Default type
		case string:
			typ, err := stringToType(v)
			if err != nil {
				return fmt.Errorf("invalid source type: %s", v)
			}
			*dest = typ
		default:
			return fmt.Errorf("unsupported type for SourceType: %T", v)
		}
		return nil
	}
}

// StringToType converts a string into the corresponding Type enum value.
//
// Types are specified ONLY for Source dictionary and are ENUMS in API.
func stringToType(typeStr string) (_go.SourceType, error) {
	switch strings.ToUpper(typeStr) {
	case "CALL":
		return _go.SourceType_CALL, nil
	case "CHAT":
		return _go.SourceType_CHAT, nil
	case "SOCIAL_MEDIA":
		return _go.SourceType_SOCIAL_MEDIA, nil
	case "EMAIL":
		return _go.SourceType_EMAIL, nil
	case "API":
		return _go.SourceType_API, nil
	case "MANUAL":
		return _go.SourceType_MANUAL, nil
	default:
		return _go.SourceType_TYPE_UNSPECIFIED, fmt.Errorf("invalid type value: %s", typeStr)
	}
}
