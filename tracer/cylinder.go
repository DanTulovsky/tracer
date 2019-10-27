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
	Shape
}

// NewDefaultCylinder returns a new cylinder
func NewDefaultCylinder() *Cylinder {
	return &Cylinder{
		Radius:  1.0,
		Minimum: math.Inf(-1),
		Maximum: math.Inf(1),
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

// IntersectWith returns the 't' values of Ray r intersecting with the Cylinder in sorted order
func (c *Cylinder) IntersectWith(r Ray) Intersections {
	t := Intersections{}

	// transform the ray by the inverse of the sphere transfrom matrix
	// instead of changing the sphere, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(c.transform.Inverse())

	// cylinder custom
	a := math.Pow(r.Dir.X(), 2) + math.Pow(r.Dir.Z(), 2)

	// ray is parallel to the y axis
	if a < constants.Epsilon {
		return t
	}

	b := 2*r.Origin.X()*r.Dir.X() + 2*r.Origin.Z()*r.Dir.Z()
	cc := math.Pow(r.Origin.X(), 2) + math.Pow(r.Origin.Z(), 2) - 1

	disc := math.Pow(b, 2) - 4*a*cc

	// ray does not intersect cylinder
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
	op := p.TimesMatrix(c.Transform().Inverse())

	// object normal, this is different for each shape
	on := NewVector(op.X(), 0, op.Z())

	// world normal
	wn := on.TimesMatrix(c.Transform().Submatrix(3, 3).Inverse().Transpose())

	return wn.Normalize()
}
