package tracer

import (
	"math"

	"golang.org/x/image/colornames"
)

// Light is the interface all the lights use
type Light interface {
	Intensity() Color
	Position() Point
	Shape() Shaper
}

// Lights is a collection of Light
type Lights []Light

// PointLight implements the light interface and is a single point light with no size
type PointLight struct {
	position  Point
	intensity Color
}

// AreaLight shines in all directions and is a shape
type AreaLight struct {
	Shaper
	intensity Color
}

// NewAreaLight returns a new area light
func NewAreaLight(s Shaper, i Color) *AreaLight {
	s.Material().Emissive = i
	s.Material().ShadowCaster = false
	s.Material().Diffuse = 0
	s.Material().Specular = 0
	s.Material().Reflective = 0
	s.Material().Transparency = 1

	return &AreaLight{
		Shaper:    s,
		intensity: i,
	}
}

// Intensity implements the Light interface
func (al *AreaLight) Intensity() Color {
	return al.intensity
}

// Position implements the Light interface
func (al *AreaLight) Position() Point {
	// TODO: Fix this
	p := al.Shaper.Bounds().Center().ToWorldSpace(al.Shaper)
	return p
}

// Shape returns the Shaper object of this light
func (al *AreaLight) Shape() Shaper {
	return al.Shaper
}

// NewPointLight returns a new point light
func NewPointLight(p Point, i Color) *PointLight {
	return &PointLight{position: p, intensity: i}
}

// Intensity returns the intensity of the light
func (pl *PointLight) Intensity() Color {
	return pl.intensity
}

// Position returns the position of the light
func (pl *PointLight) Position() Point {
	return pl.position
}

// Shape returns the Shaper object of this light, point lights have no shape
func (pl *PointLight) Shape() Shaper {
	return nil
}

// lighting returns the color for a given point
func lighting(m *Material, o Shaper, p Point, l Light, eye, normal Vector, inShadow bool, u, v float64) Color {
	var ambient, diffuse, specular Color

	var clr Color

	switch {
	case m.HasPattern():
		clr = m.Pattern.ColorAtObject(o, p)
	default:
		clr = m.Color
	}

	switch {
	case m.HasTexture():
		clr = m.ColorAtTexture(o, u, v)
	}

	// combine surface color with light's color/intensity
	effectiveColor := clr.Blend(l.Intensity())

	// compute ambient contribution
	ambient = effectiveColor.Scale(m.Ambient)

	// Compute the emissive contribution
	emissive := m.Emissive

	// light not visible, ignore diffuse and specular components
	if inShadow {
		return ambient.Add(emissive)
	}

	// find the direction to the light source
	lightv := l.Position().SubPoint(p).Normalize()

	// lightDotNormal represents the cosine of the angle between the light vector and the normal vector
	// a negaive number means the light is on the other side of the surface
	lightDotNormal := lightv.Dot(normal)

	if lightDotNormal < 0 {
		diffuse = ColorName(colornames.Black)
		specular = ColorName(colornames.Black)
	} else {
		// compute the diffuse contribution
		diffuse = effectiveColor.Scale(m.Diffuse).Scale(lightDotNormal)

		// reflectDotEye represens the cosine of the angle between the relfection vector and the eye vector
		// a negative number means the light reflects away from the eye
		reflectv := lightv.Negate().Reflect(normal)
		reflectDotEye := reflectv.Dot(eye)

		if reflectDotEye <= 0 {
			specular = ColorName(colornames.Black)
		} else {
			// compute the specular contrbution
			factor := math.Pow(reflectDotEye, m.Shininess)
			specular = l.Intensity().Scale(m.Specular).Scale(factor)
		}
	}

	return emissive.Add(ambient).Add(diffuse).Add(specular)
}

// ColorAtPoint returns the clamped color at the given point
func ColorAtPoint(m *Material, o Shaper, p Point, l Light, eye, normal Vector, inShadow bool) Color {
	return lighting(m, o, p, l, eye, normal, inShadow, 0, 0).Clamp()
}
