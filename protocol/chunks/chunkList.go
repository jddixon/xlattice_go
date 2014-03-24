package chunks

// xlattice_go/protocol/chunks/chunkList.go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"github.com/jddixon/xlattice_go/u"
	"github.com/jddixon/xlattice_go/util"
	"io"
)

var _ = fmt.Print

type ChunkList struct {
	length   int64
	hashes   [][]byte
	DigiList // contains sk, title, timestamp, digSig
}

// Create a ChunkList for an io.Reader where the length and SHA3-256
// content key of the document are already known.  If uStore is not nil,
// each chunk is inserted into uStore.  Chunks are CHUNK_BYTES in size,
// except that the last chunk may be shorter.  The last chunk is padded
// out to a multiple of WORD_BYTES to ensure alignment.  Padding bytes
// are sent over the wire and are included in the chunk hash but are
// not considered in calculating datum, the message hash.
//
func NewChunkList(sk *rsa.PublicKey, title string, timestamp util.Timestamp,
	reader io.Reader, length int64, datum []byte, uStore u.UI) (
	cl *ChunkList, err error) {

	var (
		dl     *DigiList
		header *Chunk // SCRATCH
	)
	chunkCount := uint32((length + MAX_DATA_BYTES - 1) / MAX_DATA_BYTES)
	bigD := sha3.NewKeccak256() // used to check datum
	hashes := make([][]byte, chunkCount)
	// DEBUG
	fmt.Printf("NewChunkList: file len %d; there will be %d chunk(s)\n",
		length, chunkCount)
	fmt.Printf("  number of hashes = %d\n", len(hashes)) // XXX
	// END

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

		stillToGo := length // bytes left unread at this point
		eofSeen := false
		for i := uint32(0); i < chunkCount && !eofSeen && err == nil; i++ {
			// Use a packet with no data as a scratch pad to build dummy headers
			hPacket := make([]byte, DATUM_OFFSET)
			hPacket = append(hPacket, datum...)
			header = &Chunk{packet: hPacket}
			// default length is 128KB - 80 = MAX_DATA_BYTES
			header.setLength(MAX_DATA_BYTES)

			var paddingBytes int
			header.setIndex(i)
			if i == chunkCount-1 {
				header.setLength(uint32(stillToGo))
			}
			var bytesToRead int64
			var count int

			if stillToGo <= MAX_DATA_BYTES {
				bytesToRead = stillToGo
			} else {
				bytesToRead = MAX_DATA_BYTES
			}
			data := make([]byte, bytesToRead)

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
			if bytesToRead != MAX_DATA_BYTES {
				adjLen := WORD_BYTES * ((bytesToRead + WORD_BYTES - 1) /
					WORD_BYTES)
				paddingBytes = int(adjLen - bytesToRead)
			}
			stillToGo -= int64(count) // ASSUMES NO PARTIAL READ

			d := sha3.NewKeccak256()
			d.Write(header.packet) // <-- header is included
			bigD.Write(data)
			if paddingBytes > 0 {
				padding := make([]byte, paddingBytes) // null bytes
				data = append(data, padding...)
			}
			d.Write(data)
			hashes[i] = d.Sum(nil)
			// DEBUG
			fmt.Printf("NewChunkList %d: %6d bytes, %2d bytes padding, %7d data len\n",
				i, bytesToRead, paddingBytes, len(data))
			fmt.Printf("   ==> %s\n", hex.EncodeToString(hashes[i]))
			// END

			if uStore != nil {
				// DEBUG
				fmt.Printf("writing chunk %d to store\n", i)
				// END

				// XXX FIX ME: WASTEFUL COPYING OF BYTES
				header.packet = append(header.packet, data...)
				bytesWritten, writeHash, err := uStore.PutData(header.packet, hashes[i])
				if err == nil {
					if bytesWritten != int64(len(header.packet)) {
						// DEBUG
						fmt.Printf("chunk length %d but only %d bytes written\n",
							len(header.packet), bytesWritten)
						// END
						err = WrongNumberBytesWritten
					} else if !bytes.Equal(writeHash, hashes[i]) {
						err = WriteReturnsWrongHash
					}
				}
			}
			// DEBUG
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
			// END
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
func (cl *ChunkList) HashItem(n uint) (hash []byte, err error) {

	// DEBUG
	// fmt.Printf("HashItem(%d) where size is %d\n", n, cl.Size())
	// END
	if n >= cl.Size() {
		err = NoNthItem
	} else {
		hash = cl.hashes[n]
	}
	return
}

func (self *ChunkList) Sign(key *rsa.PrivateKey) (err error) {
	return self.DigiList.Sign(key, self)
}

// Return the number of items currently in the DigiList.
func (cl *ChunkList) Size() uint {
	return uint(len(cl.hashes))
}

// Return nil if verification succeeds, otherwise the error encountered.
func (self *ChunkList) Verify() (err error) {
	return self.DigiList.Verify(self)
}

// SERIALIZATION ////////////////////////////////////////////////////

// Serialize the DigiList, terminating each field and each item
// with a CRLF.
func (cl *ChunkList) String() (str string) {

	// XXX STUB
	str = "CHUNK_LIST STRING()"

	return
}
