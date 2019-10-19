package tracer

// Tupler ...

// Object represents an physical object
type Object interface {
	IntersectWith(Ray) Intersections

	NormalAt(Point) Vector

	SetTransform(m Matrix)
	Transform() Matrix
}
