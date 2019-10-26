package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testPattern struct {
	basePattern
}

func newTestPattern() Patterner {
	return &testPattern{
		basePattern: basePattern{
			transform: IdentityMatrix(),
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
					transform: IdentityMatrix(),
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
			oTransform: IdentityMatrix().Scale(2, 2, 2),
			pTransform: IdentityMatrix(),
			args: args{
				o: NewUnitSphere(),
				p: NewPoint(1.5, 0, 0),
			},
			want: White(),
		},
		{
			name:       "pattern transform",
			pattern:    NewStripedPattern(White(), Black()),
			oTransform: IdentityMatrix(),
			pTransform: IdentityMatrix().Scale(2, 2, 2),
			args: args{
				o: NewUnitSphere(),
				p: NewPoint(1.5, 0, 0),
			},
			want: White(),
		},
		{
			name:       "object and pattern transform",
			pattern:    NewStripedPattern(White(), Black()),
			oTransform: IdentityMatrix().Scale(2, 2, 2),
			pTransform: IdentityMatrix().Translate(0.5, 0, 0),
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
					transform: IdentityMatrix(),
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
