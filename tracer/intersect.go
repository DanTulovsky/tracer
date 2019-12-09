package tracer

import (
	"fmt"
	"math"
	"sort"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/google/go-cmp/cmp"
)

// Intersection encapsulates an intersection t value an an object
type Intersection struct {
	o Shaper
	t float64
	// Intersection point on the shape (used only for triangles)
	u, v float64
}

// NewIntersection returns an intersection object
func NewIntersection(o Shaper, t float64) *Intersection {
	return &Intersection{o: o, t: t}
}

// T returns the t value for the intersection
func (i *Intersection) T() float64 {
	return i.t
}

// Object returns the object of the intersection
func (i *Intersection) Object() Shaper {
	return i.o
}

// NewIntersectionUV returns an intersection object with UV filled in
func NewIntersectionUV(o Shaper, t, u, v float64) *Intersection {
	return &Intersection{o: o, t: t, u: u, v: v}
}

// Equal returns true if the intersections are the same
func (i *Intersection) Equal(i2 *Intersection) bool {
	return i.t == i2.t &&
		i.u == i2.u &&
		i.v == i2.v &&
		cmp.Equal(i.Object(), i2.Object())
}

// Intersections is a collection of Intersections
type Intersections []*Intersection

// NewIntersections aggregates the given intersections into a sorted list
func NewIntersections(i ...*Intersection) Intersections {
	is := make(Intersections, 0, 4)

	for _, int := range i {
		is = append(is, int)
	}

	if len(is) > 1 {
		sort.Sort(byT(is))
	}
	return is
}

// Hit returns the visible intersection (lowest non-negative value)
func (i Intersections) Hit() (*Intersection, error) {

	sort.Sort(byT(i))

	for _, xs := range i {
		if xs.t >= 0 {
			return xs, nil
		}
	}

	return &Intersection{}, fmt.Errorf("no intersections")
}

// byT sorts Intersections by the t value
type byT Intersections

func (a byT) Len() int           { return len(a) }
func (a byT) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byT) Less(i, j int) bool { return a[i].t < a[j].t }

// IntersectionState holds precomputed values for an intersection
type IntersectionState struct {
	T                     float64 // How far away from Ray origin did this occur?
	Object                Shaper  // The object we intersected
	Point                 Point   // the point of intersection
	EyeV                  Vector  // eye vector
	NormalV               Vector  // normal vector
	Inside                bool    // did the hit occure inside or outside the shape?
	OverPoint, UnderPoint Point   // offset to properly render shadows and refraction due to floating point errors
	ReflectV              Vector  // reflection vector
	N1, N2                float64 // RefractiveIndex of (n1) leaving material and (n2) entering material
	U, V                  float64 // u,v values for where the intersection occured
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
func findRefractiveIndexes(hit *Intersection, xs Intersections) (n1, n2 float64) {
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
func PrepareComputations(hit *Intersection, r Ray, xs Intersections) *IntersectionState {
	var n1, n2 float64
	point := r.Position(hit.T())
	object := hit.Object()

	normalv := object.NormalAt(point, hit)
	eyev := r.Dir.Negate()
	inside := false

	// check if interscection happened inside the shape or outside of it
	if normalv.Dot(eyev) < 0 {
		inside = true
		normalv = normalv.Negate()
	}

	normalvScaled := normalv.Scale(constants.Epsilon)
	overPoint := point.AddVector(normalvScaled)
	underPoint := point.SubVector(normalvScaled)
	reflectv := r.Dir.Reflect(normalv)
	n1, n2 = findRefractiveIndexes(hit, xs)

	return &IntersectionState{
		T:          hit.T(),
		Object:     object,
		Point:      point,
		EyeV:       eyev,
		NormalV:    normalv,
		Inside:     inside,
		OverPoint:  overPoint,
		UnderPoint: underPoint,
		ReflectV:   reflectv,
		N1:         n1,
		N2:         n2,
		U:          hit.u,
		V:          hit.v,
	}
}

// Schlick returns the reflectance - the fraction of light that is reflected [0,1]
func Schlick(s *IntersectionState) float64 {
	// cos of the angle between the eye and normal vectors
	cos := s.EyeV.Dot(s.NormalV)

	// total intrnal reflection can only occur if n1 > n2
	if s.N1 > s.N2 {
		n := s.N1 / s.N2
		sin2t := n * n * (1.0 - cos*cos)
		if sin2t > 1.0 {
			return 1.0
		}

		// compute cosine of theta_t using trig identity
		cost := math.Sqrt(1.0 - sin2t)

		// when n1 > n2, use cos(theta_t) instead
		cos = cost
	}

	r0 := ((s.N1 - s.N2) / (s.N1 + s.N2)) * ((s.N1 - s.N2) / (s.N1 + s.N2))

	c := (1 - cos) * (1 - cos) * (1 - cos) * (1 - cos) * (1 - cos)
	return r0 + (1-r0)*c
}
