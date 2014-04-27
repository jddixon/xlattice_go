package reg

// xlattice_go/reg/xlreg_const_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestErrorConst(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_ERROR_CONSTS")
	}
	c.Assert(BAD_ATTRS_LINE, Equals, -1)
	c.Assert(BAD_VERSION, Equals, -2)
	c.Assert(CANT_FIND_CLUSTER_BY_ID, Equals, -3)
}
