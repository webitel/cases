package graph

var CaseLink Link

// Link graph model
type Link struct {
	// Object that describes link's object, with it's default and all fields
	Output Metadata
	// implements
	*Tuple
	// fields
	Name   *Metadata
	Url    *Metadata
	Author *Lookup
}

func init() {
	CaseLink = Link{
		Output: Metadata{
			Name: "link",
			Type: "object",
			Fields: []*Metadata{
				// name
				{Name: "name", Type: "string"},
				{Name: "url", Type: "string"},
				{Name: "author", Type: "lookup"},
			},
		},
	}
}
