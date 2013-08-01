package crypto

// xlattice_go/crypto/rsa.go

import (
	// "bytes"
	// "crypto/rsa"
	"encoding/binary"
	"math/big"
)

// convert Key, PublicKey <-> []byte

// convert rsa key, public key as []byte to and from ssh-compatible string form

// use host version of ssh_keygen to generate an RSA key of a given bit length

//

// UTILITIES FOR CONVERTING TO AND FROM WIRE FORMAT /////////////////
// Given a byte slice preceded by a 4-byte big-endian length,
// extract that many bytes as a big-endian value, and return that
// value and the remaining bytes.

// the go compiler cannot be convinced that this is a constant
var BIG_ONE = big.NewInt(1)

func parseInt(in []byte) (out *big.Int, rest []byte, ok bool) {
	contents, rest, ok := parseString(in)
	if !ok {
		return
	}
	out = new(big.Int)

	if len(contents) > 0 && (contents[0]&0x80) == 0x80 {
		// a negative number
		notBytes := make([]byte, len(contents))
		for i := range notBytes {
			notBytes[i] = ^contents[i]
		}
		out.SetBytes(notBytes)
		out.Add(out, BIG_ONE)
		out.Neg(out)
	} else {
		// a positive number
		out.SetBytes(contents)
	}
	ok = true
	return
}

// Extract a subslice from a byte slice by construing the first
// four bytes as a big-endian uint32, returning that many bytes
// and any remainder as another subslice.
func parseString(in []byte) (out, rest []byte, ok bool) {
	if len(in) < 4 {
		return // not ok by default
	}
	length := binary.BigEndian.Uint32(in)
	if uint32(len(in)) < 4+length {
		return // too short, so length is invalid
	}
	out = in[4 : 4+length]
	rest = in[4+length:]
	ok = true
	return
}
