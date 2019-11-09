package tracer

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestNewUnitCube(t *testing.T) {
	tests := []struct {
		name string
		want *Cube
	}{
		{
			name: "test1",
			want: &Cube{
				Shape: Shape{
					transform:        IdentityMatrix(),
					transformInverse: IdentityMatrix().Inverse(),
					material:         NewDefaultMaterial(),
					shape:            "cube",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.Shape.lna = tt.want.localNormalAt
			got := NewUnitCube()
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}

func TestCube_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name string
		cube *Cube
		args args
		want Intersections
	}{
		{
			name: "1",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(5, 0.5, 0), NewVector(-1, 0, 0)),
			},
			want: NewIntersections(
				NewIntersection(NewUnitCube(), 4),
				NewIntersection(NewUnitCube(), 6),
			),
		},
		{
			name: "2",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(-5, 0.5, 0), NewVector(1, 0, 0)),
			},
			want: NewIntersections(
				NewIntersection(NewUnitCube(), 4),
				NewIntersection(NewUnitCube(), 6),
			),
		},
		{
			name: "3",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0.5, 5, 0), NewVector(0, -1, 0)),
			},
			want: NewIntersections(
				NewIntersection(NewUnitCube(), 4),
				NewIntersection(NewUnitCube(), 6),
			),
		},
		{
			name: "4",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0.5, -5, 0), NewVector(0, 1, 0)),
			},
			want: NewIntersections(
				NewIntersection(NewUnitCube(), 4),
				NewIntersection(NewUnitCube(), 6),
			),
		},
		{
			name: "5",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0.5, 0, 5), NewVector(0, 0, -1)),
			},
			want: NewIntersections(
				NewIntersection(NewUnitCube(), 4),
				NewIntersection(NewUnitCube(), 6),
			),
		},
		{
			name: "6",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0.5, 0.5, -5), NewVector(0, 0, 1)),
			},
			want: NewIntersections(
				NewIntersection(NewUnitCube(), 4),
				NewIntersection(NewUnitCube(), 6),
			),
		},
		{
			name: "inside",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0, 0.5, 0), NewVector(0, 0, 1)),
			},
			want: NewIntersections(
				NewIntersection(NewUnitCube(), -1),
				NewIntersection(NewUnitCube(), 1),
			),
		},
		{
			name: "miss1",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(-2, 0, 0), NewVector(0.2673, 0.5345, 0.8018)),
			},
			want: NewIntersections(),
		},
		{
			name: "miss2",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0, -2, 0), NewVector(0.8018, 0.2673, 0.5345)),
			},
			want: NewIntersections(),
		},
		{
			name: "miss3",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0, 0, -2), NewVector(0.5345, 0.8018, 0.2673)),
			},
			want: NewIntersections(),
		},
		{
			name: "miss4",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(2, 0, 2), NewVector(0, 0, -1)),
			},
			want: NewIntersections(),
		},
		{
			name: "miss5",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(0, 2, 2), NewVector(0, -1, 0)),
			},
			want: NewIntersections(),
		},
		{
			name: "miss6",
			cube: NewUnitCube(),
			args: args{
				r: NewRay(NewPoint(2, 2, 0), NewVector(-1, 0, 0)),
			},
			want: NewIntersections(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cube.IntersectWith(tt.args.r, NewIntersections())
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}

func TestCube_NormalAt(t *testing.T) {
	type args struct {
		p  Point
		xs *Intersection
	}
	tests := []struct {
		name string
		cube *Cube
		args args
		want Vector
	}{
		{
			name: "test1",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(1, 0.5, -0.8),
			},
			want: NewVector(1, 0, 0),
		},
		{
			name: "test2",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(-1, -0.2, 0.9),
			},
			want: NewVector(-1, 0, 0),
		},
		{
			name: "test3",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(-0.4, 1, -0.1),
			},
			want: NewVector(0, 1, 0),
		},
		{
			name: "test4",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(0.3, -1, -0.7),
			},
			want: NewVector(0, -1, 0),
		},
		{
			name: "test5",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(-0.6, 0.3, 1),
			},
			want: NewVector(0, 0, 1),
		},
		{
			name: "test6",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(0.4, 0.4, -1),
			},
			want: NewVector(0, 0, -1),
		},
		{
			name: "test7",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(1, 1, 1),
			},
			want: NewVector(1, 0, 0),
		},
		{
			name: "test8",
			cube: NewUnitCube(),
			args: args{
				p: NewPoint(-1, -1, -1),
			},
			want: NewVector(-1, 0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.cube.NormalAt(tt.args.p, tt.args.xs))
		})
	}
}
