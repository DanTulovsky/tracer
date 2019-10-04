package color

import "math"

const (
	// Epsilon is used to compare floating point numbers
	Epsilon float64 = 0.00001
)

type Color struct {
	R float64
	G float64
	B float64
	// A float64
}

// NewColor returns a new RGB color
func NewColor(r, g, b float64) Color {
	return Color{R: r, G: g, B: b}
}

// Add adds two colors
func (c Color) Add(c2 Color) Color {
	return NewColor(c.R+c2.R, c.G+c2.G, c.B+c2.B)
}

// Sub subtracts two colors
func (c Color) Sub(c2 Color) Color {
	return NewColor(c.R-c2.R, c.G-c2.G, c.B-c2.B)
}

// Scale scales a color
func (c Color) Scale(s float64) Color {
	return NewColor(c.R*s, c.G*s, c.B*s)
}

// Mult multiplies two colors together
func (c Color) Mult(c2 Color) Color {
	return NewColor(c.R*c2.R, c.G*c2.G, c.B*c2.B)
}

// Equal compares two colors
func Equal(c1, c2 Color) bool {
	return math.Abs(c1.R-c2.R) < Epsilon && math.Abs(c1.G-c2.G) < Epsilon && math.Abs(c1.B-c2.B) < Epsilon
}
