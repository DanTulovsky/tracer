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

func TestPlaneMap_Map(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		pm    *PlaneMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			name: "1",
			pm:   NewPlaneMap(),
			args: args{
				p: NewPoint(0.25, 0, 0.5),
			},
			wantU: 0.25,
			wantV: 0.5,
		},
		{
			name: "2",
			pm:   NewPlaneMap(),
			args: args{
				p: NewPoint(0.25, 0, -0.25),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			name: "3",
			pm:   NewPlaneMap(),
			args: args{
				p: NewPoint(0.25, 0.5, -0.25),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			name: "4",
			pm:   NewPlaneMap(),
			args: args{
				p: NewPoint(1.25, 0, 0.5),
			},
			wantU: 0.25,
			wantV: 0.5,
		},
		{
			name: "5",
			pm:   NewPlaneMap(),
			args: args{
				p: NewPoint(0.25, 0, -1.75),
			},
			wantU: 0.25,
			wantV: 0.25,
		},
		{
			name: "6",
			pm:   NewPlaneMap(),
			args: args{
				p: NewPoint(1, 0, 1),
			},
			wantU: 0,
			wantV: 0,
		},
		{
			name: "7",
			pm:   NewPlaneMap(),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			wantU: 0,
			wantV: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, gotV := tt.pm.Map(tt.args.p)
			assert.Equal(t, tt.wantU, gotU, "should equal")
			assert.Equal(t, tt.wantV, gotV, "should equal")
		})
	}
}

func TestCylinderMap_Map(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		cm    *CylinderMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			name: "1",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(0, 0, -1),
			},
			wantU: 0,
			wantV: 0,
		},
		{
			name: "2",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(0, 0.5, -1),
			},
			wantU: 0,
			wantV: 0.5,
		},
		{
			name: "3",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(0, 1, -1),
			},
			wantU: 0,
			wantV: 0,
		},
		{
			name: "4",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(0.70711, 0.5, -0.70711),
			},
			wantU: 0.125,
			wantV: 0.5,
		},
		{
			name: "5",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(1, 0.5, 0),
			},
			wantU: 0.25,
			wantV: 0.5,
		},
		{
			name: "6",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(0.70711, 0.5, 0.70711),
			},
			wantU: 0.375,
			wantV: 0.5,
		},
		{
			name: "7",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(0, -0.25, 1),
			},
			wantU: 0.5,
			wantV: 0.75,
		},
		{
			name: "8",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(-0.70711, 0.5, 0.70711),
			},
			wantU: 0.625,
			wantV: 0.5,
		},
		{
			name: "9",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(-1, 1.25, 0),
			},
			wantU: 0.75,
			wantV: 0.25,
		},
		{
			name: "10",
			cm:   NewCylinderMap(),
			args: args{
				p: NewPoint(-0.70711, 0.5, -0.70711),
			},
			wantU: 0.875,
			wantV: 0.5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, gotV := tt.cm.Map(tt.args.p)
			assert.Equal(t, tt.wantU, gotU, "should equal")
			assert.Equal(t, tt.wantV, gotV, "should equal")
		})
	}
}

// func TestCube_Map(t *testing.T) {
// 	type args struct {
// 		p Point
// 	}
// 	tests := []struct {
// 		name  string
// 		cm    *CubeMap
// 		args  args
// 		wantU float64
// 		wantV float64
// 	}{
// 		{
// 			name: "test1",
// 			cm: NewCubeMap(),
// 			args: args{
//                p: NewPoint(-1, 0,0),
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotU, gotV := tt.cm.Map(tt.args.p)
// 			assert.Equal(t, tt.wantU, gotU, "should equal")
// 			assert.Equal(t, tt.wantV, gotV, "should equal")
// 		})
// 	}
// }

func TestCubeMap_faceFromPoint(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name string
		cm   *CubeMap
		args args
		want cubeFace
	}{
		{
			name: "left",
			cm:   &CubeMap{},
			args: args{
				p: NewPoint(-1, 0.5, -0.25),
			},
			want: cubeFaceLeft,
		},
		{
			name: "right",
			cm:   &CubeMap{},
			args: args{
				p: NewPoint(1.1, -0.75, 0.8),
			},
			want: cubeFaceRight,
		},
		{
			name: "front",
			cm:   &CubeMap{},
			args: args{
				p: NewPoint(0.1, 0.6, 0.9),
			},
			want: cubeFaceFront,
		},
		{
			name: "back",
			cm:   &CubeMap{},
			args: args{
				p: NewPoint(-0.7, 0, -2),
			},
			want: cubeFaceBack,
		},
		{
			name: "up",
			cm:   &CubeMap{},
			args: args{
				p: NewPoint(0.5, 1, 0.9),
			},
			want: cubeFaceUp,
		},
		{
			name: "down",
			cm:   &CubeMap{},
			args: args{
				p: NewPoint(-0.2, -1.3, 1.1),
			},
			want: cubeFaceDown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cm.faceFromPoint(tt.args.p)
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestCubeMap_uvFront(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		cm    *CubeMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(-0.5, 0.5, 1),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(0.5, -0.5, 1),
			},
			wantU: 0.75,
			wantV: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotu, gotv := tt.cm.uvFront(tt.args.p)
			assert.Equal(t, tt.wantU, gotu, "should equal")
			assert.Equal(t, tt.wantV, gotv, "should equal")
		})
	}
}

func TestCubeMap_uvBack(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		cm    *CubeMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(0.5, 0.5, -1),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(-0.5, -0.5, -1),
			},
			wantU: 0.75,
			wantV: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotu, gotv := tt.cm.uvBack(tt.args.p)
			assert.Equal(t, tt.wantU, gotu, "should equal")
			assert.Equal(t, tt.wantV, gotv, "should equal")
		})
	}
}

func TestCubeMap_uvLeft(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		cm    *CubeMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(-1, 0.5, -0.5),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(-1, -0.5, 0.5),
			},
			wantU: 0.75,
			wantV: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotu, gotv := tt.cm.uvLeft(tt.args.p)
			assert.Equal(t, tt.wantU, gotu, "should equal")
			assert.Equal(t, tt.wantV, gotv, "should equal")
		})
	}
}

func TestCubeMap_uvRight(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		cm    *CubeMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(1, 0.5, 0.5),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(1, -0.5, -0.5),
			},
			wantU: 0.75,
			wantV: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotu, gotv := tt.cm.uvRight(tt.args.p)
			assert.Equal(t, tt.wantU, gotu, "should equal")
			assert.Equal(t, tt.wantV, gotv, "should equal")
		})
	}
}

func TestCubeMap_uvUp(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		cm    *CubeMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(-0.5, 1, -0.5),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(0.5, 1, 0.5),
			},
			wantU: 0.75,
			wantV: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotu, gotv := tt.cm.uvUp(tt.args.p)
			assert.Equal(t, tt.wantU, gotu, "should equal")
			assert.Equal(t, tt.wantV, gotv, "should equal")
		})
	}
}

func TestCubeMap_uvDown(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name  string
		cm    *CubeMap
		args  args
		wantU float64
		wantV float64
	}{
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(-0.5, -1, 0.5),
			},
			wantU: 0.25,
			wantV: 0.75,
		},
		{
			cm: &CubeMap{},
			args: args{
				p: NewPoint(0.5, -1, -0.5),
			},
			wantU: 0.75,
			wantV: 0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotu, gotv := tt.cm.uvDown(tt.args.p)
			assert.Equal(t, tt.wantU, gotu, "should equal")
			assert.Equal(t, tt.wantV, gotv, "should equal")
		})
	}
}
