package util

import (
	"io"
	"strings"
)

// CompactSQL formats given SQL text to compact form.
// - replaces consecutive white-space(s) with single SPACE(' ')
// - suppress single-line comment(s), started with -- up to [E]nd[o]f[L]ine
// - suppress multi-line comment(s), enclosed into /* ... */ pairs
// - transmits literal '..' or ".." sources in their original form
// https://www.postgresql.org/docs/current/sql-syntax-lexical.html#SQL-SYNTAX-OPERATORS
func CompactSQL(s string) string {
	var (
		r = strings.NewReader(s)
		w strings.Builder
	)

	w.Grow(int(r.Size()))

	var (
		err  error
		char rune
		last rune
		hold rune

		isSpace = func() (is bool) {
			switch char {
			case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
				is = true
			}
			return // false
		}
		isPunct = func(char rune) (is bool) {
			switch char {
			// none; start of text
			case 0:
				is = true
			// special
			// ':' USES [squirrel] for :named parameters,
			//     so we need to keep SPACE if there were any
			case ',', '(', ')', '[', ']', ';', '\'', '"': // , ':':
				is = true
			// operators
			case '+', '-', '*', '/', '<', '>', '=', '~', '!', '@', '#', '%', '^', '&', '|':
				is = true
			}
			return // false
		}
		isQuote = func() (is bool) {
			switch char {
			case '\'', '"': // SQUOTE, DQUOTE:
				is = true
			}
			return // false
		}
		// context
		space   bool // [IN] [w]hite[sp]ace(s)
		quote   rune // [IN] [l]i[t]eral(s); *QUOTE(s)
		comment rune // [IN] [c]o[m]ment; [-|*]
		// helpers
		isComment = func() bool {
			switch comment {
			case '-':
				{
					// comment: close(\n)
					if char == '\n' { // EOL
						space = true // inject
						comment = 0  // close
						hold = 0     // clear
					}
					return true // still IN ...
				}
			case '*':
				{
					// comment: close(*/)
					if hold == 0 && char == '*' {
						// MAY: close(*/)
						hold = char
						// need more data ...
					} else if hold == '*' && char == '/' {
						space = true // inject
						comment = 0  // close
						hold = 0     // clear
					}
					return true // still IN ...
				}
				// default: 0
			}
			// NOTE: (comment == 0)
			switch hold {
			// comment: start(--)
			case '-': // single-line
				{
					if char == hold {
						hold = 0       // clear
						comment = char // start
						return true
					}
					return false
				}
			// comment: start(/*)
			case '/': // multi-line
				{
					if char == '*' {
						hold = 0       // clear
						comment = char // start
						return true
					}
					return false
				}
			case 0:
				{
					// NOTE: (hold == 0)
					switch char {
					case '-':
					case '/':
					default:
						// NOT alike ...
						return false
					}
					// need more data ...
					hold = char
					// DO NOT write(!)
					return true
				}
			default:
				{
					// NO match
					// need to write hold[ed] char
					return false
				}
			}
		}
		isLiteral = func() bool {
			if !isQuote() || last == '\\' { // ESC(\')
				return quote > 0 // We are IN ?
			}
			// close(?)
			if quote == char { // inLiteral(?)
				quote = 0
				return true // as last
			}
			// start(!)
			quote = char
			return true
		}
		// [re]write
		output = func() {
			if hold > 0 {
				w.WriteRune(hold)
				last = hold
				hold = 0
			}
			if space {
				space = false
				if !isPunct(last) && !isPunct(char) {
					w.WriteRune(' ') // INJECT SPACE(' ')
				}
			}
			w.WriteRune(char)
			last = char
		}
	)

	var e int
	for {

		char, _, err = r.ReadRune()
		if err != nil {
			break
		}
		e++ // char index position

		if isComment() {
			// suppress; DO NOT write(!)
			continue
		}

		if isLiteral() {
			// [re]write: as is (!)
			output()
			continue
		}

		if isSpace() {
			// fold sequence ...
			space = true
			continue
		}
		// [re]write: [hold]char
		output()
	}

	if err != io.EOF {
		panic(err)
	}

	return w.String()
}
