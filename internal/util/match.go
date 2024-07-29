package util

import "strings"

// Match provides methods for substring manipulation and matching.
type Match struct{}

// Substrings represents a slice of substrings.
type Substrings []string

// SubstringMask splits the input string s into substrings using the provided any and one runes as delimiters.
func (m Match) SubstringMask(s string, any, one rune) Substrings {
	if any == 0 {
		any = '*'
	}
	if one == 0 {
		one = '?'
	}
	sv := strings.Split(s, string(any))
	// omit any empty sequences: [1:len()-2]
	for i := len(sv) - 2; i > 0; i-- {
		if len(sv[i]) == 0 {
			// cut
			sv = append(sv[:i], sv[i+1:]...)
		}
	}
	return Substrings(sv)
}

// Substring splits the input string s into substrings using default delimiters '*' and '?'.
func (m Match) Substring(s string) Substrings {
	return m.SubstringMask(s, '*', '?')
}
