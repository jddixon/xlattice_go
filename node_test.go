package xlattice_go

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"github.com/bmizerany/assert"
	. "github.com/jddixon/xlattice_go/rnglib"
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
	
	privkey := node.key		// naughty
	assert.Equal( t, nil, privkey.Validate() )

	expLen := (*privkey.D).BitLen()
	fmt.Printf("bit length of private key exponent is %d\n", expLen) // DEBUG
	assert.Equal(t, true, (2040 <= expLen) && (expLen <= 2048))

	assert.Equal(t, privkey.PublicKey, *pubkey)

	// sign /////////////////////////////////////////////////////////
	msgLen := 128
	msg	   := make([]byte, msgLen)
	rng.NextBytes(&msg)

	d := sha1.New()
	d.Write(msg)
	hash := d.Sum(nil)
	
	sig, err := rsa.SignPKCS1v15(rand.Reader, node.key, crypto.SHA1, hash)
	assert.Equal(t, nil, err)

	// verify ///////////////////////////////////////////////////////
	err = rsa.VerifyPKCS1v15(pubkey, crypto.SHA1, hash, sig)
	assert.Equal(t, nil, err)
}
func TestNewNew(t *testing.T) {
	rng := makeRNG()
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
	// rng := makeRNG()

	// if  constructor assigns a nil NodeID, we should get an
	// IllegalArgument panic
	// XXX STUB

	// if assigned a nil key, the NewNode constructor should panic
	// with an IllegalArgument string
	// XXX STUB

}
