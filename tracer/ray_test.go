package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRay(t *testing.T) {
	type args struct {
		o Point
		d Vector
	}
	tests := []struct {
		name string
		args args
		want Ray
	}{
		{
			name: "test1",
			args: args{
				o: NewPoint(1, 2, 3),
				d: NewVector(4, 5, 6),
			},
			want: Ray{Origin: NewPoint(1, 2, 3), Dir: NewVector(4, 5, 6)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewRay(tt.args.o, tt.args.d), "should be equal")
		})
	}
}

func TestRay_Position(t *testing.T) {
	type args struct {
		t float64
	}
	tests := []struct {
		name string
		ray  Ray
		args args
		want Point
	}{
		{
			name: "test1",
			ray:  NewRay(NewPoint(2, 3, 4), NewVector(1, 0, 0)),
			args: args{
				t: 0,
			},
			want: NewPoint(2, 3, 4),
		},
		{
			name: "test2",
			ray:  NewRay(NewPoint(2, 3, 4), NewVector(1, 0, 0)),
			args: args{
				t: 1,
			},
			want: NewPoint(3, 3, 4),
		},
		{
			name: "test3",
			ray:  NewRay(NewPoint(2, 3, 4), NewVector(1, 0, 0)),
			args: args{
				t: -1,
			},
			want: NewPoint(1, 3, 4),
		},
		{
			name: "test4",
			ray:  NewRay(NewPoint(2, 3, 4), NewVector(1, 0, 0)),
			args: args{
				t: 2.5,
			},
			want: NewPoint(4.5, 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ray.Position(tt.args.t), "should be equal")
		})
	}
}

func TestRay_Transform(t *testing.T) {
	type args struct {
		m Matrix
	}
	tests := []struct {
		name string
		ray  Ray
		args args
		want Ray
	}{
		{
			name: "translate1",
			ray:  NewRay(NewPoint(1, 2, 3), NewVector(0, 1, 0)),
			args: args{
				m: NewTranslation(3, 4, 5),
			},
			want: NewRay(NewPoint(4, 6, 8), NewVector(0, 1, 0)),
		},
		{
			name: "scale1",
			ray:  NewRay(NewPoint(1, 2, 3), NewVector(0, 1, 0)),
			args: args{
				m: NewScaling(2, 3, 4),
			},
			want: NewRay(NewPoint(2, 6, 12), NewVector(0, 3, 0)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.ray.Transform(tt.args.m), "should be equal")
		})
	}
}
