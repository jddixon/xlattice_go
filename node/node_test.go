package node

import (
	cr "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	"github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
	"strings"
	"time"
)

const (
	VERBOSITY = 1
)

func makeNodeID(rng *rnglib.PRNG) (*NodeID, error) {
	var buffer []byte
	// quasi-random choice, whether to use an SHA1 or SHA3 nodeID
	if rng.NextBoolean() {
		buffer = make([]byte, SHA1_LEN)
	} else {
		buffer = make([]byte, SHA3_LEN)
	}
	rng.NextBytes(&buffer)
	return NewNodeID(buffer)
}

// func doKeyTests(t *testing.T, node *Node, rng *SimpleRNG) {
func (s *XLSuite) doKeyTests(c *C, node *Node, rng *rnglib.PRNG) {
	// COMMS KEY
	commsPubKey := node.GetCommsPublicKey()
	c.Assert(commsPubKey, Not(IsNil)) // NOT

	privCommsKey := node.commsKey // naughty
	c.Assert(privCommsKey.Validate(), IsNil)

	expLen := (*privCommsKey.D).BitLen()
	if VERBOSITY > 0 {
		fmt.Printf("bit length of private key exponent is %d\n", expLen)
	}
	// 2037 seen at least once
	c.Assert(true, Equals, (2036 <= expLen) && (expLen <= 2048))

	c.Assert(privCommsKey.PublicKey, Equals, *commsPubKey) // XXX FAILS

	// SIG KEY
	sigPubKey := node.GetSigPublicKey()
	c.Assert(sigPubKey, Not(IsNil)) // NOT

	privSigKey := node.sigKey // naughty
	c.Assert(privSigKey.Validate(), IsNil)

	expLen = (*privSigKey.D).BitLen()
	if VERBOSITY > 0 {
		fmt.Printf("bit length of private key exponent is %d\n", expLen)
	}
	// lowest value seen as of 2013-07-16 was 2039
	c.Assert(true, Equals, (2038 <= expLen) && (expLen <= 2048))

	c.Assert(privSigKey.PublicKey, Equals, *sigPubKey) // FOO

	// sign /////////////////////////////////////////////////////////
	msgLen := 128
	msg := make([]byte, msgLen)
	rng.NextBytes(&msg)

	d := sha1.New()
	d.Write(msg)
	hash := d.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, node.sigKey, cr.SHA1, hash)
	c.Assert(err, IsNil)

	signer := node.getSigner()
	signer.Update(msg)
	sig2, err := signer.Sign() // XXX change interface to allow arg

	lenSig := len(sig)
	lenSig2 := len(sig2)
	c.Assert(lenSig, Equals, lenSig2)

	// XXX why does this succeed?
	for i := 0; i < lenSig; i++ {
		c.Assert(sig[i], Equals, sig2[i])
	}

	// verify ///////////////////////////////////////////////////////
	err = rsa.VerifyPKCS1v15(sigPubKey, cr.SHA1, hash, sig)
	c.Assert(err, IsNil)

	// 2013-06-15, SigVerify now returns error, so nil means OK
	c.Assert(nil, Equals, xc.SigVerify(sigPubKey, msg, sig))

	s.nilArgCheck(c)
}

// XXX TODO: move these tests into crypto/sig_test.go
// func nilArgCheck(t *testing.T) {
func (s *XLSuite) nilArgCheck(c *C) {
	// the next statement should always return an error
	err := xc.SigVerify(nil, nil, nil)
	c.Assert(nil, Not(Equals), err)
}

// END OF TODO

func (s *XLSuite) TestNewConstructor(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_NEW_CONSTRUCTOR")
	}
	// if  constructor assigns a nil NodeID, we should get an
	// IllegalArgument panic
	// XXX STUB

	// if assigned a nil key, the New constructor should panic
	// with an IllegalArgument string
	// XXX STUB

}

func (s *XLSuite) shouldCreateTcpEndPoint(c *C, addr string) *xt.TcpEndPoint {
	ep, err := xt.NewTcpEndPoint(addr)
	c.Assert(err, Equals, nil)
	c.Assert(ep, Not(Equals), nil)
	return ep
}
func (s *XLSuite) TestAutoCreateOverlays(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_AUTO_CREATE_OVERLAYS")
	}
	rng := rnglib.MakeSimpleRNG()
	name := rng.NextFileName(4)
	id, err := makeNodeID(rng)
	c.Assert(err, Equals, nil)
	c.Assert(id, Not(IsNil))

	ep0 := s.shouldCreateTcpEndPoint(c, "127.0.0.0:0")
	ep1 := s.shouldCreateTcpEndPoint(c, "127.0.0.0:0")
	ep2 := s.shouldCreateTcpEndPoint(c, "127.0.0.0:0")
	e := []xt.EndPointI{ep0, ep1, ep2}

	node, err := New(name, id, "", nil, nil, nil, e, nil)
	c.Assert(err, Equals, nil)
	c.Assert(node, Not(Equals), nil)
	defer node.Close()

	c.Assert(node.SizeEndPoints(), Equals, len(e))
	c.Assert(node.SizeOverlays(), Equals, 1)

	// expect to find an acceptor for each endpoint
	// XXX STUB XXX

	// Close must close all three acceptors
	// XXX STUB XXX

}

// Return an initialized and tested host, with a NodeID, commsKey,
// and sigKey
func (s *XLSuite) makeHost(c *C, rng *rnglib.PRNG) *Node {
	// XXX names may not be unique
	name := rng.NextFileName(6)
	for {
		first := string(name[0])
		if !strings.Contains(first, "0123456789") &&
			!strings.Contains(name, "-") {
			break
		}
		name = rng.NextFileName(6)
	}
	id, err := makeNodeID(rng)
	c.Assert(err, Equals, nil)
	c.Assert(id, Not(IsNil))

	n, err2 := NewNew(name, id)
	c.Assert(n, Not(IsNil))
	c.Assert(err2, IsNil)
	c.Assert(name, Equals, n.GetName())
	actualID := n.GetNodeID()
	c.Assert(true, Equals, id.Equal(actualID))
	s.doKeyTests(c, n, rng)
	c.Assert(0, Equals, (*n).SizePeers())
	c.Assert(0, Equals, (*n).SizeOverlays())
	c.Assert(0, Equals, n.SizeConnections())
	c.Assert("", Equals, n.GetLFS())
	return n
}

// Create a Peer from information in the Node passed.  Endpoints
// (and so Overlays) must have already been added to the Node.
func (s *XLSuite) peerFromHost(c *C, n *Node) (peer *Peer) {
	var err error
	k := len(n.endPoints)
	ctors := make([]xt.ConnectorI, k)
	for i := 0; i < k; i++ {
		ctors[i], err = xt.NewTcpConnector(n.GetEndPoint(i))
		c.Assert(err, Equals, nil)
	}
	peer = &Peer{ctors, n.BaseNode}
	//peer.commsPubKey =  n.GetCommsPublicKey()
	//peer.sigPubKey	 =  n.GetSigPublicKey()

	return peer
}
func (s *XLSuite) makeAnEndPoint(c *C, rng *rnglib.PRNG, node *Node) {
	// XXX these may not be unique
	port := rng.Intn(256 * 256)
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	ep, err := xt.NewTcpEndPoint(addr)
	c.Assert(err, IsNil)
	c.Assert(ep, Not(IsNil))
	ndx, err := node.AddEndPoint(ep)
	c.Assert(err, IsNil)
	c.Assert(ndx, Equals, 0) // it's the only one
}
func (s *XLSuite) TestNodeSerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_NODE_SERIALIZATION")
	}
	rng := rnglib.MakeSimpleRNG()

	node := s.makeHost(c, rng)
	s.makeAnEndPoint(c, rng, node)
	lfs := rng.NextFileName(4)
	err := node.setLFS(lfs)
	c.Assert(err, IsNil)
	c.Assert(node.GetLFS(), Equals, lfs)

	const K = 3
	peers := make([]*Peer, K)

	for i := 0; i < K; i++ {
		host := s.makeHost(c, rng)
		s.makeAnEndPoint(c, rng, host)
		peers[i] = s.peerFromHost(c, host)
		ndx, err := node.AddPeer(peers[i])
		c.Assert(err, IsNil)
		c.Assert(ndx, Equals, i)
	}
	// we now have a node with K peers
	serialized := node.String()
	fmt.Printf("SERIALIZED NODE WITH PEERS:\n%s\n", serialized)

	// we can't deserialize the node - it contains live acceptors!
	for i := 0; i < node.SizeAcceptors(); i++ {
		node.GetAcceptor(i).Close()
	}
	// XXX parse succeeds if we sleep 100ms, fails if we sleep 10ms
	time.Sleep(50 * time.Millisecond)

	backAgain, rest, err := Parse(serialized)
	c.Assert(err, IsNil)
	c.Assert(len(rest), Equals, 0)

	reserialized := backAgain.String()
	c.Assert(reserialized, Equals, serialized)
}
