package store

import "fmt"

// Ident returns a string that represents a qualified identifier.
// For example, Ident("cc", "case_link") returns "cc.case_link".
var Ident = func(left, right string) string {
	return fmt.Sprintf("%s.%s", left, right)
}
