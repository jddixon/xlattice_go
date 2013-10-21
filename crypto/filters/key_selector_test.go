package filters

// xlattice_go/crypto/filters/key_selector_test.go

import (
	"fmt" // DEBUG
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

const (
	NUM_TEST_KEYS = 16
)

func setUpTestKS() (
	ks *KeySelector,
	m,
	k,
	v uint, // size of test values (160 = SHA1, 256 = SHA3)
	keys [][]byte,
	bOff []byte,
	wOff []uint) {

	m = 20 // default
	k = 8
	v = 20 // v is the number of bytes in a test value (key)

	// Create the array of keys to be used to test the KeySelector.
	// These are v=20 byte keys, so SHA1 hashes.
	keys = make([][]byte, NUM_TEST_KEYS)
	for i := 0; i < NUM_TEST_KEYS; i++ {
		keys[i] = make([]byte, v)
	}
	// Each of the k hash functions selects a bit in the filter;
	// this bit will be located first by its word offset within the
	// filter of 2 << m bits and then by its bit offset within
	// that word.
	bOff = make([]byte, k)
	wOff = make([]uint, k)
	return
}

func (s *XLSuite) TestBitSelection(c *C) {

	var err error
	ks, m, k, v, keys, bOff, wOff := setUpTestKS()

	// Set up bit selectors for NUM_TEST_KEYS test keys.  Each
	// bit selector is populated with a distinct value.
	for i := 0; i < NUM_TEST_KEYS; i++ {
		bitOffsets := make([]byte, k)
		for j := uint(0); j < k; j++ {
			bitOffsets[j] = byte((j*k + j) % 8)
		}
		s.setBitOffsets(c, &keys[i], bitOffsets)
	}
	ks, err = NewKeySelector(m, k, bOff, wOff)
	c.Assert(err, IsNil)
	for i := uint(0); i < NUM_TEST_KEYS; i++ {
		ks.getOffsets(keys[i])
		for j := uint(0); j < k; j++ {
			// DEBUG
			if bOff[j] != byte((j*k+j)%8) {
				fmt.Printf("i = %d, j = %d, actual %x, expected %x\n",
					i, j, bOff[j], byte((j*k+j)%8))
			}
			// END
			c.Assert(bOff[j], Equals, byte((j*k+j)%8))
		}
	}

	_ = v
}

// Set the bit selectors, which are the k KEY_SEL_BITS-bit values
// at the beginning of a key.
// @param b   key, expected to be at least 20 bytes long
// @param val array of key values, expected to be k long
func (s *XLSuite) setBitOffsets(c *C, b *[]byte, val []byte) {

	vLen := uint(len(val))
	var curBit, curByte uint

	for i := uint(0); i < vLen; i++ {
		curByte = curBit / 8                   // byte offset in b
		offsetInByte := curBit - (curByte * 8) // bit offset
		// mask value to KEY_SEL_BITS bits
		bVal := val[i] & UNMASK[KEY_SEL_BITS]

		//      // DEBUG
		//      System.out.println(
		//          "hash " + i + ": bit " + curBit + ", byte " + curByte
		//          + "; inserting " + itoh(bVal)
		//          + " into " + btoh(b[curByte]))
		//      // END

		if offsetInByte == 0 {
			// write val to left end of byte
			(*b)[curByte] |= bVal << (8 - KEY_SEL_BITS)

			//          // DEBUG
			//          System.out.println(
			//              "    current byte becomes " + btoh(b[curByte]))
			//          // END

		} else if offsetInByte < (8 - KEY_SEL_BITS + 1) {
			// it will fit in this byte
			//b[curByte] &= ( KeySelector.MASK[KEY_SEL_BITS] << ((8 - KEY_SEL_BITS) - offsetInByte) )
			(*b)[curByte] |= (bVal << ((8 - KEY_SEL_BITS) - offsetInByte))

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
			(*b)[curByte] |= valThisByte

			valNextByte := (bVal & MASK[bitsThisByte]) << (8 - KEY_SEL_BITS)

			//b[curByte+1] &= (MASK[KEY_SEL_BITS - bitsThisByte]
			//                    << ((8 - KEY_SEL_BITS) + bitsThisByte))
			(*b)[curByte+1] |= valNextByte
		}
		curBit += KEY_SEL_BITS
	}
}
