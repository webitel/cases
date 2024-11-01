package graphql

import (
	"fmt"
	"strconv"
)

// ---------- SCALAR ----------- //

type Int32 struct{}

func (Int32) Name() string { return "int32" }
func (Int32) DecodeValue(src any) (any, error) {
	if src == nil {
		return nil, nil // NULL
	}
	switch data := src.(type) {
	case int32:
		return data, nil
	case *int32:
		if data != nil {
			return data, nil
		}
		return nil, nil // NULL
	case string:
		if data == "" {
			return nil, nil // NULL
		}
		i32, err := strconv.ParseInt(data, 10, 32)
		if err != nil {
			return nil, err
		}
		return int32(i32), nil
	}
	return nil, fmt.Errorf("graphql: decode %T value %[1]v into Int32", src)
}

type Int64 struct{}

func (Int64) Name() string { return "int64" }
func (Int64) DecodeValue(src any) (any, error) {
	if src == nil {
		return nil, nil // NULL
	}
	switch data := src.(type) {
	case int64:
		return data, nil
	case *int64:
		if data != nil {
			return *data, nil
		}
		return nil, nil // NULL
	case string:
		if data == "" {
			return nil, nil // NULL
		}
		i64, err := strconv.ParseInt(data, 10, 64)
		if err != nil {
			return nil, err
		}
		return i64, nil
	}
	return nil, fmt.Errorf("graphql: decode %T value %[1]v into Int64", src)
}

type Uint32 struct{}

func (Uint32) Name() string { return "uint32" }
func (Uint32) DecodeValue(src any) (any, error) {
	if src == nil {
		return nil, nil // NULL
	}
	switch data := src.(type) {
	case uint32:
		return data, nil
	case *uint32:
		if data != nil {
			return data, nil
		}
		return nil, nil // NULL
	case string:
		if data == "" {
			return nil, nil // NULL
		}
		i32, err := strconv.ParseUint(data, 10, 32)
		if err != nil {
			return nil, err
		}
		return uint32(i32), nil
	}
	return nil, fmt.Errorf("graphql: decode %T value %[1]v into Uint32", src)
}

type String struct{}

func (String) Name() string { return "string" }
func (String) DecodeValue(src any) (any, error) {
	if src == nil {
		return nil, nil // NULL
	}
	switch data := src.(type) {
	case string:
		return data, nil
	case *string:
		if data != nil {
			return data, nil
		}
		return nil, nil // NULL
	}
	return nil, fmt.Errorf("graphql: decode %T value %[1]v into String", src)
}
