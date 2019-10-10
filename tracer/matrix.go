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

// Transpose transposes a YxY matrix
func (m Matrix) Transpose() Matrix {
	mR, mC := m.Dims()
	if mR != mC {
		panic("can only handle matrixies with same number of rows and columns")
	}

	new := NewMatrix(mR, mC)

	for r := 0; r < mR; r++ {
		for c := 0; c < mC; c++ {
			new[r][c] = m[c][r]
		}
	}

	return new
}

// Determinant returns the determinant of a 2x2 matrix
func (m Matrix) Determinant() float64 {

	mR, mC := m.Dims()
	if mR != mC {
		panic("can only handle matricies with same number of row and col")
	}

	if mR == 2 {
		return m[0][0]*m[1][1] - m[0][1]*m[1][0]
	}

	// matricies larger than 2x2
	result := 0.0

	for c := 0; c < mC; c++ {
		result = result + m[0][c]*m.Cofactor(0, c)
	}
	return result
}

// Submatrix returns a submatrix
func (m Matrix) Submatrix(row, col int) Matrix {

	mR, mC := m.Dims()
	if mR != mC {
		panic("can only handle matrixies with same number of rows and columns")
	}

	new := NewMatrix(mR-1, mC-1)

	for r := 0; r < mR; r++ {
		for c := 0; c < mC; c++ {
			if r == row || c == col {
				continue
			}
			rnew := r
			if r > row {
				rnew = r - 1
			}
			cnew := c
			if c > col {
				cnew = c - 1
			}
			new[rnew][cnew] = m[r][c]
		}
	}

	return new
}

// Minor returns the minr of a matrix
func (m Matrix) Minor(row, col int) float64 {
	return m.Submatrix(row, col).Determinant()
}

// Cofactor returns the cofactor of a matrix at row, col
func (m Matrix) Cofactor(row, col int) float64 {
	minor := m.Minor(row, col)
	if math.Mod(float64(row)+float64(col), 2) == 0 {
		return minor
	}
	return -minor
}

// IsInvertible return true if the matrix is invertible
func (m Matrix) IsInvertible() bool {
	return m.Determinant() != 0
}

// Inverse returns the inverse of a matrix
func (m Matrix) Inverse() Matrix {

	mR, mC := m.Dims()
	if mR != mC {
		panic("can only handle matrixies with same number of rows and columns")
	}

	new := NewMatrix(mR, mC)

	d := m.Determinant()

	for r := 0; r < mR; r++ {
		for c := 0; c < mC; c++ {
			cfactor := m.Cofactor(r, c)
			new[c][r] = cfactor / d
		}
	}

	return new
}
