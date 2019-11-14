package tracer

import (
	"log"
)

// LintWorld runs some common checks and prints out log messages to stdout
func (w *World) LintWorld() {
	log.Println("Linting the world...")

	// if w.Config.Antialias != 1 && !utils.IsPowerOf2(w.Config.Antialias) {
	// 	log.Fatalf("world antialias parameter must be a power of 2 (have: %v)", w.Config.Antialias)
	// }

	for _, o := range w.Objects {
		LintObject(o)
	}

}

// LintObject runs linter checks for objects
func LintObject(o Shaper) {
	if o.HasMembers() {
		for _, m := range o.(*Group).Members() {
			LintMaterial(m.Material(), m)
		}
	}
	LintMaterial(o.Material(), o)
}

// LintMaterial runs linter checks for material
func LintMaterial(m *Material, o Shaper) {

	// Transparency checks
	if m.Transparency > 0 {
		// Warn if ShadowCaster is set on transparent objects
		if m.ShadowCaster {
			log.Printf("Object [%v] has Transparency and ShadowCaster set.", o.Name())
		}
	}
}
