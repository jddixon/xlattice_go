package filters

// xlattice_go/crypto/filters/bloom_sha3_test.go

import (
	"fmt" // DEBUG
	. "launchpad.net/gocheck"
)

// Bloom filters for sets whose members are SHA3 digests.

func setUpB3() (filter *BloomSHA, n, m, k uint, keys [][]byte) {
	m = 20
	k = 8
	keys = make([][]byte, 100)
	for i := 0; i < 100; i++ {
		keys[i] = make([]byte, 20)
	}
	return
}

func (s *XLSuite) TestEmptyFilter(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_EMPTY_FILTER")
	}

	filter, n, m, k, keys := setUpB3()

	filter, err := NewBloomSHA(m, k)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)
	c.Assert(filter.Size(), Equals, uint(0))
	c.Assert(filter.Capacity(), Equals, uint(2<<(m-1)))

	_, _, _ = filter, n, keys
}

// Verify that out of range or otherwise unacceptable constructor
// parameters are caught.
func (s *XLSuite) TestParamExceptions(c *C) {

	if VERBOSITY > 0 {
		fmt.Println("TEST_PARAM_EXCEPTIONS")
	}
	// m checks

	// zero filter size exponent
	_, err := NewBloomSHA(0, 8)
	c.Assert(err, NotNil)

	// filter size exponent too large
	_, err = NewBloomSHA(25, 8)
	c.Assert(err, NotNil)

	// checks on k

	// zero hash function count
	_, err = NewBloomSHA(20, 0)
	c.Assert(err, NotNil)

	// invalid hash function count (k*m>256)
	_, err = NewBloomSHA(20, 13)
	c.Assert(err, NotNil)
}

func (s *XLSuite) doTestInserts(c *C, m, k, numKey uint) {

	var err error
	filter, n, m, k, keys := setUpB3()
	_ = n

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
	filter, err = NewBloomSHA(m, k) // default m=20, k=8
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)
	for i := uint(0); i < numKey; i++ {
		c.Assert(filter.Size(), Equals, i)
		isMember, ks, err := filter.Member(keys[i])
		_ = ks // XXX NOT YET
		c.Assert(err, IsNil)
		c.Assert(isMember, Equals, false)
		filter.Insert(keys[i])
	}
	for i := uint(0); i < numKey; i++ {
		isMember, ks, err := filter.Member(keys[i])
		_ = ks // XXX NEW, NOT YET TESTED
		c.Assert(err, IsNil)
		c.Assert(isMember, Equals, true)
	}
}
func (s *XLSuite) TestInserts(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_INSERTS")
	}
	m := uint(20)
	k := uint(8)

	s.doTestInserts(c, m, k, 16)  // default values
	s.doTestInserts(c, 14, 8, 16) // stride = 8
	s.doTestInserts(c, 13, 8, 16) // stride = 7
	s.doTestInserts(c, 12, 8, 16) // stride = 6

	s.doTestInserts(c, 14, 7, 16) // stride = 8
	s.doTestInserts(c, 13, 7, 16) // stride = 7
	s.doTestInserts(c, 12, 7, 16) // stride = 6

	s.doTestInserts(c, 14, 6, 16) // stride = 8
	s.doTestInserts(c, 13, 6, 16) // stride = 7
	s.doTestInserts(c, 12, 6, 16) // stride = 6

	s.doTestInserts(c, 14, 5, 16) // stride = 8
	s.doTestInserts(c, 13, 5, 16) // stride = 7
	s.doTestInserts(c, 12, 5, 16) // stride = 6
}
