package tracer

import "math"

// Mapper is a texture map interface
type Mapper interface {
	// Map maps a 3D point to a 2D coordinate
	Map(Point) (float64, float64)
}

// SphericalMap implements a spherical map
type SphericalMap struct {
}

// NewSphericalMap returns a new spherical map
func NewSphericalMap() *SphericalMap {
	return &SphericalMap{}
}

// Map implements the Mapper interface
func (sm *SphericalMap) Map(p Point) (float64, float64) {
	// p is assume to lie on the surface of a sphere centered at the origin

	// compute the azimuthal angle
	// -π < theta <= π
	// angle increases clockwise as viewed from above,
	// which is opposite of what we want, but we'll fix it later.
	theta := math.Atan2(p.X(), p.Z())

	// vec is the vector pointing from the sphere's origin (the world origin)
	// to the point, which will also happen to be exactly equal to the sphere's
	// radius.
	vec := NewVector(p.X(), p.Y(), p.Z())
	radius := vec.Magnitude()

	// compute the polar angle
	// 0 <= phi <= π
	phi := math.Acos(p.Y() / radius)

	// -0.5 < raw_u <= 0.5
	rawU := theta / (2 * math.Pi)

	// 0 <= u < 1
	// here's also where we fix the direction of u. Subtract it from 1,
	// so that it increases counterclockwise as viewed from above.
	u := 1 - (rawU + 0.5)

	// we want v to be 0 at the south pole of the sphere,
	// and 1 at the north pole, so we have to "flip it over"
	// by subtracting it from 1.
	v := 1 - phi/math.Pi

	return u, v
}
