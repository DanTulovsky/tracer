package tracer

import (
	"math"
	"time"

	"github.com/ojrac/opensimplex-go"
)

// Patterner is a pattern that can be applied to a material
// This does not do UV Mapping, the pattern is set on the 3D Space of the object
type Patterner interface {
	// Returns the correct color at the given point on the given object for this pattern
	ColorAtObject(Shaper, Point) Color

	Transform() Matrix
	TransformInverse() Matrix
	SetTransform(Matrix)
}

// basePattern is the base pattern for others
type basePattern struct {
	transform        Matrix
	transformInverse Matrix
}

// Transform returns the pattern transform
func (bp *basePattern) Transform() Matrix {
	return bp.transform
}

// SetTransform sets the transform
func (bp *basePattern) SetTransform(m Matrix) {
	bp.transform = m
	bp.transformInverse = m.Inverse()
}

// TransformInverse returns the inverse of the pattern transform
func (bp *basePattern) TransformInverse() Matrix {
	return bp.transformInverse
}

// ColorAt implements Patterner
func (bp *basePattern) ColorAtObject(o Shaper, p Point) Color {
	panic("need to implement ColorAt")
}

// objectSpacePoint returns the world point as object point
func (bp *basePattern) objectSpacePoint(o Shaper, p Point) Point {
	op := p.ToObjectSpace(o)
	return op.TimesMatrix(bp.TransformInverse())
}

// CubeMapPattern ...
type CubeMapPattern struct {
	basePattern
	// The uv patterner to use for each face
	left, front, right, back, up, down UVPatterner
	mapper                             Mapper
}

// NewCubeMapPattern maps a texture onto a cube
func NewCubeMapPattern(left, front, right, back, up, down UVPatterner) *CubeMapPattern {
	return &CubeMapPattern{
		left:   left,
		front:  front,
		right:  right,
		back:   back,
		up:     up,
		down:   down,
		mapper: NewCubeMap(left, front, right, back, up, down),
		basePattern: basePattern{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
		},
	}
}

// NewCubeMapPatternSame maps a texture onto a cube
func NewCubeMapPatternSame(p UVPatterner) *CubeMapPattern {
	return &CubeMapPattern{
		left:   p,
		front:  p,
		right:  p,
		back:   p,
		up:     p,
		down:   p,
		mapper: NewCubeMapSame(p),
		basePattern: basePattern{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
		},
	}
}

// faceFromPoint returns the face that the point is on
func (cm *CubeMapPattern) faceFromPoint(p Point) cubeFace {
	coord := math.Max(math.Max(math.Abs(p.x), math.Abs(p.y)), math.Abs(p.z))

	switch coord {
	case p.x:
		return cubeFaceRight
	case -p.x:
		return cubeFaceLeft
	case p.y:
		return cubeFaceUp
	case -p.y:
		return cubeFaceDown
	case p.z:
		return cubeFaceFront
	}
	return cubeFaceBack
}

// ColorAtObject returns the color for the given pattern on the given object
func (cm *CubeMapPattern) ColorAtObject(o Shaper, p Point) Color {
	return cm.colorAt(cm.objectSpacePoint(o, p))
}

// uvColorAt returns the color at the 2D coordinate (u, v)
func (cm *CubeMapPattern) uvColorAt(u, v float64) Color {
	return White()
}

// ColorAt implements Patterner
func (cm *CubeMapPattern) colorAt(p Point) Color {
	// The correct face is calculated by this function
	u, v := cm.mapper.Map(p)

	// find the face the point is on
	face := cm.faceFromPoint(p)

	switch face {
	case cubeFaceFront:
		return cm.front.UVColorAt(u, v)
	case cubeFaceBack:
		return cm.back.UVColorAt(u, v)
	case cubeFaceLeft:
		return cm.left.UVColorAt(u, v)
	case cubeFaceRight:
		return cm.right.UVColorAt(u, v)
	case cubeFaceUp:
		return cm.up.UVColorAt(u, v)
	case cubeFaceDown:
		return cm.down.UVColorAt(u, v)
	}

	// should never happen
	return Black()
}

// TextureMapPattern maps the child pattern using the provided map
type TextureMapPattern struct {
	basePattern
	pattern UVPatterner
	mapper  Mapper
}

// NewTextureMapPattern returns a new texture map pattern
// For cubes, use NewCubeMapPattern, this one works for planes, spheres, cylinders and cones
func NewTextureMapPattern(p UVPatterner, m Mapper) *TextureMapPattern {
	return &TextureMapPattern{
		pattern: p,
		mapper:  m,
		basePattern: basePattern{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
		},
	}
}

// ColorAtObject returns the color for the given pattern on the given object
func (tmp *TextureMapPattern) ColorAtObject(o Shaper, p Point) Color {
	return tmp.colorAt(tmp.objectSpacePoint(o, p))
}

// ColorAt implements Patterner
func (tmp *TextureMapPattern) colorAt(p Point) Color {
	u, v := tmp.mapper.Map(p)
	return tmp.pattern.UVColorAt(u, v)
}

// StripedPattern is a pattern that overlays stripes
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
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
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
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
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
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
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
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
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

// PertrubedPattern jitters the points before passing them on to the real pattern
// Using Opensimplex: https://github.com/ojrac/opensimplex-go
type PertrubedPattern struct {
	basePattern
	p        Patterner // real pattern to delegate to
	noise    opensimplex.Noise
	maxNoise float64
}

// NewPertrubedPattern returns a new pertrubed patterner
// maxNoise is a [0, 1] value which clamps how much the noise affects the input
func NewPertrubedPattern(p Patterner, maxNoise float64) *PertrubedPattern {

	if maxNoise < 0 || maxNoise > 1 {
		panic("maxNoise must be between 0 and 1")
	}

	n := opensimplex.NewNormalized(time.Now().Unix())

	return &PertrubedPattern{
		p:        p,
		noise:    n,
		maxNoise: maxNoise,
		basePattern: basePattern{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
		},
	}
}

// ColorAtObject returns the color for the given pattern on the given object
func (pp *PertrubedPattern) ColorAtObject(o Shaper, p Point) Color {
	// change p using opensimplex
	n := pp.noise.Eval3(p.X(), p.Y(), p.Z()) * pp.maxNoise

	// pass it to the real patterner
	return pp.p.ColorAtObject(o, p.AddScalar(n))
}

// BlendedPattern blends the output of two patterns
type BlendedPattern struct {
	basePattern
	p1 Patterner // real pattern to delegate to
	p2 Patterner // real pattern to delegate to
}

// NewBlendedPattern returns a new blended patterner
func NewBlendedPattern(p1, p2 Patterner) *BlendedPattern {

	return &BlendedPattern{
		p1: p1,
		p2: p2,
		basePattern: basePattern{
			transform:        IdentityMatrix(),
			transformInverse: IdentityMatrix().Inverse(),
		},
	}
}

// ColorAtObject returns the color for the given pattern on the given object
func (bp *BlendedPattern) ColorAtObject(o Shaper, p Point) Color {

	c1 := bp.p1.ColorAtObject(o, p)
	c2 := bp.p2.ColorAtObject(o, p)

	// blend them together
	return c1.Blend(c2)
}
