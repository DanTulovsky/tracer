package tracer

// Tupler ...
type Tupler interface {
	Equals(Tupler) bool
	X() float64
	Y() float64
	Z() float64
	W() float64
	SetX(float64)
	SetY(float64)
	SetZ(float64)
	SetW(float64)
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

// SetX sets the value
func (t Tuple) SetX(x float64) {
	t.x = x
}

// SetY sets the value
func (t Tuple) SetY(y float64) {
	t.y = y
}

// SetZ sets the value
func (t Tuple) SetZ(z float64) {
	t.z = z
}

// SetW sets the value
func (t Tuple) SetW(w float64) {
	t.w = w
}

// Equals compares tuples
func (t Tuple) Equals(s Tupler) bool {
	if Equals(t.X(), s.X()) && Equals(t.Y(), s.Y()) && Equals(t.Z(), s.Z()) && Equals(t.W(), s.W()) {
		return true
	}
	return false
}

// AddVector adds a point to a vector
func (t Tuple) AddVector(v Vector) Tupler {
	t.SetX(t.X() + v.X())
	t.SetY(t.Y() + v.Y())
	t.SetZ(t.Z() + v.Z())
	t.SetW(t.W() + v.W())

	return t
	// return NewPoint(t.X()+v.X(), t.Y()+v.Y(), t.Z()+v.Z())

}

// SubPoint subtracts points
func (t Tuple) SubPoint(p Point) Tupler {
	return NewVector(t.X()-p.X(), t.Y()-p.Y(), t.Z()-p.Z())
}

// SubVector subtracts a vector
func (t Tuple) SubVector(v Vector) Point {
	return NewPoint(t.X()-v.X(), t.Y()-v.Y(), t.Z()-v.Z())
}

// TimesMatrix multiplies a point by the matrix
func (t Tuple) TimesMatrix(m Matrix) Point {
	return NewPoint(
		m[0][0]*t.X()+m[0][1]*t.Y()+m[0][2]*t.Z()+m[0][3]*t.W(),
		m[1][0]*t.X()+m[1][1]*t.Y()+m[1][2]*t.Z()+m[1][3]*t.W(),
		m[2][0]*t.X()+m[2][1]*t.Y()+m[2][2]*t.Z()+m[2][3]*t.W())
}
