package chunks

import (
	"code.google.com/p/go.crypto/sha3"
	"encoding/binary"
	xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
)

const (
	MAGIC           = 0
	MAGIC_BYTES     = 1
	MAGIC_OFFSET    = 0
	TYPE            = 0
	TYPE_BYTES      = 1
	TYPE_OFFSET     = MAGIC_OFFSET + MAGIC_BYTES
	RESERVED_OFFSET = TYPE_OFFSET + TYPE_BYTES
	RESERVED_BYTES  = 6
	LENGTH_OFFSET   = RESERVED_OFFSET + RESERVED_BYTES
	LENGTH_BYTES    = 4
	INDEX_OFFSET    = LENGTH_OFFSET + LENGTH_BYTES
	INDEX_BYTES     = 4
	MAX_LENGTH      = 2 * 256 * 256 // exclusive
	DATUM_OFFSET    = INDEX_OFFSET + INDEX_BYTES
	DATUM_LENGTH    = 32
	DATA_OFFSET     = DATUM_OFFSET + DATUM_LENGTH
	HASH_BYTES      = xc.SHA3_LEN
	WORD_BYTES      = 16
)

type Chunk struct {
	chunk []byte
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
		id := datum.Value()
		realLen := len(data)
		adjLen := ((realLen + WORD_BYTES - 1) / WORD_BYTES) * WORD_BYTES
		paddingBytes := adjLen - realLen
		packet := make([]byte, DATUM_OFFSET)
		ch = &Chunk{packet}
		ch.setLength(uint32(realLen)) // length of the data part
		ch.setIndex(ndx)              // index of this chunk in overall message
		packet = append(packet, id...)
		packet = append(packet, data...)
		if paddingBytes > 0 {
			padding := make([]byte, paddingBytes)
			packet = append(packet, padding...)
		}
		// calculate the SHA3-256 hash of the chunk
		d := sha3.NewKeccak256()
		d.Write(packet)
		chunkHash := d.Sum(nil)

		// append that to the packet
		packet = append(packet, chunkHash...)
	}
	return
}

func (ch *Chunk) Magic() byte {
	return ch.chunk[MAGIC_OFFSET]
}

func (ch *Chunk) Type() byte {
	return ch.chunk[TYPE_OFFSET]
}

func (ch *Chunk) Reserved() []byte {
	return ch.chunk[RESERVED_OFFSET : RESERVED_OFFSET+RESERVED_BYTES]
}

func (ch *Chunk) GetLength() uint32 {
	return binary.LittleEndian.Uint32(
		ch.chunk[LENGTH_OFFSET : LENGTH_OFFSET+LENGTH_BYTES])
}

func (ch *Chunk) setLength(n uint32) {
	binary.LittleEndian.PutUint32(
		ch.chunk[LENGTH_OFFSET:LENGTH_OFFSET+LENGTH_BYTES], n)
}

func (ch *Chunk) GetIndex() uint32 {
	return binary.LittleEndian.Uint32(
		ch.chunk[INDEX_OFFSET : INDEX_OFFSET+INDEX_BYTES])
}

func (ch *Chunk) setIndex(n uint32) {
	binary.LittleEndian.PutUint32(
		ch.chunk[INDEX_OFFSET:INDEX_OFFSET+INDEX_BYTES], n)
}

// Given a byte slice, determine the length of a chunk wrapping it.
func (ch *Chunk) CalculateLength(data []byte) uint32 {
	dataLen := ((len(data) + WORD_BYTES - 1) / WORD_BYTES) * WORD_BYTES
	return uint32(DATA_OFFSET + dataLen + HASH_BYTES)
}
