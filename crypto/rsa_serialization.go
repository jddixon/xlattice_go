package crypto

// xlattice_go/crypto/rsa_serialization.go

import (
	// "bytes"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	//"encoding/binary"
	"code.google.com/p/go.crypto/ssh"
	"encoding/pem"
	//"math/big"
)

var (
	NilData                 = errors.New("Nil data")
	NotAnRSAPrivateKey      = errors.New("Not an RSA private key")
	NotAnRSAPublicKey       = errors.New("Not an RSA public key")
	PemEncodeDecodeFailure  = errors.New("Pem encode/decode failure")
	X509ParseOrMarshalError = errors.New("X509 parse/marshal error")
)

// CONVERSION TO AND FROM WIRE FORMAT ///////////////////////////////

// Serialize an RSA public key to wire format
func RSAPubKeyToWire(pubKey *rsa.PublicKey) ([]byte, bool) {

	// XXX STUB

	return nil, false
}

// Deserialize an RSA public key from wire format
func RSAPubKeyFromWire(data []byte) (*rsa.PublicKey, error) {

	// XXX STUB

	return nil, nil
} // FOO

// Serialize an RSA private key to wire format
func RSAPrivKeyToWire(pubKey *rsa.PrivateKey) ([]byte, bool) {

	// XXX STUB

	return nil, false
}

// Deserialize an RSA private key from wire format
func RSAPrivKeyFromWire(data []byte) (*rsa.PrivateKey, error) {

	// XXX STUB

	return nil, nil
} // FOO

// CONVERSION TO AND FROM DISK FORMAT ///////////////////////////////

// Serialize an RSA public key to disk format, specifically to the
// format used by SSH. Should return nil if the conversion fails.
func RSAPubKeyToDisk(pubKey *rsa.PublicKey) ([]byte, bool) {
	out := ssh.MarshalAuthorizedKey(pubKey)
	// STUB ?
	return out, true
}

// Deserialize an RSA public key from the format used in SSH
// key files
func RSAPubKeyFromDisk(data []byte) (*rsa.PublicKey, error) {
	out, comment, options, rest, ok := ssh.ParseAuthorizedKey(data)
	_, _, _ = comment, options, rest
	if ok {
		return out.(*rsa.PublicKey), nil
	} else {
		return nil, NotAnRSAPublicKey
	}
}

// Serialize an RSA private key to disk format
func RSAPrivKeyToDisk(privKey *rsa.PrivateKey) (data []byte, err error) {
	if privKey == nil {
		err = NilData
	} else {
		data509 := x509.MarshalPKCS1PrivateKey(privKey)
		if data509 == nil {
			err = X509ParseOrMarshalError
		} else {
			block := pem.Block{Bytes: data509}
			data = pem.EncodeToMemory(&block)
		}
	}
	return
}

// Deserialize an RSA private key from disk format
func RSAPrivKeyFromDisk(data []byte) (key *rsa.PrivateKey, err error) {
	if data == nil {
		err = NilData
	} else {
		block, _ := pem.Decode(data)
		if block == nil {
			err = PemEncodeDecodeFailure
		} else {
			key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		}
	}
	return
}
