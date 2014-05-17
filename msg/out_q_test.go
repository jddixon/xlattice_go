package msg

// xlattice_go/msg/out_q_test.go

import (
	"encoding/hex"
	"fmt"
	xc "github.com/jddixon/xlCrypto_go"
	xn "github.com/jddixon/xlNode_go"
	xi "github.com/jddixon/xlNodeID_go"
	xa "github.com/jddixon/xlProtocol_go/aes_cnx"
	xr "github.com/jddixon/rnglib_go"
	. "gopkg.in/check.v1"
)

var _ = fmt.Print

// HELLO --------------------------------------------------------------
// If we receive a hello on a connection but do not know recognize the
// nodeID we just drop the connection.  We only deal with known peers.
// If either the crypto public key or sig public key is wrong, we send
// an error message and close the connection.  If the nodeID, cKey, and
// sKey are correct, we advance the handler's state to HELLO_RCVD

// XXX We should probably also require that msgN be 1.GEEP

func (s *XLSuite) TestMakeHelloMsg(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MAKE_HELLO_MSG")
	}
	rng := xr.MakeSimpleRNG()
	id := make([]byte, xa.SHA1_LEN)
	rng.NextBytes(id)
	nodeID, err := xi.NewNodeID(id)
	c.Assert(err, IsNil)

	name := rng.NextFileName(8)
	lfs := "tmp/" + hex.EncodeToString(id)
	mrX, err := xn.NewNew(name, nodeID, lfs)
	c.Assert(err, IsNil)
	cPubKey := mrX.GetCommsPublicKey()
	c.Assert(cPubKey, Not(IsNil))
	sPubKey := mrX.GetSigPublicKey()
	c.Assert(sPubKey, Not(IsNil))

	// convert MrX's keys to wire form as byte slices
	wcPubKey, err := xc.RSAPubKeyToWire(cPubKey)
	c.Assert(err, IsNil)
	c.Assert(len(wcPubKey) > 0, Equals, true)
	wsPubKey, err := xc.RSAPubKeyToWire(sPubKey)
	c.Assert(err, IsNil)
	c.Assert(len(wsPubKey) > 0, Equals, true)
	c.Assert(wsPubKey, Not(IsNil))

	hello, err := MakeHelloMsg(mrX)
	c.Assert(err, IsNil)
	c.Assert(hello, Not(IsNil))

	// check NodeID
	idInMsg := hello.GetID() // a byte slice, not a NodeID
	c.Assert(id, DeepEquals, idInMsg)

	// these are byte slices
	mcPubKey := hello.GetCommsKey()
	msPubKey := hello.GetSigKey()

	c.Assert(len(mcPubKey), Equals, len(wcPubKey))
	c.Assert(len(msPubKey), Equals, len(wsPubKey)) // FAILS 0, 294

	c.Assert(wcPubKey, DeepEquals, mcPubKey)
	c.Assert(wsPubKey, DeepEquals, msPubKey)
}
