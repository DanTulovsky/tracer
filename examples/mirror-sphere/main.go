package main

import (
	"flag"
	"log"
	"math"
	"runtime"

	"golang.org/x/image/colornames"

	"github.com/DanTulovsky/tracer/tracer"
)

var (
	output = flag.String("output", "", "name of the output file, if empty, renders to screen")
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func env() *tracer.World {
	// width, height := 100.0, 100.0
	width, height := 400.0, 300.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(0, 4, -1), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 3, -4)
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
	p.SetTransform(tracer.IdentityMatrix().Translate(0, 5, 0))
	pp := tracer.NewCheckerPattern(tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func backWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, 10))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Lightgreen), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func rightWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(4, 0, 0))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Lightgreen), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func leftWall() *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(-4, 0, 0))
	pp := tracer.NewStripedPattern(tracer.ColorName(colornames.Lightgreen), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func sphere() *tracer.Sphere {
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Translate(0, 1.5, 0))
	s.Material().Ambient = 0
	s.Material().Diffuse = 0
	s.Material().Reflective = 1
	return s
}

func scene() {
	w := env()

	w.AddObject(sphere())
	w.AddObject(floor())
	w.AddObject(ceiling())
	w.AddObject(backWall())
	w.AddObject(rightWall())
	w.AddObject(leftWall())

	render(w)
}

func render(w *tracer.World) {
	if *output != "" {
		tracer.Render(w, *output)
	} else {
		tracer.RenderLive(w)
	}
}

func main() {

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	scene()
}
