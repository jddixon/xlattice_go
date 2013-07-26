package node

import (
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestBadNodeIDs(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BAD_NODE_IDS")
	}
	c.Assert(false, Equals, IsValidID(nil))
	candidate := make([]byte, SHA1_LEN-1)
	c.Assert(false, Equals, IsValidID(candidate))
	candidate = make([]byte, SHA1_LEN)
	c.Assert(true, Equals, IsValidID(candidate))
	candidate = make([]byte, SHA1_LEN+1)
	c.Assert(false, Equals, IsValidID(candidate))
	candidate = make([]byte, SHA3_LEN)
	c.Assert(true, Equals, IsValidID(candidate)) // FOO
}
func (s *XLSuite) TestThisAndThat(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_THIS_AND_THAT")
	}
	rng := rnglib.MakeSimpleRNG()
	v1 := make([]byte, SHA1_LEN)
	rng.NextBytes(&v1)
	v2 := make([]byte, SHA1_LEN)
	rng.NextBytes(&v2)
	id1 := NewNodeID(v1)
	id2 := NewNodeID(v2)
	c.Assert(id1, Not(Equals), id2)

	v1a := id1.Value()
	v2a := id2.Value()

	// XXX gocheck cannot handle these comparisons
	// c.Assert(v1, Not(DeepEquals), v1a)				// 'Deep' is for desperation
	// c.Assert(v2, Not(Equals), v2a)

	// XXX not sure that gocheck results are meaningful
	c.Assert(&v1, Not(Equals), &v1a)
	c.Assert(&v2, Not(Equals), &v2a)

	c.Assert(SHA1_LEN, Equals, len(v1a))
	c.Assert(SHA1_LEN, Equals, len(v2a))
	for i := 0; i < SHA1_LEN; i++ {
		c.Assert(v1[i], Equals, v1a[i])
		c.Assert(v2[i], Equals, v2a[i])
	}
	c.Assert(false, Equals, id1.Equal(nil))
	c.Assert(true, Equals, id1.Equal(id1))
	c.Assert(false, Equals, id1.Equal(id2)) // FOO
}

func (s *XLSuite) TestComparator(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_COMPARATOR")
	}
	rng := rnglib.MakeSimpleRNG()
	v1 := make([]byte, SHA1_LEN)
	rng.NextBytes(&v1)
	v3 := make([]byte, SHA3_LEN)
	rng.NextBytes(&v3)
	id1 := NewNodeID(v1) // SHA1
	id3 := NewNodeID(v3) // SHA3

	// make clones which will sort before and after v1
	v1a := make([]byte, SHA1_LEN) // sorts AFTER v1
	for i := 0; i < SHA1_LEN; i++ {
		if v1[i] == 255 {
			v1a[i] = 255
		} else {
			v1a[i] = v1[i] + 1
			break
		}
	}
	v1b := make([]byte, SHA1_LEN) // sorts BEFORE v1
	for i := 0; i < SHA1_LEN; i++ {
		if v1[i] == 0 {
			v1b[i] = 0
		} else {
			v1b[i] = v1[i] - 1
			break
		}
	}
	id1a := NewNodeID(v1a)
	id1b := NewNodeID(v1b)

	result, err := id1.Compare(id1) // self
	c.Assert(0, Equals, result)
	c.Assert(err, IsNil)

	result, err = id1.Compare(id1.Clone()) // identical copy
	c.Assert(0, Equals, result)
	c.Assert(err, IsNil)

	result, err = id1.Compare(nil) // nil comparand
	c.Assert(err, Not(IsNil))      // NOT

	result, err = id1.Compare(id3)
	c.Assert(err, Not(IsNil)) // different lengths	// NOT

	result, err = id1.Compare(id1a)
	c.Assert(-1, Equals, result) // id1 < id1a
	c.Assert(err, IsNil)

	result, err = id1.Compare(id1b) // id1 > id1b
	c.Assert(1, Equals, result)
	c.Assert(err, IsNil)

	result, err = id1a.Compare(id1b) // id1a > id1b
	c.Assert(1, Equals, result)
	c.Assert(err, IsNil)

}
