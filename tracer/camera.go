package tracer

// ViewTransform returns the view transform matrix given from, to and up vectors
func ViewTransform(from, to, up Vector) Matrix {
	forward := to.SubVector(from).Normalize()
	upn := up.Normalize()
	left := forward.Cross(upn)
	trueUp := left.Cross(forward)

	o := NewMatrixFromData([][]float64{
		{left.X(), left.Y(), left.Z(), 0},
		{trueUp.X(), trueUp.Y(), trueUp.Z(), 0},
		{-forward.X(), -forward.Y(), -forward.Z(), 0},
		{0, 0, 0, 1},
	})

	return o.TimesMatrix(NewTranslation(from.Negate().X(), from.Negate().Y(), from.Negate().Z()))

}

// Camera defines the camera looking at the world
type Camera struct {
	Hsize, Vsize int     // canvas size
	FoV          float64 // field of view angle in radians
	Transform    Matrix  // view transformation matrix from the above function
}

// NewCamera returns
func NewCamera(hsize, vsize int, fov float64) *Camera {
	return &Camera{
		Hsize:     hsize,
		Vsize:     vsize,
		FoV:       fov,
		Transform: IdentityMatrix(),
	}
}

// SetTransform sets the transform on the camera
func (c *Camera) SetTransform(t Matrix) {
	c.Transform = t
}
