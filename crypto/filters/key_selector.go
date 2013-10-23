package filters

import (
	"fmt"
	xc "github.com/jddixon/xlattice_go/crypto"
)

var _ = fmt.Print

const (
	BITS_PER_WORD = uint(64)
	KEY_SEL_BITS  = uint(6)
)

// Go won't let these be constants
var (
	// AND with byte to expose index-many bits */
	UNMASK = []byte{
		//0  1  2  3   4   5   6    7    8
		0, 1, 3, 7, 15, 31, 63, 127, 255}
	// AND with byte to zero out index-many bits */
	MASK = []byte{
		255, 254, 252, 248, 240, 224, 192, 128, 0}
)

// Given a key, populates arrays determining word and bit offsets into
// a Bloom filter.
type KeySelector struct {
	m, k       uint
	b          []byte // key that we are inserting into the filter
	bitOffset  []byte
	wordOffset []uint
}

// Creates a key selector for a Bloom filter.  When a key is presented
// to the getOffsets() method, the k 'hash function' values are
// extracted and used to populate bitOffset and wordOffset arrays which
// specify the k flags to be set or examined in the filter.
//
// @param m    size of the filter as a power of 2
// @param k    number of 'hash functions'
// @param bitOffset array of k bit offsets (offset of flag bit in word)
// @param wordOffset array of k word offsets (offset of word flag is in)
func NewKeySelector(m, k uint, bitOffset []byte, wordOffset []uint) (
	ks *KeySelector, err error) {

	if (m < MIN_M) || (m > MAX_M) || (k < MIN_K) || (bitOffset == nil) ||
		(wordOffset == nil) {

		err = KeySelectorArgOutOfRange
	}
	if err == nil {
		ks = &KeySelector{
			m:          m,
			k:          k,
			bitOffset:  bitOffset,
			wordOffset: wordOffset,
		}
	}
	return
}

// Given a key, populate the word and bit offset arrays, each
// of which has k elements.
//
// @param key cryptographic key used in populating the arrays
///
func (ks *KeySelector) getOffsets(key []byte) (err error) {
	if key == nil {
		err = NilKey
		if err == nil && len(key) < xc.SHA3_LEN {
			err = KeyTooShort
		}
	}
	ks.b = key
	if err == nil {
		ks.GetBitSelectors()
		ks.GetWordSelectors()
	}
	return
}

// Extracts the k bit offsets from a key, suitable for general values
// of m and k.
func (ks *KeySelector) GetBitSelectors() {

	var curBit, curByte uint
	for j := uint(0); j < ks.k; j++ {
		curByte = curBit / 8
		bitsUnused := ((curByte + 1) * 8) - curBit // left in byte

		// DEBUG
		//fmt.Printf("GBS: curByte %d, curBit %d\n", curByte, curBit)
		//fmt.Printf("    this byte 0x%x, next byte 0x%x; bitsUnused %d\n",
		//	ks.b[curByte], ks.b[curByte+1], bitsUnused)
		// END

		if bitsUnused > KEY_SEL_BITS {
			// Both Java and  >> sign-extend to the right, hence the 0xff.
			// However, b is a slice of (unsigned) bytes.

			// DEBUG
			// fmt.Printf("    case %d > KEY_SEL_BITS unused bits\n", bitsUnused)
			// END
			ks.bitOffset[j] = ((0xff & ks.b[curByte]) >>
				(bitsUnused - KEY_SEL_BITS)) & UNMASK[KEY_SEL_BITS]

			// DEBUG
			//fmt.Printf("        before shifting: 0x%x\n"+
			//	"        after shifting:  0x%x\n"+
			//	"        mask:            0x%x\n",
			//	ks.b[curByte],
			//	(0xff&ks.b[curByte])>>(bitsUnused-KEY_SEL_BITS),
			//	UNMASK[KEY_SEL_BITS])
			// END

		} else if bitsUnused == KEY_SEL_BITS {
			// DEBUG
			// fmt.Printf("    case %d = KEY_SEL_BITS unused bits\n", bitsUnused)
			// END
			ks.bitOffset[j] = ks.b[curByte] & UNMASK[KEY_SEL_BITS]

		} else {
			// DEBUB
			// fmt.Printf("    case %d < KEY_SEL_BITS unused bits\n", bitsUnused)
			// END

			ks.bitOffset[j] = (ks.b[curByte] & UNMASK[bitsUnused]) |
				(((0xff & ks.b[curByte+1]) >> (8 - KEY_SEL_BITS)) &
					MASK[bitsUnused])

			//              // DEBUG
			//              fmt.Println(
			//                "    contribution from first byte:  "
			//                + itoh(b[curByte] & UNMASK[bitsUnused])
			//            + "\n    second byte: " + btoh(b[curByte + 1])
			//            + "\n    shifted:     " + itoh((0xff & b[curByte + 1]) >> 3)
			//            + "\n    mask:        " + itoh(MASK[bitsUnused])
			//            + "\n    contribution from second byte: "
			//                + itoh((0xff & b[curByte + 1] >> 3) & MASK[bitsUnused]))
			//              // END
		}
		// DEBUG
		//fmt.Printf("    ks.bitOffset[%d] = 0x%x\n", j, ks.bitOffset[j])
		// END
		//          // DEBUG
		//          fmt.Println ("    ks.bitOffset[j] = " + ks.bitOffset[j])
		//          // END
		curBit += KEY_SEL_BITS
	}
}

// Extracts the k word offsets from a key.  Suitable for general
// values of m and k.

// Extract the k offsets into the word offset array */
func (ks *KeySelector) GetWordSelectors() {
	// the word selectors being created
	selBits := ks.m - uint(6)
	selBytes := (selBits + uint(7)) / uint(8)
	bitsLastByte := selBits - uint(8)*(selBytes-uint(1))

	// DEBUG
	fmt.Printf("k %d, selBits %d, selBytes %d, bitsLastByte %d, mask 0x%x\n",
		ks.k, selBits, selBytes, bitsLastByte, UNMASK[bitsLastByte])
	// END

	// these describe the key being inserted into the filter
	lenB := len(ks.b)
	curBit := ks.k * KEY_SEL_BITS

	for i := uint(0); i < ks.k; i++ {
		curByte := curBit / 8
		var wordSel uint // accumulate selector bits here

		if curBit%8 == 0 {
			// DEBUG
			for j := uint(0); j < selBytes; j++ {
				fmt.Printf("0x%02x ", ks.b[curByte+j])
			}
			fmt.Println()
			// END

			// byte-aligned, life is easy
			for j := uint(0); j < selBytes-1; j++ {
				wordSel |= uint(ks.b[curByte]) << (j * 8)
				// DEBUG
				fmt.Printf("%d:%d 0x%02x => wordSel = 0x%06x\n",
					i, j, ks.b[curByte], wordSel)
				// END
				curByte++
			}
			// DEBUG
			lastByte := ks.b[curByte]
			maskedLastByte := lastByte & UNMASK[bitsLastByte]
			shiftedLastByte := uint(maskedLastByte) << ((selBytes - 1) * 8)
			fmt.Printf("%02x => %02x => %05x\n",
				lastByte, maskedLastByte, shiftedLastByte)
			// END
			wordSel |= (uint(ks.b[curByte] & UNMASK[bitsLastByte])) <<
				((selBytes - 1) * 8)
			// DEBUG
			fmt.Printf("%d:%d 0x%02x => wordSel = 0x%06x\n",
				i, selBytes-1, ks.b[curByte], wordSel)
			// END
		} else {
			fmt.Printf("extracting selector %d; curBit %d, curByte %d\n", curBit, curByte)
			endBit := curBit + selBits

			// first byte in b has bits on right
			bitsLeftByte := (8 * curByte) - curBit
			wordSel = uint(ks.b[curByte]) >> (8 - bitsLeftByte)
			curBit += bitsLeftByte
			wordSelBit := bitsLeftByte

			for curBit < endBit {
				curByte = curBit / 8
				var bitsThisByte uint
				if endBit-curBit >= 8 {
					bitsThisByte = 8
				} else {
					bitsThisByte = endBit - curBit
				}
				val := uint(ks.b[curByte] & UNMASK[bitsThisByte])
				wordSel |= val << wordSelBit
				wordSelBit += bitsThisByte
				curBit += bitsThisByte
			}
		}
		ks.wordOffset[i] = wordSel

		curBit += selBits
	}
	_ = lenB
}
