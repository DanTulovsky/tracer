package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMatrix(t *testing.T) {
	type args struct {
		r int
		c int
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "4x4",
			args: args{
				r: 4,
				c: 4,
			},
			want: Matrix{
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
				{0, 0, 0, 0},
			},
		},
		{
			name: "3x3",
			args: args{
				r: 3,
				c: 3,
			},
			want: Matrix{
				{0, 0, 0},
				{0, 0, 0},
				{0, 0, 0},
			},
		},
		{
			name: "2x2",
			args: args{
				r: 2,
				c: 2,
			},
			want: Matrix{
				{0, 0},
				{0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewMatrix(tt.args.r, tt.args.c), "should be equal")
		})
	}
}

func TestNewMatrixFromData(t *testing.T) {
	type args struct {
		d [][]float64
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "4x4",
			args: args{
				d: [][]float64{
					{0, 0, 0, 0},
					{0, 0, 5, 0},
					{0, 0, 0, 7.0},
					{-3.1, 0, 0, 0},
				},
			},
			want: Matrix{
				{0, 0, 0, 0},
				{0, 0, 5, 0},
				{0, 0, 0, 7.0},
				{-3.1, 0, 0, 0},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMatrixFromData(tt.args.d)
			assert.Equal(t, tt.want, m, "should be equal")
			assert.Equal(t, 5.0, m[1][2], "should be equal")
			assert.Equal(t, -3.1, m[3][0], "should be equal")
			assert.Equal(t, 7.0, m[2][3], "should be equal")

		})
	}
}

func TestMatrix_Dims(t *testing.T) {
	tests := []struct {
		name  string
		m     Matrix
		wantR int
		wantC int
	}{
		{
			name:  "4x4",
			m:     NewMatrix(4, 4),
			wantR: 4,
			wantC: 4,
		},
		{
			name:  "2x2",
			m:     NewMatrix(2, 2),
			wantR: 2,
			wantC: 2,
		},
		{
			name:  "2x4",
			m:     NewMatrix(2, 4),
			wantR: 2,
			wantC: 4,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotR, gotC := tt.m.Dims()
			assert.Equal(t, tt.wantR, gotR, "should be equal")
			assert.Equal(t, tt.wantC, gotC, "should be equal")
		})
	}
}

func TestMatrix_Equals(t *testing.T) {
	type args struct {
		m2 Matrix
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want bool
	}{
		{
			name: "equals",
			m: NewMatrixFromData([][]float64{
				{1, 2, 3, 0},
				{4, 5, 6, 2},
				{7, 8, 9, 2},
				{7, 8, 9, 2},
			}),
			args: args{
				m2: NewMatrixFromData([][]float64{
					{1, 2, 3, 0},
					{4, 5, 6, 2},
					{7, 8, 9, 2},
					{7, 8, 9, 2},
				})},
			want: true,
		},
		{
			name: "different dimensions",
			m: NewMatrixFromData([][]float64{
				{1, 2, 3, 0},
				{4, 5, 6, 2},
				{7, 8, 9, 2},
			}),
			args: args{
				m2: NewMatrixFromData([][]float64{
					{1, 2, 3},
					{4, 5, 6},
				})},
			want: false,
		},
		{
			name: "different values",
			m: NewMatrixFromData([][]float64{
				{1, 2, 3, 0},
				{4, 5, 6, 2},
				{7, 8, 9, 2},
			}),
			args: args{
				m2: NewMatrixFromData([][]float64{
					{1, 2, 3, 0},
					{4, 9, 6, 2},
					{7, 8, 9, 2},
				})},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.m.Equals(tt.args.m2), tt.want, "should equal")
		})
	}
}

func TestMatrix_TimesMatrix(t *testing.T) {
	type args struct {
		m2 Matrix
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want Matrix
	}{
		{
			name: "valid1",
			m: NewMatrixFromData(Matrix{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 8, 7, 6},
				{5, 4, 3, 2},
			}),
			args: args{
				m2: NewMatrixFromData(Matrix{
					{-2, 1, 2, 3},
					{3, 2, 1, -1},
					{4, 3, 6, 5},
					{1, 2, 7, 8},
				}),
			},
			want: NewMatrixFromData(Matrix{
				{20, 22, 50, 48},
				{44, 54, 114, 108},
				{40, 58, 110, 102},
				{16, 26, 46, 42},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.m.TimesMatrix(tt.args.m2).Equals(tt.want), "should be true")
		})
	}

	panicTests := []struct {
		name string
		m    Matrix
		args args
		want Matrix
	}{
		{
			name: "invalid dimensions1",
			m: NewMatrixFromData(Matrix{
				{1, 2, 3, 4},
				{4, 6, 7, 8},
				{9, 8, 7, 6},
			}),
			args: args{
				m2: NewMatrixFromData(Matrix{
					{-2, 1, 2, 3},
					{3, 2, 1, -1},
					{4, 3, 6, 5},
					{1, 2, 7, 8},
				}),
			},
		},
		{
			name: "invalid dimensions2",
			m: NewMatrixFromData(Matrix{
				{1, 2, 3},
				{4, 6, 7},
				{9, 8, 7},
			}),
			args: args{
				m2: NewMatrixFromData(Matrix{
					{-2, 1, 2},
					{3, 2, 1},
					{4, 3, 6},
				}),
			},
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.m.TimesMatrix(tt.args.m2) }, "should panic")
		})
	}
}
