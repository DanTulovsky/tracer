package tracer

// Point is a single point in 3D space. p[3] is always 1. Implements Tupler.
type Point struct {
	Tuple
}

// NewPoint returns a new Point
func NewPoint(x, y, z float64) Point {
	return Point{
		Tuple{x, y, z, 1.0},
	}
}

// AddVector adds a point to a vector
func (p Point) AddVector(t Vector) Point {
	return NewPoint(p.X()+t.X(), p.Y()+t.Y(), p.Z()+t.Z())
}

// SubPoint subtracts points
func (p Point) SubPoint(t Point) Vector {
	return NewVector(p.X()-t.X(), p.Y()-t.Y(), p.Z()-t.Z())
}

// SubVector subtracts a vector
func (p Point) SubVector(t Vector) Point {
	return NewPoint(p.X()-t.X(), p.Y()-t.Y(), p.Z()-t.Z())
}

// Scale scales the point
func (p Point) Scale(s float64) Point {
	return NewPoint(p.X()*s, p.Y()*s, p.Z()*s)
}

// TimesMatrix multiplies a point by the matrix
func (p Point) TimesMatrix(m Matrix) Point {
	return NewPoint(
		m[0][0]*p.X()+m[0][1]*p.Y()+m[0][2]*p.Z()+m[0][3]*p.W(),
		m[1][0]*p.X()+m[1][1]*p.Y()+m[1][2]*p.Z()+m[1][3]*p.W(),
		m[2][0]*p.X()+m[2][1]*p.Y()+m[2][2]*p.Z()+m[2][3]*p.W())
}
