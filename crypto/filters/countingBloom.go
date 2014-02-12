package filters

// xlattice_go/crypto/filters/countingBloom.go

import (
	"sync"
)

/**
 * Counting version of the Bloom filter.  Adds a 4-bit counter to each
 * bit in the Bloom filter, enabling members to be removed from the set
 * without having to recreate the filter from scratch.
 */

type CountingBloom struct {
	nibCounter NibbleCounters
	cbMU       sync.RWMutex
	BloomSHA
}

func NewCountingBloom(m, k uint) (cb *CountingBloom, err error) {

	var (
		b3 *BloomSHA
		nc *NibbleCounters
	)
	b3, err = NewBloomSHA(m, k)
	if err == nil {
		nc = NewNibbleCounters(b3.filterWords)
		cb = &CountingBloom{
			nibCounter: *nc,
			BloomSHA:   *b3,
		}
	}
	return
}
func NewNewCountingBloom(m uint) (*CountingBloom, error) {
	return NewCountingBloom(m, 8)
}
func NewNewNewCountingBloom() (*CountingBloom, error) {
	return NewCountingBloom(20, 8)
}

/**
 * Clear both the underlying filter in the superclass and the
 * bit counters maintained here.
 *
 * XXX Possible deadlock.
 */
func (cb *CountingBloom) Clear() {
	// XXX ORDER IN WHICH LOCKS ARE OBTAINED MUST BE THE SAME EVERYWHERE.
	cb.cbMU.Lock()
	defer cb.cbMU.Unlock()
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.doClear()          // the BloomSHA1; otherwise unsynchronized
	cb.nibCounter.clear() // the nibble counters; otherwise unsynchronized
}

/**
 * Add a key to the set represented by the filter, updating counters
 * as it does so.  Overflows are silently ignored.
 *
 * @param b byte array representing a key (SHA1 digest)
 */
func (cb *CountingBloom) Insert(b []byte) {
	cb.cbMU.Lock()
	defer cb.cbMU.Unlock()
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// XXX copied from BloomSHA.Insert
	ks, err := NewKeySelector(cb.m, cb.k, b)
	if err == nil {
		for i := uint(0); i < cb.k; i++ {
			cb.Filter[ks.wordOffset[i]] |= uint64(1) << ks.bitOffset[i]
		}
		cb.count++

		// XXX DOESN'T CALL Dec !
	}
}

/**
 * Remove a key from the set, updating counters while doing so.
 * If the key is not a member of the set, no action is taken.
 * However, if it is a member (a) the count is decremented,
 * (b) all bit counters are decremented, and (c) where the bit
 * counter goes to zero the corresponding bit in the filter is
 * zeroed. [No change in code, but jdd clarified these comments
 * 2005-03-29.]
 *
 * @param b byte array representing the key to be removed.
 */
func (cb *CountingBloom) Remove(b []byte) {
	cb.cbMU.Lock()
	defer cb.cbMU.Unlock()
	cb.mu.Lock()
	defer cb.mu.Unlock()

	// XXX IGNORING POSSIBLE ERROR
	present, ks, _ := cb.IsMember(b) // calls ks.getOffsets

	if present {
		for i := uint(0); i < cb.k; i++ {
			newCount := cb.nibCounter.Dec(
				ks.wordOffset[i], uint(ks.bitOffset[i]))
			if newCount == 0 {
				cb.Filter[ks.wordOffset[i]] &= ^(1 << ks.bitOffset[i])
			}
		}
		cb.count--
	}
}
