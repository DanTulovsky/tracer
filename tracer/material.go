package tracer

import (
	"image"
	"log"
	"math"

	"golang.org/x/image/colornames"
)

// Material is a material to apply to shapes
type Material struct {
	Color                                                                            Color
	Pattern                                                                          Patterner
	Ambient, Diffuse, Specular, Shininess, Reflective, Transparency, RefractiveIndex float64

	// This material emits light
	Emissive Color

	// Some objects should not cast shadows (e.g. water in a pond)
	ShadowCaster bool

	// Some materials have textures associated with them, this is an image file read in and stored as a canvas
	Texture *Canvas

	// Used to apply perturbations to the material (changes in the normal vector)
	// Use this for BumpMaps
	perturber Perturber
}

// NewMaterial returns a new material
func NewMaterial(clr Color, a, d, sp, s, r, t, ri float64, p Perturber) *Material {

	return &Material{
		Color:           clr,
		Ambient:         a,
		Diffuse:         d,
		Specular:        sp,
		Shininess:       s,
		Reflective:      r,
		Transparency:    t,
		RefractiveIndex: ri,
		ShadowCaster:    true,
		perturber:       p,
	}
}

// NewDefaultMaterial returns a default material
func NewDefaultMaterial() *Material {

	return &Material{
		Color:           NewColor(1, 1, 1),
		Ambient:         0.1,
		Diffuse:         0.9,
		Specular:        0.9,
		Shininess:       200.0,
		Reflective:      0,
		Transparency:    0,
		RefractiveIndex: 1.0,
		ShadowCaster:    true,
	}
}

// NewDefaultGlassMaterial returns a default glass material
func NewDefaultGlassMaterial() *Material {

	return &Material{
		Color:           NewColor(1, 1, 1),
		Ambient:         0.1,
		Diffuse:         0.9,
		Specular:        0.9,
		Shininess:       200.0,
		Reflective:      0.0,
		Transparency:    1.0,
		RefractiveIndex: 1.5,
		ShadowCaster:    false,
	}
}

// PerturbNormal applies the material perturbation function to the normal n at point p
func (m *Material) PerturbNormal(normal Vector, p Point) Vector {
	if m.perturber != nil {
		return m.perturber.Perturb(normal, p)
	}
	return normal
}

// ColorAtTexture returns the color at the u,v point based on the texture attached to the material
// Only works for SmoothTriangles, used during obj import
func (m *Material) ColorAtTexture(o Shaper, u, v float64, base Color) Color {
	if m.Texture == nil {
		return ColorName(colornames.Purple) // highly visible, texture emissing
	}

	t := o.(*SmoothTriangle)

	w := 1 - u - v
	x := (u*t.VT2.x + v*t.VT3.x + w*t.VT1.x) * float64((m.Texture.Width - 1))
	y := (u*t.VT2.y + v*t.VT3.y + w*t.VT1.y) * float64((m.Texture.Height - 1))

	// wrap textures around if needed
	if x < 0 {
		x = float64(m.Texture.Width-1) + math.Mod(x, float64(m.Texture.Width-1))
	}
	if y < 0 {
		y = float64(m.Texture.Height-1) + math.Mod(y, float64(m.Texture.Height-1))
	}

	clr, err := m.Texture.Get(int(x), int(y))
	if err != nil {
		log.Println(err)
		return ColorName(colornames.Purple) // highly visible, texture missing
	}

	// - Kd - material diffuse is multiplied by the texture value
	return clr.Blend(base)
}

// AddDiffuseTexture adds a texture mapped to a Canvas
func (m *Material) AddDiffuseTexture(name string, i image.Image) error {
	log.Println("converting image (texture) to canvas...")
	canvas := imageToCanvas(i)

	m.Texture = canvas

	return nil
}

// HasPattern returns true if a material has a pattern attached to it
func (m *Material) HasPattern() bool {
	return m.Pattern != nil
}

// HasTexture returns true if a material has a texture attached to it
func (m *Material) HasTexture() bool {
	return m.Texture != nil
}

// SetPattern sets a pattern on a material
func (m *Material) SetPattern(p Patterner) {
	m.Pattern = p
}

// SetPerturber sets a perturber on a material
func (m *Material) SetPerturber(p Perturber) {
	m.perturber = p
}

// Equals return true if the materials are the same
func (m *Material) Equals(m2 *Material) bool {
	return m.Color.Equal(m2.Color) &&
		m.Pattern == m2.Pattern &&
		m.Ambient == m2.Ambient &&
		m.Diffuse == m2.Diffuse &&
		m.Specular == m2.Specular &&
		m.Shininess == m2.Shininess &&
		m.Reflective == m2.Reflective &&
		m.Transparency == m2.Transparency &&
		m.RefractiveIndex == m2.RefractiveIndex &&
		m.Emissive.Equal(m2.Emissive) &&
		m.ShadowCaster == m2.ShadowCaster &&
		m.Texture == m2.Texture &&
		m.perturber == m2.perturber
}
