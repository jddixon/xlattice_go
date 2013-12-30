package util

import (
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var (
	VERBOSITY = 0
)

func (s *XLSuite) TestThisAndThat(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_THIS_AND_THAT")
	}
	rng := rnglib.MakeSimpleRNG()

	// confirm that function handles odd bytes correctly
	for n := 6; n < 27; n++ {
		buf := make([]byte, n)
		rng.NextBytes(buf)
		c.Assert(SameBytes(buf, buf), Equals, true)
	}

	// 2013-10-04: empty slices should be equal
	var x, y []byte
	c.Assert(SameBytes(x, y), Equals, true)
}
