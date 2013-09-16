package reg

// xlattice_go/reg/helloAndReply.go

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/sha1"
	xc "github.com/jddixon/xlattice_go/crypto"
	xr "github.com/jddixon/xlattice_go/rnglib"
)

// Create an AES IV and key and an 8-byte salt, then encrypt these and
// the proposed protocol version using the server's comms public key.
func ClientEncodeHello(version1 uint32, ck *rsa.PublicKey) (
	ciphertext []byte, iv1, key1, salt1 []byte, err error) {
	rng := xr.MakeSystemRNG()

	vBytes := make([]byte, 4)
	vBytes[0] = byte(version1)
	vBytes[1] = byte(version1 >> 8)
	vBytes[2] = byte(version1 >> 16)
	vBytes[3] = byte(version1 >> 24)

	// Generate 16-byte AES IV, 32-byte AES key, and 8-byte salt for the
	// Hello and another 20 bytes as salt for the OAEP encryp
	salty := make([]byte, 3*aes.BlockSize+8+SHA1_LEN)
	rng.NextBytes(&salty)

	iv1 = salty[:aes.BlockSize]
	key1 = salty[aes.BlockSize : 3*aes.BlockSize]
	salt1 = salty[3*aes.BlockSize : 3*aes.BlockSize+8]
	oaep1 := salty[3*aes.BlockSize+8:]
	oaepSalt := bytes.NewBuffer(oaep1)

	sha := sha1.New()
	data := salty[:3*aes.BlockSize+8] // contains iv1,key1,salt1
	data = append(data, vBytes...)    // ... plus preferred protocol version

	ciphertext, err = rsa.EncryptOAEP(sha, oaepSalt, ck, data, nil)
	return
}

// Decrypt the Hello using the node's private comms key, and decode its
// contents.
func ServerDecodeHello(ciphertext []byte, ckPriv *rsa.PrivateKey) (
	iv1s, key1s, salt1s []byte, version1s uint32, err error) {

	sha := sha1.New()
	data, err := rsa.DecryptOAEP(sha, nil, ckPriv, ciphertext, nil)
	if err == nil {
		iv1s = data[0:aes.BlockSize]
		key1s = data[aes.BlockSize : 3*aes.BlockSize]
		salt1s = data[3*aes.BlockSize : 3*aes.BlockSize+8]
		vBytes := data[3*aes.BlockSize+8:]
		version1s = uint32(vBytes[0]) |
			uint32(vBytes[1])<<8 |
			uint32(vBytes[2])<<16 |
			uint32(vBytes[3])<<24
	}
	return
}

// Create and marshal using AES iv1 and key1 a reply containing iv2, key2,
// salt2, salt1 and version 2, the server-decreed protocol version number.
func ServerEncodeHelloReply(iv1, key1, salt1 []byte, version2 uint32) (
	iv2, key2, salt2, ciphertext []byte, err error) {

	var engine1a cipher.Block

	vBytes := make([]byte, 4)
	vBytes[0] = byte(version2)
	vBytes[1] = byte(version2 >> 8)
	vBytes[2] = byte(version2 >> 16)
	vBytes[3] = byte(version2 >> 24)

	rng := xr.MakeSystemRNG()
	reply := make([]byte, 3*aes.BlockSize+8)
	rng.NextBytes(&reply)

	iv2 = reply[:aes.BlockSize]
	key2 = reply[aes.BlockSize : 3*aes.BlockSize]
	salt2 = reply[3*aes.BlockSize : 3*aes.BlockSize+8]

	// add the original salt, and then vBytes, representing version2
	reply = append(reply, salt1...)
	reply = append(reply, vBytes...)

	// We need padding because the message is not an integer multiple
	// of the block size.
	paddedReply, err := xc.AddPKCS7Padding(reply, aes.BlockSize)
	if err == nil {
		// encrypt the reply using engine1a = iv1, key1
		engine1a, err = aes.NewCipher(key1) // on server
	}
	if err == nil {
		aesEncrypter1a := cipher.NewCBCEncrypter(engine1a, iv1)

		// we require that the message size be a multiple of the block size
		msgLen := len(paddedReply)
		nBlocks := (msgLen + aes.BlockSize - 1) / aes.BlockSize
		ciphertext = make([]byte, nBlocks*aes.BlockSize)
		aesEncrypter1a.CryptBlocks(ciphertext, paddedReply) // dest <- src
	}
	return
}

// Decrypt the reply using AES iv1 and key1, then decode from the reply
// iv2, key2, an 8-byte salt2, and the original salt1.

func ClientDecodeHelloReply(ciphertext, iv1, key1 []byte) (
	iv2, key2, salt2, salt1 []byte, version2 uint32, err error) {

	var unpaddedReply []byte

	engine1b, err := aes.NewCipher(key1) // on client
	if err == nil {
		aesDecrypter1b := cipher.NewCBCDecrypter(engine1b, iv1)
		plaintext := make([]byte, len(ciphertext))
		aesDecrypter1b.CryptBlocks(plaintext, ciphertext) // dest <- src
		unpaddedReply, err = xc.StripPKCS7Padding(plaintext, aes.BlockSize)
	}
	if err == nil {
		iv2 = unpaddedReply[:aes.BlockSize]
		key2 = unpaddedReply[aes.BlockSize : 3*aes.BlockSize]
		salt2 = unpaddedReply[3*aes.BlockSize : 3*aes.BlockSize+8]
		salt1 = unpaddedReply[3*aes.BlockSize+8 : 3*aes.BlockSize+16]

		vBytes2 := unpaddedReply[3*aes.BlockSize+16 : 3*aes.BlockSize+20]
		version2 = uint32(vBytes2[0]) |
			(uint32(vBytes2[1]) << 8) |
			(uint32(vBytes2[2]) << 16) |
			(uint32(vBytes2[3]) << 24)
	}
	return
}
