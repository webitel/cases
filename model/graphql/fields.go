package graphql

// Argument Input of the Field
type Argument struct {
	Name    string      `json:"name"`
	Type    Input       `json:"type"`
	About   string      `json:"description"`
	Default interface{} `json:"defaultValue"`
}

// GraphQL *like Field Descriptor
type Field struct {
	// Name of the Field
	Name string
	// OPTIONAL. Formal Arguments
	Args []*Argument
	// Type of the Field
	Type Output
	// OPTIONAL. Nested Fields of the Object Type.
	Fields
}

// GraphQL *like set of the Object's Fields
type Fields map[string]*Field

// Field[Q]uery Arguments
type FieldQ struct {
	// Field Descriptor
	*Field
	// Optional. Actual Arguments
	Args map[string]any
	// Nested Fields query
	Fields []*FieldQ
}

// InputFields parses ?fields=a.arg(val),b{c,d{e}} query string
func InputFields(schema Fields, query string) ([]*FieldQ, error) {
	panic("not implemented")
}
