package model

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type NullTime time.Time

// - database/sql.Valuer
func (t NullTime) Value() (driver.Value, error) {
	if !t.IsZero() {
		return t.Time(), nil
	}
	return nil, nil
}

// IsZero value (?)
func (t *NullTime) IsZero() bool {
	if t != nil && !t.Time().IsZero() {
		return false
	}
	return true
}

func (t *NullTime) ToMilliseconds() int64 {
	if t == nil {
		return 0
	}

	return t.Time().UnixMilli()
}

func NewNullTime(i time.Time) *NullTime {
	return (*NullTime)(&i)
}

// IsZero value (?)
func (t *NullTime) String() string {
	if !t.IsZero() {
		return t.Time().String()
	}
	return ""
}

// Time native value (!)
func (t *NullTime) Time() (stamp time.Time) {
	// v == time.Time{} // Zero(!)
	if t != nil {
		stamp = (time.Time)(*t) // shallowcopy(!)
	}
	return stamp
}

// Scan native value decoder function
func (t *NullTime) Scan(v interface{}) error {
	// scan: nullable (!)
	if v == nil {
		(t) = nil // Zero(!)
		return nil
	}
	switch v := v.(type) {
	case time.Time:
		// +OK: datetime
		(*t) = (NullTime)(v) // shallowcopy
		return nil
	default:
	}
	return fmt.Errorf("nulltime: convert %[1]T value %[1]v", v)
}
