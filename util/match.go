package util

import "strings"

// // Match provides methods for substring manipulation and matching.
// type Match struct{}

// Substrings represents a slice of substrings.
type Substrings []string

// SubstringMask splits the input string s into substrings using the provided any and one runes as delimiters.
func SubstringMask(s string, any, one rune) Substrings {
	if any == 0 {
		any = '*'
	}
	if one == 0 {
		one = '?'
	}
	// Replace '?' with '_' for SQL LIKE operator (matches exactly one character)
	s = strings.ReplaceAll(s, string(one), "_")
	// Replace '*' with '%' for SQL LIKE operator (matches zero or more characters)
	s = strings.ReplaceAll(s, string(any), "%")
	return Substrings{strings.TrimSpace(s)}
}

// Substring splits the input string s into substrings using default delimiters '*' and '?'.
func Substring(s string) Substrings {
	return SubstringMask(s, '*', '?')
}
