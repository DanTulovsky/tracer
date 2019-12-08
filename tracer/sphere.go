package tracer

import (
	"math"
	"sort"

	"github.com/DanTulovsky/tracer/utils"
)

// Sphere is a spherical object, implements Shaper
type Sphere struct {
	Center Point
	Radius float64
	Shape
}

// NewUnitSphere returns a new Sphere centered at the origin with r=1
func NewUnitSphere() *Sphere {
	s := &Sphere{
		Center: NewPoint(0, 0, 0),
		Radius: 1,
		Shape: Shape{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "sphere",
		},
	}
	s.lna = s.localNormalAt
	s.calculateBounds()
	return s
}

// Equal returns true if the spheres are equal
func (s *Sphere) Equal(s2 *Sphere) bool {
	return s.Shape.Equal(&s2.Shape) &&
		s.Center.Equal(s2.Center) &&
		s.Radius == s2.Radius
}

// NewGlassSphere returns a new Sphere centered at the origin with r=1, with a transparent material
func NewGlassSphere() *Sphere {
	s := NewUnitSphere()
	s.SetMaterial(NewDefaultGlassMaterial())
	return s
}

// IntersectWith returns the 't' values of Ray r intersecting with the Sphere in sorted order
func (s *Sphere) IntersectWith(r Ray, t Intersections) Intersections {
	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(s.transformInverse)

	// vector from sphere's center to ray origin
	sphereToRay := r.Origin.SubPoint(s.Center)

	a := r.Dir.Dot(r.Dir)
	b := 2 * r.Dir.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1

	// discriminant
	d := math.Pow(b, 2) - 4*a*c

	// no intersection
	if d < 0 {
		return t
	}

	// one intersection means ray hits at tangent
	t = append(t, NewIntersection(s, (-b-math.Sqrt(d))/(2*a)))
	t = append(t, NewIntersection(s, (-b+math.Sqrt(d))/(2*a)))

	sort.Sort(byT(t))

	return t
}

func (s *Sphere) localNormalAt(p Point, xs *Intersection) Vector {
	return p.SubPoint(Origin())
}

// calculateBounds calculates the bounding box of the shape
func (s *Sphere) calculateBounds() {
	s.bound = NewBound(NewPoint(-1, -1, -1), NewPoint(1, 1, 1))
}

// PrecomputeValues precomputes some values for render speedup
func (s *Sphere) PrecomputeValues() {
	// calculate group bounding box
	// s.calculateBounds()
}

// Includes implements includes logic
func (s *Sphere) Includes(s2 Shaper) bool {
	return s == s2
}

// RandomPosition returns a random point on the surface of the sphere
// http://mathworld.wolfram.com/SpherePointPicking.html
func (s *Sphere) RandomPosition() Point {

	u := utils.RandomFloat(0, 1)
	v := utils.RandomFloat(0, 1)

	theta := math.Pi * 2 * u
	phi := math.Acos(2*v - 1)

	x := s.Radius * math.Sin(phi) * math.Cos(theta)
	y := s.Radius * math.Sin(phi) * math.Sin(theta)
	z := s.Radius * math.Cos(phi)

	return NewPoint(x, y, z).ToWorldSpace(s)
}
