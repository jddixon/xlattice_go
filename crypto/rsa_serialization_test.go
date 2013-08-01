package crypto

// xlattice_go/crypto/rsa_test.go

import (
	cr "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	//"code.google.com/p/go.crypto/ssh"
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestRSAPubKeyToFromDisk(c *C) {
	rng := rnglib.MakeSimpleRNG()

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	c.Assert(err, Equals, nil)
	c.Assert(key, Not(Equals), nil)

	pubKey := key.PublicKey
	c.Assert(pubKey, Not(Equals), nil)

	// the public key in SSH disk format
	sshAuthKey, ok := RSAPubKeyToDisk(&pubKey)
	c.Assert(ok, Equals, true)

	// generate the public key from the serialized version
	pk2, err := RSAPubKeyFromDisk(sshAuthKey)
	c.Assert(err, Equals, nil)
	c.Assert(pk2, Not(Equals), nil)

	// a long-winded way of proving pk2 == pubKey: we create a
	// random message, sign it with the private key and verify it with pk2

	msgLen := 128
	msg := make([]byte, msgLen)
	rng.NextBytes(&msg)

	digest := sha1.New()
	digest.Write(msg)
	hash := digest.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, key, cr.SHA1, hash)
	c.Assert(err, IsNil)

	err = rsa.VerifyPKCS1v15(pk2, cr.SHA1, hash, sig)
	c.Assert(err, IsNil)
}
