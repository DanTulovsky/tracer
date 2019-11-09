package tracer

import (
	"fmt"
	"math"
	"testing"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/google/go-cmp/cmp"
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
		want *Intersection
	}{
		{
			name: "test1",
			args: args{
				o: NewUnitSphere(),
				t: 3.0,
			},
			want: &Intersection{
				o: NewUnitSphere(),
				t: 3.0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIntersection(tt.args.o, tt.args.t)
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}

func TestIntersections_Hit(t *testing.T) {
	tests := []struct {
		name    string
		i       Intersections
		want    *Intersection
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
				diff := cmp.Diff(tt.want, got)
				assert.Equal(t, "", fmt.Sprint(diff))
			}
		})
	}
}

func TestPrepareComputations(t *testing.T) {
	type args struct {
		i *Intersection
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
			xs := NewIntersections(tt.args.i)
			comps := PrepareComputations(tt.args.i, tt.args.r, xs)
			assert.InEpsilon(t, tt.want.T, comps.T, constants.Epsilon, "should equal")
			assert.True(t, tt.want.Point.Equal(comps.Point))
			assert.True(t, tt.want.EyeV.Equal(comps.EyeV))
			assert.True(t, tt.want.NormalV.Equal(comps.NormalV))
			assert.True(t, tt.want.OverPoint.Equal(comps.OverPoint))
			assert.True(t, tt.want.ReflectV.Equal(comps.ReflectV))
		})
	}
}

func TestUnderPoint(t *testing.T) {

	ray := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	shape := NewGlassSphere()
	shape.SetTransform(IdentityMatrix().Translate(0, 0, 1))

	i := NewIntersection(shape, 5)
	xs := NewIntersections(i)

	comps := PrepareComputations(i, ray, xs)

	assert.Less(t, constants.Epsilon/2, comps.UnderPoint.Z(), "should be less")
	assert.Less(t, comps.Point.Z(), comps.UnderPoint.Z(), "should be less")
}

func Test_findRefractiveIndexes(t *testing.T) {
	tests := []struct {
		name   string
		wantN1 float64
		wantN2 float64
		index  int
	}{
		{
			name:   "test1",
			index:  0,
			wantN1: 1.0,
			wantN2: 1.5,
		},
		{
			name:   "test2",
			index:  1,
			wantN1: 1.5,
			wantN2: 2.0,
		},
		{
			name:   "test3",
			index:  2,
			wantN1: 2.0,
			wantN2: 2.5,
		},
		{
			name:   "test4",
			index:  3,
			wantN1: 2.5,
			wantN2: 2.5,
		},
		{
			name:   "test5",
			index:  4,
			wantN1: 2.5,
			wantN2: 1.5,
		},
		{
			name:   "test6",
			index:  5,
			wantN1: 1.5,
			wantN2: 1.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			transformA := IdentityMatrix().Scale(2, 2, 2)
			transformB := IdentityMatrix().Translate(0, 0, -0.25)
			transformC := IdentityMatrix().Translate(0, 0, 0.25)

			glass1 := NewGlassSphere()
			glass1.SetTransform(transformA)
			glass1.Material().RefractiveIndex = 1.5

			glass2 := NewGlassSphere()
			glass2.SetTransform(transformB)
			glass2.Material().RefractiveIndex = 2.0

			glass3 := NewGlassSphere()
			glass3.SetTransform(transformC)
			glass3.Material().RefractiveIndex = 2.5

			xs := NewIntersections(
				NewIntersection(glass1, 2),
				NewIntersection(glass2, 2.75),
				NewIntersection(glass3, 3.25),
				NewIntersection(glass2, 4.75),
				NewIntersection(glass3, 5.25),
				NewIntersection(glass1, 6),
			)

			n1, n2 := findRefractiveIndexes(xs[tt.index], xs)
			assert.Equal(t, tt.wantN1, n1, "should equal")
			assert.Equal(t, tt.wantN2, n2, "should equal")
		})
	}
}

func Test_objectInList(t *testing.T) {
	o := NewUnitSphere()
	list := []Shaper{o}
	want := true
	assert.Equal(t, want, objectInList(o, list))

	o = NewUnitSphere()
	list = []Shaper{}
	want = false
	assert.Equal(t, want, objectInList(o, list))
}

func Test_findObjectInList(t *testing.T) {
	o := NewUnitSphere()
	list := []Shaper{o}
	want := 0
	got, err := findObjectInList(o, list)
	assert.Equal(t, want, got, "should equal")
	assert.Nil(t, err, "no error")

	o = NewUnitSphere()
	list = []Shaper{NewUnitSphere(), NewUnitCube()}
	want = 2
	got, err = findObjectInList(o, list)
	assert.Equal(t, want, got, "should equal")
	assert.Error(t, err, "error")

	o = NewUnitSphere()
	list = []Shaper{NewUnitCube(), o, NewPlane()}
	want = 1
	got, err = findObjectInList(o, list)
	assert.Equal(t, want, got, "should equal")
	assert.Nil(t, err, "no error")
}

func Test_delObjectFromList(t *testing.T) {
	o := NewUnitSphere()
	list := []Shaper{NewUnitCube(), o, NewPlane()}
	list = delObjectFromList(o, list)

	// check to make sure it's gon
	got, err := findObjectInList(o, list)
	want := 2
	assert.Equal(t, want, got, "should equal")
	assert.Error(t, err, "error")
}

func TestSchlick(t *testing.T) {
	type args struct {
		s *IntersectionState
	}
	tests := []struct {
		name    string
		args    args
		shape   Shaper
		ray     Ray
		xsT     []float64
		xsIndex int
		want    float64
	}{
		{
			name:    "total internal reflection",
			shape:   NewGlassSphere(),
			ray:     NewRay(NewPoint(0, 0, math.Sqrt2/2), NewVector(0, 1, 0)),
			xsT:     []float64{-math.Sqrt2 / 2, math.Sqrt2 / 2},
			xsIndex: 1,
			want:    1.0,
		},
		{
			name:    "perpendicular viewing angle",
			shape:   NewGlassSphere(),
			ray:     NewRay(NewPoint(0, 0, 0), NewVector(0, 1, 0)),
			xsT:     []float64{-1, 1},
			xsIndex: 1,
			want:    0.04,
		},
		{
			name:    "small angle and n2 > n1",
			shape:   NewGlassSphere(),
			ray:     NewRay(NewPoint(0, 0.99, -2), NewVector(0, 0, 1)),
			xsT:     []float64{1.8589},
			xsIndex: 0,
			want:    0.48873,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xs := NewIntersections()

			for i := range tt.xsT {
				xs = append(xs, NewIntersection(tt.shape, tt.xsT[i]))
			}

			state := PrepareComputations(xs[tt.xsIndex], tt.ray, xs)
			got := Schlick(state)

			assert.InDelta(t, tt.want, got, constants.Epsilon, "should equal")
		})
	}
}

func TestNewIntersectionUV(t *testing.T) {
	type args struct {
		o Shaper
		t float64
		u float64
		v float64
	}
	tests := []struct {
		name string
		args args
		want *Intersection
	}{
		{
			name: "test1",
			args: args{
				o: newTestTriangle(),
				t: 3.5,
				u: 0.2,
				v: 0.4,
			},
			want: &Intersection{
				o: newTestTriangle(),
				t: 3.5,
				u: 0.2,
				v: 0.4,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewIntersectionUV(tt.args.o, tt.args.t, tt.args.u, tt.args.v)
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
		})
	}
}
