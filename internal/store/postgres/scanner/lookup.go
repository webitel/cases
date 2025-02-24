package scanner

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgtype"
	_go "github.com/webitel/cases/api/cases"
)

type ScanFunc func(src any) error

func (s ScanFunc) Scan(src any) error {
	return s(src)
}

type TextDecoder func(src []byte) error

func (dec TextDecoder) DecodeText(_ *pgtype.ConnInfo, src []byte) error {
	if dec != nil {
		return dec(src)
	}
	return nil
}

func (dec TextDecoder) Scan(src any) error {
	if src == nil {
		return dec.DecodeText(nil, nil)
	}

	switch data := src.(type) {
	case string:

		return dec.DecodeText(nil, []byte(data))
	case []byte:
		text := make([]byte, len(data))
		copy(text, data)
		return dec.DecodeText(nil, text)
	}

	return fmt.Errorf("text_decoder: cannot scan %T value %[1]v into %T", src, dec)
}

// preprocessCompositeString normalizes composite strings by unescaping quotes
// and removing quotes around the composite values.
func preprocessCompositeString(input string) string {
	// Remove escaped quotes
	cleaned := strings.ReplaceAll(input, `\"`, `"`)

	// Check for the pattern of a quoted composite value and clean it
	if strings.HasPrefix(cleaned, "(") && strings.HasSuffix(cleaned, ")") {
		// Extract the inner part of the composite
		inner := cleaned[1 : len(cleaned)-1]

		// Split into parts and trim quotes if they exist
		parts := strings.Split(inner, ",")
		for i, part := range parts {
			// Trim leading/trailing quotes only from the name field
			if strings.HasPrefix(part, `"`) && strings.HasSuffix(part, `"`) {
				parts[i] = part[1 : len(part)-1]
			}
		}

		// Reconstruct the composite
		cleaned = fmt.Sprintf("(%s)", strings.Join(parts, ","))
	}

	return cleaned
}

func ScanRowLookup(value **_go.Lookup) any {
	return TextDecoder(func(src []byte) error {
		res := *(value)
		*(value) = nil

		if len(src) == 0 {
			return nil // NULL
		}

		if res == nil {
			res = new(_go.Lookup)
		}

		var (
			ok  bool
			str pgtype.Text
			row = []pgtype.TextDecoder{
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					id, err := strconv.ParseInt(str.String, 10, 64)
					if err != nil {
						return err
					}
					res.Id = id
					return nil
				}),
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Name = str.String
					ok = ok || (str.String != "" && str.String != "[deleted]") // && str.Status == pgtype.Present
					return nil
				}),
			}
			raw = pgtype.NewCompositeTextScanner(nil, src)
		)

		var err error
		for _, col := range row {

			raw.ScanDecoder(col)

			err = raw.Err()
			if err != nil {
				return err
			}
		}

		if ok {
			*(value) = res
		}

		return nil
	})
}

func ScanRowExtendedLookup(value **_go.ExtendedLookup) any {
	return TextDecoder(func(src []byte) error {
		res := *(value)
		*(value) = nil

		if len(src) == 0 {
			return nil // NULL
		}

		if res == nil {
			res = new(_go.ExtendedLookup)
		}

		var (
			ok  bool
			str pgtype.Text
			row = []pgtype.TextDecoder{
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					id, err := strconv.ParseInt(str.String, 10, 64)
					if err != nil {
						return err
					}
					res.Id = id
					return nil
				}),
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Name = str.String
					ok = ok || (str.String != "" && str.String != "[deleted]") // && str.Status == pgtype.Present
					return nil
				}),
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Type = str.String
					ok = ok || (str.String != "" && str.String != "[deleted]") // && str.Status == pgtype.Present
					return nil
				}),
			}
			raw = pgtype.NewCompositeTextScanner(nil, src)
		)

		var err error
		for _, col := range row {

			raw.ScanDecoder(col)

			err = raw.Err()
			if err != nil {
				return err
			}
		}

		if ok {
			*(value) = res
		}

		return nil
	})
}

func ScanRelatedCaseLookup(value **_go.RelatedCaseLookup) any {
	return TextDecoder(func(src []byte) error {
		res := *(value)
		*(value) = nil

		if len(src) == 0 {
			return nil // NULL
		}

		if res == nil {
			res = new(_go.RelatedCaseLookup)
		}

		var (
			ok  bool
			str pgtype.Text
			num pgtype.Int8
			row = []pgtype.TextDecoder{
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := num.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Id = num.Int
					return nil
				}),
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Name = str.String
					ok = ok || (str.String != "" && str.String != "[deleted]")
					return nil
				}),
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Subject = str.String
					ok = ok || (str.String != "" && str.String != "[deleted]")
					return nil
				}),
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					ver, parseErr := strconv.Atoi(str.String)
					if parseErr != nil {
						return parseErr
					}
					res.Ver = int32(ver)
					return nil
				}),
				TextDecoder(func(src []byte) error {
					if len(src) == 0 {
						return nil
					}
					err := str.DecodeText(nil, src)
					if err != nil {
						return err
					}
					res.Color = str.String
					ok = ok || str.String != ""
					return nil
				}),
			}
			raw = pgtype.NewCompositeTextScanner(nil, src)
		)

		var err error
		for _, col := range row {
			raw.ScanDecoder(col)
			err = raw.Err()
			if err != nil {
				return err
			}
		}

		if ok {
			*(value) = res
		}

		return nil
	})
}

func ScanLookupList(value *[]*_go.Lookup) any {
	return TextDecoder(func(src []byte) error {
		if len(src) == 0 {
			// If the source is empty, set the value to an empty list
			*value = []*_go.Lookup{}
			return nil
		}

		// Decode the text array from the source
		var array pgtype.TextArray
		if err := array.DecodeText(nil, src); err != nil {
			return fmt.Errorf("failed to decode text array: %w", err)
		}

		// Prepare the slice to store Lookup objects
		lookupList := make([]*_go.Lookup, len(array.Elements))
		for i, element := range array.Elements {
			if element.String == "" {
				lookupList[i] = nil
				continue
			}

			// Remove parentheses and split the composite element (e.g., "(1,Name)")
			trimmed := element.String[1 : len(element.String)-1] // Remove enclosing parentheses
			parts := splitCompositeFields(trimmed)

			if len(parts) < 2 {
				return fmt.Errorf("invalid composite format: %s", element.String)
			}

			// Parse ID and Name
			id, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse ID: %w", err)
			}
			name := parts[1]

			lookupList[i] = &_go.Lookup{
				Id:   id,
				Name: name,
			}
		}

		*value = lookupList
		return nil
	})
}

// Helper function to split composite fields, considering possible commas inside quoted strings.
func splitCompositeFields(composite string) []string {
	var parts []string
	var currentField strings.Builder
	inQuotes := false

	for i := 0; i < len(composite); i++ {
		c := composite[i]

		switch c {
		case ',':
			if inQuotes {
				currentField.WriteByte(c)
			} else {
				parts = append(parts, currentField.String())
				currentField.Reset()
			}
		case '"':
			// Check for escaped quotes (e.g., `""` inside a quoted field)
			if inQuotes && i+1 < len(composite) && composite[i+1] == '"' {
				currentField.WriteByte('"')
				i++ // Skip the next quote
			} else {
				inQuotes = !inQuotes
			}
		default:
			currentField.WriteByte(c)
		}
	}

	// Add the last field (even if empty)
	parts = append(parts, currentField.String())

	return parts
}
