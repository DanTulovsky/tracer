package tracer

import (
	"log"
)

// LintWorld runs some common checks and prints out log messages to stdout
func (w *World) LintWorld() {
	log.Println("Linting the world...")

	for _, o := range w.Objects {
		lintObject(o)
	}

	w.lintLights(w.Lights)
	log.Println()
}

// LintObject runs linter checks for objects
func lintObject(o Shaper) {
	if o.HasMembers() {
		for _, m := range o.(*Group).Members() {
			lintMaterial(m.Material(), m)
		}
	}
	lintMaterial(o.Material(), o)
}

// LintMaterial runs linter checks for material
func lintMaterial(m *Material, o Shaper) {

	// Transparency checks
	if m.Transparency > 0 {
		// Warn if ShadowCaster is set on transparent objects
		if m.ShadowCaster {
			log.Printf("[warning] Object [%v] has Transparency and ShadowCaster set.", o.Name())
		}
	}
}

func (w *World) lintLights(lights []Light) {

	haveAreaLights := false

	for _, l := range lights {
		switch l.(type) {
		case *AreaLight, *AreaSpotLight:
			haveAreaLights = true
			if !w.Config.SoftShadows {
				log.Printf("[warning] Have area lights, but soft shadows are off.")
			}
		}
	}

	if !haveAreaLights && w.Config.SoftShadows {
		log.Printf("[warning] Soft shadows are enabled, but no area lights present.")
	}

}
