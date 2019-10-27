package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCylinder(t *testing.T) {
	tests := []struct {
		name string
		want *Cylinder
	}{
		{
			name: "test1",
			want: &Cylinder{
				Radius: 1.0,
				Shape: Shape{
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "cylinder",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCylinder()
			tt.want.Shape.name = got.name // random uuid

			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}
