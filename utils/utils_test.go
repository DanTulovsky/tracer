package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(t *testing.T) {
	type args struct {
		a float64
		b float64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "equals",
			args: args{
				a: 2.345685,
				b: 2.345684,
			},
			want: true,
		},
		{
			name: "not equals",
			args: args{
				a: 2.44567,
				b: 2.34568,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Equals(tt.args.a, tt.args.b)

			assert.Equal(t, tt.want, result, "float comparison")
		})
	}
}

func TestRandomFloat(t *testing.T) {
	type args struct {
		min float64
		max float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			args: args{
				min: 0.2,
				max: 30.0,
			},
		},
		{
			args: args{
				min: 0,
				max: 0.3,
			},
		},
		{
			args: args{
				min: 100.45,
				max: 100.46,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := RandomFloat(tt.args.min, tt.args.max)

			assert.GreaterOrEqual(t, got, tt.args.min, "should be >=")
			assert.Less(t, got, tt.args.max, "should be >=")
		})
	}
}
