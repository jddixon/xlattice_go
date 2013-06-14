package rnglib

import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkRand(b *testing.B) {
	// typically 39.2 ns/op
	for i := 0; i < b.N; i++ {
		_ = rand.Int63() // Go's rand.Random package
	}
}

func BenchmarkSimpleRNG(b *testing.B) {
	t := time.Now().Unix()
	rng := NewSimpleRNG(t)
	b.ResetTimer()
	// typically 28.4 ns/op
	for i := 0; i < b.N; i++ {
		_ = rng.Int63() // Mersenne Twister
	}
}
func BenchmarkSystemRNG(b *testing.B) {
	t := time.Now().Unix()
	rng := NewSystemRNG(t)
	b.ResetTimer()
	// typically XX.X ns/op
	for i := 0; i < b.N; i++ {
		_ = rng.Int63() // /dev/urandom
	}
}
