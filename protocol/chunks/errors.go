package chunks

import (
	e "errors"
)

var (
	EmptyTitle                  = e.New("empty title parameter")
	MismatchedPublicPrivateKeys = e.New("public and private keys don't match")
	NilContentHash              = e.New("nil content hash parameter")
	NilData                     = e.New("nil data parameter")
	NilDatum                    = e.New("nil datum parameter")
	NilReader                   = e.New("nil io.Reader")
	NilRSAPrivKey               = e.New("nil RSA private key parameter")
	NilSubClass                 = e.New("nil subClass parameter")
	NoDigSig                    = e.New("no digital signature: the list is not signed")
	NoNthItem                   = e.New("no Nth item")
	ZeroLengthInput             = e.New("zero length input")
)
