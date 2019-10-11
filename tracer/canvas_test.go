package tracer

import (
	"image/color"
	"testing"

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
		{
			name: "canvas-rect",
			args: args{
				w: 10,
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
					assert.Equal(t, c.data[w][h], color.RGBA{0, 0, 0, 0xff})
				}
			}
		})
	}
}

func TestCanvas_Set(t *testing.T) {
	type args struct {
		x   int
		y   int
		clr color.Color
	}
	tests := []struct {
		name    string
		canvas  *Canvas
		args    args
		wantErr bool
	}{
		{
			name:   "set1",
			canvas: NewCanvas(40, 20),
			args: args{
				x:   10,
				y:   15,
				clr: color.RGBA{10, 10, 10, 0xff},
			},
			wantErr: false,
		},
		{
			name:   "invalid1",
			canvas: NewCanvas(40, 20),
			args: args{
				x:   70,
				y:   15,
				clr: color.RGBA{10, 10, 10, 0xff},
			},
			wantErr: true,
		},
		{
			name:   "invalid2",
			canvas: NewCanvas(40, 20),
			args: args{
				x:   15,
				y:   20,
				clr: color.RGBA{10, 10, 10, 0xff},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.wantErr {
			case true:
				assert.Error(t, tt.canvas.Set(tt.args.x, tt.args.y, tt.args.clr), "no error")
			case false:
				assert.NoError(t, tt.canvas.Set(tt.args.x, tt.args.y, tt.args.clr), "no error")
				assert.Equal(t, tt.args.clr, tt.canvas.data[tt.args.x][tt.args.y], "should be equal")

			}
		})
	}
}

func TestCanvas_SetFloat(t *testing.T) {
	type args struct {
		x   float64
		y   float64
		clr color.Color
	}
	tests := []struct {
		name    string
		canvas  *Canvas
		args    args
		wantErr bool
	}{
		{
			name:   "set1",
			canvas: NewCanvas(40, 20),
			args: args{
				x:   10.0,
				y:   15.0,
				clr: color.RGBA{10, 10, 10, 0xff},
			},
			wantErr: false,
		},
		{
			name:   "invalid1",
			canvas: NewCanvas(40, 20),
			args: args{
				x:   70.00001,
				y:   15.00001,
				clr: color.RGBA{10, 10, 10, 0xff},
			},
			wantErr: true,
		},
		{
			name:   "invalid2",
			canvas: NewCanvas(40, 20),
			args: args{
				x:   15.0,
				y:   20.0,
				clr: color.RGBA{10, 10, 10, 0xff},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.wantErr {
			case true:
				assert.Error(t, tt.canvas.SetFloat(tt.args.x, tt.args.y, tt.args.clr), "no error")
			case false:
				assert.NoError(t, tt.canvas.SetFloat(tt.args.x, tt.args.y, tt.args.clr), "no error")
				assert.Equal(t, tt.args.clr, tt.canvas.data[int(tt.args.x)][int(tt.args.y)], "should be equal")

			}
		})
	}
}
func TestCanvas_Get(t *testing.T) {
	type args struct {
		x int
		y int
	}
	tests := []struct {
		name    string
		args    args
		canvas  *Canvas
		want    color.Color
		wantErr bool
	}{
		{
			name: "valid1",
			args: args{
				x: 10,
				y: 20,
			},
			canvas:  NewCanvas(30, 50),
			want:    color.RGBA{10, 10, 10, 0xff},
			wantErr: false,
		},
		{
			name: "invalid1",
			args: args{
				x: 80,
				y: 20,
			},
			canvas:  NewCanvas(30, 50),
			want:    color.RGBA{10, 10, 10, 0xff},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got color.Color
			var err error

			switch tt.wantErr {
			case true:
				got, err = tt.canvas.Get(tt.args.x, tt.args.y)
				assert.Error(t, err, "should error")
			case false:
				tt.canvas.Set(tt.args.x, tt.args.y, tt.want)
				got, err = tt.canvas.Get(tt.args.x, tt.args.y)
				assert.NoError(t, err, "should not error")
				assert.Equal(t, tt.want, got, "should equal")
			}
		})
	}
}
