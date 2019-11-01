package tracer

import (
	"log"
	"math"
	"sort"
	"sync"

	"golang.org/x/image/colornames"
)

// WorldConfig collects various settings to configure the world
type WorldConfig struct {
	// How many times to allow the ray to bounce between two objects (controls reflections of reflections)
	MaxRecusions int
}

// NewWorldConfig returns a new world config with default settings
func NewWorldConfig() *WorldConfig {
	return &WorldConfig{
		MaxRecusions: 4,
	}
}

// World holds everything in it
type World struct {
	Objects []Shaper
	Lights  []Light
	camera  *Camera
	Config  *WorldConfig
}

// NewWorld returns a new empty world
func NewWorld(config *WorldConfig) *World {
	return &World{
		Config: config,
	}
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
		Config:  NewWorldConfig(),
	}
}

// NewDefaultTestWorld returns a world that many tests expect
func NewDefaultTestWorld() *World {
	l1 := NewPointLight(NewPoint(-10, 10, -10), ColorName(colornames.White))

	s1 := NewUnitSphere()
	s1.SetMaterial(NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200, 0, 0, 1))

	s2 := NewUnitSphere()
	s2.SetTransform(IdentityMatrix().Scale(0.5, 0.5, 0.5))

	w := NewWorld(NewWorldConfig())
	w.Objects = []Shaper{s1, s2}
	w.Lights = []Light{l1}
	return w

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
func (w *World) ColorAt(r Ray, remaining int) Color {
	xs := w.Intersections(r)
	hit, err := xs.Hit()
	if err != nil {
		return Black()
	}

	state := PrepareComputations(hit, r, xs)
	return w.shadeHit(state, remaining).Clamp()
}

// ReflectedColor returns the reflected color given an IntersectionState
// remaining controls how many times a light ray can bounce between the same objects
func (w *World) ReflectedColor(state *IntersectionState, remaining int) Color {
	if remaining <= 0 || state.Object.Material().Reflective == 0 {
		return Black()
	}

	reflectR := NewRay(state.OverPoint, state.ReflectV)
	clr := w.ColorAt(reflectR, remaining-1)

	return clr.Scale(state.Object.Material().Reflective)
}

// RefractedColor returns the refracted color given an IntersectionState
// remaining controls how many times a light ray can bounce between the same objects
func (w *World) RefractedColor(state *IntersectionState, remaining int) Color {
	if remaining <= 0 || state.Object.Material().Transparency == 0 {
		return Black()
	}

	// check for total internal reflection

	// find the ratio of the first index of refraction to the second
	nRatio := state.N1 / state.N2

	// cos(theta_i) is the same as the dot product of the two vectors
	cosi := state.EyeV.Dot(state.NormalV)

	// find sun(theta_t)^2 via trigonometric identity
	sin2t := math.Pow(nRatio, 2) * (1 - math.Pow(cosi, 2))

	if sin2t > 1 { // total internal reflection
		return Black()
	}

	// find refracted color

	// find cos(theta_t) via trigonometric identity
	cost := math.Sqrt(1.0 - sin2t)

	// compute the direction of the refracted ray
	dir := state.NormalV.Scale(nRatio*cosi - cost).SubVector(state.EyeV.Scale(nRatio))

	// create the refracted ray
	refractedRay := NewRay(state.UnderPoint, dir)

	// find the color of the refracted ray, making sure to multiply
	// by the transparency value to account for any opacity
	clr := w.ColorAt(refractedRay, remaining-1).Scale(state.Object.Material().Transparency)

	return clr

}

// ShadeHit returns the color at the intersection enapsulated by IntersectionState
func (w *World) shadeHit(state *IntersectionState, remaining int) Color {

	var result Color

	for _, l := range w.Lights {
		isShadowed := w.IsShadowed(state.OverPoint, l)

		surface := lighting(
			state.Object.Material(),
			state.Object,
			state.OverPoint,
			l,
			state.EyeV,
			state.NormalV,
			isShadowed)

		reflected := w.ReflectedColor(state, remaining)
		refracted := w.RefractedColor(state, remaining)

		m := state.Object.Material()
		if m.Reflective > 0 && m.Transparency > 0 {
			// Use Schlick approximation for the Fresnel Effect
			reflectance := Schlick(state)

			result = surface.Add(reflected.Scale(reflectance)).Add(refracted.Scale((1 - reflectance)))
		} else {

			result = result.Add(surface.Add(reflected.Add(refracted)))
		}
	}

	return result
}

// IsShadowed returns true if p is in a shadow from the given light
// TODO: Change this to be a range rather than a bool
func (w *World) IsShadowed(p Point, l Light) bool {

	v := l.Position().SubPoint(p)
	distance := v.Magnitude()
	direction := v.Normalize()

	r := NewRay(p, direction)
	intersections := w.Intersections(r)

	// Some objects do not cast shadows, so we need to look at all the objects r intersects with
	sort.Sort(byT(intersections))

	for _, it := range intersections {
		if it.t >= 0 {

			if it.t < distance && it.Object().Material().ShadowCaster {
				return true
			}
		}
	}

	return false
}

// PrecomputeValues does some initial promcomputations to speed up render speed
func (w *World) PrecomputeValues() {
	log.Println("Precomputing values for the world...")

	for _, o := range w.Objects {
		if o.HasMembers() {
			for _, om := range o.(*Group).Members() {
				om.PrecomputeValues()
			}
		}
		o.PrecomputeValues()
	}
}

func (w *World) doRender(camera *Camera, canvas *Canvas) *Canvas {

	log.Println("Running render...")

	var wg sync.WaitGroup

	for y := 0.0; y < camera.Vsize-1; y++ {
		for x := 0.0; x < camera.Hsize-1; x++ {
			wg.Add(1)
			go func(x, y float64) {
				ray := camera.RayForPixel(x, y)
				clr := w.ColorAt(ray, w.Config.MaxRecusions)
				canvas.SetFloat(x, y, clr)
				wg.Done()
			}(x, y)
		}
	}

	wg.Wait()
	return canvas
}

// ShowInfo dumps info about the world
func (w *World) ShowInfo() {
	log.Printf("Camera: %#v", w.Camera())
}

// Render renders the world using the world camera
func (w *World) Render() *Canvas {
	w.LintWorld()
	w.PrecomputeValues()

	camera := w.Camera()
	canvas := NewCanvas(int(camera.Hsize), int(camera.Vsize))

	w.ShowInfo()

	return w.doRender(camera, canvas)
}
