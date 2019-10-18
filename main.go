package main

import (
	"fmt"
	"log"
	"math"
	"os"

	"golang.org/x/image/colornames"

	"github.com/DanTulovsky/tracer/tracer"
)

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

	c.Set(x, y, colornames.Red)

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
	c.SetFloat(center.X(), center.Z(), colornames.Yellow)

	radius := 7.0 / 8.0 * center.X()
	twelve := tracer.NewPoint(0, 0, 1)

	for hour := 1.0; hour <= 12; hour++ {
		m := tracer.IdentityMatrix().RotateY(hour*(math.Pi/6.0)).Scale(radius, 1, radius).Translate(center.X(), center.Y(), center.Z())
		p := twelve.TimesMatrix(m)
		c.SetFloat(p.X(), p.Z(), colornames.Red)
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
	// first circled drawn by a ray
}

func main() {

	// testCanvas()
	// test1()
	// clock()
	circle()

}
