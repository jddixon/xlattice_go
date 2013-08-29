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

var (
	TWO   uint64 = 2
	THREE uint64 = 3
	FOUR  uint64 = 4
	FIVE  uint64 = 5
	SIX   uint64 = 6
)

func (s *XLSuite) makeANode(c *C) (badGuy *node.Node, acc xt.AcceptorI) {
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

// If we receive a hello on a connection but do not know recognize the
// nodeID we just drop the connection.  We only deal with known peers.
// If either the crypto public key or sig public key is wrong, we send
// an error message and close the connection.  If the nodeID, cKey, and
// sKey are correct, we advance the handler's state to HELLO_RCVD

func (s *XLSuite) TestHelloHandler(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_HELLO_HANDLER")
	}

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
	meAsPeer := peerNode.GetPeer(0)
	myAcc, peerAcc := accs[0], accs[1]
	_ = peerAcc // never used

	c.Assert(myAcc, Not(IsNil))
	myAccEP := myAcc.GetEndPoint()
	myCtor, err := xt.NewTcpConnector(myAccEP)
	c.Assert(err, IsNil)

	// myNode's server side
	go func() {
		for {
			cnx, err := myAcc.Accept()
			c.Assert(err, IsNil)

			// each connection handled by a separate goroutine
			go func() {
				_, _ = NewInHandler(myNode, cnx)
			}()
		}
	}()

	// -- WELL-FORMED HELLO -----------------------------------------
	// Known peer sends Hello with all parameters correct.  We reply
	// with an Ack and advance state to open.

	conn, err := myCtor.Connect(xt.ANY_TCP_END_POINT)
	c.Assert(err, IsNil)
	c.Assert(conn, Not(IsNil))
	cnx2 := conn.(*xt.TcpConnection)
	defer cnx2.Close()

	oh := &OutHandler{CnxHandler{Cnx: cnx2, Peer: meAsPeer}}

	// manually create and send a hello message -
	peerHello, err := MakeHelloMsg(peerNode)
	c.Assert(err, IsNil)
	c.Assert(peerHello, Not(IsNil))

	data, err := EncodePacket(peerHello)
	c.Assert(err, IsNil)
	c.Assert(data, Not(IsNil))
	count, err := cnx2.Write(data)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, len(data))
	oh.MsgN = ONE
	// end manual hello -------------------------

	time.Sleep(100 * time.Millisecond)

	// wait for ack
	ack, err := oh.readMsg()
	c.Assert(err, IsNil)
	c.Assert(ack, Not(IsNil))

	// verify msg returned is an ack and has the correct parameters
	c.Assert(ack.GetOp(), Equals, XLatticeMsg_Ack)
	c.Assert(ack.GetMsgN(), Equals, TWO)
	c.Assert(ack.GetYourMsgN(), Equals, ONE)

	// -- KEEP-ALIVE ------------------------------------------------
	cmd := XLatticeMsg_KeepAlive
	keepAlive := &XLatticeMsg{
		Op:   &cmd,
		MsgN: &THREE,
	}
	data, err = EncodePacket(keepAlive)
	c.Assert(err, IsNil)
	c.Assert(data, Not(IsNil))
	count, err = cnx2.Write(data)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, len(data))

	// Wait for ack.  In a better world we time out if an ack is not
	// received in some short period rather than blocking forever.
	ack, err = oh.readMsg()
	c.Assert(err, IsNil)
	c.Assert(ack, Not(IsNil))

	// verify msg returned is an ack and has the correct parameters
	c.Assert(ack.GetOp(), Equals, XLatticeMsg_Ack)
	c.Assert(ack.GetMsgN(), Equals, FOUR)
	c.Assert(ack.GetYourMsgN(), Equals, THREE)

	// -- BYE -------------------------------------------------------
	cmd = XLatticeMsg_Bye
	bye := &XLatticeMsg{
		Op:   &cmd,
		MsgN: &FIVE,
	}
	data, err = EncodePacket(bye)
	c.Assert(err, IsNil)
	c.Assert(data, Not(IsNil))
	count, err = cnx2.Write(data)
	c.Assert(err, IsNil)
	c.Assert(count, Equals, len(data))

	// Wait for ack.  In a better world we time out if an ack is not
	// received in some short period rather than blocking forever.
	ack, err = oh.readMsg()
	c.Assert(err, IsNil)
	c.Assert(ack, Not(IsNil))

	// verify msg returned is an ack and has the correct parameters
	c.Assert(ack.GetOp(), Equals, XLatticeMsg_Ack)
	c.Assert(ack.GetMsgN(), Equals, SIX)
	c.Assert(ack.GetYourMsgN(), Equals, FIVE)

	fmt.Println("bye acknowledged") // DEBUG
}
func (s *XLSuite) TestHelloFromStranger(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_HELLO_FROM_STRANGER")
	}
	myNode, myAcc := s.makeANode(c)
	defer myAcc.Close()

	c.Assert(myAcc, Not(IsNil))
	myAccEP := myAcc.GetEndPoint()
	myCtor, err := xt.NewTcpConnector(myAccEP)
	c.Assert(err, IsNil)

	// myNode's server side
	go func() {
		for {
			cnx, err := myAcc.Accept()
			c.Assert(err, IsNil)

			// each connection handled by a separate goroutine
			go func() {
				_, _ = NewInHandler(myNode, cnx)
			}()
		}
	}()

	// Create a second mock peer unknown to myNode.
	badGuy, badAcc := s.makeANode(c)
	defer badAcc.Close()
	badHello, err := MakeHelloMsg(badGuy)
	c.Assert(err, IsNil)
	c.Assert(badHello, Not(IsNil))

	time.Sleep(100 * time.Millisecond)

	// Unknown peer sends Hello.  Test node should just drop the
	// connection.  It is an error if we receive a reply.

	conn, err := myCtor.Connect(xt.ANY_TCP_END_POINT)
	c.Assert(err, IsNil)
	c.Assert(conn, Not(IsNil))
	cnx := conn.(*xt.TcpConnection)
	defer cnx.Close()

	data, err := EncodePacket(badHello)
	c.Assert(err, IsNil)
	c.Assert(data, Not(IsNil))
	count, err := cnx.Write(data)

	c.Assert(err, IsNil)
	c.Assert(count, Equals, len(data))

	fmt.Println("U")

	time.Sleep(100 * time.Millisecond)

	// XXX THIS TEST FAILS because of a deficiency in
	// transport/tcp_connection.GetState() - it does not look at
	// the state of the underlying connection
	// c.Assert(cnx.GetState(), Equals, xt.DISCONNECTED)	// GEEP

	fmt.Println("Z")
}

// -- ILL-FORMED HELLO ------------------------------------------
// Known peer sends Hello with at least one of cKey or sKey wrong.
// We expect to receive an error msg and then the connection
// should be closed.

// XXX STUB XXX

// -- SECOND WELL-FORMED HELLO ----------------------------------
// In this implementation, a second hello is an error and like all
// errors will cause the peer to close the connection.
// --------------------------------------------------------------

// send well-formed hello
// XXX STUB XXX

// wait for and check ack
// XXX STUB XXX

// send second well-formed hello
// XXX STUB XXX

// wait for and check error message
// XXX STUB XXX

// clean up: close the connection
// XXX STUB XXX				// GEEP
