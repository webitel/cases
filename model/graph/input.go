package graph

import (
	"errors"
	"fmt"
	"strconv"
)

// Input Type interface
type Input interface {
	// String type name
	GoString() string
	// Decode input data source value
	InputValue(src any) (any, error)
}

type NotNull struct {
	TypeOf Input
}

var _ Input = NotNull{}

func (c NotNull) GoString() string {
	return c.TypeOf.GoString() + "!"
}

func (c NotNull) InputValue(src any) (any, error) {
	data, err := c.TypeOf.InputValue(src)
	if err == nil && data == nil {
		err = fmt.Errorf("notnull: value required")
	}
	return data, err
}

type Bool struct{}

var _ Input = Bool{}

func (c Bool) GoString() string {
	return "bool"
}

func (c Bool) InputValue(src any) (set any, err error) {
	if src == nil {
		return // set, nil
	}
	defer func() {
		if err != nil {
			set = nil
		}
	}()
	var value bool
	parse := func(s string) {
		value, err = strconv.ParseBool(s) // bool
		// if err != nil {
		// 	err =
		// }
	}
	switch data := src.(type) {
	case string:
		if len(data) == 0 {
			return // NULL
		}
		parse(data)
	case []byte:
		if len(data) == 0 {
			return // NULL
		}
		parse(string(data))
	case rune:
		if data == 0 {
			return // NULL
		}
		parse(string(data))
	case bool:
		value = data
	}

	if err != nil {
		return // nil, err
	}

	set = value
	return // bool, nil
}

// Input Argument descriptor
type Argument struct {
	// Name of the Argument
	Name string `json:"name"`
	// Type of the Argument
	Type Input `json:"typeOf"`
	// OPTIONAL. Default Value
	Value any `json:"value,omitempty"`
}

// IsValid Argument descriptor for use
func (e *Argument) IsValid() bool {
	return e != nil && e.Name != "" && e.Type != nil
}

// Required -if- has no default value
func (e *Argument) Required() bool {
	return e != nil && e.Value == nil
}

// Optional -if- has default value
func (e *Argument) Optional() bool {
	return e != nil && e.Value != nil
}

// InputArgs maps argument name to it's input type
type InputArgs map[string]Argument

func (vs InputArgs) Len() int {
	return len(vs)
}

func (vs InputArgs) Has(name string) bool {
	_, ok := vs[name]
	return ok
}

func (vs InputArgs) Get(name string) *Argument {
	if input, ok := vs[name]; ok {
		return &input // shallowcopy
	}
	return nil
}

func (vs InputArgs) Args() []string {
	if n := len(vs); n > 0 {
		as := make([]string, 0, n)
		for name, _ := range vs {
			as = append(as, name)
		}
		return as
	}
	return nil
}

// [option]: control
// [isDefault]: means option.Value(nil) returns default value
// [isMandatory]: means this option MUST be always [re]assigned with it's default value
// func (list *OptionList) Add(option Option, isDefault, isMandatory bool) error {
func (vs InputArgs) Add(reg Argument) (err error) {
	if reg.Name == "" {
		panic(fmt.Errorf("input: register <empty> name"))
	}
	if reg.Type == nil {
		panic(fmt.Errorf("input: register <nil> type"))
	}
	if vs.Has(reg.Name) {
		panic(fmt.Errorf("input: argument (%s:) duplicate", reg.Name))
	}
	if defVal := reg.Value; defVal != nil {
		reg.Value, err = reg.Type.InputValue(defVal)
		if err == nil && reg.Value == nil {
			err = fmt.Errorf("input: argument (%s:%s); default: invalid", reg.Name, defVal)
		}
		if err != nil {
			return // err
		}
	}
	// register NEW
	vs[reg.Name] = reg
	return nil // OK
}

// Parse req.Query options
func (vs InputArgs) Parse(req *Query) error {
	// Specified any ?
	if len(vs) == 0 {
		// NoArgs: not supported !
		if n := len(req.Args); n > 0 {
			return errors.New(
				fmt.Sprintf("graphql: %s( args:[%d] ); has no arguments",
					req.Name, n),
			)
		}
		// OK:
		return nil
	}
	var (
		err    error
		params = req.Args
	)
	// input: validate
	for name, value := range params {
		input := vs.Get(name)
		if input == nil {
			return errors.New(
				fmt.Sprintf("graphql: %s( %s: ); no such argument",
					req.Name, name),
			)
		}
		value, err = input.Type.InputValue(value)
		if err != nil {
			return errors.New(
				fmt.Sprintf("graphql: %s( %s:%v ); argument: %s",
					req.Name, name, params[name], err),
			)
		}
		// reset normalized
		params[name] = value
	}
	// input: defaults
	for _, input := range vs {
		// required ?
		if input.Value == nil {
			// specified ?
			if !params.Has(input.Name) {
				return fmt.Errorf("graphql: %s(%s:) argument required", req.Name, input.Name)
			}
			// specified !
			continue // OK
		}
		// optional !
		// specified ?
		if !params.Has(input.Name) {
			// NO: provide default value !
			params.Set(input.Name, input.Value)
		}
	}
	req.Args = params
	return nil
}

// Args of the query
type Args map[string]any

func (vs Args) Has(param string) (ok bool) {
	if vs != nil && param != "" {
		_, ok = vs[param]
	}
	return ok
}

func (vs Args) Get(param string) any {
	value, _ := vs[param]
	return value
}

func (vs *Args) Set(param string, value any) error {
	if param == "" {
		return fmt.Errorf("graphql: argument (name:) required")
	}
	if *(vs) == nil {
		*vs = make(Args)
	}
	(*vs)[param] = value
	return nil
}

func (vs Args) ValueOf(param string, typeOf Input) (value any, err error) {
	if param == "" {
		return nil, fmt.Errorf("graphql: argument (name:) required")
	}
	value = vs.Get(param)
	if typeOf != nil {
		value, err = typeOf.InputValue(value)
	}
	return // value, err
}

func (vs Args) Size() int32 {
	value, _ := vs.ValueOf(
		"size", InputSize{-1, 64, 0},
	)
	size, _ := value.(int32)
	return size // default(32)
}

func (vs Args) Page() uint32 {
	value, _ := vs.ValueOf(
		"page", InputPage(0),
	)
	page, _ := value.(uint32)
	return page // default(0)
}

func (vs Args) Sort() []string {
	value, _ := vs.ValueOf(
		"sort", InputSort(nil),
	)
	sort, _ := value.([]string)
	return sort // default(nil)
}

func (vs Args) Clone() Args {
	if vs == nil {
		return nil
	}
	v2 := make(Args, len(vs))
	for param, value := range vs {
		v2[param] = value
	}
	return v2
}
