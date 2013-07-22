package node

// xlattice_go/node/localHost_test.go

import (
	"crypto/rsa"
	"github.com/jddixon/xlattice_go/rnglib"
	xo "github.com/jddixon/xlattice_go/overlay"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
)
var _ = xo.NewIPOverlay

// See cluster_test.go for a general description of these tests.  
//
// This test involves nodes executing on a single machine, with accessor
// IP addresses 127.0.0.1:P, where P represents a system-assigned unique 
// port number.

func (s *XLSuite) TestLocalHostCluster(c *C) {
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
	accs  := make([]*xt.TcpAcceptor, K)
	accEndPoints := make([]*xt.TcpEndPoint, K)
	for i:= 0; i < K; i++ {
		nodes[i], err = NewNew(nodeIDs[i])
		c.Assert(err, Equals, nil)
		accs[i],err = xt.NewTcpAcceptor("127.0.0.1:0")
		c.Assert(err, Equals, nil)
		accEndPoints[i] = accs[i].GetEndPoint()
	}
	defer func() { 
		for i := 0; i < K; i++ {
			accs[i].Close()
		}
	}()
	// XXX NO WAY TO ASSIGN ACCEPTORS TO NODES :-)

	// Collect the nodeID, public keys, and listening address from each
	// node.

	// XXX WORKING HERE:
	// XXX SIMPLIFICATION - all nodes on the same overlay for now
	// overlay,err := xo.NewIPOverlay("XO", "tcp", 1.0)
	c.Assert(err, Equals, nil)
	
	commsKeys	:= make([]*rsa.PublicKey, K)
	sigKeys		:= make([]*rsa.PublicKey, K)
	ctors		:= make([]*xt.TcpConnector, K)

	for i := 0; i < K; i++ {
		// we have nodeIDs
		commsKeys[i] = nodes[i].GetCommsPublicKey()
		sigKeys[i]   = nodes[i].GetSigPublicKey()
		ctors[i],err = xt.NewTcpConnector(accEndPoints[i])
		c.Assert(err, Equals, nil)
	}

	// overlaySlice := []*xo.OverlayI{overlay}				// AND HERE
	peers	:= make([]*Peer, K)
	for i := 0; i < K; i++ {
		ctorSlice    := []*xt.TcpConnector{ctors[i]}
		_ = ctorSlice	// AND HERE
		// peers[i],err = NewPeer(nodeIDs[i], commsKeys[i], sigKeys[i], 
		//			&overlaySlice, &ctorSlice)				// AND HERE
		c.Assert(err, Equals, nil)
	}

	// Use the information collected to configure each node.
	for i := 0; i < K; i++ {
		// err = nodes[i].addOverlay(overlay)	// all in one overlay AND HERE
		c.Assert(nodes[i].SizeOverlays(), Equals, 1)
		c.Assert(err, Equals, nil)
		for j := 0; j < K; j++ {
			if i != j {
				err = nodes[i].addPeer(peers[j])
				c.Assert(err, Equals, nil)
			}
		}
		c.Assert(nodes[i].SizePeers(), Equals, K - 1)
	}

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
}
