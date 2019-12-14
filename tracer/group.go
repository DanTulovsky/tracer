package tracer

import (
	"math"
	"sort"

	"github.com/google/go-cmp/cmp"

	"github.com/DanTulovsky/tracer/constants"
)

// Group is a collection of other groups/objects
type Group struct {
	members []Shaper
	Shape
}

// NewGroup returns a new, empty group
func NewGroup() *Group {
	g := &Group{
		members: []Shaper{},
		Shape: Shape{
			transform:        IM(),
			transformInverse: IM().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "group",
		},
	}
	g.calculateBounds()
	return g
}

// Equal returns true if the groups are equal
func (g *Group) Equal(g2 *Group) bool {
	return g.Shape.Equal(&g2.Shape) &&
		cmp.Equal(g.members, g2.members)
}

// AddMember adds a new member to this group
func (g *Group) AddMember(m Shaper) {
	g.members = append(g.members, m)
	m.SetParent(g)
	g.calculateBounds()
}

// AddMembers adds a new member to this group
func (g *Group) AddMembers(m ...Shaper) {
	for _, mem := range m {

		g.members = append(g.members, mem)
		mem.SetParent(g)
	}
	g.calculateBounds()
}

// Members returns all the direct members of this group
func (g *Group) Members() []Shaper {
	return g.members
}

// HasMembers returns true if this is a group that has members
func (g *Group) HasMembers() bool {
	return len(g.members) > 0
}

// Includes implements includes logic
func (g *Group) Includes(s Shaper) bool {
	for _, m := range g.members {
		if m.Includes(s) {
			return true
		}
	}
	return false
}

// checkAxis is a helper function for check for intersection of the group's bounding box and ray
func (g *Group) checkAxis(o, d, min, max float64) (float64, float64) {

	var tmin, tmax float64

	tminNumerator := min - o
	tmaxNumerator := max - o

	if math.Abs(d) >= constants.Epsilon {
		tmin = tminNumerator / d
		tmax = tmaxNumerator / d
	} else {
		tmin = tminNumerator * math.MaxFloat64
		tmax = tmaxNumerator * math.MaxFloat64
	}

	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

	return tmin, tmax
}

// IntersectWithBoundingBox returns true if the ray intersects with the bounding box
// min and max define the bounding box
func (g *Group) IntersectWithBoundingBox(r Ray, b Bound) bool {

	var tmin, tmax float64

	xtmin, xtmax := g.checkAxis(r.Origin.X(), r.Dir.X(), b.Min.X(), b.Max.X())
	ytmin, ytmax := g.checkAxis(r.Origin.Y(), r.Dir.Y(), b.Min.Y(), b.Max.Y())
	ztmin, ztmax := g.checkAxis(r.Origin.Z(), r.Dir.Z(), b.Min.Z(), b.Max.Z())

	tmin = math.Max(math.Max(xtmin, ytmin), ztmin)
	tmax = math.Min(math.Min(xtmax, ytmax), ztmax)

	// missed the bounding box
	if tmin > tmax {
		return false
	}
	return true
}

// IntersectWith returns the 't' values of Ray r intersecting with the group
func (g *Group) IntersectWith(r Ray, t Intersections) Intersections {
	// transform the ray by the inverse of the group transfrom matrix
	// instead of changing the group, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(g.transformInverse)

	if !g.IntersectWithBoundingBox(r, g.Bounds()) {
		// bail out early, ray does not intersect group bounding box
		return t
	}

	xs := NewIntersections()

	for _, m := range g.Members() {
		mxs := m.IntersectWith(r, xs)
		t = append(t, mxs...)
		xs = xs[:0]
	}

	// sort.Sort(byT(t))
	return t
}

// NormalAt returns the normal vector at the given point on the surface of the group
func (g *Group) NormalAt(p Point, xs *Intersection) Vector {
	panic("called NormalAt on a group")
}

func (g *Group) boundBoxFromBoundingBoxes(boxes []Bound) Bound {

	if len(boxes) <= 0 {
		return Bound{
			Min: Origin(),
			Max: Origin(),
		}
	}

	var x []float64
	var y []float64
	var z []float64

	for _, b := range boxes {
		x = append(x, b.Min.X())
		x = append(x, b.Max.X())
		y = append(y, b.Min.Y())
		y = append(y, b.Max.Y())
		z = append(z, b.Min.Z())
		z = append(z, b.Max.Z())
	}

	sort.Float64s(x)
	sort.Float64s(y)
	sort.Float64s(z)

	return NewBound(
		NewPoint(x[0], y[0], z[0]),
		NewPoint(x[len(x)-1], y[len(y)-1], z[len(z)-1]))
}

// PrecomputeValues precomputes some values for render speedup
func (g *Group) PrecomputeValues() {
	// calculate group bounding box
	// g.calculateBounds()
}

// calculateBounds sets the g.bound variable
func (g *Group) calculateBounds() {

	// combine bounding boxes for all sub-objects into one

	// convert all member bounding boxes into group space
	var all []Bound

	for _, m := range g.members {
		mb := m.Bounds()

		// transform all 8 points to World space by multiplying by the m's transformation matrix
		p1 := NewPoint(mb.Min.X(), mb.Min.Y(), mb.Min.Z()).TimesMatrix(m.Transform())
		p2 := NewPoint(mb.Max.X(), mb.Min.Y(), mb.Min.Z()).TimesMatrix(m.Transform())
		p3 := NewPoint(mb.Min.X(), mb.Max.Y(), mb.Min.Z()).TimesMatrix(m.Transform())
		p4 := NewPoint(mb.Max.X(), mb.Max.Y(), mb.Min.Z()).TimesMatrix(m.Transform())
		p5 := NewPoint(mb.Min.X(), mb.Min.Y(), mb.Max.Z()).TimesMatrix(m.Transform())
		p6 := NewPoint(mb.Max.X(), mb.Min.Y(), mb.Max.Z()).TimesMatrix(m.Transform())
		p7 := NewPoint(mb.Min.X(), mb.Max.Y(), mb.Max.Z()).TimesMatrix(m.Transform())
		p8 := NewPoint(mb.Max.X(), mb.Max.Y(), mb.Max.Z()).TimesMatrix(m.Transform())

		// now find the min and max of all the point sto get the new bounding box
		all = append(all, boundingBoxFromPoints(p1, p2, p3, p4, p5, p6, p7, p8))
	}

	// not combine all bounding boxes into one
	g.bound = g.boundBoxFromBoundingBoxes(all)
}
