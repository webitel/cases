package graph

// Link graph model
type Link struct {
	// Object that describes link's object, with it's default and all fields
	Output *Metadata
	// implements
	*Tuple
	// fields
	Name   *Metadata
	Url    *Metadata
	Author *Lookup
}

func TypeLink() Link {
	link := Link{
		Tuple: &Schema.Scalar.Tuple,
		Output: &Metadata{
			Name: "link",
			Type: "object",
		},
		Name: &Metadata{
			Name: "name",
			Type: "string",
		},
		Url: &Metadata{
			Name: "url",
			Type: "string!",
		},
		Author: Schema.Scalar.Lookup.Field("author", "author"),
	}

	link.Output.Fields = append(link.Tuple.Output.Fields, link.Name, link.Url, &link.Author.Output)
	link.Output.Default = []string{link.Name.Name, link.Url.Name, link.Author.Output.Name}

	return link
}
