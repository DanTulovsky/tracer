package tracer

import (
	"testing"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/stretchr/testify/assert"
)

func TestNewTriangle(t *testing.T) {
	type args struct {
		p1 Point
		p2 Point
		p3 Point
	}
	tests := []struct {
		name string
		args args
		want *Triangle
	}{
		{
			args: args{
				p1: NewPoint(0, 1, 0),
				p2: NewPoint(-1, 0, 0),
				p3: NewPoint(1, 0, 0),
			},
			want: &Triangle{
				P1:     NewPoint(0, 1, 0),
				P2:     NewPoint(-1, 0, 0),
				P3:     NewPoint(1, 0, 0),
				E1:     NewVector(-1, -1, 0),
				E2:     NewVector(1, -1, 0),
				Normal: NewVector(0, 0, -1),
				Shape: Shape{
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "triangle",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewTriangle(tt.args.p1, tt.args.p2, tt.args.p3))
		})
	}
}

func TestTriangle_NormalAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name string
		tri  *Triangle
		args args
	}{
		{
			tri: NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0)),
			args: args{
				p: NewPoint(0, 0.5, 0),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := tt.tri.Normal
			assert.Equal(t, want, tt.tri.NormalAt(tt.args.p))
		})
	}
}

func TestTriangle_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name   string
		tri    *Triangle
		args   args
		wantXS int
		wantT1 float64 // triangles only have one intersection at most
	}{
		{
			name: "miss1",
			tri:  NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0)),
			args: args{
				r: NewRay(NewPoint(0, -1, -2), NewVector(0, 1, 0)),
			},
			wantXS: 0,
		},
		{
			name: "miss p1-p3 edge",
			tri:  NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0)),
			args: args{
				r: NewRay(NewPoint(1, 1, -2), NewVector(0, 0, 1)),
			},
			wantXS: 0,
		},
		{
			name: "miss p1-p2 edge",
			tri:  NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0)),
			args: args{
				r: NewRay(NewPoint(-1, 1, -2), NewVector(0, 0, 1)),
			},
			wantXS: 0,
		},
		{
			name: "miss p2-p3 edge",
			tri:  NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0)),
			args: args{
				r: NewRay(NewPoint(0, -1, -2), NewVector(0, 0, 1)),
			},
			wantXS: 0,
		},
		{
			name: "hit",
			tri:  NewTriangle(NewPoint(0, 1, 0), NewPoint(-1, 0, 0), NewPoint(1, 0, 0)),
			args: args{
				r: NewRay(NewPoint(0, 0.5, -2), NewVector(0, 0, 1)),
			},
			wantXS: 1,
			wantT1: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := NewIntersections()

			if tt.wantXS > 0 {
				want = append(want, NewIntersection(tt.tri, tt.wantT1))
			}
			got := tt.tri.IntersectWith(tt.args.r)

			assert.Equal(t, len(want), len(got), "should equal")
			if tt.wantXS > 0 {
				assert.InDelta(t, tt.wantT1, got[0].T(), constants.Epsilon, "within epsilon")
			}
		})
	}
}
