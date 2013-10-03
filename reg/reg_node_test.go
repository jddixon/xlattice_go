package reg

// xlattice_go/msg/reg_node_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"time"
)

func (s *XLSuite) TestRegNode(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_REG_NODE")
	}
}

// TEST SERIALIZATION ///////////////////////////////////////////////
func (s *XLSuite) TestRegNodeSerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_REG_NODE_SERIALIZATION")
	}
	rng := xr.MakeSimpleRNG()

	node, ckPriv, skPriv := s.makeHostAndKeys(c, rng)

	// This assigns an endPoint in 127.0.0.1 to the node; it
	// also starts the corresponding acceptor listening.
	s.makeALocalEndPoint(c, node)

	// We now have a node with 0 peers and a live acceptor.

	regNode, err := NewRegNode(node, ckPriv, skPriv)

	serialized := regNode.String()

	// We can't deserialize the node - it contains a live acceptor
	// at the same endPoint.
	for i := 0; i < regNode.SizeAcceptors(); i++ {
		regNode.GetAcceptor(i).Close()
	}

	// the Node version of this fails if sleep is say 10ms
	time.Sleep(70 * time.Millisecond)

	backAgain, rest, err := ParseRegNode(serialized)

	// DEBUG
	if len(rest) > 0 {
		for i := 0; i < len(rest); i++ {
			fmt.Printf("REST: %s\n", rest[i])
		}
	}
	// END
	c.Assert(err, IsNil)
	c.Assert(len(rest), Equals, 0)

	reserialized := backAgain.String()
	c.Assert(reserialized, Equals, serialized)
}
