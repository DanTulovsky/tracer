package tracer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGroup(t *testing.T) {
	tests := []struct {
		name string
		want *Group
	}{
		{
			name: "test1",
			want: &Group{
				members: []Shaper{},
				Shape: Shape{
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "group",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGroup()
			assert.Equal(t, tt.want, g, "should equal")
			assert.Equal(t, IdentityMatrix(), g.Transform(), "should equal")
		})
	}
}

func TestGroup_AddMember(t *testing.T) {
	type args struct {
		m Shaper
	}
	tests := []struct {
		name  string
		group *Group
		args  args
		want  int // number of members
	}{
		{
			name:  "test1",
			group: NewGroup(),
			args: args{
				m: NewUnitSphere(),
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.group.AddMember(tt.args.m)
			assert.Equal(t, tt.want, len(tt.group.members))
		})
	}
}

func TestGroup_Members(t *testing.T) {
	tests := []struct {
		name    string
		group   *Group
		members []Shaper
		want    []Shaper
	}{
		{
			name:    "test1",
			group:   NewGroup(),
			members: []Shaper{NewUnitCube(), NewUnitSphere()},
		},
		{
			name:    "test2",
			group:   NewGroup(),
			members: []Shaper{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, m := range tt.members {
				tt.group.AddMember(m)
			}

			assert.Equal(t, len(tt.members), len(tt.group.Members()))

			for i, m := range tt.members {
				assert.Equal(t, m, tt.group.Members()[i])
			}
		})
	}
}

func TestGroup_Includes(t *testing.T) {
	type args struct {
		s Shaper
	}
	tests := []struct {
		name    string
		args    args
		group   *Group
		members []Shaper
	}{
		{
			name:    "test1",
			group:   NewGroup(),
			members: []Shaper{NewUnitCube(), NewUnitSphere()},
		},
		{
			name:    "test2",
			group:   NewGroup(),
			members: []Shaper{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, m := range tt.members {
				tt.group.AddMember(m)
			}

			for _, m := range tt.members {
				assert.True(t, tt.group.Includes(m), "should include")
				assert.Equal(t, tt.group, m.Parent(), "should equal")
			}
		})
	}
}

func TestGroup_IntersectWith(t *testing.T) {
	type args struct {
		r Ray
	}
	type member struct {
		s Shaper
		t Matrix
	}

	tests := []struct {
		name        string
		group       *Group
		members     []member
		args        args
		transform   Matrix
		wantXS      int
		wantShapers []int // index of the shapers we are expecting
	}{
		{
			name:      "empty group",
			group:     NewGroup(),
			members:   []member{},
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(Origin(), NewVector(0, 0, 1)),
			},
			wantXS: 0,
		},
		{
			name:  "spheres",
			group: NewGroup(),
			members: []member{
				{s: NewUnitSphere(), t: IdentityMatrix()},
				{s: NewUnitSphere(), t: IdentityMatrix().Translate(0, 0, -3)},
				{s: NewUnitSphere(), t: IdentityMatrix().Translate(5, 0, 0)},
			},
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1)),
			},
			wantXS:      4,
			wantShapers: []int{1, 1, 0, 0},
		},
		{
			name:  "group and object transform",
			group: NewGroup(),
			members: []member{
				{s: NewUnitSphere(), t: IdentityMatrix().Translate(5, 0, 0)},
			},
			transform: IdentityMatrix().Scale(2, 2, 2),
			args: args{
				r: NewRay(NewPoint(10, 0, -10), NewVector(0, 0, 1)),
			},
			wantXS:      2,
			wantShapers: []int{0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.group.SetTransform(tt.transform)

			// populate members
			for _, m := range tt.members {
				shape := m.s
				shape.SetTransform(m.t)
				tt.group.AddMember(shape)
			}
			tt.group.PrecomputeValues()

			got := tt.group.IntersectWith(tt.args.r)

			assert.Equal(t, tt.wantXS, len(got), "should be equal")

			for i, val := range tt.wantShapers {
				assert.Equal(t, tt.members[val].s, got[i].Object(), "should equal")
			}
		})
	}
}

func TestGroup_NormalAt(t *testing.T) {
	g1 := NewGroup()
	g1.SetTransform(IdentityMatrix().RotateY(math.Pi / 2))

	g2 := NewGroup()
	g2.SetTransform(IdentityMatrix().Scale(1, 2, 3))

	g1.AddMember(g2)

	s := NewUnitSphere()
	s.SetTransform(IdentityMatrix().Translate(5, 0, 0))
	g2.AddMember(s)

	point := NewPoint(1.7321, 1.1547, -5.5774)
	want := NewVector(0.285703, 0.428543, -0.8571605)
	got := s.NormalAt(point)

	assert.True(t, want.Equals(got), "should be true")

}

func TestGroup_Bounds(t *testing.T) {
	type member struct {
		s Shaper
		t Matrix
	}
	tests := []struct {
		name      string
		group     *Group
		members   []member
		transform Matrix
		want      Bound
	}{
		{
			name:      "test1",
			group:     NewGroup(),
			members:   []member{},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(0, 0, 0),
				Max: NewPoint(0, 0, 0),
			},
		},
		{
			name:  "single cube",
			group: NewGroup(),
			members: []member{
				{s: NewUnitCube(), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-1, -1, -1),
				Max: NewPoint(1, 1, 1),
			},
		},
		{
			name:  "single cube moved",
			group: NewGroup(),
			members: []member{
				{s: NewUnitCube(), t: IdentityMatrix().Translate(1, 0, 0)},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(0, -1, -1),
				Max: NewPoint(2, 1, 1),
			},
		},
		{
			name:  "double cube moved",
			group: NewGroup(),
			members: []member{
				{s: NewUnitCube(), t: IdentityMatrix().Translate(1, 0, 0)},
				{s: NewUnitCube(), t: IdentityMatrix().Translate(-1, 0, 0)},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-2, -1, -1),
				Max: NewPoint(2, 1, 1),
			},
		},
		{
			name:  "single cube moved and scaled",
			group: NewGroup(),
			members: []member{
				{s: NewUnitCube(), t: IdentityMatrix().Scale(2, 2, 2).Translate(1, 1, 1)},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-1, -1, -1),
				Max: NewPoint(3, 3, 3),
			},
		},
		{
			name:  "cube and sphere",
			group: NewGroup(),
			members: []member{
				{s: NewUnitCube(), t: IdentityMatrix().Scale(2, 2, 2).Translate(1, 1, 1)},
				{s: NewUnitSphere(), t: IdentityMatrix().Scale(2, 2, 2).Translate(-1, -1, -1)},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-3, -3, -3),
				Max: NewPoint(3, 3, 3),
			},
		},
		{
			name:  "plane",
			group: NewGroup(),
			members: []member{
				{s: NewPlane(), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-math.MaxFloat64, -0.001, -math.MaxFloat64),
				Max: NewPoint(math.MaxFloat64, 0.001, math.MaxFloat64),
			},
		},
		{
			name:  "cube , sphere, plane",
			group: NewGroup(),
			members: []member{
				{s: NewUnitCube(), t: IdentityMatrix().Scale(2, 2, 2).Translate(1, 1, 1)},
				{s: NewUnitSphere(), t: IdentityMatrix().Scale(2, 2, 2).Translate(-1, -1, -1)},
				{s: NewPlane(), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-math.MaxFloat64, -3, -math.MaxFloat64),
				Max: NewPoint(math.MaxFloat64, 3, math.MaxFloat64),
			},
		},
		{
			name:  "cube , sphere, plane, cone",
			group: NewGroup(),
			members: []member{
				{s: NewUnitCube(), t: IdentityMatrix().Scale(2, 2, 2).Translate(1, 1, 1)},
				{s: NewUnitSphere(), t: IdentityMatrix().Scale(2, 2, 2).Translate(-1, -1, -1)},
				{s: NewPlane(), t: IdentityMatrix()},
				{s: NewDefaultCone(), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-math.MaxFloat64, -math.MaxFloat64, -math.MaxFloat64),
				Max: NewPoint(math.MaxFloat64, math.MaxFloat64, math.MaxFloat64),
			},
		},
		{
			name:  "cylinder",
			group: NewGroup(),
			members: []member{
				{s: NewDefaultCylinder(), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-1, -math.MaxFloat64, -1),
				Max: NewPoint(1, math.MaxFloat64, 1),
			},
		},
		{
			name:  "capped cylinder",
			group: NewGroup(),
			members: []member{
				{s: NewClosedCylinder(5, 10), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-1, 5, -1),
				Max: NewPoint(1, 10, 1),
			},
		},
		{
			name:  "capped cylinder + cone",
			group: NewGroup(),
			members: []member{
				{s: NewClosedCylinder(5, 10), t: IdentityMatrix()},
				{s: NewClosedCone(-5, 2), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-5, -5, -5),
				Max: NewPoint(5, 10, 5),
			},
		},
		{
			name:  "cone",
			group: NewGroup(),
			members: []member{
				{s: NewClosedCone(-5, 2), t: IdentityMatrix()},
			},
			transform: IdentityMatrix(),
			want: Bound{
				Min: NewPoint(-5, -5, -5),
				Max: NewPoint(5, 2, 5),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.group.SetTransform(tt.transform)

			// populate members
			for _, m := range tt.members {
				shape := m.s
				shape.SetTransform(m.t)
				tt.group.AddMember(shape)
			}
			tt.group.PrecomputeValues()

			got := tt.group.Bounds()
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestGroup_boundingBoxFromPoints(t *testing.T) {
	type args struct {
		points []Point
	}
	tests := []struct {
		name string
		g    *Group
		args args
		want Bound
	}{
		{
			name: "test1",
			g:    NewGroup(), // unused
			args: args{
				points: []Point{
					NewPoint(-1, -1, -1),
					NewPoint(1, 1, 1),
				},
			},
			want: Bound{
				Min: NewPoint(-1, -1, -1),
				Max: NewPoint(1, 1, 1),
			},
		},
		{
			name: "test2",
			g:    NewGroup(), // unused
			args: args{
				points: []Point{
					NewPoint(-1, -1, -1),
					NewPoint(1, 1, 1),
					NewPoint(2, 2, 2),
				},
			},
			want: Bound{
				Min: NewPoint(-1, -1, -1),
				Max: NewPoint(2, 2, 2),
			},
		},
		{
			name: "test3",
			g:    NewGroup(), // unused
			args: args{
				points: []Point{
					NewPoint(-1, -1, -1),
					NewPoint(1, 1, 1),
					NewPoint(2, 2, 2),
					NewPoint(-2, 2, 14),
				},
			},
			want: Bound{
				Min: NewPoint(-2, -1, -1),
				Max: NewPoint(2, 2, 14),
			},
		},
		{
			name: "test4",
			g:    NewGroup(), // unused
			args: args{
				points: []Point{
					NewPoint(-1, -4, -7),
					NewPoint(2, 6, 9),
					NewPoint(3, 0, 2),
					NewPoint(-2, 2, 14),
				},
			},
			want: Bound{
				Min: NewPoint(-2, -4, -7),
				Max: NewPoint(3, 6, 14),
			},
		},
		{
			name: "inf test1",
			g:    NewGroup(), // unused
			args: args{
				points: []Point{
					NewPoint(-1, -4, -7),
					NewPoint(2, 6, 9),
					NewPoint(3, 0, 2),
					NewPoint(-2, 2, 14),
					NewPoint(math.Inf(-1), 3, math.Inf(1)),
				},
			},
			want: Bound{
				Min: NewPoint(math.Inf(-1), -4, -7),
				Max: NewPoint(3, 6, math.Inf(1)),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.g.boundingBoxFromPoints(tt.args.points...))
		})
	}
}

func TestGroup_boundBoxFromBoundingBoxes(t *testing.T) {
	type args struct {
		boxes []Bound
	}
	tests := []struct {
		name string
		g    *Group // not used
		args args
		want Bound
	}{
		{
			name: "zero test",
			g:    NewGroup(),
			args: args{
				boxes: []Bound{},
			},
			want: Bound{
				Min: Origin(),
				Max: Origin(),
			},
		},
		{
			name: "test1",
			g:    NewGroup(),
			args: args{
				boxes: []Bound{
					Bound{Min: NewPoint(-1, -1, -1), Max: NewPoint(1, 1, 1)},
				},
			},
			want: Bound{
				Min: NewPoint(-1, -1, -1),
				Max: NewPoint(1, 1, 1),
			},
		},
		{
			name: "test2",
			g:    NewGroup(),
			args: args{
				boxes: []Bound{
					Bound{Min: NewPoint(-1, -2, -3), Max: NewPoint(4, 5, 6)},
				},
			},
			want: Bound{
				Min: NewPoint(-1, -2, -3),
				Max: NewPoint(4, 5, 6),
			},
		},
		{
			name: "test3",
			g:    NewGroup(),
			args: args{
				boxes: []Bound{
					Bound{Min: NewPoint(-1, -2, -3), Max: NewPoint(4, 5, 6)},
					Bound{Min: NewPoint(-10, -2, -3), Max: NewPoint(43, 50, 6)},
				},
			},
			want: Bound{
				Min: NewPoint(-10, -2, -3),
				Max: NewPoint(43, 50, 6),
			},
		},
		{
			name: "test4",
			g:    NewGroup(),
			args: args{
				boxes: []Bound{
					Bound{Min: NewPoint(-1, -2, -3), Max: NewPoint(4, 5, 6)},
					Bound{Min: NewPoint(-10, -2, -3), Max: NewPoint(43, 50, 6)},
					Bound{Min: NewPoint(-10, math.Inf(-1), -3), Max: NewPoint(43, math.Inf(1), 6)},
				},
			},
			want: Bound{
				Min: NewPoint(-10, math.Inf(-1), -3),
				Max: NewPoint(43, math.Inf(1), 6),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.g.boundBoxFromBoundingBoxes(tt.args.boxes)
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestGroup_HasMembers(t *testing.T) {
	tests := []struct {
		name    string
		group   *Group
		members []Shaper
		want    bool
	}{
		{
			name:    "has members",
			group:   NewGroup(),
			members: []Shaper{NewUnitSphere()},
			want:    true,
		},
		{
			name:    "no members",
			group:   NewGroup(),
			members: []Shaper{},
			want:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// populate members
			for _, m := range tt.members {
				tt.group.AddMember(m)
			}

			assert.Equal(t, tt.want, tt.group.HasMembers())
		})
	}
}
