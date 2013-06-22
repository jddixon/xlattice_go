package overlay

// xlattice_go/addr_range/addr_range_test.go

import (
	x "github.com/jddixon/xlattice_go"
	. "launchpad.net/gocheck"
)


func (s *XLSuite) TestAddrRangeCtor(c *C) {
	rng			:= x.MakeSimpleRNG()
	
	// v4 address ---------------------------------------------------
	v4PLen		:= uint(1 + rng.Intn(32))				// in bits
	v4ByteLen	:= ((v4PLen + 7) / 8) * 8
	pBuffer		:= make([]byte, v4ByteLen)
	p, err		:= NewV4AddrRange(pBuffer, v4PLen)
	c.Assert(err,	IsNil)
	c.Assert(p,		Not(IsNil))

	//        actual          vs    expected
	c.Assert(p.PrefixLen(), Equals, v4PLen)
	c.Assert(p.AddrLen(),	Equals, uint(32))

	// very weak tests of Equal()
	c.Assert(p.Equal(p),	Equals, true)
	c.Assert(p.Equal(nil),	Equals, false)

	// a better implementation would truncate the prefix to the right
	// number of bits; a better test would test whether the truncation
	// is done correctly

	// v6 address ---------------------------------------------------
	v6PLen		:= uint(1 + rng.Intn(64))				// in bits
	v6ByteLen	:= ((v6PLen + 7) / 8) * 8
	p6Buffer	:= make([]byte, v6ByteLen)
	p6, err		:= NewV6AddrRange(p6Buffer, v6PLen)
	c.Assert(err,	IsNil)
	c.Assert(p,		Not(IsNil))

	//        actual          vs    expected
	c.Assert(p6.PrefixLen(), Equals, v6PLen)
	c.Assert(p6.AddrLen(),	Equals, uint(64))

	// very weak tests of Equal()
	c.Assert(p6.Equal(p6),	Equals, true)
	c.Assert(p6.Equal(nil),	Equals, false)

	// v4 vs v6 -----------------------------------------------------
	c.Assert(p6.Equal(p),	Equals, false)

}
