package graph

import "fmt"

// String Input Type
type String string

var _ Input = (*String)(nil)

func (String) GoString() string {
	return "string"
}

func (c String) InputValue(src any) (set any, err error) {
	if src == nil {
		if c != "" {
			// default
			return string(c), nil
		}
		return nil, nil
	}
	switch data := src.(type) {
	case string:
		return data, nil
	case *string:
		if data == nil {
			return nil, nil
		}
		return *(data), nil
	case String:
		if data != "" {
			// default
			return string(data), nil
		}
		return nil, nil
	case *String:
		if data != nil && *data != "" {
			// default
			return string(*data), nil
		}
		return nil, nil
	}
	err = fmt.Errorf("input(string): convert %T value %[1]v", src)
	return nil, err
}
