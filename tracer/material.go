package tracer

// Material is a material to apply to shapes
type Material struct {
	Color                                                                            Color
	Pattern                                                                          Patterner
	Ambient, Diffuse, Specular, Shininess, Reflective, Transparency, RefractiveIndex float64
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
		Reflective:      0,
		Transparency:    1.0,
		RefractiveIndex: 1.5,
	}
}

// HasPattern returns true if a material has a pattern attached to it
func (m *Material) HasPattern() bool {
	return m.Pattern != nil
}

// SetPattern sets a pattern on a material
func (m *Material) SetPattern(p Patterner) {
	m.Pattern = p
}
