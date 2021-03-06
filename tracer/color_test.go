package tracer

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
)

func TestColorName(t *testing.T) {
	type args struct {
		c color.Color
	}
	tests := []struct {
		name string
		args args
		want Color
	}{
		{
			name: "black",
			args: args{
				c: colornames.Black,
			},
			want: NewColor(0, 0, 0),
		},
		{
			name: "white",
			args: args{
				c: colornames.White,
			},
			want: NewColor(1, 1, 1),
		},
		{
			name: "red",
			args: args{
				c: colornames.Red,
			},
			want: NewColor(1, 0, 0),
		},
		{
			name: "green",
			args: args{
				c: colornames.Lime,
			},
			want: NewColor(0, 1, 0),
		},
		{
			name: "blue",
			args: args{
				c: colornames.Blue,
			},
			want: NewColor(0, 0, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, ColorName(tt.args.c))
		})
	}
}

func TestColor_Equal(t *testing.T) {
	type args struct {
		c2 Color
	}
	tests := []struct {
		name string
		c1   Color
		args args
		want bool
	}{
		{
			c1: White(),
			args: args{
				c2: White(),
			},
			want: true,
		},
		{
			c1: White(),
			args: args{
				c2: Black(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c1.Equal(tt.args.c2), "should equal")
		})
	}
}
