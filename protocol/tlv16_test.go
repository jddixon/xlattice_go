package protocol

import (
	// "fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"testing"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

func (s *XLSuite) TestConstructors(c *C) {
	rng := rnglib.MakeSimpleRNG()
	tType := uint16(rng.NextInt32(256 * 256))

	const BUF_SIZE = 16
	value := make([]byte, BUF_SIZE)
	rng.NextBytes(value)

	tlv := new(TLV16)
	err := tlv.Init(tType, nil) // illegal nil data buffer
	c.Assert(err, Not(IsNil))

	err = tlv.Init(tType, &value)
	c.Assert(tType, Equals, tlv.Type())
	c.Assert(BUF_SIZE, Equals, int(tlv.Length()))
	// XXX CAN'T COMPARE []uint8
	// c.Assert(value, Equals, *(tlv.Value()))
}

func (s *XLSuite) TestCFields(c *C) {
	rng := rnglib.MakeSimpleRNG()
	tType := uint16(rng.NextInt32(256 * 256))

	const BUF_SIZE = 16
	value := make([]byte, BUF_SIZE)
	rng.NextBytes(value)

	tlv := new(TLV16)
	tlv.Init(tType, &value)
	c.Assert(tType, Equals, tlv.Type())
	c.Assert(len(value), Equals, int(tlv.Length()))
	// XXX CAN'T COMPARE []uint8
	// c.Assert(value, Equals, *(tlv.Value()))

}

//// DEBUG
//func dumpBuffer( title string, p*[]byte) {
//    length := len(*p)
//    fmt.Printf("%s ", title)
//    for i := 0; i < length; i++ {
//        fmt.Printf("%02x ", (*p)[i])
//    }
//    fmt.Print("\n")
//}
//// END
func (s *XLSuite) TestReadWrite(c *C) {
	const COUNT = 16
	const MAX_TYPE = 128
	rng := rnglib.MakeSimpleRNG()
	for i := 0; i < COUNT; i++ {
		tType := uint16(rng.NextInt32(MAX_TYPE))
		bufLen := 4 * (1 + int(rng.NextInt32(16)))
		value := make([]byte, bufLen)
		rng.NextBytes(value) // just adds noise
		tlv := new(TLV16)
		tlv.Init(tType, &value)

		// create a buffer, write TLV16 into it, read it back
		buffer := make([]byte, 4+bufLen)
		// XXX offset should be randomized
		offset := uint16(0)
		tlv.Encode(&buffer, offset)
		decoded := Decode(&buffer, offset)
		c.Assert(tType, Equals, decoded.Type())
		c.Assert(bufLen, Equals, int(decoded.Length()))
		//      // DEBUG
		//      dumpBuffer("value   ", &buffer)
		//      dumpBuffer("encoded ", decoded.Value())
		//      // END
		for j := 0; j < bufLen; j++ {
			c.Assert((*tlv.Value())[j], Equals, (*decoded.Value())[j])
		}
	}
}
