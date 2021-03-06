package tracer

import "runtime"

// WorldConfig collects various settings to configure the world
type WorldConfig struct {
	// How many times to allow the ray to bounce between two objects (controls reflections of reflections)
	MaxRecusions int

	// Antialiasing support
	Antialias int

	// Parallelism, how many pixels to render at the same time
	Parallelism int

	// SoftShadow enables soft shadows
	SoftShadows bool

	// SoftShadowRays specifies how many shadow rays to cast, also used by area lights
	SoftShadowRays int

	// AreaLightRays specifies how many rays to cast for area lights
	AreaLightRays int

	// RenderPasses controls the display of the render to the screen
	RenderPasses int

	// BackfaceCulling disables drawing riangles facing away from the camera
	BackfaceCulling bool
}

// NewWorldConfig returns a new world config with default settings
func NewWorldConfig() *WorldConfig {
	return &WorldConfig{
		Antialias:       0,
		AreaLightRays:   10,
		MaxRecusions:    4,
		Parallelism:     runtime.NumCPU(),
		SoftShadows:     true,
		SoftShadowRays:  6,
		RenderPasses:    8,
		BackfaceCulling: false, // off by default, as transpaencies require it
	}
}
