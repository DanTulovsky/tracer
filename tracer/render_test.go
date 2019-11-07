package tracer

import (
	"testing"
)

func BenchmarkRenderSphere(b *testing.B) {
	width, height := 300.0, 300.0
	w := NewDefaultWorld(width, height)

	s1 := NewUnitSphere()

	w.AddObject(s1)

	for n := 0; n < b.N; n++ {
		RenderToFile(w, "/tmp/output.png")
	}
}
