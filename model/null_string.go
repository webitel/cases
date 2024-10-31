package model

import (
	"database/sql/driver"
	"fmt"
)

type NullString string

// - database/sql.Valuer
func (t NullString) Value() (driver.Value, error) {
	if !t.IsZero() {
		return t.String(), nil
	}
	return nil, nil
}

// IsZero value (?)
func (t *NullString) IsZero() bool {
	return t == nil
}

//// IsZero value (?)
//func (t *NullString) String() string {
//	if t != nil {
//		return t.Time().String()
//	}
//	return ""
//}

// Time native value (!)
func (t *NullString) String() (v string) {
	// v == time.Time{} // Zero(!)
	if t != nil {
		v = (string)(*t)
	}
	return v
}

func NewNullString(i string) *NullString {
	return (*NullString)(&i)
}

// Scan native value decoder function
func (t *NullString) Scan(v interface{}) error {
	// scan: nullable (!)
	if v == nil {
		(t) = nil // Zero(!)
		return nil
	}
	switch v := v.(type) {
	case string:
		// +OK: datetime
		(*t) = (NullString)(v) // shallowcopy
		return nil
	default:
	}
	return fmt.Errorf("nullstring: convert %[1]T value %[1]v", v)
}
