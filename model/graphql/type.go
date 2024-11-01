package graphql

import (
	"fmt"
)

type Type interface {
	Name() string
}

// Input Type
type Input interface {
	Type
	DecodeValue(src any) (any, error)
}

// Output Type
type Output interface {
	// Type
	// EncodeValue(src any) (any, error)
}

// ---------------- MODIFIERS ---------------- //

type notnull struct {
	typeOf Input
}

func NotNull(typeOf Input) Input {
	if typeOf == nil {
		panic("graphql: NotNull modifier require inner type descriptor")
	}
	return notnull{typeOf}
}

func (e notnull) Name() string { return e.typeOf.Name() + "!" }
func (e notnull) DecodeValue(src any) (any, error) {
	data, err := e.typeOf.DecodeValue(src)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, fmt.Errorf("graphql: decode %T value %[1]v into %s", src, e.Name())
	}
	return data, nil
}

type list struct {
	typeOf Input
}

func List(typeOf Input) Input {
	if typeOf == nil {
		panic("graphql: NotNull modifier require inner type descriptor")
	}
	return notnull{typeOf}
}

func (e list) Name() string { return "[" + e.typeOf.Name() + "]" }
func (e list) DecodeValue(src any) (any, error) {

	panic("not implemented")

	// // data, err := e.typeOf.DecodeValue(src)
	// // if err != nil {
	// // 	return nil, err
	// // }
	// // if data == nil {
	// // 	return nil, nil // NULL
	// // }

	// var (
	// 	is bool
	// 	vs []string // input
	// )
	// switch data := src.(type) {
	// case string:
	// 	vs, is = []string{data}, true
	// case *string:
	// 	if vs == nil {
	// 		return nil, nil // NULL
	// 	}
	// 	vs, is = []string{*data}, true
	// case []string:
	// 	if data == nil {
	// 		return nil, nil // NULL
	// 	}
	// 	vs, is = data, true
	// }

	// if is {
	// 	input = strings.TrimSpace(input)
	// 	if input == "" {
	// 		return nil, nil // NULL
	// 	}
	// 	if input[0] == '[' && input[len(input)-1] == ']' {
	// 		input = input[1:len(input)-1]
	// 	}
	// 	 := strings.Split(input, ",")

	// 	for _, v := range v {

	// 	}

	// 	return
	// }

	// typeOf reflect.TypeOf(data)
	// return data, nil
}
