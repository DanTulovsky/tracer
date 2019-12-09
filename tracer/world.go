package tracer

import (
	"log"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"golang.org/x/image/colornames"
)

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
	s1.SetMaterial(NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200, 0, 0, 1, nil))

	s2 := NewUnitSphere()
	s2.SetTransform(IM().Scale(0.5, 0.5, 0.5))

	w := NewWorld(NewWorldConfig())
	w.AddObject(s1)
	w.AddObject(s2)
	w.SetLights([]Light{l1})

	return w
}

// AddObject adds an object into the world
func (w *World) AddObject(o Shaper) {
	w.Objects = append(w.Objects, o)
}

// Intersections returns all the intersections in the world with the given ray (sorted)
func (w *World) Intersections(r Ray, xs Intersections) Intersections {
	var is Intersections

	for _, o := range w.Objects {
		iso := o.IntersectWith(r, xs)
		is = append(is, iso...)
	}

	sort.Sort(byT(is))

	return is
}

// SetLights sets the world lights
func (w *World) SetLights(l Lights) {
	w.Lights = l
	for _, l := range w.Lights {
		if l.IsVisible() {
			switch l.(type) {
			case *AreaLight:
				w.AddObject(l.(*AreaLight))
			case *AreaSpotLight:
				w.AddObject(l.(*AreaSpotLight))
			}
		}
	}
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
func (w *World) ColorAt(r Ray, remaining int, xs Intersections, rng *rand.Rand) Color {
	// First solve the visibility problem
	xs = w.Intersections(r, xs)
	hit, err := xs.Hit()
	if err != nil {
		return Black()
	}

	// Second solve the shading problem
	state := PrepareComputations(hit, r, xs)
	return w.shadeHit(state, remaining, xs, rng).Clamp()
}

// ReflectedColor returns the reflected color given an IntersectionState
// remaining controls how many times a light ray can bounce between the same objects
func (w *World) ReflectedColor(state *IntersectionState, remaining int, xs Intersections, rng *rand.Rand) Color {
	if remaining <= 0 || state.Object.Material().Reflective == 0 {
		return Black()
	}

	reflectR := NewRay(state.OverPoint, state.ReflectV)
	xs = xs[:0]
	clr := w.ColorAt(reflectR, remaining-1, xs, rng)

	return clr.Scale(state.Object.Material().Reflective)
}

// RefractedColor returns the refracted color given an IntersectionState
// remaining controls how many times a light ray can bounce between the same objects
func (w *World) RefractedColor(state *IntersectionState, remaining int, xs Intersections, rng *rand.Rand) Color {
	if remaining <= 0 || state.Object.Material().Transparency == 0 {
		return Black()
	}
	// check for total internal reflection

	// find the ratio of the first index of refraction to the second
	nRatio := state.N1 / state.N2

	// cos(theta_i) is the same as the dot product of the two vectors
	cosi := state.EyeV.Dot(state.NormalV)

	// find sun(theta_t)^2 via trigonometric identity
	sin2t := nRatio * nRatio * (1 - cosi*cosi)

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
	xs = xs[:0]
	clr := w.ColorAt(refractedRay, remaining-1, xs, rng).Scale(state.Object.Material().Transparency)

	return clr

}

// ShadeHit returns the color at the intersection enapsulated by IntersectionState
func (w *World) shadeHit(state *IntersectionState, remaining int, xs Intersections, rng *rand.Rand) Color {

	xs = xs[:0] // clear intersections
	var result Color

	for _, l := range w.Lights {
		inensity := w.IntensityAt(state.OverPoint, l, xs, rng)

		surface := lighting(
			state.Object.Material(),
			state.Object,
			state.OverPoint,
			l,
			state.EyeV,
			state.NormalV,
			inensity,
			w.Config.AreaLightRays,
			state.U,
			state.V,
			rng)

		reflected := w.ReflectedColor(state, remaining, xs, rng)
		refracted := w.RefractedColor(state, remaining, xs, rng)

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

// IntensityAt returns the intensity of the light at point p
func (w *World) IntensityAt(p Point, l Light, xs Intersections, rng *rand.Rand) float64 {
	switch l.(type) {
	case *PointLight, *SpotLight:
		if w.IsShadowed(p, l.Position(), xs) {
			return 0
		}
		return 1
	case *AreaLight, *AreaSpotLight:
		if !w.Config.SoftShadows {
			if w.IsShadowed(p, l.Position(), xs) {
				return 0
			}
			return 1
		}
		total := 0.0
		for try := 0; try < w.Config.SoftShadowRays; try++ {
			if !w.IsShadowed(p, l.RandomPosition(rng), xs) {
				total = total + 1
			}
		}
		return total / float64(w.Config.SoftShadowRays)
	}
	return 0
}

// IsShadowed returns true if p is in a shadow from the given light
func (w *World) IsShadowed(p Point, lp Point, xs Intersections) bool {
	v := lp.SubPoint(p)
	distance, direction := v.MagnitudeNormalize()

	r := NewRay(p, direction)

	intersections := w.Intersections(r, xs)

	// intersections must be sorted!
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

type pixel struct {
	x, y float64
}

// Render is the work done by the renderWorker, renders one pixel
func (p *pixel) Render(w *World, canvas *Canvas, xs Intersections, offset, l float64, rng *rand.Rand) {
	clrs := Colors{}

	// Collect colors for each sub-pixel and average them (antialias), slow and naive implementation
	for sx := 1.0; sx < l+1; sx++ {
		for sy := 1.0; sy < l+1; sy++ {
			a := p.x + offset*(sx*2-1)
			b := p.y + offset*(sy*2-1)

			ray := w.Camera().RayForPixel(a, b)
			clr := w.ColorAt(ray, w.Config.MaxRecusions, xs, rng)
			clrs = append(clrs, clr)
		}
	}

	// There is a race condition here, as canvas is also read by the GPU
	// Only true when using GPU to display the render live.
	canvas.SetFloat(p.x, p.y, clrs.Average())
}

// renderWorker processes a single pixel at a time
func (w *World) renderWorker(in chan *pixel, canvas *Canvas) {
	// One intersections list per worker, making these per pixel is very expensive
	xs := NewIntersections()

	// antialias config
	aa := float64(w.Config.Antialias)
	numSquares := 1.0
	offset := 0.5

	if aa > 0 {
		numSquares = math.Pow(2, aa)
		offset = 1.0 / (2 * aa)
	}
	rowLength := math.Sqrt(numSquares)

	rng := rand.New(rand.NewSource(time.Now().Unix()))

	for pixel := range in {
		// render the pixel
		pixel.Render(w, canvas, xs, offset, rowLength, rng)
		// clear intersections for next pixel
		xs = xs[:0]
	}
}

func (w *World) doRender(camera *Camera, canvas *Canvas) *Canvas {

	log.Println("Running render...")

	// allow this many renders to run at once
	max := w.Config.Parallelism

	// create communications channel
	pending := make(chan *pixel)
	var wg sync.WaitGroup

	// start the render goroutines
	for i := 0; i < max; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w.renderWorker(pending, canvas)
		}()
	}

	total := (camera.Vsize - 1) * (camera.Hsize - 1)
	done := 0.0
	last := 0.0

	offset := float64(w.Config.RenderPasses)
	for i := 0.0; i < offset; i++ {
		for y := i; y < camera.Vsize; y = y + offset {
			for x := 0.0; x < camera.Hsize; x++ {
				// send work to workers
				pending <- &pixel{x: x, y: y}
				done++
				last = showProgress(total, done, last)
			}
		}
	}
	close(pending)
	wg.Wait()

	log.Print("Render finished!")
	return canvas
}

// ShowInfo dumps info about the world
func (w *World) ShowInfo() {
	haveAreaLight := false

	for _, l := range w.Lights {
		switch l.(type) {
		case *AreaLight, *AreaSpotLight:
			haveAreaLight = true
		}
	}

	log.Printf("Camera HSize: %v", w.Camera().Hsize)
	log.Printf("Camera WSize: %v", w.Camera().Vsize)
	log.Printf("Camera Pixel Size: %.4f", w.Camera().PixelSize)
	log.Printf("Camera Half With: %.4f", w.Camera().HalfWidth)
	log.Printf("Camera Half Height: %.4f", w.Camera().HalfHeight)
	log.Printf("Antialiasing: %v", w.Config.Antialias)
	log.Printf("Parallelism: %v", w.Config.Parallelism)
	log.Printf("Max Recursion: %v", w.Config.MaxRecusions)
	log.Printf("Have Area Lights? -> %v", haveAreaLight)
	if haveAreaLight {
		log.Printf("  Area Light Rays: %v", w.Config.AreaLightRays)
	}
	log.Printf("Soft Shadows enabled? -> %v", w.Config.SoftShadows)
	if w.Config.SoftShadows {
		log.Printf("  Soft shadow rays: %v", w.Config.SoftShadowRays)
	}
}

// Render renders the world using the world camera
func (w *World) Render(camera *Camera, canvas *Canvas) {
	w.LintWorld()
	w.PrecomputeValues()

	w.ShowInfo()

	w.doRender(camera, canvas)
}
