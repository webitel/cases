package graph

import (
	"fmt"
	"io"
	"strings"
	"unicode"
	"unicode/utf8"

	goerrors "github.com/pkg/errors"
	errors "github.com/webitel/cases/internal/errors"
)

// type FieldsQ map[string]*Query
// NOTE: ORDER MATTERS implementation !

// Fields[Q]uery represents a set of query{ fields }.
type FieldsQ []*Query

func (fs FieldsQ) MarshalText() ([]byte, error) {
	return fs.Encode()
}

func (fs FieldsQ) String() string {
	text, err := fs.MarshalText()
	if err != nil {
		panic(err) // FIXME
	}
	return string(text)
}

// Field[Q]uery
type Query struct {
	// REQUIRED. Name of the Query Object for output.
	Name string
	// OPTIONAL. Input arguments of the Query.
	Args
	// OPTIONAL. Nested Fields of the Query Object.
	Fields FieldsQ
}

func (req *Query) Clone() *Query {
	return &Query{
		Name:   req.Name,
		Args:   req.Args.Clone(),
		Fields: req.Fields.Clone(),
	}
}

type textWriter interface {
	io.Writer
	io.StringWriter
	// io.ByteWriter
	WriteRune(r rune) (n int, err error)
}

func (req *Query) MarshalText() ([]byte, error) {
	return req.Encode()
}

func (req *Query) String() string {
	text, err := req.MarshalText()
	if err != nil {
		panic(err) // FIXME
	}
	return string(text)
}

// ---------------------- parse -------------------------- //

// ParseFieldsQ parse ?fields= query specification.
// Format: name[.func([args..]),..][{inner,..}],..
func parseFieldsQ(dst *FieldsQ, src string) error {
	var (
		err error
		fq  Query
	)
	for src != "" {
		src, fq, err = scanFieldQ(src)
		if err != nil {
			return err
		}
		if fq.Name == "" {
			panic("fieldsQ: missing field name")
		}
		// snap := fq // snapshot
		// if !dst.Add(&snap) {
		if !dst.Add(fq) {
			return errors.NewBadRequestError(
				"api.graphql.fields.invalid",
				fmt.Sprintf("graphql: duplicate %s field",
					fq.Name),
			)
		}
		if src != "" && src[0] == ',' { // MAY: '}' for field{inner,..} spec
			src = src[1:] // COMMA
			continue
		}
		break
	}
	if src != "" {
		return errors.NewBadRequestError(
			"api.graph.fields.invalid",
			fmt.Sprintf("graphql: invalid syntax; char: %c",
				src[0]),
		)
	}
	return nil
}

// Fields returns set of FieldQ.Name
func (req FieldsQ) Fields() []string {
	n := len(req)
	if n > 0 {
		names := make([]string, 0, n)
		for _, fd := range req {
			names = append(names, fd.Name)
		}
		return names
	}
	return nil
}

// -1: NOT FOUND
func (req FieldsQ) Index(name string) int {
	if name == "" {
		return -1 // NOT FOUND
	}
	var (
		e  int
		n  = len(req)
		eq = func(s, t string) bool {
			return s == t
		}
	)
	for ; e < n && !eq(name, req[e].Name); e++ {
		// lookup: by field.name ...
	}
	if e < n {
		return e // FOUND
	}
	// if e == n {
	return -1 // NOT FOUND
	// }
}

func (vs FieldsQ) Has(name string) bool {
	return !(vs.Index(name) < 0)
	// _, ok := vs[name]
	// return ok
}

// <nil>: NOT FOUND
func (vs FieldsQ) Get(name string) *Query {
	e := vs.Index(name)
	if e < 0 {
		return nil
	}
	return vs[e]
	//return vs[name]
}

func (vs *FieldsQ) Add(out Query) bool {
	if out.Name == "" || vs.Has(out.Name) {
		return false
	}
	*(vs) = append(*(vs), &out)
	return true
}

func (vs FieldsQ) Clone() FieldsQ {
	n := len(vs)
	if n == 0 {
		return nil
	}
	v2 := make(FieldsQ, n)
	for name, query := range vs {
		v2[name] = query
	}
	return v2
}

const (
	DOT      = '.'
	COMMA    = ','
	PLUS     = '+'
	HYPHEN   = '-'
	EXCLAM   = '!'
	USCORE   = '_'
	LPAREN   = '('
	SQUOTE   = '\''
	DQUOTE   = '"'
	RPAREN   = ')'
	LCURLY   = '{'
	RCURLY   = '}'
	ESCAPE   = '\\'
	ASTERISK = '*'
)

func isASCIIDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isASCIIAlpha(r rune) bool {
	switch {
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	}
	return false
}

// InlineFieldsQ split given values into set of Field(s)Q selection.
func SplitFieldsQ(inline string) []string {
	if inline == "" {
		return nil
	}
	var (
		off  int      // offset reader
		size int      // character bytes size
		char rune     // read character
		last rune     // last character
		next int      // next field(word) offset
		list []string // exploded result

		top   = -1
		stack []rune
		quote = map[rune]rune{
			SQUOTE: SQUOTE,
			DQUOTE: DQUOTE,
			LPAREN: RPAREN,
			LCURLY: RCURLY,
		}
	)
	char, size = utf8.DecodeRuneInString(inline[off:])
	for ; off < len(inline); char, size = utf8.DecodeRuneInString(inline[off:]) {
		if char == utf8.RuneError {
			last = char
			off += size
			continue
		}
		if last == ESCAPE {
			last = char
			off += size
			continue
		}
		// SPLIT: by COMMA(,) or SPACE( )
		if char == COMMA || unicode.IsSpace(char) {
			// STACK: empty(?)
			if top < 0 {
				// SPLIT: OK(!)
				if next < off {
					list = append(list, inline[next:off])
				}
				next = off + size
			}
			// ELSE { AWAIT for trailing char(s) }
		} else if !(top < 0) && stack[top] == char {
			stack = stack[0:top]
			top-- // release latest
		} else if close, open := quote[char]; open {
			stack = append(stack, close)
			top++ // await trailing char
		}

		// switch char {
		// case ESCAPE:
		// 	//
		// case COMMA:
		// 	{
		// 		// if len(stack) > 0 {
		// 		if !(top < 0) {
		// 			// NOT EMPTY
		// 			break // switch
		// 		}
		// 		if next < off {
		// 			list = append(list, values[next:off])
		// 		}
		// 		next = off + size
		// 	}
		// default:
		// 	if !(top < 0) && stack[top] == char {
		// 		stack = stack[0:top]
		// 		top-- // release latest
		// 	} else if close, open := quote[char]; open {
		// 		stack = append(stack, close)
		// 		top++
		// 	}
		// }
		last = char
		off += size
	}
	if next < off {
		list = append(list, inline[next:off])
	}
	return list
}

func scanName(spec string) (rest, name string, err error) {
	var (
		off  int  // offset: spec[0:advance]; bytes count
		size int  // count: char(r) bytes count
		char rune // rune: Unicode Character
		// last rune //
	)
	rest = spec
	for len(rest) > 0 {
		char, size = utf8.DecodeRuneInString(rest)
		if char == utf8.RuneError {
			// err = goerrors.Errorf("fields: invalid UTF-8 sequence at position %d", e+1)
			err = goerrors.New("fields: invalid UTF-8 sequence")
			rest = "" // sanitize
			return    // "", "", err
		}

		switch {
		case char == USCORE:
			// MAY: __my___name__
			// if last > 0 && last == char {}
		case isASCIIDigit(char):
			if off == 0 {
				err = goerrors.New("fields: invalid name; leading DIGIT rune")
				rest = "" // sanitize
				return    // "", "", err
			}
		case isASCIIAlpha(char):
			// OK: advance !
		default:
			// LEADING exception(s)
			switch char {
			case PLUS, HYPHEN, EXCLAM:
				if off > 0 {
					err = goerrors.Errorf("fields: invalid name; illegal %c rune", char)
					rest = "" // sanitize
					return    // "", "", err
				}
				// SINGLE ( "+" ) ==> expands to all known fields set
				// LEADING ['+'|'-'|'!'] ==> sort order field spec
				// OK
			case ASTERISK:
			default:
				// Illegal field name char; stop
				name = spec[0:off]
				return // rest, name, nil
			}
		}
		rest = rest[size:] // advance
		off += size        // position
		// last = char        // remeber
	}
	name = spec
	return // "", "", nil
}

func scanFieldQ(spec string) (rest string, field Query, err error) {
	rest, field.Name, err = scanName(spec)
	if err != nil {
		return // "", {}, err
	}
	for rest != "" && rest[0] != COMMA {
		switch rest[0] {
		case DOT:
			{
				var fname string
				rest = rest[1:] // DOT
				rest, fname, err = scanName(rest)
				if err != nil {
					return // "", {}, err
				}
				if fname == "" {
					err = goerrors.New("fields: missing field.func name specification")
					return // "", {}, err
				}
				if rest == "" || rest[0] != LPAREN {
					err = goerrors.New("fields: malformed field.func name specification")
					return // "", {}, err
				}
				rest = rest[1:] // LPAREN
				rparen := strings.IndexByte(rest, RPAREN)
				if rparen < 0 {
					err = goerrors.New("fields: malformed field.func args specification")
					return
				}
				fargs := rest[0:rparen]
				err = field.Args.Set(fname, fargs)
				if err != nil {
					return //
				}
				// if Q.Funcs == nil {
				// 	Q.Funcs = map[string][]string{
				// 		fname: {fargs},
				// 	}
				// } else {
				// 	Q.Funcs[fname] = append(
				// 		Q.Funcs[fname], fargs,
				// 	)
				// }
				rest = rest[rparen+1:]
			}
		case LCURLY:
			{
				// var inner contactFieldQ
				rest = rest[1:] // LCURLY
				rest, field.Fields, err = scanFieldsQ(rest)
				if err == nil && (rest == "" || rest[0] != RCURLY) {
					err = goerrors.New("fields: malformed field{nested,...} spec")
				}
				if err != nil {
					return
				}
				rest = rest[1:] // RCURLY
				return          // MUST: This MIGHT be the last component !
			}
		default:
			return // rest, Q, nil
		}
	}
	return // rest, Q, nil
}

// ParseFieldQ parse single field
func ParseFieldQ(s string) (Query, error) {
	rest, field, err := scanFieldQ(s)
	if err == nil && rest != "" {
		err = goerrors.New("field: invalid query specification")
	}
	return field, err
}

func scanFieldsQ(s string) (rest string, fields FieldsQ, err error) {
	rest = s
	var query Query
	for rest != "" {
		rest, query, err = scanFieldQ(rest)
		if err != nil {
			return
		}
		if query.Name == "" {
			panic("fieldsQ: missing field name")
		}
		// snap := spec // snapshot
		// fields = append(fields, &snap)
		fields.Add(query)
		if rest != "" && rest[0] == ',' { // MAY: '}' for field{inner,..} spec
			rest = rest[1:] // COMMA
			continue
		}
		break
	}
	return // rest, fields, nil
}
