package graph

// Comment model for GraphQL
type Comment struct {
	// Object that describes comment's object, with its default and all fields
	Output *Metadata
	// Implements Tuple
	*Tuple
	// Fields
	ID        *Metadata
	Version   *Metadata
	ETag      *Metadata
	CreatedBy *Lookup
	CreatedAt *Metadata
	UpdatedBy *Lookup
	UpdatedAt *Metadata
	Author    *Lookup
	Text      *Metadata
	Edited    *Metadata
}

func TypeComment() Comment {
	comment := Comment{
		// Implementing Tuple
		Tuple: &Schema.Scalar.Tuple,

		// Define Output metadata
		Output: &Metadata{
			Name: "comment",
			Type: "object",
		},

		// Define each field in Comment
		ID: &Metadata{
			Name: "id",
			Type: "int64",
		},
		Version: &Metadata{
			Name: "ver",
			Type: "int32",
		},
		ETag: &Metadata{
			Name: "etag",
			Type: "string!",
		},
		CreatedBy: Schema.Scalar.Lookup.Field("created_by", "user"),
		CreatedAt: &Metadata{
			Name: "created_at",
			Type: "int64",
		},
		UpdatedBy: Schema.Scalar.Lookup.Field("updated_by", "user"),
		UpdatedAt: &Metadata{
			Name: "updated_at",
			Type: "int64",
		},
		Author: Schema.Scalar.Lookup.Field("author", "contact-author"),
		Text: &Metadata{
			Name: "text",
			Type: "string",
		},
		Edited: &Metadata{
			Name: "edited",
			Type: "bool",
		},
	}

	// Set Fields and Default in Output metadata
	comment.Output.Fields = append(
		comment.Tuple.Output.Fields,
		comment.ID,
		comment.Version,
		comment.ETag,
		&comment.CreatedBy.Output,
		comment.CreatedAt,
		&comment.UpdatedBy.Output,
		comment.UpdatedAt,
		&comment.Author.Output,
		comment.Text,
		comment.Edited,
	)

	comment.Output.Default = []string{
		comment.ID.Name,
		comment.ETag.Name,
		comment.Author.Name.Name,
		comment.Text.Name,
		comment.Edited.Name,
	}

	return comment
}
