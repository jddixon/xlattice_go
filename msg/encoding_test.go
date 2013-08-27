package msg

// xlattice_go/msg/encoding_test.go

import (
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestEncoding(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_ENCODING")
	}
	rng := rnglib.MakeSimpleRNG()

	const K = 2
	nodes, accs := xn.MockLocalHostCluster(K) // lazy way to get keys
	defer func() {
		for i := 0; i < K; i++ {
			if accs[i] != nil {
				accs[i].Close()
			}
		}
	}()

	_ = rng // WORKING HERE

	// XXX NEED TO TEST EVERY FIELD

	peer := nodes[0].GetPeer(0)
	pID := peer.GetNodeID().Value()
	pck, err := xc.RSAPubKeyToWire(peer.GetCommsPublicKey())
	c.Assert(err, IsNil)
	psk, err := xc.RSAPubKeyToWire(peer.GetSigPublicKey())
	c.Assert(err, IsNil)

	cmd := xn.XLatticeMsg_Hello
	one := uint64(1)
	msg := &xn.XLatticeMsg{
		Op:       &cmd,
		MsgN:     &one,
		ID:       pID,
		CommsKey: pck,
		SigKey:   psk,
	}
	wired, err := EncodePacket(msg)
	c.Assert(err, IsNil)
	c.Assert(wired, Not(IsNil))

	backAgain, err := DecodePacket(wired)
	c.Assert(err, IsNil)
	c.Assert(backAgain, Not(IsNil))

	rewired, err := EncodePacket(msg)
	c.Assert(err, IsNil)
	c.Assert(rewired, Not(IsNil))

	c.Assert(len(wired), Equals, len(rewired))
	for i := 0; i < len(wired); i++ {
		c.Assert(wired[i], Equals, rewired[i])
	}
}
