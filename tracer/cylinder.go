package tracer

import "github.com/google/uuid"

// Cylinder implements a cylinder of radius 1
type Cylinder struct {
	Radius float64
	Shape
}

// NewCylinder returns a new cylinder
func NewCylinder() *Cylinder {
	return &Cylinder{
		Radius: 1.0,
		Shape: Shape{
			transform: IdentityMatrix(),
			material:  NewDefaultMaterial(),
			shape:     "cylinder",
			name:      uuid.New().String(),
		},
	}
}
