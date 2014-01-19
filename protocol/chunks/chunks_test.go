package chunks

// xlattice_go/protocol/chunks/chunks_test.go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

func (s *XLSuite) TestConstants(c *C) {
	c.Assert(MAGIC_OFFSET, Equals, 0)
	c.Assert(TYPE_OFFSET, Equals, 1)
	c.Assert(RESERVED_OFFSET, Equals, 2)
	c.Assert(LENGTH_OFFSET, Equals, 8)
	c.Assert(INDEX_OFFSET, Equals, 12)
	c.Assert(DATUM_OFFSET, Equals, 16)
	c.Assert(DATA_OFFSET, Equals, 48)
}

func (s *XLSuite) TestProperties(c *C) {
	rng := xr.MakeSimpleRNG()
	_ = rng

}

func (s *XLSuite) TestChunks(c *C) {
	rng := xr.MakeSimpleRNG()

	ndx := uint32(rng.Int31())
	datum, err := xi.New(nil)
	c.Assert(err, IsNil)
	dataLen := rng.Intn(256 * 256)
	data := make([]byte, dataLen)
	rng.NextBytes(data)
	ch, err := NewChunk(datum, ndx, data)
	c.Assert(err, IsNil)
	c.Assert(ch, NotNil)

	// field checks: magic, type, reserved
	c.Assert(ch.Magic(), Equals, byte(0))
	c.Assert(ch.Type(), Equals, byte(0))
	expectedReserved := make([]byte, 6)
	c.Assert(bytes.Equal(expectedReserved, ch.Reserved()), Equals, true)

	// field checks: length, index, datum (= hash of overall message)
	c.Assert(int(ch.GetLength()), Equals, dataLen)
	c.Assert(ch.GetIndex(), Equals, ndx)
	actualDatum := ch.GetDatum()
	c.Assert(actualDatum, NotNil)
	c.Assert(bytes.Equal(actualDatum, datum.Value()), Equals, true)

	// field checks: data, chunk hash
	c.Assert(bytes.Equal(ch.GetData(), data), Equals, true)
	d := sha3.NewKeccak256()
	d.Write(ch.packet[0 : len(ch.packet)-HASH_BYTES])
	hash := d.Sum(nil)
	c.Assert(bytes.Equal(ch.GetChunkHash(), hash), Equals, true)
}
