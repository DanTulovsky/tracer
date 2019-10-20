package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
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
		want *Material
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
			want: &Material{
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

func TestMaterial_HasPattern(t *testing.T) {
	tests := []struct {
		name string
		m    *Material
		p    Patterner
		want bool
	}{
		{
			name: "has pattern",
			m:    NewDefaultMaterial(),
			p:    NewStripedPattern(Black(), White()),
			want: true,
		},
		{
			name: "no pattern",
			m:    NewDefaultMaterial(),
			want: false,
		},
	}
	for _, tt := range tests {
		if tt.p != nil {
			tt.m.SetPattern(tt.p)
		}
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.m.HasPattern(), "should equal")
		})
	}
}
