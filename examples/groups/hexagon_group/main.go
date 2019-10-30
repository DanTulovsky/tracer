package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"

	_ "net/http/pprof"

	"github.com/DanTulovsky/tracer/tracer"
	"github.com/DanTulovsky/tracer/utils"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

func corner() *tracer.Sphere {
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Scale(0.25, 0.25, 0.25).Translate(0, 0, -1))
	return s
}

func edge() *tracer.Cylinder {
	c := tracer.NewCylinder(0, 1)
	c.SetTransform(
		tracer.IdentityMatrix().Scale(0.25, 1, 0.25).RotateZ(-math.Pi/2).RotateY(-math.Pi/6).Translate(0, 0, -1))
	return c
}

func side() *tracer.Group {
	g := tracer.NewGroup()

	g.AddMember(corner())
	g.AddMember(edge())

	return g
}

func hexagon() {

	// width, height := 100.0, 100.0
	width, height := 200.0, 200.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		// tracer.NewPointLight(tracer.NewPoint(2, 10, -1), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(-0.5, 1, 0), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 2.0, -4)
	to := tracer.NewPoint(0, 0, 0)
	up := tracer.NewVector(0, 1, 0)
	fov := math.Pi / 3.8

	camera := tracer.NewCamera(width, height, fov)
	cameraTransform := tracer.ViewTransform(from, to, up)
	camera.SetTransform(cameraTransform)

	w.SetCamera(camera)

	// hexagon
	hex := tracer.NewGroup()

	for n := 0.0; n < 6; n++ {
		s := side()
		s.SetTransform(tracer.IdentityMatrix().RotateY(n * math.Pi / 3))
		hex.AddMember(s)
	}

	hex.SetTransform(tracer.IdentityMatrix().RotateX(-math.Pi / 12))

	w.AddObject(hex)

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

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	hexagon()

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

}
