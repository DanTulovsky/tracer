package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/image/colornames"
)

func TestNewMaterial(t *testing.T) {
	type args struct {
		clr             Color
		ambient         float64
		diffuse         float64
		specular        float64
		shininess       float64
		reflective      float64
		transparency    float64
		refractiveIndex float64
		perturber       Perturber
	}
	tests := []struct {
		name string
		args args
		want *Material
	}{
		{
			name: "test1",
			args: args{
				clr:             ColorName(colornames.Red),
				ambient:         0.5,
				diffuse:         0.4,
				specular:        0.3,
				shininess:       40,
				reflective:      0.5,
				transparency:    0.6,
				refractiveIndex: 0.7,
				perturber:       nil,
			},
			want: &Material{
				Color:           ColorName(colornames.Red),
				Ambient:         0.5,
				Diffuse:         0.4,
				Specular:        0.3,
				Shininess:       40,
				Reflective:      0.5,
				Transparency:    0.6,
				RefractiveIndex: 0.7,
				ShadowCaster:    true,
				perturber:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(
				t,
				tt.want,
				NewMaterial(
					tt.args.clr, tt.args.ambient, tt.args.diffuse,
					tt.args.specular, tt.args.shininess, tt.args.reflective, tt.args.transparency, tt.args.refractiveIndex, nil))
		})
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
		t.Run(tt.name, func(t *testing.T) {
			if tt.p != nil {
				tt.m.SetPattern(tt.p)
			}
			assert.Equal(t, tt.want, tt.m.HasPattern(), "should equal")
		})
	}
}

func TestNewDefaultMaterial(t *testing.T) {
	tests := []struct {
		name string
		want *Material
	}{
		{
			name: "test1",
			want: &Material{
				Color:           NewColor(1, 1, 1),
				Ambient:         0.1,
				Diffuse:         0.9,
				Specular:        0.9,
				Shininess:       200.0,
				Reflective:      0,
				Transparency:    0,
				RefractiveIndex: 1.0,
				ShadowCaster:    true,
				perturber:       nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDefaultMaterial()
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}
