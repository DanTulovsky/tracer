package tracer

import (
	"math/rand"

	"github.com/DanTulovsky/tracer/utils"
)

// Shaper represents an physical object
type Shaper interface {
	Bounds() Bound

	HasMembers() bool

	Parent() Shaper
	HasParent() bool
	SetParent(Shaper)

	// If A is a Group, true if any child includes B
	// If A is a CSG, true if either child includes B
	// If A is anything else, true if A == B
	Includes(Shaper) bool

	IntersectWith(Ray, Intersections) Intersections
	NormalAt(Point, *Intersection) Vector
	PrecomputeValues()

	Material() *Material
	SetMaterial(*Material)

	WorldConfig() *WorldConfig
	SetWorldConfig(*WorldConfig)

	Name() string
	SetName(string)

	NumShapes() int

	// RandomPosition returns a andom point on the surface of the geometry
	RandomPosition(*rand.Rand) Point

	SetTransform(Matrix)
	Transform() Matrix

	TransformInverse() Matrix
}

// Shape is the abstract shape
type Shape struct {
	name             string
	shape            string
	transform        Matrix
	transformInverse Matrix
	material         *Material
	bound            Bound // cache the group bounding box
	wc               *WorldConfig

	parent Shaper

	// localNormalAt
	lna func(Point, *Intersection) Vector
}

// NumShapes returns the number of shapes contained in this object
func (s *Shape) NumShapes() int {
	return 1
}

// SetWorldConfig attachs the world config to this object
func (s *Shape) SetWorldConfig(wc *WorldConfig) {
	s.wc = wc
}

// WorldConfig returns the world config attached to this object
func (s *Shape) WorldConfig() *WorldConfig {
	return s.wc
}

// Equal returns true if the shapes are equal
func (s *Shape) Equal(s2 *Shape) bool {
	return s.shape == s2.shape &&
		s.name == s2.name &&
		s.transform.Equals(s2.Transform()) &&
		s.transformInverse.Equals(s2.Transform()) &&
		s.material.Equals(s2.material) &&
		// s.bound == s2.bound &&
		s.parent == s2.parent &&
		(s.lna != nil) == (s2.lna != nil)

}

// Includes implements includes logic
func (s *Shape) Includes(s2 Shaper) bool {
	panic("please implement Includes")
}

// HasMembers returns true if this is a group that has members
func (s *Shape) HasMembers() bool {
	return false
}

// PrecomputeValues precomputes some values for render speedup
func (s *Shape) PrecomputeValues() {
	// nothing by default, each shape can override
	panic("please implement PrecomputeValues")
}

// Parent returns the parent group this shape is part of
func (s *Shape) Parent() Shaper {
	return s.parent
}

// HasParent returns True if this shape has a parent
func (s *Shape) HasParent() bool {
	return s.parent != nil
}

// SetParent sets the parents of the object
func (s *Shape) SetParent(p Shaper) {
	s.parent = p
}

// IntersectWith implements Shaper interface
func (s *Shape) IntersectWith(r Ray, xs Intersections) Intersections {
	panic("must implement IntersectWith")
}

// NormalAt implements the Shaper interface
func (s *Shape) NormalAt(p Point, xs *Intersection) Vector {
	// move point to object space
	op := p.ToObjectSpace(s)

	// object normal, this is different for each shape
	on := s.lna(op, xs)

	// Apply any material perturbations to the normal
	on = s.Material().PerturbNormal(on, op)

	// world normal
	wn := on.NormalToWorldSpace(s)

	return wn.Normalize()
}

// localNormalAt returns the local normal vector at the point
func (s *Shape) localNormalAt(p Point, xs *Intersection) Vector {
	panic("must implement localNormalAt")
}

// Material returns the material of the shape
func (s *Shape) Material() *Material {
	return s.material
}

// SetMaterial sets the material of the shape
func (s *Shape) SetMaterial(m *Material) {
	s.material = m
}

// Transform returns the transformation matrix of the shape
func (s *Shape) Transform() Matrix {
	return s.transform
}

// SetTransform sets the transformation matrix of the shape
func (s *Shape) SetTransform(m Matrix) {
	s.transform = m
	s.transformInverse = m.Inverse()
}

// TransformInverse returns the inverse of the transformation matrix of the shape
func (s *Shape) TransformInverse() Matrix {
	return s.transformInverse
}

// Name returns the name of the shape
func (s *Shape) Name() string {
	return s.name
}

// SetName sets the name
func (s *Shape) SetName(n string) {
	s.name = n
}

// calculateBounds calculates the bounding box of the shape
func (s *Shape) calculateBounds() {
	panic("please implement calculateBounds!")
}

// Bounds returns the untransformed bounding box
func (s *Shape) Bounds() Bound {
	return s.bound
}

// RandomPosition returns a random point on the surface
func (s *Shape) RandomPosition(rng *rand.Rand) Point {
	minx := s.Bounds().Min.X()
	miny := s.Bounds().Min.Y()
	minz := s.Bounds().Min.Z()
	maxx := s.Bounds().Max.X()
	maxy := s.Bounds().Max.Y()
	maxz := s.Bounds().Max.Z()

	rx := utils.RandomFloat(rng, minx, maxx)
	ry := utils.RandomFloat(rng, miny, maxy)
	rz := utils.RandomFloat(rng, minz, maxz)

	p := NewPoint(rx, ry, rz).ToWorldSpace(s)
	return p
}

// Bound describes the bounding box for a shape
type Bound struct {
	Min, Max Point
	center   Point
}

// NewBound returns a new bounding box
func NewBound(min, max Point) Bound {
	b := Bound{
		Min: min,
		Max: max,
	}
	b.center = b.calculateCenter()

	return b
}

// calculateCenter calculates the center of the bounding box
func (b Bound) calculateCenter() Point {
	// TODO: Handle case when Max and Min are Inf and -Inf (adding them results in NaN)
	return NewPoint((b.Max.x+b.Min.x)/2, (b.Max.y+b.Min.y)/2, (b.Max.z+b.Min.z)/2)
}

// Center returns the center of the boundng box
func (b Bound) Center() Point {
	return b.center
}

// Equal returns true if the b == b2
func (b Bound) Equal(b2 Bound) bool {
	return b.Min == b2.Min && b.Max == b2.Max
}
