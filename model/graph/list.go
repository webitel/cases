package graph

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"
)

type (
	// Input type of the list (size:int32)
	InputSize struct {
		Minimum, Maximum, Default int32
	}
	// Input type of the list (page:uint32)
	InputPage uint32
	// Input type of the list (sort:$[-!+]field[{nested,..}],..)
	InputSort []string
)

var _ Input = InputSize{}

func (c InputSize) GoString() string {
	return "int32"
}

// Value returns `set` value of $int32;
// Zero(0) - stands for default value;
// Negative(-1) - no limit; fetch all;
// Positive(+1) - regular value;
func (c InputSize) InputValue(src any) (set any, err error) {
	defer func() { // defaultOrNull
		if err != nil {
			return // failed
		}
		size, notNull := set.(int32)
		if !notNull { // ISNULL
			if c.Default == 0 {
				set = nil
				return // NULL
			}
			size = c.Default // DEFAULT
		}
		if size < -1 {
			size = -1
		}
		// if size < c.Minimum {
		// 	size = c.Minimum
		// }
		if c.Maximum < size {
			size = c.Maximum
		}
		if size == 0 {
			if notNull { // input: NOTNULL
				err = fmt.Errorf("expect non-zero int32")
				return
			} else { // input: ISNULL
				set = nil
				return
			}
		}
		set = size
	}()
	if src == nil {
		return // set, nil
	}
	parse := func(s string) (int32, error) {
		i, err := strconv.ParseInt(s, 10, 32) // int32
		if err != nil {
			// ERR: expect integer value
			return 0, fmt.Errorf("expect non-zero int32")
		}
		return int32(i), err
	}
	var size int32 // value
	switch data := src.(type) {
	case []byte:
		if len(data) == 0 {
			return // set, nil
		}
		size, err = parse(
			string(data),
		)
	case string:
		if len(data) == 0 {
			return // set, nil
		}
		size, err = parse(
			data,
		)
	case int32:
		size = data
	case int:
		if data < math.MinInt32 || math.MaxInt32 < data {
			err = fmt.Errorf("convert %T value %[1]v into int32", src)
		}
		size = int32(data)
	default:
		err = fmt.Errorf("convert %T value %[1]v into int32", src)
	}

	if err != nil {
		return // nil, err
	}

	set = size
	return // size, nil
}

var _ Input = InputPage(0)

func (c InputPage) GoString() string {
	return "uint32"
}

// Value always return uint32 value; Default: 1.
func (c InputPage) InputValue(src any) (set any, err error) {
	defer func() { // defaultOrNull
		if err != nil {
			return // failed
		}
		page, notNull := set.(uint32)
		if !notNull { // ISNULL
			if c == 0 {
				set = nil
				return // NULL
			}
			page = uint32(c) // DEFAULT
		}
		if page == 0 { // input
			if notNull { // NOTNULL
				err = fmt.Errorf("zero")
				return
			} else { // ISNULL
				set = nil
				return
			}
		}
		set = page
	}()
	if src == nil {
		return // set, nil
	}
	parse := func(s string) (uint32, error) {
		n, err := strconv.ParseUint(s, 10, 32) // uint32
		if err != nil {
			// ERR: expect positive integer value
			return 0, fmt.Errorf("!int")
		}
		return uint32(n), nil
	}

	var page uint32 // value
	switch data := src.(type) {
	case []byte:
		page, err = parse(
			string(data),
		)
	case string:
		page, err = parse(
			data,
		)
	case int32:
		if data < 0 {
			err = fmt.Errorf("negative")
			break // switch
		}
		page = uint32(data)
	case int:
		if data < 0 {
			err = fmt.Errorf("negative")
			break // switch
		}
		if data > math.MaxInt32 {
			err = fmt.Errorf("overflow")
			break // switch
		}
		page = uint32(data)
	case uint:
		page = uint32(data)
	case uint32:
		page = data
	default:
		err = fmt.Errorf("convert %T value %[1]v into int32", src)
	}

	if err == nil && page == 0 {
		err = fmt.Errorf("zero")
	}

	if err != nil {
		return // nil, err
	}

	set = page
	return // set, nil
}

var _ Input = InputSize{}

func (c InputSort) GoString() string {
	return "[string]"
}

func (c InputSort) InputValue(src any) (set any, err error) {
	defaultOrNil := func() {
		set = nil // NULL
		if len(c) > 0 {
			set = append([]string(nil), c...)
		}
	}
	if src == nil {
		// Has Default ?
		defaultOrNil()
		return // set, nil
	}
	parse := func(s string) []string {
		return strings.FieldsFunc(
			s, func(r rune) bool {
				return r == ',' || unicode.IsSpace(r)
			},
		)
	}
	var sort []string // value
	switch data := src.(type) {
	case []byte:
		if len(data) == 0 {
			defaultOrNil()
			return
		}
		sort = parse(
			string(data),
		)
	case string:
		if len(data) == 0 {
			defaultOrNil()
			return
		}
		sort = parse(
			data,
		)
	case []string:
		if len(data) == 0 {
			defaultOrNil()
			return
		}
		for _, inline := range data {
			sort = append(sort, parse(inline)...)
		}
	default:
		err = fmt.Errorf("input(sort): convert %T value %[1]v into []string", src)
		return sort, err
	}

	// [TODO]: Validate sort(field,..) spec
	// for _, spec := range sort {

	// }

	if len(sort) == 0 {
		defaultOrNil()
	}

	return sort, nil
}
