package color

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewColor(t *testing.T) {
	type args struct {
		r float64
		g float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want Color
	}{
		{
			name: "test1",
			args: args{
				r: 0,
				g: 0,
				b: 0,
			},
			want: Color{0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewColor(tt.args.r, tt.args.g, tt.args.b), tt.want, "should be equal")
		})
	}
}

func TestColor_Add(t *testing.T) {

	tests := []struct {
		name string
		c    Color
		c2   Color
		want Color
	}{
		{
			name: "add1",
			c:    NewColor(0.9, 0.6, 0.75),
			c2:   NewColor(0.7, 0.1, 0.25),
			want: NewColor(1.6, 0.7, 1.0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.c.Add(tt.c2), tt.want, "should be equal")
		})
	}
}

func TestColor_Sub(t *testing.T) {

	tests := []struct {
		name string
		c    Color
		c2   Color
		want Color
	}{
		{
			name: "sub",
			c:    NewColor(0.9, 0.6, 0.75),
			c2:   NewColor(0.7, 0.1, 0.25),
			want: NewColor(0.2, 0.5, 0.5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, Equal(tt.c.Sub(tt.c2), tt.want), "should be equal")
		})
	}
}

func TestColor_Scale(t *testing.T) {

	tests := []struct {
		name string
		c    Color
		s    float64
		want Color
	}{
		{
			name: "sub",
			c:    NewColor(0.2, 0.3, 0.4),
			s:    2.0,
			want: NewColor(0.4, 0.6, 0.8),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, Equal(tt.c.Scale(tt.s), tt.want), "should be equal")
		})
	}
}

func TestColor_Mult(t *testing.T) {

	tests := []struct {
		name string
		c    Color
		c2   Color
		want Color
	}{
		{
			name: "sub",
			c:    NewColor(1, 0.2, 0.4),
			c2:   NewColor(0.9, 1, 0.1),
			want: NewColor(0.9, 0.2, 0.04),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, Equal(tt.c.Mult(tt.c2), tt.want), "should be equal")
		})
	}
}
