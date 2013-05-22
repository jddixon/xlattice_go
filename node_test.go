package xlattice_go

import (
	cr "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	. "github.com/jddixon/xlattice_go/rnglib"
	"strings"
	. "launchpad.net/gocheck"
	// "github.com/bmizerany/assert"
	// "testing"
)

func makeNodeID(rng *SimpleRNG) *NodeID {
	var buffer []byte
	if rng.NextBoolean() {
		buffer = make([]byte, SHA1_LEN)
	} else {
		buffer = make([]byte, SHA3_LEN)
	}
	rng.NextBytes(&buffer)
	return NewNodeID(buffer)
}

// func doKeyTests(t *testing.T, node *Node, rng *SimpleRNG) {
func (s *XLSuite) doKeyTests(c *C, node *Node, rng *SimpleRNG) {
	pubkey := node.GetPublicKey()
	c.Assert(pubkey, Not(IsNil))	// NOT

	privkey := node.key // naughty
	c.Assert(privkey.Validate(), IsNil)

	expLen := (*privkey.D).BitLen()
	fmt.Printf("bit length of private key exponent is %d\n", expLen) // DEBUG
	c.Assert(true, Equals, (2040 <= expLen) && (expLen <= 2048))

	c.Assert(privkey.PublicKey, Equals, *pubkey)

	// sign /////////////////////////////////////////////////////////
	msgLen := 128
	msg := make([]byte, msgLen)
	rng.NextBytes(&msg)

	d := sha1.New()
	d.Write(msg)
	hash := d.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, node.key, cr.SHA1, hash)
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
	err = rsa.VerifyPKCS1v15(pubkey, cr.SHA1, hash, sig)
	c.Assert(err, IsNil)

	c.Assert(true, Equals, xc.SigVerify(pubkey, msg, sig))

	s.nilArgCheck(c)
}

// XXX TODO: move these tests into crypto/sig_test.go
// func nilArgCheck(t *testing.T) {
func (s *XLSuite) nilArgCheck(c *C) {
	defer func() {
		r := recover()
		c.Assert(r, Not(IsNil))	// NOT
		str := r.(string)
		c.Assert(true, Equals, strings.HasPrefix(str, "IllegalArgument"))
	}()
	// the next statement should cause a panic
	_ = xc.SigVerify(nil, nil, nil)
	c.Assert(true, Equals, false, "you should never see this message")
}

// END OF TODO

// func TestNewNew(t *testing.T) {
func (s *XLSuite)TestNewNew(c *C) {
	rng := MakeRNG()
	_, err := NewNewNode(nil)
	c.Assert(err, Not(IsNil))	// NOT

	id := makeNodeID(rng)
	c.Assert(id, Not(IsNil))	// NOT
	n, err2 := NewNewNode(id)
	c.Assert(n, Not(IsNil))	// NOT
	c.Assert(err2, IsNil)
	actualID := n.GetNodeID()
	c.Assert(true, Equals, id.Equal(actualID))
	s.doKeyTests(c, n, rng)
	c.Assert(0, Equals, (*n).SizePeers())
	c.Assert(0, Equals, (*n).SizeOverlays())
	c.Assert(0, Equals, n.SizeConnections())
}

//func TestNewCtor(t *testing.T) {
func (s *XLSuite)TestNewCtor(c *C) {
	// rng := MakeRNG()

	// if  constructor assigns a nil NodeID, we should get an
	// IllegalArgument panic
	// XXX STUB

	// if assigned a nil key, the NewNode constructor should panic
	// with an IllegalArgument string
	// XXX STUB

}
