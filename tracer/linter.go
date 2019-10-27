package tracer

import "log"

// LintWorld runs some common checks and prints out log messages to stdout
func (w *World) LintWorld() {

	for _, o := range w.Objects {
		LintObject(o)
	}

}

// LintObject runs linter checks for objects
func LintObject(o Shaper) {

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
