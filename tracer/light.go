package tracer

import (
	"math"

	"golang.org/x/image/colornames"
)

// Light is the interface all the lights use
type Light interface {
	Intensity() Color
	Position() Point
}

// PointLight implements the light interface and is a single point light with no size
type PointLight struct {
	position  Point
	intensity Color
}

// AreaLight shines in all directions, but is a sphere with size radius
type AreaLight struct {
	radius    float64 // size of the area light
	center    Point
	intensity Color
}

// NewAreaLight returns a new area light
func NewAreaLight(p Point, i Color, r float64) *AreaLight {
	return &AreaLight{}
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

	// find the direction to the light source
	lightv := l.Position().SubPoint(p).Normalize()

	// compute ambient contribution
	ambient = effectiveColor.Scale(m.Ambient)

	// Compute the emissive contribution
	emissive := m.Emissive

	// light not visible, ignore diffuse and specular components
	if inShadow {
		return ambient.Add(emissive)
	}

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
