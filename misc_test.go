package xlattice_go

import . "github.com/jddixon/xlattice_go/rnglib"
import "time"

func MakeRNG() *SimpleRNG {
	t := time.Now().Unix()
	rng := NewSimpleRNG(t)
	return rng
}
