package crypto

// xlattice_go/crypto/rsa_serialization.go

import (
	"code.google.com/p/go.crypto/ssh"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

var _ =  fmt.Print


// CONVERSION TO AND FROM WIRE FORMAT ///////////////////////////////

// Serialize an RSA public key to wire format
func RSAPubKeyToWire(pubKey *rsa.PublicKey) ([]byte, error) {

	return x509.MarshalPKIXPublicKey(pubKey)
}

// Deserialize an RSA public key from wire format
func RSAPubKeyFromWire(data []byte) (pub *rsa.PublicKey, err error) {
	pk, err := x509.ParsePKIXPublicKey(data)
	if err == nil {
		pub = pk.(*rsa.PublicKey)
	}
	return
}

// Serialize an RSA private key to wire format
func RSAPrivateKeyToWire(privKey *rsa.PrivateKey) (data []byte, err error) {
	data = x509.MarshalPKCS1PrivateKey(privKey)
	return
}

// Deserialize an RSA private key from wire format
func RSAPrivateKeyFromWire(data []byte) (key *rsa.PrivateKey, err error) {
	return x509.ParsePKCS1PrivateKey(data)
}

// CONVERSION TO AND FROM DISK FORMAT ///////////////////////////////

// Serialize an RSA public key to disk format, specifically to the
// format used by SSH. Should return nil if the conversion fails.
func RSAPubKeyToDisk(rsaPubKey *rsa.PublicKey) (out []byte, err error) {
	pubKey, err := ssh.NewPublicKey(rsaPubKey)
	if err == nil {
		out = ssh.MarshalAuthorizedKey(pubKey)
	}
	return out, nil
}

// Deserialize an RSA public key from the format used in SSH
// key files
func RSAPubKeyFromDisk(data []byte) (*rsa.PublicKey, error) {
	// out, _, _, _, ok := ssh.ParseAuthorizedKey(data)
	out, _, _, _, ok := ParseAuthorizedKey(data)
	_ = out	// DEBUG
	if ok {
		return out, nil
	} else {
		return nil, NotAnRSAPublicKey
	}
}

// Serialize an RSA private key to disk format
func RSAPrivateKeyToDisk(privKey *rsa.PrivateKey) (data []byte, err error) {
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
func RSAPrivateKeyFromDisk(data []byte) (key *rsa.PrivateKey, err error) {
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
