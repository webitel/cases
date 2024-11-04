package graph

// Tuple Interface
type Tuple struct {
	// Output Type Metadata
	Output *Metadata
	// fields
	Id        *Metadata
	Ver       *Metadata // revision: version
	Etag      *Metadata
	CreatedAt *Metadata
	CreatedBy *Lookup // metadata
	UpdatedAt *Metadata
	UpdatedBy *Lookup // metadata
}

func IfaceTuple() Tuple {
	output := Tuple{
		Id: &Metadata{
			Name: "id",
			Type: "string!",
		},
		Ver: &Metadata{
			Name: "ver",
			Type: "int32",
		},
		Etag: &Metadata{
			Name: "etag",
			Type: "string!",
		},
		CreatedAt: &Metadata{
			Name: "created_at",
			Type: "int64!",
		},
		CreatedBy: Schema.Scalar.Lookup.Field(
			"created_by", "user",
		),
		UpdatedAt: &Metadata{
			Name: "updated_at",
			Type: "int64",
		},
		UpdatedBy: Schema.Scalar.Lookup.Field(
			"updated_by", "user",
		),
	}
	output.Output = &Metadata{
		Name: "tuple",
		Type: "interface",
		Fields: []*Metadata{
			output.Id,
			output.Ver,
			output.Etag,
			output.CreatedAt,
			&output.CreatedBy.Output,
			output.UpdatedAt,
			&output.UpdatedBy.Output,
		},
		Default: []string{
			output.Id.Name,
		},
	}
	return output
}
