package tracer

import "math"

// Matrix defines a 2 diemnsional matrix
type Matrix [][]float64

// NewMatrix returns a new matrix or r rows and c columns
func NewMatrix(r, c int) Matrix {

	m := make([][]float64, r)

	for i := range m {
		m[i] = make([]float64, c)
	}

	return Matrix(m)
}

// NewMatrixFromData returns a new matrix from the passed in data
func NewMatrixFromData(d [][]float64) Matrix {
	return Matrix(d)
}

// IdentityMatrix returns a 4x4 identity Matrix
func IdentityMatrix() Matrix {
	return Matrix{
		{1, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

}

// Dims returns the row, column dimensions of the matrix
func (m Matrix) Dims() (r, c int) {
	return len(m), len(m[0])
}

// Equals compares the matrix with another one within a margin of error for each value
func (m Matrix) Equals(m2 Matrix) bool {
	mR, mC := m.Dims()
	m2R, m2C := m2.Dims()
	if mR != m2R || mC != m2C {
		return false
	}

	for x := 0; x < mR; x++ {
		for y := 0; y < mC; y++ {
			if math.Abs(m[x][y]-m2[x][y]) > Epsilon {
				return false
			}
		}
	}

	return true
}

// TimesMatrix multiplies m by m2 and returns a new matric, currently only handles 4x4 matricies
func (m Matrix) TimesMatrix(m2 Matrix) Matrix {
	mR, mC := m.Dims()
	m2R, m2C := m2.Dims()
	if mR != m2R || mC != m2C || mR != mC || mR != 4 {
		panic("can only handle 4x4 matricies")
	}

	new := NewMatrix(4, 4)

	for x := 0; x < mR; x++ {
		for y := 0; y < mC; y++ {
			new[x][y] = m[x][0]*m2[0][y] + m[x][1]*m2[1][y] + m[x][2]*m2[2][y] + m[x][3]*m2[3][y]
		}
	}

	return new
}

// TimesVector multiplies m by a vector and returns a new vector
// func (m Matrix) TimesVector(v Vector) Vector {
// 	mR, mC := m.Dims()
// 	if mR != 4 || mC != 4 {
// 		panic("can only handle 4x4 matricies")
// 	}

// 	new := NewVector(
// 		m[0][0]*v.X()+m[0][1]*v.Y()+m[0][2]*v.Z()+m[0][3]*v.W(),
// 		m[1][0]*v.X()+m[1][1]*v.Y()+m[1][2]*v.Z()+m[1][3]*v.W(),
// 		m[2][0]*v.X()+m[2][1]*v.Y()+m[2][2]*v.Z()+m[2][3]*v.W())
// 	// m[3][0]*v.X()+m[3][1]*v.Y()+m[3][2]*v.Z()+m[3][3]*v.W())

// 	return new
// }
