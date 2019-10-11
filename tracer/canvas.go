package tracer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

// Canvas is the canvas for drawing on
type Canvas struct {
	Width, Height int
	data          [][]color.Color
}

// NewCanvas returns a pointer to a new canvas
func NewCanvas(w, h int) *Canvas {
	// Allocate the top-level slice, the same as before.
	data := make([][]color.Color, w) // One row per unit of y.

	for c := 0; c < w; c++ {
		data[c] = make([]color.Color, h)
		for r := 0; r < h; r++ {
			data[c][r] = color.RGBA{0, 0, 0, 0xff}
		}
	}

	return &Canvas{Width: w, Height: h, data: data}
}

// Set sets the color of a pixel
func (c *Canvas) Set(x, y int, clr color.Color) error {
	if x >= c.Width || y >= c.Height {
		return fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}

	c.data[x][y] = clr
	return nil
}

// SetFloat sets the color of a pixel
func (c *Canvas) SetFloat(x, y float64, clr color.Color) error {
	if int(x) >= c.Width || int(y) >= c.Height {
		return fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}

	c.data[int(x)][int(y)] = clr
	return nil
}

// Get returns the color at the given coordinates
func (c *Canvas) Get(x, y int) (color.Color, error) {
	if x >= c.Width || y >= c.Height {
		return nil, fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}
	return c.data[x][y], nil
}

// ExportToPNG exports the canvas to a png file
func (c *Canvas) ExportToPNG(w io.Writer) error {
	// create an image covering the entire canvas
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.Width, c.Height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for col := 0; col < c.Width; col++ {
		for row := 0; row < c.Height; row++ {
			img.Set(col, row, c.data[col][row])
		}
	}

	// Write
	png.Encode(w, img)

	return nil
}
