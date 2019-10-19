package tracer

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
)

func TestNewPointLight(t *testing.T) {
	type args struct {
		p Point
		i color.Color
	}
	tests := []struct {
		name string
		args args
		want PointLight
	}{
		{
			name: "test1",
			args: args{
				p: NewPoint(0, 0, 0),
				i: colornames.White,
			},
			want: PointLight{
				NewPoint(0, 0, 0),
				colornames.White,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewPointLight(tt.args.p, tt.args.i))
		})
	}
}
