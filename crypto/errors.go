package crypto

// xlattice_go/crypto/errors.go

import (
	"errors"
)

var (
	ImpossibleBlockSize     = errors.New("impossible block size")
	IncorrectPKCS7Padding   = errors.New("incorrectly padded data")
	NilData                 = errors.New("nil data argument")
	NotAnRSAPrivateKey      = errors.New("Not an RSA private key")
	NotAnRSAPublicKey       = errors.New("Not an RSA public key")
	PemEncodeDecodeFailure  = errors.New("Pem encode/decode failure")
	X509ParseOrMarshalError = errors.New("X509 parse/marshal error")
)
