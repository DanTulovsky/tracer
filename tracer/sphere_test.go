package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSphere(t *testing.T) {
	type args struct {
		c Point
		r float64
	}
	tests := []struct {
		name string
		args args
		want Sphere
	}{
		{
			name: "test1",
			args: args{
				c: NewPoint(0, 0, 0),
				r: 1.0,
			},
			want: Sphere{NewPoint(0, 0, 0), 1.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewSphere(tt.args.c, tt.args.r))
		})
	}
}
