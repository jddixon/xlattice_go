package msg

// xlattice_go/msg/in_q_test.go

import (
//	cr "crypto"
//	"crypto/rand"
//	"crypto/rsa"
//	"crypto/sha1"
	"fmt"
//	xc "github.com/jddixon/xlattice_go/crypto"
	"github.com/jddixon/xlattice_go/rnglib"
	"github.com/jddixon/xlattice_go/node"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
	"runtime"
//	"strings"
//	"time"
)

const (
	VERBOSITY   = 1
    MY_MAX_PROC = 2                 // OK for Test
    SHA1_LEN    = 20
)

// XXX DROP THIS RSN
func (s *XLSuite) TestRuntime(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_RUN_TIME")
	}
	was := runtime.GOMAXPROCS(MY_MAX_PROC)
	fmt.Printf("GOMAXPROCS was %d, has been reset to %d\n", was, MY_MAX_PROC)
	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())
}

func (s *XLSuite) makeBadGuy(c *C) (acc xt.AcceptorI, badGuy *node.Node) {
    rng := rnglib.MakeSimpleRNG()
    id  := make([]byte, SHA1_LEN)
    rng.NextBytes(&id)
    nodeID, err := node.NewNodeID(id)
    c.Assert(err, IsNil)
    name        := rng.NextFileName(8)
    badGuy, err = node.NewNew(name, nodeID)
    c.Assert(err, IsNil)
    accCount    := badGuy.SizeAcceptors()
    c.Assert(accCount, Equals, 0)
    ep, err     := xt.NewTcpEndPoint("127.0.0.1:0")
    c.Assert(err, IsNil)
    ndx, err    := badGuy.AddEndPoint(ep)
    c.Assert(err, IsNil)
    c.Assert(ndx, Equals, 0)
    acc         = badGuy.GetAcceptor(0)
    return
}
// HELLO --------------------------------------------------------------
// If we receive a hello on a connection but do not know recognize the
// nodeID we just drop the connection.  We only deal with known peers.
// If either the crypto public key or sig public key is wrong, we send 
// an error message and close the connection.  If the nodeID, cKey, and 
// sKey are correct, we advance the handler's state to HDLR_OPEN

// XXX We should probably also require that msgN be 1.

func (s *XLSuite) TestHelloHandler(c *C) {

    // Create a node and add a mock peer.  This is a cluster of 2.
    nodes, accs := node.MockLocalHostCluster(2)
    defer func() { 
        for i := 0; i < 2; i++ {
            if accs[i] != nil {
                accs[i].Close()
            }
        }
    }()
    myNode, peerNode := nodes[0], nodes[1]
    myAcc,  peerAcc  := accs[0],  accs[1]

    // Create a second mock peer unknown to myNode.
    badGuy, badAcc := s.makeBadGuy(c)
    defer badAcc.Close()

    _,_,_,_,_ = badGuy, myNode, peerNode, myAcc, peerAcc

    // Second mock peer sends Hello.  Test node should just drop the
    // connection.  It is an error if we receive a reply.

    // initial state:   IN_START
    // final state:     IN_CLOSED
    // XXX STUB XXX

    // Known peer sends Hello with at least one of cKey or sKey wrong. 
    // We expect to receive an error msg and then the connection 
    // should be closed.
    
    // initial state:   IN_START
    // final state:     IN_CLOSED
    // XXX STUB XXX

    // Known peer sends Hello with all parameters correct.  We reply 
    // with an Ack and advance state to open.
    
    // initial state:   IN_START
    // final state:     IN_OPEN
    // XXX STUB XXX


    // Clean up: close the connection.
}
