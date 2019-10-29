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
	type members struct {
		s Shaper
		t Matrix
	}

	tests := []struct {
		name        string
		group       *Group
		members     []members
		args        args
		transform   Matrix
		wantXS      int
		wantShapers []int // index of the shapers we are expecting
	}{
		{
			name:      "empty group",
			group:     NewGroup(),
			members:   []members{},
			transform: IdentityMatrix(),
			args: args{
				r: NewRay(Origin(), NewVector(0, 0, 1)),
			},
			wantXS: 0,
		},
		{
			name:  "spheres",
			group: NewGroup(),
			members: []members{
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
			members: []members{
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

	// assert.Equal(t, want, got, "should be true")
	assert.True(t, want.Equals(got), "should be true")

}
