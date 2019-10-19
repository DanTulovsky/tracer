package tracer

import (
	"sort"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

// World holds everything in it
type World struct {
	Objects []Object
	Lights  []Light
	camera  *Camera
}

// NewWorld returns a new empty world
func NewWorld() *World {
	return &World{}
}

// NewDefaultTestWorld returns a world that many tests expect
func NewDefaultTestWorld() *World {
	l1 := NewPointLight(NewPoint(-10, 10, -10), ColorName(colornames.White))

	s1 := NewUnitSphere()
	s1.SetMaterial(NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200))

	s2 := NewUnitSphere()
	s2.SetTransform(IdentityMatrix().Scale(0.5, 0.5, 0.5))

	return &World{
		Objects: []Object{s1, s2},
		Lights:  []Light{l1},
	}
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

// ShadeHit returns the color at the intersection enapsulated by IntersectionState
func (w *World) shadeHit(state IntersectionState) Color {

	var result Color

	for _, l := range w.Lights {
		c := lighting(
			state.Object.Material(),
			state.Point,
			l,
			state.EyeV,
			state.NormalV)

		result = result.Add(c)
	}

	return result
}

// SetLights sets the world lights
func (w *World) SetLights(l []Light) {
	w.Lights = l
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

// Render renders the world using the world camera
func (w *World) Render() *Canvas {
	camera := w.Camera()
	canvas := NewCanvas(int(camera.Hsize), int(camera.Vsize))

	for y := 0.0; y < camera.Vsize-1; y++ {
		for x := 0.0; x < camera.Hsize-1; x++ {
			ray := camera.RayForPixel(x, y)
			clr := w.ColorAt(ray)
			canvas.SetFloat(x, y, clr)
		}
	}
	return canvas
}
