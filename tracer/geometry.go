package tracer

// Tupler ...

// Object represents an physical object
type Object interface {
	IntersectWith(Ray) Intersections
}
