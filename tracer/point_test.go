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

func TestPoint_Add(t *testing.T) {
	type args struct {
		t Vector
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
			assert.Equal(t, tt.p.AddVector(tt.args.t), tt.want, "should be equal")
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.p.TimesMatrix(tt.args.m), tt.want, "should be equal")
		})
	}
}
