package tracer

import (
	"testing"

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
					transform: IdentityMatrix(),
					material:  NewDefaultMaterial(),
					shape:     "csg",
				},
			}
			got := NewCSG(tt.args.s1, tt.args.s2, tt.args.op)
			assert.Equal(t, want, got, "should equal")
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
	type args struct {
		xs Intersections
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			fields: fields{
				left:  NewUnitSphere(),
				right: NewUnitCube(),
				op:    Union,
			},
			args: args{
				xs: NewIntersections(
					NewIntersection(NewUnitSphere(), 1),
					NewIntersection(NewUnitCube(), 2),
					NewIntersection(NewUnitSphere(), 3),
					NewIntersection(NewUnitCube(), 4),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			csg := NewCSG(tt.fields.left, tt.fields.right, tt.fields.op)
			want := NewIntersections(
				NewIntersection(csg, 0),
				NewIntersection(csg, 3),
			)
			assert.Equal(t, want, csg.FilterIntersections(tt.args.xs), "should equal")
		})
	}
}
