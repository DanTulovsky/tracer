package tracer

import (
	"fmt"
	"math"

	"github.com/DanTulovsky/tracer/utils"
)

// Vector is a vector in 3D space. v[3] is always 0. Implements Tuple.
type Vector struct {
	x, y, z, w float64
}

// NewVector returns a new Vector
func NewVector(x, y, z float64) Vector {

	return Vector{x, y, z, 0.0}
}

// X returns the point's X coordinate
func (v Vector) X() float64 {
	return v.x
}

// Y returns the point's y coordinate
func (v Vector) Y() float64 {
	return v.y
}

// Z returns the point's Z coordinate
func (v Vector) Z() float64 {
	return v.z
}

// W returns the point's W coordinate
func (v Vector) W() float64 {
	return v.w
}

// SetX sets x
func (v Vector) SetX(a float64) {
	v.x = a
}

// SetY sets y
func (v Vector) SetY(a float64) {
	v.y = a
}

// SetZ sets y
func (v Vector) SetZ(a float64) {
	v.z = a
}

// SetW sets w
func (v Vector) SetW(a float64) {
	v.w = a
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
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z + v.w*v.w)
}

// Normalize normalizes a vector to a unit vector
func (v Vector) Normalize() Vector {
	magnitude := v.Magnitude()
	if magnitude == 0 {
		return NewVector(0, 0, 0)
	}
	m := 1 / magnitude // optimization
	return NewVector(v.x*m, v.y*m, v.z*m)
}

// MagnitudeNormalize returns both the magnitude and the normalized vector
func (v Vector) MagnitudeNormalize() (float64, Vector) {
	magnitude := v.Magnitude()
	if magnitude == 0 {
		return magnitude, NewVector(0, 0, 0)
	}
	m := 1 / magnitude // optimization
	return magnitude, NewVector(v.x*m, v.y*m, v.z*m)
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

	r, _ := m.Dims()
	var vec Vector

	// Chap 6, page 82 claims Sphere.NormalAt requires this and SubMatrix
	// but v.W() is always 0, so it should never matter
	switch r {
	case 3:
		vec = NewVector(
			m[0][0]*v.X()+m[0][1]*v.Y()+m[0][2]*v.Z(),
			m[1][0]*v.X()+m[1][1]*v.Y()+m[1][2]*v.Z(),
			m[2][0]*v.X()+m[2][1]*v.Y()+m[2][2]*v.Z())
	case 4:
		vec = NewVector(
			m[0][0]*v.X()+m[0][1]*v.Y()+m[0][2]*v.Z()+m[0][3]*v.W(),
			m[1][0]*v.X()+m[1][1]*v.Y()+m[1][2]*v.Z()+m[1][3]*v.W(),
			m[2][0]*v.X()+m[2][1]*v.Y()+m[2][2]*v.Z()+m[2][3]*v.W())
	default:
		panic("only 3x3 and 4x4 matricies supported")
	}

	return vec
}

// String returns ...
func (v Vector) String() string {
	return fmt.Sprintf("Vector(%.2f, %.2f, %.2f)", v.X(), v.Y(), v.Z())
}

// Equal compares vectors
func (v Vector) Equal(s Vector) bool {
	if utils.Equals(v.X(), s.X()) && utils.Equals(v.Y(), s.Y()) && utils.Equals(v.Z(), s.Z()) && utils.Equals(v.W(), s.W()) {
		return true
	}
	return false
}

// Reflect returns the vector reflected around another one
func (v Vector) Reflect(n Vector) Vector {
	return v.SubVector(n.Scale(2).Scale(v.Dot(n)))
}

// NormalToWorldSpace converts the given vector from object space to world space
func (v Vector) NormalToWorldSpace(s Shaper) Vector {
	n := v.TimesMatrix(s.TransformInverse().Transpose())
	n.SetW(0)
	n = n.Normalize()

	if s.HasParent() {
		n = n.NormalToWorldSpace(s.Parent())
	}

	return n
}
