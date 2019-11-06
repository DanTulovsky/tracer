package tracer

// Operation defines the operation applied to a CSG
type Operation int

const (
	// Union combines the inputs into a single shape, prseving all external surfaces
	Union Operation = iota

	// Intersect preserves the portion of th einputs that share a volume
	Intersect

	// Difference preserves only the portion of the first shape where it's overlapped by the others
	Difference
)

// CSG implements Constructive Solid Geometry shape
// Two shapes combined together via union, intersection or difference
type CSG struct {
	left, right Shaper
	op          Operation

	Shape
}

// NewCSG returns a new CSG
func NewCSG(s1, s2 Shaper, op Operation) *CSG {
	csg := &CSG{
		left:  s1,
		right: s2,
		op:    op,
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "csg",
		},
	}
	csg.left.SetParent(csg)
	csg.right.SetParent(csg)
	return csg
}

// IntersectionAllowed returns true if this is a valid intersection
func (csg *CSG) IntersectionAllowed(op Operation, lhit, inl, inr bool) bool {

	switch op {
	case Union:
		return (lhit && !inr) || (!lhit && !inl)
	case Intersect:
		return (lhit && inr) || (!lhit && inl)
	case Difference:
		return (lhit && !inr) || (!lhit && inl)
	}
	return false
}

// FilterIntersections takes a list of intersections of two shapes and returns only those valid for the current CSG
func (csg *CSG) FilterIntersections(xs Intersections) Intersections {
	result := NewIntersections()

	inl, inr := false, false

	for _, x := range xs {
		lhit := csg.left.Includes(x.Object())

		if csg.IntersectionAllowed(csg.op, lhit, inl, inr) {
			result = append(result, x)
		}

		if lhit {
			inl = !inl
		} else {
			inr = !inr
		}
	}

	return result
}

// Includes implements includes logic
func (csg *CSG) Includes(s2 Shaper) bool {
	return (csg.left.Includes(s2) || csg.right.Includes(s2))
}
