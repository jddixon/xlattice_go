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
	//                0    1   2   3   4   5   6    7    8
    UNMASK = []byte{  0,   1,  3,  7, 15, 31, 63, 127, 255 }
    // AND with byte to zero out index-many bits */
	MASK   = []byte{255, 127, 63, 31, 15,  7,  3,   2,   1 } 
)
const (
    TWO_UP_15 = 32 * 1024
)
//type BitSelectorI interface {
//    GetBitSelectors()
//}
//type WordSelectorI interface {
//    GetWordSelectors()
//}

// Given a key, populates arrays determining word and bit offsets into
// a Bloom filter.
type KeySelector struct {
	m, k uint
	b	[]byte
    bitOffset []byte
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
func NewKeySelector (m, k uint, bitOffset []byte, wordOffset[]uint) (
	ks *KeySelector, err error) {


    if  (m < 2) || (m > 20)|| (k < 1) || (bitOffset == nil) || 
		(wordOffset == nil) {

		err = KeySelectorArgOutOfRange 
    } 
	if err == nil  {
		ks = & KeySelector{ 
		    m: m,
		    k: k,
		    bitOffset: bitOffset,
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
func (ks *KeySelector) getOffsets (key []byte) (err error) {
    if (key == nil) {
		err = NilKey
		if (err == nil && len(key) < xc.SHA3_LEN) {
			err = KeyTooShort
		}
    }
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

