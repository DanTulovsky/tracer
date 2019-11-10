package tracer

import (
	"flag"
	"fmt"
<<<<<<< HEAD
	"log"
=======
>>>>>>> 4b3cd7fcf042bee4bff311887865678d8f7940ba
	"math"
)

var (
	output = flag.String("output", "", "name of the output file, if empty, renders to screen")
)

// Render runs the render
func Render(w *World) {
<<<<<<< HEAD
	log.Println(*output)
=======
>>>>>>> 4b3cd7fcf042bee4bff311887865678d8f7940ba
	if *output != "" {
		RenderToFile(w, *output)
	} else {
		RenderLive(w)
	}
}

func showProgress(total, last, height, width, x, y float64) float64 {
	every := 5.0
	done := (width*y + x) / total * 100
	if last < math.Floor(done) && math.Mod(math.Floor(done), every) == 0 {
		fmt.Printf("...%.2f", done)
		last = math.Floor(done)
	}

	return last
}
