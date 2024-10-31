package model

import (
	"database/sql/driver"
	"fmt"
	"strconv"
)

type NullInt int64

// - database/sql.Valuer
func (t NullInt) Value() (driver.Value, error) {
	if !t.IsZero() {
		return t.Int64(), nil
	}
	return nil, nil
}

// IsZero value (?)
func (t *NullInt) IsZero() bool {
	return t == nil
}

// IsZero value (?)
func (t *NullInt) Int() int {
	if t != nil {
		return (int)(*t)
	}
	return 0
}

// IsZero value (?)
func (t *NullInt) Int64() int64 {
	if t != nil {
		return (int64)(*t)
	}
	return 0
}

// IsZero value (?)
func (t *NullInt) Int32() int32 {
	if t != nil {
		return (int32)(*t)
	}
	return 0
}

func NewNullInt(i any) (*NullInt, error) {
	switch data := i.(type) {
	case int64:
		return (*NullInt)(&data), nil
	case int:
		value := int64(data)
		return (*NullInt)(&value), nil
	case int32:
		value := int64(data)
		return (*NullInt)(&value), nil
	case string:
		var (
			value int64
			err   error
		)
		value, err = strconv.ParseInt(data, 10, 64)
		if err != nil {
			return nil, err
		}
		return (*NullInt)(&value), nil
	default:
		return nil, fmt.Errorf("null_int: unknown format")
	}
}

func (t *NullInt) Scan(v interface{}) error {
	// scan: nullable (!)
	if v == nil {
		t = nil // Zero(!)
		return nil
	}

	switch v := v.(type) {
	case int:
		// +OK: int
		*t = (NullInt)(v)
		return nil
	case int64:
		// +OK: int
		newVal := (int)(v)
		*t = (NullInt)(newVal) // shallowcopy
		return nil
	case int32:
		// +OK: int
		newVal := (int)(v)
		*t = (NullInt)(newVal) // shallowcopy
		return nil

	default:
	}
	return fmt.Errorf("int: convert %[1]T value %[1]v", v)
}
