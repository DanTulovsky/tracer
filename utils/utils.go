package utils

import (
	"math"
	"os"

	"github.com/DanTulovsky/tracer/constants"
)

// Equals is used to compare floating point numbers
func Equals(a, b float64) bool {
	return math.Abs(a-b) < constants.Epsilon
}

// Homedir returns the home directory of the user
func Homedir() string {
	homedir, _ := os.UserHomeDir()
	return homedir
}
