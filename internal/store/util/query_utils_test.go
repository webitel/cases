package util

import (
	"github.com/Masterminds/squirrel"
	"reflect"
	"testing"
)

func TestFormAsCTEs(t *testing.T) {
	type args struct {
		in map[string]squirrel.Sqlizer
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   []any
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				in: map[string]squirrel.Sqlizer{
					"my_new_1": squirrel.Select("1", "2", "3", "4", "5", "6").From("table1").Where("id = ?", 1),
					"my_new_2": squirrel.Select("1", "2", "3", "4", "5", "6").From("table2").Where("ver = ?", 2),
					"my_new_3": squirrel.Select("1", "2", "3", "4", "5", "6").From("table3").Where("name = ?", "yehor"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub, a, _ := FormAsCTEs(tt.args.in)
			query, ar, _ := squirrel.Select("one").From("good").Where("name = ?", "yehor").Prefix(sub, a...).PlaceholderFormat(squirrel.Dollar).ToSql()

			t.Log(query, ar)
		})
	}
}

func TestPrepareSearchNumber(t *testing.T) {
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			"1",
			"1234567890",
			"0987654321",
		},
		{
			"2",
			"hello",
			"olleh",
		},
		{
			"3",
			"world",
			"dlrow",
		},
		{
			"4",
			"golang",
			"gnalog",
		},
		{
			"5",
			"backend",
			"dnekcab",
		},
		{
			"6",
			"developer",
			"repoleved",
		},
		{
			"7",
			"application",
			"noitacilppa",
		},
		{
			"8",
			"programming",
			"gnimmargorp",
		},
		{
			"9",
			"example",
			"elpmaxe",
		},
		{
			"10",
			"reverse",
			"esrever",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PrepareSearchNumber(tt.args); got != tt.want {
				t.Errorf("PrepareSearchNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolvePaging(t *testing.T) {
	tests := []struct {
		name      string
		size      int
		items     []int
		wantItems []int
		wantNext  bool
	}{
		{
			name: "No paging, all items returned",
			size: -1,
			items: []int{
				1, 2, 3, 4, 5,
			},
			wantItems: []int{
				1, 2, 3, 4, 5,
			},
			wantNext: false,
		},
		{
			name: "Paging with size greater than items",
			size: 10,
			items: []int{
				1, 2, 3,
			},
			wantItems: []int{
				1, 2, 3,
			},
			wantNext: false,
		},
		{
			name: "Paging with size equal to items",
			size: 4,
			items: []int{
				1, 2, 3, 4,
			},
			wantItems: []int{
				1, 2, 3, 4,
			},
			wantNext: false,
		},
		{
			name: "Paging with size less than items",
			size: 4,
			items: []int{
				1, 2, 3, 4, 5, 6,
			},
			wantItems: []int{
				1, 2, 3, 4,
			},
			wantNext: true,
		},
		{
			name: "Paging with size 1",
			size: 1,
			items: []int{
				1, 2, 3, 4, 5, 6,
			},
			wantItems: []int{
				1,
			},
			wantNext: true,
		},
		{
			name:      "Empty items",
			size:      3,
			items:     []int{},
			wantItems: []int{},
			wantNext:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items := make([]*int, len(tt.items))
			for i := range tt.items {
				items[i] = &tt.items[i]
			}
			gotItems, gotNext := ResolvePaging(tt.size, items)
			t.Logf("items slice: %p, array: %p\n", &items, items)
			t.Logf("gotItems slice: %p, array: %p\n", &gotItems, gotItems)
			if &items == &gotItems {
				t.Errorf("ResolvePaging() results mutate input parameter, input - %p , output - %p", &items, &gotItems)
			}
			// Convert `gotItems` back to values for comparison.
			gotValues := make([]int, len(gotItems))
			for i := range gotItems {
				gotValues[i] = *gotItems[i]
			}

			if !reflect.DeepEqual(gotValues, tt.wantItems) {
				t.Errorf("ResolvePaging() gotItems = %v, want %v", gotValues, tt.wantItems)
			}

			if !reflect.DeepEqual(gotValues, tt.wantItems) {
				t.Errorf("ResolvePaging() gotItems = %v, want %v", gotValues, tt.wantItems)
			}
			if gotNext != tt.wantNext {
				t.Errorf("ResolvePaging() gotNext = %v, want %v", gotNext, tt.wantNext)
			}
		})
	}
}
