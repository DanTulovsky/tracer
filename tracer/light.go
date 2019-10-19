package tracer

import "image/color"

// Light is the interface all the lights use
type Light interface {
}

// PointLight implements the light interface and is a single point light with no size
type PointLight struct {
	Position  Point
	Intensity color.Color
}

// NewPointLight returns a nw point light
func NewPointLight(p Point, i color.Color) PointLight {
	return PointLight{Position: p, Intensity: i}
}
