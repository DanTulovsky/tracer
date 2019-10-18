package utils

import (
	"math"

	"github.com/DanTulovsky/tracer/constants"
)

// Equals is used to compare floating point numbers
func Equals(a, b float64) bool {
	return math.Abs(a-b) < constants.Epsilon
}
