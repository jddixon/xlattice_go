package filters

import (
	"fmt"
)

// Extracts the k bit offsets from a key, suitable for general values
// of m and k.
func (ks *KeySelector) GetBitSelectors() {

	var curBit, curByte uint
	for j := uint(0); j < ks.k; j++ {
		curByte = curBit / 8
		bitsUnused := ((curByte + 1) * 8) - curBit // left in byte

		// DEBUG
		fmt.Printf("GBS: curByte %d, curBit %d\n", curByte, curBit)
		fmt.Printf("    this byte 0x%x, next byte 0x%x; bitsUnused %d\n",
			ks.b[curByte], ks.b[curByte+1], bitsUnused)
		// END
		if bitsUnused > 5 {
			// Both Java and  >> sign-extend to the right, hence the 0xff.
			// However, b is a slice of (unsigned) bytes.
			ks.bitOffset[j] = ((0xff & ks.b[curByte]) >>
				(bitsUnused - 5)) & UNMASK[5]
			// DEBUG
			fmt.Printf("    case %d > 5 unused bits\n", bitsUnused)
			fmt.Printf("        before shifting: 0x%x\n"+
				"        after shifting:  0x%x\n"+
				"        mask:            0x%x\n",
				ks.b[curByte],
				(0xff&ks.b[curByte])>>(bitsUnused-5),
				UNMASK[5])
			// END
		} else if bitsUnused == 5 {
			fmt.Printf("    case %d = 5 unused bits\n", bitsUnused) // DEBUG
			ks.bitOffset[j] = ks.b[curByte] & UNMASK[5]
		} else {
			fmt.Printf("    case %d < 5 unused bits\n", bitsUnused) // DEBUG
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
		fmt.Printf("    ks.bitOffset[%d] = 0x%x\n", j, ks.bitOffset[j]) //DEBUG
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
