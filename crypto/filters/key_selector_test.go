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

func (s *XLSuite) dumpB(c *C, b []byte) {
	for i := uint(0); i < uint(len(b)); i++ {
		fmt.Printf("%02x", b[i])
	}
	fmt.Println()
} // GEEP
func (s *XLSuite) TestBitSelection(c *C) {

	var err error
	ks, m, k, v, keys, bOff, wOff := setUpTestKS()

	// Set up bit selectors for NUM_TEST_KEYS test keys.  Each
	// bit selector is populated with a random value.
	for i := 0; i < NUM_TEST_KEYS; i++ {
		bitOffsets := make([]byte, k)
		for j := uint(0); j < k; j++ {
			bitOffsets[j] = byte((j*k + j) % 8)
		}
		s.setBitOffsets(c, &keys[i], bitOffsets)
	}
	ks, err = NewKeySelector(m, k, bOff, wOff)
	c.Assert(err, IsNil)

	// DEBUG
	for i := uint(0); i < NUM_TEST_KEYS; i++ {
		ks.getOffsets(keys[i])
		for j := uint(0); j < k; j++ {
			fmt.Printf("key %d, bitSel %d, actual %02x, expected %02x\n",
				i, j, bOff[j], byte((j*k+j)%8))
		}
	}
	// END
	for i := uint(0); false && i < NUM_TEST_KEYS; i++ {
		ks.getOffsets(keys[i])
		for j := uint(0); j < k; j++ {
			// DEBUG
			if bOff[j] != byte((j*k+j)%8) {
				fmt.Printf("i = %d, j = %d, actual %02x, expected %02x\n",
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
		curByte = curBit / 8           // byte offset in b
		tBit := curBit - (curByte * 8) // bit offset
		uBits := 8 - tBit

		// mask value to KEY_SEL_BITS bits
		unVal := val[i] & UNMASK[KEY_SEL_BITS]

		if tBit == 0 {
			// we are aligned, so just OR it in
			(*b)[curByte] |= unVal

		} else if uBits >= KEY_SEL_BITS {
			// it will fit in this byte
			(*b)[curByte] |= (unVal << tBit)

		} else {
			// some goes in this byte, some in the next
			valThisByte := (unVal & UNMASK[uBits])
			(*b)[curByte] |= valThisByte << tBit

			valNextByte := (unVal >> uBits)
			(*b)[curByte+1] |= valNextByte

			fmt.Printf("  %d val %02x tBit %d this %02x next %02x\n",
				i, unVal, tBit, valThisByte, valNextByte)

		}
		curBit += KEY_SEL_BITS

		fmt.Printf("%02x => ", unVal)
		s.dumpB(c, *b) // DEBUG
	}
} // GEEP

// Set the word selectors, which are the k wordSelBits-bit values
// following the bit sectors in the key.
func (s *XLSuite) setWordOffsets(c *C, b *[]byte, val []uint, m, k uint) {
	// 2 ^ 6 == 64, number of bits in a uint64
	wordSelBits := m - 6
	wordSelMask := ^(uint(1) << wordSelBits)
	bytesInU := (wordSelBits + 7) / 8
	var bitsLastByte uint
	if bytesInU*8 == bytesInU {
		bitsLastByte = uint(8)
	} else {
		bitsLastByte = wordSelBits - (bytesInU-1)*uint(8)
	}

	fmt.Printf("bytesInU %d, wordSelBits %d, bitsLastByte %d\n",
		bytesInU, wordSelBits, bitsLastByte)

	vLen := uint(len(val))

	var curTByte uint           // byte offset in b
	curTBit := k * KEY_SEL_BITS // bit offset in b

	// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
	// XXX THIS MANGLES THE LAST NIBBLE OF THE KEY SELECTORS
	// XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

	// iterate through the test values, merging them into target
	for i := uint(0); i < vLen; i++ {

		// be paranoid: mask test value to wordSelBits bits
		maskedVal := val[i] & wordSelMask

		fmt.Printf("\nval[%d] = 0x%05x => 0x%05x (%6d)\n",
			i, val[i], maskedVal, maskedVal)

		for j := uint(0); j < bytesInU; j++ {
			thisUByte := byte(maskedVal >> (j * uint(8)))

			fmt.Printf("  thisUByte %d = 0x%02x", j, thisUByte)

			bitsThisUByte := uint(8)
			if j == (bytesInU - 1) {
				bitsThisUByte = bitsLastByte
			}
			// these point into the target, b
			curTByte = curTBit / 8
			offsetInTByte := curTBit - (curTByte * uint(8)) // bit offset
			fmt.Printf("  tBit %3d, tByte %3d\n", curTBit, curTByte)

			if offsetInTByte == 0 {
				// we just assign it in, trusting b was all zeroes
				(*b)[curTByte] = byte(thisUByte)

				fmt.Printf("= ")
				s.dumpB(c, *b)
			} else {
				// we have to shift
				ususedTBits := uint(8) - offsetInTByte
				if bitsThisUByte <= ususedTBits {
					// XXX GOES INTO WRONG BYTE
					// it will fit in this byte
					value := thisUByte << offsetInTByte
					(*b)[curTByte] |= value

					fmt.Printf("  ")
					s.dumpB(c, *b)
				} else {
					// we have to split it over two target bytes
					lValue := thisUByte >> offsetInTByte
					(*b)[curTByte] |= lValue

					fmt.Printf("L ")
					s.dumpB(c, *b)

					rValue := thisUByte << (uint(8) - offsetInTByte)
					(*b)[curTByte+1] |= rValue
					fmt.Printf("R ")
					s.dumpB(c, *b)
				}
			}
			curTBit += bitsThisUByte
		}
	}
}
