package filters

// xlattice_go/crypto/filters/countingBloom_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

//  private CountingBloom     filter
//  private int n;          // number of strings in set
//  private int m;          // size of set expressed as a power of two
//  private int k;          // number of filters
//  private byte[][] keys

//  public void setUp () {
//      filter = null
//      m = 20;             // default
//      k = 8
//      keys = new byte[100][20]
//  }

func (s *XLSuite) TestEmptyCountingBloom(c *C) {

	m := uint(20)
	k := uint(8)
	filter, err := NewCountingBloom(m, k)
	c.Assert(err, IsNil)
	c.Assert(filter.Size(), Equals, uint(0))
	c.Assert(filter.Capacity(), Equals, uint(2<<(m-1)))
}

func (s *XLSuite) doTestCBInserts(c *C, m, k, numKey uint) {
	keys := make([][]byte, numKey) // ][20]
	// set up distinct keys
	for i := uint(0); i < numKey; i++ {
		keys[i] = make([]byte, 20)
		for j := uint(0); j < 20; j++ {
			keys[i][j] = byte(0xff & (i + j + 100))
		}
	}
	filter, err := NewCountingBloom(m, k)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)
	for i := uint(0); i < numKey; i++ {
		c.Assert(filter.Size(), Equals, i)
		isAMember, ks, err := filter.IsMember(keys[i])
		c.Assert(err, IsNil)
		c.Assert(ks, NotNil)
		c.Assert(isAMember, Equals, false)
		filter.Insert(keys[i])
	}
	for i := uint(0); i < numKey; i++ {
		isAMember, ks, err := filter.IsMember(keys[i])
		c.Assert(err, IsNil)
		c.Assert(ks, NotNil)
		c.Assert(isAMember, Equals, true)
	}
}

func (s *XLSuite) doTestCBRemovals(c *C, m, k, numKey uint) {
	keys := make([][]byte, numKey) // ][20]
	// set up distinct keys
	for i := uint(0); i < numKey; i++ {
		keys[i] = make([]byte, 20)
		for j := uint(0); j < 20; j++ {
			keys[i][j] = byte(0xff & (i + j + 100))
		}
	}
	filter, err := NewCountingBloom(m, k)
	c.Assert(err, IsNil)
	c.Assert(filter, NotNil)
	for i := uint(0); i < numKey; i++ {
		c.Assert(filter.Size(), Equals, i)
		isAMember, ks, err := filter.IsMember(keys[i])
		c.Assert(err, IsNil)
		c.Assert(ks, NotNil)
		c.Assert(isAMember, Equals, false)
		filter.Insert(keys[i])
	}
	for i := uint(0); i < numKey; i++ {
		isAMember, ks, err := filter.IsMember(keys[i])
		c.Assert(err, IsNil)
		c.Assert(ks, NotNil)
		c.Assert(isAMember, Equals, true)
	}
	for i := uint(0); i < numKey; i++ {
		filter.Remove(keys[i])
		isAMember, ks, err := filter.IsMember(keys[i])
		c.Assert(err, IsNil)
		c.Assert(ks, IsNil)
		c.Assert(isAMember, Equals, false)
	}
}

func (s *XLSuite) testInserts(c *C) {
	const m = uint(20)
	const k = uint(8)
	s.doTestCBInserts(c, m, k, 16)  // default values
	s.doTestCBInserts(c, 14, k, 16) // stride = 9
	s.doTestCBInserts(c, 13, k, 16) // stride = 8
	s.doTestCBInserts(c, 12, k, 16) // stride = 7
	s.doTestCBInserts(c, 12, 7, 16)
}

func (s *XLSuite) testRemovals(c *C) {
	const m = uint(20)
	const k = uint(8)
	s.doTestCBRemovals(c, m, k, 16)  // default values
	s.doTestCBRemovals(c, 14, k, 16) // stride = 9
	s.doTestCBRemovals(c, 13, k, 16) // stride = 8
	s.doTestCBRemovals(c, 12, k, 16) // stride = 7
	s.doTestCBRemovals(c, 12, 5, 16)
}
