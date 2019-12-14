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
	s.SetTransform(IM().Scale(.75, .75, .75).Translate(0, 1.75, 0))
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
	s.SetTransform(IM().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
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
	p.SetTransform(IM().Translate(0, y, 0))
	pp := NewCheckerPattern(ColorName(colornames.Gray), ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

// room returns a room with all walls of the provided sizes
func room(left, front, right, back, clng, flr float64) *Group {
	g := NewGroup()
	g.AddMember(floor(flr))
	g.AddMember(backWallGhost(back))
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
	p.SetTransform(IM().Translate(0, y, 0))
	pp := NewGradientPattern(ColorName(colornames.Blue), ColorName(colornames.Red))
	pp.SetTransform(IM().Scale(10, 1, 1).Translate(-15, 0, 0))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0
	p.Material().Ambient = 0.15

	return p
}
func backWall(z float64) *Plane {
	p := NewPlane()
	p.SetTransform(
		IM().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	p.Material().Color = ColorName(colornames.Lightpink)
	p.Material().Specular = 0

	return p
}
func backWallGhost(z float64) *Plane {
	p := NewPlane()
	p.SetTransform(
		IM().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	ppuv, err := NewUVImagePattern("/Users/dant/go/src/github.com/DanTulovsky/tracer/images/ghost.png")
	if err != nil {
		log.Fatal(err)
	}
	pp := NewTextureMapPattern(ppuv, NewPlaneMap())
	pp.SetTransform(IM().Scale(10, 5, 5).RotateY(math.Pi/2).Translate(0, 0, -3))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func frontWall(z float64) *Plane {
	p := NewPlane()
	p.SetTransform(
		IM().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	uvpp := NewUVCheckersPattern(4, 4,
		ColorName(colornames.Orange), ColorName(colornames.White))
	pp := NewTextureMapPattern(uvpp, NewPlaneMap())
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func rightWall(x float64) *Plane {
	p := NewPlane()
	p.SetTransform(IM().RotateZ(math.Pi/2).Translate(x, 0, 0))
	pp := NewGradientPattern(
		ColorName(colornames.Orange), ColorName(colornames.White))
	pp.SetTransform(IM().Scale(10, 1, 1).Translate(-5, 0, 0))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func leftWall(x float64) *Plane {
	p := NewPlane()
	p.SetTransform(IM().RotateZ(math.Pi/2).Translate(x, 0, 0))
	pp := NewStripedPattern(
		ColorName(colornames.Lightskyblue), ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func sphere() *Sphere {
	s := NewUnitSphere()
	s.SetTransform(IM().Scale(.75, .75, .75).Translate(0, 1.75, 0))
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
	s.SetTransform(IM().Translate(0, 2, 0))
	sp := NewCheckerPattern(ColorName(colornames.Green), ColorName(colornames.Violet))
	s.Material().SetPattern(sp)
	return s
}

func background() *Group {
	g := NewGroup()
	g.AddMember(cone())
	g.SetTransform(IM().Translate(0, 1, 6))
	return g
}

func group(s ...Shaper) *Group {
	g := NewGroup()

	for _, s := range s {
		g.AddMember(s)
	}

	return g
}

func testenv() *World {
	width, height := 640.0, 480.0

	// setup world, default light and camera
	w := NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5
	w.Config.SoftShadows = false
	w.Config.Antialias = 4

	// override light here
	w.SetLights([]Light{
		NewPointLight(NewPoint(1, 3, -1), NewColor(1, 1, 1)),
		// NewPointLight(NewPoint(-9, 10, 10), NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := NewPoint(0, 1.7, -4.7)
	to := NewPoint(0, 0, 10)
	up := NewVector(0, 1, 0)
	cameraTransform := ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	return w
}

func scene() *World {
	w := testenv()

	w.AddObject(backWall(10))
	w.AddObject(frontWall(-5))
	w.AddObject(rightWall(4))
	w.AddObject(leftWall(-4))
	w.AddObject(ceiling(4))
	w.AddObject(floor(0))

	w.AddObject(group(sphere(), pedestal()))
	w.AddObject(background())

	return w
}

func testenvxy(width, height float64) *World {
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
	from := NewPoint(0, 10, -30)
	to := NewPoint(0, 0, 40)
	up := NewVector(0, 1, 0)
	cameraTransform := ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 3.5)

	return w
}
func BenchmarkRenderMonkey(b *testing.B) {
	dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/bench/"))
	f := path.Join(dir, "obj", "monkey.obj")
	w := testenvxy(640, 480)
	w.Config.Antialias = 3
	w.Config.SoftShadows = false

	g, err := ParseOBJ(f)
	if err != nil {
		log.Fatalln(err)
	}

	g.SetTransform(IM().Scale(6.5, 6.5, 6.5).RotateY(math.Pi/7).Translate(0, 4, 0))

	w.AddObject(g)
	for n := 0; n < b.N; n++ {
		output := path.Join(dir, "output", "monkey.png")
		RenderToFile(w, output)
	}

}
func BenchmarkRenderSmoothMonkey(b *testing.B) {
	dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/bench/"))
	f := path.Join(dir, "obj", "monkey-smooth.obj")
	w := testenvxy(640, 480)
	w.Config.Antialias = 3
	w.Config.SoftShadows = false

	g, err := ParseOBJ(f)
	if err != nil {
		log.Fatalln(err)
	}

	g.SetTransform(IM().Scale(6.5, 6.5, 6.5).RotateY(math.Pi/7).Translate(0, 4, 0))

	w.AddObject(g)
	for n := 0; n < b.N; n++ {
		output := path.Join(dir, "output", "monkey-smooth.png")
		RenderToFile(w, output)
	}

}
func BenchmarkRenderEmmisive(b *testing.B) {
	dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/bench/"))
	w := envxy(640, 480)
	w.Config.Antialias = 3
	w.Config.SoftShadows = true
	w.Config.SoftShadowRays = 10
	w.Config.AreaLightRays = 10

	l := NewAreaLight(NewUnitSphere(),
		ColorName(colornames.White), true)
	l.SetTransform(
		IM().Scale(0.2, 1, 0.2).Translate(2, 1, 2))
	l.SetIntensity(l.Intensity().Scale(0.5))

	l2 := NewAreaLight(NewUnitCube(),
		ColorName(colornames.White), true)
	l2.SetTransform(
		IM().Scale(0.2, 1, 0.2).Translate(-2, 1, 2))
	l2.SetIntensity(l.Intensity().Scale(0.5))

	w.SetLights(Lights{l, l2})

	g := mirrorSphereOnPedestal()
	g.SetTransform(IM().Translate(0, 0, 2.5))

	w.AddObject(g)
	w.AddObject(defaultroom())

	for n := 0; n < b.N; n++ {
		output := path.Join(dir, "output", "emmisive.png")
		RenderToFile(w, output)
	}
}

func BenchmarkRenderSphere(b *testing.B) {
	dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/bench/"))
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
		output := path.Join(dir, "output", "sphere.png")
		RenderToFile(w, output)
	}
}

func BenchmarkRenderGlassSphere(b *testing.B) {
	dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/bench/"))
	w := scene()
	for n := 0; n < b.N; n++ {
		output := path.Join(dir, "output", "glasssphere.png")
		RenderToFile(w, output)
	}
}
