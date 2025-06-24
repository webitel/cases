package utils

import (
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/model"
	"reflect"
	"time"
)

func UnmarshalLookup[K model.Lookup](lp *_go.Lookup, lookup K) K {
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

func MarshalLookup(lp model.Lookup) *_go.Lookup {
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