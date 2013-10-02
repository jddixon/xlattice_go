package util

import (
	"fmt"
)

var _ = fmt.Print

type BitMap64 struct {
	Bits uint64
}

func NewBitMap64(bits uint64) (bm *BitMap64) {
	return &BitMap64{Bits: bits}
}

// An array of a bit maps with the low order N bits set.

// In a better world this would be a constant.

var lowNMap = [...]uint64{
	0x0000000000000000,
	0x0000000000000001, 0x0000000000000003, 0x0000000000000007, 0x000000000000000f,
	0x000000000000001f, 0x000000000000003f, 0x000000000000007f, 0x00000000000000ff,

	0x00000000000001ff, 0x00000000000003ff, 0x00000000000007ff, 0x0000000000000fff,
	0x0000000000001fff, 0x0000000000003fff, 0x0000000000007fff, 0x000000000000ffff,

	0x000000000001ffff, 0x000000000003ffff, 0x000000000007ffff, 0x00000000000fffff,
	0x00000000001fffff, 0x00000000003fffff, 0x00000000007fffff, 0x0000000000ffffff,

	0x0000000001ffffff, 0x0000000003ffffff, 0x0000000007ffffff, 0x000000000fffffff,
	0x000000001fffffff, 0x000000003fffffff, 0x000000007fffffff, 0x00000000ffffffff,

	0x00000001ffffffff, 0x00000003ffffffff, 0x00000007ffffffff, 0x0000000fffffffff,
	0x0000001fffffffff, 0x0000003fffffffff, 0x0000007fffffffff, 0x000000ffffffffff,

	0x000001ffffffffff, 0x000003ffffffffff, 0x000007ffffffffff, 0x00000fffffffffff,
	0x00001fffffffffff, 0x00003fffffffffff, 0x00007fffffffffff, 0x0000ffffffffffff,

	0x0001ffffffffffff, 0x0003ffffffffffff, 0x0007ffffffffffff, 0x000fffffffffffff,
	0x001fffffffffffff, 0x003fffffffffffff, 0x007fffffffffffff, 0x00ffffffffffffff,

	0x01ffffffffffffff, 0x03ffffffffffffff, 0x07ffffffffffffff, 0x0fffffffffffffff,
	0x1fffffffffffffff, 0x3fffffffffffffff, 0x7fffffffffffffff, 0xffffffffffffffff,
}

// Returns a bit map with the low order N bits set.  If N is 0, the map
// is empty - that is, no bits are set.

func LowNMap(n uint) (bm *BitMap64) {
	// XXX Constrain 0 <= n <= 63
	var u BitMap64
	bm = &u

	if n < 0 {
		n = 0
	} else if n > 64 {
		n = 64
	}

	u = BitMap64{lowNMap[n]}

	return
}

// OTHER FUNCTIONS //////////////////////////////////////////////////

// Return true if all bits are set
func (bm *BitMap64) All() bool {
	return bm.Bits == 0xffffffffffffffff
}

// Return true if any bits are set
func (bm *BitMap64) Any() bool {
	return bm.Bits != 0
}

// Clear bit N, setting it to zero.
func (bm *BitMap64) Clear(n uint) *BitMap64 {
	// XXX Constrain 0 <= n <= 63
	return &BitMap64{Bits: bm.Bits & ^(uint64(1) << n)}
}

// Set the low order N bits to zero.
func (bm *BitMap64) ClearLowN(n uint) *BitMap64 {
	// XXX Constrain 0 <= n <= 63
	return &BitMap64{Bits: bm.Bits & ^lowNMap[n+1]}
}

// Return a clone of this bit map.
func (bm *BitMap64) Clone() *BitMap64 {
	return NewBitMap64(bm.Bits)
}

// Returns a bit map in which all of the bits in this map have been flipped.
func (bm *BitMap64) Complement() *BitMap64 {
	return &BitMap64{Bits: ^bm.Bits}
}

// Returns a count of the bits set (equal to 1) in the bit map.
func (bm *BitMap64) Count() int {
	return popCount3(bm.Bits)
}

// Returns the difference between two bit maps.
func (bm *BitMap64) Difference(other *BitMap64) *BitMap64 {
	b := bm.Bits & ^other.Bits
	return &BitMap64{Bits: b}
}

// Whether 'other' is a bit map and has the same bits set.
func (bm *BitMap64) Equal(any interface{}) bool {
	if any == nil {
		return false
	}
	if any == bm {
		return true
	}
	switch v := any.(type) {
	case *BitMap64:
		_ = v
	default:
		return false
	}
	other := any.(*BitMap64)
	return bm.Bits == other.Bits
}

// Flip the Nth bit in the map, where 0 <= n <= 63
func (bm *BitMap64) Flip(n uint) *BitMap64 {
	// XXX Constrain 0 <= n <= 63
	b := bm.Bits ^ (1 << n)
	return &BitMap64{b}
}

// Return the intersection of the two bit maps -- that is, a map
// in which all of the bits set in both input sets are set.

func (bm *BitMap64) Intersection(other *BitMap64) *BitMap64 {
	b := bm.Bits & other.Bits
	return &BitMap64{b}
}

// Return whether none of the bits in the map is set.
func (bm *BitMap64) None() bool {
	return bm.Bits == 0
}

// Return a map identical to this one except that the Nth bit is set.
func (bm *BitMap64) Set(n uint) *BitMap64 {
	// XXX Constrain 0 <= n <= 63
	b := bm.Bits | (uint64(1) << n)
	return &BitMap64{b}
}

// Return a map which is the XOR of the two inputs.
func (bm *BitMap64) SymmetricDifference(other *BitMap64) *BitMap64 {
	b := bm.Bits ^ other.Bits
	return &BitMap64{b}
}

// Test the Nth bit in the map
func (bm *BitMap64) Test(n uint) bool {
	if n > 63 {
		return false
	}
	return bm.Bits&(1<<n) != 0
}

// Return a map which is the union of the two inputs -- that is,
// where all of the bits which are set in either of the two inputs
// is set in the output.
func (bm *BitMap64) Union(other *BitMap64) *BitMap64 {
	b := bm.Bits | other.Bits
	return &BitMap64{Bits: b}
}

// POP_COUNT AND SUCH ///////////////////////////////////////////////

// See Wikipedia: http://en.wikipedia.org/wiki/Hamming_weight.

const (
	m1  = 0x5555555555555555
	m2  = 0x3333333333333333
	m4  = 0x0f0f0f0f0f0f0f0f
	h01 = 0x0101010101010101
)

// Code suitable for machines with a fast multiply operation.
func popCount3(x uint64) (count int) {
	x -= (x >> 1) & m1
	x = (x & m2) + ((x >> 2) & m2)
	x = (x + (x >> 4)) & m4
	return int((x * h01) >> 56)
}

// Better for cases where few bits are non-zero
func popCount4(x uint64) (count int) {
	for count = 0; x != 0; count++ {
		x &= x - 1
	}
	return
}
