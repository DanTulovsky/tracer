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
