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
		ks.getBitSelectors()
		ks.getWordSelectors()
	}
	return
}

// Extracts the k bit offsets from a key, suitable for general values
// of m and k.
func (ks *KeySelector) getBitSelectors() {

	var curBit, curByte uint
	for j := uint(0); j < ks.k; j++ {
		var keySel byte
		curByte = curBit / 8
		tBit := curBit - 8*curByte // bit offset this byte
		uBits := 8 - tBit          // unused, left in byte

		if curBit%8 == 0 {
			keySel = ks.b[curByte] & UNMASK[KEY_SEL_BITS]
		} else if uBits >= KEY_SEL_BITS {
			// it's all in this byte
			keySel = (ks.b[curByte] >> tBit) & UNMASK[KEY_SEL_BITS]
		} else {
			// the selector spans two bytes
			rBits := KEY_SEL_BITS - uBits
			lSide := (ks.b[curByte] >> tBit) & UNMASK[uBits]
			rSide := (ks.b[curByte+1] & UNMASK[rBits]) << uBits
			keySel = lSide | rSide
		}
		ks.bitOffset[j] = keySel
		curBit += KEY_SEL_BITS
	}
}

// Extracts the k word offsets from a key.  Suitable for general
// values of m and k.

// Extract the k offsets into the word offset array */
func (ks *KeySelector) getWordSelectors() {
	// the word selectors being created
	selBits := ks.m - uint(6)
	selBytes := (selBits + 7) / 8
	bitsLastByte := selBits - 8*(selBytes-1)

	// bit offset into ks.b, the key being inserted into the filter
	curBit := ks.k * KEY_SEL_BITS

	for i := uint(0); i < ks.k; i++ {
		curByte := curBit / 8

		var wordSel uint // accumulate selector bits here

		if curBit%8 == 0 {
			// byte-aligned, life is easy
			for j := uint(0); j < selBytes-1; j++ {
				wordSel |= uint(ks.b[curByte]) << (j * 8)
				curByte++
			}
			wordSel |= (uint(ks.b[curByte] & UNMASK[bitsLastByte])) <<
				((selBytes - 1) * 8)
			curBit += selBits

		} else {
			endBit := curBit + selBits
			usedBits := curBit - (8 * curByte)
			wordSel = uint(ks.b[curByte]) >> usedBits
			curBit += (8 - usedBits)
			wordSelBit := 8 - usedBits

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
	}
}
