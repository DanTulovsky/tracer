package tracer

import (
	"math"
	"sort"

	"github.com/DanTulovsky/tracer/constants"
)

// Cylinder implements a cylinder of radius 1
type Cylinder struct {
	Radius float64

	// used to truncate the cylinder, min and max values on the y axis (-y, y)
	Minimum, Maximum float64
	Closed           bool // if true, cylinder is capped on both ends
	Shape
}

// NewDefaultCylinder returns a new cylinder
func NewDefaultCylinder() *Cylinder {
	c := &Cylinder{
		Radius:  1.0,
		Minimum: -math.MaxFloat64,
		Maximum: math.MaxFloat64,
		Closed:  false,
		Shape: Shape{
			transform:        IM(),
			transformInverse: IM().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "cylinder",
		},
	}
	c.lna = c.localNormalAt
	c.calculateBounds()
	return c
}

// NewCylinder returns a new cylinder capped at min and max (y axis values, exclusive)
func NewCylinder(min, max float64) *Cylinder {
	if min > max {
		panic("cylinder min is greter than its max")
	}
	c := NewDefaultCylinder()
	c.Minimum = min
	c.Maximum = max
	c.calculateBounds()
	return c
}

// NewClosedCylinder returns a new closed cylinder capped at min and max (y axis values, exclusive)
func NewClosedCylinder(min, max float64) *Cylinder {
	c := NewCylinder(min, max)
	c.Closed = true
	return c
}

// Equal returns true if the cylinders are equal
func (c *Cylinder) Equal(c2 *Cylinder) bool {
	return c.Shape.Equal(&c2.Shape) &&
		c.Radius == c2.Radius &&
		c.Minimum == c2.Minimum &&
		c.Closed == c2.Closed
}

// checkCap checks to see if the intersection at t is within the radius of the cylinder from the
// y axis
func (c *Cylinder) checkCap(r Ray, t float64) bool {
	x := r.Origin.X() + t*r.Dir.X()
	z := r.Origin.Z() + t*r.Dir.Z()

	return (x*x + z*z) <= c.Radius
}

func (c *Cylinder) intersectCaps(r Ray, xs Intersections) Intersections {
	// caps only matter if the cylinder is closed and might possibly be intersected by the ray
	if !c.Closed || math.Abs(r.Dir.Y()) < constants.Epsilon {
		return xs
	}

	// check for an intersection with the lower end cap by intersecting the ray
	// with the plan a y=c.min
	t := (c.Minimum - r.Origin.Y()) / r.Dir.Y()
	if c.checkCap(r, t) {
		xs = append(xs, NewIntersection(c, t))
	}

	// check for an intersection with the upper end cap by intersecting the ray
	// with the plan a y=c.max
	t = (c.Maximum - r.Origin.Y()) / r.Dir.Y()
	if c.checkCap(r, t) {
		xs = append(xs, NewIntersection(c, t))
	}

	return xs
}

// IntersectWith returns the 't' values of Ray r intersecting with the Cylinder in sorted order
func (c *Cylinder) IntersectWith(r Ray, t Intersections) Intersections {
	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(c.transformInverse)

	// check for intersections with the caps
	t = c.intersectCaps(r, t)

	// cylinder custom
	a := r.Dir.X()*r.Dir.X() + r.Dir.Z()*r.Dir.Z()

	// ray is parallel to the y axis
	if a < constants.Epsilon {
		return t
	}

	b := 2*r.Origin.X()*r.Dir.X() + 2*r.Origin.Z()*r.Dir.Z()
	cc := r.Origin.X()*r.Origin.X() + r.Origin.Z()*r.Origin.Z() - 1

	disc := b*b - 4*a*cc

	// ray does not intersect cylinder itself
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

func (c *Cylinder) localNormalAt(p Point, xs *Intersection) Vector {
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
		on = NewVector(p.X(), 0, p.Z())
	}
	return on
}

// calculateBounds calculates the bounding box of the shape
func (c *Cylinder) calculateBounds() {
	c.bound = NewBound(NewPoint(-1, c.Minimum, -1), NewPoint(1, c.Maximum, 1))
}

// PrecomputeValues precomputes some values for render speedup
func (c *Cylinder) PrecomputeValues() {
	// calculate group bounding box
	// c.calculateBounds()
}

// Includes implements includes logic
func (c *Cylinder) Includes(s Shaper) bool {
	return c == s
}
