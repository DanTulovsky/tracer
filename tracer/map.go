package tracer

import (
	"math"
)

// Mapper is a texture map interface
type Mapper interface {
	// Map maps a 3D point to a 2D coordinate
	Map(Point) (float64, float64)
}

type cubeFace int

const (
	cubeFaceLeft cubeFace = iota
	cubeFaceRight
	cubeFaceFront
	cubeFaceBack
	cubeFaceUp
	cubeFaceDown
)

// CubeMap implements a cube map
type CubeMap struct {
	facepatterns map[cubeFace]UVPatterner
}

// NewCubeMapSame returns a new cube map that applies the same UVPattern to all faces
func NewCubeMapSame(p UVPatterner) *CubeMap {
	return &CubeMap{
		facepatterns: map[cubeFace]UVPatterner{
			cubeFaceLeft:  p,
			cubeFaceFront: p,
			cubeFaceRight: p,
			cubeFaceBack:  p,
			cubeFaceUp:    p,
			cubeFaceDown:  p,
		},
	}
}

// NewCubeMap returns a new cube map
func NewCubeMap(left, front, right, back, up, down UVPatterner) *CubeMap {
	return &CubeMap{
		facepatterns: map[cubeFace]UVPatterner{
			cubeFaceLeft:  left,
			cubeFaceFront: front,
			cubeFaceRight: right,
			cubeFaceBack:  back,
			cubeFaceUp:    up,
			cubeFaceDown:  down,
		},
	}
}

// faceFromPoint returns the face that the point is on
func (cm *CubeMap) faceFromPoint(p Point) cubeFace {
	coord := math.Max(math.Max(math.Abs(p.x), math.Abs(p.y)), math.Abs(p.z))

	switch coord {
	case p.x:
		return cubeFaceRight
	case -p.x:
		return cubeFaceLeft
	case p.y:
		return cubeFaceUp
	case -p.y:
		return cubeFaceDown
	case p.z:
		return cubeFaceFront
	}
	return cubeFaceBack
}

func (cm *CubeMap) uvFront(p Point) (float64, float64) {
	u := math.Mod((p.x+1), 2) / 2
	v := math.Mod((p.y+1), 2) / 2
	return u, 1 - v
}

func (cm *CubeMap) uvBack(p Point) (float64, float64) {
	u := math.Mod((1-p.x), 2) / 2
	v := math.Mod((p.y+1), 2) / 2
	return u, 1 - v
}
func (cm *CubeMap) uvLeft(p Point) (float64, float64) {
	u := math.Mod((p.z+1), 2) / 2
	v := math.Mod((p.y+1), 2) / 2
	return u, 1 - v
}
func (cm *CubeMap) uvRight(p Point) (float64, float64) {
	u := math.Mod((1-p.z), 2) / 2
	v := math.Mod((p.y+1), 2) / 2
	return u, 1 - v
}
func (cm *CubeMap) uvUp(p Point) (float64, float64) {
	u := math.Mod((p.x+1), 2) / 2
	v := math.Mod((1-p.z), 2) / 2
	return u, 1 - v
}
func (cm *CubeMap) uvDown(p Point) (float64, float64) {
	u := math.Mod((p.x+1), 2) / 2
	v := math.Mod((p.z+1), 2) / 2
	return u, 1 - v
}

// Map implements the Mapper interface
func (cm *CubeMap) Map(p Point) (float64, float64) {
	// first find the face the point is on
	face := cm.faceFromPoint(p)

	switch face {
	case cubeFaceFront:
		return cm.uvFront(p)
	case cubeFaceBack:
		return cm.uvBack(p)
	case cubeFaceLeft:
		return cm.uvLeft(p)
	case cubeFaceRight:
		return cm.uvRight(p)
	case cubeFaceUp:
		return cm.uvUp(p)
	case cubeFaceDown:
		return cm.uvDown(p)
	}

	return 0, 0
}

// PlaneMap implements a plane map
type PlaneMap struct{}

// NewPlaneMap returns a new plane map
func NewPlaneMap() *PlaneMap {
	return &PlaneMap{}
}

// Map implements the Mapper interface
func (pm *PlaneMap) Map(p Point) (float64, float64) {
	u := math.Mod(p.x, 1)
	if p.x < 0 {
		u = 1 + u
	}

	v := math.Mod(p.z, 1)
	if p.z < 0 {
		v = 1 + v
	}
	return u, v
}

// CylinderMap implements a cylinder map
type CylinderMap struct{}

// NewCylinderMap returns a new plane map
func NewCylinderMap() *CylinderMap {
	return &CylinderMap{}
}

// Map implements the Mapper interface
// TODO: Does not handle cap ends.
func (cm *CylinderMap) Map(p Point) (float64, float64) {
	theta := math.Atan2(p.x, p.z)
	rawU := theta / (2 * math.Pi)
	u := 1 - (rawU + 0.5)
	v := math.Mod(p.y, 1)
	if v < 0 {
		v = 1 + v
	}
	return u, v
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
	theta := math.Atan2(p.x, p.z)

	// vec is the vector pointing from the sphere's origin (the world origin)
	// to the point, which will also happen to be exactly equal to the sphere's
	// radius.
	vec := NewVector(p.x, p.y, p.z)
	radius := vec.Magnitude()

	// compute the polar angle
	// 0 <= phi <= π
	phi := math.Acos(p.y / radius)

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

	return 1 - u, 1 - v
}
