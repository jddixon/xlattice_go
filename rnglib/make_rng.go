package rnglib

// XXX It seems to be an undocumented feature of Go that functions in files
// names like *_test.go are visible to go test when run the same directory
// but not visible elsewhere.  This this file, which used to be called
// misc_test.go is now make_rng.go, and I can use MakeRNG() in other directories.

import (
	"time"
)

func MakeSimpleRNG() *PRNG {
	t := time.Now().Unix()
	rng := NewSimpleRNG(t)
	return rng
}
