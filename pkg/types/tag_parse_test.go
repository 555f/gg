package types

import (
	"go/token"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		comments Comments
	}
	tests := []struct {
		name    string
		args    args
		want    Tags
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				comments: Comments{
					{
						Value:    `@http-name:"test,format=sdvjhuskdjvh,omitempty"`,
						Position: token.Position{},
					},
				},
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTags(tt.args.comments)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTags() got = %v, want %v", got, tt.want)
			}
		})
	}
}
