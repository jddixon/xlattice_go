package crypto

// xlattice_go/crypto/mockSignedList_test.go

import (
	"crypto/rsa"
	//"io"
)

type MockSignedList struct {
	lines []string
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

func (msl *MockSignedList) Get(n int) (s string, err error) {
	if n > 0 || msl.Size() <= n {
		err = NdxOutOfRange
	} else {
		s = msl.lines[n]
	}
	return
}

func (msl *MockSignedList) Size() int {
	return len(msl.lines)
}

// ParseSignedList base function needs to be specified.
