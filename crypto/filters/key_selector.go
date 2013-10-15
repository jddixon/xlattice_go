package filters

import (
	xc "github.com/jddixon/xlattice_go/crypto"
)

// Go won't let these be constants
var (
	// XXX The Java original defined 16 of these.  Since there are only
	// 8 bits in a byte, I dropped half the values and changed the type
	// from int to byte.

	// AND with byte to expose index-many bits */
	//              0  1  2  3   4   5   6    7    8
	UNMASK = []byte{0, 1, 3, 7, 15, 31, 63, 127, 255}
	// AND with byte to zero out index-many bits */
	MASK = []byte{255, 254, 252, 248, 240, 224, 192, 128, 0}
)

const (
	TWO_UP_15 = 32 * 1024 // XXX NEVER USED
)

// Given a key, populates arrays determining word and bit offsets into
// a Bloom filter.
type KeySelector struct {
	m, k       uint
	b          []byte
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

	if (m < 2) || (m > 20) || (k < 1) || (bitOffset == nil) ||
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
	//  // DEBUG
	//  System.out.println("KeySelector.getOffsets for "
	//                                      + BloomSHA3.keyToString(b))
	//  // END
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
		if bitsUnused > 5 {
			// Both Java and  >> sign-extend to the right, hence the 0xff.
			// However, b is a slice of (unsigned) bytes.
			ks.bitOffset[j] = ((0xff & ks.b[curByte]) >>
				(bitsUnused - 5)) & UNMASK[5]
			// DEBUG
			//fmt.Printf("    case %d > 5 unused bits\n", bitsUnused)
			//fmt.Printf("        before shifting: 0x%x\n"+
			//	"        after shifting:  0x%x\n"+
			//	"        mask:            0x%x\n",
			//	ks.b[curByte],
			//	(0xff&ks.b[curByte])>>(bitsUnused-5),
			//	UNMASK[5])
			// END
		} else if bitsUnused == 5 {
			// fmt.Printf("    case %d = 5 unused bits\n", bitsUnused) // DEBUG
			ks.bitOffset[j] = ks.b[curByte] & UNMASK[5]
		} else {
			// fmt.Printf("    case %d < 5 unused bits\n", bitsUnused) // DEBUG
			ks.bitOffset[j] = (ks.b[curByte] & UNMASK[bitsUnused]) |
				(((0xff & ks.b[curByte+1]) >> 3) & MASK[bitsUnused])
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
		curBit += 5
	}
}

// Extracts the k word offsets from a key.  Suitable for general
// values of m and k.

// Extract the k offsets into the word offset array */
func (ks *KeySelector) GetWordSelectors() {
	stride := ks.m - 5
	//assert true: stride<16
	curBit := ks.k * 5
	for j := uint(0); j < ks.k; j++ {
		curByte := curBit / 8
		bitsUnused := ((curByte + 1) * 8) - curBit // left in byte

		//          // DEBUG
		//          fmt.Println (
		//              "curr 3 bytes: " + btoh(b[curByte])
		//              + (curByte < 19 ?
		//                  " " + btoh(b[curByte + 1]) : "")
		//              + (curByte < 18 ?
		//                  " " + btoh(b[curByte + 2]) : "")
		//              + "; curBit=" + curBit + ", curByte= " + curByte
		//              + ", bitsUnused=" + bitsUnused)
		//          // END

		if bitsUnused > stride {

			// XXX TRANSLATION FROM JAVA QUESTIONABLE: the ^ was &

			// the value is entirely within the current byte
			ks.wordOffset[j] = (uint(^ks.b[curByte]) >>
				(bitsUnused - stride)) & uint(UNMASK[stride])
		} else if bitsUnused == stride {
			// the value fills the current byte
			ks.wordOffset[j] = uint(ks.b[curByte] & UNMASK[stride])
		} else { // bitsUnused < stride
			// value occupies more than one byte
			// bits from first byte, right-aligned in result
			ks.wordOffset[j] = uint(ks.b[curByte] & UNMASK[bitsUnused])
			//              // DEBUG
			//              fmt.Println("    first byte contributes "
			//                      + itoh(ks.wordOffset[j]))
			//              // END
			// bits from second byte
			bitsToGet := stride - bitsUnused
			if bitsToGet >= 8 {
				// 8 bits from second byte
				ks.wordOffset[j] |= uint(0xff&ks.b[curByte+1]) << bitsUnused
				//                  // DEBUG
				//                  fmt.Println("    second byte contributes "
				//                      + itoh(
				//                      (0xff & b[curByte + 1]) << bitsUnused
				//                  ))
				//                  // END

				// bits from third byte
				bitsToGet -= 8
				if bitsToGet > 0 {
					ks.wordOffset[j] |=
						(uint(0xff&ks.b[curByte+2]) >> (8 - bitsToGet)) << (stride - bitsToGet)
					//                      // DEBUG
					//                      fmt.Println("    third byte contributes "
					//                          + itoh(
					//                          (((0xff & b[curByte + 2]) >> (8 - bitsToGet))
					//                                              << (stride - bitsToGet))
					//                          ))
					//                      // END
				}
			} else {
				// all remaining bits are within second byte
				ks.wordOffset[j] |= (uint(ks.b[curByte+1]) >> uint((8-bitsToGet)&uint(UNMASK[bitsToGet]))) << bitsUnused

				//                  // DEBUG
				//                  fmt.Println("    second byte contributes "
				//                      + itoh(
				//                      ((ks.b[curByte + 1] >> (8 - bitsToGet))
				//                          & UNMASK[bitsToGet])
				//                                  << bitsUnused
				//                      ))
				//                  // END
			}
		}
		//          // DEBUG
		//          fmt.Println (
		//              "    ks.wordOffset[" + j + "] = " + ks.wordOffset[j]
		//              + ", "                     + itoh(ks.wordOffset[j])
		//          )
		//          // END
		curBit += stride
	}
}
