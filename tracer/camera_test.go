package tracer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewTransform(t *testing.T) {
	type args struct {
		from Vector
		to   Vector
		up   Vector
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "default",
			args: args{
				from: NewVector(0, 0, 0),
				to:   NewVector(0, 0, -1),
				up:   NewVector(0, 1, 0),
			},
			want: IdentityMatrix(),
		},
		{
			name: "positive z",
			args: args{
				from: NewVector(0, 0, 0),
				to:   NewVector(0, 0, 1),
				up:   NewVector(0, 1, 0),
			},
			want: IdentityMatrix().Scale(-1, 1, -1),
		},
		{
			name: "view transform moves the world",
			args: args{
				from: NewVector(0, 0, 8),
				to:   NewVector(0, 0, 0),
				up:   NewVector(0, 1, 0),
			},
			want: IdentityMatrix().Translate(0, 0, -8),
		},
		{
			name: "arbitrary view transformation",
			args: args{
				from: NewVector(1, 3, 2),
				to:   NewVector(4, -2, 8),
				up:   NewVector(1, 1, 0),
			},
			want: NewMatrixFromData([][]float64{
				{-0.50709, 0.50709, 0.67612, -2.36643},
				{0.76772, 0.60609, 0.12122, -2.82843},
				{-0.35857, 0.59761, -0.71714, 0},
				{0, 0, 0, 1.0},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.want.Equals(ViewTransform(tt.args.from, tt.args.to, tt.args.up)))
		})
	}
}

func TestNewCamera(t *testing.T) {
	type args struct {
		hsize int
		vsize int
		fov   float64
	}
	tests := []struct {
		name string
		args args
		want *Camera
	}{
		{
			name: "test1",
			args: args{
				hsize: 160,
				vsize: 120,
				fov:   math.Pi / 2,
			},
			want: &Camera{
				Hsize:     160,
				Vsize:     120,
				FoV:       math.Pi / 2,
				Transform: IdentityMatrix(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewCamera(tt.args.hsize, tt.args.vsize, tt.args.fov))
		})
	}
}
