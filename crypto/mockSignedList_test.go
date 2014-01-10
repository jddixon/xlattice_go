package crypto

// xlattice_go/crypto/mockSignedList_test.go
// The file has the _test suffix to limit MockSignedList's visibility
// to test runs.

import (
	"bufio"
	"bytes"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"io"
	"strings"
)

var _ = fmt.Print

type MockSignedList struct {
	content []string
	SignedList
}

func NewMockSignedList(pubKey *rsa.PublicKey, title string) (
	msl *MockSignedList, err error) {

	sl, err := NewSignedList(pubKey, title)
	if err == nil {
		msl = &MockSignedList{
			SignedList: *sl,
		}
	}
	return
}

func (msl *MockSignedList) AddItem(s string) (n uint) {
	n = uint(len(msl.content)) // index of this item
	msl.content = append(msl.content, s)
	return
}

// Return the Nth content item in string form, without any CRLF.
func (msl *MockSignedList) Get(n uint) (s string, err error) {
	if n < 0 || msl.Size() <= n {
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
func (msl *MockSignedList) Size() uint {
	return uint(len(msl.content))
}

/**
 * Serialize the entire document.  All lines are CRLF-terminated.
 * If any error is encountered, this function silently returns an
 * empty string.
 */
func (msl *MockSignedList) String() (s string) {

	var (
		err error
		ss  []string
	)
	pk, title, timestamp := msl.SignedList.Strings()
	ss = append(ss, title)
	ss = append(ss, timestamp)

	// content lines ----------------------------------
	ss = append(ss, string(CONTENT_START))
	for i := uint(0); err == nil && i < msl.Size(); i++ {
		var line string
		line, err = msl.Get(i)
		if err == nil || err == io.EOF {
			ss = append(ss, line)
			if err == io.EOF {
				err = nil
				break
			}
		}
	}
	if err == nil {
		ss = append(ss, string(CONTENT_END))

		myDigSig := base64.StdEncoding.EncodeToString(msl.digSig)
		ss = append(ss, myDigSig)
		s = string(pk) + strings.Join(ss, CRLF) + CRLF
	}
	return
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
			if err == nil || err == io.EOF {
				digSig, err = base64.StdEncoding.DecodeString(string(line))
			}
			if err == nil || err == io.EOF {
				msl.digSig = digSig
			}
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}
