package tracer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
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

func TestWorld_shadeHit(t *testing.T) {
	type args struct {
		i         Intersection
		r         Ray
		material  Material
		transform Matrix
		lights    []Light
	}
	tests := []struct {
		name  string
		world *World
		args  args
		want  Color
	}{
		{
			name:  "test1",
			world: NewDefaultTestWorld(),
			args: args{
				i:         NewIntersection(NewUnitSphere(), 4),
				r:         NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
				material:  NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200),
				transform: IdentityMatrix(),
				lights: []Light{
					NewPointLight(NewPoint(-10, 10, -10), ColorName(colornames.White)),
				},
			},
			want: NewColor(0.38066, 0.47583, 0.2855),
		},
		{
			name:  "inside",
			world: NewDefaultTestWorld(),
			args: args{
				i:         NewIntersection(NewUnitSphere(), 0.5),
				r:         NewRay(NewPoint(0, 0, 0), NewVector(0, 0, 1)),
				material:  NewDefaultMaterial(),
				transform: IdentityMatrix().Scale(0.5, 0.5, 0.5),
				lights: []Light{
					NewPointLight(NewPoint(0, 0.25, 0), ColorName(colornames.White)),
				},
			},
			want: NewColor(0.90498, 0.90498, 0.90498),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.i.Object().SetMaterial(tt.args.material)
			tt.args.i.Object().SetTransform(tt.args.transform)
			tt.world.SetLights(tt.args.lights)

			state := PrepareComputations(tt.args.i, tt.args.r)
			assert.True(t, tt.want.Equal(tt.world.shadeHit(state)), "should equal")
		})
	}
}

func TestWorld_ColorAt(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name   string
		world  *World
		args   args
		m1, m2 Material
		want   Color
	}{
		{
			name:  "ray misses",
			world: NewDefaultTestWorld(),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 1, 0)),
			},
			m1:   NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200),
			m2:   NewDefaultMaterial(),
			want: ColorName(colornames.Black),
		},
		{
			name:  "ray hits",
			world: NewDefaultTestWorld(),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			m1:   NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200),
			m2:   NewDefaultMaterial(),
			want: NewColor(0.38066, 0.47583, 0.2855),
		},
		{
			name:  "color with an intersection behind the ray",
			world: NewDefaultTestWorld(),
			args: args{
				r: NewRay(NewPoint(0, 0, 0.75), NewVector(0, 0, -1)),
			},
			m1:   NewMaterial(NewColor(0.8, 1.0, 0.6), 1, 0.7, 0.2, 200),
			m2:   NewMaterial(NewColor(1.0, 1.0, 1.0), 1, 0.9, 0.9, 200),
			want: NewColor(1, 1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.world.Objects[0].SetMaterial(tt.m1)
			tt.world.Objects[1].SetMaterial(tt.m2)

			assert.True(t, tt.want.Equal(tt.world.ColorAt(tt.args.r)), "should equal")
		})
	}
}

func TestWorld_Render(t *testing.T) {
	tests := []struct {
		name      string
		world     *World
		camera    *Camera
		transform Matrix
	}{
		{
			name:      "test1",
			world:     NewDefaultTestWorld(),
			camera:    NewCamera(11, 11, math.Pi/2),
			transform: ViewTransform(NewPoint(0, 0, -5), NewPoint(0, 0, 0), NewVector(0, 1, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.world.SetCamera(tt.camera)
			tt.world.Camera().SetTransform(tt.transform)
			canvas := tt.world.Render()
			got, err := canvas.Get(5, 5)
			if err != nil {
				t.Error(err)
			}
			assert.True(t, NewColor(0.38066, 0.47583, 0.2855).Equal(got), "should equal")

		})
	}
}
