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
	}
	csg.left.SetParent(csg)
	csg.right.SetParent(csg)
	return csg
}
