package tracer

import (
	"math"

	"github.com/DanTulovsky/tracer/constants"
	"github.com/google/go-cmp/cmp"
)

// TriangleMesh is a mesh made up entirely of triangles
type TriangleMesh struct {
	// vertices of the mesh
	V []Point
	// per vertex normal
	Vn []Vector
	// per vertex texture coordinate, only x,y is used
	Vt []Point
	// each number is the index in 'v', three consecutive numbers form a triangle
	TrisIndex []int

	Shape
}

// NewMesh generates a new polygon mesh by triangulating the input
// numFaces: total number of faces
// faceIndex: how many vertices each face is made of
// vertexIndex:  lists the vertecies (indexed into verts) for each face
// verts: list of vertices
func NewMesh(numFaces int, faceIndex, vertexIndex []int, verts []Point, normals []Vector, textures []Point) *TriangleMesh {
	// how many triangles we need to create
	var numTris int
	// total number of vertices
	var k int

	// largest vertex index in vertexIndex
	var maxVertIndex int

	for i := 0; i < numFaces; i++ {
		numTris = numTris + faceIndex[i] - 2
		for j := 0; j < faceIndex[i]; j++ {
			if vertexIndex[k+j] > maxVertIndex {
				maxVertIndex = vertexIndex[k+j]
			}
		}
		k = k + faceIndex[i]
	}
	maxVertIndex = maxVertIndex + 1

	// store only those vertices we use
	v := make([]Point, maxVertIndex)
	vn := make([]Vector, maxVertIndex)
	vt := make([]Point, maxVertIndex)
	for i := 0; i < maxVertIndex; i++ {
		v[i] = verts[i]
		vn[i] = normals[i]
		vt[i] = textures[i]
	}

	trisIndex := make([]int, numTris*3)

	var l int
	k = 0
	for i := 0; i < numFaces; i++ { // for each face
		for j := 0; j < faceIndex[i]-2; j++ { // for each triangle
			trisIndex[l] = vertexIndex[k]
			trisIndex[l+1] = vertexIndex[k+j+1]
			trisIndex[l+2] = vertexIndex[k+j+2]
			l = l + 3
		}
		k = k + faceIndex[i]
	}

	m := &TriangleMesh{
		V:         v,
		Vn:        vn,
		Vt:        vt,
		TrisIndex: trisIndex,
		Shape: Shape{
			transform:        IM(),
			transformInverse: IM().Inverse(),
			material:         NewDefaultMaterial(),
			shape:            "trimesh",
		},
	}
	m.calculateBounds()
	return m
}

// Equal returns true if the meshes are equal
func (m *TriangleMesh) Equal(m2 *TriangleMesh) bool {
	return m.Shape.Equal(&m2.Shape) &&
		cmp.Equal(m.V, m2.V) &&
		cmp.Equal(m.Vn, m2.Vn) &&
		cmp.Equal(m.Vt, m2.Vt) &&
		cmp.Equal(m.TrisIndex, m2.TrisIndex)
}

// calculateBounds sets the m.bound variable
func (m *TriangleMesh) calculateBounds() {
	m.bound = boundingBoxFromPoints(m.V...)
}

// Includes implements includes logic
func (m *TriangleMesh) Includes(m2 Shaper) bool {
	return m == m2
}

// checkAxis is a helper function for check for intersection of the group's bounding box and ray
func (m *TriangleMesh) checkAxis(o, d, min, max float64) (float64, float64) {

	var tmin, tmax float64

	tminNumerator := min - o
	tmaxNumerator := max - o

	if math.Abs(d) >= constants.Epsilon {
		tmin = tminNumerator / d
		tmax = tmaxNumerator / d
	} else {
		tmin = tminNumerator * math.MaxFloat64
		tmax = tmaxNumerator * math.MaxFloat64
	}

	if tmin > tmax {
		tmin, tmax = tmax, tmin
	}

	return tmin, tmax
}

// IntersectWithBoundingBox returns true if the ray intersects with the bounding box
// min and max define the bounding box
func (m *TriangleMesh) IntersectWithBoundingBox(r Ray, b Bound) bool {

	var tmin, tmax float64

	xtmin, xtmax := m.checkAxis(r.Origin.X(), r.Dir.X(), b.Min.X(), b.Max.X())
	ytmin, ytmax := m.checkAxis(r.Origin.Y(), r.Dir.Y(), b.Min.Y(), b.Max.Y())
	ztmin, ztmax := m.checkAxis(r.Origin.Z(), r.Dir.Z(), b.Min.Z(), b.Max.Z())

	tmin = math.Max(math.Max(xtmin, ytmin), ztmin)
	tmax = math.Min(math.Min(xtmax, ytmax), ztmax)

	// missed the bounding box
	if tmin > tmax {
		return false
	}
	return true
}

// NormalAt returns the normal vector at the given point on the surface of the mesh
func (m *TriangleMesh) NormalAt(p Point, xs *Intersection) Vector {
	panic("called NormalAt on a mesh")
}

// PrecomputeValues precomputes some values for render speedup
func (m *TriangleMesh) PrecomputeValues() {
}

// IntersectWith returns the 't' values of Ray r intersecting with the mesh
func (m *TriangleMesh) IntersectWith(r Ray, t Intersections) Intersections {
	// transform the ray by the inverse of the group transfrom matrix
	// instead of changing the group, we change the ray coming from the camera
	// by the inverse, which achieves the same thing
	r = r.Transform(m.transformInverse)

	if !m.IntersectWithBoundingBox(r, m.Bounds()) {
		// bail out early, ray does not intersect group bounding box
		return t
	}

	xs := NewIntersections()

	numTris := len(m.TrisIndex) / 3.0
	// check for intersection with every triangle
	j := 0
	for i := 0; i < numTris; i++ {
		tri := NewSmoothTriangle(
			m.V[m.TrisIndex[j]], m.V[m.TrisIndex[j+1]], m.V[m.TrisIndex[j+2]],
			m.Vn[m.TrisIndex[j]], m.Vn[m.TrisIndex[j+1]], m.Vn[m.TrisIndex[j+2]],
			m.Vt[m.TrisIndex[j]], m.Vt[m.TrisIndex[j+1]], m.Vt[m.TrisIndex[j+2]],
		)
		j = j + 3

		txs := tri.IntersectWith(r, xs)
		t = append(t, txs...)
		xs = xs[:0]
	}

	return t
}
