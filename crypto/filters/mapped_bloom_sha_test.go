package filters

// xlattice_go/crypto/filters/mapped_bloom_sha3_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	xf "github.com/jddixon/xlattice_go/util/lfs"
	. "launchpad.net/gocheck"
	"os"
)

// Bloom filters for sets whose members are SHA3 digests.

func setUpMB3(c *C, rng *xr.PRNG) (
	filter *MappedBloomSHA, m, k uint, keys [][]byte, backingFile string) {

	m = 20
	k = 8
	keys = make([][]byte, 100)
	for i := 0; i < 100; i++ {
		keys[i] = make([]byte, 20)
	}
	backingFile = "tmp/" + rng.NextFileName(8)
	// make sure the file does not already exist
	found, err := xf.PathExists(backingFile)
	c.Assert(err, IsNil)
	for found {
		backingFile = "tmp/" + rng.NextFileName(8)
		found, err = xf.PathExists(backingFile)
		c.Assert(err, IsNil)
	}
	return
}

func (s *XLSuite) doTestMappedInserts(c *C, m, k, numKey uint) {

	var err error

	rng := xr.MakeSimpleRNG()
	filter, m, k, keys, backingFile := setUpMB3(c, rng)

	keys = make([][]byte, numKey)
	for i := uint(0); i < numKey; i++ {
		keys[i] = make([]byte, 20)
	}
	// set up distinct keys
	for i := uint(0); i < numKey; i++ {
		for j := uint(0); j < 20; j++ {
			keys[i][j] = byte(0xff & (i + j + 100))
		}
	}
	filter, err = NewMappedBloomSHA(m, k, backingFile)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)
	for i := uint(0); i < numKey; i++ {
		c.Assert(filter.Size(), Equals, i)
		// before you insert the key, it's not there
		found, err := filter.Member(keys[i])
		c.Assert(err, IsNil)
		c.Assert(found, Equals, false)
		filter.Insert(keys[i])
	}
	for i := uint(0); i < numKey; i++ {
		// the keys we just inserted are in the filter
		found, err := filter.Member(keys[i])
		c.Assert(err, IsNil)
		c.Assert(found, Equals, true)
	}
}
func (s *XLSuite) TestMappedInserts(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MAPPED_INSERTS")
	}
	err := os.MkdirAll("tmp", 755)
	c.Assert(err, IsNil)

	m := uint(20)
	k := uint(8)

	s.doTestMappedInserts(c, m, k, 16)  // default values
	s.doTestMappedInserts(c, 14, 8, 16) // stride = 9
	s.doTestMappedInserts(c, 13, 8, 16) // stride = 8
	s.doTestMappedInserts(c, 12, 8, 16) // stride = 7

	s.doTestMappedInserts(c, 14, 7, 16) // stride = 9
	s.doTestMappedInserts(c, 13, 7, 16) // stride = 8
	s.doTestMappedInserts(c, 12, 7, 16) // stride = 7

	s.doTestMappedInserts(c, 14, 6, 16) // stride = 9
	s.doTestMappedInserts(c, 13, 6, 16) // stride = 8
	s.doTestMappedInserts(c, 12, 6, 16) // stride = 7

	s.doTestMappedInserts(c, 14, 5, 16) // stride = 9
	s.doTestMappedInserts(c, 13, 5, 16) // stride = 8
	s.doTestMappedInserts(c, 12, 5, 16) // stride = 7
}
