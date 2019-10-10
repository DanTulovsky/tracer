package tracer

// Tupler ...
type Tupler interface {
	Equals(Tupler) bool
	X() float64
	Y() float64
	Z() float64
	W() float64
}

// Tuple is a base struct for Points and Vectors
type Tuple struct {
	x, y, z, w float64
}

// NewTuple returns a new Tuple
func NewTuple(x, y, z, w float64) Tuple {

	return Tuple{x, y, z, w}
}

// X returns the point's X coordinate
func (t Tuple) X() float64 {
	return t.x
}

// Y returns the point's y coordinate
func (t Tuple) Y() float64 {
	return t.y
}

// Z returns the point's Z coordinate
func (t Tuple) Z() float64 {
	return t.z
}

// W returns the point's W coordinate
func (t Tuple) W() float64 {
	return t.w
}

// Equals compares tuples
func (t Tuple) Equals(s Tupler) bool {
	if Equals(t.X(), s.X()) && Equals(t.Y(), s.Y()) && Equals(t.Z(), s.Y()) && Equals(t.W(), s.W()) {
		return true
	}
	return false
}
