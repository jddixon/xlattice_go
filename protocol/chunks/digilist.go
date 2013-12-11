package chunks

// xlattice_go/protocol/chunks/digilist.go

import (
	"crypto/rsa"
	"crypto/sha1"
)

// A DigiList is an abstract class.  Such a list is associated with an
// RSA key.  The public part of that key is part of this structure.
// The list has a title and a timestamp.  The owner of the key is
// responsible for making these meaningful.
//
// Items are added to the list using the DigiListI interface's
// HashItem (n int, interface{}) function.  When all items have been
// added, DigiList.Sign() is used to add the digital signature to
// the data structure.
type DigiList struct {
	sk        *rsa.PublicKey
	title     string
	timestamp int64
	digSig    []byte
}

func NewDigiList(sk *rsa.PublicKey, title string, timestamp int64) (
	dl *DigiList, err error) {

	// nil public key is acceptable
	if title == "" {
		err = EmptyTitle
	} else {
		dl = &DigiList{
			sk:        sk,
			title:     title,
			timestamp: timestamp,
		}
	}
	return
}

func (dl *DigiList) PublicKey() *rsa.PublicKey {
	return dl.sk // may be nil
}

func (dl *DigiList) Title() string {
	return dl.title
}

func (dl *DigiList) Timestamp() int64 {
	return dl.timestamp
}

func (dl *DigiList) DigSig() []byte {
	return dl.digSig
}

// DigiListI INTERFACE //////////////////////////////////////////////

// If dl.sk is not nil, return an error if it does not match the
// public part of the key.
//
// If there are any items in the DigiList, sign it.  If this succeeds,
// any existing signature is overwritten and the public part of the key
// is written to the data structure, to dl.sk.
//
func (dl *DigiList) Sign(key *rsa.PrivateKey) (err error) {
	// XXX STUB
	return
}

// If the DigiList has been signed, verify the digital signature.
// Otherwise return false.
func Verify() (ok bool) {

	// XXx STUB
	return
}

// SERIALIZATION ////////////////////////////////////////////////////

func ParseDigiList(str string) (dl *DigiList, err error) {

	// XXX STUB
	return
}

// Serialize the DigiList, terminating each field and each item
// with a CRLF.  This is a default implementation; subclasses
// satisfying DigiListI can override.
func (dl *DigiList) String() (str string) {
	// XXX STUB
	return
}
