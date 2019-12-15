package tracer

import (
	"math"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/rcrowley/go-metrics"
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
			transform:        IM(),
			transformInverse: IM().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "triangle",
		},
	}

	t.Normal = t.E2.Cross(t.E1).Normalize()
	t.lna = t.localNormalAt
	t.calculateBounds()
	return t
}

// Equal returns true if the triangles are equal
func (t *Triangle) Equal(t2 *Triangle) bool {
	return t.Shape.Equal(&t2.Shape) &&
		t.P1.Equal(t2.P1) &&
		t.P2.Equal(t2.P2) &&
		t.P3.Equal(t2.P3) &&
		t.E1.Equal(t2.E1) &&
		t.E2.Equal(t2.E2)
}

// sharedIntersectWith returns the tval, u, v, or an error if there is no intersection
// if no intersection is found, the last bool is false
func (t *Triangle) sharedIntersectWith(r Ray) (float64, float64, float64, bool) {

	dirCrossE2 := r.Dir.Cross(t.E2)
	d := t.E1.Dot(dirCrossE2)

	// if back culling is enabled, ignore back faces
	if t.WorldConfig().BackfaceCulling && d < 0 {
		metrics.GetOrRegisterCounter("num_backfaces_culled", nil).Inc(1)
		return 0, 0, 0, false
	}

	if math.Abs(d) < constants.Epsilon {
		// ray parallel to surface of triangle
		return 0, 0, 0, false
	}

	f := 1.0 / d
	p1ToOrigin := r.Origin.SubPoint(t.P1)
	u := f * p1ToOrigin.Dot(dirCrossE2)
	if u < 0 || u > 1 {
		// ray passes beyond p1-p3 edge
		return 0, 0, 0, false
	}

	// ray misses p1-p2 and p2-p3 edges
	oCrossE1 := p1ToOrigin.Cross(t.E1)
	v := f * r.Dir.Dot(oCrossE1)
	if v < 0 || (u+v) > 1 {
		return 0, 0, 0, false
	}

	return f * t.E2.Dot(oCrossE1), u, v, true
}

// IntersectWith returns the 't' value of Ray r intersecting with the triangle in sorted order
func (t *Triangle) IntersectWith(r Ray, xs Intersections) Intersections {
	r = r.Transform(t.transformInverse)

	// u, v not used here
	tval, _, _, found := t.sharedIntersectWith(r)
	if !found {
		return xs
	}
	xs = append(xs, NewIntersection(t, tval))
	return xs
}

func (t *Triangle) localNormalAt(unused Point, xs *Intersection) Vector {
	return t.Normal
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

	t.bound = NewBound(NewPoint(minX, minY, minZ), NewPoint(maxX, maxY, maxZ))
}

// PrecomputeValues precomputes some values for render speedup
func (t *Triangle) PrecomputeValues() {
	// calculate group bounding box
	// t.calculateBounds()
}

// Includes implements includes logic
func (t *Triangle) Includes(s Shaper) bool {
	return t == s
}
