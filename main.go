package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"runtime/pprof"
	"sync"

	_ "net/http/pprof"

	"github.com/lucasb-eyer/go-colorful"

	"golang.org/x/image/colornames"

	"github.com/DanTulovsky/tracer/tracer"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

type projectile struct {
	Position tracer.Point
	Velocity tracer.Vector
}

type environment struct {
	Gravity tracer.Vector
	Wind    tracer.Vector
}

func tick(e environment, p projectile) projectile {
	position := p.Position.AddVector(p.Velocity)
	velocity := p.Velocity.AddVector(e.Gravity).AddVector(e.Wind)

	return projectile{Position: position, Velocity: velocity}
}

func addToCanvas(c *tracer.Canvas, p projectile) error {
	pos := p.Position

	x := int(pos.X())
	y := c.Height - int(pos.Y())

	c.Set(x, y, tracer.ColorName(colornames.Red))

	return nil
}

func testCanvas() {
	// Canvas
	c := tracer.NewCanvas(900, 550)

	ticks := 0
	vScale := 11.25
	startiPosition := tracer.NewPoint(0, 1, 0)
	initialVelocity := tracer.NewVector(1, 1.8, 0).Normalize().Scale(vScale)
	gravity := tracer.NewVector(0, -0.1, 0)
	wind := tracer.NewVector(-0.01, 0, 0)

	p := projectile{Position: startiPosition, Velocity: initialVelocity}
	e := environment{Gravity: gravity, Wind: wind}

	fmt.Printf("position: %2f\n", p.Position)
	for p.Position.Y() > 0 {
		p = tick(e, p)
		fmt.Printf("position: %2f\n", p.Position)
		ticks++
		addToCanvas(c, p)
	}
	fmt.Printf("Total Ticks: %v\n", ticks)

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	c.ExportToPNG(f)
}

func test1() {

	p := tracer.NewPoint(1, 1.000000001, 1)
	p2 := tracer.NewPoint(2, 4, 6)
	p3 := tracer.NewPoint(2, 4, 6.000000001)

	log.Printf("%#v", p.Equals(p2))
	log.Printf("%#v", p2.Equals(p3))
}

func clock() {

	c := tracer.NewCanvas(550, 600)

	// center
	center := tracer.NewVector(275, 0, 300)
	c.SetFloat(center.X(), center.Z(), tracer.ColorName(colornames.Yellow))

	radius := 7.0 / 8.0 * center.X()
	twelve := tracer.NewPoint(0, 0, 1)

	for hour := 1.0; hour <= 12; hour++ {
		m := tracer.IdentityMatrix().RotateY(hour*(math.Pi/6.0)).Scale(radius, 1, radius).Translate(center.X(), center.Y(), center.Z())
		p := twelve.TimesMatrix(m)
		c.SetFloat(p.X(), p.Z(), tracer.ColorName(colornames.Red))
	}

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	c.ExportToPNG(f)
}

func circle() {
	// first circle drawn by a ray
	canvasX := 200
	canvasY := canvasX

	c := tracer.NewCanvas(canvasX, canvasY)

	// camera location
	camera := tracer.NewPoint(0, 0, -5)

	type wall struct {
		Z    float64
		Size float64
	}

	// wall is parallel to the y-axis, on negative z
	// size is large enough to sho a unit spere at the origin from the camera
	w := wall{Z: 10, Size: 7}

	// size of a world pixel
	pixelSize := w.Size / float64(canvasX)

	// transform matrix
	m := tracer.IdentityMatrix().Scale(1, 0.5, 1).RotateZ(math.Pi/4).Shear(1, 0, 0, 0, 0, 0)

	shape := tracer.NewUnitSphere()
	shape.SetTransform(m)

	clr := tracer.ColorName(colorful.HappyColor())

	// for each row of pixels on the canvas
	for y := 0.0; y < float64(canvasY); y++ {
		// world coordinate of y
		wy := w.Size/2 - pixelSize*y

		for x := 0.0; x < float64(canvasX); x++ {

			// world y coordinate of x
			wx := -w.Size/2 + pixelSize*x

			// point on the wall the ray is targetting
			target := tracer.NewPoint(wx, wy, w.Z)

			// the ray from camera to the world target
			ray := tracer.NewRay(camera, target.SubPoint(camera).Normalize())

			if _, err := shape.IntersectWith(ray).Hit(); err == nil {
				c.SetFloat(x, y, clr)
			}
		}
	}

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	c.ExportToPNG(f)
}

func sphere() {

	// first sphere drawn by a ray
	canvasX := 200
	canvasY := canvasX

	c := tracer.NewCanvas(canvasX, canvasY)

	// camera location
	camera := tracer.NewPoint(0, 0, -5)

	type wall struct {
		Z    float64
		Size float64
	}

	// wall is parallel to the y-axis, on negative z
	// size is large enough to sho a unit spere at the origin from the camera
	w := wall{Z: 10, Size: 7}

	// size of a world pixel
	pixelSize := w.Size / float64(canvasX)

	// transform matrix
	// m := tracer.IdentityMatrix().Scale(1, 0.5, 1).RotateZ(math.Pi/4).Shear(1, 0, 0, 0, 0, 0)
	m := tracer.IdentityMatrix()

	// material
	mat := tracer.NewDefaultMaterial()
	mat.Color = tracer.ColorName(colornames.Yellow)
	mat.Ambient = 0.1
	mat.Diffuse = 0.9
	mat.Specular = 0.9
	mat.Shininess = 30.0

	shape := tracer.NewUnitSphere()
	shape.SetTransform(m)
	shape.SetMaterial(mat)

	// light source
	light := tracer.NewPointLight(tracer.NewPoint(-10, 10, -10), tracer.ColorName(colornames.White))

	var wg sync.WaitGroup

	// for each row of pixels on the canvas
	for y := 0.0; y < float64(canvasY); y++ {
		for x := 0.0; x < float64(canvasX); x++ {

			wg.Add(1)

			go func(x, y float64) {

				// world coordinate of y
				wy := w.Size/2 - pixelSize*y
				// world y coordinate of x
				wx := -w.Size/2 + pixelSize*x

				// point on the wall the ray is targetting
				target := tracer.NewPoint(wx, wy, w.Z)

				// the ray from camera to the world target
				ray := tracer.NewRay(camera, target.SubPoint(camera).Normalize())

				if hit, err := shape.IntersectWith(ray).Hit(); err == nil {

					comp := tracer.PrepareComputations(hit, ray)
					clr := tracer.ColorAtPoint(comp.Object.Material(), comp.Object, comp.Point, light, comp.EyeV, comp.NormalV, false)

					c.SetFloat(x, y, clr)
				}
				wg.Done()

			}(x, y)
		}
	}

	wg.Wait()

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	c.ExportToPNG(f)
}

func scene() {

	width, height := 300.0, 300.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)

	// second light
	l2 := tracer.NewPointLight(tracer.NewPoint(10, 10, -10), tracer.ColorName(colornames.Pink))
	w.AddLight(l2)
	// w.SetLights([]tracer.Light{l2})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 1.5, -5)
	to := tracer.NewPoint(0, 1, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	var material *tracer.Material

	// floor
	floor := tracer.NewPlane()
	material = floor.Material()
	material.Color = tracer.ColorName(colornames.Beige)
	material.Specular = 0
	w.AddObject(floor)

	// left wall
	leftWall := tracer.NewPlane()
	leftWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(-15, 0, 0))
	leftWall.SetMaterial(floor.Material())
	w.AddObject(leftWall)

	// right wall
	rightWall := tracer.NewPlane()
	rightWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(15, 0, 0))
	rightWall.SetMaterial(floor.Material())
	w.AddObject(rightWall)

	// sphere
	middle := tracer.NewUnitSphere()
	middle.SetTransform(tracer.IdentityMatrix().Translate(-0.5, 1, 0.5))
	material = middle.Material()
	material.Color = tracer.ColorName(colornames.Greenyellow)
	material.Diffuse = 0.7
	material.Specular = 0.3
	p1 := tracer.NewStripedPattern(tracer.ColorName(colornames.Red), tracer.Black())
	p1.SetTransform(tracer.IdentityMatrix().Scale(0.3, 0.1, 0.3).RotateX(math.Pi / 1.5).RotateY(math.Pi / 5))
	material.SetPattern(p1)
	w.AddObject(middle)

	// another sphere
	right := tracer.NewUnitSphere()
	right.SetTransform(tracer.IdentityMatrix().Scale(0.5, 0.5, 0.5).Translate(1.5, 0.5, -0.5))
	material = right.Material()
	material.Color = tracer.ColorName(colornames.Lime)
	material.Diffuse = 0.7
	material.Specular = 0.3
	w.AddObject(right)

	// cube
	left := tracer.NewUnitCube()
	left.SetTransform(
		tracer.IdentityMatrix().Scale(0.33, 0.33, 0.33).RotateX(math.Pi/4).RotateY(math.Pi/4).RotateZ(math.Pi/4).Translate(-1.5, 2, -0.5))
	material = left.Material()
	material.Color = tracer.ColorName(colornames.Lightblue)
	material.Diffuse = 0.7
	material.Specular = 0.3
	w.AddObject(left)

	canvas := w.Render()

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)

}

func main() {

	flag.Parse()
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

	// testCanvas()
	// test1()
	// clock()
	// circle()
	// sphere()
	scene()

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
