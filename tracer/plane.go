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
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "plane",
			// name:      uuid.New().String(),
		},
	}
}

// NormalAt returns the normal vector at the given point on the surface of the plane
func (pl *Plane) NormalAt(p Point, xs Intersection) Vector {
	// move point to object space
	op := p.ToObjectSpace(pl)

	// object normal, this is different for each shape
	on := pl.localNormalAt(op, xs)

	// world normal
	wn := on.NormalToWorldSpace(pl)

	return wn.Normalize()
}

func (pl *Plane) localNormalAt(p Point, xs Intersection) Vector {
	return NewVector(0, 1, 0)
}

// IntersectWith returns the 't' values of Ray r intersecting with the plane in sorted order
func (pl *Plane) IntersectWith(r Ray) Intersections {

	t := Intersections{}

	//  common calculation for all shapes
	r = r.Transform(pl.transformInverse)

	// parallel or coplanar
	if math.Abs(r.Dir.Y()) < constants.Epsilon {
		return t
	}

	t = append(t, NewIntersection(pl, -r.Origin.Y()/r.Dir.Y()))
	sort.Sort(byT(t))

	return t
}

// calculateBounds calculates the bounding box of the shape
func (pl *Plane) calculateBounds() {
	pl.bound = Bound{
		Min: NewPoint(-math.MaxFloat64, -0.001, -math.MaxFloat64),
		Max: NewPoint(math.MaxFloat64, 0.001, math.MaxFloat64),
	}
}

// PrecomputeValues precomputes some values for render speedup
func (pl *Plane) PrecomputeValues() {
	// calculate group bounding box
	pl.calculateBounds()
}

// Includes implements includes logic
func (pl *Plane) Includes(s Shaper) bool {
	return pl == s
}
