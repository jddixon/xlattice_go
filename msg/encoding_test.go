package msg

// xlattice_go/msg/encoding_test.go

import (
	"fmt"
	xc "github.com/jddixon/xlCrypto_go"
	xn "github.com/jddixon/xlNode_go"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestEncoding(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_ENCODING")
	}
	const K = 2
	nodes, accs := xn.MockLocalHostCluster(K) // lazy way to get keys
	defer func() {
		for i := 0; i < K; i++ {
			if accs[i] != nil {
				accs[i].Close()
			}
		}
	}()

	peer := nodes[0].GetPeer(0)
	pID := peer.GetNodeID().Value()
	pck, err := xc.RSAPubKeyToWire(peer.GetCommsPublicKey())
	c.Assert(err, IsNil)
	psk, err := xc.RSAPubKeyToWire(peer.GetSigPublicKey())
	c.Assert(err, IsNil)

	cmd := XLatticeMsg_Hello
	one := uint64(1)
	msg := &XLatticeMsg{
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

	// DEBUG
	//fmt.Printf("    len ck %d bytes\n", len(msg.GetCommsKey()))	// 294
	//fmt.Printf("    len sk %d bytes\n", len(msg.GetSigKey()))		// 294
	//fmt.Println("    end TestEncoding")
	// END
}
