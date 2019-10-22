package tracer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"sync"

	"golang.org/x/image/colornames"
)

// Canvas is the canvas for drawing on
type Canvas struct {
	Width, Height int
	data          [][]Color
}

// NewCanvas returns a pointer to a new canvas
func NewCanvas(w, h int) *Canvas {
	// Allocate the top-level slice, the same as before.
	data := make([][]Color, w) // One row per unit of y.

	for c := 0; c < w; c++ {
		data[c] = make([]Color, h)
		for r := 0; r < h; r++ {
			data[c][r] = ColorName(colornames.Black)
		}
	}

	return &Canvas{Width: w, Height: h, data: data}
}

// Set sets the color of a pixel
func (c *Canvas) Set(x, y int, clr Color) error {
	if x >= c.Width || y >= c.Height {
		return fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}

	c.data[x][y] = clr
	return nil
}

// SetFloat sets the color of a pixel
func (c *Canvas) SetFloat(x, y float64, clr Color) error {
	if int(x) >= c.Width || int(y) >= c.Height {
		return fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}

	c.data[int(x)][int(y)] = clr
	return nil
}

// Get returns the color at the given coordinates
func (c *Canvas) Get(x, y int) (Color, error) {
	if x >= c.Width || y >= c.Height {
		return ColorName(colornames.Black), fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}
	return c.data[x][y], nil
}

// ExportToPNG exports the canvas to a png file
func (c *Canvas) ExportToPNG(w io.Writer) error {
	// create an image covering the entire canvas
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.Width, c.Height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var wg sync.WaitGroup

	// Almost certainly not needed here, but good exercise anyway
	// TODO: Limit this to MAXProcs at a time - channels?
	// e.g. https://medium.com/@zufolo/a-pattern-for-limiting-the-number-of-goroutines-in-execution-56e13b226e72
	for col := 0; col < c.Width; col++ {
		for row := 0; row < c.Height; row++ {
			wg.Add(1)
			go func(img *image.RGBA, col, row int, clr color.Color) {
				img.Set(col, row, clr)
				wg.Done()
			}(img, col, row, c.data[col][row])
		}
	}

	wg.Wait()

	// Write
	if err := png.Encode(w, img); err != nil {
		fmt.Println(err)
	}

	return nil
}
