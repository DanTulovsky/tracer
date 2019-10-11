package tracer

// Sphere is a spherical object, implement Object
type Sphere struct {
	Center Point
	Radius float64
}

// NewSphere returns a new Sphere
func NewSphere(c Point, r float64) Sphere {
	return Sphere{Center: c, Radius: r}
}

// IntersectWith returns the 't' values of Ray r intersecting with the Sphere in sorted order
func (s Sphere) IntersectWith(r Ray) []float64 {
	return []float64{0}
}
