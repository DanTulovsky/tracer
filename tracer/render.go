package tracer

import (
	"flag"
	"fmt"
	"math"

	"github.com/rcrowley/go-metrics"
)

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

func showProgress(total, finished, last float64) float64 {
	every := 5.0
	// done := (width*y + x) / total * 100
	done := finished / total * 100
	if last < math.Floor(done) && math.Mod(math.Floor(done), every) == 0 {
		metrics.GetOrRegisterGaugeFloat64("total_progress_pcnt", nil).Update(done)
		last = math.Floor(done)

		if 100-every <= last {
			fmt.Println()
		}
	}

	return last
}
