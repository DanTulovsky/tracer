package tracer

// Tuple is the base for Vector and Point
type Tuple interface {
	X() float64
	Y() float64
	Z() float64
	W() float64
	Equals(Tuple) bool
	Add(Tuple) Tuple
	Sub(Tuple) Tuple
}

// Point is a single point in 3D space. p[3] is always 1. Implements Tuple.
type Point [4]float64

// NewPoint returns a new Point
func NewPoint(x, y, z float64) Point {

	return Point{x, y, z, 1.0}
}

// X returns the point's X coordinate
func (p Point) X() float64 {
	return p[0]
}

// Y returns the point's y coordinate
func (p Point) Y() float64 {
	return p[1]
}

// Z returns the point's Z coordinate
func (p Point) Z() float64 {
	return p[2]
}

// W returns the point's W coordinate
func (p Point) W() float64 {
	return 1
}

// Equals compares a point with another tuple
func (p Point) Equals(t Tuple) bool {
	if Equals(p.X(), t.X()) && Equals(p.Y(), t.Y()) && Equals(p.Z(), t.Y()) && Equals(p.W(), t.W()) {
		return true
	}
	return false
}

// Add adds a point to a vector
func (p Point) Add(t Tuple) Tuple {
	// t is a point, this is an error
	if t.W() == 1 {
		panic("cannot add two points")
	}
	return NewPoint(p.X()+t.X(), p.Y()+t.Y(), p.Z()+t.Z())

}

// Sub subtracts points or vectors
func (p Point) Sub(t Tuple) Tuple {

	// t is a vector, return a point
	if t.W() == 0 {
		return NewPoint(p.X()-t.X(), p.Y()-t.Y(), p.Z()-t.Z())
	}

	// t is a point, return a vector
	return NewVector(p.X()-t.X(), p.Y()-t.Y(), p.Z()-t.Z())
}

// Vector is a vector in 3D space. v[3] is always 0. Implements Tuple.
type Vector [4]float64

// NewVector returns a new Vector
func NewVector(x, y, z float64) Vector {

	return Vector{x, y, z, 0.0}
}

// X returns the vector's X coordinate
func (v Vector) X() float64 {
	return v[0]
}

// Y returns the vector's y coordinate
func (v Vector) Y() float64 {
	return v[1]
}

// Z returns the vector's Z coordinate
func (v Vector) Z() float64 {
	return v[2]
}

// W returns the vector's W coordinate
func (v Vector) W() float64 {
	return 0
}

// Equals compares a vector with another tuple
func (v Vector) Equals(t Tuple) bool {
	if Equals(v.X(), t.X()) && Equals(v.Y(), t.Y()) && Equals(v.Z(), t.Y()) && Equals(v.W(), t.W()) {
		return true
	}
	return false
}

// Add adds a point to a vector or a vector to a vector
func (v Vector) Add(t Tuple) Tuple {
	// t is a vector, return a vector
	if t.W() == 0 {
		return NewVector(v.X()+t.X(), v.Y()+t.Y(), v.Z()+t.Z())
	}

	// t is a point, return a point
	return NewPoint(v.X()+t.X(), v.Y()+t.Y(), v.Z()+t.Z())

}

// Sub subtracts vectors
func (v Vector) Sub(t Tuple) Tuple {

	// t is a point, this is an error
	if t.W() == 1 {
		panic("cannot subtract point from a vector")
	}

	// t is a point, return a vector
	return NewVector(v.X()-t.X(), v.Y()-t.Y(), v.Z()-t.Z())
}
