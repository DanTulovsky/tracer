package tracer

import (
	"image"
	"log"

	"golang.org/x/image/colornames"

	"github.com/google/go-cmp/cmp"
)

// Material is a material to apply to shapes
type Material struct {
	Color                                                                            Color
	Pattern                                                                          Patterner
	Ambient, Diffuse, Specular, Shininess, Reflective, Transparency, RefractiveIndex float64

	// Some objects should not cast shadows (e.g. water in a pond)
	ShadowCaster bool

	// Some materials have textures associated with them, this is an image file read in and stored as a canvas
	Texture *Canvas
}

// NewMaterial returns a new material
func NewMaterial(clr Color, a, d, sp, s, r, t, ri float64) *Material {

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

// ColorAtTexture returns the color at the u,v point based on the texture attached to the material
func (m *Material) ColorAtTexture(o Shaper, u, v float64) Color {
	if m.Texture == nil {
		return ColorName(colornames.Purple) // highly visible, texture emissing
	}

	t := o.(*SmoothTriangle)

	w := 1 - u - v
	x := (u*t.VT2.x + v*t.VT3.x + w*t.VT1.x) * float64((m.Texture.Width - 1))
	y := (u*t.VT2.y + v*t.VT3.y + w*t.VT1.y) * float64((m.Texture.Height - 1))

	clr, err := m.Texture.Get(int(x), int(y))
	if err != nil {
		log.Println(err)
		return ColorName(colornames.Purple) // highly visible, texture emissing
	}

	return clr
}

// AddTexture adds a texture mapped to a Canvas
func (m *Material) AddTexture(name string, i image.Image) error {
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

// Equals return true if the materials are the same
func (m *Material) Equals(m2 *Material) bool {
	return cmp.Equal(m, m2)
}
