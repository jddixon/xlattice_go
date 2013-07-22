package node

// xlattice_go/node/localHost_test.go

import (
	. "launchpad.net/gocheck"
	"github.com/jddixon/xlattice_go/rnglib"
)

// See cluster_test.go for a general description of these tests.  
//
// This test involves nodes executing on a single machine, with accessor
// IP addresses 127.0.0.1:P, where P represents a system-assigned unique 
// port number.

func (s *XLSuite) TestLocalHostCluster(c *C) {
	rng := rnglib.MakeSimpleRNG()
	_ = rng

	// Create N nodes, each with a NodeID, two RSA private keys (sig and
	// comms), and two RSA public keys.  Each node creates a TcpAcceptor
	// running on 127.0.0.1 and a random (= system-supplied) port.
	// XXX STUB XXX

	// Collect the nodeID, public keys, and listening address from each
	// node and use this information to configure the []Peer data structure
	// on each node.
	// XXX STUB XXX

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
