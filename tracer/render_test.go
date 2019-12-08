package tracer

import (
	"fmt"
	"log"
	"math"
	"path"
	"testing"

	"github.com/DanTulovsky/tracer/utils"
	"golang.org/x/image/colornames"
)

func envxy(width, height float64) *World {
	// setup world, default light and camera
	w := NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]Light{
		// NewPointLight(NewPoint(0, 4, 5), NewColor(1, 1, 1)),
		// NewPointLight(NewPoint(2, -10, -10), NewColor(1, 1, 1)),
		NewPointLight(NewPoint(-6, 10, -10), NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := NewPoint(0, 4, -9)
	to := NewPoint(0, 0, 20)
	up := NewVector(0, 1, 0)
	cameraTransform := ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 4)

	return w
}

func mirorsphere() *Sphere {
	s := NewUnitSphere()
	s.SetTransform(IdentityMatrix().Scale(.75, .75, .75).Translate(0, 1.75, 0))
	s.Material().Ambient = 0
	s.Material().Diffuse = 0
	s.Material().Reflective = 1.0
	s.Material().Transparency = 0
	s.Material().ShadowCaster = true
	// s.Material().RefractiveIndex = 1.573

	return s
}

func pedestal() *Cube {
	s := NewUnitCube()
	s.SetTransform(IdentityMatrix().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
	s.Material().Color = ColorName(colornames.Gold)
	up := NewUVCheckersPattern(8, 8,
		ColorName(colornames.White), ColorName(colornames.Violet))
	cp := NewTextureMapPattern(up, NewCubeMapSame(up))
	p := NewPerturbedPattern(cp, 0.09)
	s.Material().SetPattern(p)

	return s
}

func floor(y float64) *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().Translate(0, y, 0))
	pp := NewCheckerPattern(ColorName(colornames.Gray), ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

// room returns a room with all walls of the provided sizes
func room(left, front, right, back, clng, flr float64) *Group {
	g := NewGroup()
	g.AddMember(floor(flr))
	g.AddMember(backWall(back))
	g.AddMember(leftWall(left))
	g.AddMember(rightWall(right))
	g.AddMember(frontWall(front))
	g.AddMember(ceiling(clng))
	return g
}

func defaultroom() *Group {
	left, right := -5.0, 5.0
	front, back := -10.0, 10.0
	floor, ceiling := 0.0, 5.0
	return room(left, front, right, back, ceiling, floor)
}

func mirrorSphereOnPedestal() *Group {
	g := NewGroup()
	g.AddMembers(mirorsphere(), pedestal())
	return g
}

func ceiling(y float64) *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().Translate(0, y, 0))
	pp := NewGradientPattern(ColorName(colornames.Blue), ColorName(colornames.Red))
	pp.SetTransform(IdentityMatrix().Scale(10, 1, 1).Translate(-15, 0, 0))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0
	p.Material().Ambient = 0.15

	return p
}
func backWall(z float64) *Plane {
	p := NewPlane()
	p.SetTransform(
		IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	p.Material().Color = ColorName(colornames.Lightpink)
	p.Material().Specular = 0

	return p
}
func frontWall(z float64) *Plane {
	p := NewPlane()
	p.SetTransform(
		IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	uvpp := NewUVCheckersPattern(4, 4,
		ColorName(colornames.Orange), ColorName(colornames.White))
	pp := NewTextureMapPattern(uvpp, NewPlaneMap())
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func rightWall(x float64) *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().RotateZ(math.Pi/2).Translate(x, 0, 0))
	pp := NewGradientPattern(
		ColorName(colornames.Orange), ColorName(colornames.White))
	pp.SetTransform(IdentityMatrix().Scale(10, 1, 1).Translate(-5, 0, 0))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func leftWall(x float64) *Plane {
	p := NewPlane()
	p.SetTransform(IdentityMatrix().RotateZ(math.Pi/2).Translate(x, 0, 0))
	pp := NewStripedPattern(
		ColorName(colornames.Lightskyblue), ColorName(colornames.White))
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
		// NewPointLight(NewPoint(-9, 10, 10), NewColor(1, 1, 1)),
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

	w.AddObject(backWall(0))
	w.AddObject(frontWall(0))
	w.AddObject(rightWall(0))
	w.AddObject(leftWall(0))
	w.AddObject(ceiling(0))
	w.AddObject(floor(0))

	w.AddObject(group(sphere(), pedestal()))
	w.AddObject(background())

	return w
}

func BenchmarkRenderEmmisive(b *testing.B) {
	w := envxy(640, 480)
	w.Config.Antialias = 3
	w.Config.SoftShadows = true
	w.Config.SoftShadowRays = 10
	w.Config.AreaLightRays = 10

	l := NewAreaLight(NewUnitSphere(),
		ColorName(colornames.White), true)
	l.SetTransform(
		IdentityMatrix().Scale(0.2, 1, 0.2).Translate(2, 1, 2))
	l.SetIntensity(l.Intensity().Scale(0.5))

	l2 := NewAreaLight(NewUnitCube(),
		ColorName(colornames.White), true)
	l2.SetTransform(
		IdentityMatrix().Scale(0.2, 1, 0.2).Translate(-2, 1, 2))
	l2.SetIntensity(l.Intensity().Scale(0.5))

	w.SetLights(Lights{l, l2})

	g := mirrorSphereOnPedestal()
	g.SetTransform(IdentityMatrix().Translate(0, 0, 2.5))

	w.AddObject(g)
	w.AddObject(defaultroom())

	for n := 0; n < b.N; n++ {
		RenderToFile(w, "/tmp/output.png")
	}
}

func BenchmarkRenderSphere(b *testing.B) {
	width, height := 300.0, 300.0
	w := NewDefaultWorld(width, height)

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

func BenchmarkRenderGlassSphere(b *testing.B) {

	w := scene()
	for n := 0; n < b.N; n++ {
		RenderToFile(w, "/tmp/output.png")
	}
}

func benchmarkRenderObjParse(filename string, b *testing.B) {
	// func BenchmarkRenderObjParse(b *testing.B) {
	// for profiling
	// filename := "complex-smooth1.obj"

	// width, height := 640.0, 480.0
	width, height := 1400.0, 1000.0

	// setup world, default light and camera
	w := NewDefaultWorld(width, height)

	// override light here
	w.SetLights([]Light{
		NewPointLight(NewPoint(3, 4, -30), NewColor(1, 1, 1)),
		// NewPointLight(NewPoint(-5, 4, -1), NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := NewPoint(0, 6, -8)
	to := NewPoint(0, 0, 4)
	up := NewVector(0, 1, 0)
	cameraTransform := ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 3)

	dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/obj"))
	g, err := ParseOBJ(path.Join(dir, filename))
	if err != nil {
		log.Fatalln(err)
	}

	g.SetTransform(IdentityMatrix().Translate(0, 2, 0))

	w.AddObject(g)

	for n := 0; n < b.N; n++ {
		RenderToFile(w, "/tmp/output.png")
	}
}

func BenchmarkRenderObjParse0(b *testing.B) { benchmarkRenderObjParse("complex-smooth.obj", b) }
func BenchmarkRenderObjParse1(b *testing.B) { benchmarkRenderObjParse("complex-smooth1.obj", b) }
func BenchmarkRenderObjParse2(b *testing.B) { benchmarkRenderObjParse("complex-smooth2.obj", b) }
func BenchmarkRenderObjParse3(b *testing.B) { benchmarkRenderObjParse("monkey", b) }
