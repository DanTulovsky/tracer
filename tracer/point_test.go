package tracer

import (
	"math"
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
			want: Point{Tuple{0, 0, 0, 1}},
		},
		{
			name: "point1",
			args: args{4.3, -4.2, 3.1},
			want: Point{Tuple{4.3, -4.2, 3.1, 1}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPoint(tt.args.x, tt.args.y, tt.args.z)
			assert.Equal(t, p, tt.want, "should be equal")
			assert.Equal(t, p.W(), 1.0, "w should be 1")
			assert.Equal(t, p.X(), tt.want.x, "should be equal")
			assert.Equal(t, p.Y(), tt.want.y, "should be equal")
		})
	}
}

func TestPoint_Add(t *testing.T) {
	type args struct {
		t Vector
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Point
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
			assert.Equal(t, tt.want, tt.p.AddVector(tt.args.t), "should be equal")
		})
	}

}

func TestPoint_SubPoint(t *testing.T) {
	type args struct {
		t Point
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Vector
	}{
		{
			name: "point - point",
			p:    NewPoint(3, 2, 1),
			args: args{
				t: NewPoint(5, 6, 7),
			},
			want: NewVector(-2, -4, -6),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.p.SubPoint(tt.args.t), tt.want, "should be equal")
		})
	}
}

func TestPoint_SubVector(t *testing.T) {
	type args struct {
		t Vector
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Point
	}{
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
			assert.Equal(t, tt.p.SubVector(tt.args.t), tt.want, "should be equal")
		})
	}
}

func TestPoint_TimesMatrix(t *testing.T) {
	type args struct {
		m Matrix
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Point
	}{
		{
			name: "test1",
			p:    NewPoint(1, 2, 3),
			args: args{
				m: NewMatrixFromData([][]float64{
					{1, 2, 3, 4},
					{2, 4, 4, 2},
					{8, 6, 4, 1},
					{0, 0, 0, 1},
				}),
			},
			want: NewPoint(18, 24, 33),
		},
		{
			name: "translation",
			p:    NewPoint(-3, 4, 5),
			args: args{
				m: NewTranslation(5, -3, 2),
			},
			want: NewPoint(2, 1, 7),
		},
		{
			name: "translation-inverse",
			p:    NewPoint(-3, 4, 5),
			args: args{
				m: NewTranslation(5, -3, 2).Inverse(),
			},
			want: NewPoint(-8, 7, 3),
		},
		{
			name: "scaling",
			p:    NewPoint(-4, 6, 8),
			args: args{
				m: NewScaling(2, 3, 4),
			},
			want: NewPoint(-8, 18, 32),
		},
		{
			name: "reflection",
			p:    NewPoint(-1, 1, 1),
			args: args{
				m: NewScaling(2, 3, 4),
			},
			want: NewPoint(-2, 3, 4),
		},
		{
			name: "rotateX",
			p:    NewPoint(0, 1, 0),
			args: args{
				m: NewRotationX(math.Pi / 4),
			},
			want: NewPoint(0, math.Sqrt2/2, math.Sqrt2/2),
		},
		{
			name: "rotateX.2",
			p:    NewPoint(0, 1, 0),
			args: args{
				m: NewRotationX(math.Pi / 2),
			},
			want: NewPoint(0, 0, 1),
		},
		{
			name: "rotateX inverse",
			p:    NewPoint(0, 1, 0),
			args: args{
				m: NewRotationX(math.Pi / 4).Inverse(),
			},
			want: NewPoint(0, math.Sqrt2/2, -math.Sqrt2/2),
		},
		{
			name: "rotateY",
			p:    NewPoint(0, 0, 1),
			args: args{
				m: NewRotationY(math.Pi / 4),
			},
			want: NewPoint(math.Sqrt2/2, 0, math.Sqrt2/2),
		},
		{
			name: "rotateY.2",
			p:    NewPoint(0, 0, 1),
			args: args{
				m: NewRotationY(math.Pi / 2),
			},
			want: NewPoint(1, 0, 0),
		},
		{
			name: "rotateZ",
			p:    NewPoint(0, 1, 0),
			args: args{
				m: NewRotationZ(math.Pi / 4),
			},
			want: NewPoint(-math.Sqrt2/2, math.Sqrt2/2, 0),
		},
		{
			name: "rotateZ.2",
			p:    NewPoint(0, 1, 0),
			args: args{
				m: NewRotationZ(math.Pi / 2),
			},
			want: NewPoint(-1, 0, 0),
		},
		{
			name: "shear",
			p:    NewPoint(2, 3, 4),
			args: args{
				m: NewShearing(1, 0, 0, 0, 0, 0),
			},
			want: NewPoint(5, 3, 4),
		},
		{
			name: "shear.2",
			p:    NewPoint(2, 3, 4),
			args: args{
				m: NewShearing(0, 1, 0, 0, 0, 0),
			},
			want: NewPoint(6, 3, 4),
		},
		{
			name: "shear.3",
			p:    NewPoint(2, 3, 4),
			args: args{
				m: NewShearing(0, 0, 1, 0, 0, 0),
			},
			want: NewPoint(2, 5, 4),
		},
		{
			name: "shear.4",
			p:    NewPoint(2, 3, 4),
			args: args{
				m: NewShearing(0, 0, 0, 1, 0, 0),
			},
			want: NewPoint(2, 7, 4),
		},
		{
			name: "shear.5",
			p:    NewPoint(2, 3, 4),
			args: args{
				m: NewShearing(0, 0, 0, 0, 1, 0),
			},
			want: NewPoint(2, 3, 6),
		},
		{
			name: "shear.6",
			p:    NewPoint(2, 3, 4),
			args: args{
				m: NewShearing(0, 0, 0, 0, 0, 1),
			},
			want: NewPoint(2, 3, 7),
		},
		{
			name: "sequence",
			p:    NewPoint(1, 0, 1),
			args: args{
				m: IdentityMatrix().RotateX(math.Pi/2).Scale(5, 5, 5).Translate(10, 5, 7),
			},
			want: NewPoint(15, 0, 7),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.want.Equals(tt.p.TimesMatrix(tt.args.m)), "should be true")
		})
	}
}

func TestPoint_Scale(t *testing.T) {
	type args struct {
		s float64
	}
	tests := []struct {
		name string
		p    Point
		args args
		want Point
	}{
		{
			name: "scale by > 0",
			p:    NewPoint(1, -2, 3),
			args: args{
				s: 3.5,
			},
			want: NewPoint(3.5, -7, 10.5),
		},
		{
			name: "scale by fraction",
			p:    NewPoint(1, -2, 3),
			args: args{
				s: 0.5,
			},
			want: NewPoint(0.5, -1, 1.5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.p.Scale(tt.args.s), "should be equal")
		})
	}
}
