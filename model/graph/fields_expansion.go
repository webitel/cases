package graph

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	goerrors "errors"
)

// FieldExpansion syntax options.
type FieldExpansion struct {
	// Invalidate field.query(args) expansion
	NoArgs bool
	// Invalidate field{nested,..} expansion
	NoNested bool
	// Default fields set if none specified
	// Used to reassign Metadata.Default set.
	Default []string
	// Sorting fields specification mode ?
	//
	// sort   ::= spec *( "," spec )
	// spec   ::= [ ORDER ] name / nest
	// nest   ::= name "{" nest / sort "}"
	// name   ::= ALPHA / DIGIT / USCORE ; field name
	//
	// ALPHA    = %x41-5A / %x61-7A  ; "A"-"Z" / "a"-"z"
	// DIGIT    = %x30-39            ; "0"-"9"
	// USCORE   = %x5F ; underscore  ; "_"
	//
	// Examples:
	//
	// ?sort=+id,names{!common_name,id}
	// ?sort=id&sort=names{-common_name,id}; affected as above
	//
	// Invalid:
	// ?sort=+id,!names{-common_name,id}; illegal "!"; (nested only) fields MAY be sorted
	//
	// If false - regular GET fields specification mode is used.
	Sorting bool
}

// FieldEncoding as an expansion option
type FieldEncoding func(syntax *FieldExpansion)

func NewFieldExpansion(options ...FieldEncoding) (syntax FieldExpansion) {
	for _, option := range options {
		option(&syntax)
	}
	return // encoding
}

// // FieldEncoding Option
// type FieldExpansion interface {
// 	configure(config *fieldExpansion)
// }

// type fieldExpansionOption func(config *fieldExpansion)

// func (setup fieldExpansionOption) configure(config *fieldExpansion) {
// 	if setup != nil {
// 		setup(config)
// 	}
// }

// NoArgs FieldEncoding Option.
// Disallow field[.arg(input)]* component(s)
func NoArgs() FieldEncoding {
	return func(syntax *FieldExpansion) {
		syntax.NoArgs = true
	}
}

// NoQuery FieldEncoding Option.
// Does NOT expose .query(),.. component(s)
func NoNested() FieldEncoding {
	return func(syntax *FieldExpansion) {
		syntax.NoNested = true
	}
}

// Sorting FieldEncoding Option.
// Cause to validate ?sort= fields spec
func Sorting() FieldEncoding {
	return func(syntax *FieldExpansion) {
		syntax.Sorting = true
		syntax.NoArgs = true
	}
}

// DefaultFields FieldEncoding Option.
// If no ?fields= were specified,
// this set of `q` will be used.
func DefaultFields(q ...string) FieldEncoding {
	return func(spec *FieldExpansion) {
		var e, n int
		for _, s := range q {
			_, name, err := scanName(s)
			if err != nil {
				panic(fmt.Errorf("fields: default %s invalid spec", s))
			}
			for e, n = 0, len(spec.Default); e < n; e++ {
				if strings.HasPrefix(spec.Default[e], name) {
					break // duplicate found
				}
			}
			if e == n {
				spec.Default = append(
					spec.Default, s,
				)
			}
		}
	}
}

func scanField(syntax *FieldExpansion, text string) (rest string, field Query, err error) {
	rest, field.Name, err = scanName(text)
	// Validate field.Name
	if err == nil {
		if field.Name == "" {
			abrev := text
			if len(abrev) > 5 {
				abrev = abrev[0:3] + ".."
			}
			err = fmt.Errorf("fields: invalid spec; want = %s, got = %s", "name", abrev)
		} else {
			if !syntax.Sorting && field.Name == string(PLUS) {
				// NO MORE .option(args){nested,..} allowed
				return // rest, field, nil
			}
			switch field.Name[0] {
			case PLUS, HYPHEN, EXCLAM:
				if !syntax.Sorting {
					err = fmt.Errorf("fields: invalid spec; want = name, got = %s", field.Name)
				}
			}
		}
	}
	if err != nil {
		return // "", {}, err
	}
	for rest != "" && rest[0] != COMMA {
		switch rest[0] {
		case DOT:
			{
				if syntax.NoArgs || syntax.Sorting {
					err = fmt.Errorf("field: query .arg(input) not allowed")
					return // rest, field, err
				}
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
				if syntax.NoNested {
					err = fmt.Errorf("field: query {nested} not allowed")
					return // rest, field, err
				}
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

// parseInlineFields parse ?fields= query specification.
// Format: name[.func([args..]),..][{inner,..}],..
func parseFields(spec *FieldExpansion, into *FieldsQ, s string) error {
	var (
		err error
		fq  Query
	)
	for s != "" {
		s, fq, err = scanField(spec, s)
		if err != nil {
			return err
		}
		if fq.Name == "" {
			panic("fieldsQ: missing field name")
		}
		// if extend && into.Has(fq.Name) {
		// 	continue // omit: defined
		// }
		// snap := fq // snapshot
		// if !into.Add(&snap) {
		if !into.Add(fq) {
			return errors.New(
				fmt.Sprintf("fields: duplicate %s field",
					fq.Name),
			)
		}
		if s != "" && s[0] == COMMA { // MAY: '}' for field{inner,..} spec
			s = s[1:] // COMMA
			continue
		}
		break
	}
	if s != "" {
		return errors.New(
			fmt.Sprintf("fields: invalid syntax; char: %c",
				s[0]),
		)
	}
	return nil
}

func (fs *FieldsQ) decode(spec *FieldExpansion, vs []string) (err error) {
	// Default ?
	if len(vs) == 0 {
		vs = spec.Default
	}
	for _, s := range vs {
		err = parseFields(spec, fs, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fq *Query) encode(syntax *FieldExpansion, text textWriter) (n int, err error) {
	if fq == nil || fq.Name == "" {
		return // 0, nil
	}
	// text.Grow()
	var (
		c         int
		writeRune = func(r rune) error {
			c, err = text.WriteRune(r)
			n += c
			return err
		}
		writeString = func(s string) error {
			c, err = text.WriteString(s)
			n += c
			return err
		}
	)

	err = writeString(fq.Name)
	if err != nil {
		return // n, err
	}
	if !syntax.NoArgs && len(fq.Args) > 0 {
		for param, value := range fq.Args {
			if param == "" {
				continue
			}
			err = writeRune(DOT)
			if err != nil {
				return // n, err
			}
			err = writeString(param)
			if err != nil {
				return // n, err
			}
			err = writeRune(LPAREN)
			if err != nil {
				return // n, err
			}
			if value != nil {
				c, err = fmt.Fprintf(text, "%v", value)
				if n += c; err != nil {
					return // n, err
				}
			}
			err = writeRune(RPAREN)
			if err != nil {
				return // n, err
			}
		}
	}
	if !syntax.NoNested && len(fq.Fields) > 0 {
		err = writeRune(LCURLY)
		if err != nil {
			return // n, err
		}
		c = 0 // track: something was written ?
		for _, fx := range fq.Fields {
			if c > 0 {
				err = writeRune(COMMA)
				if err != nil {
					return // n, err
				}
			}
			c, err = fx.encode(syntax, text)
			if n += c; err != nil {
				return // n, err
			}
		}
		err = writeRune(RCURLY)
		if err != nil {
			return // n, err
		}
	}
	return // n, nil
}

func (fs FieldsQ) encode(syntax *FieldExpansion, text textWriter) (n int, err error) {
	if len(fs) == 0 {
		return // 0, nil
	}
	var c int
	for _, fq := range fs {
		if c > 0 {
			c, err = text.WriteRune(COMMA)
			if n += c; err != nil {
				return // n, err
			}
		}
		c, err = fq.encode(syntax, text)
		if n += c; err != nil {
			return // n, err
		}
	}
	return // n, nil
}

// ParseFields decodes ?fields=name[.query([args..]),..][{nested,..}][,..] string specification.
func ParseFields(s string, decode ...FieldEncoding) (fields FieldsQ, err error) {
	err = fields.Parse(s, decode...)
	return // fields, err
}

// ParseFieldsQuery acts like ParseFields method, but decodes ?fields=[&fields=] query specification.
func ParseFieldsQuery(vs []string, decode ...FieldEncoding) (fields FieldsQ, err error) {
	err = fields.ParseQuery(vs, decode...)
	return // fields, err
}

func (fs *FieldsQ) Parse(s string, decode ...FieldEncoding) error {
	syntax := NewFieldExpansion(decode...)
	return fs.decode(&syntax, []string{s})
}

func (fs *FieldsQ) ParseQuery(vs []string, decode ...FieldEncoding) error {
	syntax := NewFieldExpansion(decode...)
	return fs.decode(&syntax, vs)
}

func (fd *Query) Encode(syntax ...FieldEncoding) (text []byte, err error) {
	if fd == nil || fd.Name == "" {
		return // 0, nil
	}
	out := bytes.NewBuffer(text)
	enc := NewFieldExpansion(syntax...)
	_, err = fd.encode(&enc, out)
	if err == nil {
		text = out.Bytes()
	}
	return // text, err
}

func (fs FieldsQ) Encode(syntax ...FieldEncoding) (text []byte, err error) {
	if len(fs) == 0 {
		return // nil, nil
	}
	out := bytes.NewBuffer(text)
	enc := NewFieldExpansion(syntax...)
	_, err = fs.encode(&enc, out)
	if err == nil {
		text = out.Bytes()
	}
	return // text, err
}
