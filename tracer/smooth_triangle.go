package tracer

// SmoothTriangle is a triangle defined by 3 points in 3d space and the normals at those points
type SmoothTriangle struct {
	P1, P2, P3 Point

	// Normals at each point
	N1, N2, N3 Vector

	// Texture coordinates at each point
	VT1, VT2, VT3 Point

	Triangle
}

// NewSmoothTriangle returns a new triangle
func NewSmoothTriangle(p1, p2, p3 Point, n1, n2, n3 Vector, vt1, vt2, vt3 Point) *SmoothTriangle {
	t := &SmoothTriangle{
		N1:  n1,
		N2:  n2,
		N3:  n3,
		VT1: vt1,
		VT2: vt2,
		VT3: vt3,

		Triangle: Triangle{
			P1: p1,
			P2: p2,
			P3: p3,
			E1: p2.SubPoint(p1),
			E2: p3.SubPoint(p1),
			Shape: Shape{
				transform:        IM(),
				transformInverse: IM().Inverse(),
				material:         NewDefaultMaterial(),
				shape:            "smooth-triangle",
			},
		},
	}
	t.lna = t.localNormalAt
	t.calculateBounds()
	return t
}

// Equal returns true if the mooth triangles are equal
func (t *SmoothTriangle) Equal(t2 *SmoothTriangle) bool {
	return t.Shape.Equal(&t2.Shape) &&
		t.Triangle.Equal(&t2.Triangle) &&
		t.N1 == t2.N1 &&
		t.N2 == t2.N2 &&
		t.N3 == t2.N3
}

// IntersectWith returns the 't' value of Ray r intersecting with the triangle in sorted order
func (t *SmoothTriangle) IntersectWith(r Ray, xs Intersections) Intersections {
	r = r.Transform(t.transformInverse)

	tval, u, v, found := t.sharedIntersectWith(r)
	if !found {
		return xs
	}

	xs = append(xs, NewIntersectionUV(t, tval, u, v))
	return xs
}

func (t *SmoothTriangle) localNormalAt(unused Point, hit *Intersection) Vector {
	return t.N2.Scale(hit.u).AddVector(t.N3.Scale(hit.v)).AddVector(t.N1.Scale(1 - hit.u - hit.v))
}

// Includes implements includes logic
func (t *SmoothTriangle) Includes(s Shaper) bool {
	return t == s
}
