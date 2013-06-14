package rnglib

// xlattice_go/rnglib/system_rng.go

import (
	"bufio"
	"encoding/binary"
	"io"
	"math/rand"
	"os"
	"sync"
)

// SystemRNG draws values from /dev/urandom and runs about 30x slower
// than SimpleRNG, which relies upon the 64-bit Mersenne Twister.  Both
// share the same interface and implement Go's math.rand.

type uReader struct {
	name string
	f    io.Reader
	mu   sync.Mutex
}

func (r *uReader) Read(b []byte) (n int, err error) {
	// with locking, about 1050ns/op; without, about 900ns/op
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.f == nil {
		f, err := os.Open(r.name)
		if f == nil {
			return 0, err
		}
		r.f = bufio.NewReader(f)
	}
	return r.f.Read(b)
}

// rand.Source interface ////////////////////////////////////////////

func (r *uReader) Seed(seed int64) {
	_ = seed
}
func (r *uReader) Int63() int64 {
	var n uint64
	err := binary.Read(r, binary.LittleEndian, &n)
	if err != nil {
		panic("error reading from /dev/urandom")
	}
	val := int64(n >> 1)
	return val
}

// SystemRNG ////////////////////////////////////////////////////////

func NewSystemSource() rand.Source {
	var dReader uReader
	dReader.name = "/dev/urandom"
	return &dReader
}
func NewSystemRNG(seed int64) *PRNG {
	s := new(PRNG)
	src := NewSystemSource()
	s.rng = rand.New(src)
	return s
}
