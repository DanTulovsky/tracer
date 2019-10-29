package tracer

import "sort"

// Group is a collection of other groups/objects
type Group struct {
	members []Shaper
	Shape
}

// NewGroup returns a new, empty group
func NewGroup() *Group {
	return &Group{
		members: []Shaper{},
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "group",
		},
	}
}

// AddMember adds a new member to this group
func (g *Group) AddMember(m Shaper) {
	g.members = append(g.members, m)
	m.SetParent(g)
}

// Members returns all the direct members of this group
func (g *Group) Members() []Shaper {
	return g.members
}

// Includes returns true if the Shaper is part of this group
func (g *Group) Includes(s Shaper) bool {
	for _, m := range g.members {
		if m == s {
			return true
		}
	}
	return false
}

// IntersectWith returns the 't' values of Ray r intersecting with the group in sorted order
func (g *Group) IntersectWith(r Ray) Intersections {
	t := Intersections{}

	// transform the ray by the inverse of the group transfrom matrix
	// instead of changing the group, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(g.transform.Inverse())

	for _, m := range g.Members() {
		mxs := m.IntersectWith(r)
		t = append(t, mxs...)
	}

	sort.Sort(byT(t))
	return t
}

// NormalAt returns the normal vector at the given point on the surface of the group
func (g *Group) NormalAt(p Point) Vector {
	panic("called NormalAt on a group")
}
