package node

import (
	cr "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	// "strings"
	// "github.com/bmizerany/assert"
	// "testing"
)

func makeNodeID(rng *rnglib.PRNG) *NodeID {
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
	commsPubKey := node.GetCommsPublicKey()
	c.Assert(commsPubKey, Not(IsNil)) // NOT

	privCommsKey := node.commsKey // naughty
	c.Assert(privCommsKey.Validate(), IsNil)

	expLen := (*privCommsKey.D).BitLen()
	fmt.Printf("bit length of private key exponent is %d\n", expLen) // DEBUG
	c.Assert(true, Equals, (2038 <= expLen) && (expLen <= 2048))

	c.Assert(privCommsKey.PublicKey, Equals, *commsPubKey) // XXX FAILS

	sigPubKey := node.GetSigPublicKey()
	c.Assert(sigPubKey, Not(IsNil)) // NOT

	privSigKey := node.sigKey // naughty
	c.Assert(privSigKey.Validate(), IsNil)

	expLen = (*privSigKey.D).BitLen()
	fmt.Printf("bit length of private key exponent is %d\n", expLen) // DEBUG
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
} // GEEP

// XXX TODO: move these tests into crypto/sig_test.go
// func nilArgCheck(t *testing.T) {
func (s *XLSuite) nilArgCheck(c *C) {
	// the next statement should always return an error
	err := xc.SigVerify(nil, nil, nil)
	c.Assert(nil, Not(Equals), err)
}

// END OF TODO

// func TestNewNew(t *testing.T) {
func (s *XLSuite) TestNewNew(c *C) {
	rng := rnglib.MakeSimpleRNG()
	_, err := NewNew(nil)
	c.Assert(err, Not(IsNil)) // NOT

	id := makeNodeID(rng)
	c.Assert(id, Not(IsNil)) // NOT
	n, err2 := NewNew(id)
	c.Assert(n, Not(IsNil)) // NOT
	c.Assert(err2, IsNil)
	actualID := n.GetNodeID()
	c.Assert(true, Equals, id.Equal(actualID))
	s.doKeyTests(c, n, rng)
	c.Assert(0, Equals, (*n).SizePeers())
	c.Assert(0, Equals, (*n).SizeOverlays())
	c.Assert(0, Equals, n.SizeConnections())
}

//func TestNewCtor(t *testing.T) {
func (s *XLSuite) TestNewCtor(c *C) {
	// rng := rnglib.MakeSimpleRNG()

	// if  constructor assigns a nil NodeID, we should get an
	// IllegalArgument panic
	// XXX STUB

	// if assigned a nil key, the New constructor should panic
	// with an IllegalArgument string
	// XXX STUB

}
