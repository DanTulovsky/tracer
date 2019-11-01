package tracer

import (
	"math"

	"github.com/DanTulovsky/tracer/constants"
)

// Triangle is a triangle defined by 3 points in 3d space
type Triangle struct {
	P1, P2, P3 Point

	// edge1, edge2 and the normal
	E1, E2, Normal Vector

	Shape
}

// NewTriangle returns a new triangle
func NewTriangle(p1, p2, p3 Point) *Triangle {
	t := &Triangle{
		P1: p1,
		P2: p2,
		P3: p3,
		E1: p2.SubPoint(p1),
		E2: p3.SubPoint(p1),
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "triangle",
		},
	}

	t.Normal = t.E2.Cross(t.E1).Normalize()

	return t
}

// IntersectWith returns the 't' value of Ray r intersecting with the triangle in sorted order
func (t *Triangle) IntersectWith(r Ray) Intersections {
	xs := NewIntersections()

	dirCrossE2 := r.Dir.Cross(t.E2)
	d := t.E1.Dot(dirCrossE2)
	if math.Abs(d) < constants.Epsilon {
		// ray parallel to surface of triangle
		return xs
	}

	f := 1.0 / d
	p1ToOrigin := r.Origin.SubPoint(t.P1)
	u := f * p1ToOrigin.Dot(dirCrossE2)
	if u < 0 || u > 1 {
		// ray passes beyond p1-p3 edge
		return xs
	}

	// ray misses p1-p2 and p2-p3 edges
	oCrossE1 := p1ToOrigin.Cross(t.E1)
	v := f * r.Dir.Dot(oCrossE1)
	if v < 0 || (u+v) > 1 {
		return xs
	}

	tval := f * t.E2.Dot(oCrossE1)
	xs = append(xs, NewIntersection(t, tval))
	return xs
}

// NormalAt returns the normal of the triangle at the given point
func (t *Triangle) NormalAt(p Point) Vector {
	// world normal
	return t.Normal.NormalToWorldSpace(t)
}

// Bounds returns the untransformed bounding box
func (t *Triangle) Bounds() Bound {
	return t.bound
}

// calculateBounds calculates the bounding box of the shape
func (t *Triangle) calculateBounds() {
	minX := math.Min(t.P1.X(), math.Min(t.P2.X(), t.P3.X()))
	minY := math.Min(t.P1.Y(), math.Min(t.P2.Y(), t.P3.Y()))
	minZ := math.Min(t.P1.Z(), math.Min(t.P2.Z(), t.P3.Z()))

	maxX := math.Max(t.P1.X(), math.Max(t.P2.X(), t.P3.X()))
	maxY := math.Max(t.P1.Y(), math.Max(t.P2.Y(), t.P3.Y()))
	maxZ := math.Max(t.P1.Z(), math.Max(t.P2.Z(), t.P3.Z()))

	t.bound = Bound{
		Min: NewPoint(minX, minY, minZ),
		Max: NewPoint(maxX, maxY, maxZ),
	}
}

// PrecomputeValues precomputes some values for render speedup
func (t *Triangle) PrecomputeValues() {
	// calculate group bounding box
	t.calculateBounds()
}
