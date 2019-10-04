package main

import (
	"fmt"

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

func main() {

	ticks := 0
	vScale := 1.0

	fmt.Println(tracer.NewVector(1, 1, 0).Normalize())
	fmt.Println(tracer.NewVector(1, 1, 0).Magnitude())

	p := projectile{Position: tracer.NewPoint(0, 1, 0), Velocity: tracer.NewVector(1, 0, 0).Normalize().Scale(vScale)}
	e := environment{Gravity: tracer.NewVector(0, -0.1, 0), Wind: tracer.NewVector(-0.01, 0, 0)}

	fmt.Printf("position: %2f\n", p.Position)
	for p.Position.Y() > 0 {
		p = tick(e, p)
		fmt.Printf("position: %2f\n", p.Position)
		ticks++
	}
	fmt.Printf("Total Ticks: %v\n", ticks)

}
