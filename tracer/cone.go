package tracer

import (
	"math"
	"sort"

	"github.com/DanTulovsky/tracer/constants"
)

// Cone implements a double-napped cone
type Cone struct {

	// used to truncate the cone, min and max values on the y axis (-y, y)
	Minimum, Maximum float64
	Closed           bool // if true, cone is capped on both ends
	Shape
}

// NewDefaultCone returns a new default cone
func NewDefaultCone() *Cone {
	return &Cone{
		Minimum: -math.MaxFloat64,
		Maximum: math.MaxFloat64,
		Closed:  false,
		Shape: Shape{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "cone",
		},
	}
}

// NewCone returns a new cone capped at min and max (y axis values, exclusive)
func NewCone(min, max float64) *Cone {
	if min > max {
		panic("cone min is greter than its max")
	}
	c := NewDefaultCone()
	c.Minimum = min
	c.Maximum = max
	return c
}

// NewClosedCone returns a new closed cone capped at min and max (y axis values, exclusive)
func NewClosedCone(min, max float64) *Cone {
	c := NewCone(min, max)
	c.Closed = true
	return c
}

// checkCap checks to see if the intersection at t is within the radius of the cone from the
// y axis
func (c *Cone) checkCap(r Ray, t, y float64) bool {
	x := r.Origin.X() + t*r.Dir.X()
	z := r.Origin.Z() + t*r.Dir.Z()

	return (math.Pow(x, 2) + math.Pow(z, 2)) <= y
}

// intersectCaps returns intersections with the caps
func (c *Cone) intersectCaps(r Ray) Intersections {
	xs := NewIntersections()

	// caps only matter if the cone is closed and might possibly be intersected by the ray
	if !c.Closed || math.Abs(r.Dir.Y()) < constants.Epsilon {
		return xs
	}

	// check for an intersection with the lower end cap by intersecting the ray
	// with the plane at y=c.min
	t := (c.Minimum - r.Origin.Y()) / r.Dir.Y()
	if c.checkCap(r, t, math.Abs(c.Minimum)) {
		xs = append(xs, NewIntersection(c, t))
	}

	// check for an intersection with the upper end cap by intersecting the ray
	// with the plane at y=c.max
	t = (c.Maximum - r.Origin.Y()) / r.Dir.Y()
	if c.checkCap(r, t, math.Abs(c.Maximum)) {
		xs = append(xs, NewIntersection(c, t))
	}

	return xs
}

// IntersectWith returns the 't' values of Ray r intersecting with the cone in sorted order
func (c *Cone) IntersectWith(r Ray) Intersections {
	t := Intersections{}

	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(c.transformInverse)

	// cylinder custom

	// check for intersections with the caps
	t = append(t, c.intersectCaps(r)...)

	a := math.Pow(r.Dir.X(), 2) - math.Pow(r.Dir.Y(), 2) + math.Pow(r.Dir.Z(), 2)
	b := 2*r.Origin.X()*r.Dir.X() - 2*r.Origin.Y()*r.Dir.Y() + 2*r.Origin.Z()*r.Dir.Z()
	cc := math.Pow(r.Origin.X(), 2) - math.Pow(r.Origin.Y(), 2) + math.Pow(r.Origin.Z(), 2)

	// ray misses the cone
	if math.Abs(a) < constants.Epsilon {
		if math.Abs(b) < constants.Epsilon {
			// misses the cone
			return t
		}

		// single point of intersect
		t0 := -cc / (2 * b)
		t = append(t, NewIntersection(c, t0))
		return t
	}

	disc := math.Pow(b, 2) - 4*a*cc

	// ray does not intersect cone itself
	if disc < 0 {
		return t
	}

	t0 := (-b - math.Sqrt(disc)) / (2 * a)
	t1 := (-b + math.Sqrt(disc)) / (2 * a)

	if t0 > t1 {
		t0, t1 = t1, t0
	}

	y0 := r.Origin.Y() + t0*r.Dir.Y()
	if c.Minimum < y0 && y0 < c.Maximum {
		t = append(t, NewIntersection(c, t0))
	}

	y1 := r.Origin.Y() + t1*r.Dir.Y()
	if c.Minimum < y1 && y1 < c.Maximum {
		t = append(t, NewIntersection(c, t1))
	}

	sort.Sort(byT(t))
	return t
}

// NormalAt returns the normal vector at the given point on the surface of the cone
func (c *Cone) NormalAt(p Point, xs Intersection) Vector {
	// move point to object space
	op := p.ToObjectSpace(c)

	// object normal, this is different for each shape
	var on Vector

	// compute the square of the distance form the y-axis
	dist := math.Pow(op.X(), 2) + math.Pow(op.Z(), 2)

	switch {
	case dist < 1 && op.Y() >= c.Maximum-constants.Epsilon:
		on = NewVector(0, 1, 0)
	case dist < 1 && op.Y() <= c.Minimum+constants.Epsilon:
		on = NewVector(0, -1, 0)
	default:
		y := math.Sqrt(math.Pow(op.X(), 2) + math.Pow(op.Z(), 2))
		if op.Y() > 0 {
			y = -y
		}
		on = NewVector(op.X(), y, op.Z())
	}

	// world normal
	wn := on.NormalToWorldSpace(c)

	return wn.Normalize()
}

// calculateBounds calculates the bounding box of the shape
func (c *Cone) calculateBounds() {
	min := -math.Max(math.Abs(c.Maximum), math.Abs(c.Minimum))

	c.bound = Bound{
		Min: NewPoint(min, c.Minimum, min),
		Max: NewPoint(-min, c.Maximum, -min),
	}
}

// PrecomputeValues precomputes some values for render speedup
func (c *Cone) PrecomputeValues() {
	// calculate group bounding box
	c.calculateBounds()
}

// Includes implements includes logic
func (c *Cone) Includes(s Shaper) bool {
	return c == s
}
