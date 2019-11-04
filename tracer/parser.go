package tracer

// Parser for OBJ files
// Triangulates the incoming faces and stores them as groups of trianges
// TODO: Reimplement as a TriangleMesh
// https://www.scratchapixel.com/lessons/3d-basic-rendering/ray-tracing-polygon-mesh

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/mokiat/go-data-front/decoder/mtl"
	"github.com/mokiat/go-data-front/decoder/obj"
)

const (
	maxMaterials = 10
)

func parseMaterials(model *obj.Model, dir string) (*mtl.Library, error) {

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

	lib, err := parseMaterials(model, dir)
	if err != nil {
		return nil, nil, err
	}

	return model, lib, nil
}

// convertData converts the parsed model to *Group instance
func convertData(model *obj.Model, lib *mtl.Library) (*Group, error) {
	g := NewGroup()

	for _, o := range model.Objects {
		log.Printf("Object:%v", o.Name)
		for _, m := range o.Meshes {
			log.Printf("  material: %v\n", m.MaterialName)
			mat, ok := lib.FindMaterial(m.MaterialName)
			if !ok {
				return nil, fmt.Errorf("Unable to find material %v in lib", m.MaterialName)
			}
			log.Printf("    %v\n", mat.Name)
			log.Printf("    Diffuse: %v\n", mat.DiffuseColor)
			log.Printf("    Ambient: %v\n", mat.AmbientColor)
			log.Printf("    Specular: %v\n", mat.SpecularColor)
			log.Printf("    Specular Exp: %v\n", mat.SpecularExponent)
			log.Println("  Faces:")
			for i, f := range m.Faces {
				log.Printf("  (%v)", i)
				for _, r := range f.References {
					log.Printf("    vertex: %v", model.GetVertexFromReference(r))
					log.Printf("    normal: %v", model.GetNormalFromReference(r))
					log.Printf("    texture: %v", model.GetTexCoordFromReference(r))
					log.Println()
				}
			}
		}
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

	g, err := convertData(model, lib)
	if err != nil {
		return nil, err
	}

	return g, nil
}
