package msg

// xlattice_go/msg/in_q_test.go

import (
	"fmt"
	"github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
	"time"
)

var _ = fmt.Print

const (
	VERBOSITY = 1
	SHA1_LEN  = 20
)

func (s *XLSuite) makeBadGuy(c *C) (badGuy *node.Node, acc xt.AcceptorI) {
	rng := rnglib.MakeSimpleRNG()
	id := make([]byte, SHA1_LEN)
	rng.NextBytes(&id)
	nodeID, err := xi.NewNodeID(id)
	c.Assert(err, IsNil)
	name := rng.NextFileName(8)
	badGuy, err = node.NewNew(name, nodeID)
	c.Assert(err, IsNil)
	accCount := badGuy.SizeAcceptors()
	c.Assert(accCount, Equals, 0)
	ep, err := xt.NewTcpEndPoint("127.0.0.1:0")
	c.Assert(err, IsNil)
	ndx, err := badGuy.AddEndPoint(ep)
	c.Assert(err, IsNil)
	c.Assert(ndx, Equals, 0)
	acc = badGuy.GetAcceptor(0)
	return
}

// HELLO --------------------------------------------------------------
// If we receive a hello on a connection but do not know recognize the
// nodeID we just drop the connection.  We only deal with known peers.
// If either the crypto public key or sig public key is wrong, we send
// an error message and close the connection.  If the nodeID, cKey, and
// sKey are correct, we advance the handler's state to HELLO_RCVD

// XXX We should probably also require that msgN be 1.

func (s *XLSuite) TestHelloHandler(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_HELLO_HANDLER")
	}

	const TWO = 2

	// Create a node and add a mock peer.  This is a cluster of TWO.
	nodes, accs := node.MockLocalHostCluster(TWO)
	defer func() {
		for i := 0; i < TWO; i++ {
			if accs[i] != nil {
				accs[i].Close()
			}
		}
	}()
	myNode, peerNode := nodes[0], nodes[1]
	myAcc, peerAcc := accs[0], accs[1]

	c.Assert(myAcc, Not(IsNil))
	myAccEP := myAcc.GetEndPoint()
	myCtor, err := xt.NewTcpConnector(myAccEP)
	c.Assert(err, IsNil)

	// myNode's server side

	fmt.Println("STARTING SERVER")
	go func() {
		for {
			cnx, err := myAcc.Accept()
			c.Assert(err, IsNil)
			fmt.Printf("CONNECTION\n")

			go func() {
				_, _ = NewInHandler(myNode, cnx)
			}()
		}

	}()

	// Create a second mock peer unknown to myNode.
	badGuy, badAcc := s.makeBadGuy(c)
	defer badAcc.Close()
	badHello, err := MakeHelloMsg(badGuy)
	c.Assert(err, IsNil)
	c.Assert(badHello, Not(IsNil))

	_, _, _, _, _ = badGuy, myNode, peerNode, myAcc, peerAcc

	time.Sleep(100 * time.Millisecond)

	// Unknown peer sends Hello.  Test node should just drop the
	// connection.  It is an error if we receive a reply.

	conn, err := myCtor.Connect(xt.ANY_TCP_END_POINT)
	c.Assert(err, IsNil)
	c.Assert(conn, Not(IsNil))
	cnx := conn.(*xt.TcpConnection)

	data, err := EncodePacket(badHello)
	c.Assert(err, IsNil)
	c.Assert(data, Not(IsNil))
	fmt.Println("SENDING BADHELLO")
	count, err := cnx.Write(data)
	fmt.Println("BADHELLO SENT")

	c.Assert(err, IsNil)
	c.Assert(count, Equals, len(data))
	time.Sleep(100 * time.Millisecond)

	// XXX THIS TEST FAILS because of a deficiency in
	// transport/tcp_connection.GetState() - it does not look at
	// the state of the underlying connection
	// c.Assert(cnx.GetState(), Equals, xt.DISCONNECTED)

	// Known peer sends Hello with at least one of cKey or sKey wrong.
	// We expect to receive an error msg and then the connection
	// should be closed.

	// XXX STUB XXX

	// Known peer sends Hello with all parameters correct.  We reply
	// with an Ack and advance state to open.

	peerHello, err := MakeHelloMsg(peerNode)
	c.Assert(err, IsNil)
	c.Assert(peerHello, Not(IsNil))

	conn, err = myCtor.Connect(xt.ANY_TCP_END_POINT)
	c.Assert(err, IsNil)
	c.Assert(conn, Not(IsNil))
	cnx = conn.(*xt.TcpConnection)

	data, err = EncodePacket(peerHello)
	c.Assert(err, IsNil)
	c.Assert(data, Not(IsNil))
	fmt.Println("SENDING PEERHELLO")
	count, err = cnx.Write(data)
	fmt.Println("PEER_HELLO SENT")

	c.Assert(err, IsNil)
	c.Assert(count, Equals, len(data))
	time.Sleep(100 * time.Millisecond)

	// WORKING HERE

	fmt.Println("should be waiting for ack in reply to hello")

	// XXX STUB XXX

	// wait for ack
	// XXX STUB XXX

	// verify msg returned is an ack and has the correct parameters
	// XXX STUB XXX

	// send bye
	// XXX STUB XXX

	// wait for ack
	// XXX STUB XXX

	// Clean up: close the connection.
	// XXX STUB XXX

}
