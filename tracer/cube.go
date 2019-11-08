package tracer

import (
	"math"
	"sort"

	"github.com/DanTulovsky/tracer/constants"
)

// Cube implements an AABB cube
type Cube struct {
	Shape
}

// NewUnitCube returns a 1x1x1 cube
func NewUnitCube() *Cube {
	return &Cube{
		Shape: Shape{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "cube",
			// name:      uuid.New().String(),
		},
	}
}

// checkAxis is a helper function for check for intersection of the cube and ray
func (c *Cube) checkAxis(o, d float64) (float64, float64) {

	var tmin, tmax float64

	tminNumerator := -1 - o
	tmaxNumerator := 1 - o

	if math.Abs(d) >= constants.Epsilon {
		tmin = tminNumerator / d
		tmax = tmaxNumerator / d
	} else {
		tmin = tminNumerator * math.MaxFloat64
		tmax = tmaxNumerator * math.MaxFloat64
	}

	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

	return tmin, tmax
}

// IntersectWith returns the 't' values of Ray r intersecting with the Cube in sorted order
func (c *Cube) IntersectWith(r Ray) Intersections {

	t := Intersections{}

	// common to all shapes
	r = r.Transform(c.transformInverse)

	// Cube specific
	var tmin, tmax float64

	xtmin, xtmax := c.checkAxis(r.Origin.X(), r.Dir.X())
	ytmin, ytmax := c.checkAxis(r.Origin.Y(), r.Dir.Y())
	ztmin, ztmax := c.checkAxis(r.Origin.Z(), r.Dir.Z())

	tmin = math.Max(math.Max(xtmin, ytmin), ztmin)
	tmax = math.Min(math.Min(xtmax, ytmax), ztmax)

	// missed the cube
	if tmin > tmax {
		return t
	}

	t = append(t, NewIntersection(c, tmin))
	t = append(t, NewIntersection(c, tmax))

	sort.Sort(byT(t))

	return t
}

// NormalAt returns the normal vector at the given point on the surface of the cube
func (c *Cube) NormalAt(p Point, xs Intersection) Vector {

	// move point to object space
	op := p.ToObjectSpace(c)

	// object normal, this is different for each shape
	on := c.localNormalAt(op, xs)

	// world normal
	wn := on.NormalToWorldSpace(c)

	return wn.Normalize()
}

func (c *Cube) localNormalAt(p Point, xs Intersection) Vector {
	var on Vector
	maxc := math.Max(math.Max(math.Abs(p.X()), math.Abs(p.Y())), math.Abs(p.Z()))

	switch maxc {
	case math.Abs(p.X()):
		on = NewVector(p.X(), 0, 0)
	case math.Abs(p.Y()):
		on = NewVector(0, p.Y(), 0)
	default:
		on = NewVector(0, 0, p.Z())
	}
	return on
}

// calculateBounds calculates the bounding box of the shape
func (c *Cube) calculateBounds() {
	c.bound = Bound{
		Min: NewPoint(-1, -1, -1),
		Max: NewPoint(1, 1, 1),
	}
}

// PrecomputeValues precomputes some values for render speedup
func (c *Cube) PrecomputeValues() {
	// calculate group bounding box
	c.calculateBounds()
}

// Includes implements includes logic
func (c *Cube) Includes(s Shaper) bool {
	return c == s
}
