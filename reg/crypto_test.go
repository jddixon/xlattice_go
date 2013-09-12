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

func (s *XLSuite) makeAnID(c *C, rng *xr.PRNG) (id []byte) {
	id = make([]byte, SHA3_LEN)
	rng.NextBytes(&id)
	return
}
func (s *XLSuite) makeANodeID(c *C, rng *xr.PRNG) (nodeID *xi.NodeID) {
	id := s.makeAnID(c, rng)
	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))
	return
}
func (s *XLSuite) makeAnRSAKey(c *C) (key *rsa.PrivateKey) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)
	c.Assert(key, Not(IsNil))
	return key
}
func (s *XLSuite) TestCrytpo(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CRYPTO")
	}
	rng := xr.MakeSimpleRNG()
	nodeID := s.makeANodeID(c, rng)

	ckPriv := s.makeAnRSAKey(c)
	skPriv := s.makeAnRSAKey(c)

	node, err := xn.New("foo", nodeID, "", ckPriv, skPriv, nil, nil, nil)
	c.Assert(err, IsNil)
	c.Assert(node, Not(IsNil))

	ck := node.GetCommsPublicKey()

	// XXX ADDING 4 BYTE VERSION NUMBER

	// Generate 16-byte AES IV, 32-byte AES key, and 8-byte salt
	// and 4-byte version number for the Hello and another 20 bytes as salt
	// for the OAEP encrypt. For testing purposes these need not be crypto
	// grade.

	salty := make([]byte, 3*aes.BlockSize+8+4+SHA1_LEN)
	rng.NextBytes(&salty)

	iv1 := salty[:aes.BlockSize]
	key1 := salty[aes.BlockSize : 3*aes.BlockSize]
	salt1 := salty[3*aes.BlockSize : 3*aes.BlockSize+8]
	vBytes := salty[3*aes.BlockSize+8 : 3*aes.BlockSize+12]
	version1 := uint32(
		(0xff & vBytes[3] << 24) | (0xff & vBytes[2] << 16) |
			(0xff & vBytes[1] << 8) | (0xff & vBytes[0]))
	_ = version1 // DEBUG
	oaep1 := salty[3*aes.BlockSize+12:]

	oaepSalt := bytes.NewBuffer(oaep1)

	// == HELLO =====================================================
	// On the client side:
	//     create and marshal a hello message containing AES iv1, key1, salt1.
	// There is no reason at all to use protobufs for this purpose.
	// Just encrypt iv1 + key1 + salt1
	sha := sha1.New()
	data := salty[:3*aes.BlockSize+8+4] // contains iv1,key1,salt1,version1
	c.Assert(len(data), Equals, 60)

	ciphertext, err := rsa.EncryptOAEP(sha, oaepSalt, ck, data, nil)
	c.Assert(err, IsNil)
	c.Assert(ciphertext, Not(IsNil))

	// On the server side: ------------------------------------------
	// decrypt the hello using the node's private comms key
	plaintext, err := rsa.DecryptOAEP(sha, nil, ckPriv, ciphertext, nil)
	c.Assert(err, IsNil)
	c.Assert(plaintext, Not(IsNil))

	// verify that iv1, key1, salt1, version1 are the same
	c.Assert(data, DeepEquals, plaintext)

	// == HELLO REPLY ===============================================
	// On the server side:
	//     create, marshal a reply containing iv2, key2, salt2, salt1, version2
	// This could be done as a Protobuf message, but is handled as a simple
	// byte slice instead.

	// create the session iv + key plus salt2
	reply := make([]byte, 3*aes.BlockSize+8)
	rng.NextBytes(&reply)

	iv2 := reply[:aes.BlockSize]
	key2 := reply[aes.BlockSize : 3*aes.BlockSize]
	salt2 := reply[3*aes.BlockSize]

	reply = append(reply, salt1...)
	reply = append(reply, vBytes...)

	// We need padding because the message is not an integer multiple
	// of the block size.

	paddedReply, err := AddPKCS7Padding(reply, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(paddedReply, Not(IsNil))

	// encrypt the reply using engine1a = iv1, key1
	engine1a, err := aes.NewCipher(key1) // on server
	c.Assert(err, IsNil)
	c.Assert(engine1a, Not(IsNil))

	aesEncrypter1a := cipher.NewCBCEncrypter(engine1a, iv1)
	c.Assert(err, IsNil)
	c.Assert(aesEncrypter1a, Not(IsNil))

	// we require that the message size be a multiple of the block size
	c.Assert(aesEncrypter1a.BlockSize(), Equals, aes.BlockSize)
	msgLen := len(paddedReply)
	nBlocks := (msgLen + aes.BlockSize - 1) / aes.BlockSize
	c.Assert(msgLen, Equals, nBlocks*aes.BlockSize)

	ciphertext = make([]byte, nBlocks*aes.BlockSize)
	aesEncrypter1a.CryptBlocks(ciphertext, paddedReply) // dest <- src

	// On the client side: ------------------------------------------
	//     decrypt the reply using engine1b = iv1, key1

	engine1b, err := aes.NewCipher(key1) // on client
	c.Assert(err, IsNil)
	c.Assert(engine1b, Not(IsNil))

	aesDecrypter1b := cipher.NewCBCDecrypter(engine1b, iv1)
	c.Assert(err, IsNil)
	c.Assert(aesDecrypter1b, Not(IsNil))

	plaintext = make([]byte, nBlocks*aes.BlockSize)
	aesDecrypter1b.CryptBlocks(plaintext, ciphertext) // dest <- src

	c.Assert(plaintext, DeepEquals, paddedReply)
	unpaddedReply, err := StripPKCS7Padding(paddedReply, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(unpaddedReply, DeepEquals, reply)

	_ = salt2 // WE DON'T USE THIS YET

	// == CLIENT ====================================================
	// On the client side:

	// create and marshal client name, specs, salt2, digsig over that
	clientName := rng.NextFileName(8)

	// create and marshal a token containing attrs, id, ck, sk, myEnd*
	attrs := uint64(947)
	ckBytes, err := xc.RSAPrivateKeyToWire(ckPriv)
	c.Assert(err, IsNil)
	skBytes, err := xc.RSAPrivateKeyToWire(skPriv)
	c.Assert(err, IsNil)

	myEnd := []string{"127.0.0.1:4321"}
	token := &XLRegMsg_Token{
		Attrs:    &attrs,
		ID:       nodeID.Value(),
		CommsKey: ckBytes,
		SigKey:   skBytes,
		MyEnd:    myEnd,
	}

	op := XLRegMsg_Client
	clientMsg := XLRegMsg{
		Op:          &op,
		ClientName:  &clientName,
		ClientSpecs: token,
	}
	// -- CLIENT-SIDE AES SETUP -----------------
	// encrypt the client msg using engine2a = iv2, key2

	engine2a, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engine2a, Not(IsNil))

	aesEncrypter2a := cipher.NewCBCEncrypter(engine2a, iv2)
	c.Assert(err, IsNil)
	c.Assert(aesEncrypter2a, Not(IsNil))

	// we require that the message size be a multiple of the block size
	c.Assert(aesEncrypter2a.BlockSize(), Equals, aes.BlockSize)

	// -- BEGIN encode, pad, and encrypt --------
	cData, err := EncodePacket(&clientMsg)
	c.Assert(err, IsNil)
	c.Assert(cData, Not(IsNil))

	paddedCData, err := AddPKCS7Padding(cData, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(paddedCData, Not(IsNil))

	msgLen = len(paddedCData)
	nBlocks = (msgLen + aes.BlockSize - 2) / aes.BlockSize
	c.Assert(msgLen, Equals, nBlocks*aes.BlockSize)

	ciphertext = make([]byte, nBlocks*aes.BlockSize)
	aesEncrypter2a.CryptBlocks(ciphertext, paddedCData) // dest <- src

	// On the server side: ------------------------------------------

	// -- SERVER-SIDE AES SETUP -----------------
	engine2b, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engine2b, Not(IsNil))

	aesDecrypter2b := cipher.NewCBCDecrypter(engine2b, iv2)
	c.Assert(err, IsNil)
	c.Assert(aesDecrypter2b, Not(IsNil))

	// we require that the message size be a multiple of the block size
	c.Assert(aesDecrypter2b.BlockSize(), Equals, aes.BlockSize)

	// -- BEGIN decrypt, unpad, and decode ------

	// decrypt the join using engine2b = iv2, key2
	plaintext = make([]byte, nBlocks*aes.BlockSize)
	aesDecrypter2b.CryptBlocks(plaintext, ciphertext) // dest <- src

	unpaddedCData, err := StripPKCS7Padding(plaintext, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(unpaddedCData, DeepEquals, cData)

	clientMsg2, err := DecodePacket(unpaddedCData)
	c.Assert(err, IsNil)
	c.Assert(clientMsg2, Not(IsNil))
	// -- END decrypt, unpad, and decode --------

	// verify that id, ck, sk, myEnd* survive the trip unchanged

	name2 := clientMsg2.GetClientName()
	c.Assert(name2, Equals, clientName)

	clientSpecs2 := clientMsg2.GetClientSpecs()
	c.Assert(clientSpecs2, Not(IsNil))

	attrs2 := clientSpecs2.GetAttrs()
	id2 := clientSpecs2.GetID()
	ckBytes2 := clientSpecs2.GetCommsKey()
	skBytes2 := clientSpecs2.GetSigKey()
	myEnd2 := clientSpecs2.GetMyEnd() // a string array

	c.Assert(attrs2, Equals, attrs)
	c.Assert(id2, DeepEquals, nodeID.Value())
	c.Assert(ckBytes2, DeepEquals, ckBytes)
	c.Assert(skBytes2, DeepEquals, skBytes)
	c.Assert(myEnd2, DeepEquals, myEnd)

	// OK TO HERE ///

	// == CLIENT OK =================================================
	// on the server side:

	// == CREATE ====================================================

	// == JOIN ======================================================
	// On the client side:
	// == MEMBERS ===================================================
	// On the server side:

	// create and marshal a set of 3=5 tokens each containing attrs,
	// nodeID, clusterID
	// XXX STUB XXX

	// encrypt the join using engine2a = iv2, key2
	// XXX STUB XXX

	// decrypt the join using engine2b = iv2, key2
	// XXX STUB XXX

	// verify that id, ck, sk, myEnd* survive the trip unchanged
	// XXX STUB XXX

	// == MEMBER LIST  ==============================================

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
