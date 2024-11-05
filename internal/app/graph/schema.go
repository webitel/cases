package graph

type CaseTypes struct {
	Link    Link
	Comment Comment
	Inited  bool
}

type ScalarTypes struct {
	Lookup Lookup
	Tuple  Tuple
	Inited bool
}

type schema struct {
	// types
	Scalar ScalarTypes // default
	Case   CaseTypes   // case related types
}

// The standard types for the cases
var Schema schema

func init() {
	InitCaseTypes()
}

func InitScalarTypes() {
	if Schema.Scalar.Inited {
		return
	}
	Schema.Scalar.Lookup = TypeLookup()
	Schema.Scalar.Tuple = IfaceTuple()

	Schema.Scalar.Inited = true // required to not repeat initiation
}

func InitCaseTypes() {
	if Schema.Scalar.Inited {
		return
	}
	InitScalarTypes() // required

	Schema.Case.Link = TypeLink()

	Schema.Scalar.Inited = true // required to not repeat initiation
}
