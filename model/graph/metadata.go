package graph

import (
	"fmt"
	"strings"

	"errors"
)

// GraphQL Metadata about an Object.
type Metadata struct {
	// Name of the Object
	Name string `json:"name"`
	// Type of the Output data
	Type string `json:"typeOf"`
	// OPTIONAL. Input Arguments supported
	Args InputArgs `json:"input,omitempty"`
	// OPTIONAL. Output Fields of the Object Type
	Fields []*Metadata `json:"fields,omitempty"`
	// OPTIONAL. Default set of output Fields; expansion of ('*')
	Default []string `json:"output,omitempty"`
	// Resolve Output data or an error.
	Resolve OutputFunc `json:"-"`
}

// FieldsQ validates and normalize given `fs` query fields
func (md *Metadata) GetField(name string) *Metadata {
	var (
		e int // index
		n = len(md.Fields)
	)
	for ; e < n && md.Fields[e].Name != name; e++ {
		// lookup: field by name
	}
	if e < n {
		return md.Fields[e]
	}
	return nil
}

// GetQuery validates and normalize given `req` FieldQ of this Metadata configuration.
func (md *Metadata) FieldQ(req *Query, syntax ...FieldEncoding) error {
	spec := NewFieldExpansion(syntax...)
	return md.fieldQ(&spec, req, false)
}

// FieldQ validates and normalize given `req` query of this Metadata.
// <all> - means to disclose ALL('+') nested fields into `req`.
func (md *Metadata) fieldQ(spec *FieldExpansion, req *Query, all bool) (err error) {
	// name
	// if req.Name != md.Name {
	if !strings.HasSuffix(req.Name, md.Name) || len(req.Name)-len(md.Name) > 1 {
		return errors.New(
			fmt.Sprintf("graphql: %s; expect: %s",
				req.Name, md.Name),
		)
		// return fmt.Errorf("fields: invalid setup; want = %s, got = %s", md.Name, req.Name)
	}
	// arguments
	if spec.NoArgs {
		if n := len(req.Args); n > 0 {
			return errors.New(
				fmt.Sprintf("graphql: %s( args:[%d] ); expect: no arguments",
					md.Name, n),
			)
		}
	} else {
		// validate & normalize & defaults
		err = md.Args.Parse(req)
		if err != nil {
			return // err
		}
	}
	// nested
	if spec.NoNested {
		if len(req.Fields) > 0 {
			return errors.New(
				fmt.Sprintf("graphql: %s{ fields.. }; expect: no nested fields",
					md.Name),
			)
			// return fmt.Errorf("query: %s; input: fields {nested,..} not allowed", req.Name)
		}
	} else {
		err = md.fieldsQ(spec, &req.Fields, all)
		if err != nil {
			return // err
		}
	}
	// ok
	return nil
	// panic("not implemented")
}

// <sup> ==> def(*)
// <nil> ==> all(+)
func (md *Metadata) extendQ(spec *FieldExpansion, req *FieldsQ, sup []string) (err error) {

	var (
		fields = *(req)
		nested FieldsQ
		nestAs = *(spec) // snap
	)
	if len(sup) > 0 {
		for _, input := range sup {
			// parse input complex, e.g.: name{given,common,updated_by{name}}
			err = parseFields(spec, &nested, input)
			if err != nil {
				return err
			}
			// extend(name, nil)
		}
		for e := 0; e < len(nested); e++ {
			// remove duplicates ...
			if fields.Has(nested[e].Name) {
				nested = append(nested[0:e], nested[e+1:]...)
				e--
				continue
			}
		}
		// Validate & normalize
		nestAs.Default = md.Default
		err = md.fieldsQ(&nestAs, &nested, false) // all: auto-detect
		if err != nil {
			return err // invalid set of fields: spec.Default -or- Metadata.Default
		}
		fields = append(fields, nested...)
	} else {

		for _, fd := range md.Fields {
			// extend(fd.Name, fd)
			if fields.Has(fd.Name) {
				// omit; already has
				continue
			}
			with := &Query{
				Name: fd.Name,
			}
			nestAs.Default = fd.Default
			err = fd.fieldQ(&nestAs, with, false)
			if err != nil {
				return err
			}
			fields = append(fields, with)
		}
	}

	*(req) = fields
	return nil
}

// FieldQ validates and normalize given `req` query of this Metadata.
func (md *Metadata) fieldsQ(spec *FieldExpansion, req *FieldsQ, all bool) (err error) {
	var (
		fields = *(req)  // early binding
		nested = *(spec) // field expansion
	)
	// FIXME: def.NoNested is omitted !
	if len(md.Fields) == 0 {
		if len(fields) > 0 {
			return errors.New(
				fmt.Sprintf("graphql: %s{ fields.. }; scalar: has no fields",
					md.Name),
				// "%s( fields:.. ); scalar( %s ): has no fields",
				// md.Name, md.Type,
			)
		}
		// Has NO inner spec !
		return nil
	}

	var (
		fx rune      // prefix: [ + | - | ! ]
		fn string    // field.name
		fd *Metadata // field.spec
		ex rune      // supplement level: [ '*' | '+' ]
	)
	// for e, fq := range vs {
next:
	// for name, query := range fields {
	for e := 0; e < len(fields); e++ {
		// fq := list[e]
		query := fields[e]
		// valid: field name ?
		fx, fn = 0, query.Name
		if len(fn) > 0 {
			fx = rune(fn[0])
		}
		switch fx {
		case PLUS, HYPHEN, EXCLAM:
			// fx, fn = rune(fn[0]), fn[1:]
			if spec.Sorting {
				fn = fn[1:] // accept; ignore while lookup field
			} else {
				// ?fields=%2b ; expands to ALL known md.Fields
				if len(fn) == 1 && fx == PLUS {
					// TODO: Expand to ALL known md.Fields
					ex = PLUS // ALL
					// all = true
					// // delete(fields, name)
					fields = append(fields[0:e], fields[e+1:]...)
					e--
					continue next
				}
			}
		case ASTERISK:
			// ?fields=* ; expands to DEFAULT set of fields
			if !spec.Sorting && len(fn) == 1 {
				// TODO: Expand to DEFAULT set of md.Fields
				if ex != PLUS {
					ex = ASTERISK // DEF
					// sup := spec.Default
					// if len(sup) == 0 {
					// 	sup = md.Default
					// }
					// vs = append(vs[0:e], vs[e+1:]...)
					// err = vs.decode(specreq.FieldsQ, sup) // FIXME: duplicate fields
					// if err != nil {
					// 	return err
					// }
					//
				}
				// delete(fields, name)
				fields = append(fields[0:e], fields[e+1:]...)
				e--
				continue next
				// break // OK: continue
			}
		}

		if fn == "" {
			return errors.New(
				fmt.Sprintf("graphql: %s( fields:[%d] ); missing",
					md.Name, e+1),
			)
		}
		fd = md.GetField(fn)
		if fd == nil {
			return errors.New(
				fmt.Sprintf("graphql: %s{ %s }; no such field",
					md.Name, query.Name),
			)
		}
		nested.Default = fd.Default
		err = fd.fieldQ(&nested, query, all)
		if err != nil {
			return err
		}
	}

	// if len(vs) == 0 && len(md.Default) > 0 {
	// 	// Default: core fields
	// 	var (
	// 		node *FieldQ
	// 		page = make([]FieldQ, len(md.Default))
	// 	)
	// 	if cap(vs) < len(md.Default) {
	// 		vs = make(FieldsQ, 0, len(md.Default))
	// 	}
	// 	for _, name := range md.Default {
	// 		spec := md.GetField(name)
	// 		if len(page) > 0 {
	// 			node = &page[0]
	// 			page = page[1:]
	// 		} else {
	// 			node = new(FieldQ)
	// 		}
	// 		node.Name = name
	// 		err = spec.fieldQ(spec, node)
	// 		if err != nil {
	// 			// Default: MUST not trigger an error !
	// 			panic(err)
	// 		}
	// 		vs = append(vs, node)
	// 	}
	// }

	// Default('*') if nothing specified
	if len(fields) == 0 && ex < ASTERISK {
		ex = ASTERISK
	}

	switch ex {
	case PLUS: // ALL KNOWN
		err = md.extendQ(spec, &fields, nil) // TODO: expand all fields recurcivly
	case ASTERISK: // DEFAULT
		sup := spec.Default
		if len(sup) == 0 {
			sup = md.Default
		}
		if len(sup) > 0 {
			err = md.extendQ(spec, &fields, sup)
		}
	}

	if err != nil {
		return err
	}

	*(req) = fields
	return nil
	// panic("not implemented")
}

// FieldsQ presets ( validate & normalize & defaults )
// given FieldsQ `req` according to this `md` Metadata.Fields
// and optional `syntax` options configuration.
func (md *Metadata) FieldsQ(req *FieldsQ, syntax ...FieldEncoding) error {
	_ = *(req) // early binding
	spec := NewFieldExpansion(syntax...)
	return md.fieldsQ(&spec, req, false) // all: auto-detect
}

// ParseFields decodes ?fields=[&fields=] query value(s) spec.
// Validate & normalize expansion due to `md`.Fields configuration.
func (md *Metadata) ParseFields(vs []string, decode ...FieldEncoding) (fields FieldsQ, err error) {
	// Parse given vs inline ?fields=[&fields=] expansion
	spec := NewFieldExpansion(decode...)
	// fields, err := ParseFieldsQuery(vs, decode...)
	if len(vs) == 0 {
		vs = spec.Default // manual
		if len(vs) == 0 {
			vs = md.Default // defined
		}
	}
	err = fields.decode(&spec, vs)
	if err == nil {
		// Validate & normalize
		err = md.fieldsQ(&spec, &fields, false) // all: auto-detect
	}
	return fields, err
}
