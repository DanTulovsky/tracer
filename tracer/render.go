package tracer

import "flag"

var (
	output = flag.String("output", "", "name of the output file, if empty, renders to screen")
)

// Render runs the render
func Render(w *World) {
	if *output != "" {
		RenderToFile(w, *output)
	} else {
		RenderLive(w)
	}
}
