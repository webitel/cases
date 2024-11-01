package graph

import (
	"reflect"
	"testing"
)

func TestParseFieldsQ(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantFields FieldsQ
		wantErr    bool
	}{
		// TODO: Add test cases.
		{
			name: "none",
			args: args{
				s: "",
			},
			wantFields: nil,
			wantErr:    false,
		},
		{
			name: "simple",
			args: args{
				s: "id",
			},
			wantFields: FieldsQ{
				{Name: "id"},
			},
			wantErr: false,
		},
		{
			name: "paging",
			args: args{
				// s: "photos.sort(primary).offset(0).limit(1)",
				s: "photos.sort(primary).page(0).size(1)",
			},
			wantFields: FieldsQ{
				{
					Name: "photos",
					Args: Args{
						"size": int32(1),
						"page": uint32(0),
						"sort": []string{
							"primary",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "filters",
			args: args{
				s: "comments.since(1678803531).id(4642-4647).created_by(3,7){text,version}",
			},
			wantFields: FieldsQ{
				{
					Name: "comments",
					Args: Args{
						"id": []string{
							"4642-4647",
						},
						"since": []string{
							"1678803531",
						},
						"created_by": []string{
							"3", "7",
						},
					},
					Fields: FieldsQ{
						{
							Name: "text",
						},
						{
							Name: "version",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "nested",
			args: args{
				s: "id,name,emails.limit(3){id,from,email},phones.limit(3).sort(!created_at){number,created_by.id(7){name}},deleted_at",
			},
			wantFields: FieldsQ{
				{Name: "id"},
				{Name: "name"},
				{
					Name: "emails",
					Args: Args{
						"limit": int32(3),
					},
					Fields: FieldsQ{
						{Name: "id"},
						{Name: "from"},
						{Name: "email"},
					},
				},
				{
					Name: "phones",
					Args: Args{
						"limit": 3,
						"sort":  []string{"!created_at"},
					},
					Fields: FieldsQ{
						{Name: "number"},
						{
							Name: "created_by",
							Args: Args{
								"id": []string{
									"7",
								},
							},
							Fields: FieldsQ{
								{Name: "name"},
							},
						},
					},
				},
				{Name: "deleted_at"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFields, err := ParseFields(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFieldsQ() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotFields, tt.wantFields) {
				t.Errorf("ParseFieldsQ() = %v, want %v", gotFields, tt.wantFields)
				// return
			}
			normFields := gotFields.String()
			// if normFields != tt.args.s {
			// 	t.Errorf("FieldsQ.String() = %v, want %v", normFields, tt.args.s)
			// 	return
			// } else {
			t.Logf("FieldsQ.String() = %v", normFields)
			// }
		})
	}
}

func TestSplitFieldsQ(t *testing.T) {
	type args struct {
		fields string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
		{
			name: "",
			args: args{
				fields: "id,,name{common_name,first_name},emails.limit(3).order(-updated_at){id,email,created_by{id,name}}",
			},
			want: []string{
				"id",
				"name{common_name,first_name}",
				"emails.limit(3).order(-updated_at){id,email,created_by{id,name}}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SplitFieldsQ(tt.args.fields); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FieldsExpansion() = %v, want %v", got, tt.want)
			}
		})
	}
}
