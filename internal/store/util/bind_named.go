package util

import (
	"strconv"
	"unicode"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

// Allowed runes for named parameters
var allowedBindRunes = []*unicode.RangeTable{unicode.Letter, unicode.Digit}

// compileNamedQuery converts a named query into a positional query and extracts parameter names.
func compileNamedQuery(qs []byte, bindType int) (query string, names []string, err error) {
	names = make([]string, 0, 10)
	rebound := make([]byte, 0, len(qs))
	inName := false
	name := make([]byte, 0, 10)

	for i := 0; i < len(qs); i++ {
		b := qs[i]
		if b == ':' {
			// Check for "::" and append as is
			if i+1 < len(qs) && qs[i+1] == ':' {
				if inName {
					// Close the parameter name before appending "::"
					inName = false
					names = append(names, string(name))
					switch bindType {
					case sqlx.DOLLAR:
						rebound = append(rebound, '$')
						rebound = append(rebound, []byte(strconv.Itoa(len(names)))...)
					default:
						rebound = append(rebound, '?')
					}
				}
				// Append "::" directly
				rebound = append(rebound, ':', ':')
				i++ // Skip the next ':'
				continue
			}
			// Start of a named parameter
			if inName {
				err = errors.New("unexpected `:` while reading named param")
				return query, names, err
			}
			inName = true
			name = []byte{}
		} else if inName && (unicode.IsOneOf(allowedBindRunes, rune(b)) || b == '_' || b == '.') {
			name = append(name, b)
		} else if inName {
			// End of the parameter name
			inName = false
			names = append(names, string(name))
			switch bindType {
			case sqlx.DOLLAR:
				rebound = append(rebound, '$')
				rebound = append(rebound, []byte(strconv.Itoa(len(names)))...)
			default:
				rebound = append(rebound, '?')
			}
			rebound = append(rebound, b)
		} else {
			rebound = append(rebound, b)
		}
	}

	if inName {
		// Close the final parameter name
		names = append(names, string(name))
		if bindType == sqlx.DOLLAR {
			rebound = append(rebound, '$')
			rebound = append(rebound, []byte(strconv.Itoa(len(names)))...)
		} else {
			rebound = append(rebound, '?')
		}
	}

	return string(rebound), names, nil
}

// bindMapArgs maps the values of named parameters into a slice of arguments.
func bindMapArgs(names []string, params map[string]interface{}) ([]interface{}, error) {
	arglist := make([]interface{}, 0, len(names))
	for _, name := range names {
		val, ok := params[name]
		if !ok {
			return nil, errors.Errorf("missing parameter: %s", name)
		}
		arglist = append(arglist, val)
	}
	return arglist, nil
}

// BindNamed transforms a named query into a positional query with its arguments.
func BindNamed(query string, params map[string]interface{}) (string, []interface{}, error) {
	bound, names, err := compileNamedQuery([]byte(query), sqlx.DOLLAR)
	if err != nil {
		return "", nil, err
	}

	arglist, err := bindMapArgs(names, params)
	return bound, arglist, err
}
