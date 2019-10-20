package tracer

// Shaper represents an physical object
type Shaper interface {
	IntersectWith(Ray) Intersections

	NormalAt(Point) Vector

	Material() *Material
	SetMaterial(*Material)

	SetTransform(Matrix)
	Transform() Matrix
}

// Shape is the abstract shape
type Shape struct {
	transform Matrix
	material  *Material
}

// IntersectWith implements Shaper interface
func (s *Shape) IntersectWith(r Ray) Intersections {
	return Intersections{}
}

// NormalAt implements the Shaper interface
func (s *Shape) NormalAt(p Point) Vector {
	return NewVector(1, 1, 1)
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
}
