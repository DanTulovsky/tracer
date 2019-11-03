package utils

import (
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/DanTulovsky/tracer/constants"
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Random returns a random int in [min, max)
func Random(min, max int) int {
	return rand.Intn(max-min) + min
}

// RandomFloat returns a random float in [min, max)
func RandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// Equals is used to compare floating point numbers
func Equals(a, b float64) bool {
	return math.Abs(a-b) < constants.Epsilon
}

// Homedir returns the home directory of the user
func Homedir() string {
	homedir, _ := os.UserHomeDir()
	return homedir
}

// AT x (in the range [a, b] to a number in [c, d]
func AT(x, a, b, c, d float64) float64 {
	// log.Printf("in: %v [%v, %v] -> [%v, %v]", x, a, b, c, d)
	if x < a {
		log.Print("invalid input into AffineTransform, returning min.")
		log.Printf("AffineTransform -> in: %v [%v, %v] -> [%v, %v]", x, a, b, c, d)
		return c
	}
	if x > b {
		log.Print("invalid input into AffineTransform, returning max.")
		log.Printf("AffineTransform -> in: %v [%v, %v] -> [%v, %v]", x, a, b, c, d)
		return d
	}
	return (x-a)*((d-c)/(b-a)) + c
}
