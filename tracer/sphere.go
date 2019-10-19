package tracer

import (
	"math"
	"sort"
)

// Sphere is a spherical object, implement Object
type Sphere struct {
	Center    Point
	Radius    float64
	transform Matrix
	material  Material
}

// NewUnitSphere returns a new Sphere centered at the origin with r=1
func NewUnitSphere() *Sphere {
	return &Sphere{Center: NewPoint(0, 0, 0),
		Radius:    1,
		transform: IdentityMatrix(),
		material:  NewDefaultMaterial(),
	}
}

// NewSphere returns a new Sphere
func NewSphere(c Point, r float64) *Sphere {
	return &Sphere{Center: c,
		Radius:    r,
		transform: IdentityMatrix(),
		material:  NewDefaultMaterial(),
	}
}

// IntersectWith returns the 't' values of Ray r intersecting with the Sphere in sorted order
func (s *Sphere) IntersectWith(r Ray) Intersections {

	t := Intersections{}

	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing

	r = r.Transform(s.transform.Inverse())

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

// Transform returns the transformation matrix of the Sphere
func (s *Sphere) Transform() Matrix {
	return s.transform
}

// SetTransform sets the transformation matrix of the Sphere
func (s *Sphere) SetTransform(m Matrix) {
	s.transform = m
}

// Material returns the material of the sphere
func (s *Sphere) Material() Material {
	return s.material
}

// SetMaterial sets the material of the sphere
func (s *Sphere) SetMaterial(m Material) {
	s.material = m
}

// NormalAt returns the normal vector at the given point on the surface of the sphere
func (s *Sphere) NormalAt(p Point) Vector {

	// move point to object space
	op := p.TimesMatrix(s.Transform().Inverse())
	// object normal
	on := op.SubPoint(Origin())

	// world normal
	wn := on.TimesMatrix(s.Transform().Submatrix(3, 3).Inverse().Transpose())

	return wn.Normalize()
}
