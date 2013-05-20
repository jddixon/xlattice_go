package xlattice_go

import (
	cr "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"github.com/bmizerany/assert"
	xc "github.com/jddixon/xlattice_go/crypto"
	. "github.com/jddixon/xlattice_go/rnglib"
	"strings"
	"testing"
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

func doKeyTests(t *testing.T, node *Node, rng *SimpleRNG) {
	pubkey := node.GetPublicKey()
	assert.NotEqual(t, nil, pubkey)

	privkey := node.key // naughty
	assert.Equal(t, nil, privkey.Validate())

	expLen := (*privkey.D).BitLen()
	fmt.Printf("bit length of private key exponent is %d\n", expLen) // DEBUG
	assert.Equal(t, true, (2040 <= expLen) && (expLen <= 2048))

	assert.Equal(t, privkey.PublicKey, *pubkey)

	// sign /////////////////////////////////////////////////////////
	msgLen := 128
	msg := make([]byte, msgLen)
	rng.NextBytes(&msg)

	d := sha1.New()
	d.Write(msg)
	hash := d.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, node.key, cr.SHA1, hash)
	assert.Equal(t, nil, err)

	signer := node.getSigner()
	signer.Update(msg)
	sig2, err := signer.Sign() // XXX change interface to allow arg

	lenSig := len(sig)
	lenSig2 := len(sig2)
	assert.Equal(t, lenSig, lenSig2)

	// XXX why does this succeed?
	for i := 0; i < lenSig; i++ {
		assert.Equal(t, sig[i], sig2[i])
	}

	// verify ///////////////////////////////////////////////////////
	err = rsa.VerifyPKCS1v15(pubkey, cr.SHA1, hash, sig)
	assert.Equal(t, nil, err)

	assert.Equal(t, true, xc.SigVerify(pubkey, msg, sig))

	nilArgCheck(t)
}

// XXX TODO: move these tests into crypto/sig_test.go
func nilArgCheck(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotEqual(t, r, nil)
		str := r.(string)
		assert.Equal(t, true, strings.HasPrefix(str, "IllegalArgument"))
	}()
	// the next statement should cause a panic
	_ = xc.SigVerify(nil, nil, nil)
	assert.Equal(t, true, false, "you should never see this message")
}

// END OF TODO

func TestNewNew(t *testing.T) {
	rng := MakeRNG()
	_, err := NewNewNode(nil)
	assert.NotEqual(t, nil, err)

	id := makeNodeID(rng)
	assert.NotEqual(t, nil, id)
	n, err2 := NewNewNode(id)
	assert.NotEqual(t, nil, n)
	assert.Equal(t, nil, err2)
	actualID := n.GetNodeID()
	assert.Equal(t, true, id.Equal(actualID))
	doKeyTests(t, n, rng)
	assert.Equal(t, 0, (*n).SizePeers())
	assert.Equal(t, 0, (*n).SizeOverlays())
	assert.Equal(t, 0, n.SizeConnections())
}

func TestNewCtor(t *testing.T) {
	// rng := MakeRNG()

	// if  constructor assigns a nil NodeID, we should get an
	// IllegalArgument panic
	// XXX STUB

	// if assigned a nil key, the NewNode constructor should panic
	// with an IllegalArgument string
	// XXX STUB

}
