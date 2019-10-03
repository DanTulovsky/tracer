package tracer

import "math"

// Equals is used to compare floating point numbers
func Equals(a, b float64) bool {
	if math.Abs(a-b) < Epsilon {
		return true
	}
	return false
}
