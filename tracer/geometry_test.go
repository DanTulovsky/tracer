package tracer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// newTestShape returns a test shape for testing
func newTestShape() *Shape {
	return &Shape{
		transform: IdentityMatrix(),
		material:  NewDefaultMaterial(),
	}
}

func TestShape_Material(t *testing.T) {
	tests := []struct {
		name  string
		shape *Shape
		want  *Material
	}{
		{
			name:  "default",
			shape: newTestShape(),
			want:  NewDefaultMaterial(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.shape.Material(), "should equal")
		})
	}
}

func TestShape_SetMaterial(t *testing.T) {
	type args struct {
		m *Material
	}
	tests := []struct {
		name  string
		shape *Shape
		args  args
		want  *Material
	}{
		{
			name:  "test1",
			shape: newTestShape(),
			args: args{
				m: NewDefaultMaterial(),
			},
			want: NewDefaultMaterial(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.shape.SetMaterial(tt.args.m)
			assert.Equal(t, tt.want, tt.shape.Material(), "should equal")

		})
	}
}

func TestShape_Transform(t *testing.T) {
	tests := []struct {
		name  string
		shape *Shape
		want  Matrix
	}{
		{
			name:  "identity by default",
			shape: newTestShape(),
			want:  IdentityMatrix(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.shape.Transform(), "should equal")
		})
	}
}

func TestShape_SetTransform(t *testing.T) {
	type args struct {
		m Matrix
	}
	tests := []struct {
		name  string
		shape *Shape
		args  args
		want  Matrix
	}{
		{
			name:  "test1",
			shape: newTestShape(),
			args: args{
				m: NewTranslation(2, 3, 4),
			},
			want: NewTranslation(2, 3, 4),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.shape.SetTransform(tt.args.m)
			assert.Equal(t, tt.want, tt.shape.Transform(), "should equal")

		})
	}
}

func TestShape_Parent(t *testing.T) {
	tests := []struct {
		name  string
		shape *Shape
		want  Shaper
	}{
		{
			name:  "default group is nil",
			shape: &Shape{},
			want:  nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.shape.Parent(), "should equal")
		})
	}
}

func TestShape_HasParent(t *testing.T) {
	tests := []struct {
		name  string
		shape *Shape
		want  bool
	}{
		{
			name:  "default group is nil",
			shape: &Shape{},
			want:  false,
		},
		{
			name: "has parent",
			shape: &Shape{
				parent: NewGroup(),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.shape.HasParent(), "should equal")
		})
	}
}

func TestBound_Center(t *testing.T) {
	type fields struct {
		Min Point
		Max Point
	}
	tests := []struct {
		name   string
		fields fields
		want   Point
	}{
		{
			fields: fields{
				Min: NewPoint(-1, -1, -1),
				Max: NewPoint(1, 1, 1),
			},
			want: NewPoint(0, 0, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := Bound{
				Min: tt.fields.Min,
				Max: tt.fields.Max,
			}
			got := b.Center()
			assert.Equal(t, tt.want, got, "should equal")
		})
	}
}
