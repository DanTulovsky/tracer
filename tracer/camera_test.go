package tracer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestViewTransform(t *testing.T) {
	type args struct {
		from Point
		to   Point
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
				from: NewPoint(0, 0, 0),
				to:   NewPoint(0, 0, -1),
				up:   NewVector(0, 1, 0),
			},
			want: IdentityMatrix(),
		},
		{
			name: "positive z",
			args: args{
				from: NewPoint(0, 0, 0),
				to:   NewPoint(0, 0, 1),
				up:   NewVector(0, 1, 0),
			},
			want: IdentityMatrix().Scale(-1, 1, -1),
		},
		{
			name: "view transform moves the world",
			args: args{
				from: NewPoint(0, 0, 8),
				to:   NewPoint(0, 0, 0),
				up:   NewVector(0, 1, 0),
			},
			want: IdentityMatrix().Translate(0, 0, -8),
		},
		{
			name: "arbitrary view transformation",
			args: args{
				from: NewPoint(1, 3, 2),
				to:   NewPoint(4, -2, 8),
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
			got := ViewTransform(tt.args.from, tt.args.to, tt.args.up)
			assert.True(t, tt.want.Equals(got))
		})
	}
}

func TestNewCamera(t *testing.T) {
	type args struct {
		hsize float64
		vsize float64
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
				Hsize:            160,
				Vsize:            120,
				fov:              math.Pi / 2,
				Transform:        IdentityMatrix(),
				TransformInverse: IdentityMatrix().Inverse(),
				HalfWidth:        1,
				HalfHeight:       0.75,
				PixelSize:        0.0125,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewCamera(tt.args.hsize, tt.args.vsize, tt.args.fov))
		})
	}
}

func TestCamera_PixelSize(t *testing.T) {
	tests := []struct {
		name   string
		camera *Camera
		want   float64
	}{
		{
			name:   "horizontal",
			camera: NewCamera(200, 125, math.Pi/2),
			want:   0.01,
		},
		{
			name:   "vertical",
			camera: NewCamera(125, 200, math.Pi/2),
			want:   0.01,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.camera.PixelSize)
		})
	}
}

func TestCamera_RayForPixel(t *testing.T) {
	type args struct {
		x float64
		y float64
	}
	tests := []struct {
		name      string
		camera    *Camera
		transform Matrix
		args      args
		want      Ray
	}{
		{
			name:   "through center of canvas",
			camera: NewCamera(201, 101, math.Pi/2),
			args: args{
				x: 100.5,
				y: 50.5,
			},
			transform: IdentityMatrix(),
			want:      NewRay(NewPoint(0, 0, 0), NewVector(0, 0, -1)),
		},
		{
			name:   "through corner of canvas",
			camera: NewCamera(201, 101, math.Pi/2),
			args: args{
				x: 0.5,
				y: 0.5,
			},
			transform: IdentityMatrix(),
			want:      NewRay(NewPoint(0, 0, 0), NewVector(0.66519, 0.33259, -0.66851)),
		},
		{
			name:   "ray through transformed camera",
			camera: NewCamera(201, 101, math.Pi/2),
			args: args{
				x: 100.5,
				y: 50.5,
			},
			// transform: IdentityMatrix().RotateY(math.Pi/4).Translate(0, -2, 5),
			transform: IdentityMatrix().Translate(0, -2, 5).RotateY(math.Pi / 4),
			want:      NewRay(NewPoint(0, 2, -5), NewVector(math.Sqrt2/2, 0, -math.Sqrt2/2)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.camera.SetTransform(tt.transform)
			assert.True(t, tt.want.Equal(tt.camera.RayForPixel(tt.args.x, tt.args.y)), "should equal")
		})
	}
}
