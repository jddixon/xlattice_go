package chunks

// xlattice_go/protocol/chunks/digilistI.go

import (
	"crypto/rsa"
)

type DigiListI interface {

	// Return the SHA3-256 hash of the Nth item in the DigiList.  Return an
	// error if there is no such item.
	//
	// There may be a requirement that this interface be called in order,
	// beginning at zero, without skipping any items.  It might be an error
	// to call this function more than once for the Nth item.
	//
	// If the DigiList has been signed, a call to this function will clear
	// the digital signature.
	HashItem(n int) ([]byte, error)

	// If there are any items in the DigiList, sign it.  Any existing
	// signature is overwritten.  The public part of the key is written
	// to the data structure.
	Sign(key *rsa.PrivateKey) error

	// Serialize the DigiList, terminating each field and each item
	// with a CRLF.
	String() string

	// Return the number of items currently in the DigiList.
	Size() int

	// If the DigiList has been signed, verify the digital signature.
	// Otherwise return false.
	Verify() bool
}
