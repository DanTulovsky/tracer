package tracer

import (
	"sort"

	"golang.org/x/exp/shiny/materialdesign/colornames"
)

// World holds everything in it
type World struct {
	Objects []Object
	Lights  []Light
}

// NewWorld returns a new empty world
func NewWorld() *World {
	return &World{}
}

// NewDefaultTestWorld returns a world that many tests expect
func NewDefaultTestWorld() *World {
	l1 := NewPointLight(NewPoint(-10, 10, -10), ColorName(colornames.White))

	s1 := NewUnitSphere()
	s1.SetMaterial(NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200))

	s2 := NewUnitSphere()
	s2.SetTransform(IdentityMatrix().Scale(0.5, 0.5, 0.5))

	return &World{
		Objects: []Object{s1, s2},
		Lights:  []Light{l1},
	}
}

// Intersections returns all the intersections in the world with the given ray
func (w *World) Intersections(r Ray) Intersections {
	var is Intersections

	for _, o := range w.Objects {
		iso := o.IntersectWith(r)
		is = append(is, iso...)
	}

	sort.Sort(byT(is))

	return is
}
