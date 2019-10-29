package tracer

import (
	"math"
	"sort"

	"github.com/DanTulovsky/tracer/constants"
)

// Plane implment a plan in xz extending infinitely in both x and z dimensions
type Plane struct {
	Shape
}

// NewPlane returns a new plane
func NewPlane() *Plane {

	return &Plane{
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "plane",
			// name:      uuid.New().String(),
		},
	}
}

// NormalAt returns the normal vector at the given point on the surface of the plane
func (pl *Plane) NormalAt(p Point) Vector {
	on := NewVector(0, 1, 0)

	// common calculation to all shapes
	// world normal
	wn := on.NormalToWorldSpace(pl)

	return wn.Normalize()
}

// IntersectWith returns the 't' values of Ray r intersecting with the plane in sorted order
func (pl *Plane) IntersectWith(r Ray) Intersections {

	t := Intersections{}

	//  common calculation for all shapes
	r = r.Transform(pl.transform.Inverse())

	// parallel or coplanar
	if math.Abs(r.Dir.Y()) < constants.Epsilon {
		return t
	}

	t = append(t, NewIntersection(pl, -r.Origin.Y()/r.Dir.Y()))
	sort.Sort(byT(t))

	return t
}

// Bounds returns the untransformed bounding box
func (pl *Plane) Bounds() Bound {
	return Bound{
		Min: NewPoint(-math.MaxFloat64, -0.001, -math.MaxFloat64),
		Max: NewPoint(math.MaxFloat64, 0.001, math.MaxFloat64),
	}
}
