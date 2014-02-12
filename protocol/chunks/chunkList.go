package chunks

// xlattice_go/protocol/chunks/chunkList.go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/rsa"
	// "encoding/hex"		// DEBUG
	"fmt"
	"io"
)

var _ = fmt.Print

type ChunkList struct {
	length   int64
	hashes   [][]byte
	DigiList // contains sk, title, timestamp, digSig
}

// Create a ChunkList for an io.Reader where the length and SHA3-256
// content key of the document are already known.
//
func NewChunkList(sk *rsa.PublicKey, title string, timestamp int64,
	reader io.Reader, length int64, datum []byte) (
	cl *ChunkList, err error) {

	var (
		dl     *DigiList
		header *Chunk // SCRATCH
	)
	chunkCount := uint32((length + MAX_CHUNK_BYTES - 1) / MAX_CHUNK_BYTES)
	bigD := sha3.NewKeccak256() // used to check datum
	hashes := make([][]byte, chunkCount)

	if reader == nil {
		err = NilReader
	} else if length == 0 {
		err = ZeroLengthInput
	} else if datum == nil {
		err = NilDatum
	} else if len(datum) != DATUM_BYTES {
		err = BadDatumLength
	} else {
		dl, err = NewDigiList(sk, title, timestamp) // checks parameters
	}
	if err == nil {

		// We use a packet with no data as a scratch pad to build dummy headers
		hPacket := make([]byte, DATUM_OFFSET)
		hPacket = append(hPacket, datum...)
		header = &Chunk{packet: hPacket}
		// default length is 128KB, which is stored as 128K - 1, 0x01ff
		header.setLength(MAX_CHUNK_BYTES)

		stillToGo := length // bytes left unread at this point
		eofSeen := false
		for i := uint32(0); i < chunkCount && !eofSeen; i++ {
			var paddingBytes int
			header.setIndex(i)
			if i == chunkCount-1 {
				header.setLength(uint32(stillToGo))
			}
			var bytesToRead int64
			var count int
			data := make([]byte, MAX_CHUNK_BYTES)

			if stillToGo <= MAX_CHUNK_BYTES {
				bytesToRead = stillToGo
			} else {
				bytesToRead = MAX_CHUNK_BYTES
			}
			// XXX DOES NOT ALLOW FOR PARTIAL READS
			count, err = reader.Read(data)
			if err != nil {
				if err == io.EOF {
					err = nil
					eofSeen = true
				} else {
					break
				}
			}
			if bytesToRead != MAX_CHUNK_BYTES {
				data = data[0:bytesToRead]
				adjLen := WORD_BYTES * ((bytesToRead + WORD_BYTES - 1) /
					WORD_BYTES)
				paddingBytes = int(adjLen - bytesToRead)
			}
			stillToGo -= int64(count) // ASSUMES NO PARTIAL READ

			d := sha3.NewKeccak256()
			d.Write(header.packet) // <-- header is included
			bigD.Write(data)
			if paddingBytes > 0 {
				padding := make([]byte, paddingBytes)
				data = append(data, padding...)
			}
			d.Write(data)
			hashes[i] = d.Sum(nil)
			//// DEBUG
			//fmt.Printf("NewChunkList %d: %6d bytes, %2d bytes padding\n",
			//	i, bytesToRead, paddingBytes)
			//fmt.Printf("   ==> %s\n", hex.EncodeToString(hashes[i]))
			//if i == chunkCount - 1 {
			//	fmt.Println("DUMP OF LAST CHUNK" )
			//	for j := 0; j < len(header.packet); j += 16 {
			//		fmt.Printf("%6d %s\n", j,
			//			hex.EncodeToString(header.packet[j: j+16]))
			//	}
			//	for j := 0; j < len(data); j += 16 {
			//		fmt.Printf("%6d %s\n", j + 48,
			//			hex.EncodeToString(data[j: j+16]))
			//	}
			//}
			//// END
		}
	}
	if err == nil {
		contentHash := bigD.Sum(nil)
		// DEBUG
		// fmt.Printf("datum2: %s\n", hex.EncodeToString(contentHash))
		// END
		if !bytes.Equal(contentHash, datum) {
			err = BadDatum
		}
	}
	if err == nil {
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
// digiList.go
func (cl *ChunkList) String() (str string) {
	// XXX STUB
	return
}
