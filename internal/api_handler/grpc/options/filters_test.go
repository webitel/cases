package options

import (
	"testing"

	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func Test_extractIdentifier(t *testing.T) {
	type args struct {
		expr *expr.Expr
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{

		// "id"
		{
			name: "Singe identifier",
			args: args{
				expr: &expr.Expr{
					ExprKind: &expr.Expr_IdentExpr{
						IdentExpr: &expr.Expr_Ident{Name: "id"},
					},
				},
			},
			want:    "id",
			wantErr: false,
		},
		// "case.id"
		{
			name: "Nested identifier",
			args: args{
				expr: &expr.Expr{
					ExprKind: &expr.Expr_SelectExpr{
						SelectExpr: &expr.Expr_Select{
							Operand: &expr.Expr{
								ExprKind: &expr.Expr_IdentExpr{
									IdentExpr: &expr.Expr_Ident{Name: "case"},
								},
							},
							Field: "id",
						},
					},
				},
			},
			want:    "id",
			wantErr: false,
		},
		// "case.status_condition.id"
		{
			name: "Triple nested identifier",
			args: args{
				expr: &expr.Expr{
					ExprKind: &expr.Expr_SelectExpr{
						SelectExpr: &expr.Expr_Select{
							Operand: &expr.Expr{
								ExprKind: &expr.Expr_SelectExpr{
									SelectExpr: &expr.Expr_Select{
										Operand: &expr.Expr{
											ExprKind: &expr.Expr_IdentExpr{
												IdentExpr: &expr.Expr_Ident{Name: "case"},
											},
										},
										Field: "status_condition",
									},
								},
							},
							Field: "id",
						},
					},
				},
			},
			want:    "status_condition.id",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractIdentifier(tt.args.expr)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractIdentifier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractIdentifier() got = %v, want %v", got, tt.want)
			}
		})
	}
}
