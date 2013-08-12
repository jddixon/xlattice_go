package node

// xlattice_go/node/localHost_test.go

import (
	"encoding/hex"
	"fmt"
	xo "github.com/jddixon/xlattice_go/overlay"
	"github.com/jddixon/xlattice_go/rnglib"
	"io/ioutil"
	. "launchpad.net/gocheck"
	"os"
	"path"
)

var _ = fmt.Print
var _ = xo.NewIPOverlay

// See cluster_test.go for a general description of these tests.
//
// This test involves nodes executing on a single machine, with accessor
// IP addresses 127.0.0.1:P, where P represents a system-assigned unique
// port number.

// Accept connections from peers until a message is received on stopCh.
// For each message received from a peer, calculate its SHA3-256 hash,
// send that as a reply, and close the connection.  Send on stoppedCh
// when all replies have been sent.
func (s *XLSuite) nodeAsServer(c *C, node *Node, stopCh, stoppedCh chan bool) {
	// XXX STUB XXX
	// select
	//	cnx <- acceptor
	//		go
	//			read
	//			calculate hash
	//			send hash
	//			close cnx
	//	stopCh <-
	//		stoppedCh <- true
	//
}

// Send Q messages to each peer, expecting to receive an SHA3-256 hash
// back.  When all are received and verified, send on doneCh.

func (s *XLSuite) nodeAsClient(c *C, node *Node, Q int, doneCh chan bool) {
	// XXX STUB XXX
	//	for i := 0; i < Q; i++
	//		for j := 0; j < K; j++
	//			if j == myNdx
	//				continue
	//			open cnx to peer j
	//			make msg, random length, random content
	//			send msg
	//			calculate hash
	//			wait for reply
	//			complain if hash differs from reply
	//	doneCh <- true
}

func (s *XLSuite) TestLocalHostTcpCluster(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_LOCAL_HOST_TCP_CLUSTER")
	}
	var err error
	const K = 5
	rng := rnglib.MakeSimpleRNG()
	_ = rng // XXX

	nodes, accs := MockLocalHostCluster(K)
	defer func() {
		for i := 0; i < K; i++ {
			if accs[i] != nil {
				accs[i].Close()
			}
		}
	}()

	// AT THIS POINT we have K nodes, each with K-1 peers.
	// Save the configurations
	pathsToCfg := make([]string, K)
	for i := 0; i < K; i++ {
		hexNodeID := hex.EncodeToString(nodes[i].GetNodeID().Value())
		pathsToCfg[i] = path.Join("tmp", hexNodeID, ".xlattice")

		err = os.MkdirAll(pathsToCfg[i], 0755)
		c.Assert(err, IsNil)
		cfgFileName := path.Join(pathsToCfg[i], "config")

		fmt.Printf("WRITING CONFIG FILE %s\n", cfgFileName)

		cfg := nodes[i].String()
		err = ioutil.WriteFile(cfgFileName, []byte(cfg), 0644)
		c.Assert(err, IsNil)
	}

	// Start each node running in a separate goroutine.
	doneCh := make(chan (bool), K)
	stopCh := make(chan (bool), K)
	stoppedCh := make(chan (bool), K)
	_, _, _ = doneCh, stopCh, stoppedCh // DEBUG

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
	// XXX STUB XXX

}
