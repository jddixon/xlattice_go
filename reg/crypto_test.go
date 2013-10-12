package reg

// xlattice_go/reg/crypto_test.go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha1"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xn "github.com/jddixon/xlattice_go/node"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

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

	version1 := uint32(rng.Int31n(255 * 255)) // in effect an unsigned short
	vBytes := make([]byte, 4)
	vBytes[0] = byte(version1)
	c.Assert(vBytes[0], Equals, byte(0xff&version1)) // XXX trivial
	vBytes[1] = byte(version1 >> 8)
	c.Assert(vBytes[1], Equals, byte(0xff&(version1>>8))) // XXX trivial

	// Generate 16-byte AES IV, 32-byte AES key, and 8-byte salt
	// and 4-byte version number for the Hello and another 20 bytes as salt
	// for the OAEP encrypt. For testing purposes these need not be crypto
	// grade.

	salty := make([]byte, 3*aes.BlockSize+8+SHA1_LEN)
	rng.NextBytes(&salty)

	iv1 := salty[:aes.BlockSize]
	key1 := salty[aes.BlockSize : 3*aes.BlockSize]
	salt1 := salty[3*aes.BlockSize : 3*aes.BlockSize+8]
	oaep1 := salty[3*aes.BlockSize+8:]
	oaepSalt := bytes.NewBuffer(oaep1)

	// == HELLO =====================================================
	// On the client side:
	//     create and marshal a hello message containing AES iv1, key1, salt1.
	// There is no reason at all to use protobufs for this purpose.
	// Just encrypt iv1 + key1 + salt1 + version1, where version1 is the
	// version proposed for communications.

	sha := sha1.New()
	data := salty[:3*aes.BlockSize+8] // contains iv1,key1,salt1
	data = append(data, vBytes...)    // ... plus preferred protocol version
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
	iv1s := data[0:aes.BlockSize]
	key1s := data[aes.BlockSize : 3*aes.BlockSize]
	salt1s := data[3*aes.BlockSize : 3*aes.BlockSize+8]
	vBytes1s := data[3*aes.BlockSize+8:]
	c.Assert(iv1s, DeepEquals, iv1)
	c.Assert(key1s, DeepEquals, key1)
	c.Assert(salt1s, DeepEquals, salt1)
	c.Assert(vBytes1s, DeepEquals, vBytes)

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
	salt2 := reply[3*aes.BlockSize : 3*aes.BlockSize+8]

	// add the original salt, and then vBytes, representing version2
	reply = append(reply, salt1...)
	reply = append(reply, vBytes...)

	// We need padding because the message is not an integer multiple
	// of the block size.

	paddedReply, err := xc.AddPKCS7Padding(reply, aes.BlockSize)
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
	unpaddedReply, err := xc.StripPKCS7Padding(paddedReply, aes.BlockSize)
	c.Assert(err, IsNil)
	c.Assert(unpaddedReply, DeepEquals, reply)
	_ = salt2 // we don't explicitly use this

	iv2c := unpaddedReply[:aes.BlockSize]
	key2c := unpaddedReply[aes.BlockSize : 3*aes.BlockSize]
	salt2c := unpaddedReply[3*aes.BlockSize : 3*aes.BlockSize+8]
	salt1c := unpaddedReply[3*aes.BlockSize+8 : 3*aes.BlockSize+16]
	vBytes2c := unpaddedReply[3*aes.BlockSize+16 : 3*aes.BlockSize+20]
	c.Assert(iv2c, DeepEquals, iv2)
	c.Assert(key2c, DeepEquals, key2)
	c.Assert(salt2c, DeepEquals, salt2)
	c.Assert(salt1c, DeepEquals, salt1)
	c.Assert(vBytes2c, DeepEquals, vBytes)
	version2c := uint32(vBytes2c[0]) |
		(uint32(vBytes2c[1]) << 8) |
		(uint32(vBytes2c[2]) << 16) |
		(uint32(vBytes2c[3]) << 24)
	c.Assert(version2c, Equals, version1)

	// AES HANDLING FOR ALL FURTHER MESSAGES ========================

	// -- CLIENT-SIDE AES SETUP -----------------
	// encrypt the client msg using engineC = iv2, key2

	engineC, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engineC, Not(IsNil))

	encrypterC := cipher.NewCBCEncrypter(engineC, iv2)
	c.Assert(err, IsNil)
	c.Assert(encrypterC, Not(IsNil))

	decrypterC := cipher.NewCBCDecrypter(engineC, iv2)
	c.Assert(err, IsNil)
	c.Assert(decrypterC, Not(IsNil))

	// we require that the message size be a multiple of the block size
	c.Assert(encrypterC.BlockSize(), Equals, aes.BlockSize)
	c.Assert(decrypterC.BlockSize(), Equals, aes.BlockSize)

	// -- SERVER-SIDE AES SETUP -----------------
	engineS, err := aes.NewCipher(key2)
	c.Assert(err, IsNil)
	c.Assert(engineS, Not(IsNil))

	encrypterS := cipher.NewCBCEncrypter(engineS, iv2)
	c.Assert(err, IsNil)
	c.Assert(encrypterS, Not(IsNil))

	decrypterS := cipher.NewCBCDecrypter(engineS, iv2)
	c.Assert(err, IsNil)
	c.Assert(decrypterS, Not(IsNil))

	// we require that the message size be a multiple of the block size
	c.Assert(encrypterS.BlockSize(), Equals, aes.BlockSize)
	c.Assert(decrypterS.BlockSize(), Equals, aes.BlockSize)

	// == CLIENT MSG ================================================
	// On the client side:

	// create and marshal client name, specs, salt2, digsig over that
	clientName := rng.NextFileName(8)

	// create and marshal a token containing attrs, id, ck, sk, myEnds*
	attrs := uint64(947)
	ckBytes, err := xc.RSAPrivateKeyToWire(ckPriv)
	c.Assert(err, IsNil)
	skBytes, err := xc.RSAPrivateKeyToWire(skPriv)
	c.Assert(err, IsNil)

	myEnds := []string{"127.0.0.1:4321"}
	token := &XLRegMsg_Token{
		Attrs:    &attrs,
		ID:       nodeID.Value(),
		CommsKey: ckBytes,
		SigKey:   skBytes,
		MyEnds:   myEnds,
	}

	op := XLRegMsg_Client
	clientMsg := XLRegMsg{
		Op:          &op,
		ClientName:  &clientName,
		ClientSpecs: token,
	}

	ciphertext, err = EncodePadEncrypt(&clientMsg, encrypterC)
	c.Assert(err, IsNil)

	// On the server side: ------------------------------------------
	clientMsg2, err := DecryptUnpadDecode(ciphertext, decrypterS)
	c.Assert(err, IsNil)

	c.Assert(clientMsg2.GetOp(), Equals, XLRegMsg_Client)

	// verify that id, ck, sk, myEnds* survive the trip unchanged

	name2 := clientMsg2.GetClientName()
	c.Assert(name2, Equals, clientName)

	clientSpecs2 := clientMsg2.GetClientSpecs()
	c.Assert(clientSpecs2, Not(IsNil))

	attrs2 := clientSpecs2.GetAttrs()
	id2 := clientSpecs2.GetID()
	ckBytes2 := clientSpecs2.GetCommsKey()
	skBytes2 := clientSpecs2.GetSigKey()
	myEnds2 := clientSpecs2.GetMyEnds() // a string array

	c.Assert(attrs2, Equals, attrs)
	c.Assert(id2, DeepEquals, nodeID.Value())
	c.Assert(ckBytes2, DeepEquals, ckBytes)
	c.Assert(skBytes2, DeepEquals, skBytes)
	c.Assert(myEnds2, DeepEquals, myEnds)

	// == CLIENT OK =================================================
	// on the server side:

	clientID := s.makeAnID(c, rng)
	attrsBack := uint64(479)

	op = XLRegMsg_ClientOK
	clientOKMsg := XLRegMsg{
		Op:       &op,
		ClientID: clientID,
		ClientAttrs:    &attrsBack,
	}
	ciphertext, err = EncodePadEncrypt(&clientOKMsg, encrypterS)
	c.Assert(err, IsNil)

	// on the client side -------------------------------------------
	clientOK2, err := DecryptUnpadDecode(ciphertext, decrypterC)
	c.Assert(err, IsNil)

	c.Assert(clientOK2.GetOp(), Equals, XLRegMsg_ClientOK)
	clientID2 := clientOK2.GetClientID()
	c.Assert(clientID2, DeepEquals, clientID)
	attrsBack2 := clientOK2.GetClientAttrs()
	c.Assert(attrsBack2, Equals, attrsBack)

	// == CREATE ====================================================
	// on the client side:
	clusterName := rng.NextFileName(8)
	clusterSize := uint32(2 + rng.Intn(60))

	op = XLRegMsg_Create
	createMsg := XLRegMsg{
		Op:          &op,
		ClusterName: &clusterName,
		ClusterSize: &clusterSize,
	}
	ciphertext, err = EncodePadEncrypt(&createMsg, encrypterC)
	c.Assert(err, IsNil)

	// on the server side -------------------------------------------
	createMsg2, err := DecryptUnpadDecode(ciphertext, decrypterS)
	c.Assert(err, IsNil)

	c.Assert(createMsg2.GetOp(), Equals, XLRegMsg_Create)
	clusterName2 := createMsg2.GetClusterName()
	c.Assert(clusterName2, Equals, clusterName)
	clusterSize2 := createMsg2.GetClusterSize()
	c.Assert(clusterSize2, Equals, clusterSize)

	// == CREATE REPLY ==============================================
	// on the server side:
	clusterID := s.makeAnID(c, rng)
	sizeBack := uint32(2 + rng.Intn(60))

	op = XLRegMsg_CreateReply
	createReplyMsg := XLRegMsg{
		Op:          &op,
		ClusterID:   clusterID,
		ClusterSize: &sizeBack,
	}
	ciphertext, err = EncodePadEncrypt(&createReplyMsg, encrypterS)
	c.Assert(err, IsNil)

	// on the client side -------------------------------------------
	createReply2, err := DecryptUnpadDecode(ciphertext, decrypterC)
	c.Assert(err, IsNil)

	c.Assert(createReply2.GetOp(), Equals, XLRegMsg_CreateReply)
	clusterID2 := createReply2.GetClusterID()
	c.Assert(clusterID2, DeepEquals, clusterID)
	sizeBack2 := createReply2.GetClusterSize()
	c.Assert(sizeBack2, Equals, sizeBack)

	// == JOIN ======================================================
	// On the client side:
	op = XLRegMsg_Join
	joinMsg := XLRegMsg{
		Op:        &op,
		ClusterID: clusterID,
	}
	ciphertext, err = EncodePadEncrypt(&joinMsg, encrypterC)
	c.Assert(err, IsNil)

	// on the server side -------------------------------------------
	join2, err := DecryptUnpadDecode(ciphertext, decrypterS)
	c.Assert(err, IsNil)

	c.Assert(join2.GetOp(), Equals, XLRegMsg_Join)
	clusterID2 = join2.GetClusterID()
	c.Assert(clusterID2, DeepEquals, clusterID) // GEEP

	// == JOIN REPLY ================================================
	// on the server side:
	op = XLRegMsg_JoinReply
	joinReplyMsg := XLRegMsg{
		Op:          &op,
		ClusterID:   clusterID,
		ClusterSize: &clusterSize,
	}
	ciphertext, err = EncodePadEncrypt(&joinReplyMsg, encrypterS)
	c.Assert(err, IsNil)

	// on the client side -------------------------------------------
	joinReply2, err := DecryptUnpadDecode(ciphertext, decrypterC)
	c.Assert(err, IsNil)

	c.Assert(joinReply2.GetOp(), Equals, XLRegMsg_JoinReply)
	clusterID2 = joinReply2.GetClusterID()
	c.Assert(clusterID2, DeepEquals, clusterID)
	sizeBack = joinReply2.GetClusterSize()
	c.Assert(sizeBack, Equals, clusterSize)

	// == GET =======================================================
	// On the client side:

	// on the server side -------------------------------------------

	// == MEMBERS ===================================================
	// On the server side:

	// on the client side -------------------------------------------

	// create and marshal a set of 3=5 tokens each containing attrs,
	// nodeID, clusterID
	// XXX STUB XXX

	// encrypt the msg using engineC = iv2, key2
	// XXX STUB XXX

	// decrypt the msg using engineS = iv2, key2
	// XXX STUB XXX

	// verify that the various tokens (id, ck, sk, myEnds*) survive the
	// trip unchanged
	// XXX STUB XXX

	// == BYE =======================================================
	// On the client side:
	op = XLRegMsg_Bye
	byeMsg := XLRegMsg{
		Op: &op,
	}
	ciphertext, err = EncodePadEncrypt(&byeMsg, encrypterC)
	c.Assert(err, IsNil)

	// on the server side -------------------------------------------
	bye2, err := DecryptUnpadDecode(ciphertext, decrypterS)
	c.Assert(err, IsNil)
	c.Assert(bye2.GetOp(), Equals, XLRegMsg_Bye)

	// == ACK =======================================================
	// on the server side:
	op = XLRegMsg_Ack
	ackMsg := XLRegMsg{
		Op: &op,
	}
	ciphertext, err = EncodePadEncrypt(&ackMsg, encrypterS)
	c.Assert(err, IsNil)

	// on the client side -------------------------------------------
	ack2, err := DecryptUnpadDecode(ciphertext, decrypterC)
	c.Assert(err, IsNil)
	c.Assert(ack2.GetOp(), Equals, XLRegMsg_Ack)
}
