package chunks

// xlattice_go/protocol/chunks/chunklist.go

import (
	"crypto/rsa"
	"io"
)

/////////////////////////////////////////////////////////////////////
// XXX THIS IS JUST WRONG.  We want information about the message being
// chunked, so that we can build the hash table.
/////////////////////////////////////////////////////////////////////

type ChunkList struct {
	length   int64
	hashes   [][]byte
	DigiList // contains sk, title, timestamp, digSig
}

// Create a ChunkList for an io.Reader where the length and SHA3-256
// content key of the document are already known.
//
func NewChunkList(sk *rsa.PublicKey, title string, timestamp int64,
	reader io.Reader, length int64, hash []byte) (
	cl *ChunkList, err error) {

	var header *Chunk // SCRATCH

	dl, err := NewDigiList(sk, title, timestamp) // checks parameters
	if err == nil {
		if reader == nil {
			err = NilReader
		} else if length == 0 {
			err = ZeroLengthInput
		} else if hash == nil {
			err = NilContentHash
		}
	}
	if err == nil {
		var hashes [][]byte
		cl = &ChunkList{
			hashes:   hashes,
			DigiList: *dl,
		}
	}
	if err == nil {
		// A packet with no data.
		packet := make([]byte, DATA_OFFSET)
		header = &Chunk{packet: packet}
		// default length is 128KB, which is stored as 128K - 1, 0x01ff
		header.setLength(MAX_CHUNK_BYTES)

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
func (cl *ChunkList) HashItem(n uint) (hash []byte, err error) {

	if n >= cl.Size() {
		err = NoNthItem
	}

	// XXX STUB

	hash = cl.hashes[n]

	return
}

func (self *ChunkList) Sign(key *rsa.PrivateKey) (err error) {
	return self.DigiList.Sign(key, self)
}

// Return the number of items currently in the DigiList.
func (cl *ChunkList) Size() uint {
	return uint(len(cl.hashes))
}

// SERIALIZATION ////////////////////////////////////////////////////

// Serialize the DigiList, terminating each field and each item
// with a CRLF.  This implementation should override the code in
// digilist.go
func (cl *ChunkList) String() (str string) {
	// XXX STUB
	return
}
