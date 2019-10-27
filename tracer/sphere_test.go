package tracer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
				Center: NewPoint(0, 0, 0),
				Radius: 1,
				Shape: Shape{
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "sphere",
				},
			},
		},
	}
	for _, tt := range tests {
		got := NewUnitSphere()
		tt.want.Shape.name = got.name // random uuid

		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, got, "should equal")
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

func TestSphere_NormalAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name   string
		sphere *Sphere
		args   args
		m      Matrix // transform matrix
		want   Vector
	}{
		{
			name:   "x-axis",
			sphere: NewUnitSphere(),
			args: args{
				p: NewPoint(1, 0, 0),
			},
			m:    IdentityMatrix(),
			want: NewVector(1, 0, 0),
		},
		{
			name:   "y-axis",
			sphere: NewUnitSphere(),
			args: args{
				p: NewPoint(0, 1, 0),
			},
			m:    IdentityMatrix(),
			want: NewVector(0, 1, 0),
		},
		{
			name:   "z-axis",
			sphere: NewUnitSphere(),
			args: args{
				p: NewPoint(0, 0, 1),
			},
			m:    IdentityMatrix(),
			want: NewVector(0, 0, 1),
		},
		{
			name:   "non-axial",
			sphere: NewUnitSphere(),
			args: args{
				p: NewPoint(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3),
			},
			m:    IdentityMatrix(),
			want: NewVector(math.Sqrt(3)/3, math.Sqrt(3)/3, math.Sqrt(3)/3),
		},
		{
			name:   "translated",
			sphere: NewUnitSphere(),
			args: args{
				p: NewPoint(0, 1.70711, -0.70711),
			},
			m:    IdentityMatrix().Translate(0, 1, 0),
			want: NewVector(0, 0.70711, -0.70711),
		},
		{
			name:   "transform",
			sphere: NewUnitSphere(),
			args: args{
				p: NewPoint(0, math.Sqrt2/2, -math.Sqrt2/2),
			},
			m:    IdentityMatrix().RotateZ(math.Pi/5).Scale(1, 0.5, 1),
			want: NewVector(0, 0.97014, -0.24254),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sphere.SetTransform(tt.m)
			v := tt.sphere.NormalAt(tt.args.p)

			assert.True(t, tt.want.Equals(v), "should equal")
			assert.True(t, v.Equals(v.Normalize()), "should equal")
		})
	}
}

func TestSphere_Material(t *testing.T) {
	tests := []struct {
		name   string
		sphere *Sphere
		want   *Material
	}{
		{
			name:   "default",
			sphere: NewUnitSphere(),
			want:   NewDefaultMaterial(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.sphere.Material(), "should equal")
		})
	}
}

func TestSphere_SetMaterial(t *testing.T) {
	type args struct {
		m *Material
	}
	tests := []struct {
		name   string
		sphere *Sphere
		args   args
		want   *Material
	}{
		{
			name:   "test1",
			sphere: NewUnitSphere(),
			args: args{
				m: NewDefaultMaterial(),
			},
			want: NewDefaultMaterial(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.sphere.SetMaterial(tt.args.m)
			assert.Equal(t, tt.want, tt.sphere.Material(), "should equal")

		})
	}
}

func TestNewGlassSphere(t *testing.T) {
	tests := []struct {
		name                    string
		wantTranparency, wantRI float64
	}{
		{
			name:            "test1",
			wantTranparency: 1.0,
			wantRI:          1.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewGlassSphere()
			assert.Equal(t, tt.wantTranparency, s.Material().Transparency, "should equal")
			assert.Equal(t, tt.wantRI, s.Material().RefractiveIndex, "should equal")
		})
	}
}
