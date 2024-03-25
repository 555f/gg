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
		v    *VarParser
		args args
		want []Var
	}{
		{
			name: "",
			v:    &VarParser{},
			args: args{
				s: "vsdvsdvsdv sdvsdvsdv {test} dvsdvsdv sdvsdv {title}",
			},
			want: []Var{
				{ID: "test", Pos: VarPos{21, 26}},
				{ID: "title", Pos: VarPos{53, 60}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &VarParser{}
			if got := v.Parse(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("VarParser.Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
