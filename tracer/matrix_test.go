package tracer

import (
	"math"
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

func TestIdentityMatrix(t *testing.T) {
	tests := []struct {
		name string
		want Matrix
	}{
		{
			name: "test1",
			want: Matrix{
				{1, 0, 0, 0},
				{0, 1, 0, 0},
				{0, 0, 1, 0},
				{0, 0, 0, 1},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IdentityMatrix(), "should be equal")
		})
	}
}

func TestMatrix_Transpose(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want Matrix
	}{
		{
			name: "test1",
			m: NewMatrixFromData([][]float64{
				{0, 9, 3, 0},
				{9, 8, 0, 8},
				{1, 8, 5, 3},
				{0, 0, 5, 8},
			}),
			want: NewMatrixFromData([][]float64{
				{0, 9, 1, 0},
				{9, 8, 8, 0},
				{3, 0, 5, 5},
				{0, 8, 3, 8},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.Transpose(), "should be equal")
		})
	}

	panicTests := []struct {
		name string
		m    Matrix
	}{
		{
			name: "fail1",
			m: NewMatrixFromData([][]float64{
				{9, 8, 0, 8},
				{1, 8, 5, 3},
				{0, 0, 5, 8},
			}),
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.m.Transpose() }, "should panic")
		})
	}
}

func TestMatrix_Determinant(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want float64
	}{
		{
			name: "2x2",
			m: NewMatrixFromData([][]float64{
				{1, 5},
				{-3, 2},
			}),
			want: 17.0,
		},
		{
			name: "3x3",
			m: NewMatrixFromData([][]float64{
				{1, 2, 6},
				{-5, 8, -4},
				{2, 6, 4},
			}),
			want: -196,
		},
		{
			name: "4x4",
			m: NewMatrixFromData([][]float64{
				{-2, -8, 3, 5},
				{-3, 1, 7, 3},
				{1, 2, -9, 6},
				{-6, 7, 7, -9},
			}),
			want: -4071,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.Determinant(), "should be equal")
		})
	}

	panicTests := []struct {
		name string
		m    Matrix
	}{
		{
			name: "fail1",
			m: NewMatrixFromData([][]float64{
				{9, 8, 0, 8},
				{1, 8, 5, 3},
				{0, 0, 5, 8},
			}),
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.m.Determinant() }, "should panic")
		})
	}
}

func TestMatrix_Submatrix(t *testing.T) {
	type args struct {
		r int
		c int
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want Matrix
	}{
		{
			name: "3x3",
			m: NewMatrixFromData([][]float64{
				{1, 5, 0},
				{-3, 2, 7},
				{0, 6, -3},
			}),
			args: args{
				r: 0,
				c: 2,
			},
			want: NewMatrixFromData([][]float64{
				{-3, 2},
				{0, 6},
			}),
		},
		{
			name: "4x4",
			m: NewMatrixFromData([][]float64{
				{-6, 1, 1, 6},
				{-8, 5, 8, 6},
				{-1, 0, 8, 2},
				{-7, 1, -1, 1},
			}),
			args: args{
				r: 2,
				c: 1,
			},
			want: NewMatrixFromData([][]float64{
				{-6, 1, 6},
				{-8, 8, 6},
				{-7, -1, 1},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.Submatrix(tt.args.r, tt.args.c), "should be equal")
		})
	}

	panicTests := []struct {
		name string
		m    Matrix
		args args
	}{
		{
			name: "fail1",
			m: NewMatrixFromData([][]float64{
				{9, 8, 0, 8},
				{1, 8, 5, 3},
				{0, 0, 5, 8},
			}),
		},
	}
	for _, tt := range panicTests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.m.Submatrix(tt.args.r, tt.args.c) }, "should panic")
		})
	}
}

func TestMatrix_Minor(t *testing.T) {
	type args struct {
		r int
		c int
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want float64
	}{
		{
			name: "test1",
			m: NewMatrixFromData([][]float64{
				{3, 5, 0},
				{2, -1, -7},
				{6, -1, 5},
			}),
			args: args{
				r: 1,
				c: 0,
			},
			want: 25.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.Minor(tt.args.r, tt.args.c), "should be equal")
		})
	}
}

func TestMatrix_Cofactor(t *testing.T) {
	type args struct {
		row int
		col int
	}
	tests := []struct {
		name string
		m    Matrix
		args args
		want float64
	}{
		{
			name: "test1",
			m: NewMatrixFromData([][]float64{
				{3, 5, 0},
				{2, -1, -7},
				{6, -1, 5},
			}),
			args: args{
				row: 0,
				col: 0,
			},
			want: -12,
		},
		{
			name: "test2",
			m: NewMatrixFromData([][]float64{
				{3, 5, 0},
				{2, -1, -7},
				{6, -1, 5},
			}),
			args: args{
				row: 1,
				col: 0,
			},
			want: -25,
		},
		{
			name: "3x3",
			m: NewMatrixFromData([][]float64{
				{1, 2, 6},
				{-5, 8, -4},
				{2, 6, 4},
			}),
			args: args{
				row: 0,
				col: 0,
			},
			want: 56,
		},
		{
			name: "3x3.2",
			m: NewMatrixFromData([][]float64{
				{1, 2, 6},
				{-5, 8, -4},
				{2, 6, 4},
			}),
			args: args{
				row: 0,
				col: 1,
			},
			want: 12,
		},
		{
			name: "3x3.3",
			m: NewMatrixFromData([][]float64{
				{1, 2, 6},
				{-5, 8, -4},
				{2, 6, 4},
			}),
			args: args{
				row: 0,
				col: 2,
			},
			want: -46,
		},
		{
			name: "4x4",
			m: NewMatrixFromData([][]float64{
				{-2, -8, 3, 5},
				{-3, 1, 7, 3},
				{1, 2, -9, 6},
				{-6, 7, 7, -9},
			}),
			args: args{
				row: 0,
				col: 0,
			},
			want: 690,
		},
		{
			name: "4x4.2",
			m: NewMatrixFromData([][]float64{
				{-2, -8, 3, 5},
				{-3, 1, 7, 3},
				{1, 2, -9, 6},
				{-6, 7, 7, -9},
			}),
			args: args{
				row: 0,
				col: 1,
			},
			want: 447,
		},
		{
			name: "4x4.3",
			m: NewMatrixFromData([][]float64{
				{-2, -8, 3, 5},
				{-3, 1, 7, 3},
				{1, 2, -9, 6},
				{-6, 7, 7, -9},
			}),
			args: args{
				row: 0,
				col: 2,
			},
			want: 210,
		},
		{
			name: "4x4.4",
			m: NewMatrixFromData([][]float64{
				{-2, -8, 3, 5},
				{-3, 1, 7, 3},
				{1, 2, -9, 6},
				{-6, 7, 7, -9},
			}),
			args: args{
				row: 0,
				col: 3,
			},
			want: 51,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.Cofactor(tt.args.row, tt.args.col), "should equal")
		})
	}
}

func TestMatrix_IsInvertible(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want bool
	}{
		{
			name: "invertible",
			m: NewMatrixFromData([][]float64{
				{6, 4, 4, 4},
				{5, 5, 7, 6},
				{4, -9, 3, -7},
				{9, 1, 7, -6},
			}),
			want: true,
		},
		{
			name: "not invertible",
			m: NewMatrixFromData([][]float64{
				{-4, 2, -2, -3},
				{9, 6, 2, 6},
				{0, -5, 1, -5},
				{0, 0, 0, 0},
			}),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.IsInvertible(), "should equal")
		})
	}
}

func TestMatrix_Inverse(t *testing.T) {
	tests := []struct {
		name string
		m    Matrix
		want Matrix
	}{
		{
			name: "test1",
			m: NewMatrixFromData([][]float64{
				{-5, 2, 6, -8},
				{1, -5, 1, 8},
				{7, 7, -6, -7},
				{1, -3, 7, 4},
			}),
			want: NewMatrixFromData([][]float64{
				{0.21805, 0.45113, 0.24060, -0.04511},
				{-0.80827, -1.45677, -0.44361, 0.52068},
				{-0.07895, -0.22368, -0.05263, 0.19737},
				{-0.52256, -0.81391, -0.30075, 0.30639},
			}),
		},
		{
			name: "test2",
			m: NewMatrixFromData([][]float64{
				{8, -5, 9, 2},
				{7, 5, 6, 1},
				{-6, 0, 9, 6},
				{-3, 0, -9, -4},
			}),
			want: NewMatrixFromData([][]float64{
				{-0.15385, -0.15385, -0.28205, -0.53846},
				{-0.07692, 0.12308, 0.02564, 0.03077},
				{0.35897, 0.35897, 0.43590, 0.92308},
				{-0.69231, -0.69231, -0.76923, -1.92308},
			}),
		},
		{
			name: "test3",
			m: NewMatrixFromData([][]float64{
				{9, 3, 0, 9},
				{-5, -2, -6, -3},
				{-4, 9, 6, 4},
				{-7, 6, 6, 2},
			}),
			want: NewMatrixFromData([][]float64{
				{-0.04074, -0.07778, 0.14444, -0.22222},
				{-0.07778, 0.03333, 0.36667, -0.33333},
				{-0.02901, -0.14630, -0.10926, 0.12963},
				{0.17778, 0.06667, -0.26667, 0.33333},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.True(t, tt.m.Inverse().Equals(tt.want), "should equal")
		})
	}
}

func TestMatrix_Inverse2(t *testing.T) {
	tests := []struct {
		name string
		a    Matrix
		b    Matrix
	}{
		{
			name: "test1",
			a: NewMatrixFromData([][]float64{
				{3, -9, 7, 3},
				{3, -8, 2, -9},
				{-4, 4, 4, 1},
				{-6, 5, -1, 1},
			}),
			b: NewMatrixFromData([][]float64{
				{8, 2, 2, 2},
				{3, -1, 7, 0},
				{7, 0, 5, 4},
				{6, -2, 0, 5},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.a.TimesMatrix(tt.b)
			assert.True(t, tt.a.Equals(c.TimesMatrix(tt.b.Inverse())), "should equal")
		})
	}
}

func TestNewTranslation(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "test1",
			args: args{
				x: 5,
				y: -3,
				z: 2,
			},
			want: NewMatrixFromData([][]float64{
				{1, 0, 0, 5},
				{0, 1, 0, -3},
				{0, 0, 1, 2},
				{0, 0, 0, 1},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewTranslation(tt.args.x, tt.args.y, tt.args.z), tt.want, "should equal")
		})
	}
}

func TestNewScaling(t *testing.T) {
	type args struct {
		x float64
		y float64
		z float64
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "test1",
			args: args{
				x: 5,
				y: -3,
				z: 2,
			},
			want: NewMatrixFromData([][]float64{
				{5, 0, 0, 0},
				{0, -3, 0, 0},
				{0, 0, 2, 0},
				{0, 0, 0, 1},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewScaling(tt.args.x, tt.args.y, tt.args.z), tt.want, "should equal")
		})
	}
}

func TestNewRotationX(t *testing.T) {
	type args struct {
		r float64
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "test1",
			args: args{
				r: math.Pi,
			},
			want: NewMatrixFromData([][]float64{
				{1, 0, 0, 0},
				{0, math.Cos(math.Pi), -math.Sin(math.Pi), 0},
				{0, math.Sin(math.Pi), math.Cos(math.Pi), 0},
				{0, 0, 0, 1},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewRotationX(math.Pi), tt.want, "should equal")
		})
	}
}

func TestNewRotationY(t *testing.T) {
	type args struct {
		r float64
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "test1",
			args: args{
				r: math.Pi,
			},
			want: NewMatrixFromData([][]float64{
				{math.Cos(math.Pi), 0, math.Sin(math.Pi), 0},
				{0, 1, 0, 0},
				{-math.Sin(math.Pi), 0, math.Cos(math.Pi), 0},
				{0, 0, 0, 1},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewRotationY(math.Pi), tt.want, "should equal")
		})
	}
}

func TestNewRotationZ(t *testing.T) {
	type args struct {
		r float64
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "test1",
			args: args{
				r: math.Pi,
			},
			want: NewMatrixFromData([][]float64{
				{math.Cos(math.Pi), -math.Sin(math.Pi), 0, 0},
				{math.Sin(math.Pi), math.Cos(math.Pi), 0, 0},
				{0, 0, 1, 0},
				{0, 0, 0, 1},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewRotationZ(math.Pi), tt.want, "should equal")
		})
	}
}

func TestNewShearing(t *testing.T) {
	type args struct {
		xy float64
		xz float64
		yx float64
		yz float64
		zx float64
		zy float64
	}
	tests := []struct {
		name string
		args args
		want Matrix
	}{
		{
			name: "test1",
			args: args{
				xy: 1,
				xz: 2,
				yx: 3,
				yz: 4,
				zx: 5,
				zy: 6,
			},
			want: NewMatrixFromData([][]float64{
				{1, 1, 2, 0},
				{3, 1, 4, 0},
				{5, 6, 1, 0},
				{0, 0, 0, 1},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, NewShearing(tt.args.xy, tt.args.xz, tt.args.yx, tt.args.yz, tt.args.zx, tt.args.zy), tt.want, "should equal")
		})
	}
}
