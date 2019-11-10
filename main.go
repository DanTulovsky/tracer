package main

import (
	"flag"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"runtime/pprof"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "net/http/pprof"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"

	"github.com/DanTulovsky/tracer/tracer"
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

func sceneold() {

	// width, height := 300.0, 300.0
	width, height := 1200.0, 1000.0

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

	tracer.Render(w)
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

	width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

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
	floor.Material().Color = tracer.ColorName(colornames.Darkblue)
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
	cube1.Material().Color = tracer.ColorName(colornames.Black)
	w.AddObject(cube1)

	// mirror2
	cube2 := tracer.NewUnitCube()
	cube2.SetTransform(
		tracer.IdentityMatrix().Scale(0.001, 1, 5).Translate(2, 2, 0))
	cube2.Material().Reflective = 1
	cube2.Material().Color = tracer.ColorName(colornames.Black)
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

	tracer.Render(w)
}

func mirror() {

	// width, height := 300.0, 300.0
	width, height := 1200.0, 1000.0

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

	tracer.Render(w)
}

func cube() {

	width, height := 1000.0, 1000.0

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

	tracer.Render(w)
}

func glass() {

	// width, height := 500.0, 500.0
	width, height := 1200.0, 1000.0

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

	tracer.Render(w)
}

func window() {

	// width, height := 500.0, 500.0
	width, height := 1000.0, 1000.0

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

	tracer.Render(w)
}

func pond() {

	// width, height := 100.0, 100.0
	// width, height := 400.0, 400.0
	width, height := 1400.0, 1000.0

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

	tracer.Render(w)
}

func cylinder() {

	// width, height := 100.0, 100.0
	// width, height := 400.0, 400.0
	width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(3, 20, -10), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 11, -10)
	to := tracer.NewPoint(0, 3, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// floor
	floor := tracer.NewPlane()
	floor.Material().Specular = 0
	floor.Material().Reflective = 0
	// floor.Material().Transparency = 1.0
	// floor.Material().RefractiveIndex = 1.5
	floorP := tracer.NewCheckerPattern(
		tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	w.AddObject(floor)

	// closed
	c := tracer.NewClosedCylinder(0, 8)
	c.Material().Color = tracer.ColorName(colornames.Lightgreen)
	c.SetTransform(tracer.IdentityMatrix().Translate(0, 0, 0))
	w.AddObject(c)

	// open
	c2 := tracer.NewCylinder(0, 8)
	c2.Material().Color = tracer.ColorName(colornames.Lightblue)
	c2.SetTransform(tracer.IdentityMatrix().Translate(-3, 0, 0))
	w.AddObject(c2)

	// infinite
	c3 := tracer.NewDefaultCylinder()
	c3.Material().Color = tracer.ColorName(colornames.Lightcoral)
	c3.SetTransform(tracer.IdentityMatrix().Translate(3, 0, 0))
	w.AddObject(c3)

	// flipped & glass
	c4 := tracer.NewClosedCylinder(-4, 4)
	c4.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/3).RotateY(-math.Pi/4).Translate(0, 5.7, -4))
	c4.Material().Color = tracer.ColorName(colornames.Darkolivegreen)
	c4.Material().Transparency = 0.8
	c4.Material().Reflective = 0.5
	c4.Material().RefractiveIndex = 1.75
	c4.Material().Ambient = 0.1
	c4.Material().Diffuse = 0.1
	c4.Material().ShadowCaster = false
	w.AddObject(c4)

	// flipped & glass sphere
	// s := tracer.NewUnitSphere()
	// s.SetTransform(tracer.IdentityMatrix().Scale(2, 2, 2).Translate(0, 5.7, -4))
	// s.Material().Color = tracer.ColorName(colornames.Darkolivegreen)
	// s.Material().Transparency = 0.8
	// s.Material().Reflective = 0.5
	// s.Material().RefractiveIndex = 1.7 // force fish-eye affect
	// s.Material().Ambient = 0.1
	// s.Material().Diffuse = 0.1
	// s.Material().ShadowCaster = false
	// w.AddObject(s)

	tracer.Render(w)
}

func spherewarp() {

	// width, height := 200.0, 200.0
	width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(0, 10, 0), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 12, 0)
	to := tracer.NewPoint(0, 0, 0)
	up := tracer.NewVector(0, 0, 1) // note up vector change when looking straight down!
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	// floor
	floor := tracer.NewPlane()
	floor.Material().Specular = 0
	floor.Material().Reflective = 0
	// floor.Material().Transparency = 1.0
	// floor.Material().RefractiveIndex = 1.5
	floorP := tracer.NewCheckerPattern(
		tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	w.AddObject(floor)

	// glass sphere
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Scale(2.5, 2.5, 2.5).Translate(0, 4, 0))
	s.Material().Color = tracer.ColorName(colornames.White)
	s.Material().Transparency = 1.0
	s.Material().Reflective = 1.0
	s.Material().RefractiveIndex = 3
	// s.Material().Ambient = 0.1
	s.Material().Diffuse = 0
	s.Material().Specular = 0
	s.Material().ShadowCaster = false
	w.AddObject(s)

	tracer.Render(w)
}

func cone() {

	// width, height := 100.0, 100.0
	// width, height := 200.0, 200.0
	width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		// tracer.NewPointLight(tracer.NewPoint(2, 10, -1), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(-0.5, 1, 0), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 4.0, -4)
	to := tracer.NewPoint(0, 2.0, 0)
	up := tracer.NewVector(0, 1, 0)
	fov := math.Pi / 3.0

	camera := tracer.NewCamera(width, height, fov)
	cameraTransform := tracer.ViewTransform(from, to, up)
	camera.SetTransform(cameraTransform)

	w.SetCamera(camera)

	// floor
	floor := tracer.NewPlane()
	floor.Material().Specular = 0
	floor.Material().Reflective = 0.5
	// floor.Material().Transparency = 1.0
	// floor.Material().RefractiveIndex = 1.5
	floorP := tracer.NewCheckerPattern(
		tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	w.AddObject(floor)

	// mirror
	mirror := tracer.NewPlane()
	mirror.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).RotateY(math.Pi/4).Translate(0, 0, 4))
	mirror.Material().Reflective = 0.8
	mirror.Material().Diffuse = 0.01
	mirror.Material().Specular = 1
	mirror.Material().Ambient = 0.01
	mirror.Material().Color = tracer.Black()
	w.AddObject(mirror)

	// cone
	c := tracer.NewClosedCone(-1, 1)
	c.SetTransform(tracer.IdentityMatrix().Translate(-1, 1, 0))
	cp := tracer.NewStripedPattern(tracer.ColorName(colornames.Red), tracer.ColorName(colornames.White))
	cp.SetTransform(tracer.IdentityMatrix().Scale(0.1, 0.1, 0.1))
	cpp := tracer.NewPertrubedPattern(cp, 0.4)
	c.Material().SetPattern(cpp)
	w.AddObject(c)

	// orb
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Translate(-1, 3, 0))
	sp := tracer.NewStripedPattern(tracer.ColorName(colornames.Red), tracer.ColorName(colornames.Blue))
	sp.SetTransform(tracer.IdentityMatrix().Scale(0.3, 0.3, 0.3).RotateZ(math.Pi / 2))
	spp := tracer.NewPertrubedPattern(sp, 0.4)
	s.Material().SetPattern(spp)
	s.Material().Transparency = 0.8
	s.Material().Ambient = 0.1
	s.Material().Diffuse = 0.1
	s.Material().ShadowCaster = false
	s.Material().RefractiveIndex = 1.53
	s.Material().Reflective = 1
	w.AddObject(s)

	tracer.Render(w)
}

func group() {

	// width, height := 100.0, 100.0
	// width, height := 400.0, 400.0
	width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(3, 20, -10), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 11, -10)
	to := tracer.NewPoint(0, 3, 0)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	g := tracer.NewGroup()

	// floor
	floor := tracer.NewPlane()
	floor.Material().Specular = 0
	floor.Material().Reflective = 0
	// floor.Material().Transparency = 1.0
	// floor.Material().RefractiveIndex = 1.5
	floorP := tracer.NewCheckerPattern(
		tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.Yellow))
	floor.Material().SetPattern(floorP)
	g.AddMember(floor)

	// closed
	c := tracer.NewClosedCylinder(0, 8)
	c.Material().Color = tracer.ColorName(colornames.Lightgreen)
	c.SetTransform(tracer.IdentityMatrix().Translate(0, 0, 0))
	g.AddMember(c)

	// open
	c2 := tracer.NewCylinder(0, 8)
	c2.Material().Color = tracer.ColorName(colornames.Lightblue)
	c2.SetTransform(tracer.IdentityMatrix().Translate(-3, 0, 0))
	g.AddMember(c2)

	// infinite
	c3 := tracer.NewDefaultCylinder()
	c3.Material().Color = tracer.ColorName(colornames.Lightcoral)
	c3.SetTransform(tracer.IdentityMatrix().Translate(3, 0, 0))
	g.AddMember(c3)

	// flipped & glass
	c4 := tracer.NewClosedCylinder(-4, 4)
	c4.SetTransform(
		tracer.IdentityMatrix().RotateZ(math.Pi/3).RotateY(-math.Pi/4).Translate(0, 5.7, -4))
	c4.Material().Color = tracer.ColorName(colornames.Darkolivegreen)
	c4.Material().Transparency = 0.8
	c4.Material().Reflective = 0.5
	c4.Material().RefractiveIndex = 1.75
	c4.Material().Ambient = 0.1
	c4.Material().Diffuse = 0.1
	c4.Material().ShadowCaster = false
	g.AddMember(c4)

	// flipped & glass sphere
	// s := tracer.NewUnitSphere()
	// s.SetTransform(tracer.IdentityMatrix().Scale(2, 2, 2).Translate(0, 5.7, -4))
	// s.Material().Color = tracer.ColorName(colornames.Darkolivegreen)
	// s.Material().Transparency = 0.8
	// s.Material().Reflective = 0.5
	// s.Material().RefractiveIndex = 1.7 // force fish-eye affect
	// s.Material().Ambient = 0.1
	// s.Material().Diffuse = 0.1
	// s.Material().ShadowCaster = false
	// g.AddMember(s)

	// g.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi / 2))
	w.AddObject(g)

	tracer.Render(w)
}

func triangle() {
	// width, height := 400.0, 300.0
	width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(3, 4, -30), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-5, 4, -1), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 3, -4)
	to := tracer.NewPoint(0, 0, 4)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)

	floor := tracer.NewPlane()
	floor.SetTransform(tracer.IdentityMatrix().Translate(0, -3, 0))
	floorp := tracer.NewRingPattern(tracer.ColorName(colornames.Red), tracer.White())
	floor.Material().SetPattern(floorp)
	w.AddObject(floor)

	ceiling := tracer.NewPlane()
	ceiling.SetTransform(tracer.IdentityMatrix().Translate(0, 8, 0))
	ceiling.Material().Color = tracer.ColorName(colornames.Lightskyblue)
	w.AddObject(ceiling)

	backWall := tracer.NewPlane()
	backWall.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).Translate(0, 0, 40))
	backWallp := tracer.NewStripedPattern(tracer.ColorName(colornames.Blue), tracer.White())
	backWall.Material().SetPattern(backWallp)
	backWall.Material().Specular = 0.2
	w.AddObject(backWall)

	s1 := tracer.NewUnitSphere()
	s1.SetTransform(tracer.IdentityMatrix().Scale(2, 2, 2).Translate(0, 1, 3))
	s1.Material().Color = tracer.ColorName(colornames.Black)
	s1.Material().Reflective = 1
	w.AddObject(s1)

	g1 := tracer.NewGroup()
	g1.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/8).Translate(0.7, 0.4, 0))
	w.AddObject(g1)

	t1 := tracer.NewTriangle(tracer.NewPoint(0, 0, 0), tracer.NewPoint(2, 0, 0), tracer.NewPoint(1, 2, 0))
	t1.Material().Color = tracer.ColorName(colornames.Darkred)
	t1.Material().Transparency = 0.4
	t1.Material().Diffuse = 0.1
	t1.Material().Ambient = 0.1
	t1.Material().ShadowCaster = false
	g1.AddMember(t1)

	t2 := tracer.NewTriangle(tracer.NewPoint(0, 0, 0), tracer.NewPoint(-2, 0, 0), tracer.NewPoint(-1, 2, 0))
	t2.Material().Color = tracer.ColorName(colornames.Darkblue)
	t2.Material().Transparency = 0.5
	t2.Material().Diffuse = 0.1
	t2.Material().Ambient = 0.1
	t2.Material().ShadowCaster = false
	g1.AddMember(t2)

	t3 := tracer.NewTriangle(tracer.NewPoint(0, 0, 0), tracer.NewPoint(-1, 2, 0), tracer.NewPoint(1, 2, 0))
	t3.Material().Color = tracer.ColorName(colornames.Darkgreen)
	t3.Material().Transparency = 1
	t3.Material().Diffuse = 0.1
	t3.Material().Ambient = 0.1
	t3.Material().ShadowCaster = false
	g1.AddMember(t3)

	tracer.Render(w)
}

func objParse(f string) {

	width, height := 640.0, 480.0
	// width, height := 1400.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(0, 3, -10), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-10, -3, -5), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 2, -8)
	to := tracer.NewPoint(0, 0, 4)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 3)

	g, err := tracer.ParseOBJ(f)
	if err != nil {
		log.Fatalln(err)
	}

	g.SetTransform(tracer.IdentityMatrix().Translate(0, 2, 0))

	w.AddObject(g)
	tracer.Render(w)
}

func csg() {
	// width, height := 640.0, 480.0
	width, height := 1400.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(0, 20, -35), tracer.NewColor(1, 1, 1)),
		// tracer.NewPointLight(tracer.NewPoint(-10, -4, -1), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 6, -8)
	to := tracer.NewPoint(0, 0, 4)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 3)

	s1 := tracer.NewClosedCylinder(-5, 5)
	s1.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi / 2))
	s1.Material().Color = tracer.ColorName(colornames.Lightcyan)

	s2 := tracer.NewUnitSphere()
	s2.SetTransform(tracer.IdentityMatrix().Scale(1, 2, 1))
	s2.Material().Color = tracer.ColorName(colornames.Lightcoral)

	op := tracer.Difference
	csg := tracer.NewCSG(s1, s2, op)
	csg.SetTransform(tracer.IdentityMatrix().Scale(1, 0.5, 1).Translate(0, 2, 0))

	w.AddObject(csg)

	tracer.Render(w)
}

func env() *tracer.World {
	width, height := 640.0, 480.0
	// width, height := 1000.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(3, 4, -3), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(-9, 10, 10), tracer.NewColor(1, 1, 1)),
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

func glassplane() *tracer.Plane {
	p := tracer.NewPlane()
	p.Material().Specular = 0.0
	p.Material().Diffuse = 0.1
	p.Material().Ambient = 0.1
	p.Material().Reflective = 1
	p.Material().Transparency = 0.6
	p.Material().RefractiveIndex = 1.3442
	p.Material().Color = tracer.ColorName(colornames.White)
	p.Material().ShadowCaster = false

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

	tracer.Render(w)
}

func plane() {
	w := env()
	w.AddObject(floor())
	w.AddObject(rightWall())

	tracer.Render(w)
}

func simplecone() {
	w := env()

	cone1 := tracer.NewClosedCone(-2, 0)
	cone1.SetTransform(tracer.IdentityMatrix().Scale(1, 3, 1).Translate(0, 2, 2))

	w.AddObject(cone1)
	w.AddObject(floor())
	tracer.Render(w)
}

func simplecylinder() {
	w := env()

	cylinder1 := tracer.NewClosedCylinder(0, 2)
	cylinder1.SetTransform(tracer.IdentityMatrix().Scale(0.5, 1, 0.5).RotateY(math.Pi/2).Translate(0, 0, 0))

	w.AddObject(cylinder1)
	w.AddObject(floor())
	tracer.Render(w)
}

func glasssphere() *tracer.Sphere {
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Scale(.75, .75, .75).Translate(0, 1.75, 0))
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
	s.SetTransform(tracer.IdentityMatrix().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
	s.Material().Color = tracer.ColorName(colornames.Gold)
	up := tracer.NewUVCheckersPattern(8, 8,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Violet))
	cp := tracer.NewTextureMapPattern(up, tracer.NewCubeMapSame(up))
	p := tracer.NewPertrubedPattern(cp, 0.09)
	s.Material().SetPattern(p)

	return s
}

func shapes() {

	w := envxy(2000, 1600)

	floory := -3.3

	sphere1 := tracer.NewUnitSphere()
	sphere1.SetTransform(tracer.IdentityMatrix().Translate(-4.5, 1, 5))
	mapper := tracer.NewSphericalMap()
	uvpattern := tracer.NewUVCheckersPattern(20, 10,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Black))
	pattern := tracer.NewTextureMapPattern(uvpattern, mapper)
	sphere1.Material().SetPattern(pattern)

	cube1 := tracer.NewUnitCube()
	cube1.SetTransform(tracer.IdentityMatrix().RotateY(math.Pi/4).Translate(5.8, 1, 9))
	left := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Yellow),
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Blue),
		tracer.ColorName(colornames.Brown))
	front := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Yellow),
		tracer.ColorName(colornames.Brown),
		tracer.ColorName(colornames.Green))
	right := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Yellow),
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Green),
		tracer.ColorName(colornames.White))
	back := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Green),
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.White),
		tracer.ColorName(colornames.Blue))
	up := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Brown),
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Yellow))
	down := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Brown),
		tracer.ColorName(colornames.Green),
		tracer.ColorName(colornames.Blue),
		tracer.ColorName(colornames.White))
	pattern2 := tracer.NewCubeMapPattern(left, front, right, back, up, down)
	cube1.Material().SetPattern(pattern2)

	cylinder1 := tracer.NewClosedCylinder(-4, 4)
	cylinder1.SetTransform(
		tracer.IdentityMatrix().Scale(0.5, 0.5, 0.5).RotateZ(math.Pi/2).Translate(0, 0.5, 0))
	mapper3 := tracer.NewCylinderMap()
	uvpattern3 := tracer.NewUVCheckersPattern(10, 2,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Blue))
	pattern3 := tracer.NewTextureMapPattern(uvpattern3, mapper3)
	cylinder1.Material().SetPattern(pattern3)

	backWall1 := glassplane()
	// backWall1 := floor()
	backWall1.SetTransform(
		tracer.IdentityMatrix().RotateX(math.Pi/2).Translate(0, 0, 20))
	mapper4 := tracer.NewPlaneMap()
	uvpattern4 := tracer.NewUVCheckersPattern(2, 2,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Orange))
	pattern4 := tracer.NewTextureMapPattern(uvpattern4, mapper4)
	pattern4.SetTransform(tracer.IdentityMatrix().Scale(5, 5, 5))
	backWall1.Material().SetPattern(pattern4)

	cone1 := tracer.NewClosedCone(-2, 0)
	cone1.SetTransform(
		tracer.IdentityMatrix().Scale(0.3, 1, 0.3).Translate(-2.3, 2, 7))
	mapper5 := tracer.NewCylinderMap()
	uvpattern5 := tracer.NewUVCheckersPattern(2, 2,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Goldenrod))
	pattern5 := tracer.NewTextureMapPattern(uvpattern5, mapper5)
	pattern5.SetTransform(tracer.IdentityMatrix().Scale(0.5, 0.5, 0.5))
	cone1.Material().SetPattern(pattern5)

	csgMember1 := tracer.NewUnitCube()
	pattern6 := tracer.NewCubeMapPattern(left, front, right, back, up, down)
	csgMember1.Material().SetPattern(pattern6)

	csgMember2 := tracer.NewUnitSphere()
	csgMember2.SetTransform(tracer.IdentityMatrix().Translate(0, 0.2, 0))
	mapper7 := tracer.NewSphericalMap()
	uvpattern7 := tracer.NewUVCheckersPattern(30, 15,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Goldenrod))
	pattern7 := tracer.NewTextureMapPattern(uvpattern7, mapper7)
	pattern7.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi / 4))
	csgMember2.Material().SetPattern(pattern7)

	csg1 := tracer.NewCSG(csgMember1, csgMember2, tracer.Intersect)
	csg1.SetTransform(tracer.IdentityMatrix().Translate(4.0, 1, 2))

	// earth
	earth := tracer.NewUnitSphere()
	earth.SetTransform(tracer.IdentityMatrix().RotateY(math.Pi/2).Translate(-4, 3.5, 25))
	image := "images/earthmap1k.jpg"
	earthup, err := tracer.NewUVImagePattern(image)
	if err != nil {
		log.Fatal(err)
	}
	mapperearth := tracer.NewSphericalMap()
	p := tracer.NewTextureMapPattern(earthup, mapperearth)
	earth.Material().SetPattern(p)

	flr := floor()
	flr.SetTransform(tracer.IdentityMatrix().Translate(0, floory, 0))

	g := tracer.NewGroup()
	g.AddMembers(csg1, cone1, sphere1, cube1, cylinder1)
	g.SetTransform(tracer.IdentityMatrix().Translate(0, floory, 4))

	g2 := tracer.NewGroup()
	g2.AddMembers(glasssphere(), pedestal())
	g2.SetTransform(tracer.IdentityMatrix().Translate(0, 0, 4))

	// skybox
	w.AddObject(skyboxcube("field1"))
	w.AddObject(earth)
	w.AddObject(g)
	w.AddObject(g2)
	// w.AddObject(csg1)
	// w.AddObject(cone1)
	// w.AddObject(sphere1)
	// w.AddObject(cube1)
	// w.AddObject(cylinder1)
	w.AddObject(flr)
	// w.AddObject(ceiling())
	w.AddObject(backWall1)

	tracer.Render(w)

}

func simplesphere() {
	w := env()

	sphere1 := tracer.NewUnitSphere()
	sphere1.SetTransform(tracer.IdentityMatrix().Translate(0, 1, 1))
	mapper := tracer.NewSphericalMap()
	uvpattern := tracer.NewUVCheckersPattern(20, 10,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Black))
	pattern := tracer.NewTextureMapPattern(uvpattern, mapper)
	sphere1.Material().SetPattern(pattern)

	floor := tracer.NewPlane()
	mapper2 := tracer.NewPlaneMap()
	uvpattern2 := tracer.NewUVCheckersPattern(2, 2,
		tracer.ColorName(colornames.Red), tracer.ColorName(colornames.Blue))
	pattern2 := tracer.NewTextureMapPattern(uvpattern2, mapper2)
	floor.Material().SetPattern(pattern2)

	w.AddObject(sphere1)
	w.AddObject(floor)
	tracer.Render(w)
}

func texttureMap() {

	w := env()

	tracer.Render(w)
}

func cubeMap() {

	w := env()

	left := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Yellow),
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Blue),
		tracer.ColorName(colornames.Brown))
	front := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Yellow),
		tracer.ColorName(colornames.Brown),
		tracer.ColorName(colornames.Green))
	right := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Yellow),
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Green),
		tracer.ColorName(colornames.White))
	back := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Green),
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.White),
		tracer.ColorName(colornames.Blue))
	up := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Brown),
		tracer.ColorName(colornames.Cyan),
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Red),
		tracer.ColorName(colornames.Yellow))
	down := tracer.NewUVAlignCheckPattern(
		tracer.ColorName(colornames.Purple),
		tracer.ColorName(colornames.Brown),
		tracer.ColorName(colornames.Green),
		tracer.ColorName(colornames.Blue),
		tracer.ColorName(colornames.White))

	cube := tracer.NewUnitCube()
	cube.SetTransform(
		tracer.IdentityMatrix().Scale(0.8, 0.8, 0.8).RotateX(math.Pi/2).RotateY(math.Pi/6).RotateZ(math.Pi/6).Translate(0, 1.7, 0))
	p := tracer.NewCubeMapPattern(left, front, right, back, up, down)
	cube.Material().SetPattern(p)

	w.AddObject(cube)
	tracer.Render(w)
}

func image1() {
	w := env()
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().RotateY(math.Pi/2).Translate(0, 2, 0))

	image := "images/earthmap1k.jpg"

	up, err := tracer.NewUVImagePattern(image)
	if err != nil {
		log.Fatal(err)
	}
	mapper := tracer.NewSphericalMap()
	p := tracer.NewTextureMapPattern(up, mapper)
	s.Material().SetPattern(p)

	w.AddObject(s)
	tracer.Render(w)
}

func skyboxcube(folder string) *tracer.Cube {

	sb := tracer.NewUnitCube()
	left, _ := tracer.NewUVImagePattern(path.Join("images/skybox/", folder, "negx.jpg"))
	right, _ := tracer.NewUVImagePattern(path.Join("images/skybox/", folder, "posx.jpg"))
	front, _ := tracer.NewUVImagePattern(path.Join("images/skybox", folder, "negz.jpg"))
	back, _ := tracer.NewUVImagePattern(path.Join("images/skybox/", folder, "posz.jpg"))
	up, _ := tracer.NewUVImagePattern(path.Join("images/skybox", folder, "posy.jpg"))
	down, _ := tracer.NewUVImagePattern(path.Join("images/skybox", folder, "negy.jpg"))

	p := tracer.NewCubeMapPattern(left, front, right, back, up, down)
	sb.Material().Ambient = 1
	sb.Material().Specular = 0
	sb.Material().Diffuse = 0
	sb.Material().SetPattern(p)
	sb.SetTransform(tracer.IdentityMatrix().Scale(50, 50, 50))
	return sb
}

func skyboxcube1(folder string) {
	w := envxy(1000, 1000)

	sb := tracer.NewUnitCube()
	left, _ := tracer.NewUVImagePattern(path.Join("images/skybox/", folder, "negx.jpg"))
	right, _ := tracer.NewUVImagePattern(path.Join("images/skybox/", folder, "posx.jpg"))
	front, _ := tracer.NewUVImagePattern(path.Join("images/skybox", folder, "negz.jpg"))
	back, _ := tracer.NewUVImagePattern(path.Join("images/skybox/", folder, "posz.jpg"))
	up, _ := tracer.NewUVImagePattern(path.Join("images/skybox", folder, "posy.jpg"))
	down, _ := tracer.NewUVImagePattern(path.Join("images/skybox", folder, "negy.jpg"))

	p := tracer.NewCubeMapPattern(left, front, right, back, up, down)
	sb.Material().Ambient = 1
	sb.Material().Specular = 0
	sb.Material().Diffuse = 0
	sb.Material().SetPattern(p)
	sb.SetTransform(tracer.IdentityMatrix().Scale(10, 10, 10))

	sphere := tracer.NewGlassSphere()
	sphere.Material().Diffuse = 0.1
	sphere.Material().Reflective = 0.5
	sphere.Material().Color = tracer.Black()
	sphere.SetTransform(tracer.IdentityMatrix().Translate(0, 2, 0))

	w.AddObject(sphere)
	w.AddObject(sb)

	tracer.Render(w)
}

func skyboxsphere1(input string) {
	w := envxy(1600, 1000)

	sb := tracer.NewUnitSphere()
	sb.Material().Ambient = 1
	sb.Material().Specular = 0
	sb.Material().Diffuse = 0
	sb.SetTransform(tracer.IdentityMatrix().Scale(10, 10, 10))

	filename := path.Join("images/hdri", input)
	m := tracer.HDRToImage(filename)

	up, err := tracer.NewUVImagePatternImage(m)
	if err != nil {
		log.Fatal(err)
	}
	p := tracer.NewTextureMapPattern(up, tracer.NewSphericalMap())
	sb.Material().SetPattern(p)

	sphere := tracer.NewGlassSphere()
	sphere.Material().Diffuse = 0.1
	sphere.Material().Reflective = 0.5
	sphere.Material().Color = tracer.Black()
	sphere.SetTransform(tracer.IdentityMatrix().Translate(0, 2, 0))

	w.AddObject(sphere)
	w.AddObject(sb)
	tracer.Render(w)
}

func envxy(width, height float64) *tracer.World {
	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(3, 10, -3), tracer.NewColor(1, 1, 1)),
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
func main() {

	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// skyboxcube1("field1")
	// skyboxsphere1("shanghai_bund_4k.hdr")
	// image1()
	// textureMap()
	// cubeMap()

	// scene()
	// plane()
	shapes()

	// colors()
	// mirrors()
	// mirror()
	// cube()
	// glass()
	// window()
	// pond()
	// spherewarp()
	// cylinder()
	// cone()
	// simplecone()
	// simplecylinder()
	// group()
	// triangle()
	// https://octolinker-demo.now.sh/mokiat/go-data-front
	// csg()
	// simplesphere()

	// dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/obj"))
	// f := path.Join(dir, "complex-smooth4.obj")
	// objParse(f)

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
