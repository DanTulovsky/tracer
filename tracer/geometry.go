package tracer

// Shaper represents an physical object
type Shaper interface {
	Bounds() Bound

	HasMembers() bool

	Parent() Shaper
	HasParent() bool
	SetParent(Shaper)

	// If A is a Group, true if any child includes B
	// If A is a CSG, true if either child includes B
	// If A is anything else, true if A == B
	Includes(Shaper) bool

	IntersectWith(Ray, Intersections) Intersections
	NormalAt(Point, *Intersection) Vector
	PrecomputeValues()

	Material() *Material
	SetMaterial(*Material)

	Name() string
	SetName(string)

	SetTransform(Matrix)
	Transform() Matrix

	TransformInverse() Matrix
}

// Shape is the abstract shape
type Shape struct {
	name             string
	shape            string
	transform        Matrix
	transformInverse Matrix
	material         *Material
	bound            Bound // cache the group bounding box

	parent Shaper

	// localNormalAt
	lna func(Point, *Intersection) Vector
}

// Equal returns true if the shapes are equal
func (s *Shape) Equal(s2 *Shape) bool {
	return s.shape == s2.shape &&
		s.name == s2.name &&
		s.transform.Equals(s2.Transform()) &&
		s.transformInverse.Equals(s2.Transform()) &&
		s.material.Equals(s2.material) &&
		// s.bound == s2.bound &&
		s.parent == s2.parent &&
		(s.lna != nil) == (s2.lna != nil)

}

// Includes implements includes logic
func (s *Shape) Includes(s2 Shaper) bool {
	panic("please implement Includes")
}

// HasMembers returns true if this is a group that has members
func (s *Shape) HasMembers() bool {
	return false
}

// PrecomputeValues precomputes some values for render speedup
func (s *Shape) PrecomputeValues() {
	// nothing by default, each shape can override
	panic("please implement PrecomputeValues")
}

// Parent returns the parent group this shape is part of
func (s *Shape) Parent() Shaper {
	return s.parent
}

// HasParent returns True if this shape has a parent
func (s *Shape) HasParent() bool {
	return s.parent != nil
}

// SetParent sets the parents of the object
func (s *Shape) SetParent(p Shaper) {
	s.parent = p
}

// IntersectWith implements Shaper interface
func (s *Shape) IntersectWith(r Ray, xs Intersections) Intersections {
	panic("must implement IntersectWith")
}

// NormalAt implements the Shaper interface
func (s *Shape) NormalAt(p Point, xs *Intersection) Vector {
	// move point to object space
	op := p.ToObjectSpace(s)

	// object normal, this is different for each shape
	on := s.lna(op, xs)

	// world normal
	wn := on.NormalToWorldSpace(s)

	return wn.Normalize()
}

// localNormalAt return the local normal vector at the point
func (s *Shape) localNormalAt(p Point, xs *Intersection) Vector {
	panic("must implement localNormalAt")
}

// Material returns the material of the shape
func (s *Shape) Material() *Material {
	return s.material
}

// SetMaterial sets the material of the shape
func (s *Shape) SetMaterial(m *Material) {
	s.material = m
}

// Transform returns the transformation matrix of the shape
func (s *Shape) Transform() Matrix {
	return s.transform
}

// SetTransform sets the transformation matrix of the shape
func (s *Shape) SetTransform(m Matrix) {
	s.transform = m
	s.transformInverse = m.Inverse()
}

// TransformInverse returns the inverse of the transformation matrix of the shape
func (s *Shape) TransformInverse() Matrix {
	return s.transformInverse
}

// Name returns the name of the shape
func (s *Shape) Name() string {
	return s.name
}

// SetName sets the name
func (s *Shape) SetName(n string) {
	s.name = n
}

// calculateBounds calculates the bounding box of the shape
func (s *Shape) calculateBounds() {
	panic("please implement calculateBounds!")
}

// Bounds returns the untransformed bounding box
func (s *Shape) Bounds() Bound {
	return s.bound
}

// Bound describes the bounding box for a shape
type Bound struct {
	Min, Max Point
}
