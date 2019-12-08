package tracer

import (
	"image"
	"log"
	"math"
	"os"
	"time"

	"github.com/ojrac/opensimplex-go"
)

// Perturber defines an object that can be used to apply perturbation to the normal vectors of shapes
type Perturber interface {
	Perturb(Vector, Point) Vector

	Transform() Matrix
	TransformInverse() Matrix
	SetTransform(Matrix)
}

type basePerturb struct {
	transform        Matrix
	transformInverse Matrix
}

// Transform returns the pattern transform
func (bp *basePerturb) Transform() Matrix {
	return bp.transform
}

// SetTransform sets the transform
func (bp *basePerturb) SetTransform(m Matrix) {
	bp.transform = m
	bp.transformInverse = m.Inverse()
}

// TransformInverse returns the inverse of the perturber transform
func (bp *basePerturb) TransformInverse() Matrix {
	return bp.transformInverse
}

// Perturb implements the Perturber interface
func (bp *basePerturb) Perturb(v Vector, p Point) Vector {
	panic("need to implement Perturb")
}

// NoisePerturber perturbes based on noise
type NoisePerturber struct {
	basePerturb

	n opensimplex.Noise

	// maxNoise generally controls the "vertical" (along the original normal) height
	maxNoise float64
}

// NewNoisePerturber returns a perturber that makes waves on the shape
func NewNoisePerturber(maxNoise float64) *NoisePerturber {
	return &NoisePerturber{
		n: opensimplex.NewNormalized(time.Now().Unix()),

		// These two parameters control the size and frequency of the bumps
		maxNoise: maxNoise, // 1 is a nice value here
		basePerturb: basePerturb{
			transform:        IM(),
			transformInverse: IM().Inverse(),
		},
	}
}

// Perturb implements the Perturber interface
func (np *NoisePerturber) Perturb(v Vector, p Point) Vector {
	return np.perturb(v, p.TimesMatrix(np.TransformInverse()))
}

// perturb is the local perturb function
func (np *NoisePerturber) perturb(n Vector, p Point) Vector {

	epsilon := 0.001
	f0 := np.n.Eval3(p.X(), p.Y(), p.Z()) * np.maxNoise
	fx := np.n.Eval3(p.X()+epsilon, p.Y(), p.Z()) * np.maxNoise
	fy := np.n.Eval3(p.X(), p.Y()+epsilon, p.Z()) * np.maxNoise
	fz := np.n.Eval3(p.X(), p.Y(), p.Z()+epsilon) * np.maxNoise

	df := NewVector((fx-f0)/epsilon, (fy-f0)/epsilon, (fz-f0)/epsilon)

	return n.SubVector(df).Normalize()
}

// SetNoise allows setting the noise generator (mostly used for tests)
func (np *NoisePerturber) SetNoise(n opensimplex.Noise) {
	np.n = n
}

// SinePerturber perturbes based on the Sine wave
type SinePerturber struct {
	basePerturb

	n opensimplex.Noise
}

// NewSinePerturber returns a perturber that makes waves on the shape
func NewSinePerturber() *SinePerturber {
	return &SinePerturber{
		basePerturb: basePerturb{
			transform:        IM(),
			transformInverse: IM().Inverse(),
		},
	}
}

// Perturb implements the Perturber interface
func (sp *SinePerturber) Perturb(v Vector, p Point) Vector {
	return sp.perturb(v, p.TimesMatrix(sp.TransformInverse()))
}

// Perturb implements the Perturber interface
func (sp *SinePerturber) perturb(n Vector, p Point) Vector {
	epsilon := 0.001

	f0 := math.Sin(p.Y())
	fy := math.Sin(p.Y() + epsilon)

	df := NewVector(0, (fy-f0)/epsilon, 0)

	return n.AddVector(df).Normalize()
}

// ImageHeightmapPerturber perturbes based on a height image
type ImageHeightmapPerturber struct {
	basePerturb

	canvas *Canvas

	// convert (x,y,z) -> (u,v)
	mapper Mapper
}

// NewImageHeightmapPerturber returns a perturber that uses an image to simulate bumps
func NewImageHeightmapPerturber(filename string, mapper Mapper) (*ImageHeightmapPerturber, error) {
	// read in image
	reader, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Decode image
	m, _, err := image.Decode(reader)
	if err != nil {
		return nil, err
	}

	// convert to canvas
	canvas := imageToCanvas(m)

	p := &ImageHeightmapPerturber{
		canvas: canvas,
		mapper: mapper,
		basePerturb: basePerturb{
			transform:        IM(),
			transformInverse: IM().Inverse(),
		},
	}

	return p, nil
}

// UVColorAt returns the color at the 2D coordinate (u, v)
func (ip *ImageHeightmapPerturber) UVColorAt(u, v float64) Color {
	// v = 1 - v
	// u = 1 - u

	x := u * (float64(ip.canvas.Width) - 1)
	y := v * (float64(ip.canvas.Height) - 1)

	c, err := ip.canvas.Get(int(x), int(y))
	if err != nil {
		log.Fatal(err)
	}
	return c
}

// Perturb implements the Perturber interface
func (ip *ImageHeightmapPerturber) Perturb(normal Vector, p Point) Vector {
	return ip.perturb(normal, p.TimesMatrix(ip.TransformInverse()))
}

// perturb is the local perturb function
func (ip *ImageHeightmapPerturber) perturb(n Vector, p Point) Vector {

	u, v := ip.mapper.Map(p)
	epsilon := 0.001

	north := ip.UVColorAt(u, math.Min(v+epsilon, 1))
	northwest := ip.UVColorAt(math.Max(u-epsilon, 0), math.Min(v+epsilon, 1))
	west := ip.UVColorAt(math.Max(u-epsilon, 0), v)
	southwest := ip.UVColorAt(math.Max(u-epsilon, 0), math.Max(v-epsilon, 0))
	south := ip.UVColorAt(u, math.Max(v-epsilon, 0))
	southeast := ip.UVColorAt(math.Min(u+epsilon, 1), math.Max(v-epsilon, 0))
	east := ip.UVColorAt(math.Min(u+epsilon, 1), v)
	northeast := ip.UVColorAt(math.Min(u+epsilon, 1), math.Min(v+epsilon, 1))

	// gradient vector
	dydx := ((northeast.GreyScale() + 2*east.GreyScale() + southeast.GreyScale()) - (northwest.GreyScale() + 2*west.GreyScale() + southwest.GreyScale())) / 2
	dydz := ((southwest.GreyScale() + 2*south.GreyScale() + southeast.GreyScale()) - (northwest.GreyScale() + 2*north.GreyScale() + northeast.GreyScale())) / 2

	new := NewVector(dydx, 0, dydz).Normalize()

	// return new
	return n.SubVector(new)
}
