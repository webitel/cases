package scanner

import (
	"errors"
	"strings"

	_go "github.com/webitel/cases/api/cases"
)

// SourceTypeScanner is a custom scanner for converting a text-based enum to _go.SourceType
type SourceTypeScanner struct {
	SourceType *_go.SourceType
}

// Scan implements the sql.Scanner interface
func (s *SourceTypeScanner) Scan(src any) error {
	if src == nil {
		*s.SourceType = _go.SourceType_TYPE_UNSPECIFIED
		return nil
	}

	// Ensure the incoming value is a string
	str, ok := src.(string)
	if !ok {
		return errors.New("SourceTypeScanner: expected string, got different type")
	}

	// Convert string to enum
	typ, err := stringToType(str)
	if err != nil {
		*s.SourceType = _go.SourceType_TYPE_UNSPECIFIED
		return nil
	}

	*s.SourceType = typ
	return nil
}

// stringToType converts a string into the corresponding SourceType enum
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
		return _go.SourceType_TYPE_UNSPECIFIED, errors.New("invalid type value: " + typeStr)
	}
}
