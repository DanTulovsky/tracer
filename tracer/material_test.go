package tracer

import (
	"testing"

	"golang.org/x/image/colornames"

	"github.com/stretchr/testify/assert"
)

func TestNewMaterial(t *testing.T) {
	type args struct {
		clr       Color
		ambient   float64
		diffuse   float64
		specular  float64
		shininess float64
	}
	tests := []struct {
		name string
		args args
		want Material
	}{
		{
			name: "test1",
			args: args{
				clr:       ColorName(colornames.Red),
				ambient:   0.5,
				diffuse:   0.4,
				specular:  0.3,
				shininess: 40,
			},
			want: Material{
				Color:     ColorName(colornames.Red),
				Ambient:   0.5,
				Diffuse:   0.4,
				Specular:  0.3,
				Shininess: 40,
			},
		},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, NewMaterial(tt.args.clr, tt.args.ambient, tt.args.diffuse, tt.args.specular, tt.args.shininess))
	}
}
