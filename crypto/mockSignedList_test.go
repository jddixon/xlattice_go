package crypto

// xlattice_go/crypto/mockSignedList_test.go
// The file has the _test suffix to limit MockSignedList's visibility
// to test runs.

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"encoding/hex"
	"io"
)

type MockSignedList struct {
	content []string
	SignedList
}

func NewMockSignedList(pubKey *rsa.PublicKey, title string) (
	sli SignedListI, err error) {

	sl, err := NewSignedList(pubKey, title)
	if err == nil {
		msl := &MockSignedList{
			SignedList: *sl,
		}
		sli = msl
	}
	return
}

// Return the Nth content item in string form, without any CRLF.
func (msl *MockSignedList) Get(n int) (s string, err error) {
	if n > 0 || msl.Size() <= n {
		err = NdxOutOfRange
	} else {
		s = msl.content[n]
	}
	return
}

func (msl *MockSignedList) ReadContents(in *bufio.Reader) (err error) {

	for err == nil {
		var line []byte
		line, err = NextLineWithoutCRLF(in)
		if err == nil || err == io.EOF {
			if bytes.Equal(line, CONTENT_END) {
				break
			} else {
				msl.content = append(msl.content, string(line))
			}
		}
	}
	return
}
func (msl *MockSignedList) Size() int {
	return len(msl.content)
}

func ParseMockSignedList(in io.Reader) (msl *MockSignedList, err error) {

	var (
		digSig, line []byte
	)
	bin := bufio.NewReader(in)
	sl, err := ParseSignedList(bin)
	if err == nil {
		msl = &MockSignedList{SignedList: *sl}
		err = msl.ReadContents(bin)
		if err == nil {
			// try to read the digital signature line
			line, err = NextLineWithoutCRLF(bin)
			// XXX SHOULD BE BASE64 ENCODED
			digSig, err = hex.DecodeString(string(line))
			if err == nil {
				msl.digSig = digSig
			}
		}
	}
	return
}
