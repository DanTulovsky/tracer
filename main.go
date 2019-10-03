package main

import (
	"fmt"

	"github.com/DanTulovsky/tracer/tracer"
)

func main() {
	p := tracer.NewPoint(0, 0, 0)
	v := tracer.NewVector(1, 1, 1)

	fmt.Printf("%v == %v? %v\n", p, v, p.Equals(v))
}
