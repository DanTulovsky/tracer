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
	data          map[int]map[int]color.Color
}

// NewCanvas returns a pointer to a new canvas
func NewCanvas(w, h int) *Canvas {
	data := make(map[int]map[int]color.Color)

	for x := 0; x < w; x++ {
		data[x] = make(map[int]color.Color)
		for y := 0; y < h; y++ {
			data[x][y] = color.RGBA{0, 0, 0, 0xff}
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

	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			img.Set(x, y, c.data[x][y])
		}
	}

	// Write
	png.Encode(w, img)

	return nil
}
