package crypto

// xlattice_go/crypto/rsa_test.go

import (
	cr "crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	//. "launchpad.net/gocheck"
)

var _ = fmt.Print

func (s *XLSuite) TestRSAPubKeyToFromDisk(c *C) {
	rng := rnglib.MakeSimpleRNG()

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	c.Assert(err, Equals, nil)
	c.Assert(key, Not(Equals), nil)

	pubKey := key.PublicKey
	c.Assert(pubKey, Not(Equals), nil)

	// the public key in SSH disk format
	sshAuthKey, err := RSAPubKeyToDisk(&pubKey)
	c.Assert(err, Equals, nil)

	// generate the public key from the serialized version
	pk2, err := RSAPubKeyFromDisk(sshAuthKey)
	c.Assert(err, Equals, nil)
	c.Assert(pk2, Not(Equals), nil)

	// a long-winded way of proving pk2 == pubKey: we create a
	// random message, sign it with the private key and verify it with pk2

	msgLen := 128
	msg := make([]byte, msgLen)
	rng.NextBytes(msg)

	digest := sha1.New()
	digest.Write(msg)
	hash := digest.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, key, cr.SHA1, hash)
	c.Assert(err, IsNil)

	err = rsa.VerifyPKCS1v15(pk2, cr.SHA1, hash, sig)
	c.Assert(err, IsNil)
} // GEEP

func (s *XLSuite) TestRSAPrivateKeyToFromDisk(c *C) {

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	c.Assert(err, Equals, nil)
	c.Assert(key, Not(Equals), nil)

	pubKey := key.PublicKey
	c.Assert(pubKey, Not(Equals), nil)

	// the public key in SSH disk format
	sshAuthKey, err := RSAPubKeyToDisk(&pubKey)
	c.Assert(err, Equals, nil)

	// serialize private key
	serPrivateKey, err := RSAPrivateKeyToDisk(key)
	c.Assert(err, Equals, nil)

	// deserialize it
	key2, err := RSAPrivateKeyFromDisk(serPrivateKey)
	c.Assert(err, Equals, nil)
	c.Assert(key2, Not(Equals), nil)

	// compare serialized versions of public keys
	pubKey2 := key2.PublicKey
	sshAuthKey2, err := RSAPubKeyToDisk(&pubKey2)
	c.Assert(err, Equals, nil)
	c.Assert(string(sshAuthKey), Equals, string(sshAuthKey2))
}

func (s *XLSuite) TestRSAPubKeyToFromWire(c *C) {
	rng := rnglib.MakeSimpleRNG()

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	c.Assert(err, Equals, nil)
	c.Assert(key, Not(Equals), nil)

	pubKey := key.PublicKey
	c.Assert(pubKey, Not(Equals), nil)

	// the public key in wire format
	wirePubKey, err := RSAPubKeyToWire(&pubKey)
	c.Assert(err, Equals, nil)

	// generate the public key from the serialized version
	pk2, err := RSAPubKeyFromWire(wirePubKey)
	c.Assert(err, Equals, nil)
	c.Assert(pk2, Not(Equals), nil)

	// a long-winded way of proving pk2 == pubKey: we create a
	// random message, sign it with the private key and verify it with pk2

	msgLen := 128
	msg := make([]byte, msgLen)
	rng.NextBytes(msg)

	digest := sha1.New()
	digest.Write(msg)
	hash := digest.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, key, cr.SHA1, hash)
	c.Assert(err, IsNil)

	err = rsa.VerifyPKCS1v15(pk2, cr.SHA1, hash, sig)
	c.Assert(err, IsNil)
}
func (s *XLSuite) TestRSAPrivateKeyToFromWire(c *C) {
	rng := rnglib.MakeSimpleRNG()

	key, err := rsa.GenerateKey(rand.Reader, 1024)
	c.Assert(err, Equals, nil)
	c.Assert(key, Not(Equals), nil)

	pubKey := key.PublicKey
	c.Assert(pubKey, Not(Equals), nil)

	// the private key in wire format
	wirePrivateKey, err := RSAPrivateKeyToWire(key)
	c.Assert(err, Equals, nil)

	// generate the public key from the serialized version
	k2, err := RSAPrivateKeyFromWire(wirePrivateKey)
	c.Assert(err, Equals, nil)
	c.Assert(k2, Not(Equals), nil)

	pk2 := k2.PublicKey

	// a long-winded way of proving k2 == key: we create a
	// random message, sign it with the private key and verify it with k2

	msgLen := 128
	msg := make([]byte, msgLen)
	rng.NextBytes(msg)

	digest := sha1.New()
	digest.Write(msg)
	hash := digest.Sum(nil)

	sig, err := rsa.SignPKCS1v15(rand.Reader, key, cr.SHA1, hash)
	c.Assert(err, IsNil)

	err = rsa.VerifyPKCS1v15(&pk2, cr.SHA1, hash, sig)
	c.Assert(err, IsNil)
}
