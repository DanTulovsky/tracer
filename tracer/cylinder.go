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
	return &Cylinder{
		Radius:  1.0,
		Minimum: -math.MaxFloat64,
		Maximum: math.MaxFloat64,
		Closed:  false,
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "cylinder",
		},
	}
}

// NewCylinder returns a new cylinder capped at min and max (y axis values, exclusive)
func NewCylinder(min, max float64) *Cylinder {
	if min > max {
		panic("cylinder min is greter than its max")
	}
	c := NewDefaultCylinder()
	c.Minimum = min
	c.Maximum = max
	return c
}

// NewClosedCylinder returns a new closed cylinder capped at min and max (y axis values, exclusive)
func NewClosedCylinder(min, max float64) *Cylinder {
	c := NewCylinder(min, max)
	c.Closed = true
	return c
}

// checkCap checks to see if the intersection at t is within the radius of the cylinder from the
// y axis
func (c *Cylinder) checkCap(r Ray, t float64) bool {
	x := r.Origin.X() + t*r.Dir.X()
	z := r.Origin.Z() + t*r.Dir.Z()

	return (math.Pow(x, 2) + math.Pow(z, 2)) <= c.Radius
}

func (c *Cylinder) intersectCaps(r Ray) Intersections {
	xs := NewIntersections()

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
func (c *Cylinder) IntersectWith(r Ray) Intersections {
	t := Intersections{}

	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(c.transform.Inverse())

	// check for intersections with the caps
	t = append(t, c.intersectCaps(r)...)

	// cylinder custom
	a := math.Pow(r.Dir.X(), 2) + math.Pow(r.Dir.Z(), 2)

	// ray is parallel to the y axis
	if a < constants.Epsilon {
		return t
	}

	b := 2*r.Origin.X()*r.Dir.X() + 2*r.Origin.Z()*r.Dir.Z()
	cc := math.Pow(r.Origin.X(), 2) + math.Pow(r.Origin.Z(), 2) - 1

	disc := math.Pow(b, 2) - 4*a*cc

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

// NormalAt returns the normal vector at the given point on the surface of the cylinder
func (c *Cylinder) NormalAt(p Point) Vector {
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
		on = NewVector(op.X(), 0, op.Z())
	}

	// world normal
	wn := on.NormalToWorldSpace(c)

	return wn.Normalize()
}

// calculateBounds calculates the bounding box of the shape
func (c *Cylinder) calculateBounds() {
	c.bound = Bound{
		Min: NewPoint(-1, c.Minimum, -1),
		Max: NewPoint(1, c.Maximum, 1),
	}
}

// PrecomputeValues precomputes some values for render speedup
func (c *Cylinder) PrecomputeValues() {
	// calculate group bounding box
	c.calculateBounds()
}
