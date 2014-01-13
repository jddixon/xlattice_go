package crypto

// xlattice_go/crypto/errors.go

import (
	e "errors"
)

var (
	CantAddToSignedList  = e.New("can't add, list has been signed")
	EmptyTitle              = e.New("empty title parameter")
	ImpossibleBlockSize     = e.New("impossible block size")
	IncorrectPKCS7Padding   = e.New("incorrectly padded data")
	ListAlreadySigned       = e.New("list has already been signed")
	MissingContentStart     = e.New("missing CONTENT START line")
	NdxOutOfRange           = e.New("list index out of range")
	NilData                 = e.New("nil data argument")
	NilPrivateKey           = e.New("nil private key parameter")
	NilPublicKey            = e.New("nil public key parameter")
	NotAnRSAPrivateKey      = e.New("Not an RSA private key")
	NotImplemented          = e.New("not implemented")
	NotAnRSAPublicKey       = e.New("Not an RSA public key")
	PemEncodeDecodeFailure  = e.New("Pem encode/decode failure")
	UnsignedList            = e.New("list has not been signed")
	X509ParseOrMarshalError = e.New("X509 parse/marshal error")
)
