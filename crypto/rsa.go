package crypto

// xlattice_go/crypto/rsa.go

import (
	"bytes"
	"code.google.com/p/go.crypto/ssh"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/big"
)

var _ = fmt.Print

// -- rsaPublicKey --------------------------------------------------
type rsaPublicKey rsa.PublicKey

// The presence of these methods make it possible to cast
// *rsa.PublicKey as an ssh Public Key.
func (r *rsaPublicKey) PrivateKeyAlgo() string {
	return "ssh-rsa"
}
func (r *rsaPublicKey) PublicKeyAlgo() string {
	return r.PrivateKeyAlgo()
}

// ------------------------------------------------------------------

// man 8 sshd
func ParseAuthorizedKey(in []byte) (out *rsa.PublicKey,
	comment string, options []string, rest []byte, ok bool) {

	for len(in) > 0 {
		end := bytes.IndexByte(in, '\n')
		if end != -1 {
			rest = in[end+1:]
			in = in[:end]
		} else {
			rest = nil
		}

		end = bytes.IndexByte(in, '\r')
		if end != -1 {
			in = in[:end]
		}

		in = bytes.TrimSpace(in)
		if len(in) == 0 || in[0] == '#' {
			in = rest
			continue
		}

		i := bytes.IndexAny(in, " \t")
		if i == -1 {
			in = rest
			continue
		}

		if out, comment, ok = parseSSHAuthorizedKey(in[i:]); ok {
			return
		}

		// No key type recognised. Maybe there's an options field at
		// the beginning.
		var b byte
		inQuote := false
		var candidateOptions []string
		optionStart := 0
		for i, b = range in {
			isEnd := !inQuote && (b == ' ' || b == '\t')
			if (b == ',' && !inQuote) || isEnd {
				if i-optionStart > 0 {
					candidateOptions = append(candidateOptions,
						string(in[optionStart:i]))
				}
				optionStart = i + 1
			}
			if isEnd {
				break
			}
			if b == '"' && (i == 0 || (i > 0 && in[i-1] != '\\')) {
				inQuote = !inQuote
			}
		}
		// skip whitespace = blanks and tabs
		for i < len(in) && (in[i] == ' ' || in[i] == '\t') {
			i++
		}
		if i == len(in) {
			// Invalid line: unmatched quote
			in = rest
			continue
		}

		in = in[i:]
		i = bytes.IndexAny(in, " \t")
		if i == -1 {
			in = rest
			continue
		}

		if out, comment, ok = parseSSHAuthorizedKey(in[i:]); ok {
			options = candidateOptions
			return
		}

		in = rest
		continue
	}
	return
}

// Parse a public key in OpenSSH authorized_keys format
// (see man 8 sshd) once the options and key type fields have been
// removed.
func parseSSHAuthorizedKey(in []byte) (
	out *rsa.PublicKey, comment string, ok bool) {

	in = bytes.TrimSpace(in)
	i := bytes.IndexAny(in, " \t")
	if i == -1 {
		i = len(in)
	}
	base64Key := in[:i]

	key := make([]byte, base64.StdEncoding.DecodedLen(len(base64Key)))
	n, err := base64.StdEncoding.Decode(key, base64Key)
	if err != nil {
		return
	}
	key = key[:n]
	out, _, ok = ParseSSHPublicKey(key)
	if !ok {
		return nil, "", false
	}
	comment = string(bytes.TrimSpace(in[i:]))
	return
}

// ParseSSHPublicKey - see RFC 4253, section 6.6.
func ParseSSHPublicKey(in []byte) (
	out *rsa.PublicKey, rest []byte, ok bool) {

	algo, rest, ok := ParseLenHeadedString(in)
	if ok {
		out, rest, ok = parsePubKeyByAlgo(rest, string(algo))
	}
	return
}

// Parse a public key of the given algorithm.
func parsePubKeyByAlgo(in []byte, algo string) (
	pubKey *rsa.PublicKey, rest []byte, ok bool) {

	if algo == ssh.KeyAlgoRSA {
		return ParseBareRSAPublicKey(in)
	} else {
		return nil, nil, false
	}
}

// See RFC 4253, section 6.6.
func ParseBareRSAPublicKey(in []byte) (
	key *rsa.PublicKey, rest []byte, ok bool) {

	key = new(rsa.PublicKey)
	bigE, rest, ok := parseInt(in)
	if !ok || bigE.BitLen() > 24 {
		return
	}
	e := bigE.Int64()
	if e >= 3 && e&1 != 0 {
		key.E = int(e)
		key.N, rest, ok = parseInt(rest)
	}
	return key, rest, ok
}

var BIG_ONE = big.NewInt(1)

func parseInt(in []byte) (out *big.Int, rest []byte, ok bool) {
	contents, rest, ok := ParseLenHeadedString(in)
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
func ParseLenHeadedString(in []byte) (out, rest []byte, ok bool) {
	if len(in) >= 4 {
		// first four bytes are big-endian length
		byteCount := binary.BigEndian.Uint32(in)
		if uint32(len(in)) >= 4+byteCount {
			out = in[4 : 4+byteCount]
			rest = in[4+byteCount:]
			ok = true
		}
	}
	return
}
