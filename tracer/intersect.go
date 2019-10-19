package tracer

import (
	"fmt"
	"sort"
)

// Intersection encapsulates an intersection t value an an object
type Intersection struct {
	o Object
	t float64
}

// NewIntersection returns an intersection object
func NewIntersection(o Object, t float64) Intersection {
	return Intersection{o: o, t: t}
}

// T returns the t value for the intersection
func (i Intersection) T() float64 {
	return i.t
}

// Object returns the object of the intersection
func (i Intersection) Object() Object {
	return i.o
}

// Intersections is a collection of Intersections
type Intersections []Intersection

// NewIntersections aggregates the given intersections into a sorted list
func NewIntersections(i ...Intersection) Intersections {
	is := Intersections{}

	for _, int := range i {
		is = append(is, int)
	}

	sort.Sort(byT(is))
	return is
}

// Hit returns the visible intersection (lowest non-negative value)
func (i Intersections) Hit() (Intersection, error) {

	sort.Sort(byT(i))

	for _, int := range i {
		if int.t >= 0 {
			return int, nil
		}
	}

	return Intersection{}, fmt.Errorf("no intersections")
}

// byT sorts Intersections by the t value
type byT Intersections

func (a byT) Len() int           { return len(a) }
func (a byT) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byT) Less(i, j int) bool { return a[i].t < a[j].t }

// IntersectionState holds precomputed values for an intersection
type IntersectionState struct {
	T       float64 // How far away from Ray origin did this occur?
	Object  Object  // The object we intersected
	Point   Point
	EyeV    Vector
	NormalV Vector
	Inside  bool // did the hit occure inside or outside the shape?
}

// PrepareComputations prepopulates the IntersectionState structure
func PrepareComputations(i Intersection, r Ray) IntersectionState {
	point := r.Position(i.T())
	object := i.Object()
	normalv := object.NormalAt(point)
	eyev := r.Dir.Negate()
	inside := false

	// check if interscection happened inside the shape or outside of it
	if normalv.Dot(eyev) < 0 {
		inside = true
		normalv = normalv.Negate()
	}

	return IntersectionState{
		T:       i.T(),
		Object:  object,
		Point:   point,
		EyeV:    eyev,
		NormalV: normalv,
		Inside:  inside,
	}
}
