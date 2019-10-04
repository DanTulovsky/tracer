package tracer

// Tuple is the base for Vector and Point
type Tuple interface {
	X() float64
	Y() float64
	Z() float64
	W() float64
	Equals(Tuple) bool
	Add(Tuple) Tuple
	Sub(Tuple) Tuple
}
