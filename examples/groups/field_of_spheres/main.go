package main

import (
	"flag"
	"log"
	"math"
	"runtime"

	"golang.org/x/image/colornames"

	_ "net/http/pprof"

	"github.com/DanTulovsky/tracer/tracer"
	"github.com/DanTulovsky/tracer/utils"
	"github.com/lucasb-eyer/go-colorful"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
)

func init() {
	// This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

// return a sphere at this point, scaled by size
func sphere(p tracer.Point, size float64) *tracer.Sphere {
	transform := tracer.IM().Scale(size, size, size).Translate(p.X(), p.Y(), p.Z())

	s := tracer.NewUnitSphere()
	s.SetTransform(transform)
	s.Material().Color = tracer.ColorName(colorful.FastHappyColor())

	return s
}

// return a glass sphere at this point, scaled by size
func glassSphere(p tracer.Point, size float64) *tracer.Sphere {
	randomT := float64(utils.Random(80, 90)) / 100 // [0.5 - 0.7)]
	randomR := float64(utils.Random(60, 80)) / 100

	transform := tracer.IM().Scale(size, size, size).Translate(p.X(), p.Y(), p.Z())

	s := tracer.NewUnitSphere()
	s.SetTransform(transform)
	// s.Material().Color = tracer.ColorName(colorful.FastHappyColor())
	s.Material().Ambient = 0.1
	s.Material().Diffuse = 0.1
	s.Material().Transparency = randomT
	s.Material().RefractiveIndex = 1.34
	s.Material().Reflective = randomR

	return s
}
func scene() {

	// width, height := 100.0, 100.0
	// width, height := 400.0, 200.0
	width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		// tracer.NewPointLight(tracer.NewPoint(2, 10, -1), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(0, 4, 20), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 2.0, -0.2)
	to := tracer.NewPoint(0, 0, 8)
	up := tracer.NewVector(0, 1, 0)
	fov := math.Pi / 2

	camera := tracer.NewCamera(width, height, fov)
	cameraTransform := tracer.ViewTransform(from, to, up)
	camera.SetTransform(cameraTransform)

	w.SetCamera(camera)

	// floor
	plane := tracer.NewPlane()
	// pp := tracer.NewCheckerPattern(tracer.ColorName(colornames.White), tracer.ColorName(colornames.Red))
	// plane.Material().SetPattern(pp)
	plane.Material().Specular = 0
	plane.Material().Reflective = 1
	plane.Material().Diffuse = 0
	w.AddObject(plane)

	// back wall
	backWall := tracer.NewPlane()
	pbw := tracer.NewCheckerPattern(tracer.ColorName(colornames.Lightgray), tracer.ColorName(colornames.Gray))
	backWall.Material().SetPattern(pbw)
	backWall.SetTransform(tracer.IM().RotateX(math.Pi/2).Translate(0, 0, 30))
	backWall.Material().Specular = 0
	backWall.Material().Diffuse = 1
	w.AddObject(backWall)

	// ceiling wall
	ceiling := tracer.NewPlane()
	pc := tracer.NewCheckerPattern(tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.Yellow))
	pc.SetTransform(tracer.IM().Scale(3, 3, 3))
	ceiling.Material().SetPattern(pc)
	ceiling.SetTransform(tracer.IM().Translate(0, 5, 0))
	ceiling.Material().Specular = 0
	ceiling.Material().Diffuse = 1
	w.AddObject(ceiling)

	// groups
	g1 := tracer.NewGroup()

	numSpheres := 4

	// spheres
	for i := 0; i < numSpheres; i++ {
		x := utils.RandomFloat(-4, 4)
		z := utils.RandomFloat(1, 15)

		size := utils.RandomFloat(0.5, 0.9)
		center := tracer.NewPoint(x, size, z)
		var s *tracer.Sphere

		if utils.Random(0, 100) > 50 {
			s = sphere(center, size)
		} else {
			s = glassSphere(center, size)
		}
		g1.AddMember(s)
	}

	w.AddObject(g1)

	tracer.Render(w)
}

func main() {

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	scene()

}
