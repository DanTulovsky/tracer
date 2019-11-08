package tracer

import (
	"fmt"
	"math"
	"testing"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultCone(t *testing.T) {
	tests := []struct {
		name string
		want *Cone
	}{
		{
			name: "test1",
			want: &Cone{
				Minimum: -math.MaxFloat64,
				Maximum: math.MaxFloat64,
				Closed:  false,
				Shape: Shape{
					transform:        IdentityMatrix(),
					transformInverse: IdentityMatrix().Inverse(),
					material:         NewDefaultMaterial(),
					shape:            "cone",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.Shape.lna = tt.want.localNormalAt
			got := NewDefaultCone()
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}

func TestNewCone(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want *Cone
	}{
		{
			name: "test1",
			want: &Cone{
				Minimum: -5,
				Maximum: 4,
				Closed:  false,
				Shape: Shape{
					transform:        IdentityMatrix(),
					transformInverse: IdentityMatrix().Inverse(),
					material:         NewDefaultMaterial(),
					shape:            "cone",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.Shape.lna = tt.want.localNormalAt
			got := NewCone(-5, 4)
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}

func TestNewClosedCone(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want *Cone
	}{
		{
			name: "test1",
			want: &Cone{
				Minimum: -5,
				Maximum: 4,
				Closed:  true,
				Shape: Shape{
					transform:        IdentityMatrix(),
					transformInverse: IdentityMatrix().Inverse(),
					material:         NewDefaultMaterial(),
					shape:            "cone",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.Shape.lna = tt.want.localNormalAt
			got := NewClosedCone(-5, 4)
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}

func TestCone_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name   string
		c      *Cone
		args   args
		wantXS int // how many intersections
		wantT  []float64
	}{
		{
			name: "test1",
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1).Normalize()),
			},
			c:     NewDefaultCone(),
			wantT: []float64{5, 5},
		},
		{
			name: "test2",
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(1, 1, 1).Normalize()),
			},
			c:     NewDefaultCone(),
			wantT: []float64{8.66025, 8.66025},
		},
		{
			name: "test3",
			args: args{
				r: NewRay(NewPoint(1, 1, -5), NewVector(-0.5, -1, 1).Normalize()),
			},
			c:     NewDefaultCone(),
			wantT: []float64{4.55006, 49.44994},
		},
		{
			name: "test4",
			args: args{
				r: NewRay(NewPoint(0, 0, -1), NewVector(0, 1, 1).Normalize()),
			},
			c:     NewDefaultCone(),
			wantT: []float64{0.35355},
		},
		{
			name: "end cap1",
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 1, 0).Normalize()),
			},
			c:     NewClosedCone(-0.5, 0.5),
			wantT: []float64{},
		},
		{
			name: "end cap2",
			args: args{
				r: NewRay(NewPoint(0, 0, -0.25), NewVector(0, 1, 1).Normalize()),
			},
			c:     NewClosedCone(-0.5, 0.5),
			wantT: []float64{0.707106, 0.088388},
		},
		{
			name: "end cap3",
			args: args{
				r: NewRay(NewPoint(0, 0, -0.25), NewVector(0, 1, 0).Normalize()),
			},
			c:     NewClosedCone(-0.5, 0.5),
			wantT: []float64{-0.5, -0.25, 0.25, 0.5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := NewIntersections()

			for _, tvalue := range tt.wantT {
				want = append(want, NewIntersection(tt.c, tvalue))
			}

			got := tt.c.IntersectWith(tt.args.r, NewIntersections())

			assert.Equal(t, len(want), len(got), "should equal")

			for i, tvalue := range tt.wantT {
				assert.InDelta(t, tvalue, got[i].T(), constants.Epsilon, "within epsilon")
			}
		})
	}
}

func TestCone_NormalAt(t *testing.T) {
	type args struct {
		p  Point
		xs Intersection
	}
	tests := []struct {
		name string
		cone *Cone
		args args
		want Vector
	}{
		{
			name: "test1",
			cone: NewClosedCone(-0.5, 0.5),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: NewVector(0, 0, 0),
		},
		{
			name: "test2",
			cone: NewClosedCone(-0.5, 0.5),
			args: args{
				p: NewPoint(1, 1, 1),
			},
			want: NewVector(1, -math.Sqrt2, 1),
		},
		{
			name: "test3",
			cone: NewClosedCone(-0.5, 0.5),
			args: args{
				p: NewPoint(-1, -1, 0),
			},
			want: NewVector(-1, 1, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.want.Normalize().Equal(tt.cone.NormalAt(tt.args.p, tt.args.xs)), "should equal")
		})
	}
}

func TestCone_Bounds(t *testing.T) {
	tests := []struct {
		name string
		c    *Cone
		want Bound
	}{
		{
			name: "default inf",
			c:    NewDefaultCone(),
			want: Bound{
				Min: NewPoint(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64),
				Max: NewPoint(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64),
			},
		},
		{
			name: "capped",
			c:    NewCone(-5, 5),
			want: Bound{
				Min: NewPoint(-5, -5, -5),
				Max: NewPoint(5, 5, 5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.PrecomputeValues()
			assert.Equal(t, tt.want, tt.c.Bounds(), "should equal")
		})
	}
}
