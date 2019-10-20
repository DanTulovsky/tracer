package tracer

import "math"

// Patterner is a pattern that can be applied to a material
// This does not do UV Mapping, the pattern is set on the 3D Space of the object
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
func (bp *basePattern) ColorAtObject(o Shaper, p Point) Color {
	panic("need to implement ColorAt")
}

// objectSpacePoint returns the world point as object point
func (bp *basePattern) objectSpacePoint(o Shaper, p Point) Point {
	op := p.TimesMatrix(o.Transform().Inverse())
	return op.TimesMatrix(bp.Transform().Inverse())
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
	return sp.colorAt(sp.objectSpacePoint(o, p))
}

// ColorAt implements Patterner
func (sp *StripedPattern) colorAt(p Point) Color {
	if int(math.Floor(p.X()))%2 == 0 {
		return sp.a
	}
	return sp.b
}

// GradientPattern implements a gradient pattern
type GradientPattern struct {
	basePattern
	a, b Color
}

// NewGradientPattern returns a new gradient pattern with the given colors
func NewGradientPattern(c1, c2 Color) *GradientPattern {
	return &GradientPattern{
		a: c1,
		b: c2,
		basePattern: basePattern{
			transform: IdentityMatrix(),
		},
	}
}

// ColorAtObject returns the color for the given pattern on the given object
func (gp *GradientPattern) ColorAtObject(o Shaper, p Point) Color {
	return gp.colorAt(gp.objectSpacePoint(o, p))
}

// ColorAt implements Patterner
func (gp *GradientPattern) colorAt(p Point) Color {
	d := gp.b.Sub(gp.a)
	f := p.X() - math.Floor(p.X())

	return gp.a.Add(d.Scale(f))
}

// RingPattern implements a ring pattern
type RingPattern struct {
	basePattern
	a, b Color
}

// NewRingPattern returns a new ring pattern with the given colors
func NewRingPattern(c1, c2 Color) *RingPattern {
	return &RingPattern{
		a: c1,
		b: c2,
		basePattern: basePattern{
			transform: IdentityMatrix(),
		},
	}
}

// ColorAtObject returns the color for the given pattern on the given object
func (rp *RingPattern) ColorAtObject(o Shaper, p Point) Color {
	return rp.colorAt(rp.objectSpacePoint(o, p))
}

// ColorAt implements Patterner
func (rp *RingPattern) colorAt(p Point) Color {
	if int(math.Floor(math.Sqrt(math.Pow(p.X(), 2)+math.Pow(p.Z(), 2))))%2 == 0 {
		return rp.a
	}
	return rp.b
}

// CheckerPattern implements a checker pattern
type CheckerPattern struct {
	basePattern
	a, b Color
}

// NewCheckerPattern returns a new checker pattern with the given colors
func NewCheckerPattern(c1, c2 Color) *CheckerPattern {
	return &CheckerPattern{
		a: c1,
		b: c2,
		basePattern: basePattern{
			transform: IdentityMatrix(),
		},
	}
}

// ColorAtObject returns the color for the given pattern on the given object
func (cp *CheckerPattern) ColorAtObject(o Shaper, p Point) Color {
	return cp.colorAt(cp.objectSpacePoint(o, p))
}

// ColorAt implements Patterner
func (cp *CheckerPattern) colorAt(p Point) Color {
	if int(math.Floor(p.X())+math.Floor(p.Y())+math.Floor(p.Z()))%2 == 0 {
		return cp.a
	}
	return cp.b
}
