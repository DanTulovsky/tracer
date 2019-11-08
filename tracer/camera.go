package tracer

import "math"

// ViewTransform returns the view transform matrix given from, to and up vectors
func ViewTransform(from, to Point, up Vector) Matrix {
	forward := to.SubPoint(from).Normalize()
	upn := up.Normalize()
	left := forward.Cross(upn)
	trueUp := left.Cross(forward)

	o := NewMatrixFromData([][]float64{
		{left.X(), left.Y(), left.Z(), 0},
		{trueUp.X(), trueUp.Y(), trueUp.Z(), 0},
		{-forward.X(), -forward.Y(), -forward.Z(), 0},
		{0, 0, 0, 1},
	})

	return o.TimesMatrix(NewTranslation(-from.X(), -from.Y(), -from.Z()))

}

// Camera defines the camera looking at the world
type Camera struct {
	Hsize, Vsize          float64 // canvas size
	fov                   float64 // field of view angle in radians
	Transform             Matrix  // view transformation matrix from the above function
	TransformInverse      Matrix  // pre-cache the inverse as it's called for each pixel
	HalfWidth, HalfHeight float64
	PixelSize             float64
}

// NewCamera returns
func NewCamera(hsize, vsize, fov float64) *Camera {
	c := &Camera{
		Hsize:            hsize,
		Vsize:            vsize,
		fov:              fov,
		Transform:        IdentityMatrix(),
		TransformInverse: IdentityMatrix().Inverse(),
	}

	c.setPixelSize()
	return c
}

// SetTransform sets the transform on the camera
func (c *Camera) SetTransform(t Matrix) {
	c.Transform = t
	c.TransformInverse = t.Inverse()
}

// SetFoV sets the field of view and recalculates the pixel size
func (c *Camera) SetFoV(fov float64) {
	c.fov = fov
	c.setPixelSize()
}

// setPixelSize sets the world-space pixel size and half view values into the camera
// Assumes the canvas is one unit away
func (c *Camera) setPixelSize() {
	halfView := math.Tan(c.fov / 2)
	aspect := c.Hsize / c.Vsize

	switch {
	case aspect >= 1:
		c.HalfWidth = halfView
		c.HalfHeight = halfView / aspect
	default:
		c.HalfWidth = halfView * aspect
		c.HalfHeight = halfView
	}

	c.PixelSize = (c.HalfWidth * 2) / c.Hsize
}

// RayForPixel returns a ray that starts at the camera and passes through x,y on the canvas
func (c *Camera) RayForPixel(x, y float64) Ray {

	// the offset from the edge of the canvas to the pixel's center
	xoffset := (x + 0.5) * c.PixelSize
	yoffset := (y + 0.5) * c.PixelSize

	// untransformed coordinates of the pixel in world space
	// camera looks toward -z, so +x is to the left
	wx := c.HalfWidth - xoffset
	wy := c.HalfHeight - yoffset

	// transform the canvas point and the origin using the camera's matrix
	pixel := NewPoint(wx, wy, -1).TimesMatrix(c.TransformInverse)
	origin := Origin().TimesMatrix(c.TransformInverse)
	direction := pixel.SubPoint(origin).Normalize()

	return NewRay(origin, direction)

}
