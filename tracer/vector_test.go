package tracer

import (
	"math"
	"testing"

	"github.com/DanTulovsky/tracer/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewVector(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want Vector
	}{
		{
			name: "origin",
			args: args{0.0, 0.0, 0.0},
			want: Vector{0.0, 0.0, 0.0, 0.0},
		},
		{
			name: "vector1",
			args: args{4.3, -4.2, 3.1},
			want: Vector{4.3, -4.2, 3.1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewVector(tt.args.x, tt.args.y, tt.args.z)

			assert.Equal(t, tt.want, v, "should be equal")
			assert.Equal(t, 0.0, v.W(), "w should be 0")
		})
	}
}

func TestVector_AddVector(t *testing.T) {
	type args struct {
		t Vector
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Vector
	}{
		{
			name: "add vectors",
			v:    NewVector(1, 1, 1),
			args: args{
				t: NewVector(2, 2, 2),
			},
			want: NewVector(3, 3, 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.AddVector(tt.args.t), "should be equal")
		})
	}
}

func TestVector_AddPoint(t *testing.T) {
	type args struct {
		t Point
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Point
	}{
		{
			name: "add vector and point",
			v:    NewVector(-2, 3, 1),
			args: args{
				t: NewPoint(3, -2, 5),
			},
			want: NewPoint(1, 1, 6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.AddPoint(tt.args.t), "should be equal")
		})
	}
}
func TestVector_SubVector(t *testing.T) {
	type args struct {
		t Vector
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Vector
	}{
		{
			name: "vector - vector",
			v:    NewVector(3, 2, 1),
			args: args{
				t: NewVector(5, 6, 7),
			},
			want: NewVector(-2, -4, -6),
		},
		{
			name: "vector - zero vector",
			v:    NewVector(0, 0, 0),
			args: args{
				t: NewVector(1, -2, 3),
			},
			want: NewVector(-1, 2, -3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.SubVector(tt.args.t), "should be equal")
		})
	}
}

func TestVector_Negate(t *testing.T) {
	tests := []struct {
		name string
		v    Vector
		want Vector
	}{
		{
			name: "negate",
			v:    NewVector(1, -2, 3),
			want: NewVector(-1, 2, -3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.Negate(), "should be equal")
		})
	}
}

func TestVector_Scale(t *testing.T) {
	type args struct {
		s float64
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Vector
	}{
		{
			name: "scale by > 0",
			v:    NewVector(1, -2, 3),
			args: args{
				s: 3.5,
			},
			want: NewVector(3.5, -7, 10.5),
		},
		{
			name: "scale by fraction",
			v:    NewVector(1, -2, 3),
			args: args{
				s: 0.5,
			},
			want: NewVector(0.5, -1, 1.5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.Scale(tt.args.s), "should be equal")
		})
	}
}

func TestVector_Magnitude(t *testing.T) {
	tests := []struct {
		name string
		v    Vector
		want float64
	}{
		{
			name: "test1",
			v:    NewVector(1, 0, 0),
			want: 1.0,
		},
		{
			name: "test2",
			v:    NewVector(0, 1, 0),
			want: 1.0,
		},
		{
			name: "test3",
			v:    NewVector(0, 0, 1),
			want: 1.0,
		},
		{
			name: "test4",
			v:    NewVector(1, 2, 3),
			want: math.Sqrt(14),
		},
		{
			name: "test5",
			v:    NewVector(-1, -2, -3),
			want: math.Sqrt(14),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, utils.Equals(tt.v.Magnitude(), tt.want), "should be equal")
		})
	}
}

func TestVector_Normalize(t *testing.T) {
	tests := []struct {
		name string
		v    Vector
		want Vector
	}{
		{
			name: "test1",
			v:    NewVector(4, 0, 0),
			want: NewVector(1, 0, 0),
		},
		{
			name: "test2",
			v:    NewVector(1, 2, 3),
			want: NewVector(1/math.Sqrt(14), 2/math.Sqrt(14), 3/math.Sqrt(14)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.Normalize(), "should be equal")
			assert.Equal(t, 1.0, tt.v.Normalize().Magnitude(), "should be equal")
		})
	}

	panicTests := []struct {
		name string
		v    Vector
	}{
		{
			name: "test1",
			v:    NewVector(0, 0, 0),
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.v.Normalize() }, "should panic")
		})
	}
}

func TestVector_Dot(t *testing.T) {
	type args struct {
		w Vector
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want float64
	}{
		{
			name: "test1",
			v:    NewVector(1, 2, 3),
			args: args{
				w: NewVector(2, 3, 4),
			},
			want: 20.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.Dot(tt.args.w), "should be equal")
		})
	}
}

func TestVector_Cross(t *testing.T) {
	type args struct {
		w Vector
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Vector
	}{
		{
			name: "test1",
			v:    NewVector(1, 2, 3),
			args: args{
				w: NewVector(2, 3, 4),
			},
			want: NewVector(-1, 2, -1),
		},
		{
			name: "test2",
			v:    NewVector(2, 3, 4),
			args: args{
				w: NewVector(1, 2, 3),
			},
			want: NewVector(1, -2, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.Cross(tt.args.w), "should be equal")
		})
	}
}

func TestVector_TimesMatrix(t *testing.T) {
	type args struct {
		m Matrix
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Vector
	}{
		{
			name: "test1",
			v:    NewVector(1, 2, 3),
			args: args{
				m: NewMatrixFromData([][]float64{
					{1, 2, 3, 4},
					{2, 4, 4, 2},
					{8, 6, 4, 1},
					{0, 0, 0, 1},
				}),
			},
			want: NewVector(14, 22, 32),
		},
		{
			name: "translation",
			v:    NewVector(-3, 4, 5),
			args: args{
				m: NewTranslation(5, -3, 2),
			},
			want: NewVector(-3, 4, 5),
		},
		{
			name: "scaling",
			v:    NewVector(-4, 6, 8),
			args: args{
				m: NewScaling(2, 3, 4),
			},
			want: NewVector(-8, 18, 32),
		},
		{
			name: "scaling inverse",
			v:    NewVector(-4, 6, 8),
			args: args{
				m: NewScaling(2, 3, 4).Inverse(),
			},
			want: NewVector(-2, 2, 2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.v.TimesMatrix(tt.args.m), "should be equal")
		})
	}
}
