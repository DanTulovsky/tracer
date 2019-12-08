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
	c := &Cone{
		Minimum: -math.MaxFloat64,
		Maximum: math.MaxFloat64,
		Closed:  false,
		Shape: Shape{
			transform:        IM(),
			transformInverse: IM().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "cone",
		},
	}
	c.lna = c.localNormalAt
	c.calculateBounds()
	return c
}

// NewCone returns a new cone capped at min and max (y axis values, exclusive)
func NewCone(min, max float64) *Cone {
	if min > max {
		panic("cone min is greter than its max")
	}
	c := NewDefaultCone()
	c.Minimum = min
	c.Maximum = max
	c.calculateBounds()
	return c
}

// NewClosedCone returns a new closed cone capped at min and max (y axis values, exclusive)
func NewClosedCone(min, max float64) *Cone {
	c := NewCone(min, max)
	c.Closed = true
	return c
}

// Equal returns true if the cones are equal
func (c *Cone) Equal(c2 *Cone) bool {
	return c.Shape.Equal(&c2.Shape) &&
		c.Minimum == c2.Minimum &&
		c.Maximum == c2.Maximum &&
		c.Closed == c2.Closed
}

// checkCap checks to see if the intersection at t is within the radius of the cone from the
// y axis
func (c *Cone) checkCap(r Ray, t, y float64) bool {
	x := r.Origin.X() + t*r.Dir.X()
	z := r.Origin.Z() + t*r.Dir.Z()

	return (x*x + z*z) <= y
}

// intersectCaps returns intersections with the caps
func (c *Cone) intersectCaps(r Ray, xs Intersections) Intersections {
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
func (c *Cone) IntersectWith(r Ray, t Intersections) Intersections {
	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(c.transformInverse)

	// check for intersections with the caps
	t = c.intersectCaps(r, t)

	a := r.Dir.X()*r.Dir.X() - r.Dir.Y()*r.Dir.Y() + r.Dir.Z()*r.Dir.Z()
	b := 2*r.Origin.X()*r.Dir.X() - 2*r.Origin.Y()*r.Dir.Y() + 2*r.Origin.Z()*r.Dir.Z()
	cc := r.Origin.X()*r.Origin.X() - r.Origin.Y()*r.Origin.Y() + r.Origin.Z()*r.Origin.Z()

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

	disc := b*b - 4*a*cc

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

func (c *Cone) localNormalAt(p Point, xs *Intersection) Vector {
	// object normal, this is different for each shape
	var on Vector

	// compute the square of the distance form the y-axis
	dist := p.X()*p.X() + p.Z()*p.Z()

	switch {
	case dist < 1 && p.Y() >= c.Maximum-constants.Epsilon:
		on = NewVector(0, 1, 0)
	case dist < 1 && p.Y() <= c.Minimum+constants.Epsilon:
		on = NewVector(0, -1, 0)
	default:
		y := math.Sqrt(p.X()*p.X() + p.Z()*p.Z())
		if p.Y() > 0 {
			y = -y
		}
		on = NewVector(p.X(), y, p.Z())
	}
	return on
}

// calculateBounds calculates the bounding box of the shape
func (c *Cone) calculateBounds() {
	min := -math.Max(math.Abs(c.Maximum), math.Abs(c.Minimum))

	c.bound = NewBound(
		NewPoint(min, c.Minimum, min),
		NewPoint(-min, c.Maximum, -min))
}

// PrecomputeValues precomputes some values for render speedup
func (c *Cone) PrecomputeValues() {
	// calculate group bounding box
	// c.calculateBounds()
}

// Includes implements includes logic
func (c *Cone) Includes(s Shaper) bool {
	return c == s
}
