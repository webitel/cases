package store

import (
	"github.com/Masterminds/squirrel"
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
