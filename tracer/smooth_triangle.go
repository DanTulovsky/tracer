package tracer

import "sort"

// SmoothTriangle is a triangle defined by 3 points in 3d space and the normals at those points
type SmoothTriangle struct {
	P1, P2, P3 Point

	// Normals at each point
	N1, N2, N3 Vector

	Triangle
}

// NewSmoothTriangle returns a new triangle
func NewSmoothTriangle(p1, p2, p3 Point, n1, n2, n3 Vector) *SmoothTriangle {
	t := &SmoothTriangle{
		N1: n1,
		N2: n2,
		N3: n3,

		Triangle: Triangle{
			P1: p1,
			P2: p2,
			P3: p3,
			E1: p2.SubPoint(p1),
			E2: p3.SubPoint(p1),
			Shape: Shape{
				transform:        IdentityMatrix(),
				transformInverse: IdentityMatrix().Inverse(),
				material:         NewDefaultMaterial(),
				shape:            "smooth-triangle",
			},
		},
	}

	return t
}

// IntersectWith returns the 't' value of Ray r intersecting with the triangle in sorted order
func (t *SmoothTriangle) IntersectWith(r Ray) Intersections {
	xs := NewIntersections()

	r = r.Transform(t.transformInverse)

	tval, u, v, err := t.sharedIntersectWith(r)
	if err != nil {
		return xs
	}

	xs = append(xs, NewIntersectionUV(t, tval, u, v))
	sort.Sort(byT(xs))
	return xs
}

// NormalAt returns the normal of the triangle at u,v stored in hit
func (t *SmoothTriangle) NormalAt(p Point, hit Intersection) Vector {
	v := t.N2.Scale(hit.u).AddVector(t.N3.Scale(hit.v)).AddVector(t.N1.Scale(1 - hit.u - hit.v))
	return v.NormalToWorldSpace(t)
}

// Includes implements includes logic
func (t *SmoothTriangle) Includes(s Shaper) bool {
	return t == s
}
