package tracer

import (
	"github.com/DanTulovsky/tracer/utils"
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
	n        opensimplex.Noise
	maxNoise float64
}

// NewNoisePerturber returns a perturber that makes waves on the shape
func NewNoisePerturber() *NoisePerturber {
	return &NoisePerturber{
		// n:        opensimplex.NewNormalized(time.Now().Unix()),
		n:        opensimplex.NewNormalized(1),
		maxNoise: 1,
	}
}

// Perturb implements the Perturber interface
func (p *NoisePerturber) Perturb(n Vector, pt Point) Vector {

	// n = n.Normalize()
	epsilon := 0.0001
	f0 := p.n.Eval3(pt.X(), pt.Y(), pt.Z()) * p.maxNoise
	fx := p.n.Eval3(pt.X()+epsilon, pt.Y(), pt.Z()) * p.maxNoise
	fy := p.n.Eval3(pt.X(), pt.Y()+epsilon, pt.Z()) * p.maxNoise
	fz := p.n.Eval3(pt.X(), pt.Y(), pt.Z()+epsilon) * p.maxNoise

	f0 = utils.AT(f0, 0, 1, -1, 1)
	fx = utils.AT(fx, 0, 1, -1, 1)
	fy = utils.AT(fy, 0, 1, -1, 1)
	fz = utils.AT(fz, 0, 1, -1, 1)

	df := NewVector((fx-f0)/epsilon, (fy-f0)/epsilon, (fz-f0)/epsilon)
	// log.Println(f0)
	// log.Println(fx)
	// log.Println(fy)
	// log.Println(fz)
	// log.Println(df)
	// log.Println()

	return n.SubVector(df).Normalize()
}
