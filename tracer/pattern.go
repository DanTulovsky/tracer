package tracer

import "math"

// Patterner is a pattern that can be applied to a material
type Patterner interface {
	// Returns the correct color at the given point on the given object for this pattern
	ColorAtObject(Shaper, Point) Color

	Transform() Matrix
	SetTransform(Matrix)
}

// basePattern is the base pattern for others
type basePattern struct {
	transform Matrix
}

// Transform returns the pattern transform
func (bp *basePattern) Transform() Matrix {
	return bp.transform
}

// SetTransform sets the transform
func (bp *basePattern) SetTransform(m Matrix) {
	bp.transform = m
}

// ColorAt implements Patterner
func (bp *basePattern) ColorAtObject(p Point) Color {
	panic("need to implement ColorAt")
}

// StripedPattern is a patternt at overlays stripes
type StripedPattern struct {
	basePattern
	a, b Color
}

// NewStripedPattern returns a new striped pattern with the given colors
func NewStripedPattern(c1, c2 Color) *StripedPattern {
	return &StripedPattern{
		a: c1,
		b: c2,
		basePattern: basePattern{
			transform: IdentityMatrix(),
		},
	}
}

// ColorAtObject returns the color for the given pattern on the given object
func (sp *StripedPattern) ColorAtObject(o Shaper, p Point) Color {
	// move point to object space
	op := p.TimesMatrix(o.Transform().Inverse())
	pp := op.TimesMatrix(sp.Transform().Inverse())

	return sp.colorAt(pp)
}

// ColorAt implements Patterner
func (sp *StripedPattern) colorAt(p Point) Color {
	if int(math.Floor(p.X()))%2 == 0 {
		return sp.a
	}
	return sp.b
}
