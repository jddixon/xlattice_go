package node

// xlattice_go/node/localHost_test.go

import (
	"crypto/rsa"
	"fmt"
	xo "github.com/jddixon/xlattice_go/overlay"
	"github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print
var _ = xo.NewIPOverlay

// See cluster_test.go for a general description of these tests.
//
// This test involves nodes executing on a single machine, with accessor
// IP addresses 127.0.0.1:P, where P represents a system-assigned unique
// port number.

func (s *XLSuite) TestLocalHostTcpCluster(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_LOCAL_HOST_TCP_CLUSTER")
	}
	var err error
	const K = 5
	rng := rnglib.MakeSimpleRNG()

	// Create K nodes, each with a NodeID, two RSA private keys (sig and
	// comms), and two RSA public keys.  Each node creates a TcpAcceptor
	// running on 127.0.0.1 and a random (= system-supplied) port.
	nodeIDs := make([]*NodeID, K)
	for i := 0; i < K; i++ {
		val := make([]byte, SHA1_LEN)
		rng.NextBytes(&val)
		nodeIDs[i] = NewNodeID(val)
	}
	nodes := make([]*Node, K)
	accs := make([]*xt.TcpAcceptor, K)
	accEndPoints := make([]*xt.TcpEndPoint, K)
	for i := 0; i < K; i++ {
		nodes[i], err = NewNew(nodeIDs[i])
		c.Assert(err, Equals, nil)
	}
	defer func() {
		for i := 0; i < K; i++ {
			if accs[i] != nil {
				accs[i].Close()
			}
		}
	}()

	// Collect the nodeID, public keys, and listening address from each
	// node.

	// all nodes on the same overlay
	ar, err := xo.NewCIDRAddrRange("127.0.0.0/8")
	c.Assert(err, Equals, nil)
	overlay, err := xo.NewIPOverlay("XO", ar, "tcp", 1.0)
	c.Assert(err, Equals, nil)

	// add an endpoint to each node
	for i := 0; i < K; i++ {
		ep, err := xt.NewTcpEndPoint("127.0.0.1:0")
		c.Assert(err, Equals, nil)
		ndx, err := nodes[i].AddEndPoint(ep)
		c.Assert(err, Equals, nil)
		c.Assert(ndx, Equals, 0)
		endPoint := nodes[i].GetEndPoint(0).(*xt.TcpEndPoint)
		accs[i] = nodes[i].GetAcceptor(0).(*xt.TcpAcceptor)
		accEndPoints[i] = accs[i].GetEndPoint().(*xt.TcpEndPoint)
		myAccEnd := accEndPoints[i]
		c.Assert(endPoint.Equal(myAccEnd), Equals, true) // FAILS

		// adding the endPoint added an acceptor and an overlay
		c.Assert(nodes[i].SizeEndPoints(), Equals, 1)
		c.Assert(nodes[i].SizeAcceptors(), Equals, 1)
		c.Assert(nodes[i].SizeOverlays(), Equals, 1) // FAILS

		// XXX we should verify that each node has the same overlay
		// as calculated above

	}

	commsKeys := make([]*rsa.PublicKey, K)
	sigKeys := make([]*rsa.PublicKey, K)
	ctors := make([]*xt.TcpConnector, K)

	for i := 0; i < K; i++ {
		// we have nodeIDs
		commsKeys[i] = nodes[i].GetCommsPublicKey()
		sigKeys[i] = nodes[i].GetSigPublicKey()
		ctors[i], err = xt.NewTcpConnector(accEndPoints[i])
		c.Assert(err, Equals, nil)
	}

	overlaySlice := []xo.OverlayI{overlay}
	peers := make([]*Peer, K)
	for i := 0; i < K; i++ {
		ctorSlice := []xt.ConnectorI{ctors[i]}
		_ = ctorSlice
		peers[i], err = NewPeer(nodeIDs[i], commsKeys[i], sigKeys[i],
			overlaySlice, ctorSlice)
		c.Assert(err, Equals, nil)
	}

	// Use the information collected to configure each node.
	for i := 0; i < K; i++ {
		// This is not necessary, because the overlay should have
		// been auto-created by AddEndPoint()
		ndx, err := nodes[i].AddOverlay(overlay)
		c.Assert(err, Equals, nil)
		c.Assert(ndx, Equals, 0)
		// Despite our adding an overlay, the count hasn't changed.
		c.Assert(nodes[i].SizeOverlays(), Equals, 1)
		for j := 0; j < K; j++ {
			if i != j {
				ndx, err := nodes[i].AddPeer(peers[j])
				c.Assert(err, Equals, nil)
				var expectedNdx int
				if j < i {
					expectedNdx = j
				} else {
					expectedNdx = j - 1
				}
				c.Assert(ndx, Equals, expectedNdx)
			}
		}
		c.Assert(nodes[i].SizeAcceptors(), Equals, 1)
		// XXX NOT IMPLEMENTED !
		// c.Assert(nodes[i].SizeConnectors(),Equals, K-1)
		c.Assert(nodes[i].SizeEndPoints(), Equals, 1)
		c.Assert(nodes[i].SizeOverlays(), Equals, 1)
		c.Assert(nodes[i].SizePeers(), Equals, K-1)
	}
	// AT THIS POINT we have K nodes, each with K-1 peers.

	// Start each node running in a separate goroutine.
	// XXX STUB XXX

	// Each node will in a somewhat randomized fashion send N messages
	// to every other node, expecting to receive back from the peer a
	// digital signature for the message.  As each response = digital
	// signature comes back it is validated.  When all messages have
	// been validated, the node sends a 'done' message on a boolean
	// channel to the supervisor.
	// XXX STUB XXX

	// When all nodes have signaled that they are done, the supervisor
	// sends on stopCh, the stop command channel.
	// XXX STUB XXX

	// Each node will send a reply to the supervisor on stoppedCh.
	// and then terminate.
	// XXX STUB XXX

	// When the supervisor has received stopped signals from all nodes,
	// it summarize results and terminates.

	for i := 0; i < K; i++ {
		accs[i].Close()
	}
}
