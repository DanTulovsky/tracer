package tracer

import (
	"github.com/ojrac/opensimplex-go"
)

// Perturber defines an object that can be used to apply perturbation to the normal vectors of shapes
type Perturber interface {
	Perturb(Vector, Point) Vector
}

// DefaultPerturber is a no-op petruber
type DefaultPerturber struct {
}

// NewDefaultPerturber returns a default, no-op, petruber
func NewDefaultPerturber() *DefaultPerturber {
	return &DefaultPerturber{}
}

// Perturb implements the Perturber interface
func (p *DefaultPerturber) Perturb(v Vector, unused Point) Vector {
	return v
}

// NoisePerturber perturbes based on noise
type NoisePerturber struct {
	n opensimplex.Noise

	// maxNoise generally controls the "vertical" (along the original normal) height
	maxNoise float64

	// scale controls the "horizontal" size of the pattern
	// a larger value here means more bumps are visible
	scale float64
}

// NewNoisePerturber returns a perturber that makes waves on the shape
func NewNoisePerturber(maxNoise, scale float64) *NoisePerturber {
	return &NoisePerturber{
		// n:        opensimplex.NewNormalized(time.Now().Unix()),
		n: opensimplex.NewNormalized(1),

		// These two parameters control the size and frequency of the bumps
		maxNoise: maxNoise, // 1 is a nice value here
		scale:    scale,    // 6 is a nice value for a -1,1 shape
	}
}

// Perturb implements the Perturber interface
func (p *NoisePerturber) Perturb(n Vector, pt Point) Vector {

	pt = pt.Scale(p.scale)
	epsilon := 0.001
	f0 := p.n.Eval3(pt.X(), pt.Y(), pt.Z()) * p.maxNoise
	fx := p.n.Eval3(pt.X()+epsilon, pt.Y(), pt.Z()) * p.maxNoise
	fy := p.n.Eval3(pt.X(), pt.Y()+epsilon, pt.Z()) * p.maxNoise
	fz := p.n.Eval3(pt.X(), pt.Y(), pt.Z()+epsilon) * p.maxNoise

	df := NewVector((fx-f0)/epsilon, (fy-f0)/epsilon, (fz-f0)/epsilon)

	return n.SubVector(df).Normalize()
}

// SetNoise allows setting the noise generator (mostly used for tests)
func (p *NoisePerturber) SetNoise(n opensimplex.Noise) {
	p.n = n
}
