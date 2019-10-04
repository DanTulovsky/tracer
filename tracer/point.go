package tracer

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
