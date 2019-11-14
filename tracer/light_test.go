package tracer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
)

func TestNewPointLight(t *testing.T) {
	type args struct {
		p Point
		i Color
	}
	tests := []struct {
		name string
		args args
		want *PointLight
	}{
		{
			name: "test1",
			args: args{
				p: NewPoint(0, 0, 0),
				i: ColorName(colornames.White),
			},
			want: &PointLight{
				NewPoint(0, 0, 0),
				ColorName(colornames.White),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewPointLight(tt.args.p, tt.args.i))
		})
	}
}

func Test_lighting(t *testing.T) {
	type args struct {
		m        *Material
		p        Point
		l        *PointLight
		eye      Vector
		normal   Vector
		inShadow bool
	}
	tests := []struct {
		name string
		args args
		o    Shaper
		want Color
	}{
		{
			name: "eye between light and surface",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 0, -10), ColorName(colornames.White)),
				eye:    NewVector(0, 0, -1),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(1.9, 1.9, 1.9),
		},
		{
			name: "surface in shadow",
			args: args{
				m:        NewDefaultMaterial(),
				p:        NewPoint(0, 0, 0),
				l:        NewPointLight(NewPoint(0, 0, -10), ColorName(colornames.White)),
				eye:      NewVector(0, 0, -1),
				normal:   NewVector(0, 0, -1),
				inShadow: true,
			},
			o:    NewUnitSphere(),
			want: NewColor(0.1, 0.1, 0.1),
		},
		{
			name: "eye between light and surface, eye offset 45 degrees",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 0, -10), ColorName(colornames.White)),
				eye:    NewVector(0, math.Sqrt2/2, -math.Sqrt2/2),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(1.0, 1.0, 1.0),
		},
		{
			name: "eye opposite surface, light offset 45 degrees",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 10, -10), ColorName(colornames.White)),
				eye:    NewVector(0, 0, -1),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(0.7364, 0.7364, 0.7364),
		},
		{
			name: "eye in the path of the reflection vector",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 10, -10), ColorName(colornames.White)),
				eye:    NewVector(0, -math.Sqrt2/2, -math.Sqrt2/2),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(1.6364, 1.6364, 1.6364),
		},
		{
			name: "light behind the surface",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 0, 10), ColorName(colornames.White)),
				eye:    NewVector(0, 0, -1),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(0.1, 0.1, 0.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.want.Equal(
				lighting(tt.args.m, tt.o, tt.args.p, tt.args.l,
					tt.args.eye, tt.args.normal, tt.args.inShadow, 0, 0)))
		})
	}
}

func TestColorAtPoint(t *testing.T) {
	type args struct {
		m        *Material
		p        Point
		l        *PointLight
		eye      Vector
		normal   Vector
		inShadow bool
	}
	tests := []struct {
		name string
		args args
		o    Shaper
		want Color
	}{
		{
			name: "eye between light and surface",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 0, -10), ColorName(colornames.White)),
				eye:    NewVector(0, 0, -1),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(1.0, 1.0, 1.0),
		},
		{
			name: "eye between light and surface, eye offset 45 degrees",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 0, -10), ColorName(colornames.White)),
				eye:    NewVector(0, math.Sqrt2/2, -math.Sqrt2/2),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(1.0, 1.0, 1.0),
		},
		{
			name: "eye opposite surface, light offset 45 degrees",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 10, -10), ColorName(colornames.White)),
				eye:    NewVector(0, 0, -1),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(0.7364, 0.7364, 0.7364),
		},
		{
			name: "eye in the path of the reflection vector",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 10, -10), ColorName(colornames.White)),
				eye:    NewVector(0, -math.Sqrt2/2, -math.Sqrt2/2),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(1.0, 1.0, 1.0),
		},
		{
			name: "light behind the surface",
			args: args{
				m:      NewDefaultMaterial(),
				p:      NewPoint(0, 0, 0),
				l:      NewPointLight(NewPoint(0, 0, 10), ColorName(colornames.White)),
				eye:    NewVector(0, 0, -1),
				normal: NewVector(0, 0, -1),
			},
			o:    NewUnitSphere(),
			want: NewColor(0.1, 0.1, 0.1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.want.Equal(ColorAtPoint(tt.args.m, tt.o, tt.args.p, tt.args.l, tt.args.eye, tt.args.normal, tt.args.inShadow)))
		})
	}
}
