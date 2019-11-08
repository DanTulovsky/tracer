package tracer

import (
	"math"
	"testing"

	"golang.org/x/image/colornames"
)

func BenchmarkRenderSphere(b *testing.B) {
	width, height := 300.0, 300.0
	w := NewDefaultWorld(width, height)

	// where the camera is and where it's pointing; also which way is "up"
	from := NewPoint(0, 1.7, -4.7)
	to := NewPoint(0, -1, 10)
	up := NewVector(0, 1, 0)
	cameraTransform := ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	s1 := NewUnitSphere()

	w.AddObject(s1)

	for n := 0; n < b.N; n++ {
		RenderToFile(w, "/tmp/output.png")
	}
}

func floor() *Plane {
	p := NewPlane()
	pp := NewCheckerPattern(ColorName(colornames.Red), ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func ceiling() *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().Translate(0, 5, 0))
	pp := NewCheckerPattern(ColorName(colornames.Blue), ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func backWall() *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, 10))
	pp := NewStripedPattern(ColorName(colornames.Teal), ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func frontWall() *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, -5))
	pp := NewStripedPattern(ColorName(colornames.Purple), ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func rightWall() *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().RotateZ(math.Pi/2).Translate(4, 0, 0))
	pp := NewStripedPattern(ColorName(colornames.Peachpuff), ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func leftWall() *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().RotateZ(math.Pi/2).Translate(-4, 0, 0))
	pp := NewStripedPattern(ColorName(colornames.Lightgreen), ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func sphere() *Sphere {
	s := NewUnitSphere()
	s.SetTransform(IdentityMatrix().Scale(.75, .75, .75).Translate(0, 1.75, 0))
	s.Material().Ambient = 0
	s.Material().Diffuse = 0
	s.Material().Reflective = 1.0
	s.Material().Transparency = 1.0
	s.Material().ShadowCaster = false
	s.Material().RefractiveIndex = 1.573

	return s
}

func pedestal() *Cube {
	s := NewUnitCube()
	s.SetTransform(IdentityMatrix().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
	s.Material().Color = ColorName(colornames.Gold)

	return s
}

func cone() *Cone {
	s := NewClosedCone(-2, 0)
	s.SetTransform(IdentityMatrix().Translate(0, 2, 0))
	sp := NewCheckerPattern(ColorName(colornames.Green), ColorName(colornames.Violet))
	s.Material().SetPattern(sp)
	return s
}

func background() *Group {
	g := NewGroup()
	g.AddMember(cone())
	g.SetTransform(IdentityMatrix().Translate(0, 1, 6))
	return g
}

func group(s ...Shaper) *Group {
	g := NewGroup()

	for _, s := range s {
		g.AddMember(s)
	}

	return g
}

func env() *World {
	// width, height := 150.0, 100.0
	width, height := 400.0, 300.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]Light{
		NewPointLight(NewPoint(1, 4, -1), NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := NewPoint(0, 1.7, -4.7)
	to := NewPoint(0, -1, 10)
	up := NewVector(0, 1, 0)
	cameraTransform := ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	return w
}

func scene() *World {
	w := env()

	w.AddObject(backWall())
	w.AddObject(frontWall())
	w.AddObject(rightWall())
	w.AddObject(leftWall())
	w.AddObject(ceiling())
	w.AddObject(floor())

	w.AddObject(group(sphere(), pedestal()))
	w.AddObject(background())

	return w
}

func BenchmarkRenderGlassSphere(b *testing.B) {

	w := scene()
	for n := 0; n < b.N; n++ {
		RenderToFile(w, "/tmp/output.png")
	}
}
