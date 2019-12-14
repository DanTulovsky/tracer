package tracer

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestNewMesh(t *testing.T) {
	type args struct {
		numFaces    int
		faceIndex   []int
		vertexIndex []int
		p           []Point
		normals     []Vector
		textures    []Point
	}
	tests := []struct {
		name string
		args args
		want *TriangleMesh
	}{
		{
			name: "one",
			args: args{
				numFaces:  2,
				faceIndex: []int{4, 4},
				vertexIndex: []int{
					0, 1, 2, 3, // first face
					0, 3, 4, 5, // second face
				},
				p: []Point{
					NewPoint(-5, -5, 5),
					NewPoint(5, -5, 5),
					NewPoint(5, -5, -5),
					NewPoint(-5, -5, -5),
					NewPoint(-5, 5, -5),
					NewPoint(5, -5, -5),
				},
				normals: []Vector{
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
				},
				textures: []Point{
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
				},
			},
			want: &TriangleMesh{
				V: []Point{
					NewPoint(-5, -5, 5),
					NewPoint(5, -5, 5),
					NewPoint(5, -5, -5),
					NewPoint(-5, -5, -5),
					NewPoint(-5, 5, -5),
					NewPoint(5, -5, -5),
				},
				Vn: []Vector{
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
					NewVector(0, 0, 0),
				},
				Vt: []Point{
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
					NewPoint(0, 0, 0),
				},
				TrisIndex: []int{
					0, 1, 2, // tri1
					0, 2, 3, // tri2
					0, 3, 4, // tri3
					0, 4, 5, // tri4
				},
				Shape: Shape{
					transform:        IM(),
					transformInverse: IM().Inverse(),
					material:         NewDefaultMaterial(),
					shape:            "trimesh",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewMesh(tt.args.numFaces, tt.args.faceIndex, tt.args.vertexIndex, tt.args.p, tt.args.normals, tt.args.textures)
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))
			assert.Equal(t, 4*3, len(got.TrisIndex)) // 4 triangles * 3 vertices

			bound := Bound{
				Min: NewPoint(-5, -5, -5),
				Max: NewPoint(5, 5, 5),
			}
			assert.Equal(t, bound, got.Bounds(), "should equal")
		})
	}
}
