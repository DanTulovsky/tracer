package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPoint(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want Point
	}{
		{
			name: "origin",
			args: args{0, 0, 0},
			want: Point{0, 0, 0, 1},
		},
		{
			name: "point1",
			args: args{4.3, -4.2, 3.1},
			want: Point{4.3, -4.2, 3.1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPoint(tt.args.x, tt.args.y, tt.args.z)
			assert.Equal(t, p, tt.want, "should be equal")
			assert.Equal(t, p.W(), 1.0, "w should be 1")
			assert.Equal(t, p.X(), tt.want[0], "should be equal")
			assert.Equal(t, p.Y(), tt.want[1], "should be equal")
		})
	}
}

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

			assert.Equal(t, v, tt.want, "should be equal")
			assert.Equal(t, v.W(), 0.0, "w should be 0")
		})
	}
}

func TestPoint_Equals(t *testing.T) {
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		p    Point
		args args
		want bool
	}{
		{
			name: "equals",
			p:    NewPoint(1, 1, 1),
			args: args{
				t: NewPoint(1, 1, 1),
			},
			want: true,
		},
		{
			name: "not equals: point",
			p:    NewPoint(1, 1, 1),
			args: args{
				t: NewPoint(1, 2, 1),
			},
			want: false,
		},
		{
			name: "not equals: vector",
			p:    NewPoint(1, 1, 1),
			args: args{
				t: NewVector(1, 1, 1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.p.Equals(tt.args.t), tt.want, "should be equal")
		})
	}
}

func TestVector_Equals(t *testing.T) {
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want bool
	}{
		{
			name: "equals",
			v:    NewVector(1, 1, 1),
			args: args{
				t: NewVector(1, 1, 1),
			},
			want: true,
		},
		{
			name: "not equals: point",
			v:    NewVector(1, 1, 1),
			args: args{
				t: NewPoint(1, 1, 1),
			},
			want: false,
		},
		{
			name: "not equals: vector",
			v:    NewVector(1, 1, 1),
			args: args{
				t: NewVector(1, 2, 1),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.v.Equals(tt.args.t), tt.want, "should be equal")
		})
	}
}

func TestPoint_Add(t *testing.T) {
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Tuple
	}{
		{
			name: "add point and vector",
			p:    NewPoint(1, 1, 1),
			args: args{
				t: NewVector(2, 2, 2),
			},
			want: NewPoint(3, 3, 3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.p.Add(tt.args.t), tt.want, "should be equal")
		})
	}

	panicTests := []struct {
		name string
		p    Point
		args args
	}{
		{
			name: "add point and point",
			p:    NewPoint(1, 1, 1),
			args: args{
				t: NewPoint(2, 2, 2),
			},
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.p.Add(tt.args.t) }, "should panic")
		})
	}
}

func TestVector_Add(t *testing.T) {
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Tuple
	}{
		{
			name: "add vectors",
			v:    NewVector(1, 1, 1),
			args: args{
				t: NewVector(2, 2, 2),
			},
			want: NewVector(3, 3, 3),
		},
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
			assert.Equal(t, tt.v.Add(tt.args.t), tt.want, "should be equal")
		})
	}
}

func TestPoint_Sub(t *testing.T) {
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Tuple
	}{
		{
			name: "point - point",
			p:    NewPoint(3, 2, 1),
			args: args{
				t: NewPoint(5, 6, 7),
			},
			want: NewVector(-2, -4, -6),
		},
		{
			name: "point - vector",
			p:    NewPoint(3, 2, 1),
			args: args{
				t: NewVector(5, 6, 7),
			},
			want: NewPoint(-2, -4, -6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.p.Sub(tt.args.t), tt.want, "should be equal")
		})
	}
}

func TestVector_Sub(t *testing.T) {
	type args struct {
		t Tuple
	}
	tests := []struct {
		name string
		v    Vector
		args args
		want Tuple
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
			assert.Equal(t, tt.v.Sub(tt.args.t), tt.want, "should be equal")
		})
	}

	panicTests := []struct {
		name string
		v    Vector
		args args
	}{
		{
			name: "vector - point",
			v:    NewVector(3, 2, 1),
			args: args{
				t: NewPoint(5, 6, 7),
			},
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.v.Sub(tt.args.t) }, "should panic")
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
			assert.Equal(t, tt.v.Negate(), tt.want, "should be equal")
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
			assert.Equal(t, tt.v.Scale(tt.args.s), tt.want, "should be equal")
		})
	}
}
