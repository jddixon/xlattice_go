package reg

// xlattice_go/msg/in_handler_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
	"strings"
)

func (s *XLSuite) TestInHandler(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_IN_HANDLER")
	}

	parts := strings.Split(VERSION, ".")
	c.Assert(len(parts), Equals, 3)
	for i := 0; i < 3; i++ {
		if len(parts[i]) == 1 {
			parts[i] = "0" + parts[i]
		}
	}
	joinedParts := strings.Join(parts, "")
	paddedSV := fmt.Sprintf("%06d", serverVersion)
	c.Assert(paddedSV, Equals, joinedParts)

	// These are the tags that InHandler will accept from a client.

	c.Assert(op2tag(XLRegMsg_Client), Equals, MIN_TAG)

	c.Assert(op2tag(XLRegMsg_Client), Equals, 0)
	c.Assert(op2tag(XLRegMsg_Create), Equals, 1)
	c.Assert(op2tag(XLRegMsg_Join), Equals, 2)
	c.Assert(op2tag(XLRegMsg_Get), Equals, 3)
	c.Assert(op2tag(XLRegMsg_Bye), Equals, 4)

	c.Assert(op2tag(XLRegMsg_Bye), Equals, MAX_TAG)
}
