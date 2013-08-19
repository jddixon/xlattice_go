package node

// xlattice_go/node/mock_localHost_cluster_test.go

import (
	"crypto/rsa"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print
var _ = xo.NewIPOverlay

// XXX ROUGH CODE, NEEDS REVIEW.  A test was split into two parts.
// One became MockLocalHostCluster() and this is the other part.

func (s *XLSuite) TestMockLocalHostTcpCluster(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MOCK_LOCAL_HOST_TCP_CLUSTER")
	}
	var err error
	const K = 5
	nodes, accs := MockLocalHostCluster(K)

	defer func() {
		for i := 0; i < K; i++ {
			if accs[i] != nil {
				accs[i].Close()
			}
		}
	}()

	for i := 0; i < K; i++ {
		c.Assert(nodes, Not(IsNil))
		c.Assert(accs, Not(IsNil))
	}
	nameSet := make(map[string]bool)
	names := make([]string, K)
	nodeIDs := make([]*xi.NodeID, K)
	for i := 0; i < K; i++ {
		names[i] = nodes[i].GetName()
		_, ok := nameSet[names[i]]
		c.Assert(ok, Equals, false)
		nameSet[names[i]] = true

		// XXX should also verify nodeIDs are unique
		nodeIDs[i] = nodes[i].GetNodeID()
	}
	ar, err := xo.NewCIDRAddrRange("127.0.0.0/8")
	c.Assert(err, Equals, nil)
	overlay, err := xo.NewIPOverlay("XO", ar, "tcp", 1.0)
	c.Assert(err, Equals, nil)

	_ = overlay

	accEndPoints := make([]*xt.TcpEndPoint, K)
	for i := 0; i < K; i++ {
		accEndPoints[i] = accs[i].GetEndPoint().(*xt.TcpEndPoint)
		c.Assert(accEndPoints[i], Not(IsNil))

		c.Assert(nodes[i].SizeEndPoints(), Equals, 1)
		c.Assert(nodes[i].SizeAcceptors(), Equals, 1)
		c.Assert(nodes[i].SizeOverlays(), Equals, 1)
		c.Assert(overlay.Equal(nodes[i].GetOverlay(0)), Equals, true)
	}

	// XXX NEEDS CHECKING FROM HERE

	commsKeys := make([]*rsa.PublicKey, K)
	sigKeys := make([]*rsa.PublicKey, K)
	ctors := make([]*xt.TcpConnector, K)

	for i := 0; i < K; i++ {
		commsKeys[i] = nodes[i].GetCommsPublicKey()
		sigKeys[i] = nodes[i].GetSigPublicKey()
		ctors[i], err = xt.NewTcpConnector(accEndPoints[i])
		c.Assert(err, Equals, nil)
	}

	//overlaySlice := []xo.OverlayI{overlay}
	// peers := make([]*Peer, K)
	for i := 0; i < K; i++ {
		//ctorSlice := []xt.ConnectorI{ctors[i]}
		//_ = ctorSlice
		//peers[i], err = NewPeer(names[i], nodeIDs[i], commsKeys[i], sigKeys[i],
		//	overlaySlice, ctorSlice)
		//c.Assert(err, Equals, nil)
	}

	// Use the information collected to configure each node.
	for i := 0; i < K; i++ {
		//for j := 0; j < K; j++ {
		//	if i != j {
		//		ndx, err := nodes[i].AddPeer(peers[j])
		//		c.Assert(err, Equals, nil)
		//		var expectedNdx int
		//		if j < i {
		//			expectedNdx = j
		//		} else {
		//			expectedNdx = j - 1
		//		}
		//		c.Assert(ndx, Equals, expectedNdx)
		//	}
		//} // GEEP
		c.Assert(nodes[i].SizeAcceptors(), Equals, 1)
		// XXX WRONG APPROACH - SizeConnectors() is a Peer
		// function, and in this case should return 1 for each peer.
		// c.Assert(nodes[i].SizeConnectors(),Equals, K-1)
		c.Assert(nodes[i].SizeEndPoints(), Equals, 1)
		c.Assert(nodes[i].SizeOverlays(), Equals, 1)
		c.Assert(nodes[i].SizePeers(), Equals, K-1)

	} // GEEP
}
