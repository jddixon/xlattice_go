package rnglib

import (
	"math/rand"
)

// SimpleRNG uses the 64-bit Mersenne Twister as a source of random
// numbers and runs about 35x faster than SystemRNG, which draws values
// from /dev/urandom.  Both share the same interface, including Go's
// math.rand functions.

// SimpleRNG is entirely deterministic and will always produce the same
// sequence of values if given the same seed.

// rand.Source interface ////////////////////////////////////////////
// See mt19937-64.go

// SimpleRNG ////////////////////////////////////////////////////////
func NewMTSource(seed int64) rand.Source {
	var mt64 MT64
	mt64.Seed(seed)
	return &mt64
}

func NewSimpleRNG(seed int64) *PRNG {
	s := new(PRNG) // allocates
	src := NewMTSource(seed)
	s.rng = rand.New(src)
	s.Seed(seed)
	return s
}
