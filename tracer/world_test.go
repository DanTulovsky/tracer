package tracer

import (
	"log"
	"math"
	"testing"

	"github.com/DanTulovsky/tracer/constants"
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
			want: &World{
				Config: NewWorldConfig(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewWorld(NewWorldConfig()))
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
			is := tt.world.Intersections(tt.args.r, NewIntersections())
			assert.Equal(t, 4, len(is))

			for x := 0; x < len(is); x++ {
				assert.Equal(t, tt.want[x], is[x].T())
			}
		})
	}
}

func TestWorld_shadeHit(t *testing.T) {
	type args struct {
		i         *Intersection
		r         Ray
		material  *Material
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
				material:  NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200, 0, 0, 1, nil),
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

			xs := NewIntersections(tt.args.i)

			state := PrepareComputations(tt.args.i, tt.args.r, xs)
			assert.True(t, tt.want.Equal(tt.world.shadeHit(state, 1, NewIntersections())), "should equal")
		})
	}
}

func TestWorld_shadeHitShadow(t *testing.T) {
	w := NewDefaultTestWorld()
	w.SetLights([]Light{NewPointLight(NewPoint(0, 0, -10), ColorName(colornames.White))})

	s1 := NewUnitSphere()
	w.AddObject(s1)

	s2 := NewUnitSphere()
	s2.SetTransform(IdentityMatrix().Translate(0, 0, 10))
	w.AddObject(s2)

	r := NewRay(NewPoint(0, 0, 5), NewVector(0, 0, 1))
	xs := NewIntersections(
		NewIntersection(s2, 4))
	i := xs[0]

	state := PrepareComputations(i, r, xs)
	c := w.shadeHit(state, 1, NewIntersections())
	assert.Equal(t, NewColor(0.1, 0.1, 0.1), c, "should equal")
}

func TestWorld_shadeHitOffset(t *testing.T) {
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))
	shape := NewUnitSphere()
	shape.SetTransform(IdentityMatrix().Translate(0, 0, 1))

	xs := NewIntersections(
		NewIntersection(shape, 5))
	i := xs[0]

	state := PrepareComputations(i, r, xs)

	assert.Less(t, state.OverPoint.Z(), -constants.Epsilon/2)
	assert.Greater(t, state.Point.Z(), state.OverPoint.Z())
}

func TestWorld_ColorAt(t *testing.T) {
	type args struct {
		r Ray
	}
	tests := []struct {
		name   string
		world  *World
		args   args
		m1, m2 *Material
		want   Color
	}{
		{
			name:  "ray misses",
			world: NewDefaultTestWorld(),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 1, 0)),
			},
			m1:   NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200, 0, 0, 1, nil),
			m2:   NewDefaultMaterial(),
			want: ColorName(colornames.Black),
		},
		{
			name:  "ray hits",
			world: NewDefaultTestWorld(),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			m1:   NewMaterial(NewColor(0.8, 1.0, 0.6), 0.1, 0.7, 0.2, 200, 0, 0, 1, nil),
			m2:   NewDefaultMaterial(),
			want: NewColor(0.38066, 0.47583, 0.2855),
		},
		{
			name:  "color with an intersection behind the ray",
			world: NewDefaultTestWorld(),
			args: args{
				r: NewRay(NewPoint(0, 0, 0.75), NewVector(0, 0, -1)),
			},
			m1:   NewMaterial(NewColor(0.8, 1.0, 0.6), 1, 0.7, 0.2, 200, 0, 0, 1, nil),
			m2:   NewMaterial(NewColor(1.0, 1.0, 1.0), 1, 0.9, 0.9, 200, 0, 0, 1, nil),
			want: NewColor(1, 1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			tt.world.Objects[0].SetMaterial(tt.m1)
			tt.world.Objects[1].SetMaterial(tt.m2)

			assert.True(t, tt.want.Equal(tt.world.ColorAt(tt.args.r, 1, NewIntersections())), "should equal")
		})
	}
}

func TestWorld_Render(t *testing.T) {
	tests := []struct {
		name      string
		world     *World
		camera    *Camera
		canvas    *Canvas
		transform Matrix
	}{
		{
			name:      "test1",
			world:     NewDefaultTestWorld(),
			camera:    NewCamera(11, 11, math.Pi/2),
			canvas:    NewCanvas(100, 100),
			transform: ViewTransform(NewPoint(0, 0, -5), NewPoint(0, 0, 0), NewVector(0, 1, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.world.SetCamera(tt.camera)
			tt.world.Camera().SetTransform(tt.transform)
			tt.world.Render(tt.camera, tt.canvas)
			got, err := tt.canvas.Get(5, 5)
			if err != nil {
				t.Error(err)
			}
			assert.True(t, NewColor(0.38066, 0.47583, 0.2855).Equal(got), "should equal")

		})
	}
}

func TestWorld_IsShadowed(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		world *World
		args  args
		want  bool
	}{
		{
			name:  "nothing collinear with point and light; no shadow",
			world: NewDefaultTestWorld(),
			args: args{
				p: NewPoint(0, 10, 0),
			},
			want: false,
		},
		{
			name:  "object between point and light; shadow",
			world: NewDefaultTestWorld(),
			args: args{
				p: NewPoint(10, -10, 10),
			},
			want: true,
		},
		{
			name:  "object is behind light; no shadow",
			world: NewDefaultTestWorld(),
			args: args{
				p: NewPoint(-20, 20, -20),
			},
			want: false,
		},
		{
			name:  "object is behind the point; no shadow",
			world: NewDefaultTestWorld(),
			args: args{
				p: NewPoint(-2, 2, -2),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// default world has only one light
			assert.Equal(t, tt.want,
				tt.world.IsShadowed(tt.args.p, tt.world.Lights[0].Position(), NewIntersections()))
		})
	}
}

func TestWorld_ReflectedColor_NotReflective(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(Origin(), NewVector(0, 0, 1))

	shape := w.Objects[1]
	shape.Material().Ambient = 1

	xs := NewIntersections(
		NewIntersection(shape, math.Sqrt2))
	i := xs[0]

	state := PrepareComputations(i, r, xs)
	clr := w.ReflectedColor(state, 1, NewIntersections())

	assert.Equal(t, Black(), clr, "should equal")
}

func TestWorld_ReflectedColor_Reflective(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))

	shape := NewPlane()
	shape.Material().Reflective = 0.5
	shape.SetTransform(IdentityMatrix().Translate(0, -1, 0))
	w.AddObject(shape)

	xs := NewIntersections(
		NewIntersection(shape, math.Sqrt2))
	i := xs[0]

	state := PrepareComputations(i, r, xs)
	clr := w.ReflectedColor(state, 1, NewIntersections())
	expected := NewColor(0.19033, 0.23791, 0.142749)

	assert.True(t, expected.Equal(clr), "should equal")
}

func TestWorld_shadeHit_Reflective(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))

	shape := NewPlane()
	shape.Material().Reflective = 0.5
	shape.SetTransform(IdentityMatrix().Translate(0, -1, 0))
	w.AddObject(shape)

	xs := NewIntersections(
		NewIntersection(shape, math.Sqrt2))
	i := xs[0]

	state := PrepareComputations(i, r, xs)
	clr := w.shadeHit(state, 1, NewIntersections())
	expected := NewColor(0.876757, 0.924340, 0.829174)

	assert.True(t, expected.Equal(clr), "should equal")
}

func TestWorld_AvoidInfRecursion(t *testing.T) {
	w := NewDefaultTestWorld()
	w.Objects = []Shaper{}

	w.SetLights([]Light{NewPointLight(Origin(), NewColor(1, 1, 1))})

	lower := NewPlane()
	lower.Material().Reflective = 1
	lower.SetTransform(IdentityMatrix().Translate(0, -1, 0))
	w.AddObject(lower)

	upper := NewPlane()
	upper.Material().Reflective = 1
	upper.SetTransform(IdentityMatrix().Translate(0, 1, 0))
	w.AddObject(upper)

	r := NewRay(NewPoint(0, 0, 0), NewVector(0, 1, 0))

	clr := w.ColorAt(r, 4, NewIntersections())
	assert.NotNil(t, clr)
}

func TestWorld_shadeHit_MaxRecursiveReflected(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))

	shape := NewPlane()
	shape.Material().Reflective = 0.5
	shape.SetTransform(IdentityMatrix().Translate(0, -1, 0))
	w.AddObject(shape)

	xs := NewIntersections(
		NewIntersection(shape, math.Sqrt2))
	i := xs[0]

	state := PrepareComputations(i, r, xs)
	clr := w.ReflectedColor(state, 0, NewIntersections())
	expected := Black()
	log.Println(clr)

	assert.True(t, expected.Equal(clr), "should equal")
}

func TestWorld_RefractedColor_Opaque(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))

	shape := w.Objects[0]

	xs := NewIntersections(
		NewIntersection(shape, 4),
		NewIntersection(shape, 6))
	i := xs[0]

	state := PrepareComputations(i, r, xs)
	clr := w.RefractedColor(state, 5, NewIntersections())

	assert.Equal(t, Black(), clr, "should equal")
}

func TestWorld_efractedColor_MaxRecursion(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))

	shape := w.Objects[0]
	shape.Material().Transparency = 1
	shape.Material().RefractiveIndex = 1.5

	xs := NewIntersections(
		NewIntersection(shape, 4),
		NewIntersection(shape, 6))
	i := xs[0]

	state := PrepareComputations(i, r, xs)
	clr := w.RefractedColor(state, 0, NewIntersections())

	assert.Equal(t, Black(), clr, "should equal")
}

func TestWorldRefractedColor_TotalInternalReflection(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(NewPoint(0, 0, math.Sqrt2/2), NewVector(0, 1, 0))

	shape := w.Objects[0]
	shape.Material().Transparency = 1
	shape.Material().RefractiveIndex = 1.5

	xs := NewIntersections(
		NewIntersection(shape, -math.Sqrt2/2),
		NewIntersection(shape, math.Sqrt2/2))
	i := xs[1] // inside the sphere

	state := PrepareComputations(i, r, xs)
	clr := w.RefractedColor(state, 5, NewIntersections())

	assert.Equal(t, Black(), clr, "should equal")
}

func TestWorldRefractedColor_Normal(t *testing.T) {
	w := NewDefaultTestWorld()
	r := NewRay(NewPoint(0, 0, 0.1), NewVector(0, 1, 0))

	a := w.Objects[0]
	a.Material().Ambient = 1
	a.Material().SetPattern(newTestPattern())

	b := w.Objects[1]
	b.Material().Transparency = 1.0
	b.Material().RefractiveIndex = 1.5

	xs := NewIntersections(
		NewIntersection(a, -0.9899),
		NewIntersection(b, -0.4899),
		NewIntersection(b, 0.4899),
		NewIntersection(a, 0.9899))

	state := PrepareComputations(xs[2], r, xs)
	clr := w.RefractedColor(state, 5, NewIntersections())

	assert.True(t, NewColor(0, 0.998874, 0.047218).Equal(clr), "should be true")
}

func TestWorld_shadeHit_Transparent(t *testing.T) {
	w := NewDefaultTestWorld()

	floor := NewPlane()
	floor.SetTransform(IdentityMatrix().Translate(0, -1, 0))
	floor.Material().Transparency = 0.5
	floor.Material().RefractiveIndex = 1.5
	w.AddObject(floor)

	ball := NewUnitSphere()
	ball.Material().Color = NewColor(1, 0, 0)
	ball.Material().Ambient = 0.5
	ball.SetTransform(IdentityMatrix().Translate(0, -3.5, -0.5))
	w.AddObject(ball)

	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	xs := NewIntersections(NewIntersection(floor, math.Sqrt2))

	state := PrepareComputations(xs[0], r, xs)
	clr := w.shadeHit(state, 5, NewIntersections())

	assert.True(t, NewColor(0.936425, 0.686425, 0.686425).Equal(clr), "should be true")
}

func TestWorld_shadeHit_Schlick(t *testing.T) {
	w := NewDefaultTestWorld()

	floor := NewPlane()
	floor.SetTransform(IdentityMatrix().Translate(0, -1, 0))
	floor.Material().Transparency = 0.5
	floor.Material().Reflective = 0.5
	floor.Material().RefractiveIndex = 1.5
	w.AddObject(floor)

	ball := NewUnitSphere()
	ball.Material().Color = NewColor(1, 0, 0)
	ball.Material().Ambient = 0.5
	ball.SetTransform(IdentityMatrix().Translate(0, -3.5, -0.5))
	w.AddObject(ball)

	r := NewRay(NewPoint(0, 0, -3), NewVector(0, -math.Sqrt2/2, math.Sqrt2/2))
	xs := NewIntersections(NewIntersection(floor, math.Sqrt2))

	state := PrepareComputations(xs[0], r, xs)
	clr := w.shadeHit(state, 5, NewIntersections())

	assert.True(t, NewColor(0.93391, 0.69643, 0.69243).Equal(clr), "should be true")

}
