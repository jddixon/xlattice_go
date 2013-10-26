// A Bloom filter for sets of SHA1 and SHA3 digests.  A Bloom filter uses
// a set of k hash functions to determine set membership.  Each hash function
// produces a value in the range 0..M-1.  The filter is of size M.  To
// add a member to the set, apply each function to the new member and
// set the corresponding bit in the filter.  For M very large relative
// to k, this will normally set k bits in the filter.  To check whether
// x is a member of the set, apply each of the k hash functions to x
// and check whether the corresponding bits are set in the filter.  If
// any are not set, x is definitely not a member.  If all are set, x
// may be a member.  The probability of error (the false positive rate)
// is f = (1 - e^(-kN/M))^k, where N is the number of set members.
//
// This class takes advantage of the fact that SHA1/3 digests are good-
// quality pseudo-random numbers.  The k hash functions are the values
// of distinct sets of bits taken from the 20-byte SHA1 or 32-byte SHA3
// hash.  The number of bits in the filter, M, is constrained to be a power
// of 2; M == 2^m.  The number of bits in each hash function may not
// exceed floor(m/k).
//
// This class is designed to be thread-safe, but this has not been
// exhaustively tested.
package filters

import (
	"fmt" // DEBUG
	"math"
	"sync"
)

var _ = fmt.Print

type BloomSHAI interface {
	Capacity() uint
	Clear()
	Close()
	FalsePositives() float64
	FalsePositivesN(n uint) float64
	Insert(b []byte) (err error)
	Member(b []byte) (bool, error)
	Size() uint
}

type BloomSHA struct {
	m     uint // protected final int m
	k     uint // protected final int k
	count uint

	Filter []uint64
	//ks         *KeySelector
	//wordOffset []uint
	//bitOffset  []byte

	// convenience variables
	filterBits  uint
	filterWords uint

	mu sync.RWMutex
}

// Creates a filter with 2^m bits and k 'hash functions', where
// each hash function is a portion of the 160- or 256-bit SHA hash.

// @param m determines number of bits in filter, defaults to 20
//  @param k number of hash functions, defaults to 8
func NewBloomSHA(m, k uint) (b3 *BloomSHA, err error) {

	// XXX need to devise more reasonable set of checks
	if m < MIN_M || m > MAX_M {
		err = MOutOfRange
	}
	// XXX what is this based on??
	if err == nil && (k < MIN_K || (k*m > MAX_MK_PRODUCT)) {
		// too many hash functions for filter size
		err = TooManyHashFunctions
	}
	if err == nil {
		filterBits := uint(1) << m
		filterWords := filterBits / BITS_PER_WORD
		b3 = &BloomSHA{
			m:           m,
			k:           k,
			filterBits:  filterBits,
			filterWords: filterWords,
			Filter:      make([]uint64, filterWords),
		}
		b3.doClear() // no lock
		// offsets into the filter
		//ks, err = NewKeySelector(m, k, b3.bitOffset, b3.wordOffset)
		//if err == nil {
		//	b3.ks = ks
		//} else {
		//	b3 = nil
		//}
	}
	return
}

// Creates a filter of 2^m bits, with the number of 'hash functions"
// k defaulting to 8.
func NewNewBloomSHA(m uint) (*BloomSHA, error) {
	return NewBloomSHA(m, 8)
}

// Creates a filter of 2^20 bits with k defaulting to 8.

func NewNewNewBloomSHA() (*BloomSHA, error) {
	return NewBloomSHA(20, 8)
}

// Clear the filter, unsynchronized
func (b3 *BloomSHA) doClear() {
	for i := uint(0); i < b3.filterWords; i++ {
		b3.Filter[i] = 0
	}
}

// Synchronized version */
func (b3 *BloomSHA) Clear() {
	b3.mu.Lock()
	b3.doClear()
	b3.count = 0 // jdd added 2005-02-19
	b3.mu.Unlock()
}

// Returns the number of keys which have been inserted.  This
// class (BloomSHA) does not guarantee uniqueness in any sense; if the
// same key is added N times, the number of set members reported
// will increase by N.
func (b3 *BloomSHA) Size() uint {
	b3.mu.Lock()
	defer b3.mu.Unlock()
	return b3.count
}

// Capacity returns the number of bits in the filter.

func (b3 *BloomSHA) Capacity() uint {
	return b3.filterBits
}

// Add a key to the set represented by the filter.
//
// XXX This version does not maintain 4-bit counters, it is not
// a counting Bloom filter.
func (b3 *BloomSHA) Insert(b []byte) (err error) {
	b3.mu.Lock()
	defer b3.mu.Unlock()

	ks, err := NewKeySelector(b3.m, b3.k, b)
	if err == nil {
		for i := uint(0); i < b3.k; i++ {
			b3.Filter[ks.wordOffset[i]] |= uint64(1) << ks.bitOffset[i]
		}
		b3.count++
	}
	return
}

// Returns whether a key is in the filter.
func (b3 *BloomSHA) isMember(b []byte) (whether bool, err error) {
	ks, err := NewKeySelector(b3.m, b3.k, b)
	if err == nil {
		whether = true
		for i := uint(0); i < b3.k; i++ {
			if !((b3.Filter[ks.wordOffset[i]] & (1 << ks.bitOffset[i])) != 0) {
				whether = false
				break
			}
		}
	}
	return
}

// Whether a key is in the filter.  External interface, internally
// synchronized.
//
// @param b byte array representing a key (SHA3 digest)
// @return true if b is in the filter
func (b3 *BloomSHA) Member(b []byte) (bool, error) {
	b3.mu.RLock()
	defer b3.mu.RUnlock()

	return b3.isMember(b)
}

// For n the number of set members, return approximate false positive rate.
func (b3 *BloomSHA) FalsePositivesN(n uint) float64 {
	// (1 - e(-kN/M))^k

	fK := float64(b3.k)
	fN := float64(n)
	fB := float64(b3.filterBits)
	return math.Pow((1.0 - math.Exp(-fK*fN/fB)), fK)
}

func (b3 *BloomSHA) FalsePositives() float64 {
	return b3.FalsePositivesN(b3.count)
}
func (b3 *BloomSHA) Close() {
	// a no-op
}
