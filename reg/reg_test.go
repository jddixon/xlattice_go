package reg

// xlattice_go/msg/reg_test.go

import (
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	//xn "github.com/jddixon/xlattice_go/node"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestReg(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_REG")
	}
	c.Assert(EPHEMERAL, Equals, 1)
	c.Assert(FOO, Equals, 2)
	c.Assert(BAR, Equals, 4)
}
