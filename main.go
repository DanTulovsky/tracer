package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"runtime"
	"runtime/pprof"

	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "net/http/pprof"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/colornames"

	"github.com/DanTulovsky/tracer/tracer"
	"github.com/DanTulovsky/tracer/utils"
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
	// p := tracer.NewPerturbedPattern(
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
	pRightWall := tracer.NewPerturbedPattern(
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
	p3 := tracer.NewPerturbedPattern(p2, 0.6)
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
	borderP := tracer.NewPerturbedPattern(borderStripes, 0.1)

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
	surfacePP := tracer.NewPerturbedPattern(surfaceBlendedP, 0.4)
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
	cpp := tracer.NewPerturbedPattern(cp, 0.4)
	c.Material().SetPattern(cpp)
	w.AddObject(c)

	// orb
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Translate(-1, 3, 0))
	sp := tracer.NewStripedPattern(tracer.ColorName(colornames.Red), tracer.ColorName(colornames.Blue))
	sp.SetTransform(tracer.IdentityMatrix().Scale(0.3, 0.3, 0.3).RotateZ(math.Pi / 2))
	spp := tracer.NewPerturbedPattern(sp, 0.4)
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

	// width, height := 640.0, 480.0
	width, height := 1200.0, 1000.0

	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)

	// override light here
	w.SetLights([]tracer.Light{
		tracer.NewPointLight(tracer.NewPoint(10, 50, -30), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(-10, -3, -5), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(1, 7, -14)
	to := tracer.NewPoint(0, 0, 4)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 3)

	g, err := tracer.ParseOBJ(f)
	if err != nil {
		log.Fatalln(err)
	}

	// g.SetTransform(tracer.IdentityMatrix().RotateY(math.Pi/5).RotateX(math.Pi/3).Translate(0, 2, 0))
	g.SetTransform(tracer.IdentityMatrix().Scale(2.5, 2.5, 2.5).RotateY(math.Pi/7).Translate(0, 2, 0))

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

func floor(y float64) *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().Translate(0, y, 0))
	pp := tracer.NewCheckerPattern(
		tracer.ColorName(colornames.Gray), tracer.ColorName(colornames.White))
	p.Material().SetPattern(pp)

	return p
}

func ceiling(y float64) *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().Translate(0, y, 0))
	pp := tracer.NewGradientPattern(
		tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.Red))
	pp.SetTransform(tracer.IdentityMatrix().Scale(10, 1, 1).Translate(-15, 0, 0))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0
	p.Material().Ambient = 0.15

	return p
}

func backWallGhost(z float64) *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(
		tracer.IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	ppuv, err := tracer.NewUVImagePattern("/Users/dant/Downloads/ghost.png")
	if err != nil {
		log.Fatal(err)
	}
	pp := tracer.NewTextureMapPattern(ppuv, tracer.NewPlaneMap())
	pp.SetTransform(tracer.IdentityMatrix().Scale(10, 5, 5).RotateY(math.Pi/2).Translate(0, 0, -3))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}

func backWall(z float64) *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(
		tracer.IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	p.Material().Color = tracer.ColorName(colornames.Lightpink)
	p.Material().Specular = 0

	return p
}
func frontWall(z float64) *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(
		tracer.IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, z))
	uvpp := tracer.NewUVCheckersPattern(4, 4,
		tracer.ColorName(colornames.Orange), tracer.ColorName(colornames.White))
	pp := tracer.NewTextureMapPattern(uvpp, tracer.NewPlaneMap())
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func rightWall(x float64) *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(x, 0, 0))
	pp := tracer.NewGradientPattern(
		tracer.ColorName(colornames.Orange), tracer.ColorName(colornames.White))
	pp.SetTransform(tracer.IdentityMatrix().Scale(10, 1, 1).Translate(-5, 0, 0))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	return p
}
func leftWall(x float64) *tracer.Plane {
	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/2).Translate(x, 0, 0))
	pp := tracer.NewStripedPattern(
		tracer.ColorName(colornames.Lightskyblue), tracer.ColorName(colornames.White))
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

func simplecone() {
	w := envxy(640, 480)
	w.Config.SoftShadows = true

	cone1 := tracer.NewClosedCone(-2, 0)
	cone1.SetTransform(tracer.IdentityMatrix().Scale(1, 3, 1).Translate(0, 2, 2))

	w.AddObject(cone1)
	w.AddObject(floor(0))
	w.AddObject(backWall(10))
	tracer.Render(w)
}

func simplecylinder() {
	w := env()

	cylinder1 := tracer.NewClosedCylinder(0, 2)
	cylinder1.SetTransform(
		tracer.IdentityMatrix().Scale(0.5, 1, 0.5).RotateY(math.Pi/2).Translate(0, 0, 1))
	uvp := tracer.NewUVCheckersPattern(12, 6, tracer.Black(), tracer.White())
	p := tracer.NewTextureMapPattern(uvp, tracer.NewCylinderMap())
	cylinder1.Material().SetPattern(p)

	w.AddObject(cylinder1)
	w.AddObject(floor(0))
	tracer.Render(w)
}

func cylindertextures() {
	w := env()

	cylinder1 := tracer.NewClosedCylinder(0, 2)
	cylinder1.SetTransform(
		tracer.IdentityMatrix().Scale(0.5, 1, 0.5).RotateY(math.Pi/2).Translate(0, 0, 1))
	uvp1, _ := tracer.NewUVImagePattern("images/checker.jpg")
	p1 := tracer.NewTextureMapPattern(uvp1, tracer.NewCylinderMap())
	cylinder1.Material().SetPattern(p1)

	cylinder2 := tracer.NewClosedCylinder(0, 2)
	cylinder2.SetTransform(
		tracer.IdentityMatrix().Scale(0.5, 1, 0.5).RotateY(math.Pi/2).Translate(-1.5, 0, 1))
	uvp2 := tracer.NewUVCheckersPattern(12, 6, tracer.ColorName(colornames.Green), tracer.White())
	p2 := tracer.NewTextureMapPattern(uvp2, tracer.NewSphericalMap())
	cylinder2.Material().SetPattern(p2)

	cylinder3 := tracer.NewClosedCylinder(0, 2)
	cylinder3.SetTransform(
		tracer.IdentityMatrix().Scale(0.5, 1, 0.5).RotateY(math.Pi/2).Translate(1.5, 0, 1))
	uvp3 := tracer.NewUVCheckersPattern(2, 2, tracer.ColorName(colornames.Blue), tracer.White())
	p3 := tracer.NewTextureMapPattern(uvp3, tracer.NewPlaneMap())
	cylinder3.Material().SetPattern(p3)

	w.AddObject(cylinder1)
	w.AddObject(cylinder2)
	w.AddObject(cylinder3)
	w.AddObject(floor(0))

	tracer.Render(w)
}
func mirorsphere() *tracer.Sphere {
	s := tracer.NewUnitSphere()
	s.SetTransform(tracer.IdentityMatrix().Scale(.75, .75, .75).Translate(0, 1.75, 0))
	s.Material().Ambient = 0
	s.Material().Diffuse = 0
	s.Material().Reflective = 1.0
	s.Material().Transparency = 0
	s.Material().ShadowCaster = true
	// s.Material().RefractiveIndex = 1.573

	return s
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

// mirror cube at x,y,z scaled by xs ys, zs and roated by rx, ry, rz
func mirrorcube(x, y, z, xs, ys, zs, rx, ry, rz float64) tracer.Shaper {
	c := tracer.NewUnitCube()
	c.SetTransform(tracer.IdentityMatrix().Scale(xs, ys, zs).RotateX(rx).RotateY(ry).RotateZ(rz).Translate(x, y, z))

	c.Material().Ambient = 0
	c.Material().Diffuse = 0
	c.Material().Reflective = 1.0
	c.Material().Transparency = 0
	c.Material().ShadowCaster = true

	return c
}

func pedestal() *tracer.Cube {
	s := tracer.NewUnitCube()
	s.SetTransform(tracer.IdentityMatrix().Scale(0.5, 0.5, 0.5).Translate(0, 0.5, 0))
	s.Material().Color = tracer.ColorName(colornames.Gold)
	up := tracer.NewUVCheckersPattern(8, 8,
		tracer.ColorName(colornames.White), tracer.ColorName(colornames.Violet))
	cp := tracer.NewTextureMapPattern(up, tracer.NewCubeMapSame(up))
	p := tracer.NewPerturbedPattern(cp, 0.09)
	s.Material().SetPattern(p)

	return s
}

func sphereOnPedestal() *tracer.Group {
	g := tracer.NewGroup()
	g.AddMembers(glasssphere(), pedestal())
	return g
}

func mirrorSphereOnPedestal() *tracer.Group {
	g := tracer.NewGroup()
	g.AddMembers(mirorsphere(), pedestal())
	return g
}
func shapes() {

	w := envxy(800, 600)
	w.Camera().SetFoV(math.Pi / 2.5)

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
	// backWall1.Material().Specular = 0
	// backWall1.Material().Diffuse = 1
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
		tracer.IdentityMatrix().Scale(0.3, 1, 0.3).Translate(-2.8, 2, 7))
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

	flr := floor(0)
	flr.SetTransform(tracer.IdentityMatrix().Translate(0, floory, 0))

	// mirror on the left
	lwall := tracer.NewUnitCube()
	lwall.SetTransform(tracer.IdentityMatrix().Scale(0.1, 20, 10).Translate(-7, 0, 10))
	lwall.Material().Color = tracer.ColorName(colornames.Black)
	lwall.Material().Reflective = 0.7

	g := tracer.NewGroup()
	g.AddMembers(csg1, cone1, sphere1, cube1, cylinder1)
	g.SetTransform(tracer.IdentityMatrix().Translate(0, floory, 4))

	g2 := tracer.NewGroup()
	g2.AddMembers(glasssphere(), pedestal())
	g2.SetTransform(tracer.IdentityMatrix().Scale(2, 2, 2).RotateY(math.Pi/4).Translate(0, floory, 7))

	// skybox
	w.AddObject(skyboxcube("field1"))
	w.AddObject(earth)
	w.AddObject(g)
	w.AddObject(g2)
	w.AddObject(flr)
	w.AddObject(lwall)
	// w.AddObject(ceiling())
	w.AddObject(backWall1)

	tracer.Render(w)

}

func simplesphere() {
	w := envxy(640, 480)
	// w.Config.Parallelism = 1
	// w.Camera().SetFoV(math.Pi / 2.0)

	sphere1 := tracer.NewUnitSphere()
	sphere1.SetTransform(
		tracer.IdentityMatrix().Scale(2.3, 2.3, 2.3).Translate(0, 2.3, 1))
	// mapper := tracer.NewSphericalMap()
	// uvpattern := tracer.NewUVCheckersPattern(20, 10,
	// 	tracer.ColorName(colornames.White), tracer.ColorName(colornames.Gray))
	// pattern := tracer.NewTextureMapPattern(uvpattern, mapper)
	// sphere1.Material().SetPattern(pattern)
	sphere1.Material().Color = tracer.ColorName(colornames.Red)
	pert := tracer.NewNoisePerturber(sphere1, 1)
	pert.SetTransform(tracer.IdentityMatrix().Scale(.15, .15, .15))

	sphere1.Material().SetPerturber(pert)

	w.AddObject(sphere1)
	w.AddObject(floor(0))
	w.AddObject(backWall(50))
	tracer.Render(w)
}

func heightmapplane(filename string) {
	w := envxy(640, 480)

	plane := tracer.NewPlane()
	plane.SetTransform(tracer.IdentityMatrix().Scale(30, 30, 30).Translate(-13, 0, 8))
	plane.Material().Specular = 0
	plane.Material().Color = tracer.ColorName(colornames.Lightblue)
	mapper := tracer.NewPlaneMap()
	pert, err := tracer.NewImageHeightmapPerturber(filename, mapper, plane)
	if err != nil {
		log.Fatal(err)
	}
	pert.SetTransform(tracer.IdentityMatrix().Scale(1, 1, 1))
	plane.Material().SetPerturber(pert)

	w.AddObject(plane)
	// w.AddObject(floor(0))
	w.AddObject(backWall(50))
	tracer.Render(w)
}

func brickwall(dir string) {
	w := envxy(1024, 768)
	basecolor := path.Join(dir, "basecolor.jpg")
	heightmap := path.Join(dir, "height.jpg")

	plane := tracer.NewPlane()
	plane.SetTransform(tracer.IdentityMatrix().Scale(3, 3, 3).RotateX(math.Pi/2).Translate(0, 0, 0))
	plane.Material().Specular = 0
	mapper := tracer.NewPlaneMap()

	// The image pattern
	pp, err := tracer.NewUVImagePattern(basecolor)
	if err != nil {
		log.Fatal(err)
	}
	pattern := tracer.NewTextureMapPattern(pp, mapper)
	plane.Material().SetPattern(pattern)

	// Heightmap
	pert, err := tracer.NewImageHeightmapPerturber(heightmap, mapper, plane)
	if err != nil {
		log.Fatal(err)
	}
	pert.SetTransform(tracer.IdentityMatrix().Scale(1, 1, 1))
	plane.Material().SetPerturber(pert)

	w.AddObject(plane)
	// w.AddObject(floor(0))
	// w.AddObject(backWall(50))
	tracer.Render(w)
}

func heightmapsphere(filename string) {
	w := envxy(640, 480)
	// w.Config.Parallelism = 1
	// w.Camera().SetFoV(math.Pi / 2.0)

	sphere1 := tracer.NewUnitSphere()
	sphere1.SetTransform(
		tracer.IdentityMatrix().Scale(2.3, 2.3, 2.3).Translate(0, 2.3, 1))
	mapper := tracer.NewSphericalMap()
	// uvpattern := tracer.NewUVCheckersPattern(20, 10,
	// 	tracer.ColorName(colornames.White), tracer.ColorName(colornames.Gray))
	// pattern := tracer.NewTextureMapPattern(uvpattern, mapper)
	// sphere1.Material().SetPattern(pattern)
	sphere1.Material().Color = tracer.ColorName(colornames.Lightgoldenrodyellow)
	pert, err := tracer.NewImageHeightmapPerturber(filename, mapper, sphere1)
	if err != nil {
		log.Fatal(err)
	}
	sphere1.Material().SetPerturber(pert)

	w.AddObject(sphere1)
	w.AddObject(floor(0))
	w.AddObject(backWall(50))
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
	w.Config.Antialias = 4

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
	w.Config.Antialias = 2

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

func movedgroup() {

	w := envxy(800, 600)

	g := tracer.NewGroup()

	g.AddMembers(glasssphere(), pedestal())
	g.SetTransform(tracer.IdentityMatrix().Translate(-2, 0, 4))

	w.AddObject(floor(0))
	w.AddObject(g)

	tracer.Render(w)

}

func groupingroup() {

	w := envxy(800, 600)

	g := tracer.NewGroup()

	g.AddMembers(glasssphere(), pedestal())
	g.SetTransform(tracer.IdentityMatrix().Translate(-2, 0, 4))

	gouter := tracer.NewGroup()
	gouter.AddMember(g)

	w.AddObject(floor(0))
	w.AddObject(gouter)

	tracer.Render(w)

}
func texturetri() {
	w := envxy(1024, 768)
	w.Config.Antialias = 2

	floor := tracer.NewPlane()
	floor.SetTransform(tracer.IdentityMatrix().Translate(0, -3, 0))
	floorp := tracer.NewRingPattern(tracer.ColorName(colornames.Red), tracer.White())
	floor.Material().SetPattern(floorp)
	w.AddObject(floor)

	backWall := tracer.NewPlane()
	backWall.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).Translate(0, 0, 40))
	backWallp := tracer.NewStripedPattern(tracer.ColorName(colornames.Blue), tracer.White())
	backWall.Material().SetPattern(backWallp)
	backWall.Material().Specular = 0.2
	w.AddObject(backWall)

	g1 := tracer.NewGroup()
	g1.SetTransform(tracer.IdentityMatrix().RotateZ(math.Pi/8).Translate(0.7, 0.4, 0))
	w.AddObject(g1)

	t1 := tracer.NewTriangle(tracer.NewPoint(0, 0, 0), tracer.NewPoint(2, 0, 0), tracer.NewPoint(1, 2, 0))
	t1.Material().Color = tracer.ColorName(colornames.Darkred)

	t2 := tracer.NewTriangle(tracer.NewPoint(0, 0, 0), tracer.NewPoint(-2, 0, 0), tracer.NewPoint(-1, 2, 0))
	t2.Material().Color = tracer.ColorName(colornames.Darkblue)
	// t2.Material().Transparency = 0.5
	// t2.Material().Diffuse = 0.1
	// t2.Material().Ambient = 0.1
	// t2.Material().ShadowCaster = false

	t3 := tracer.NewTriangle(tracer.NewPoint(0, 0, 0), tracer.NewPoint(-1, 2, 0), tracer.NewPoint(1, 2, 0))
	t3.Material().Color = tracer.ColorName(colornames.Darkgreen)
	// t3.Material().Transparency = 1
	// t3.Material().Diffuse = 0.1
	// t3.Material().Ambient = 0.1
	// t3.Material().ShadowCaster = false

	g1.AddMembers(t1, t2, t3)
	image := "images/earthmap1k.jpg"
	up, err := tracer.NewUVImagePattern(image)
	if err != nil {
		log.Fatal(err)
	}
	mapper := tracer.NewSphericalMap()
	p := tracer.NewTextureMapPattern(up, mapper)
	g1.Material().SetPattern(p)

	tracer.Render(w)
}

func antialias1() {
	w := envxy(1024, 768)
	w.Camera().SetFoV(math.Pi / 5)
	w.Config.Antialias = 3

	// s1 := tracer.NewUnitSphere()
	s1 := glasssphere()
	s1.SetTransform(tracer.IdentityMatrix().Translate(0, 1.85, 0))

	w.AddObject(s1)
	w.AddObject(backWall(10))

	tracer.Render(w)
}

func hollowsphere(wallWidth float64) *tracer.Group {
	outer := tracer.NewUnitSphere()
	outer.Material().Transparency = 0.9
	outer.Material().Reflective = 0.9
	outer.Material().ShadowCaster = false
	// outer.Material().Color = tracer.ColorName(colornames.Red)
	outer.Material().RefractiveIndex = 1.55
	outer.Material().Diffuse = 0
	outer.Material().Specular = 0.8

	inner := tracer.NewUnitSphere()
	inner.SetTransform(tracer.IdentityMatrix().Scale(1-wallWidth, 1-wallWidth, 1-wallWidth))
	inner.Material().Transparency = 0.9
	// inner.Material().Reflective = 0.9
	inner.Material().ShadowCaster = false
	// inner.Material().Color = tracer.ColorName(colornames.Black)
	inner.Material().RefractiveIndex = 1.0
	inner.Material().Diffuse = 0
	inner.Material().Specular = 0.8

	g := tracer.NewGroup()
	g.AddMembers(inner, outer)

	return g
}

func hollowsphere1() {
	w := envxy(1024, 768)
	w.Config.Antialias = 1

	// width of the sphere wall: (0, 1)
	wallWidth := 0.02
	sphere := hollowsphere(wallWidth)

	innercube := tracer.NewUnitCube()
	innercube.SetTransform(
		tracer.IdentityMatrix().Scale(0.4, 0.4, 0.4).RotateX(math.Pi / 4).RotateZ(math.Pi / 4).RotateY(math.Pi / 4))
	// icuvp := tracer.NewUVCheckersPattern(4, 4,
	// 	tracer.ColorName(colornames.Blue), tracer.ColorName(colornames.Yellow))
	// icp := tracer.NewCubeMapPatternSame(icuvp)
	// innercube.Material().SetPattern(icp)
	innercube.Material().Ambient = 0
	innercube.Material().Diffuse = 0
	innercube.Material().Reflective = 0.9

	g := tracer.NewGroup()
	g.AddMembers(sphere, innercube)
	g.SetTransform(tracer.IdentityMatrix().Scale(1.7, 1.7, 1.7).Translate(0, 1.7, 2))

	w.AddObject(g)

	w.AddObject(floor(0))
	w.AddObject(backWall(10))

	tracer.Render(w)
}

func emissive() {
	w := envxy(640, 480)
	w.Config.Antialias = 0
	w.Config.SoftShadows = false
	w.Config.SoftShadowRays = 1
	// w.Camera().SetFoV(math.Pi / 4)

	l := tracer.NewAreaLight(tracer.NewUnitSphere(),
		tracer.ColorName(colornames.White), true)
	l.SetTransform(
		tracer.IdentityMatrix().Scale(0.2, 1, 0.2).Translate(2, 1, 2))
	l.SetIntensity(l.Intensity().Scale(0.5))

	l2 := tracer.NewAreaLight(tracer.NewUnitCube(),
		tracer.ColorName(colornames.White), true)
	l2.SetTransform(
		tracer.IdentityMatrix().Scale(0.2, 1, 0.2).Translate(-2, 1, 2))
	l2.SetIntensity(l.Intensity().Scale(0.5))

	w.SetLights(tracer.Lights{l, l2})

	// g := sphereOnPedestal()
	g := mirrorSphereOnPedestal()
	g.SetTransform(tracer.IdentityMatrix().Translate(0, 0, 2.5))

	w.AddObject(g)
	w.AddObject(defaultroom())

	tracer.Render(w)
}

// returns a visible spherical area light set at x,y,z, scaled by s, of intensity i
func spherearealight(x, y, z, s float64, c color.Color) tracer.Light {

	l := tracer.NewAreaLight(tracer.NewUnitSphere(), tracer.ColorName(c), true)
	l.SetTransform(tracer.IdentityMatrix().Scale(s, s, s).Translate(x, y-s, z))
	return l
}

// returns a plane area light raised by y, scaled by xs,xy,xz, of color c
func flatarealight(x, y, z, xs, ys, zs float64, c color.Color) tracer.Light {
	l := tracer.NewAreaLight(tracer.NewUnitCube(), tracer.ColorName(c), true)
	l.SetTransform(tracer.IdentityMatrix().Scale(xs, ys, zs).Translate(x, y-2*ys, z))
	return l
}

func simpleroom() {
	w := envxy(640, 480)
	w.Config.Antialias = 0
	w.Config.SoftShadows = false
	w.Config.SoftShadowRays = 10
	w.Camera().SetFoV(math.Pi / 3)

	// w.SetLights(tracer.Lights{spherearealight(0, 4.95, 5, 0.2, colornames.White)})
	w.SetLights(tracer.Lights{flatarealight(0, 5, 5, 3, 0.1, 1, colornames.White)})

	w.AddObject(defaultroom())

	s := sphereOnPedestal()
	s.SetTransform(tracer.IdentityMatrix().Scale(1.5, 1.5, 1.5).Translate(0, 0, 3))
	w.AddObject(s)

	mirrorRight := mirrorcube(5-0.02, 2.5, 3.2, 3.3, 1.5, 0.02, 0, math.Pi/2, 0)
	w.AddObject(mirrorRight)

	mirrorLeft := mirrorcube(-5+0.02, 2.5, 3.2, 3.3, 1.5, 0.02, 0, math.Pi/2, 0)
	w.AddObject(mirrorLeft)

	tracer.Render(w)
}

func defaultroom() *tracer.Group {
	left, right := -5.0, 5.0
	front, back := -10.0, 10.0
	floor, ceiling := 0.0, 5.0
	return room(left, front, right, back, ceiling, floor)
}

// room returns a room with all walls of the provided sizes
func room(left, front, right, back, clng, flr float64) *tracer.Group {
	g := tracer.NewGroup()
	g.AddMember(floor(flr))
	g.AddMember(backWall(back))
	g.AddMember(leftWall(left))
	g.AddMember(rightWall(right))
	g.AddMember(frontWall(front))
	g.AddMember(ceiling(clng))
	return g
}

func simpletexturewall(filename string) {
	w := envxy(640, 480)

	p := tracer.NewPlane()
	p.SetTransform(tracer.IdentityMatrix().RotateX(math.Pi/2).RotateZ(math.Pi/2).Translate(0, 0, 10))
	ppuv, err := tracer.NewUVImagePattern(filename)
	if err != nil {
		log.Fatal(err)
	}
	pp := tracer.NewTextureMapPattern(ppuv, tracer.NewPlaneMap())
	pp.SetTransform(tracer.IdentityMatrix().Scale(5, 5, 5).RotateY(math.Pi / 2))
	p.Material().SetPattern(pp)
	p.Material().Specular = 0

	w.AddObject(p)
	tracer.Render(w)
}

func envxy2(width, height float64) *tracer.World {
	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5
	w.Config.SoftShadows = false

	// override light here
	w.SetLights([]tracer.Light{
		// tracer.NewPointLight(tracer.NewPoint(0, 4, 5), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(2, 10, -10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 0, -4)
	to := tracer.NewPoint(0, 0, 10)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 4)

	return w
}
func envxy(width, height float64) *tracer.World {
	// setup world, default light and camera
	w := tracer.NewDefaultWorld(width, height)
	w.Config.MaxRecusions = 5

	// override light here
	w.SetLights([]tracer.Light{
		// tracer.NewPointLight(tracer.NewPoint(0, 4, 5), tracer.NewColor(1, 1, 1)),
		tracer.NewPointLight(tracer.NewPoint(2, 10, -10), tracer.NewColor(1, 1, 1)),
	})

	// where the camera is and where it's pointing; also which way is "up"
	from := tracer.NewPoint(0, 0, -9)
	to := tracer.NewPoint(0, 0, 20)
	up := tracer.NewVector(0, 1, 0)
	cameraTransform := tracer.ViewTransform(from, to, up)
	w.Camera().SetTransform(cameraTransform)
	w.Camera().SetFoV(math.Pi / 4)

	return w
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

	// dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/images/heightmaps"))
	// heightmapplane(path.Join(dir, "volcano.gif"))
	dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/images/brickwall"))
	brickwall(dir)
	// simplesphere()
	// heightmapsphere(path.Join(dir, "brick_bump.png"))
	// simpleroom()
	// emissive()
	// simpletexturewall(path.Join(dir, "brick_bump.png"))
	// simplecone()
	// simplecylinder()
	// hollowsphere1()
	// groupingroup()
	// antialias1()
	// texturetri()
	// shapes()
	// movedgroup()
	// skyboxcube1("field1")
	// skyboxsphere1("rooitou_park_4k.hdr")
	// image1()
	// textureMap()
	// cubeMap()
	// mirrors()
	// mirror()
	// cube()
	// glass()

	// window()
	// pond()
	// spherewarp()
	// cylinder()
	// cone()
	// cylindertextures()
	// group()
	// triangle()
	// https://octolinker-demo.now.sh/mokiat/go-data-front
	// csg()

	// dir := fmt.Sprintf(path.Join(utils.Homedir(), "go/src/github.com/DanTulovsky/tracer/obj"))
	// // f := path.Join(dir, "cubes2.obj")
	// f := path.Join(dir, "monkey-smooth2.obj")
	// // f := path.Join(dir, "texture2.obj")
	// objParse(f)

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
