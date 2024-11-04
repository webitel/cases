package graph

import (
	"github.com/webitel/cases/model/graph"
	"strings"
)

type Types struct {
	// interfaces
	Lookup Lookup
	Tuple  Tuple
	// type:output
	// any other types needed for APIs
	// case links, comments, related comments
}

type schema struct {
	// types
	Types Types
	// query
	GetContact  *Metadata
	ContactQin  *Metadata // ?qin=$fields; extra configuration of contact.fields that support ?q= term comparison
	ListContact *Metadata
	// mutation
}

var (
	Schema schema
)

// func Schema() schema

func init() {

	Schema.Types.Lookup = TypeLookup()
	Schema.Types.Tuple = IfaceTuple()

	// ---------- Schema.Fields ------------- //

	// GET /contacts
	types := &Schema.Types
	Schema.GetContact = Operation(
		"getContact", "contact", typeOf.Metadata,
		DefaultFields(
			typeOf.Id.Name,
			typeOf.Ver.Name,
			typeOf.Etag.Name,
			// Schema.Types.Contact.Name.Name,
			"name{"+strings.Join([]string{
				Schema.Types.Name.GivenName.Name,
				Schema.Types.Name.MiddleName.Name,
				Schema.Types.Name.FamilyName.Name,
				Schema.Types.Name.CommonName.Name,
			}, ",")+"}",
			typeOf.About.Name,
		),
	)

	// GET /contacts?q=&qin=
	Schema.ContactQin = &Metadata{
		Name: "qin", // [q]uery[in]fields
		Fields: []*graph.Metadata{
			{
				Name: typeOf.Name.Name, // "name"
				Fields: []*graph.Metadata{
					(&types.Name).GivenName,
					(&types.Name).MiddleName,
					(&types.Name).FamilyName,
					(&types.Name).CommonName,
				},
				// Default: []string{
				// 	(&types.Name).CommonName.Name,
				// },
			},
			{
				Name: typeOf.About.Name,
			},
			{
				Name: typeOf.Labels.Name, // "labels"
				Fields: []*graph.Metadata{
					(&types.Label).Label, // "tag"
				},
				// Default: []string{
				// 	(&types.Label).Label.Name,
				// },
			},
			{
				Name: typeOf.Emails.Name, // "emails"
				Fields: []*graph.Metadata{
					(&types.Email).Email, // "email"
					{
						Name: (&types.Email).Type.Output.Name, // "type"
						Fields: []*graph.Metadata{
							(&types.Email).Type.Name, // "name"
						},
					},
				},
				// Default: []string{
				// 	(&types.Email).Email.Name,
				// },
			},
			{
				Name: typeOf.Phones.Name, // "phones"
				Fields: []*graph.Metadata{
					(&types.Phone).Number, // "number"
					{
						Name: (&types.Phone).Type.Output.Name, // "type"
						Fields: []*graph.Metadata{
							(&types.Phone).Type.Name, // "name"
						},
					},
				},
				// Default: []string{
				// 	(&types.Timezone).Timezone.Name.Name,
				// },
			},
			{
				Name: typeOf.Groups.Name,
				Fields: []*graph.Metadata{
					(&types.Group).Group.Name,
				},
			},
			{
				Name: typeOf.Managers.Name, // "managers"
				Fields: []*graph.Metadata{
					(&types.Manager).User.Name, // "name"
				},
				// Default: []string{
				// 	(&types.Manager).User.Name.Name,
				// },
			},
			{
				Name: typeOf.Timezones.Name, // "timezones"
				Fields: []*graph.Metadata{
					(&types.Timezone).Timezone.Name, // "name"
				},
				// Default: []string{
				// 	(&types.Timezone).Timezone.Name.Name,
				// },
			},
			{
				Name: typeOf.Variables.Name, // "variables"
				Fields: []*graph.Metadata{
					(&types.Variable).Key,   // "key"
					(&types.Variable).Value, // "value"
				},
				// Default: []string{
				// 	(&types.Variable).Value.Name,
				// },
			},
			{
				Name: typeOf.IMClients.Name,
				Fields: []*graph.Metadata{
					&(&types.IMClient).ExternalUser.Output,
					&(&types.IMClient).App.Output,
				},
			},
		},
		Default: []string{
			typeOf.Name.Name, // "name"
		},
	}
	// GET /contacts
	Schema.ListContact = Operation(
		"listContacts", "[contact!]", typeOf.Metadata,
		InputArgs(graph.InputArgs{
			"page": {Name: "page", Type: graph.InputPage(1), Value: uint32(1)},
			"size": {Name: "size", Type: graph.InputSize{1, 64, 32}, Value: int32(32)},
			"sort": {Name: "sort", Type: graph.InputSort{"name{common_name}", "!id"}, Value: []string{"name{common_name}", "!id"}},
		}),
		DefaultFields(
			Schema.Types.Contact.Id.Name,
			Schema.Types.Contact.Name.Name, // "name{common_name}",
		),
	)

}
