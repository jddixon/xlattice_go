package reg

// xlattice_go/msg/in_handler_test.go

import (
	"fmt"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestInHandler(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_IN_HANDLER")
	}

	// These are the tags that InHandler will accept from a client.

	c.Assert(op2tag(XLRegMsg_Client), Equals, MIN_TAG)

	c.Assert(op2tag(XLRegMsg_Client), Equals, 0)
	c.Assert(op2tag(XLRegMsg_Create), Equals, 1)
	c.Assert(op2tag(XLRegMsg_Join), Equals, 2)
	c.Assert(op2tag(XLRegMsg_GetCluster), Equals, 3)
	c.Assert(op2tag(XLRegMsg_Bye), Equals, 4)

	c.Assert(op2tag(XLRegMsg_Bye), Equals, MAX_TAG)
}
