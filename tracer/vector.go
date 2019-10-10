package tracer

import (
	"fmt"
	"math"
)

// Vector is a vector in 3D space. v[3] is always 0. Implements Tuple.
type Vector struct {
	Tuple
}

// NewVector returns a new Vector
func NewVector(x, y, z float64) Vector {

	return Vector{
		Tuple{x, y, z, 0.0},
	}
}

// X returns the vector's X coordinate
func (v Vector) X() float64 {
	return v.x
}

// Y returns the vector's y coordinate
func (v Vector) Y() float64 {
	return v.y
}

// Z returns the vector's Z coordinate
func (v Vector) Z() float64 {
	return v.z
}

// W returns the vector's W coordinate
func (v Vector) W() float64 {
	return v.w
}

// AddVector adds a vector  to a vector
func (v Vector) AddVector(t Vector) Vector {
	return NewVector(v.X()+t.X(), v.Y()+t.Y(), v.Z()+t.Z())
}

// AddPoint adds a point to a vector
func (v Vector) AddPoint(t Point) Point {
	return NewPoint(v.X()+t.X(), v.Y()+t.Y(), v.Z()+t.Z())
}

// SubVector subtracts vectors
func (v Vector) SubVector(t Vector) Vector {

	return NewVector(v.X()-t.X(), v.Y()-t.Y(), v.Z()-t.Z())
}

// Negate negates the vector (subtracts it from the zero vector)
func (v Vector) Negate() Vector {
	return NewVector(0, 0, 0).SubVector(v)
}

// Scale scales the vector
func (v Vector) Scale(s float64) Vector {
	return NewVector(v.X()*s, v.Y()*s, v.Z()*s)
}

// Magnitude computes the magnitude of the vector
func (v Vector) Magnitude() float64 {
	return math.Sqrt(math.Pow(v.X(), 2) + math.Pow(v.Y(), 2) + math.Pow(v.Z(), 2) + math.Pow(v.W(), 2))
}

// Normalize normalizes a vector to a unit vector
func (v Vector) Normalize() Vector {
	magnitude := v.Magnitude()
	if magnitude == 0 {
		panic(fmt.Sprintf("magnitude of %v is 0!", v))
	}
	return NewVector(v.X()/magnitude, v.Y()/magnitude, v.Z()/magnitude)
}

// Dot returns the dot product of the two vectors
// This is the cosine of the angle between two unit vectors
func (v Vector) Dot(w Vector) float64 {
	return v.X()*w.X() + v.Y()*w.Y() + v.Z()*w.Z() + v.W()*w.W()
}

// Cross returns the cross product of two vectors
// This returns a vector perpendicular to both of the original vectors
func (v Vector) Cross(w Vector) Vector {
	return NewVector(v.Y()*w.Z()-v.Z()*w.Y(), v.Z()*w.X()-v.X()*w.Z(), v.X()*w.Y()-v.Y()*w.X())
}

// TimesMatrix multiplies the vector by the matrix
func (v Vector) TimesMatrix(m Matrix) Vector {
	return NewVector(
		m[0][0]*v.X()+m[0][1]*v.Y()+m[0][2]*v.Z()+m[0][3]*v.W(),
		m[1][0]*v.X()+m[1][1]*v.Y()+m[1][2]*v.Z()+m[1][3]*v.W(),
		m[2][0]*v.X()+m[2][1]*v.Y()+m[2][2]*v.Z()+m[2][3]*v.W())
}
