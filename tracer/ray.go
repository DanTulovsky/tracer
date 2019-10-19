package tracer

// Ray describes a light ray
type Ray struct {
	Origin Point
	Dir    Vector
}

// NewRay returns a new ray
func NewRay(o Point, d Vector) Ray {
	return Ray{Origin: o, Dir: d}
}

// Position returns the position of the point, set at r.Origin, following this ray at time t
func (r Ray) Position(t float64) Point {
	return r.Origin.AddVector(r.Dir.Scale(t))
}

// Transform returns a new ray transformed by the matrix
func (r Ray) Transform(m Matrix) Ray {
	return Ray{Origin: r.Origin.TimesMatrix(m), Dir: r.Dir.TimesMatrix(m)}
}

// Equals returns true if rays are equal within Epsilon of each other
func (r Ray) Equals(s Ray) bool {
	if !r.Origin.Equals(s.Origin) {
		return false
	}

	if !r.Dir.Equals(s.Dir) {
		return false
	}

	return true
}
