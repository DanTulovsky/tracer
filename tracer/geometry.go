package tracer

import "fmt"

// Shaper represents an physical object
type Shaper interface {
	IntersectWith(Ray) Intersections

	NormalAt(Point) Vector

	Material() *Material
	SetMaterial(*Material)

	Name() string

	SetTransform(Matrix)
	Transform() Matrix
}

// Shape is the abstract shape
type Shape struct {
	name      string
	shape     string
	transform Matrix
	material  *Material
}

// IntersectWith implements Shaper interface
func (s *Shape) IntersectWith(r Ray) Intersections {
	panic("must implement IntersectWith")
}

// NormalAt implements the Shaper interface
func (s *Shape) NormalAt(p Point) Vector {
	panic("must implement NormalAt")
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

// Name returns the name of the shape
func (s *Shape) Name() string {
	return fmt.Sprintf("%s (%s)", s.name, s.shape)
}
