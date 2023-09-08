package random

import (
	"math/rand"
	"time"
)

// Seed the random number generator at package import time.
func init() {
	rand.Seed(time.Now().UnixNano())
}
