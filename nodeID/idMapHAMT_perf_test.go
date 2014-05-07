package nodeID

// xlattice_go/nodeID/idMapHAMT_perf_test.go

/////////////////////////////////////////////////////////////////////
// THIS NEEDS TO BE RUN WITH
//   go test -gocheck.b
/////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"fmt"
	//gh "github.com/jddixon/hamt_go"
	//xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	"time"
)

var _ = fmt.Print

// -- utilities -----------------------------------------------------

// Create N random-ish K-byte values.  It takes about 2 us to create
// a value (21.2 ms for 10K values, 2.008s for 1M values)

//func makeSomeKeys(N, K int) (keys [][]byte) {
//
//	rng := xr.MakeSimpleRNG()
//	keys = make([][]byte, N)
//
//	for i := 0; i < N; i++ {
//		keys[i] = make([]byte, K)
//		rng.NextBytes(keys[i])
//	}
//	return
//}

// -- tests proper --------------------------------------------------

func (s *XLSuite) BenchmarkWithHAMTKeys(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("BENCHMARK_WITH_HAMT_KEYS")
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

	// build an IDMap to put them in
	m := NewIDMapHAMT(5, 16)

	c.ResetTimer()
	c.StartTimer()
	// HAMT results: ???? ns/op for a run of 1 million insertions
	for i := 0; i < c.N; i++ {
		_ = m.Insert(keys[i], keys[i])
	}
	c.StopTimer()

	// verify that the keys are present in the map
	for i := 0; i < N; i++ {
		value, err := m.Find(keys[i])
		c.Assert(err, IsNil)
		c.Assert(value, NotNil)
		if value == nil {
			break
		}
		val := value.([]byte)
		c.Assert(bytes.Equal(val, keys[i]), Equals, true)

	}
}
