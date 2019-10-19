package tracer

import (
	"math"

	"golang.org/x/exp/shiny/materialdesign/colornames"
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

// NewPointLight returns a nw point light
func NewPointLight(p Point, i Color) PointLight {
	return PointLight{position: p, intensity: i}
}

// Intensity returns the intensity of the light
func (pl PointLight) Intensity() Color {
	return pl.intensity
}

// Position returns the position of the light
func (pl PointLight) Position() Point {
	return pl.position
}

// lighting returns the color for a given point
func lighting(m *Material, p Point, l Light, eye, normal Vector) Color {
	var ambient, diffuse, specular Color

	// combine surface color with light's color/intensity
	effectiveColor := m.Color.Blend(l.Intensity())

	// find the direction to the light source
	lightv := l.Position().SubPoint(p).Normalize()

	// compute ambient contribution
	ambient = effectiveColor.Scale(m.Ambient)

	// lightDotNormal represnets the cosine of the angle btween the light vector and the normal vector
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

	return ambient.Add(diffuse).Add(specular)
}

// ColorAtPoint returns the clamped color at the givn point
func ColorAtPoint(m *Material, p Point, l PointLight, eye, normal Vector) Color {
	return lighting(m, p, l, eye, normal).Clamp()
}
