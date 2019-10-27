package tracer

import (
	"math"
	"testing"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultCylinder(t *testing.T) {
	tests := []struct {
		name string
		want *Cylinder
	}{
		{
			name: "test1",
			want: &Cylinder{
				Radius:  1.0,
				Minimum: math.Inf(-1),
				Maximum: math.Inf(1),
				Shape: Shape{
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "cylinder",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDefaultCylinder()
			tt.want.Shape.name = got.name // random uuid

			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestCylinder_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name           string
		c              *Cylinder
		args           args
		wantXS         int // how many intersections
		wantT1, wantT2 float64
	}{
		{
			name: "ray misses a cylinder 1",
			args: args{
				r: NewRay(NewPoint(1, 0, 0), NewVector(0, 1, 0).Normalize()),
			},
			c:      NewDefaultCylinder(),
			wantXS: 0,
		},
		{
			name: "ray misses a cylinder 2",
			args: args{
				r: NewRay(NewPoint(0, 0, 0), NewVector(0, 1, 0).Normalize()),
			},
			c:      NewDefaultCylinder(),
			wantXS: 0,
		},
		{
			name: "ray misses a cylinder 3",
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(1, 1, 1).Normalize()),
			},
			c:      NewDefaultCylinder(),
			wantXS: 0,
		},
		{
			name: "ray intersects cylinder 1",
			args: args{
				r: NewRay(NewPoint(1, 0, -5), NewVector(0, 0, 1).Normalize()),
			},
			c:      NewDefaultCylinder(),
			wantXS: 2,
			wantT1: 5,
			wantT2: 5,
		},
		{
			name: "ray intersects cylinder 2",
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1).Normalize()),
			},
			c:      NewDefaultCylinder(),
			wantXS: 2,
			wantT1: 4,
			wantT2: 6,
		},
		{
			name: "ray intersects cylinder 3",
			args: args{
				r: NewRay(NewPoint(0.5, 0, -5), NewVector(0.1, 1, 1).Normalize()),
			},
			c:      NewDefaultCylinder(),
			wantXS: 2,
			wantT1: 6.80798,
			wantT2: 7.08872,
		},
		{
			name: "truncated1",
			args: args{
				r: NewRay(NewPoint(0, 1.5, 0), NewVector(0.1, 1, 0).Normalize()),
			},
			c:      NewCylinder(1, 2),
			wantXS: 0,
		},
		{
			name: "truncated2",
			args: args{
				r: NewRay(NewPoint(0, 3, -5), NewVector(0, 0, 1).Normalize()),
			},
			c:      NewCylinder(1, 2),
			wantXS: 0,
		},
		{
			name: "truncated3",
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1).Normalize()),
			},
			c:      NewCylinder(1, 2),
			wantXS: 0,
		},
		{
			name: "truncated4",
			args: args{
				r: NewRay(NewPoint(0, 2, -5), NewVector(0, 0, 1).Normalize()),
			},
			c:      NewCylinder(1, 2),
			wantXS: 0,
		},
		{
			name: "truncated5",
			args: args{
				r: NewRay(NewPoint(0, 1, -5), NewVector(0, 0, 1).Normalize()),
			},
			c:      NewCylinder(1, 2),
			wantXS: 0,
		},
		{
			name: "truncated6 intersect",
			args: args{
				r: NewRay(NewPoint(0, 1.5, -2), NewVector(0, 0, 1).Normalize()),
			},
			c:      NewCylinder(1, 2),
			wantXS: 2,
			wantT1: 1,
			wantT2: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := NewIntersections()

			if tt.wantXS > 0 {
				want = append(want, NewIntersection(NewDefaultCylinder(), tt.wantT1))
				want = append(want, NewIntersection(NewDefaultCylinder(), tt.wantT2))
			}
			got := tt.c.IntersectWith(tt.args.r)

			assert.Equal(t, len(want), len(got), "should equal")
			if tt.wantXS > 0 {
				assert.InDelta(t, tt.wantT1, got[0].T(), constants.Epsilon, "within epsilon")
				assert.InDelta(t, tt.wantT2, got[1].T(), constants.Epsilon, "within epsilon")
			}
		})
	}
}

func TestCylinder_NormalAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name string
		c    *Cylinder
		args args
		want Vector
	}{
		{
			name: "test1",
			c:    NewDefaultCylinder(),
			args: args{
				p: NewPoint(1, 0, 0),
			},
			want: NewVector(1, 0, 0),
		},
		{
			name: "test2",
			c:    NewDefaultCylinder(),
			args: args{
				p: NewPoint(0, 5, -1),
			},
			want: NewVector(0, 0, -1),
		},
		{
			name: "test3",
			c:    NewDefaultCylinder(),
			args: args{
				p: NewPoint(0, -2, 1),
			},
			want: NewVector(0, 0, 1),
		},
		{
			name: "test3",
			c:    NewDefaultCylinder(),
			args: args{
				p: NewPoint(-1, 1, 0),
			},
			want: NewVector(-1, 0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.c.NormalAt(tt.args.p), "should equal")
		})
	}
}

func TestNewCylinder(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name      string
		args      args
		want      *Cylinder
		wantPanic bool
	}{
		{
			name: "test1",
			args: args{
				min: -2,
				max: 3,
			},
			want: &Cylinder{
				Radius:  1,
				Minimum: -2,
				Maximum: 3,
				Shape: Shape{
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "cylinder",
				},
			},
			wantPanic: false,
		},
		{
			name: "invalid",
			args: args{
				min: 2,
				max: -3,
			},
			wantPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				assert.Panics(t, func() { NewCylinder(tt.args.min, tt.args.max) }, "should panic")

			} else {
				assert.Equal(t, tt.want, NewCylinder(tt.args.min, tt.args.max), "should equal")
			}
		})
	}
}
