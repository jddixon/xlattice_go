package chunks

// xlattice_go/protocol/chunks/digiListI.go

import (
//	"crypto/rsa"
)

type DigiListI interface {

	// Return the SHA3-256 hash of the Nth item in the DigiList.  Return an
	// error if there is no such item.
	//
	// There may be a requirement that this interface be called in order,
	// beginning at zero, without skipping any items.  It might be an error
	// to call this function more than once for the Nth item.
	//
	// If the DigiList has been signed, a call to this function may clear
	// the digital signature.
	HashItem(n uint) ([]byte, error)

	// Return the number of items in the list.
	Size() uint
}
