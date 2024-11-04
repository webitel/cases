package graph

import (
	"github.com/webitel/cases/model/graph"
)

type (
	// Metadata Type Descriptor
	Metadata = graph.Metadata
	// Option
	Option func(md *Metadata)
)

func Definition(options ...Option) *Metadata {
	descriptor := &Metadata{}
	for _, option := range options {
		option(descriptor)
	}
	return descriptor
}

func TypeOf(output Metadata) Option {
	return func(md *graph.Metadata) {
		*(md) = output
	}
}

func NameOf(name, typeOf string) Option {
	return func(md *graph.Metadata) {
		md.Name = name
		if typeOf != "" {
			md.Type = typeOf
		}
	}
}

func InputArgs(input graph.InputArgs) Option {
	return func(md *graph.Metadata) {
		md.Args = input
		// if len(md.Args) == 0 {
		// 	md.Args = input
		// 	return
		// }
		// for _, param := range input {
		// 	md.Args.Add(param)
		// }
	}
}

func OutputFields(output ...*graph.Metadata) Option {
	return func(md *graph.Metadata) {
		md.Fields = output
		// if len(md.Fields) == 0 {
		// 	md.Fields = output
		// 	return
		// }
		// for _, field := range output {
		// 	if md.GetField(field.Name) != nil {
		// 		panic(fmt.Errorf("graphql: descriptor(%s).output(%s); field duplicate", md.Name, field.Name))
		// 	}
		// 	md.Fields = append(md.Fields, field)
		// }
	}
}

func DefaultFields(output ...string) Option {
	return func(md *graph.Metadata) {
		md.Default = output
	}
}

func ResolveData(output graph.OutputFunc) Option {
	return func(md *graph.Metadata) {
		md.Resolve = output
	}
}

func Operation(name, output string, typeOf Metadata, options ...Option) *Metadata {
	return Definition(func(operation *Metadata) {
		TypeOf(typeOf)(operation)
		operation.Type = typeOf.Name
		operation.Name = name
		for _, option := range options {
			option(operation)
		}
	})
}
