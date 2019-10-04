package tracer

import (
	"testing"

	"image/color"

	"github.com/stretchr/testify/assert"
)

func TestNewCanvas(t *testing.T) {
	type args struct {
		w int
		h int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "canvas-square",
			args: args{
				w: 30,
				h: 30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCanvas(tt.args.w, tt.args.h)

			// assert all pixels are black
			for w := 0; w < c.Width; w++ {
				for h := 0; h < c.Height; h++ {
					assert.Equal(t, c.Data[w][h], color.RGBA{0, 0, 0, 0})
				}
			}
		})
	}
}
