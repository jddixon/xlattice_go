package reg

// xlattice_go/msg/crypto_test.go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"

	"errors"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

// TODO: MOVE THIS TO crypto/ =======================================

// PKCS7 padding (RFC 5652) pads a message out to a whole multiple
// of the block size, with the value of each byte being the number
// of bytes of padding.  If the data passed is nil, the function
// returns a full block of padding.

func PKCS7Padding(data []byte, blockSize int) (padding []byte) {
	var length int
	if data == nil {
		length = 0
	} else {
		length = len(data)
	}
	// we want from 1 to blockSize bytes of padding
	nBlocks := (length + blockSize - 1) / blockSize
	rem := nBlocks*blockSize - length
	if rem == 0 {
		rem = blockSize
	}
	padding = make([]byte, rem)
	for i := 0; i < rem; i++ {
		padding[i] = byte(rem)
	}
	return
}

var (
	ImpossibleBlockSize   = errors.New("impossible block size")
	IncorrectPKCS7Padding = errors.New("incorrectly padded data")
	NilData               = errors.New("nil data argument")
)

func AddPKCS7Padding(data []byte, blockSize int) (out []byte, err error) {
	if blockSize <= 1 {
		err = ImpossibleBlockSize
	} else {
		padding := PKCS7Padding(data, blockSize)
		if data == nil {
			out = padding
		} else {
			out = append(data, padding...)
		}
	}
	return
}

// The data passed is presumed to have PKCS7 padding.  If possible, return
// a copy of the data without the padding.  Return an error if the padding
// is incorrect.

func StripPKCS7Padding(data []byte, blockSize int) (out []byte, err error) {
	if blockSize <= 1 {
		err = ImpossibleBlockSize
	} else if data == nil {
		err = NilData
	}
	if err == nil {
		lenData := len(data)
		if lenData < blockSize {
			err = IncorrectPKCS7Padding
		} else {
			// examine the very last byte: it must be padding and must
			// contain the number of padding bytes added
			lenPadding := int(data[lenData-1])
			if lenPadding < 1 || lenData < lenPadding {
				err = IncorrectPKCS7Padding
			} else {
				out = data[:lenData-lenPadding]
			}
		}
	}
	return
}
func (s *XLSuite) TestPKCS7Padding(c *C) {
	rng := xr.MakeSimpleRNG()
	seven := make([]byte, 7)
	rng.NextBytes(&seven)

	fifteen := make([]byte, 15)
	rng.NextBytes(&fifteen)

	sixteen := make([]byte, 16)
	rng.NextBytes(&sixteen)

	seventeen := make([]byte, 17)
	rng.NextBytes(&seventeen)

	padding := PKCS7Padding(seven, aes.BlockSize)
	c.Assert(len(padding), Equals, aes.BlockSize-7)
	c.Assert(padding[0], Equals, byte(aes.BlockSize-7))

	padding = PKCS7Padding(fifteen, aes.BlockSize)
	c.Assert(len(padding), Equals, aes.BlockSize-15)
	c.Assert(padding[0], Equals, byte(aes.BlockSize-15))

	padding = PKCS7Padding(sixteen, aes.BlockSize)
	c.Assert(len(padding), Equals, aes.BlockSize)
	c.Assert(padding[0], Equals, byte(16))

	padding = PKCS7Padding(seventeen, aes.BlockSize)
	expectedLen := 2*aes.BlockSize - 17
	c.Assert(len(padding), Equals, expectedLen)
	c.Assert(padding[0], Equals, byte(expectedLen))

	paddedSeven, err := AddPKCS7Padding(seven, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedSeven, err := StripPKCS7Padding(paddedSeven, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(seven, DeepEquals, unpaddedSeven)

	paddedFifteen, err := AddPKCS7Padding(fifteen, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedFifteen, err := StripPKCS7Padding(paddedFifteen, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(fifteen, DeepEquals, unpaddedFifteen)

	paddedSixteen, err := AddPKCS7Padding(sixteen, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedSixteen, err := StripPKCS7Padding(paddedSixteen, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(sixteen, DeepEquals, unpaddedSixteen)

	paddedSeventeen, err := AddPKCS7Padding(seventeen, aes.BlockSize)
	c.Assert(err, IsNil)
	unpaddedSeventeen, err := StripPKCS7Padding(paddedSeventeen, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(seventeen, DeepEquals, unpaddedSeventeen)
}

// END MOVE THIS TO crypto/ =========================================

func (s *XLSuite) TestCrytpo(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CRYPTO")
	}

	rng := xr.MakeSimpleRNG()
	id := make([]byte, SHA3_LEN)
	rng.NextBytes(&id)
	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))

	ckPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)
	c.Assert(ckPriv, Not(IsNil))
	skPriv, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)
	c.Assert(skPriv, Not(IsNil))

	node, err := xn.New("foo", nodeID, "", ckPriv, skPriv, nil, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(node, Not(IsNil))

	ck := node.GetCommsPublicKey()

	// Generate 16-byte AES IV, 32-byte AES key, and 8-byte salt
	// for the Hello and another 20 bytes as salt for the OAEP encrypt
	// For testing purposes these need not be crypto grade.

	salty := make([]byte, 3*aes.BlockSize+8+SHA1_LEN)
	rng.NextBytes(&salty)

	iv1 := salty[:aes.BlockSize]
	key1 := salty[aes.BlockSize : 3*aes.BlockSize]
	salt1 := salty[3*aes.BlockSize : 3*aes.BlockSize+8]
	oaep1 := salty[3*aes.BlockSize+8:]

	oaepSalt := bytes.NewBuffer(oaep1)

	// -- HELLO -----------------------------------------------------
	// On the client side:
	//     create and marshal a hello message containing AES iv1, key1, salt1.
	// There is no reason at all to use protobufs for this purpose.
	// Just encrypt iv1 + key1 + salt1
	sha := sha1.New()
	data := salty[:3*aes.BlockSize+8] // contains iv1,key1,salt1
	c.Assert(len(data), Equals, 56)

	ciphertext, err := rsa.EncryptOAEP(sha, oaepSalt, ck, data, nil)
	c.Assert(err, IsNil)
	c.Assert(ciphertext, Not(IsNil))

	// On the server side:
	// decrypt the hello using the node's private comms key
	plaintext, err := rsa.DecryptOAEP(sha, nil, ckPriv, ciphertext, nil)
	c.Assert(err, IsNil)
	c.Assert(plaintext, Not(IsNil))

	// verify that iv1, key1, and salt1 are the same
	c.Assert(data, DeepEquals, plaintext)

	// -- HELLO REPLY -----------------------------------------------
	// On the server side:
	//     create and marshal a hello reply containing iv2, key2, salt2, salt1
	reply := make([]byte, 3*aes.BlockSize+8)
	rng.NextBytes(&reply)

	iv2 := reply[:aes.BlockSize]
	key2 := reply[aes.BlockSize : 3*aes.BlockSize]
	salt2 := reply[3*aes.BlockSize]

	_, _, _, _ = iv2, key2, salt2, skPriv // DEBUG

	reply = append(reply, salt1...)

	// THERE IS NO NEED FOR PADDING because we have made the message
	// an integer multiple of the block size

	// encrypt the reply using engine1a = iv1, key1
	engine1a, err := aes.NewCipher(key1) // on server
	c.Assert(err, IsNil)
	c.Assert(engine1a, Not(IsNil))

	aesEncrypter1a := cipher.NewCBCEncrypter(engine1a, iv1)
	c.Assert(err, IsNil)
	c.Assert(aesEncrypter1a, Not(IsNil))

	// we require that the message size be a multiple of the block size
	c.Assert(aesEncrypter1a.BlockSize(), Equals, aes.BlockSize)
	msgLen := len(reply)
	nBlocks := (msgLen + aes.BlockSize - 1) / aes.BlockSize
	c.Assert(msgLen, Equals, nBlocks*aes.BlockSize)

	ciphertext = make([]byte, nBlocks*aes.BlockSize)
	aesEncrypter1a.CryptBlocks(ciphertext, reply) // dest <- src

	// On the client side:
	//     decrypt the reply using engine1b = iv1, key1

	engine1b, err := aes.NewCipher(key1) // on client
	c.Assert(err, IsNil)
	c.Assert(engine1b, Not(IsNil))

	aesDecrypter1b := cipher.NewCBCDecrypter(engine1b, iv1)
	c.Assert(err, IsNil)
	c.Assert(aesDecrypter1b, Not(IsNil))

	plaintext = make([]byte, nBlocks*aes.BlockSize)
	aesDecrypter1b.CryptBlocks(plaintext, ciphertext) // dest <- src

	c.Assert(plaintext, DeepEquals, reply)

	// -- JOIN ------------------------------------------------------
	// On the client side:

	// create and marshal a join message containing id, ck, sk, myEnd*
	attrs := uint64(947)
	ckBytes, err := xc.RSAPrivateKeyToWire(ckPriv)
	c.Assert(err, IsNil)
	skBytes, err := xc.RSAPrivateKeyToWire(skPriv)
	c.Assert(err, IsNil)

	myEnd := []string{"127.0.0.1:4321"}
	token := &XLRegMsg_Token{
		Attrs:    &attrs,
		ID:       id,
		CommsKey: ckBytes,
		SigKey:   skBytes,
		MyEnd:    myEnd,
	}

	op := XLRegMsg_Join
	joinMsg := XLRegMsg{
		Op:      &op,
		MySpecs: token,
	}
	data, err = EncodePacket(&joinMsg)
	c.Assert(err, IsNil)
	c.Assert(data, Not(IsNil))

	// XXX MISSING CRYPTO BIT ;-)

	decoded, err := DecodePacket(data)
	c.Assert(err, IsNil)
	c.Assert(decoded, Not(IsNil))

	// encrypt the join using engine2a = iv2, key2
	// XXX STUB XXX

	engine2a, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engine2a, Not(IsNil))

	// On the server side:
	engine2b, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engine2b, Not(IsNil))

	// decrypt the join using engine2b = iv2, key2
	// XXX STUB XXX

	// verify that id, ck, sk, myEnd* survive the trip unchanged
	// XXX STUB XXX

	// -- MEMBERS ---------------------------------------------------
	// On the server side:

	// create and marshal a set of 3-5 tokens each containing attrs,
	// nodeID, clusterID
	// XXX STUB XXX

	// encrypt the join using engine2a = iv2, key2
	// XXX STUB XXX

	// decrypt the join using engine2b = iv2, key2
	// XXX STUB XXX

	// verify that id, ck, sk, myEnd* survive the trip unchanged
	// XXX STUB XXX

	// -- MEMBER LIST  ----------------------------------------------

	// LOOP ANOTHER N TIMES WITH BLOCKS OF RANDOM DATA -----------
	for i := 0; i < 8; i++ {
		// create block of data
		// XXX STUB XXX

		// encrypt it with engine2a (iv2, key2)
		// XXX STUB XXX

		// decrypt with engine2b (iv2, key2)
		// XXX STUB XXX

		// verify same data
		// XXX STUB XXX

	}
}
