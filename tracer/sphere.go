package tracer

import (
	"math"
	"sort"
)

// Sphere is a spherical object, implement Object
type Sphere struct {
	Center Point
	Radius float64
}

// NewUnitSphere returns a new Sphere centered at the origin with r=1
func NewUnitSphere() Sphere {
	return Sphere{Center: NewPoint(0, 0, 0), Radius: 1}
}

// NewSphere returns a new Sphere
func NewSphere(c Point, r float64) Sphere {
	return Sphere{Center: c, Radius: r}
}

// IntersectWith returns the 't' values of Ray r intersecting with the Sphere in sorted order
func (s Sphere) IntersectWith(r Ray) Intersections {

	t := Intersections{}

	// vecto from sphere's center to ray origin
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
