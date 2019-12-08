package tracer

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestNewCSG(t *testing.T) {
	type args struct {
		s1 Shaper
		s2 Shaper
		op Operation
	}
	tests := []struct {
		name string
		args args
	}{
		{
			args: args{
				s1: NewUnitCube(),
				s2: NewUnitSphere(),
				op: Union,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := &CSG{
				left:  tt.args.s1,
				right: tt.args.s2,
				op:    Union,
				Shape: Shape{
					transform:        IM(),
					transformInverse: IM().Inverse(),
					material:         NewDefaultMaterial(),
					shape:            "csg",
				},
			}
			want.Shape.lna = want.localNormalAt
			got := NewCSG(tt.args.s1, tt.args.s2, tt.args.op)
			diff := cmp.Diff(want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
			assert.Equal(t, got, tt.args.s1.Parent(), "should equal")
			assert.Equal(t, got, tt.args.s2.Parent(), "should equal")
		})
	}
}

func TestCSG_IntersectionAllowed(t *testing.T) {
	type args struct {
		op   Operation
		lhit bool
		inl  bool
		inr  bool
	}
	tests := []struct {
		name string
		csg  *CSG
		args args
		want bool
	}{
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: true,
				inl:  true,
				inr:  true,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: true,
				inl:  true,
				inr:  false,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: true,
				inl:  false,
				inr:  true,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: true,
				inl:  false,
				inr:  false,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: false,
				inl:  true,
				inr:  true,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: false,
				inl:  true,
				inr:  false,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: false,
				inl:  false,
				inr:  true,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Union,
				lhit: false,
				inl:  false,
				inr:  false,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: true,
				inl:  true,
				inr:  true,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: true,
				inl:  true,
				inr:  false,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: true,
				inl:  false,
				inr:  true,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: true,
				inl:  false,
				inr:  false,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: false,
				inl:  true,
				inr:  true,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: false,
				inl:  true,
				inr:  false,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: false,
				inl:  false,
				inr:  true,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Intersect,
				lhit: false,
				inl:  false,
				inr:  false,
			},
			want: false,
		},

		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: true,
				inl:  true,
				inr:  true,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: true,
				inl:  true,
				inr:  false,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: true,
				inl:  false,
				inr:  true,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: true,
				inl:  false,
				inr:  false,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: false,
				inl:  true,
				inr:  true,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: false,
				inl:  true,
				inr:  false,
			},
			want: true,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: false,
				inl:  false,
				inr:  true,
			},
			want: false,
		},
		{
			csg: NewCSG(NewUnitSphere(), NewUnitCube(), Union),
			args: args{
				op:   Difference,
				lhit: false,
				inl:  false,
				inr:  false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.csg.IntersectionAllowed(tt.args.op, tt.args.lhit, tt.args.inl, tt.args.inr))
		})
	}
}

func TestCSG_FilterIntersections(t *testing.T) {
	type fields struct {
		left  Shaper
		right Shaper
		op    Operation
	}
	tests := []struct {
		name         string
		fields       fields
		want1, want2 int
	}{
		{
			fields: fields{
				left:  NewUnitSphere(),
				right: NewUnitCube(),
				op:    Union,
			},
			want1: 0,
			want2: 3,
		},
		{
			fields: fields{
				left:  NewUnitSphere(),
				right: NewUnitCube(),
				op:    Intersect,
			},
			want1: 1,
			want2: 2,
		},
		{
			fields: fields{
				left:  NewUnitSphere(),
				right: NewUnitCube(),
				op:    Difference,
			},
			want1: 0,
			want2: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			xs := NewIntersections(
				NewIntersection(tt.fields.left, 1),
				NewIntersection(tt.fields.right, 2),
				NewIntersection(tt.fields.left, 3),
				NewIntersection(tt.fields.right, 4),
			)
			csg := NewCSG(tt.fields.left, tt.fields.right, tt.fields.op)
			want := NewIntersections(
				NewIntersection(csg, xs[tt.want1].t),
				NewIntersection(csg, xs[tt.want2].t),
			)
			assert.Equal(t, want[0].t, csg.FilterIntersections(xs)[0].t, "should equal")
			assert.Equal(t, want[1].t, csg.FilterIntersections(xs)[1].t, "should equal")
		})
	}
}

func TestCSG_IntersectWith1(t *testing.T) {
	c := NewCSG(NewUnitSphere(), NewUnitCube(), Union)
	r := NewRay(NewPoint(0, 2, -5), NewVector(0, 0, 1))

	xs := c.IntersectWith(r, NewIntersections())
	assert.Equal(t, 0, len(xs), "should be 0")
}
func TestCSG_IntersectWith2(t *testing.T) {

	s1 := NewUnitSphere()
	s2 := NewUnitSphere()
	s2.SetTransform(IM().Translate(0, 0, 0.5))

	c := NewCSG(s1, s2, Union)
	r := NewRay(NewPoint(0, 0, -5), NewVector(0, 0, 1))

	xs := c.IntersectWith(r, NewIntersections())
	assert.Equal(t, 2, len(xs), "should be 0")

	assert.Equal(t, 4.0, xs[0].t, "should equal")
	assert.Equal(t, s1, xs[0].Object(), "should equal")

	assert.Equal(t, 6.5, xs[1].t, "should equal")
	assert.Equal(t, s2, xs[1].Object(), "should equal")
}
