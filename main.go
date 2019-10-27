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
					xs := tracer.NewIntersections(hit)

					comp := tracer.PrepareComputations(hit, ray, xs)
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

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(-10, 10, -10), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(10, 10, -10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 1.5, -7)
	to := tracer.NewPoint(0, 1, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	var material *tracer.Material

	// floor
	floor := tracer.NewPlane()
	material = floor.Material()
	material.Color = tracer.NewColor(1, 1, 1)
	material.Specular = 0
	material.Reflective = 0.5
	// p := tracer.NewRingPattern(tracer.ColorName(colornames.Fuchsia), tracer.ColorName(colornames.Blue))
	// p := tracer.NewPertrubedPattern(
	// 	tracer.NewRingPattern(
	// 		tracer.ColorName(colornames.Fuchsia), tracer.ColorName(colornames.Blue)),
	// 	0.9)
	bp1 := tracer.NewStripedPattern(tracer.ColorName(colornames.Green), tracer.ColorName(colornames.White))
	bp2 := tracer.NewStripedPattern(tracer.ColorName(colornames.Green), tracer.ColorName(colornames.White))
	// rotate bp2 by 90 degrees
	bp2.SetTransform(tracer.IdentityMatrix().RotateY(math.Pi / 2))

	p := tracer.NewBlendedPattern(bp1, bp2)
	floor.Material().SetPattern(p)
	w.AddObject(floor)

	wallMaterial := tracer.NewDefaultMaterial()
	wallMaterial.Color = tracer.ColorName(colornames.Whitesmoke)

	// left wall
	leftWall := tracer.NewPlane()
	leftWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(-15, 0, 0))
	leftWall.SetMaterial(wallMaterial)
	w.AddObject(leftWall)

	// right wall
	rightWall := tracer.NewPlane()
	rightWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(15, 0, 0))
	pRightWall := tracer.NewPertrubedPattern(
		tracer.NewCheckerPattern(
			tracer.ColorName(colornames.Fuchsia), tracer.ColorName(colornames.Blue)),
		0.5)
	rightWall.Material().SetPattern(pRightWall)
	// rightWall.Material().Color = tracer.ColorName(colornames.Lightseagreen)
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
	right.SetTransform(tracer.IdentityMatrix().Scale(1, 1, 1).Translate(1, 2, -0.5))
	material = right.Material()
	material.Color = tracer.ColorName(colornames.Lime) // ignored when pattern
	material.Diffuse = 0.7
	material.Specular = 0.3
	p2 := tracer.NewStripedPattern(
		tracer.ColorName(colornames.Red), tracer.ColorName(colornames.Green))
	p2.SetTransform(tracer.IdentityMatrix().Scale(0.1, 0.1, 0.1).RotateX(math.Pi / 4))
	p3 := tracer.NewPertrubedPattern(p2, 0.6)
	material.SetPattern(p3)
	w.AddObject(right)

	// cube
	left := tracer.NewUnitCube()
	left.SetTransform(
		tracer.IdentityMatrix().Scale(0.33, 0.33, 0.33).RotateX(math.Pi/4).RotateY(math.Pi/4).RotateZ(math.Pi/4).Translate(-1.5, 2, -0.5))
	material = left.Material()
	material.Color = tracer.ColorName(colornames.Lightblue)
	material.Diffuse = 0.2
	material.Specular = 0.8
	// p4 := tracer.NewGradientPattern(tracer.ColorName(colornames.Black), tracer.ColorName(colornames.White))
	// p4.SetTransform(tracer.IdentityMatrix().Scale(2, 1, 1))
	// material.SetPattern(p4)
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

func colors() {

	c1 := tracer.NewColor(1, 0, 0)
	c2 := tracer.NewColor(10, 0, 0)

	log.Printf("NewColor(1, 0, 0): %v", c1)
	log.Printf("NewColor(10, 0, 0): %v", c2)

	c1mc, _ := colorful.MakeColor(c1)
	c2mc, _ := colorful.MakeColor(c2)

	log.Printf("colorful.MakeColor(NewColor(1, 0, 0)): %v", c1mc)
	log.Printf("colorful.MakeColor(NewColor(10, 0, 0)): %v", c2mc)
}

func mirrors() {

	width, height := 300.0, 300.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 3

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(0, 10, 0), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(3, 2, -10)
	to := tracer.NewPoint(-4.5, 1, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// floor
	floor := tracer.NewPlane()
	floor.Material().Color = tracer.ColorName(colornames.White)
	floor.Material().Specular = 0
	floor.Material().Reflective = 0.5
	w.AddObject(floor)

	leftWall := tracer.NewPlane()
	leftWall.Material().Color = tracer.ColorName(colornames.White)
	leftWall.Material().Specular = 0
	leftWall.Material().Reflective = 0
	leftWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(-15, 0, 0))
	leftWall.Material().Color = tracer.ColorName(colornames.Lightblue)
	w.AddObject(leftWall)

	rightWall := tracer.NewPlane()
	rightWall.Material().Color = tracer.ColorName(colornames.White)
	rightWall.Material().Specular = 0
	rightWall.Material().Reflective = 0
	rightWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(15, 0, 0))
	rightWall.Material().Color = tracer.ColorName(colornames.Lightcoral)
	w.AddObject(rightWall)

	// mirror1
	cube1 := tracer.NewUnitCube()
	cube1.SetTransform(
		tracer.IdentityMatrix().Scale(0.001, 1, 10).Translate(-2, 2, 0))
	cube1.Material().Reflective = 1
	// cube1.Material().Color = tracer.ColorName(colornames.Black)
	w.AddObject(cube1)

	// mirror2
	cube2 := tracer.NewUnitCube()
	cube2.SetTransform(
		tracer.IdentityMatrix().Scale(0.001, 1, 5).Translate(2, 2, 0))
	cube2.Material().Reflective = 1
	// cube2.Material().Color = tracer.ColorName(colornames.Black)
	w.AddObject(cube2)

	// sphere1
	sphere1 := tracer.NewUnitSphere()
	sphere1.SetTransform(
		tracer.IdentityMatrix().Scale(.5, .5, .5).Translate(0, 2, 2))
	sphere1.Material().Color = tracer.ColorName(colornames.Yellow)
	sphere1pattern := tracer.NewStripedPattern(tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.Purple))
	sphere1pattern.SetTransform(tracer.IdentityMatrix().Scale(0.2, 1, 1))
	sphere1.Material().SetPattern(sphere1pattern)
	w.AddObject(sphere1)

	canvas := w.Render()

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}

func mirror() {

	// width, height := 300.0, 300.0
	width, height := 400.0, 400.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 1

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(10, 8, -10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(3.5, 3.8, -5.7)
	to := tracer.NewPoint(-2, 0, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// floor
	floor := tracer.NewPlane()
	// floor.Material().Color = tracer.ColorName(colornames.Gray)
	floor.Material().Specular = 0
	floor.Material().Reflective = 0
	floorP := tracer.NewCheckerPattern(tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	w.AddObject(floor)

	leftWall := tracer.NewPlane()
	leftWall.Material().Color = tracer.ColorName(colornames.Lightskyblue)
	leftWall.Material().Specular = 0
	leftWall.Material().Reflective = 0
	leftWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(-15, 0, 0))
	w.AddObject(leftWall)

	rightWall := tracer.NewPlane()
	rightWall.Material().Color = tracer.ColorName(colornames.Lightgreen)
	rightWall.Material().Specular = 0
	rightWall.Material().Reflective = 0
	rightWall.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(15, 0, 0))
	w.AddObject(rightWall)

	// mirror1
	cube1 := tracer.NewUnitCube()
	cube1.SetTransform(
		tracer.IdentityMatrix().Scale(0.01, 1.5, 3).Translate(-2, 1.9, 0))
	cube1.Material().Reflective = 1
	cube1.Material().Color = tracer.ColorName(colornames.Black)
	w.AddObject(cube1)

	// border
	borderStripes := tracer.NewStripedPattern(
		tracer.ColorName(colornames.Lightgray), tracer.ColorName(colornames.White))
	borderStripes.SetTransform(tracer.IdentityMatrix().Scale(0.1, 1, 1).RotateY(math.Pi / 2))
	borderP := tracer.NewPertrubedPattern(borderStripes, 0.1)

	// top border
	topBorder := tracer.NewUnitCube()
	topBorder.SetTransform(tracer.IdentityMatrix().Scale(0.01, .2, 3).Translate(-2, 3.6, 0))
	topBorder.Material().SetPattern(borderP)
	w.AddObject(topBorder)

	// bottom border
	bottomBorder := tracer.NewUnitCube()
	bottomBorder.SetTransform(tracer.IdentityMatrix().Scale(0.01, .2, 3).Translate(-2, 0.2, 0))
	bottomBorder.Material().SetPattern(borderP)
	w.AddObject(bottomBorder)

	// left border
	leftBorder := tracer.NewUnitCube()
	leftBorder.SetTransform(tracer.IdentityMatrix().Scale(0.01, 1.9, 0.2).Translate(-2, 1.9, -3.2))
	leftBorder.Material().SetPattern(borderP)
	w.AddObject(leftBorder)

	// right border
	rightBorder := tracer.NewUnitCube()
	rightBorder.SetTransform(tracer.IdentityMatrix().Scale(0.01, 1.9, 0.2).Translate(-2, 1.9, 3.2))
	rightBorder.Material().SetPattern(borderP)
	w.AddObject(rightBorder)

	// table
	table := tracer.NewUnitCube()
	table.SetTransform(
		tracer.IdentityMatrix().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
	table.Material().Reflective = 0
	table.Material().Color = tracer.ColorName(colornames.Lightslategray)
	w.AddObject(table)

	// sphere1
	sphere1 := tracer.NewUnitSphere()
	sphere1.SetTransform(
		tracer.IdentityMatrix().Scale(.5, .5, .5).Translate(0, 1.5, 0)) // half sphere + full cube (scaled by half())
	// sphere1.Material().Color = tracer.ColorName(colornames.Yellow)
	sphere1pattern := tracer.NewStripedPattern(
		tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.Purple))
	sphere1pattern.SetTransform(tracer.IdentityMatrix().Scale(0.2, 1, 1))
	sphere1.Material().SetPattern(sphere1pattern)
	w.AddObject(sphere1)

	canvas := w.Render()

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}

func cube() {

	// width, height := 300.0, 300.0
	width, height := 500.0, 500.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 1

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(0, 10, -2), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(6, 2, -7)
	to := tracer.NewPoint(-3.5, 1, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// floor
	floor := tracer.NewPlane()
	// floor.Material().Color = tracer.ColorName(colornames.Gray)
	floor.Material().Specular = 0
	floor.Material().Reflective = 0
	floorP := tracer.NewCheckerPattern(tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	w.AddObject(floor)

	cube := tracer.NewUnitCube()
	cube.Material().Color = tracer.ColorName(colornames.Lightgreen)
	// cube.SetTransform(tracer.IdentityMatrix().Translate(0, 1, 0).Scale(1, .5, 1))
	cube.SetTransform(tracer.IdentityMatrix().Scale(1, .5, 1).Translate(0, 1, 0))
	w.AddObject(cube)

	canvas := w.Render()

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}

func glass() {

	width, height := 500.0, 500.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(2, 10, -2), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 4, -5)
	to := tracer.NewPoint(0, 0, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// floor
	floor := tracer.NewPlane()
	floor.Material().Specular = 0
	floor.Material().Reflective = 0
	// floor.Material().Transparency = 1.0
	// floor.Material().RefractiveIndex = 1.5
	floorP := tracer.NewCheckerPattern(tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	w.AddObject(floor)

	ball := tracer.NewGlassSphere()
	ball.SetTransform(tracer.IdentityMatrix().Translate(0, 1, 0))
	ball.Material().Color = tracer.Black()
	ball.Material().Diffuse = 0.0
	ball.Material().Ambient = 0.1
	ball.Material().Reflective = 0.0
	ball.Material().RefractiveIndex = 1.5
	ball.Material().Transparency = 1
	w.AddObject(ball)

	canvas := w.Render()

	// Export
	f, err := os.Create("image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}

func window() {

	// width, height := 200.0, 200.0
	width, height := 500.0, 500.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(2, 10, -2), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(3, 1, -7)
	to := tracer.NewPoint(-1, 0, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// floor
	floor := tracer.NewPlane()
	floor.Material().Specular = 0
	floor.Material().Reflective = 0
	// floor.Material().Transparency = 1.0
	// floor.Material().RefractiveIndex = 1.5
	floorP := tracer.NewCheckerPattern(tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	w.AddObject(floor)

	cube := tracer.NewUnitCube()
	cube.SetTransform(tracer.IdentityMatrix().Scale(0.2, 0.2, .2).Translate(-1.5, 0.2, -4))
	cube.Material().Color = tracer.ColorName(colornames.Red)
	w.AddObject(cube)

	// window
	wind := tracer.NewUnitCube()
	wind.SetTransform(tracer.IdentityMatrix().Scale(3.6, 1, 0.01).Translate(-1.5, 0, -3))
	wind.Material().Transparency = 1
	wind.Material().Reflective = 1
	wind.Material().RefractiveIndex = 1.5
	wind.Material().Ambient = 0.1
	wind.Material().Diffuse = 0.1
	wind.Material().Color = tracer.ColorName(colornames.Black)
	w.AddObject(wind)

	ball := tracer.NewUnitSphere()
	ball.SetTransform(tracer.IdentityMatrix().Translate(0, 1, 0))
	ball.Material().Color = tracer.ColorName(colornames.Burlywood)
	w.AddObject(ball)

	canvas := w.Render()

	// Export
	f, err := os.Create("/Users/dant/Downloads/image.png")
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Exporting canvas to %v", f.Name())
	canvas.ExportToPNG(f)
}

func pond() {

	// width, height := 100.0, 100.0
	width, height := 400.0, 400.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(3, 20, -35), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(3, 4, -38)
	to := tracer.NewPoint(0.7, 0, -33)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// surface
	surface := tracer.NewPlane()
	surface.Material().Specular = 0.0
	surface.Material().Diffuse = 0.1
	surface.Material().Ambient = 0.1
	surface.Material().Reflective = 1
	surface.Material().Transparency = 0.6
	surface.Material().RefractiveIndex = 1.3442
	surface.Material().Color = tracer.ColorName(colornames.White)
	surface.Material().ShadowCaster = false
	surfaceRealP1 := tracer.NewStripedPattern(
		tracer.ColorName(colornames.Lightgray), tracer.ColorName(colornames.Lightskyblue))
	surfaceRealP2 := tracer.NewStripedPattern(
		tracer.ColorName(colornames.Lightgray), tracer.ColorName(colornames.Lightskyblue))
	surfaceRealP2.SetTransform(tracer.IdentityMatrix().RotateY(math.Pi / 2))
	surfaceBlendedP := tracer.NewBlendedPattern(surfaceRealP1, surfaceRealP2)
	surfacePP := tracer.NewPertrubedPattern(surfaceBlendedP, 0.4)
	surface.Material().SetPattern(surfacePP)
	w.AddObject(surface)

	// bottom
	bottom := tracer.NewPlane()
	bottom.Material().Specular = 0
	bottom.Material().Color = tracer.ColorName(colornames.White)
	bottom.SetTransform(tracer.IdentityMatrix().Translate(0, -8, 0))
	bottomP := tracer.NewCheckerPattern(tracer.ColorName(colornames.Lightcoral),
		tracer.ColorName(colornames.Lightgray))
	// bottomP := tracer.NewGradientPattern(
	// 	tracer.ColorName(colornames.Lightgray), tracer.ColorName(colornames.Darkgrey))
	// bottomP.SetTransform(tracer.IdentityMatrix().Scale(2.5, 2.5, 2.5))
	bottom.Material().SetPattern(bottomP)
	w.AddObject(bottom)

	leftWall := tracer.NewPlane()
	leftWall.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(-40, 0, 0))
	leftWall.Material().Color = tracer.ColorName(colornames.Lightskyblue)
	leftWall.Material().Specular = 0
	leftWall.Material().Shininess = 200
	leftWall.Material().Ambient = 0.3
	leftWall.Material().Diffuse = 0
	leftWallP := tracer.NewRingPattern(tracer.ColorName(colornames.Lightsteelblue), tracer.White())
	leftWall.Material().SetPattern(leftWallP)
	w.AddObject(leftWall)

	backWall := tracer.NewPlane()
	backWall.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).Translate(0, 0, 4))
	backWall.Material().Color = tracer.ColorName(colornames.Lightskyblue)
	backWall.Material().Specular = 0
	backWall.Material().Shininess = 200
	backWall.Material().Ambient = 0.3
	backWall.Material().Diffuse = 0
	backWallP := tracer.NewRingPattern(tracer.ColorName(colornames.Lightsteelblue), tracer.White())
	// backWallP.SetTransform(tracer.IdentityMatrix().)
	backWall.Material().SetPattern(backWallP)
	w.AddObject(backWall)

	// below water red cube
	cube := tracer.NewUnitCube()
	cube.SetTransform(
		tracer.IdentityMatrix().Scale(0.4, 0.4, 0.4).RotateX(math.Pi/4).RotateY(math.Pi/4).RotateZ(math.Pi/4).Translate(1.5, -4, -34))
	cube.Material().Color = tracer.ColorName(colornames.Red)
	w.AddObject(cube)

	// half submerged yellow cube
	cube3 := tracer.NewUnitCube()
	cube3.SetTransform(tracer.IdentityMatrix().Scale(0.4, 0.4, 0.4).Translate(-0.5, 0, -34))
	cube3.Material().Color = tracer.ColorName(colornames.Yellow)
	w.AddObject(cube3)

	// below water yellow sphere
	ball := tracer.NewUnitSphere()
	ball.SetTransform(tracer.IdentityMatrix().Scale(0.8, 0.8, 0.8).Translate(4, -4, -30))
	ball.Material().Color = tracer.ColorName(colornames.Yellow)
	w.AddObject(ball)

	// above water lightblue cube
	cube2 := tracer.NewUnitCube()
	cube2.SetTransform(tracer.IdentityMatrix().Scale(0.4, 0.4, .4).Translate(1.7, 1, -32))
	cube2.Material().Color = tracer.ColorName(colornames.Lightblue)
	w.AddObject(cube2)

	canvas := w.Render()

	// Export
	f, err := os.Create("/Users/dant/Downloads/image.png")
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

	// testCanvas()
	// test1()
	// clock()
	// circle()
	// sphere()
	// scene()
	// colors()
	// mirrors()
	// mirror()
	// cube()
	// glass()
	// window()
	pond()

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
