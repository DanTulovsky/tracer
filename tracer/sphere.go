package tracer

import (
	"math"
	"sort"
)

// Sphere is a spherical object, implements Shaper
type Sphere struct {
	Center Point
	Radius float64
	Shape
}

// NewUnitSphere returns a new Sphere centered at the origin with r=1
func NewUnitSphere() *Sphere {
	return &Sphere{
		Center: NewPoint(0, 0, 0),
		Radius: 1,
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "sphere",
			// name:      uuid.New().String(),
		},
	}
}

// NewGlassSphere returns a new Sphere centered at the origin with r=1, with a transparent material
func NewGlassSphere() *Sphere {
	return &Sphere{
		Center: NewPoint(0, 0, 0),
		Radius: 1,
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultGlassMaterial(),
		},
	}
}

// IntersectWith returns the 't' values of Ray r intersecting with the Sphere in sorted order
func (s *Sphere) IntersectWith(r Ray) Intersections {

	t := Intersections{}

	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(s.transform.Inverse())

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

// NormalAt returns the normal vector at the given point on the surface of the sphere
func (s *Sphere) NormalAt(p Point) Vector {

	// move point to object space
	op := p.ToObjectSpace(s)

	// object normal, this is different for each shape
	on := op.SubPoint(Origin())

	// world normal
	wn := on.NormalToWorldSpace(s)

	return wn.Normalize()

}
