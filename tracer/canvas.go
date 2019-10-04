package tracer

import (
	"image"
	"image/color"
	"image/png"
	"io"
)

// Canvas is the canvas for drawing on
type Canvas struct {
	Width, Height int
	Data          map[int]map[int]color.Color
}

// NewCanvas returns a pointer to a new canvas
func NewCanvas(w, h int) *Canvas {
	data := make(map[int]map[int]color.Color)

	for x := 0; x < w; x++ {
		data[x] = make(map[int]color.Color)
		for y := 0; y < h; y++ {
			data[x][y] = color.RGBA{0, 0, 0, 0}
		}
	}
	return &Canvas{Width: w, Height: h, Data: data}
}

// ExportToPNG exports the canvas to a png file
func (c *Canvas) ExportToPNG(w io.Writer) error {
	// create an image covering the entire canvas
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.Width, c.Height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			img.Set(x, y, c.Data[x][y])
		}
	}

	// Write
	png.Encode(w, img)

	return nil
}
