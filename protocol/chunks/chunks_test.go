package chunks

// xlattice_go/protocol/chunks/chunks_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

func (s *XLSuite) TestConstants(c *C) {
	c.Assert(MAGIC_OFFSET, Equals, 0)
	c.Assert(TYPE_OFFSET, Equals, 1)
	c.Assert(RESERVED_OFFSET, Equals, 2)
	c.Assert(LENGTH_OFFSET, Equals, 8)
	c.Assert(INDEX_OFFSET, Equals, 12)
	c.Assert(DATUM_OFFSET, Equals, 16)
	c.Assert(DATA_OFFSET, Equals, 48)
}

func (s *XLSuite) TestProperties(c *C) {
	rng := xr.MakeSimpleRNG()
	_ = rng

}

func (s *XLSuite) TestChunks(c *C) {
	rng := xr.MakeSimpleRNG()
	_ = rng

}
