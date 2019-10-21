package tracer

import (
	"math"
	"testing"

	"github.com/DanTulovsky/tracer/constants"

	"github.com/stretchr/testify/assert"
)

func TestNewIntersection(t *testing.T) {
	type args struct {
		o Shaper
		t float64
	}
	tests := []struct {
		name string
		args args
		want Intersection
	}{
		{
			name: "test1",
			args: args{
				o: NewUnitSphere(),
				t: 3.0,
			},
			want: Intersection{
				o: NewUnitSphere(),
				t: 3.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewIntersection(tt.args.o, tt.args.t), "should be equal")
		})
	}
}

func TestIntersections_Hit(t *testing.T) {
	tests := []struct {
		name    string
		i       Intersections
		want    Intersection
		wantErr bool
	}{
		{
			name: "t all positive",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), 1),
				NewIntersection(NewUnitSphere(), 2),
			),
			want:    NewIntersection(NewUnitSphere(), 1),
			wantErr: false,
		},
		{
			name: "t some negative",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), -1),
				NewIntersection(NewUnitSphere(), 1),
			),
			want:    NewIntersection(NewUnitSphere(), 1),
			wantErr: false,
		},
		{
			name: "t negative",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), -2),
				NewIntersection(NewUnitSphere(), -1),
			),
			wantErr: true,
		},
		{
			name: "t some negative 2",
			i: NewIntersections(
				NewIntersection(NewUnitSphere(), 5),
				NewIntersection(NewUnitSphere(), 7),
				NewIntersection(NewUnitSphere(), -3),
				NewIntersection(NewUnitSphere(), 2),
			),
			want:    NewIntersection(NewUnitSphere(), 2),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.i.Hit()

			if tt.wantErr {
				assert.Error(t, err, "should error")
			} else {
				assert.Equal(t, tt.want, got, "should equal")
			}
		})
	}
}

func TestPrepareComputations(t *testing.T) {
	type args struct {
		i Intersection
		r Ray
	}
	tests := []struct {
		name string
		args args
		want IntersectionState
	}{
		{
			name: "state1",
			args: args{
				i: NewIntersection(NewUnitSphere(), 4),
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			want: IntersectionState{
				T:         4,
				Object:    NewUnitSphere(),
				Point:     NewPoint(0, 0, -1),
				EyeV:      NewVector(0, 0, -1),
				NormalV:   NewVector(0, 0, -1),
				Inside:    false,
				OverPoint: NewPoint(0, 0, -1.00001),
				ReflectV:  NewVector(0, 0, -1),
			},
		},
		{
			name: "outside",
			args: args{
				i: NewIntersection(NewUnitSphere(), 4),
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			want: IntersectionState{
				T:         4,
				Object:    NewUnitSphere(),
				Point:     NewPoint(0, 0, -1),
				EyeV:      NewVector(0, 0, -1),
				NormalV:   NewVector(0, 0, -1),
				Inside:    false,
				OverPoint: NewPoint(0, 0, -1.00001),
				ReflectV:  NewVector(0, 0, -1),
			},
		},
		{
			name: "inside",
			args: args{
				i: NewIntersection(NewUnitSphere(), 1),
				r: NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1)),
			},
			want: IntersectionState{
				T:         1,
				Object:    NewUnitSphere(),
				Point:     NewPoint(0, 0, 1),
				EyeV:      NewVector(0, 0, -1),
				NormalV:   NewVector(0, 0, -1),
				Inside:    true,
				OverPoint: NewPoint(0, 0, 0.99999),
				ReflectV:  NewVector(0, 0, -1),
			},
		},
		{
			name: "reflection vector",
			args: args{
				i: NewIntersection(NewPlane(), math.Sqrt2),
				r: NewRay(NewPoint(0, 1, -1), NewVector(0, -math.Sqrt2/2, math.Sqrt2/2)),
			},
			want: IntersectionState{
				T:         math.Sqrt2,
				Object:    NewUnitSphere(),
				Point:     NewPoint(0, 0, 0),
				EyeV:      NewVector(0, math.Sqrt2/2, -math.Sqrt2/2),
				NormalV:   NewVector(0, 1, 0),
				Inside:    false,
				OverPoint: NewPoint(0, 0.00001, 0),
				ReflectV:  NewVector(0, math.Sqrt2/2, math.Sqrt2/2),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			comps := PrepareComputations(tt.args.i, tt.args.r)
			// assert.Equal(t, tt.want, comps, "should equal")
			assert.InEpsilon(t, tt.want.T, comps.T, constants.Epsilon, "should equal")
			assert.True(t, tt.want.Point.Equals(comps.Point))
			assert.True(t, tt.want.EyeV.Equals(comps.EyeV))
			assert.True(t, tt.want.NormalV.Equals(comps.NormalV))
			assert.True(t, tt.want.OverPoint.Equals(comps.OverPoint))
			assert.True(t, tt.want.ReflectV.Equals(comps.ReflectV))
		})
	}
}
