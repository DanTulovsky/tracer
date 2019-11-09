package tracer

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSphericalMap(t *testing.T) {
	tests := []struct {
		name string
		want *SphericalMap
	}{
		{
			want: &SphericalMap{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewSphericalMap(), "should equal")
		})
	}
}

func TestSphericalMap_Map(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		sm    *SphericalMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			sm: NewSphericalMap(),
			args: args{
				p: NewPoint(0, 0, -1),
			},
			wantU: 0.0,
			wantV: 0.5,
		},
		{
			sm: NewSphericalMap(),
			args: args{
				p: NewPoint(1, 0, 0),
			},
			wantU: 0.25,
			wantV: 0.5,
		},
		{
			sm: NewSphericalMap(),
			args: args{
				p: NewPoint(0, 0, 1),
			},
			wantU: 0.5,
			wantV: 0.5,
		},
		{
			sm: NewSphericalMap(),
			args: args{
				p: NewPoint(-1, 0, 0),
			},
			wantU: 0.75,
			wantV: 0.5,
		},
		{
			sm: NewSphericalMap(),
			args: args{
				p: NewPoint(0, 1, 0),
			},
			wantU: 0.5,
			wantV: 1.0,
		},
		{
			sm: NewSphericalMap(),
			args: args{
				p: NewPoint(0, -1, 0),
			},
			wantU: 0.5,
			wantV: 0,
		},
		{
			sm: NewSphericalMap(),
			args: args{
				p: NewPoint(math.Sqrt2/2, math.Sqrt2/2, 0),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, gotV := tt.sm.Map(tt.args.p)
			assert.Equal(t, tt.wantU, gotU, "should equal")
			assert.Equal(t, tt.wantV, gotV, "should equal")
		})
	}
}
