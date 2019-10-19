package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWorld(t *testing.T) {
	tests := []struct {
		name string
		want *World
	}{
		{
			name: "empty world",
			want: &World{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewWorld())
		})
	}
}

func TestNewDefaultTestWorld(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "default1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewDefaultTestWorld()
			assert.Equal(t, 2, len(w.Objects), "should equal")
			assert.Equal(t, 1, len(w.Lights), "should equal")
		})
	}
}

func TestWorld_Intersections(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name  string
		world *World
		args  args
		want  []float64 // t values of intersections
	}{
		{
			name:  "test1",
			world: NewDefaultTestWorld(),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			want: []float64{4, 4.5, 5.5, 6},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := tt.world.Intersections(tt.args.r)
			assert.Equal(t, 4, len(is))

			for x := 0; x < len(is); x++ {
				assert.Equal(t, tt.want[x], is[x].T())
			}
		})
	}
}
