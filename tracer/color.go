package tracer

import (
	"image/color"
	"math"

	"github.com/DanTulovsky/tracer/utils"
	"github.com/lucasb-eyer/go-colorful"
)

// Color represents a color with values between 0-1
type Color struct {
	R, G, B float64
}

// NewColor returns a new color
func NewColor(r, g, b float64) Color {
	return Color{r, g, b}
}

// ColorName returns a Color object given a color.Color interface
func ColorName(c color.Color) Color {
	clr, _ := colorful.MakeColor(c)
	return Color{clr.R, clr.G, clr.B}
}

// Black returns the color black
func Black() Color {
	return NewColor(0, 0, 0)
}

// White returns the color white
func White() Color {
	return NewColor(1, 1, 1)
}

// Clamp returns the color with values <0 set to 0 and values >1 set to 1
func (c Color) Clamp() Color {
	var r, g, b float64

	switch {
	case c.R < 0:
		r = 0
	case c.R > 1:
		r = 1
	default:
		r = c.R
	}

	switch {
	case c.G < 0:
		g = 0
	case c.G > 1:
		g = 1
	default:
		g = c.G
	}

	switch {
	case c.B < 0:
		b = 0
	case c.B > 1:
		b = 1
	default:
		b = c.B
	}
	return NewColor(r, g, b)
}

// RGBA implements the color.Color interface
func (c Color) RGBA() (r, g, b, a uint32) {
	cl := colorful.Color{R: c.R, G: c.G, B: c.B}

	return cl.RGBA()
}

// Add adds to colors together
func (c Color) Add(c2 Color) Color {
	return Color{c.R + c2.R, c.G + c2.G, c.B + c2.B}
}

// Sub subtracs one color from the other
func (c Color) Sub(c2 Color) Color {
	return Color{c.R - c2.R, c.G - c2.G, c.B - c2.B}
}

// Scale scales a color by a scalar
func (c Color) Scale(s float64) Color {
	return Color{c.R * s, c.G * s, c.B * s}
}

// Blend blends two colors together by multiplying the rgb components by each other
// This is the Hadamard product (or Shur product)
func (c Color) Blend(c2 Color) Color {
	return Color{c.R * c2.R, c.G * c2.G, c.B * c2.B}
}

// Equal compares two colors to within Epsilon
func (c Color) Equal(c2 Color) bool {
	if utils.Equals(c.R, c2.R) && utils.Equals(c.G, c2.G) && utils.Equals(c.B, c2.B) {
		return true
	}
	return false
}

// GreyScale returns the GreyScale value of this color
// from https://stackoverflow.com/questions/17615963/standard-rgb-to-grayscale-conversion
func (c Color) GreyScale() float64 {
	var csrgb float64

	clin := 0.2126*c.R + 0.7152*c.G + 0.0722*c.B
	switch {
	case clin <= 0.0031308:
		csrgb = 12.92 * clin
	default:
		csrgb = 1.055*math.Pow(clin, 1/2.4) - 0.055
	}

	return csrgb
}

// Colors is a list of Color
type Colors []Color

// Average returns the average of all the colors in the list
func (cl Colors) Average() Color {
	var r, g, b float64

	for _, c := range cl {
		r = r + c.R
		g = g + c.G
		b = b + c.B
	}

	l := float64(len(cl))
	return NewColor(r/l, g/l, b/l)
}
