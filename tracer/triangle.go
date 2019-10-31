package tracer

import (
	"log"
	"math"
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

// NormalAt returns the normal of the triangle at the given point
func (t *Triangle) NormalAt(p Point) Vector {
	// world normal
	log.Println(t.Normal)
	return t.Normal.NormalToWorldSpace(t)
}

// Bounds returns the untransformed bounding box
func (t *Triangle) Bounds() Bound {
	minX := math.Min(t.P1.X(), math.Min(t.P2.X(), t.P3.X()))
	minY := math.Min(t.P1.Y(), math.Min(t.P2.Y(), t.P3.Y()))
	minZ := math.Min(t.P1.Z(), math.Min(t.P2.Z(), t.P3.Z()))

	maxX := math.Max(t.P1.X(), math.Max(t.P2.X(), t.P3.X()))
	maxY := math.Max(t.P1.Y(), math.Max(t.P2.Y(), t.P3.Y()))
	maxZ := math.Max(t.P1.Z(), math.Max(t.P2.Z(), t.P3.Z()))

	return Bound{
		Min: NewPoint(minX, minY, minZ),
		Max: NewPoint(maxX, maxY, maxZ),
	}
}
