package tracer

import (
	"fmt"

	"github.com/DanTulovsky/tracer/utils"
)

// Point is a single point in 3D space. p[3] is always 1. Implements Tupler.
type Point struct {
	x, y, z, w float64
}

// Origin returns the point at the origin
func Origin() Point {
	return Point{0, 0, 0, 1.0}
}

// NewPoint returns a new Point
func NewPoint(x, y, z float64) Point {
	return Point{x, y, z, 1.0}
}

// X returns the point's X coordinate
func (p Point) X() float64 {
	return p.x
}

// Y returns the point's y coordinate
func (p Point) Y() float64 {
	return p.y
}

// Z returns the point's Z coordinate
func (p Point) Z() float64 {
	return p.z
}

// W returns the point's W coordinate
func (p Point) W() float64 {
	return p.w
}

// SetW sets w
func (p Point) SetW(a float64) {
	p.w = a
}

// AddScalar adds a scalar to all the values of the point
func (p Point) AddScalar(t float64) Point {
	return NewPoint(p.X()+t, p.Y()+t, p.Z()+t)
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

// String returns ...
func (p Point) String() string {
	return fmt.Sprintf("Point(%.2f, %.2f, %.2f)", p.X(), p.Y(), p.Z())
}

// Equal compares points
func (p Point) Equal(s Point) bool {
	if utils.Equals(p.X(), s.X()) && utils.Equals(p.Y(), s.Y()) && utils.Equals(p.Z(), s.Z()) && utils.Equals(p.W(), s.W()) {
		return true
	}
	return false
}

// ToObjectSpace converts the given point from world space to object space
func (p Point) ToObjectSpace(s Shaper) Point {
	res := p

	if s.HasParent() {
		res = p.ToObjectSpace(s.Parent())
	}

	return res.TimesMatrix(s.TransformInverse())
}

// ToWorldSpace converts the given point from object space to world space
func (p Point) ToWorldSpace(s Shaper) Point {
	res := p

	if s.HasParent() {
		res = p.ToWorldSpace(s.Parent())
	}

	res = res.TimesMatrix(s.Transform())
	res.SetW(1)

	return res
}
