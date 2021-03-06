package tracer

import (
	"testing"

	"golang.org/x/image/colornames"

	"github.com/stretchr/testify/assert"
)

type testPattern struct {
	basePattern
}

func newTestPattern() Patterner {
	return &testPattern{
		basePattern: basePattern{
			transform: IM(),
		},
	}
}

func (tp *testPattern) ColorAtObject(o Shaper, p Point) Color {
	return NewColor(p.X(), p.Y(), p.Z())
}

func TestNewStripedPattern(t *testing.T) {
	type args struct {
		c1 Color
		c2 Color
	}
	tests := []struct {
		name string
		args args
		want *StripedPattern
	}{
		{
			name: "test1",
			args: args{
				c1: Black(),
				c2: White(),
			},
			want: &StripedPattern{
				a: Black(),
				b: White(),
				basePattern: basePattern{
					transform:        IM(),
					transformInverse: IM().Inverse(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewStripedPattern(tt.args.c1, tt.args.c2))
		})
	}
}

func Test_basePattern_ColorAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name string
		bp   *basePattern
		o    Shaper
		args args
		want Color
	}{
		{
			name: "panic1",
			bp:   &basePattern{},
			o:    NewUnitSphere(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Panics(t, func() { tt.bp.ColorAtObject(tt.o, NewPoint(0, 0, 0)) }, "should panic")
		})
	}
}

func TestStripedPattern_colorAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name    string
		pattern *StripedPattern
		args    args
		want    Color
	}{
		{
			name:    "test1",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "test2",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 1, 0),
			},
			want: White(),
		},
		{
			name:    "test3",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 2, 0),
			},
			want: White(),
		},
		{
			name:    "test4",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "test5",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 1),
			},
			want: White(),
		},
		{
			name:    "test6",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 2),
			},
			want: White(),
		},
		{
			name:    "test7",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "test8",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(0.9, 0, 0),
			},
			want: White(),
		},
		{
			name:    "test9",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(1, 0, 0),
			},
			want: Black(),
		},
		{
			name:    "test10",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(-0.1, 0, 0),
			},
			want: Black(),
		},
		{
			name:    "test11",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(-1, 0, 0),
			},
			want: Black(),
		},
		{
			name:    "test12",
			pattern: NewStripedPattern(White(), Black()),
			args: args{
				p: NewPoint(-1.1, 0, 0),
			},
			want: White(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pattern.colorAt(tt.args.p))
		})
	}
}

func TestStripedPattern_ColorAtObject(t *testing.T) {
	type args struct {
		o Shaper
		p Point
	}
	tests := []struct {
		name       string
		pattern    Patterner
		oTransform Matrix // object transform
		pTransform Matrix // pattern transform
		args       args
		want       Color
	}{
		{
			name:       "object transformation",
			pattern:    NewStripedPattern(White(), Black()),
			oTransform: IM().Scale(2, 2, 2),
			pTransform: IM(),
			args: args{
				o: NewUnitSphere(),
				p: NewPoint(1.5, 0, 0),
			},
			want: White(),
		},
		{
			name:       "pattern transform",
			pattern:    NewStripedPattern(White(), Black()),
			oTransform: IM(),
			pTransform: IM().Scale(2, 2, 2),
			args: args{
				o: NewUnitSphere(),
				p: NewPoint(1.5, 0, 0),
			},
			want: White(),
		},
		{
			name:       "object and pattern transform",
			pattern:    NewStripedPattern(White(), Black()),
			oTransform: IM().Scale(2, 2, 2),
			pTransform: IM().Translate(0.5, 0, 0),
			args: args{
				o: NewUnitSphere(),
				p: NewPoint(1.5, 0, 0),
			},
			want: White(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.o.SetTransform(tt.oTransform)
			tt.pattern.SetTransform(tt.pTransform)
			assert.Equal(t, tt.want, tt.pattern.ColorAtObject(tt.args.o, tt.args.p))
		})
	}
}

func TestNewGradientPattern(t *testing.T) {
	type args struct {
		c1 Color
		c2 Color
	}
	tests := []struct {
		name string
		args args
		want *GradientPattern
	}{
		{
			name: "test1",
			args: args{
				c1: White(),
				c2: Black(),
			},
			want: &GradientPattern{
				a: White(),
				b: Black(),
				basePattern: basePattern{
					transform:        IM(),
					transformInverse: IM().Inverse(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, NewGradientPattern(tt.args.c1, tt.args.c2))
		})
	}
}

func TestGradientPattern_colorAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name    string
		pattern *GradientPattern
		args    args
		want    Color
	}{
		{
			name:    "test1",
			pattern: NewGradientPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "test2",
			pattern: NewGradientPattern(White(), Black()),
			args: args{
				p: NewPoint(0.25, 0, 0),
			},
			want: NewColor(0.75, 0.75, 0.75),
		},
		{
			name:    "test3",
			pattern: NewGradientPattern(White(), Black()),
			args: args{
				p: NewPoint(0.5, 0, 0),
			},
			want: NewColor(0.5, 0.5, 0.5),
		},
		{
			name:    "test4",
			pattern: NewGradientPattern(White(), Black()),
			args: args{
				p: NewPoint(0.75, 0, 0),
			},
			want: NewColor(0.25, 0.25, 0.25),
		},
		{
			name:    "test5",
			pattern: NewGradientPattern(White(), Black()),
			args: args{
				p: NewPoint(1, 0, 0),
			},
			want: White(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pattern.colorAt(tt.args.p))
		})
	}
}

func TestRingPattern_colorAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name    string
		pattern *RingPattern
		args    args
		want    Color
	}{
		{
			name:    "test1",
			pattern: NewRingPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "test2",
			pattern: NewRingPattern(White(), Black()),
			args: args{
				p: NewPoint(1, 0, 0),
			},
			want: Black(),
		},
		{
			name:    "test3",
			pattern: NewRingPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 1),
			},
			want: Black(),
		},
		{
			name:    "test4",
			pattern: NewRingPattern(White(), Black()),
			args: args{
				p: NewPoint(0.708, 0, 0.708),
			},
			want: Black(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pattern.colorAt(tt.args.p))
		})
	}
}

func TestCheckerPattern_colorAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name    string
		pattern *CheckerPattern
		args    args
		want    Color
	}{
		{
			name:    "x1",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "x2",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0.99, 0, 0),
			},
			want: White(),
		},
		{
			name:    "x3",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(1.01, 0, 0),
			},
			want: Black(),
		},
		{
			name:    "y1",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "y2",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0.99, 0),
			},
			want: White(),
		},
		{
			name:    "y3",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 1.01, 0),
			},
			want: Black(),
		},
		{
			name:    "z1",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0),
			},
			want: White(),
		},
		{
			name:    "z2",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 0.99),
			},
			want: White(),
		},
		{
			name:    "z3",
			pattern: NewCheckerPattern(White(), Black()),
			args: args{
				p: NewPoint(0, 0, 1.01),
			},
			want: Black(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.pattern.colorAt(tt.args.p))
		})
	}
}

func TestNewUVCheckersPattern(t *testing.T) {
	type args struct {
		w float64
		h float64
		a Color
		b Color
	}
	tests := []struct {
		name string
		args args
		want *UVCheckersPattern
	}{
		{
			args: args{
				a: Black(),
				b: White(),
				w: 10,
				h: 20,
			},
			want: &UVCheckersPattern{
				a: Black(),
				b: White(),
				w: 10,
				h: 20,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUVCheckersPattern(tt.args.w, tt.args.h, tt.args.a, tt.args.b)
			assert.Equal(t, tt.want, got, "should equal")

		})
	}
}

func TestUVCheckersPattern_uvColorAt(t *testing.T) {
	type args struct {
		u float64
		v float64
	}
	tests := []struct {
		name string
		p    *UVCheckersPattern
		args args
		want Color
	}{
		{
			p: NewUVCheckersPattern(2, 2, Black(), White()),
			args: args{
				u: 0.0,
				v: 0.0,
			},
			want: Black(),
		},
		{
			p: NewUVCheckersPattern(2, 2, Black(), White()),
			args: args{
				u: 0.5,
				v: 0.0,
			},
			want: White(),
		},
		{
			p: NewUVCheckersPattern(2, 2, Black(), White()),
			args: args{
				u: 0.0,
				v: 0.5,
			},
			want: White(),
		},
		{
			p: NewUVCheckersPattern(2, 2, Black(), White()),
			args: args{
				u: 0.5,
				v: 0.5,
			},
			want: Black(),
		},
		{
			p: NewUVCheckersPattern(2, 2, Black(), White()),
			args: args{
				u: 1.0,
				v: 1.0,
			},
			want: Black(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.UVColorAt(tt.args.u, tt.args.v)
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestTextureMapSpherical_colorAt(t *testing.T) {
	type args struct {
		p Point
	}
	tests := []struct {
		name string
		tm   *TextureMapPattern
		args args
		want Color
	}{
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(0.4315, 0.4670, 0.7719),
			},
			want: White(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(-0.9654, 0.2552, -0.0534),
			},
			want: Black(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(0.1039, 0.7090, 0.6975),
			},
			want: White(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(-0.4986, -0.7856, -0.3663),
			},
			want: Black(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(-0.0317, -0.9395, 0.3411),
			},
			want: Black(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(0.4809, -0.7721, 0.4154),
			},
			want: Black(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(0.0285, -0.9612, -0.2745),
			},
			want: Black(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(-0.5734, -0.2162, -0.7903),
			},
			want: White(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(0.7688, -0.1470, 0.6223),
			},
			want: Black(),
		},
		{
			tm: NewTextureMapPattern(NewUVCheckersPattern(16, 8, Black(), White()), NewSphericalMap()),
			args: args{
				p: NewPoint(-0.7652, 0.2175, 0.6060),
			},
			want: Black(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tm.colorAt(tt.args.p)
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestUVAlignCheckPattern_uvColorAt(t *testing.T) {
	type args struct {
		u, v float64
	}
	tests := []struct {
		name    string
		pattern *UVAlignCheckPattern
		args    args
		want    Color
	}{
		{
			name: "lmain",
			pattern: NewUVAlignCheckPattern(
				NewColor(1, 1, 1),
				NewColor(1, 0, 0),
				NewColor(1, 1, 0),
				NewColor(0, 1, 0),
				NewColor(0, 1, 1)),
			args: args{
				u: 0.5,
				v: 0.5,
			},
			want: NewColor(1, 1, 1),
		},
		{
			name: "ul",
			pattern: NewUVAlignCheckPattern(
				NewColor(1, 1, 1),
				NewColor(1, 0, 0),
				NewColor(1, 1, 0),
				NewColor(0, 1, 0),
				NewColor(0, 1, 1)),
			args: args{
				u: 0.1,
				v: 0.9,
			},
			want: NewColor(1, 0, 0),
		},
		{
			name: "ur",
			pattern: NewUVAlignCheckPattern(
				NewColor(1, 1, 1),
				NewColor(1, 0, 0),
				NewColor(1, 1, 0),
				NewColor(0, 1, 0),
				NewColor(0, 1, 1)),
			args: args{
				u: 0.9,
				v: 0.9,
			},
			want: NewColor(1, 1, 0),
		},
		{
			name: "bl",
			pattern: NewUVAlignCheckPattern(
				NewColor(1, 1, 1),
				NewColor(1, 0, 0),
				NewColor(1, 1, 0),
				NewColor(0, 1, 0),
				NewColor(0, 1, 1)),
			args: args{
				u: 0.1,
				v: 0.1,
			},
			want: NewColor(0, 1, 0),
		},
		{
			name: "br",
			pattern: NewUVAlignCheckPattern(
				NewColor(1, 1, 1),
				NewColor(1, 0, 0),
				NewColor(1, 1, 0),
				NewColor(0, 1, 0),
				NewColor(0, 1, 1)),
			args: args{
				u: 0.9,
				v: 0.1,
			},
			want: NewColor(0, 1, 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.pattern.UVColorAt(tt.args.u, tt.args.v)
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}

func TestCubeMapPattern_colorAt(t *testing.T) {

	left := NewUVAlignCheckPattern(
		ColorName(colornames.Yellow),
		ColorName(colornames.Cyan),
		ColorName(colornames.Red),
		ColorName(colornames.Blue),
		ColorName(colornames.Brown))
	front := NewUVAlignCheckPattern(
		ColorName(colornames.Cyan),
		ColorName(colornames.Red),
		ColorName(colornames.Yellow),
		ColorName(colornames.Brown),
		ColorName(colornames.Green))
	right := NewUVAlignCheckPattern(
		ColorName(colornames.Red),
		ColorName(colornames.Yellow),
		ColorName(colornames.Purple),
		ColorName(colornames.Green),
		ColorName(colornames.White))
	back := NewUVAlignCheckPattern(
		ColorName(colornames.Green),
		ColorName(colornames.Purple),
		ColorName(colornames.Cyan),
		ColorName(colornames.White),
		ColorName(colornames.Blue))
	up := NewUVAlignCheckPattern(
		ColorName(colornames.Brown),
		ColorName(colornames.Cyan),
		ColorName(colornames.Purple),
		ColorName(colornames.Red),
		ColorName(colornames.Yellow))
	down := NewUVAlignCheckPattern(
		ColorName(colornames.Purple),
		ColorName(colornames.Brown),
		ColorName(colornames.Green),
		ColorName(colornames.Blue),
		ColorName(colornames.White))
	cm := NewCubeMapPattern(left, front, right, back, up, down)

	type args struct {
		p Point
	}
	tests := []struct {
		name string
		args args
		want Color
	}{
		{
			args: args{
				p: NewPoint(-1, 0, 0),
			},
			want: ColorName(colornames.Yellow),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cm.colorAt(tt.args.p)
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}
