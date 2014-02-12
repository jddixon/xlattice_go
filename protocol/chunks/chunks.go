package chunks

import (
	"code.google.com/p/go.crypto/sha3"
	"encoding/binary"
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
)

var _ = fmt.Print

const (
	MAGIC           = 0
	MAGIC_BYTES     = 1
	MAGIC_OFFSET    = 0
	TYPE            = 0
	TYPE_BYTES      = 1
	TYPE_OFFSET     = MAGIC_OFFSET + MAGIC_BYTES
	RESERVED_OFFSET = TYPE_OFFSET + TYPE_BYTES
	RESERVED_BYTES  = 6
	// We actually store the length - 1 at this offset; and the length
	// is big-endian.
	LENGTH_OFFSET = RESERVED_OFFSET + RESERVED_BYTES
	LENGTH_BYTES  = 4
	INDEX_OFFSET  = LENGTH_OFFSET + LENGTH_BYTES
	INDEX_BYTES   = 4
	DATUM_OFFSET  = INDEX_OFFSET + INDEX_BYTES
	DATUM_BYTES   = 32
	DATA_OFFSET   = DATUM_OFFSET + DATUM_BYTES
	HASH_BYTES    = xc.SHA3_LEN
	WORD_BYTES    = 16	// we pad to likely cpu cache length

	MAX_CHUNK_BYTES = 128 * 1024 // 128 KB
)

type Chunk struct {
	packet     []byte
	dataLen    uint32 // a convenience
	hashOffset uint32 // -ditto-
}

// Datum is declared a NodeID to restrict its value to certain byteslice
// lengths.
func NewChunk(datum *xi.NodeID, ndx uint32, data []byte) (
	ch *Chunk, err error) {

	if datum == nil {
		err = NilDatum
	} else if data == nil {
		err = NilData
	} else {
		msgHash := datum.Value()
		realLen := len(data)
		adjLen := ((realLen + WORD_BYTES - 1) / WORD_BYTES) * WORD_BYTES
		paddingBytes := adjLen - realLen
		packet := make([]byte, DATUM_OFFSET)
		ch = &Chunk{packet: packet}
		ch.setLength(uint32(realLen)) // length of the data part
		ch.setIndex(ndx)              // index of this chunk in overall message
		ch.packet = append(ch.packet, msgHash...)
		ch.packet = append(ch.packet, data...)
		if paddingBytes > 0 {
			padding := make([]byte, paddingBytes)
			ch.packet = append(ch.packet, padding...)
		}
		// DEBUG
		//fmt.Printf("NewChunk %d: %6d bytes, padding %2d\n",
		//	ndx, realLen, paddingBytes)
		// END
		// calculate the SHA3-256 hash of the chunk
		d := sha3.NewKeccak256()
		d.Write(ch.packet)
		chunkHash := d.Sum(nil)

		// append that to the packet
		ch.packet = append(ch.packet, chunkHash...)
	}
	return
}

func (ch *Chunk) Magic() byte {
	return ch.packet[MAGIC_OFFSET]
}

func (ch *Chunk) Type() byte {
	return ch.packet[TYPE_OFFSET]
}

func (ch *Chunk) Reserved() []byte {
	return ch.packet[RESERVED_OFFSET : RESERVED_OFFSET+RESERVED_BYTES]
}

// Return the length encoded in the packet.  This is the actual length
// of the data in bytes, excluding any padding added.  The value actually
// stored is the length less one.
//
func (ch *Chunk) GetLength() uint32 {
	return binary.BigEndian.Uint32(
		ch.packet[LENGTH_OFFSET:LENGTH_OFFSET+LENGTH_BYTES]) + 1
}

// We store the length - 1, so it is a serious error if the length is zero.
func (ch *Chunk) setLength(n uint32) {
	binary.BigEndian.PutUint32(
		ch.packet[LENGTH_OFFSET:LENGTH_OFFSET+LENGTH_BYTES], n-1)
}

// We store the actual value of the zero-based index.
func (ch *Chunk) GetIndex() uint32 {
	return binary.BigEndian.Uint32(
		ch.packet[INDEX_OFFSET : INDEX_OFFSET+INDEX_BYTES])
}

func (ch *Chunk) setIndex(n uint32) {
	binary.BigEndian.PutUint32(
		ch.packet[INDEX_OFFSET:INDEX_OFFSET+INDEX_BYTES], n)
}

// Given a byte slice, determine the length of a chunk wrapping it:
// /header + data + chunk hash.
func (ch *Chunk) CalculateLength(data []byte) uint32 {
	dataLen := ((len(data) + WORD_BYTES - 1) / WORD_BYTES) * WORD_BYTES
	return uint32(DATA_OFFSET + dataLen + HASH_BYTES)
}

func (ch *Chunk) GetDatum() []byte {
	return ch.packet[DATUM_OFFSET : DATUM_OFFSET+DATUM_BYTES]
}

func (ch *Chunk) GetData() []byte {
	return ch.packet[DATA_OFFSET : DATA_OFFSET+ch.GetLength()]
}
func (ch *Chunk) GetChunkHash() []byte {
	if ch.dataLen == 0 {
		ch.dataLen = ch.GetLength()
		ch.hashOffset = ((ch.dataLen + DATA_OFFSET + WORD_BYTES - 1) /
			WORD_BYTES) * WORD_BYTES
	}
	return ch.packet[ch.hashOffset : ch.hashOffset+HASH_BYTES]
}
