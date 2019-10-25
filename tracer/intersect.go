package tracer

import (
	"fmt"
	"sort"

	"github.com/DanTulovsky/tracer/constants"
)

// Intersection encapsulates an intersection t value an an object
type Intersection struct {
	o Shaper
	t float64
}

// NewIntersection returns an intersection object
// TODO: Return as pointer
func NewIntersection(o Shaper, t float64) Intersection {
	return Intersection{o: o, t: t}
}

// T returns the t value for the intersection
func (i Intersection) T() float64 {
	return i.t
}

// Object returns the object of the intersection
func (i Intersection) Object() Shaper {
	return i.o
}

// Intersections is a collection of Intersections
// TODO: Return as list of pointers
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
	T         float64 // How far away from Ray origin did this occur?
	Object    Shaper  // The object we intersected
	Point     Point   // the point of intersection
	EyeV      Vector  // eye vector
	NormalV   Vector  // normal vector
	Inside    bool    // did the hit occure inside or outside the shape?
	OverPoint Point   // offset to properly render shadows due to floating point errors
	ReflectV  Vector  // reflection vector
	N1, N2    float64 // RefractiveIndex of (n1) leaving material and (n2) entering material
}

func objectInList(o Shaper, list []Shaper) bool {
	for _, obj := range list {
		if o == obj {
			return true
		}
	}
	return false
}

// findObjectInList returns the index of o in list, or error if it not in list
func findObjectInList(o Shaper, list []Shaper) (int, error) {
	for i, n := range list {
		if o == n {
			return i, nil
		}
	}
	return len(list), fmt.Errorf("%v not found in list", o)
}

// delObjectFromList deletes the object from the list
func delObjectFromList(o Shaper, list []Shaper) []Shaper {
	for i, obj := range list {
		if o == obj {
			copy(list[i:], list[i+1:]) // Shift a[i+1:] left one index.
			list[len(list)-1] = nil    // Erase last element (write zero value).
			list = list[:len(list)-1]  // Truncate slice.
		}
	}
	return list
}

// findRefractiveIndexes returns the refractive indexes of leaving material and entering material
func findRefractiveIndexes(hit Intersection, xs Intersections) (n1, n2 float64) {
	var containers []Shaper

	for _, i := range xs {
		if i == hit {
			if len(containers) == 0 {
				n1 = 1.0
			} else {
				n1 = containers[len(containers)-1].Material().RefractiveIndex
			}
		}

		if objectInList(i.Object(), containers) {
			containers = delObjectFromList(i.Object(), containers)
		} else {
			containers = append(containers, i.Object())
		}

		if i == hit {
			if len(containers) == 0 {
				n2 = 1.0
			} else {
				n2 = containers[len(containers)-1].Material().RefractiveIndex
			}
			return n1, n2
		}
	}

	return n1, n2
}

// PrepareComputations prepopulates the IntersectionState structure
func PrepareComputations(i Intersection, r Ray, xs Intersections) *IntersectionState {
	var n1, n2 float64
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

	overPoint := point.AddVector(normalv.Scale(constants.Epsilon))
	reflectv := r.Dir.Reflect(normalv)
	n1, n2 = findRefractiveIndexes(i, xs)

	return &IntersectionState{
		T:         i.T(),
		Object:    object,
		Point:     point,
		EyeV:      eyev,
		NormalV:   normalv,
		Inside:    inside,
		OverPoint: overPoint,
		ReflectV:  reflectv,
		N1:        n1,
		N2:        n2,
	}
}
