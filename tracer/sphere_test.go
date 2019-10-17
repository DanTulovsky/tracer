package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSphere(t *testing.T) {
	type args struct {
		c Point
		r float64
	}
	tests := []struct {
		name string
		args args
		want Sphere
	}{
		{
			name: "test1",
			args: args{
				c: NewPoint(0, 0, 0),
				r: 1.0,
			},
			want: Sphere{NewPoint(0, 0, 0), 1.0, IdentityMatrix()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, &tt.want, NewSphere(tt.args.c, tt.args.r), "should equal")
		})
	}
}

func TestSphere_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name      string
		sphere    *Sphere
		args      args
		transform Matrix
		want      []float64
	}{
		{
			name:      "2 point intersect",
			sphere:    NewUnitSphere(),
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			want: []float64{4.0, 6.0},
		},
		{
			name:      "1 point intersect",
			sphere:    NewUnitSphere(),
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(NewPoint(0, 1, -5), NewVector(0, 0, 1)),
			},
			want: []float64{5.0, 5.0},
		},
		{
			name:      "0 point intersect",
			sphere:    NewUnitSphere(),
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(NewPoint(0, 2, -5), NewVector(0, 0, 1)),
			},
			want: []float64{},
		},
		{
			name:      "ray inside sphere",
			sphere:    NewUnitSphere(),
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1)),
			},
			want: []float64{-1.0, 1.0},
		},
		{
			name:      "sphere behind ray",
			sphere:    NewUnitSphere(),
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(NewPoint(0, 0, 5), NewVector(0, 0, 1)),
			},
			want: []float64{-6.0, -4.0},
		},
		{
			name:      "scaled sphere",
			sphere:    NewUnitSphere(),
			transform: NewScaling(2, 2, 2),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			want: []float64{3.0, 7.0},
		},
		{
			name:      "translated sphere",
			sphere:    NewUnitSphere(),
			transform: NewTranslation(5, 0, 0),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			want: []float64{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			want := Intersections{}
			for _, int := range tt.want {
				want = append(want, NewIntersection(tt.sphere, int))
			}

			tt.sphere.SetTransform(tt.transform)
			assert.Equalf(t, want, tt.sphere.IntersectWith(tt.args.r), "should equal")
		})
	}
}

func TestNewUnitSphere(t *testing.T) {
	tests := []struct {
		name string
		want *Sphere
	}{
		{
			name: "test1",
			want: &Sphere{
				Center:    NewPoint(0, 0, 0),
				Radius:    1,
				transform: IdentityMatrix(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewUnitSphere())
		})
	}
}

func TestSphere_Transform(t *testing.T) {
	tests := []struct {
		name   string
		sphere *Sphere
		want   Matrix
	}{
		{
			name:   "identity by default",
			sphere: NewUnitSphere(),
			want:   IdentityMatrix(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sphere.Transform(), "should equal")
		})
	}
}

func TestSphere_SetTransform(t *testing.T) {
	type args struct {
		m Matrix
	}
	tests := []struct {
		name   string
		sphere *Sphere
		args   args
		want   Matrix
	}{
		{
			name:   "test1",
			sphere: NewUnitSphere(),
			args: args{
				m: NewTranslation(2, 3, 4),
			},
			want: NewTranslation(2, 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sphere.SetTransform(tt.args.m)
			assert.Equal(t, tt.want, tt.sphere.Transform(), "should equal")

		})
	}
}
