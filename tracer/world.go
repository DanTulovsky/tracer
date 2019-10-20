package tracer

import (
	"math"
	"sort"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

// World holds everything in it
type World struct {
	Objects []Shaper
	Lights  []Light
	camera  *Camera
}

// NewWorld returns a new empty world
func NewWorld() *World {
	return &World{}
}

// NewDefaultWorld returns a default world
func NewDefaultWorld(width, height float64) *World {
	defaultLight := NewPointLight(NewPoint(-10, 10, -10), ColorName(colornames.White))
	camera := NewCamera(width, height, math.Pi/3)
	viewTransform := ViewTransform(NewPoint(0, 1.5, -5), NewPoint(0, 1, 0), NewVector(0, 1, 0))
	camera.SetTransform(viewTransform)

	return &World{
		Objects: []Shaper{},
		Lights:  []Light{defaultLight},
		camera:  camera,
	}
}

// NewDefaultTestWorld returns a world that many tests expect
func NewDefaultTestWorld() *World {
	l1 := NewPointLight(NewPoint(-10, 10, -10), ColorName(colornames.White))

	s1 := NewUnitSphere()
	s1.SetMaterial(NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200))

	s2 := NewUnitSphere()
	s2.SetTransform(IdentityMatrix().Scale(0.5, 0.5, 0.5))

	return &World{
		Objects: []Shaper{s1, s2},
		Lights:  []Light{l1},
	}
}

// AddObject adds an object into the world
func (w *World) AddObject(o Shaper) {
	w.Objects = append(w.Objects, o)
}

// Intersections returns all the intersections in the world with the given ray
func (w *World) Intersections(r Ray) Intersections {
	var is Intersections

	for _, o := range w.Objects {
		iso := o.IntersectWith(r)
		is = append(is, iso...)
	}

	sort.Sort(byT(is))

	return is
}

// SetLights sets the world lights
func (w *World) SetLights(l []Light) {
	w.Lights = l
}

// AddLight adds a new light to the world
func (w *World) AddLight(l Light) {
	w.Lights = append(w.Lights, l)
}

// SetCamera sets the world camera
func (w *World) SetCamera(c *Camera) {
	w.camera = c
}

// Camera returns the world camera
func (w *World) Camera() *Camera {
	return w.camera
}

// ColorAt returns the color in the world where the given ray hits
func (w *World) ColorAt(r Ray) Color {
	is := w.Intersections(r)
	hit, err := is.Hit()
	if err != nil {
		return ColorName(colornames.Black)
	}

	state := PrepareComputations(hit, r)
	return w.shadeHit(state).Clamp()
}

// ShadeHit returns the color at the intersection enapsulated by IntersectionState
func (w *World) shadeHit(state IntersectionState) Color {

	var result Color

	for _, l := range w.Lights {
		isShadowed := w.IsShadowed(state.OverPoint, l)

		c := lighting(
			state.Object.Material(),
			state.OverPoint,
			l,
			state.EyeV,
			state.NormalV,
			isShadowed)

		result = result.Add(c)
	}

	return result
}

// IsShadowed returns true if p is in a shadow from the given light
func (w *World) IsShadowed(p Point, l Light) bool {
	var inShadow bool

	v := l.Position().SubPoint(p)
	distance := v.Magnitude()
	direction := v.Normalize()

	r := NewRay(p, direction)
	intersections := w.Intersections(r)

	h, err := intersections.Hit()
	if err == nil && h.T() < distance {
		inShadow = true
	}

	return inShadow
}

// Render renders the world using the world camera
func (w *World) Render() *Canvas {
	camera := w.Camera()
	canvas := NewCanvas(int(camera.Hsize), int(camera.Vsize))

	// var wg sync.WaitGroup

	for y := 0.0; y < camera.Vsize-1; y++ {
		for x := 0.0; x < camera.Hsize-1; x++ {
			// wg.Add(1)
			// go func(x, y float64) {
			ray := camera.RayForPixel(x, y)
			clr := w.ColorAt(ray)
			canvas.SetFloat(x, y, clr)
			// wg.Done()
			// }(x, y)
		}
	}

	// wg.Wait()

	return canvas
}
