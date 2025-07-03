package utils

import (
	_go "github.com/webitel/cases/api/cases"
	"reflect"
	"time"
)

type Lookup interface {
	SetId(int)
	GetId() *int
	SetName(string)
	GetName() *string
}

type ExtendedLookup interface {
	Lookup
	SetType(typ string)
	GetType() *string
}

func UnmarshalLookup[K Lookup](lp *_go.Lookup, lookup K) K {
	if lp == nil {
		var res K
		return res
	}
	if lp.Id != 0 {
		lookup.SetId(int(lp.Id))
	}
	if lp.Name != "" {
		lookup.SetName(lp.Name)

	}
	return lookup
}

func MarshalLookup(lp Lookup) *_go.Lookup {
	if lp == nil {
		return nil
	}
	val := reflect.ValueOf(lp)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil
	}
	var res _go.Lookup
	if id := lp.GetId(); id != nil {
		res.Id = int64(*id)
	}
	if name := lp.GetName(); name != nil {
		res.Name = *name
	}

	return &res
}

func UnmarshalExtendedLookup[K ExtendedLookup](lp *_go.ExtendedLookup, lookup K) K {
	if lp == nil {
		var res K
		return res
	}
	if lp.Id != 0 {
		lookup.SetId(int(lp.Id))
	}
	if lp.Name != "" {
		lookup.SetName(lp.Name)
	}
	if lp.Type != "" {
		lookup.SetType(lp.Type)
	}
	return lookup
}

func MarshalExtendedLookup(lp ExtendedLookup) *_go.ExtendedLookup {
	if lp == nil {
		return nil
	}
	val := reflect.ValueOf(lp)
	if val.Kind() == reflect.Ptr && val.IsNil() {
		return nil
	}
	var res _go.ExtendedLookup
	if id := lp.GetId(); id != nil {
		res.Id = int64(*id)
	}
	if name := lp.GetName(); name != nil {
		res.Name = *name
	}
	if typ := lp.GetType(); typ != nil {
		res.Type = *typ
	}

	return &res
}

func Dereference[T any](lp *T) T {
	if lp == nil {
		return *new(T)
	}
	return *lp
}

func MarshalTime(t *time.Time) int64 {
	if t == nil || t.IsZero() {
		return 0
	}
	return t.UnixMilli()
}

func TimePtr(ms int64) *time.Time {
	if ms == 0 {
		return nil
	}
	t := time.UnixMilli(ms)
	return &t
}
