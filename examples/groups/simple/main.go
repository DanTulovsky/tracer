package main

import (
	"flag"
	"log"
	"runtime"

	"github.com/DanTulovsky/tracer/tracer"
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
	from := tracer.NewPoint(0, 2, -7)
	to := tracer.NewPoint(0, -1, 10)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	return w
}

func sphere() *tracer.Sphere {
	s := tracer.NewUnitSphere()
	return s
}

func cube() *tracer.Cube {
	s := tracer.NewUnitCube()
	return s
}

func cone() *tracer.Cone {
	s := tracer.NewDefaultCone()
	return s
}

func cylinder() *tracer.Cylinder {
	s := tracer.NewClosedCylinder(0, 5)
	return s
}

func tri() *tracer.Triangle {
	t := tracer.NewTriangle(tracer.NewPoint(-1, 0, 0), tracer.NewPoint(0, 2, 0), tracer.NewPoint(1, 0, 0))
	return t
}

func scene() {
	w := env()

	g := tracer.NewGroup()
	g.AddMember(cube())
	w.AddObject(g)

	tracer.Render(w)
}

func main() {

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	scene()
}
