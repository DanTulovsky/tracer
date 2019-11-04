package tracer

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"runtime"
	"sync"

	"github.com/DanTulovsky/tracer/utils"
	"golang.org/x/image/colornames"
)

// Canvas is the canvas for drawing on
type Canvas struct {
	Width, Height int
	colors        [][]Color
	// points stores all the points in an array
	// for OpenGL consumption; normalized to [-1, 1]
	// this is a list where each 3 numbers are a vertex
	points    []float32
	oglColors []float32
}

// NewCanvas returns a pointer to a new canvas
func NewCanvas(w, h int) *Canvas {
	// Allocate the top-level slice, the same as before.
	colors := make([][]Color, w) // One row per unit of y.
	points := make([]float32, w*h*3)
	oglColors := make([]float32, w*h*3)

	var i int
	// column major order
	for c := 0; c < w; c++ {
		colors[c] = make([]Color, h)
		for r := 0; r < h; r++ {
			colors[c][r] = ColorName(colornames.Black)

			// normalize c, r to be in [-1, 1]
			// canvas [0, 0] is top left
			points[i] = float32(utils.AT(float64(c), 0.0, float64(w), -1, 1))
			points[i+1] = -1 * float32(utils.AT(float64(r), 0.0, float64(h), -1, 1))
			points[i+2] = 0 // z is always 0

			oglColors[i] = 0   // R
			oglColors[i+1] = 0 // G
			oglColors[i+2] = 0 // B

			i = i + 3
		}
	}

	return &Canvas{
		Width:     w,
		Height:    h,
		colors:    colors,
		points:    points,
		oglColors: oglColors,
	}
}

// Points returns the array of vertex points
func (c *Canvas) Points() []float32 {
	return c.points
}

// Colors returns the array of vertex colors
func (c *Canvas) Colors() []float32 {
	return c.oglColors
}

// Set sets the color of a pixel
func (c *Canvas) Set(x, y int, clr Color) error {
	if x >= c.Width || y >= c.Height {
		return fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}

	c.colors[x][y] = clr

	// colum major order; (height*c+r)*3
	i := ((c.Height * x) + y) * 3
	c.oglColors[i] = float32(clr.R)
	c.oglColors[i+1] = float32(clr.G)
	c.oglColors[i+2] = float32(clr.B)
	return nil
}

// SetFloat sets the color of a pixel
func (c *Canvas) SetFloat(x, y float64, clr Color) error {
	if int(x) >= c.Width || int(y) >= c.Height {
		return fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}

	intx := int(x)
	inty := int(y)

	c.colors[intx][inty] = clr

	// colum major order; (height*c+r)*3
	i := ((c.Height * intx) + inty) * 3
	c.oglColors[i] = float32(clr.R)
	c.oglColors[i+1] = float32(clr.G)
	c.oglColors[i+2] = float32(clr.B)
	return nil
}

// Get returns the color at the given coordinates
func (c *Canvas) Get(x, y int) (Color, error) {
	if x >= c.Width || y >= c.Height {
		return ColorName(colornames.Black), fmt.Errorf("coordinates [%v, %v] are outside the canvas", x, y)
	}
	return c.colors[x][y], nil
}

// ExportToPNG exports the canvas to a png file
func (c *Canvas) ExportToPNG(w io.Writer) error {
	// create an image covering the entire canvas
	upLeft := image.Point{0, 0}
	lowRight := image.Point{c.Width, c.Height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var wg sync.WaitGroup
	sem := make(chan bool, runtime.NumCPU())

	for col := 0; col < c.Width; col++ {
		for row := 0; row < c.Height; row++ {
			sem <- true
			go func(img *image.RGBA, col, row int, clr color.Color) {
				defer func() { <-sem }()
				img.Set(col, row, clr)
				wg.Done()
			}(img, col, row, c.colors[col][row])
		}
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	// Write
	if err := png.Encode(w, img); err != nil {
		fmt.Println(err)
	}

	return nil
}
