package graph

import (
	"context"
)

// ResolveArgs Context
type ResolveArgs struct {
	// Node to output
	Node any
	// Query of the Node
	*Query
	// Context bindings
	context.Context
}

// OutputFunc must resolve Node-* related data for output
type OutputFunc func(output *ResolveArgs) (data any, err error)

// func (md *Metadata) Output(ctx context.Context, data any) (any, error) {
// 	if md.Resolve == nil {
// 		return data, nil
// 	}
// 	return md.Resolve(&ResolveArgs{

// 	})
// }

// // // OutputMap fields for output
// // type OutputMap map[string]OutputFunc

// // output data field execution context
// type outputFd struct {
// 	context ResolveArgs
// 	resolve OutputFunc
// }

// // output data fields execution plan
// type OutputMap []outputFd

// func (plan *OutputMap) Resolve(query *Query, output OutputFunc) {
// 	*(plan) = append(*(plan), outputFd{
// 		resolve: output,
// 		context: ResolveArgs{
// 			Query: query,
// 		},
// 	})
// }

// // executes output fields plan for given node
// func (plan OutputMap) Execute(data interface{}) {
// 	// for _, out := range ctx {
// 	// 	out.context.node = node // populate
// 	// 	out.resolve(&out.context)
// 	// }
// 	var fd *outputFd
// 	for e, n := 0, len(plan); e < n; e++ {
// 		fd = &plan[e]
// 		// populate
// 		fd.context.Node = data
// 		// execute
// 		fd.resolve(&fd.context)
// 	}
// }

// func OutputPlan(output *Metadata, query *Query, data interface{}) OutputMap {
// 	if query == nil || len(query.Fields) == 0 { // no query
// 		return nil // return: as is
// 	}
// 	var (
// 		typeOf = data
// 		fields = query.Fields
// 		output OutputMap                  // []outputFd
// 		except []protoreflect.FieldNumber // .FieldDescriptor
// 	)
// 	// plan output fields manipulations
// 	for _, input := range fields {
// 		field := operation.GetField(input.Name) //  .fields[input.Name]
// 		if field == nil {
// 			panic(fmt.Errorf("output: %s{%s} no such field", query.Name, input.Name))
// 		}
// 		if field.Resolve != nil {
// 			output.Resolve(input, field.Resolve)
// 		}
// 		// descr := fields.ByName(protoreflect.Name(field.Name))
// 	}
// 	return output
// }
