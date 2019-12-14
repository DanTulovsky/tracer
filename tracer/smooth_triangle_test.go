package tracer

import (
	"fmt"
	"testing"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func newTestTriangle() *SmoothTriangle {
	return NewSmoothTriangle(
		NewPoint(0, 1, 0),
		NewPoint(-1, 0, 0),
		NewPoint(1, 0, 0),
		NewVector(0, 1, 0),
		NewVector(-1, 0, 0),
		NewVector(1, 0, 0),
		NewPoint(0, 1, 0),
		NewPoint(-1, 0, 0),
		NewPoint(1, 0, 0),
	)
}
func TestNewSmoothTriangle(t *testing.T) {
	type args struct {
		p1  Point
		p2  Point
		p3  Point
		n1  Vector
		n2  Vector
		n3  Vector
		vt1 Point
		vt2 Point
		vt3 Point
	}
	tests := []struct {
		name string
		args args
		want *SmoothTriangle
	}{
		{
			args: args{
				p1:  NewPoint(0, 1, 0),
				p2:  NewPoint(-1, 0, 0),
				p3:  NewPoint(1, 0, 0),
				n1:  NewVector(0, 1, 0),
				n2:  NewVector(-1, 0, 0),
				n3:  NewVector(1, 0, 0),
				vt1: NewPoint(0, 1, 0),
				vt2: NewPoint(-1, 0, 0),
				vt3: NewPoint(1, 0, 0),
			},
			want: &SmoothTriangle{
				N1:  NewVector(0, 1, 0),
				N2:  NewVector(-1, 0, 0),
				N3:  NewVector(1, 0, 0),
				VT1: NewPoint(0, 1, 0),
				VT2: NewPoint(-1, 0, 0),
				VT3: NewPoint(1, 0, 0),
				Triangle: Triangle{
					P1: NewPoint(0, 1, 0),
					P2: NewPoint(-1, 0, 0),
					P3: NewPoint(1, 0, 0),
					E1: NewVector(-1, -1, 0),
					E2: NewVector(1, -1, 0),
					Shape: Shape{
						transform:        IM(),
						transformInverse: IM().Inverse(),
						material:         NewDefaultMaterial(),
						shape:            "smooth-triangle",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.want.Shape.lna = tt.want.localNormalAt
			got := NewSmoothTriangle(
				tt.args.p1, tt.args.p2, tt.args.p3,
				tt.args.n1, tt.args.n2, tt.args.n3,
				tt.args.vt1, tt.args.vt2, tt.args.vt3)
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}

func TestSmoothTriangle_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name         string
		tri          *SmoothTriangle
		args         args
		wantu, wantv float64
	}{
		{
			tri: newTestTriangle(),
			args: args{
				r: NewRay(NewPoint(-0.2, 0.3, -2), NewVector(0, 0, 1)),
			},
			wantu: 0.45,
			wantv: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tri.IntersectWith(tt.args.r, NewIntersections())
			if assert.NotEqual(t, 0, len(got), "expected intersections") {
				assert.InDelta(t, tt.wantu, got[0].u, constants.Epsilon, "should equal")
				assert.InDelta(t, tt.wantv, got[0].v, constants.Epsilon, "should equal")
			}
		})
	}
}

func TestSmoothTriangle_NormalAt(t *testing.T) {
	type args struct {
		p  Point
		xs *Intersection
	}
	tests := []struct {
		name string
		tri  *SmoothTriangle
		args args
		want Vector
	}{
		{
			tri: newTestTriangle(),
			args: args{
				p:  NewPoint(0, 0, 0), // not used
				xs: NewIntersectionUV(newTestTriangle(), 1, 0.45, 0.25),
			},
			want: NewVector(-0.55470, 0.832050, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.want.Equal(tt.tri.NormalAt(tt.args.p, tt.args.xs)), "should equal")
		})
	}
}
