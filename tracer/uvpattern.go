package tracer

import (
	"image"
	"log"
	"math"
	"os"
)

// UVPatterner is a pattern that acceps UV coordinates
type UVPatterner interface {
	UVColorAt(float64, float64) Color
}

// UVImagePattern maps an image to the surface of an object
type UVImagePattern struct {
	canvas *Canvas
}

// imageToCanvas converts an image in
func imageToCanvas(m image.Image) *Canvas {

	bounds := m.Bounds()
	canvas := NewCanvas(bounds.Max.X, bounds.Max.Y)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			clr := ColorName(m.At(x, y))
			canvas.Set(x, y, clr)
		}
	}
	return canvas
}

// NewUVImagePatternImage returns a new image pattern, input is an image.Image interface
func NewUVImagePatternImage(m image.Image) (*UVImagePattern, error) {
	// convert to canvas
	canvas := imageToCanvas(m)

	p := &UVImagePattern{
		canvas: canvas,
	}

	return p, nil
}

// NewUVImagePattern returns a new image pattern
func NewUVImagePattern(filename string) (*UVImagePattern, error) {
	// read in image
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Decode image
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// convert to canvas
	canvas := imageToCanvas(m)

	p := &UVImagePattern{
		canvas: canvas,
	}

	return p, nil
}

// UVColorAt returns the color at the 2D coordinate (u, v)
func (uvip *UVImagePattern) UVColorAt(u, v float64) Color {
	v = 1 - v

	x := u * (float64(uvip.canvas.Width) - 1)
	y := v * (float64(uvip.canvas.Height) - 1)

	c, err := uvip.canvas.Get(int(x), int(y))
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// UVCheckersPattern maps checkers to the surface of the object
type UVCheckersPattern struct {
	a, b Color
	w, h float64 // how many squares along width and height
}

// NewUVCheckersPattern returns a new UV mapped checkers pattern
// If you want your checkers to look "square" on the sphere,
// be sure and set the width to twice the height. This is because of
// how the spherical map is implemented. While both u and v go from 0 to 1,
// v maps 1 to π, but u maps 1 to 2π.
func NewUVCheckersPattern(w, h float64, a, b Color) *UVCheckersPattern {
	return &UVCheckersPattern{
		a: a,
		b: b,
		w: w,
		h: h,
	}
}

// UVColorAt returns the color at the 2D coordinate (u, v)
func (ucp *UVCheckersPattern) UVColorAt(u, v float64) Color {
	u2 := math.Floor(u * ucp.w)
	v2 := math.Floor(v * ucp.h)

	if math.Mod((u2+v2), 2) == 0 {
		return ucp.a
	}

	return ucp.b
}

// UVAlignCheckPattern ...
type UVAlignCheckPattern struct {
	main, ul, ur, bl, br Color
}

// NewUVAlignCheckPattern returns a new ...
func NewUVAlignCheckPattern(main, ul, ur, bl, br Color) *UVAlignCheckPattern {
	return &UVAlignCheckPattern{
		main: main,
		ul:   ul,
		ur:   ur,
		bl:   bl,
		br:   br,
	}
}

// UVColorAt returns the color at the 2D coordinate (u, v)
func (uap *UVAlignCheckPattern) UVColorAt(u, v float64) Color {
	switch {
	case v > 0.8:
		switch {
		case u < 0.2:
			return uap.ul
		case u > 0.8:
			return uap.ur
		}
	case v < 0.2:
		switch {
		case u < 0.2:
			return uap.bl
		case u > 0.8:
			return uap.br
		}

	}
	return uap.main
}
