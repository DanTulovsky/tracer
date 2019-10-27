package tracer

import (
	"math"
	"sort"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/google/uuid"
)

// Cube implements an AABB cube
type Cube struct {
	Shape
}

// NewUnitCube returns a 1x1x1 cube
func NewUnitCube() *Cube {
	return &Cube{
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "cube",
			name:      uuid.New().String(),
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
		tmin = tminNumerator * math.Inf(1)
		tmax = tmaxNumerator * math.Inf(1)
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
	r = r.Transform(c.transform.Inverse())

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
func (c *Cube) NormalAt(p Point) Vector {

	// move point to object space
	op := p.TimesMatrix(c.Transform().Inverse())

	// object normal, this is different for each shape
	var on Vector
	maxc := math.Max(math.Max(math.Abs(op.X()), math.Abs(op.Y())), math.Abs(op.Z()))

	switch maxc {
	case math.Abs(op.X()):
		on = NewVector(op.X(), 0, 0)
	case math.Abs(op.Y()):
		on = NewVector(0, op.Y(), 0)
	default:
		on = NewVector(0, 0, op.Z())
	}

	// world normal
	wn := on.TimesMatrix(c.Transform().Submatrix(3, 3).Inverse().Transpose())

	return wn.Normalize()
}
