package nodeID

// xlattice_go/nodeID/idMap_perf_test.go

/////////////////////////////////////////////////////////////////////
// THIS NEEDS TO BE RUN WITH
//   go test -gocheck.b
/////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/binary"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	"time"
)

var _ = fmt.Print

// -- utilities -----------------------------------------------------

// Convert the first 8 bytes of a []byte into a uint64.
func slice2hash(b []byte) (hc uint64, err error) {
	buf := bytes.NewReader(b)
	err = binary.Read(buf, binary.LittleEndian, &hc)
	return
}

// Create N random-ish K-byte values.  It takes about 2 us to create
// a value (21.2 ms for 10K values, 2.008s for 1M values).  The keys
// must be unique within the set.  Interpret this as a requirement that
// the keys be unique in their first 64 bits.

func makeSomeKeys(N, K int) (keys [][]byte) {

	rng := xr.MakeSimpleRNG()
	keys = make([][]byte, N)
	keyMap := make(map[uint64]bool)

	for i := 0; i < N; i++ {
		key := make([]byte, K)
		for {
			rng.NextBytes(key)
			hc, _ := slice2hash(key)
			_, ok := keyMap[hc]
			if !ok { // value not in map
				keyMap[hc] = true
				break
			}
		}
		keys[i] = key
	}
	return
}

// -- tests proper --------------------------------------------------

func (s *XLSuite) BenchmarkWithIDMapKeys(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("BENCHMARK_WITH_ID_MAP_KEYS")
	}

	const MAX_KEY_DEPTH = 16 // bytes

	// build an array of N random-ish K-byte keys
	K := 32
	N := c.N
	t0 := time.Now()
	keys := makeSomeKeys(N, K)
	t1 := time.Now()
	deltaT := t1.Sub(t0)
	fmt.Printf("setup time for %d %d-byte keys: %v\n", N, K, deltaT)

	// build an IDMap to put them in (ignoring any error)
	m, _ := NewIDMap(MAX_KEY_DEPTH)

	c.ResetTimer()
	c.StartTimer()
	// my results: 1603 ns/op for a run of 1 million insertions
	for i := 0; i < c.N; i++ {
		_ = m.Insert(keys[i], keys[i])
	}
	c.StopTimer()

	// verify that the keys are present in the map
	for i := 0; i < N; i++ {
		value, err := m.Find(keys[i])
		c.Assert(err, IsNil)
		val := value.([]byte)
		c.Assert(bytes.Equal(val, keys[i]), Equals, true)

	}
}
