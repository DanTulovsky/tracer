package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"golang.org/x/image/colornames"

	_ "net/http/pprof"

	"github.com/DanTulovsky/tracer/tracer"
	"github.com/DanTulovsky/tracer/utils"
	"github.com/lucasb-eyer/go-colorful"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

// return a sphere at this point, scaled by size
func sphere(p tracer.Point, size float64) *tracer.Sphere {
	transform := tracer.IdentityMatrix().Scale(size, size, size).Translate(p.X(), p.Y(), p.Z())

	s := tracer.NewUnitSphere()
	s.SetTransform(transform)
	s.Material().Color = tracer.ColorName(colorful.FastHappyColor())

	return s
}

// return a glass sphere at this point, scaled by size
func glassSphere(p tracer.Point, size float64) *tracer.Sphere {
	randomT := float64(utils.Random(50, 70)) / 100 // [0.5 - 0.7)]
	randomR := float64(utils.Random(60, 80)) / 100

	transform := tracer.IdentityMatrix().Scale(size, size, size).Translate(p.X(), p.Y(), p.Z())

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
	width, height := 400.0, 200.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		// tracer.NewPointLight(tracer.NewPoint(2, 10, -1), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(0, 20, -1), tracer.NewColor(1, 1, 1)),
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

	// plane
	plane := tracer.NewPlane()
	pp := tracer.NewCheckerPattern(tracer.ColorName(colornames.White), tracer.ColorName(colornames.Red))
	plane.Material().SetPattern(pp)
	w.AddObject(plane)

	// groups
	g1 := tracer.NewGroup()

	for i := 0; i < 30; i++ {
		x := utils.RandomFloat(-4, 4)
		z := utils.RandomFloat(1, 15)

		size := utils.RandomFloat(0.1, 0.7)
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

	render(w)
}

func render(w *tracer.World) {
	canvas := w.Render()

	// Export
	f, err := os.Create(fmt.Sprintf("%s/Downloads/image.png", utils.Homedir()))
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}
func main() {

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	scene()

}
