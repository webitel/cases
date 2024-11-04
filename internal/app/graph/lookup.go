package graph

import "github.com/webitel/cases/model/graph"

// Lookup Object Type
type Lookup struct {
	// Output Type Metadata
	Output Metadata
	// fields
	Id   *Metadata
	Type *Metadata
	Name *Metadata
}

func TypeLookup() Lookup {
	output := Lookup{
		Id: &Metadata{
			Name: "id",
			Type: "string!",
		},
		Type: &Metadata{
			Name: "type",
			Type: "string",
		},
		Name: &Metadata{
			Name: "name",
			Type: "string!",
		},
	}
	output.Output = Metadata{
		Name: "lookup",
		Type: "object",
		Fields: []*Metadata{
			output.Id,
			output.Name,
			// LookupType.Type.Metadata(),
		},
		Default: []string{
			output.Id.Name,
			output.Name.Name,
		},
	}
	return output
}

// Field descriptor of Lookup Type.
func (typo Lookup) Field(name, typeOf string) *Lookup {
	// shallowcopy; value type
	typo.Output = graph.Metadata{
		Name:    name,
		Type:    typeOf,
		Fields:  typo.Output.Fields,
		Default: typo.Output.Default,
	}
	// // ($type: string = "typeOf")
	// if typeOf != "" && typo.Type != nil {
	// 	typo.Output.Args = graph.InputArgs{
	// 		typo.Type.Name: graph.Argument{
	// 			Name:  typo.Type.Name,
	// 			Type:  graph.String(typeOf),
	// 			Value: typeOf,
	// 		},
	// 	}
	// }
	if typeOf == "" {
		typeOf = "lookup"
	} else {
		typeOf = "lookup<" + typeOf + ">"
	}
	typo.Output.Type = typeOf
	return &typo
}
