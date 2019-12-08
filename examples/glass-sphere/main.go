package main

import (
	"flag"
	"log"
	"math"
	"runtime"

	"golang.org/x/image/colornames"

	"github.com/DanTulovsky/tracer/tracer"
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func env() *tracer.World {
	// width, height := 150.0, 100.0
	width, height := 400.0, 300.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(1, 4, -1), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 1.7, -4.7)
	to := tracer.NewPoint(0, -1, 10)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	return w
}

func floor() *tracer.Plane {
	p := tracer.NewPlane()
	pp := tracer.NewCheckerPattern(tracer.ColorName(colornames.Red), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func ceiling() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IM().Translate(0, 5, 0))
	pp := tracer.NewCheckerPattern(tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func backWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IM().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, 10))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Teal), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func frontWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IM().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, -5))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Purple), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func rightWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IM().RotateZ(math.Pi/2).Translate(4, 0, 0))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Peachpuff), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func leftWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IM().RotateZ(math.Pi/2).Translate(-4, 0, 0))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Lightgreen), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func sphere() *tracer.Sphere {
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IM().Scale(.75, .75, .75).Translate(0, 1.75, 0))
	s.Material().Ambient = 0
	s.Material().Diffuse = 0
	s.Material().Reflective = 1.0
	s.Material().Transparency = 1.0
	s.Material().ShadowCaster = false
	s.Material().RefractiveIndex = 1.573

	return s
}

func pedestal() *tracer.Cube {
	s := tracer.NewUnitCube()
	s.SetTransform(tracer.IM().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
	s.Material().Color = tracer.ColorName(colornames.Gold)

	return s
}

func cone() *tracer.Cone {
	s := tracer.NewClosedCone(-2, 0)
	s.SetTransform(tracer.IM().Translate(0, 2, 0))
	sp := tracer.NewCheckerPattern(tracer.ColorName(colornames.Green), tracer.ColorName(colornames.Violet))
	s.Material().SetPattern(sp)
	return s
}

func background() *tracer.Group {
	g := tracer.NewGroup()
	g.AddMember(cone())
	g.SetTransform(tracer.IM().Translate(0, 1, 6))
	return g
}

func group(s ...tracer.Shaper) *tracer.Group {
	g := tracer.NewGroup()

	for _, s := range s {
		g.AddMember(s)
	}

	return g
}

func scene() {
	w := env()

	w.AddObject(backWall())
	w.AddObject(frontWall())
	w.AddObject(rightWall())
	w.AddObject(leftWall())
	w.AddObject(ceiling())
	w.AddObject(floor())

	w.AddObject(group(sphere(), pedestal()))
	w.AddObject(background())

	tracer.Render(w)
}

func main() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	scene()
}
