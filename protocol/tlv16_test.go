package protocol

import (
	// "fmt"
	"github.com/bmizerany/assert"
	// "github.com/jddixon/xlattice_go"
	x ".." // accepted and properly interpreted
	// "github.com/jddixon/xlattice_go/rnglib"
	"testing"
	// "time"
)

// We now use this function from ../make_rng.go

//func MakeRNG() *rnglib.SimpleRNG {
//	t := time.Now().Unix()
//	rng := rnglib.NewSimpleRNG(t)
//	return rng
//}

func TestConstructors(t *testing.T) {
	rng := x.MakeRNG()
	tType := uint16(rng.NextInt32(256 * 256))

	const BUF_SIZE = 16
	value := make([]byte, BUF_SIZE)
	rng.NextBytes(&value)

	tlv := new(TLV16)
	err := tlv.Init(tType, nil) // illegal nil data buffer
	assert.NotEqual(t, nil, err)

	err = tlv.Init(tType, &value)
	assert.Equal(t, tType, tlv.Type())
	assert.Equal(t, BUF_SIZE, int(tlv.Length()))
	assert.Equal(t, value, *(tlv.Value()))
}

func TestFields(t *testing.T) {
	rng := x.MakeRNG()
	tType := uint16(rng.NextInt32(256 * 256))

	const BUF_SIZE = 16
	value := make([]byte, BUF_SIZE)
	rng.NextBytes(&value)

	tlv := new(TLV16)
	tlv.Init(tType, &value)
	assert.Equal(t, tType, tlv.Type())
	assert.Equal(t, len(value), int(tlv.Length()))
	assert.Equal(t, value, *(tlv.Value()))

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
func TestReadWrite(t *testing.T) {
	const COUNT = 16
	const MAX_TYPE = 128
	rng := x.MakeRNG()
	for i := 0; i < COUNT; i++ {
		tType := uint16(rng.NextInt32(MAX_TYPE))
		bufLen := 4 * (1 + int(rng.NextInt32(16)))
		value := make([]byte, bufLen)
		rng.NextBytes(&value) // just adds noise
		tlv := new(TLV16)
		tlv.Init(tType, &value)

		// create a buffer, write TLV16 into it, read it back
		buffer := make([]byte, 4+bufLen)
		// XXX offset should be randomized
		offset := uint16(0)
		tlv.Encode(&buffer, offset)
		decoded := Decode(&buffer, offset)
		assert.Equal(t, tType, decoded.Type())
		assert.Equal(t, bufLen, int(decoded.Length()))
		//      // DEBUG
		//      dumpBuffer("value   ", &buffer)
		//      dumpBuffer("encoded ", decoded.Value())
		//      // END
		for j := 0; j < bufLen; j++ {
			assert.Equal(t, (*tlv.Value())[j], (*decoded.Value())[j])
		}
	}
}
