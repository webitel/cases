package utils

import (
	_go "github.com/webitel/cases/api/cases"
	"github.com/webitel/cases/internal/model"
	"reflect"
	"testing"
	"time"
)

func GetIntPointer(value int) *int {
	if value < 0 {
		return nil
	}
	return &value
}

func GetStringPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func TestDereference(t *testing.T) {
	type args[T any] struct {
		lp *T
	}
	type testCase[T any] struct {
		name string
		args args[T]
		want T
	}
	timeTests := []testCase[time.Time]{
		{
			name: "not nil",
			args: args[time.Time]{
				lp: &time.Time{},
			},
			want: time.Time{},
		},
		{
			name: "nil",
			args: args[time.Time]{
				lp: nil,
			},
			want: time.Time{},
		},
	}
	intTests := []testCase[int]{
		{
			name: "not nil",
			args: args[int]{
				lp: GetIntPointer(1),
			},
			want: 1,
		},
		{
			name: "nil",
			args: args[int]{
				lp: GetIntPointer(-1),
			},
			want: 0,
		},
	}
	stringTests := []testCase[string]{
		{
			name: "not nil",
			args: args[string]{
				lp: GetStringPointer("sup"),
			},
			want: "sup",
		},
		{
			name: "nil",
			args: args[string]{
				lp: GetStringPointer(""),
			},
			want: "",
		},
	}
	for _, tt := range timeTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dereference(tt.args.lp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dereference() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dereference(tt.args.lp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dereference() = %v, want %v", got, tt.want)
			}
		})
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Dereference(tt.args.lp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dereference() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMarshalLookup(t *testing.T) {
	type args struct {
		lp model.Lookup
	}
	tests := []struct {
		name string
		args args
		want *_go.Lookup
	}{
		{
			name: "nil",
			args: args{
				lp: nil,
			},
			want: nil,
		},
		{
			name: "empty lookup",
			args: args{
				lp: &model.Author{
					Id:   nil,
					Name: nil,
				},
			},
			want: &_go.Lookup{
				Id:   0,
				Name: "",
			},
		},
		{
			name: "non-empty lookup",
			args: args{
				lp: &model.Author{
					Id:   GetIntPointer(1),
					Name: GetStringPointer("sup"),
				},
			},
			want: &_go.Lookup{
				Id:   1,
				Name: "sup",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MarshalLookup(tt.args.lp); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MarshalLookup() = %v, want %v", got, tt.want)
			}
		})
	}
}

//func TestUnmarshalLookup(t *testing.T) {
//	type args[K model.Lookup] struct {
//		lp     *_go.Lookup
//		lookup K
//	}
//	type testCase[K model.Lookup] struct {
//		name string
//		args args[K]
//		want K
//	}
//	tests := []testCase[ /* TODO: Insert concrete types here */ ]{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := UnmarshalLookup(tt.args.lp, tt.args.lookup); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("UnmarshalLookup() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
