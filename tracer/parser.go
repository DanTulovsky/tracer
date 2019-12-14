package tracer

// Parser for OBJ files
// Triangulates the incoming faces and stores them as groups of trianges
// TODO: Reimplement as a TriangleMesh
// https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-polygon-mesh

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"

	"github.com/mokiat/go-data-front/decoder/mtl"
	"github.com/mokiat/go-data-front/decoder/obj"
)

const (
	maxMaterials = 30
)

func parseMTL(model *obj.Model, dir string) (*mtl.Library, error) {

	lib := &mtl.Library{
		Materials: []*mtl.Material{},
	}

	libDecoder := mtl.NewDecoder(mtl.DecodeLimits{MaxMaterialCount: maxMaterials})

	for _, ml := range model.MaterialLibraries {
		f, err := os.Open(path.Join(dir, ml))
		if err != nil {
			return nil, err
		}
		l, err := libDecoder.Decode(f)
		if err != nil {
			return nil, err
		}
		lib.Materials = append(lib.Materials, l.Materials...)
	}
	return lib, nil
}

// parseOBJ implements OBJ parsing and returns the model
// dir is the directory that holds the .mtl files
func parseOBJ(f *os.File, dir string) (*obj.Model, *mtl.Library, error) {
	limits := obj.DefaultLimits()
	limits.MaxReferenceCount = 128
	decoder := obj.NewDecoder(limits)

	model, err := decoder.Decode(f)
	if err != nil {
		return nil, nil, err
	}

	log.Printf("Model has %d vertices.\n", len(model.Vertices))
	log.Printf("Model has %d texture coordinates.\n", len(model.TexCoords))
	log.Printf("Model has %d normals.\n", len(model.Normals))
	log.Printf("Model has %d objects.\n", len(model.Objects))
	log.Printf("Model has %d material libs.\n", len(model.MaterialLibraries))
	for _, ml := range model.MaterialLibraries {
		log.Printf("  %v", ml)
	}

	lib, err := parseMTL(model, dir)
	if err != nil {
		return nil, nil, err
	}

	return model, lib, nil
}

// boundingBoxFromPoints returns the bounding box given a list of points
func boundingBoxFromPoints(points ...Point) Bound {

	var x []float64
	var y []float64
	var z []float64

	for _, p := range points {
		x = append(x, p.X())
		y = append(y, p.Y())
		z = append(z, p.Z())
	}

	sort.Float64s(x)
	sort.Float64s(y)
	sort.Float64s(z)

	return Bound{
		Min: NewPoint(x[0], y[0], z[0]),
		Max: NewPoint(x[len(x)-1], y[len(y)-1], z[len(z)-1]),
	}
}

// normalizeOBJ resizes the vertecies to all live in a box (-1, -1, -1) - (1, 1, 1)
func normalizeOBJ(vertecies []Point) []Point {
	log.Println("normalizing obj input...")
	result := []Point{}

	bbox := boundingBoxFromPoints(vertecies...)

	sx := bbox.Max.x - bbox.Min.x
	sy := bbox.Max.y - bbox.Min.y
	sz := bbox.Max.z - bbox.Min.z

	scale := math.Max(math.Max(sx, sy), sz) / 2

	for _, v := range vertecies {
		new := NewPoint(0, 0, 0)
		new.x = (v.x - (bbox.Min.x + sx/2)) / scale
		new.y = (v.y - (bbox.Min.y + sy/2)) / scale
		new.z = (v.z - (bbox.Min.z + sz/2)) / scale
		result = append(result, new)
	}
	return result

}

// triangulate converts a face into a list of triangles
func triangulate(model *obj.Model, f *obj.Face, mat *Material) []Shaper {
	var tri []Shaper
	var vertecies []Point
	var normals []Vector
	var textures []Point

	for _, r := range f.References {
		v := model.GetVertexFromReference(r)
		// negate Z because OBJ uses right-handed coordinates, and we use left-handed coordinates
		vertecies = append(vertecies, NewPoint(v.X, v.Y, -v.Z))

		n := model.GetNormalFromReference(r)
		normals = append(normals, NewVector(n.X, n.Y, -n.Z))

		t := model.GetTexCoordFromReference(r)
		textures = append(textures, NewPoint(t.U, 1-t.V, t.W))
	}

	// TODO: Run this for ALL vertecies at the same time, here it's just one face at a time
	// http://forum.raytracerchallenge.com/thread/27/triangle-mesh-normalization
	// vertecies = normalizeOBJ(vertecies)

	for i := 1; i < len(vertecies)-1; i++ {
		t := NewSmoothTriangle(
			vertecies[0], vertecies[i], vertecies[i+1],
			normals[0], normals[i], normals[i+1],
			textures[0], textures[i], textures[i+1])
		t.SetMaterial(mat)
		tri = append(tri, t)
	}

	return tri
}

// processIllum sets various material settings based on the illum parameter
// TODO: Implement this
func processIllum(mat *mtl.Material, illum int64) *mtl.Material {
	return mat
}

// convertMaterial converts OBJ material to *Material
func convertMaterial(mat *mtl.Material, dir string) (*Material, error) {
	// https://people.sc.fsu.edu/~jburkardt/data/mtl/mtl.html

	// defines the ambient color of the material to be (r,g,b).
	ka := mat.AmbientColor
	// defines the diffuse color of the material to be (r,g,b)
	kd := mat.DiffuseColor
	// defines the specular color of the material to be (r,g,b). This color shows up in highlights.
	ks := mat.SpecularColor

	kaColor := NewColor(ka.R, ka.G, ka.B)
	kdColor := NewColor(kd.R, kd.G, kd.B)
	ksColor := NewColor(ks.R, ks.G, ks.B)
	// ke := mat.EmissiveCoefficient  // not implemented in library

	// Dissolve indicates how much an object should blend.
	// The value should range between 0.0 (fully transparent)
	// and 1.0 (opaque).  In Blender this is the Alpha value in the BSDF shader.
	d := 1 - mat.Dissolve

	// defines the shininess of the material to be s. The default is 0.0;
	ns := mat.SpecularExponent

	illum := mat.Illum

	log.Printf("    %v\n", mat.Name)
	log.Printf("    Diffuse: %v\n", kdColor)
	log.Printf("    Diffuse texture: %v\n", mat.DiffuseTexture)
	log.Printf("    Bump texture: %v\n", mat.BumpTexture)
	log.Printf("    Ambient: %v\n", kaColor)
	log.Printf("    Specular: %v\n", ksColor)
	log.Printf("    Specular Exp: %v\n", ns)
	log.Printf("    Transparency: %v\n", d)
	log.Printf("    Illumination: %v\n", illum)

	m := NewDefaultMaterial()
	m.Color = kdColor
	m.Shininess = ns
	// d = 0 is fully transparent; the reverse of what we use
	m.Transparency = d
	if m.Transparency > 0 {
		m.ShadowCaster = false
	}

	// TODO: Implement support for illum, probably in lighting()
	// http://paulbourke.net/dataformats/mtl/
	mat = processIllum(mat, illum)

	// If there is a bump map present, use it
	if mat.BumpTexture != "" {
		log.Println("Reading in bump map textures...")

		imageFile := path.Join(dir, mat.BumpTexture)

		pert, err := NewImageHeightmapPerturber(imageFile, NewPlaneMap())
		if err != nil {
			return nil, err
		}
		m.SetPerturber(pert)
	}

	// If there is a texture present, use it
	if mat.DiffuseTexture != "" {
		log.Println("Reading in material textures...")

		imageFile := path.Join(dir, mat.DiffuseTexture)
		f, err := os.Open(imageFile)
		if err != nil {
			return nil, err
		}

		decode, format, err := image.Decode(f)
		if err != nil {
			return nil, err
		}
		log.Printf("decoded image format %v", format)

		// store the texture in the material, only used by smooth triangles
		m.AddDiffuseTexture(mat.Name, decode)
	}

	return m, nil
}

// convertData converts the parsed model to *Group instance
func convertData(model *obj.Model, lib *mtl.Library, dir string) (*Group, error) {
	g := NewGroup()

	for _, o := range model.Objects {
		log.Printf("Object:%v", o.Name)
		for _, m := range o.Meshes {
			log.Printf("  material: %v\n", m.MaterialName)
			mat, ok := lib.FindMaterial(m.MaterialName)
			if !ok {
				return nil, fmt.Errorf("Unable to find material %v in lib", m.MaterialName)
			}

			omat, err := convertMaterial(mat, dir)
			if err != nil {
				return nil, err
			}

			log.Println("  Faces:")
			for _, f := range m.Faces {
				tri := triangulate(model, f, omat)
				g.AddMembers(tri...)
			}
		}
	}
	return g, nil
}

// toMesh converts an object to a TriangleMesh
func toMesh(model *obj.Model, o *obj.Object, lib *mtl.Library, dir string) (*TriangleMesh, error) {
	var tri *TriangleMesh
	var vertices []Point
	var normals []Vector
	var textures []Point
	var faceIndex []int
	var vertexIndex []int

	for _, m := range o.Meshes {
		log.Printf("  material: %v\n", m.MaterialName)
		// TODO: handle materials
		// mat, ok := lib.FindMaterial(m.MaterialName)
		// if !ok {
		// 	return nil, fmt.Errorf("Unable to find material %v in lib", m.MaterialName)
		// }

		// omat, err := convertMaterial(mat, dir)
		// if err != nil {
		// 	return nil, err
		// }

		log.Println("  Faces:")
		for _, f := range m.Faces {
			vcount := 0
			for _, r := range f.References {
				v := model.GetVertexFromReference(r)
				// negate Z because OBJ uses right-handed coordinates, and we use left-handed coordinates
				vertices = append(vertices, NewPoint(v.X, v.Y, -v.Z))

				n := model.GetNormalFromReference(r)
				normals = append(normals, NewVector(n.X, n.Y, -n.Z))

				t := model.GetTexCoordFromReference(r)
				textures = append(textures, NewPoint(t.U, 1-t.V, t.W))
				vertexIndex = append(vertexIndex, vcount)

				vcount++
			}
			faceIndex = append(faceIndex, vcount)
		}
		numFaces := len(m.Faces)
		tri = NewMesh(numFaces, faceIndex, vertexIndex, vertices, normals, textures)
	}

	return tri, nil
}

// convertToMesh converts the parsed model to a group of TriangleMesh objects
func convertDataToMesh(model *obj.Model, lib *mtl.Library, dir string) (*Group, error) {
	g := NewGroup()

	for _, o := range model.Objects {
		log.Printf("Object:%v", o.Name)
		m, err := toMesh(model, o, lib, dir)
		if err != nil {
			panic(err)
		}
		g.AddMember(m)
	}
	return g, nil
}

// ParseOBJ parses an OBJ file and returns the result as a group
func ParseOBJ(f string) (*Group, error) {
	file, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	model, lib, err := parseOBJ(file, filepath.Dir(f))
	if err != nil {
		return nil, err
	}

	// g, err := convertData(model, lib, filepath.Dir(f))
	g, err := convertDataToMesh(model, lib, filepath.Dir(f))
	if err != nil {
		return nil, err
	}

	return g, nil
}
