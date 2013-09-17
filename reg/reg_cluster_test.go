package reg

// xlattice_go/msg/reg_cluster_test.go

import (
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	//xn "github.com/jddixon/xlattice_go/node"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestCluster(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER")
	}
	c.Assert(CLUSTER_DELETED, Equals, 1)
	c.Assert(FOO, Equals, 2)
	c.Assert(BAR, Equals, 4)
}
