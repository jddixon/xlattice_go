package filters

// xlattice_go/crypto/filters/key_selector_test.go

import (
	"fmt"
	. "launchpad.net/gocheck"
)

func setUpTestKS() (
	ks *KeySelector,
	m, k uint,
	keys [][]byte,
	bOff []byte,
	wOff []uint) {

	m = 20 // default
	k = 8
	// 32 keys by default
	keys = make([][]byte, 32)
	for i := 0; i < 32; i++ {
		keys[i] = make([]byte, 20)
	}
	bOff = make([]byte, 20)
	wOff = make([]uint, 20)
	return
}

func (s *XLSuite) TestBitSelection(c *C) {

	var err error
	ks, m, k, keys, bOff, wOff := setUpTestKS()

	// set up 32 test keys
	for i := 0; i < 32; i++ {
		bitOffsets := []byte{
			byte(i % 32), byte(i + 1%32), byte(i + 2%32), byte(i + 3%32),
			byte(i + 4%32), byte(i + 5%32), byte(i + 6%32), byte(i + 7%32)}
		s.setBitOffsets(c, keys[i], bitOffsets)
	}
	ks, err = NewKeySelector(m, k, bOff, wOff) // default m=20, k=8
	c.Assert(err, IsNil)
	for i := uint(0); i < 32; i++ {
		fmt.Printf("testing, i = %d\n", i) // DEBUG
		ks.getOffsets(keys[i])
		for j := uint(0); j < k; j++ {
			fmt.Printf("        j = %d: (i + j %% 32 = %d, bOff[%d] = %d\n",
				j, (i+j)%32, j, bOff[j])
			c.Assert(byte((i+j)%32), Equals, bOff[j])
		}
	}

	_, _, _, _, _, _ = ks, m, k, keys, bOff, wOff
}

// Set the bit selectors, which are 5-bit values packed at
// the beginning of a key.
// @param b   key, expected to be at least 20 bytes long
// @param val array of key values, expected to be k long
func (s *XLSuite) setBitOffsets(c *C, b []byte, val []byte) {
	// bLen := uint(len(b))
	vLen := uint(len(val))
	var curBit, curByte uint

	for i := uint(0); i < vLen; i++ {
		curByte = curBit / 8
		offsetInByte := curBit - (curByte * 8)
		bVal := val[i] & UNMASK[5] // mask value to 5 bits

		//      // DEBUG
		//      System.out.println(
		//          "hash " + i + ": bit " + curBit + ", byte " + curByte
		//          + "; inserting " + itoh(bVal)
		//          + " into " + btoh(b[curByte]))
		//      // END
		if offsetInByte == 0 {
			// write val to left end of byte
			//b[curByte] &= 0xf1
			b[curByte] |= (bVal << 3)
			//          // DEBUG
			//          System.out.println(
			//              "    current byte becomes " + btoh(b[curByte]))
			//          // END
		} else if offsetInByte < 4 {
			// it will fit in this byte
			//b[curByte] &= ( KeySelector.MASK[5] << (3 - offsetInByte) )
			b[curByte] |= (bVal << (3 - offsetInByte))
			//          // DEBUG
			//          System.out.println(
			//              "    offsetInByte " + offsetInByte
			//          + "\n    current byte becomes " + btoh(b[curByte]))
			//          // END
		} else {
			// some goes in this byte, some in the next
			bitsThisByte := 8 - offsetInByte
			//          // DEBUG
			//          System.out.println(
			//              "SPLIT VALUE: "
			//              + "bit " + curBit + ", byte " + curByte
			//              + ", offsetInByte " + offsetInByte
			//              + ", bitsThisByte = " + bitsThisByte)
			//          // END
			valThisByte := (bVal & UNMASK[bitsThisByte])
			//b[curByte] &= MASK[bitsThisByte]
			b[curByte] |= valThisByte

			valNextByte := (bVal & MASK[bitsThisByte]) << 3

			//b[curByte+1] &= (MASK[5 - bitsThisByte]
			//                    << (3 + bitsThisByte))
			b[curByte+1] |= valNextByte
		}
		curBit += 5
	}
}
