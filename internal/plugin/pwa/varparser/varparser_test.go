package varparser

import (
	"reflect"
	"testing"
)

func TestVarParser_Parse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want []Var
	}{
		{
			name: "",
			args: args{
				s: "vsdvsdvsdv sdvsdvsdv {test} dvsdvsdv sdvsdv {title}",
			},
			want: []Var{
				{ID: "test", Pos: VarPos{21, 27}},
				{ID: "title", Pos: VarPos{44, 51}},
			},
		},
		{
			name: "",
			args: args{
				s: "{a} hello {b} world {c}",
			},
			want: []Var{
				{ID: "a", Pos: VarPos{0, 3}},
				{ID: "b", Pos: VarPos{10, 13}},
				{ID: "c", Pos: VarPos{20, 23}},
			},
		},
		{
			name: "",
			args: args{
				s: "vddb {items[i].children[i]}",
			},
			want: []Var{
				{ID: "items[i].children[i]", Pos: VarPos{5, 27}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VarParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
