package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPlane(t *testing.T) {
	tests := []struct {
		name string
		want *Plane
	}{
		{
			name: "test1",
			want: &Plane{
				Shape: Shape{
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "plane",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewPlane()
			tt.want.Shape.name = got.name // random uuid
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestPlane_NormalAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		plane *Plane
		args  args
		want  Vector
	}{
		{
			name:  "test1",
			plane: NewPlane(),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: NewVector(0, 1, 0),
		},
		{
			name:  "test2",
			plane: NewPlane(),
			args: args{
				p: NewPoint(10, 0, -10),
			},
			want: NewVector(0, 1, 0),
		},
		{
			name:  "test3",
			plane: NewPlane(),
			args: args{
				p: NewPoint(-5, 0, 150),
			},
			want: NewVector(0, 1, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.plane.NormalAt(tt.args.p))
		})
	}
}

func TestPlane_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name  string
		plane *Plane
		args  args
		want  Intersections
	}{
		{
			name:  "ray parallel to the plane",
			plane: NewPlane(),
			args: args{
				r: NewRay(NewPoint(0, 10, 0), NewVector(0, 0, 1)),
			},
			want: NewIntersections(),
		},
		{
			name:  "ray coplanar with plane",
			plane: NewPlane(),
			args: args{
				r: NewRay(NewPoint(0, 10, 0), NewVector(0, 0, 1)),
			},
			want: NewIntersections(),
		},
		{
			name:  "ray intersects from above",
			plane: NewPlane(),
			args: args{
				r: NewRay(NewPoint(0, 1, 0), NewVector(0, -1, 0)),
			},
			want: NewIntersections(
				NewIntersection(NewPlane(), 1),
			),
		},
		{
			name:  "ray intersects from below",
			plane: NewPlane(),
			args: args{
				r: NewRay(NewPoint(0, -1, 0), NewVector(0, 1, 0)),
			},
			want: NewIntersections(
				NewIntersection(NewPlane(), 1),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// fix random uuids
			for _, i := range tt.want {
				i.Object().SetName(tt.plane.Name())
			}

			assert.Equal(t, tt.want, tt.plane.IntersectWith(tt.args.r))
		})
	}
}
