package tracer

import (
	"log"
	"math"
	"math/rand"

	"golang.org/x/image/colornames"
)

// Light is the interface all the lights use
type Light interface {
	Intensity() Color
	SetIntensity(Color)
	IsVisible() bool

	// Center position on the light
	Position() Point
	// Random postion on the light
	RandomPosition(*rand.Rand) Point
	Shape() Shaper
}

// Lights is a collection of Light
type Lights []Light

// AreaLight shines in all directions and is a shape
type AreaLight struct {
	Shaper
	intensity Color
	visible   bool
}

// NewAreaLight returns a new area light
func NewAreaLight(s Shaper, i Color, v bool) *AreaLight {
	s.Material().Emissive = i
	s.Material().ShadowCaster = false
	s.Material().Diffuse = 0
	s.Material().Specular = 0
	s.Material().Reflective = 0
	s.Material().Transparency = 1

	return &AreaLight{
		Shaper:    s,
		intensity: i,
		visible:   v,
	}
}

// Intensity implements the Light interface
func (al *AreaLight) Intensity() Color {
	return al.intensity
}

// Position implements the Light interface
func (al *AreaLight) Position() Point {
	return al.Shaper.Bounds().Center().ToWorldSpace(al.Shaper)
}

// RandomPosition implements the Light interface
func (al *AreaLight) RandomPosition(rng *rand.Rand) Point {
	return al.Shaper.RandomPosition(rng)
}

// SetIntensity sets the intensity of the light
func (al *AreaLight) SetIntensity(c Color) {
	al.intensity = c
}

// Shape returns the Shaper object of this light
func (al *AreaLight) Shape() Shaper {
	return al.Shaper
}

// IsVisible returns true if the  light shape should be visible.
func (al *AreaLight) IsVisible() bool {
	return al.visible
}

// PointLight implements the light interface and is a single point light with no size
type PointLight struct {
	position  Point
	intensity Color
}

// NewPointLight returns a new point light
func NewPointLight(p Point, i Color) *PointLight {
	return &PointLight{position: p, intensity: i}
}

// Intensity returns the intensity of the light
func (pl *PointLight) Intensity() Color {
	return pl.intensity
}

// SetIntensity sets the intensity of the light
func (pl *PointLight) SetIntensity(c Color) {
	pl.intensity = c
}

// Position returns the position of the light
func (pl *PointLight) Position() Point {
	return pl.position
}

// RandomPosition returns the position of the light
func (pl *PointLight) RandomPosition(rng *rand.Rand) Point {
	return pl.position
}

// Shape returns the Shaper object of this light, point lights have no shape
func (pl *PointLight) Shape() Shaper {
	return nil
}

// IsVisible returns true if light is visible, point lights are never visible.
func (pl *PointLight) IsVisible() bool {
	return false
}

// AreaSpotLight shines in a direction and is a shape
type AreaSpotLight struct {
	Shaper
	intensity Color
	visible   bool
	angle     float64
	direction Vector
}

// NewAreaSpotLight returns a new area light
func NewAreaSpotLight(s Shaper, i Color, v bool, angle float64, to Point) *AreaSpotLight {
	s.Material().Emissive = i
	s.Material().ShadowCaster = false
	s.Material().Diffuse = 0
	s.Material().Specular = 0
	s.Material().Reflective = 0
	s.Material().Transparency = 1

	asl := &AreaSpotLight{
		Shaper:    s,
		intensity: i,
		visible:   v,
		angle:     angle,
	}
	asl.direction = to.SubPoint(asl.Position()).Normalize()
	return asl
}

// Intensity implements the Light interface
func (al *AreaSpotLight) Intensity() Color {
	return al.intensity
}

// Position implements the Light interface
func (al *AreaSpotLight) Position() Point {
	return al.Shaper.Bounds().Center().ToWorldSpace(al.Shaper)
}

// RandomPosition implements the Light interface
func (al *AreaSpotLight) RandomPosition(rng *rand.Rand) Point {
	return al.Shaper.RandomPosition(rng)
}

// SetIntensity sets the intensity of the light
func (al *AreaSpotLight) SetIntensity(c Color) {
	al.intensity = c
}

// Shape returns the Shaper object of this light
func (al *AreaSpotLight) Shape() Shaper {
	return al.Shaper
}

// IsVisible returns true if the  light shape should be visible.
func (al *AreaSpotLight) IsVisible() bool {
	return al.visible
}

// Direction returns the direction of the spotlight
func (al *AreaSpotLight) Direction() Vector {
	// return to.SubPoint(asl.Position()).Normalize()
	return al.direction
}

// Angle returns the angle of the spot light
func (al *AreaSpotLight) Angle() float64 {
	return al.angle
}

// SpotLight implements the light interface and is a single point light with no size
type SpotLight struct {
	position  Point
	intensity Color
	angle     float64
	direction Vector
}

// NewSpotLight returns a new spot light
func NewSpotLight(from Point, i Color, angle float64, to Point) *SpotLight {
	sl := &SpotLight{
		position:  from,
		intensity: i,
		angle:     angle,
	}
	sl.direction = to.SubPoint(from).Normalize()
	return sl
}

// Intensity returns the intensity of the light
func (l *SpotLight) Intensity() Color {
	return l.intensity
}

// SetIntensity sets the intensity of the light
func (l *SpotLight) SetIntensity(c Color) {
	l.intensity = c
}

// Position returns the position of the light
func (l *SpotLight) Position() Point {
	return l.position
}

// RandomPosition returns the position of the light
func (l *SpotLight) RandomPosition(rng *rand.Rand) Point {
	return l.position
}

// Shape returns the Shaper object of this light, spot lights have no shape
func (l *SpotLight) Shape() Shaper {
	return nil
}

// IsVisible returns true if light is visible, spot lights are never visible.
func (l *SpotLight) IsVisible() bool {
	return false
}

// Direction returns the direction of the spotlight
func (l *SpotLight) Direction() Vector {
	return l.direction
}

// Angle returns the angle of the spot light
func (l *SpotLight) Angle() float64 {
	return l.angle
}

// lighting returns the color for a given point
func lighting(m *Material, o Shaper, p Point, l Light, eye, normal Vector, intensity float64, rays int, u, v float64, rng *rand.Rand) Color {
	var ambient, diffuse, specular Color
	clr := m.Color

	if m.HasPattern() {
		clr = clr.Blend(m.Pattern.ColorAtObject(o, p))
	}

	if m.HasTexture() {
		// Texture blends with the base color, so pass it in here
		// - Kd - material diffuse is multiplied by the texture value
		// This is only used by Smooth Triangles, to apply textures to other shapes,
		// use the ImageTexturePattern
		// Awkward... consider fixing
		switch o.(type) {
		case *SmoothTriangle:
			clr = clr.Blend(m.ColorAtTexture(o, u, v))
		default:
			log.Fatal("Texture attached to non Smooth Triangle, use an ImagePattern instead.")
		}
	}

	// combine surface color with light's color/intensity
	effectiveColor := clr.Blend(l.Intensity())

	// compute ambient contribution
	ambient = effectiveColor.Scale(m.Ambient)

	// Compute the emissive contribution
	emissive := m.Emissive

	// light not visible, ignore diffuse and specular components
	if intensity == 0 {
		return ambient.Add(emissive)
	}

	sum := Black()

	switch l.(type) {
	case *PointLight:
		rays = 1 // randomposition on pointlight always returns the same
	case *SpotLight:
		rays = 1 // RandomPosition on SpotLight always returns the same
	}

	for try := 0; try < rays; try++ {
		// find the direction to the light source
		lightv := l.RandomPosition(rng).SubPoint(p).Normalize()

		// lightDotNormal represents the cosine of the angle between the light vector and the normal vector
		// a negative number means the light is on the other side of the surface
		// a very small number here means the angle is very close to 90 degree, this number is used
		// to scale the diffuse contribution
		lightDotNormal := lightv.Dot(normal)

		visible := true

		if lightDotNormal < 0 {
			diffuse = ColorName(colornames.Black)
			specular = ColorName(colornames.Black)
			visible = false
		} else {
			// handle spotlights
			switch l.(type) {
			case *SpotLight:
				sl := l.(*SpotLight)
				// calculate the angle between the light direction vector and lightv
				dp := sl.Direction().Dot(lightv) // cosine of the angle
				angle := math.Pi - math.Acos(dp)
				// compare to the angle of the spotlight
				if angle > sl.Angle()/2 {
					diffuse = ColorName(colornames.Black)
					specular = ColorName(colornames.Black)
					visible = false
				}
				// TODO: Fix this to remove repetition
			case *AreaSpotLight:
				sl := l.(*AreaSpotLight)
				// calculate the angle between the light direction vector and lightv
				dp := sl.Direction().Dot(lightv) // cosine of the angle
				angle := math.Pi - math.Acos(dp)
				// compare to the angle of the spotlight
				if angle > sl.Angle()/2 {
					diffuse = ColorName(colornames.Black)
					specular = ColorName(colornames.Black)
					visible = false
				}
			}
		}

		if visible {
			// compute the diffuse contribution
			diffuse = effectiveColor.Scale(m.Diffuse).Scale(lightDotNormal)

			// reflectDotEye represens the cosine of the angle between the relfection vector and the eye vector
			// a negative number means the light reflects away from the eye
			reflectv := lightv.Negate().Reflect(normal)
			reflectDotEye := reflectv.Dot(eye)

			if reflectDotEye <= 0 {
				specular = ColorName(colornames.Black)
			} else {
				// compute the specular contrbution
				factor := math.Pow(reflectDotEye, m.Shininess)
				specular = l.Intensity().Scale(m.Specular).Scale(factor)
			}
		}
		sum = sum.Add(diffuse).Add(specular)
	}
	return emissive.Add(ambient).Add(sum.Scale(1.0 / float64(rays)).Scale(intensity))
}

// ColorAtPoint returns the clamped color at the given point
func ColorAtPoint(m *Material, o Shaper, p Point, l Light, eye, normal Vector, inShadow float64, rng *rand.Rand) Color {
	return lighting(m, o, p, l, eye, normal, inShadow, 1, 0, 0, rng).Clamp()
}
