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
