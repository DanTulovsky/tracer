package tracer

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestNewMesh(t *testing.T) {
	type args struct {
		numFaces      int
		faceIndex     []int
		vertexIndex   []int
		normalIndex   []int
		textureIndex  []int
		materialIndex []int
		p             []Point
		normals       []Vector
		textures      []Point
		materials     []*Material
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
				normalIndex: []int{
					0, 1, 2, 3, // first face
					0, 3, 4, 5, // second face
				},
				textureIndex: []int{
					0, 1, 2, 3, // first face
					0, 3, 4, 5, // second face
				},
				materialIndex: []int{0, 1},
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
				materials: []*Material{
					NewDefaultMaterial(),
					NewDefaultMaterial(),
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
			got := NewMesh(tt.args.numFaces,
				tt.args.faceIndex, tt.args.vertexIndex, tt.args.normalIndex, tt.args.textureIndex, tt.args.materialIndex,
				tt.args.p, tt.args.normals, tt.args.textures, tt.args.materials)
			diff := cmp.Diff(tt.want, got)
			assert.Equal(t, "", fmt.Sprint(diff))

			bound := Bound{
				Min: NewPoint(-5, -5, -5),
				Max: NewPoint(5, 5, 5),
			}
			assert.Equal(t, bound, got.Bounds(), "should equal")
		})
	}
}
