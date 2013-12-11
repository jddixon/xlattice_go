package chunks

// xlattice_go/protocol/chunks/chunklist.go

import (
	"crypto/rsa"
)

type ChunkList struct {
	hashes [][]byte
	DigiList
}

func NewChunkList() (sk *rsa.PublicKey, title string, timestamp int64,
	cl *ChunkList, err error) {

	dl, err := NewDigiList(sk, title, timestamp)
	if err == nil {
		var hashes [][]byte
		cl = &ChunkList{
			hashes:   hashes,
			DigiList: *dl,
		}
	}
	return
}

// Return the SHA3-256 hash of the Nth item in the DigiList.  Return an
// error if there is no such item.
//
// There may be a requirement that this interface be called in order,
// beginning at zero, without skipping any items.  It might be an error
// to call this function more than once for the Nth item.
//
// If the DigiList has been signed, a call to this function will clear
// the digital signature.
func (cl *ChunkList) HashItem(n int) (hash []byte, err error) {

	// XXX STUB
	return
}

// Return the number of items currently in the DigiList.
func (cl *ChunkList) Size() int {
	return len(cl.hashes)
}

// SERIALIZATION ////////////////////////////////////////////////////

// Serialize the DigiList, terminating each field and each item
// with a CRLF.  This implementation should override the code in
// digilist.go
func (cl *ChunkList) String() (str string) {
	// XXX STUB
	return
}
