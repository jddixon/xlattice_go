package reg

// xlattice_go/reg/helloAndReply_test.go

import (
	"crypto/aes"
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

// // DO NOT DELETE THESE - at least not yet
// func (s *XLSuite) makeAnID(c *C, rng *xr.PRNG) (id []byte) {
// 	id = make([]byte, SHA3_LEN)
// 	rng.NextBytes(&id)
// 	return
// }
// func (s *XLSuite) makeANodeID(c *C, rng *xr.PRNG) (nodeID *xi.NodeID) {
// 	id := s.makeAnID(c, rng)
// 	nodeID, err := xi.New(id)
// 	c.Assert(err, IsNil)
// 	c.Assert(nodeID, Not(IsNil))
// 	return
// }
// func (s *XLSuite) makeAnRSAKey(c *C) (key *rsa.PrivateKey) {
// 	key, err := rsa.GenerateKey(rand.Reader, 2048)
// 	c.Assert(err, IsNil)
// 	c.Assert(key, Not(IsNil))
// 	return key
// } // FOO

func (s *XLSuite) TestHelloAndReply(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_HELLO_AND_REPLY")
	}
	rng := xr.MakeSimpleRNG()
	nodeID := s.makeANodeID(c, rng)

	ckPriv := s.makeAnRSAKey(c)
	skPriv := s.makeAnRSAKey(c)

	node, err := xn.New("foo", nodeID, "", ckPriv, skPriv, nil, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(node, Not(IsNil))

	ck := node.GetCommsPublicKey()

	version1 := uint32(rng.Int31n(255 * 255)) // in effect an unsigned short

	// == HELLO =====================================================
	// On the client side, create and marshal a hello message containing
	// AES iv1, key1, salt1 in addition to the client-proposed protocol
	// version.

	ciphertext, iv1, key1, salt1, err := ClientEncodeHello(version1, ck)
	c.Assert(err, IsNil)
	c.Assert(len(iv1), Equals, aes.BlockSize)
	c.Assert(len(key1), Equals, 2*aes.BlockSize)
	c.Assert(len(salt1), Equals, 8)

	// On the server side: ------------------------------------------
	// Decrypt the hello using the node's private comms key, unpack.
	iv1s, key1s, salt1s, version1s, err := ServerDecodeHello(ciphertext, ckPriv)
	c.Assert(err, IsNil)

	c.Assert(iv1s, DeepEquals, iv1)
	c.Assert(key1s, DeepEquals, key1)
	c.Assert(salt1s, DeepEquals, salt1)
	c.Assert(version1s, Equals, version1)

	// == HELLO REPLY ===============================================
	// On the server side create, marshal a reply containing iv2, key2, salt2,
	// salt1, version2
	version2 := version1 // server accepts client proposal
	iv2, key2, salt2, ciphertext, err := ServerEncodeHelloReply(
		iv1, key1, salt1, version2)
	c.Assert(err, IsNil)

	// On the client side: ------------------------------------------
	//     decrypt the reply using engine1b = iv1, key1

	iv2c, key2c, salt2c, salt1c, version2c, err := ClientDecodeHelloReply(
		ciphertext, iv1, key1)

	c.Assert(err, IsNil)

	c.Assert(iv2c, DeepEquals, iv2)
	c.Assert(key2c, DeepEquals, key2)
	c.Assert(salt2c, DeepEquals, salt2)
	c.Assert(salt1c, DeepEquals, salt1)
	c.Assert(version2c, Equals, version1)

}
