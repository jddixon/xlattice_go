package filters

// xlattice_go/crypto/filters/nibbleCounters.go

import (
	"fmt"
)

var _ = fmt.Print

const (
	NIBBLE_MASK = uint16(0xfff0)
)

/**
 * As it stands, this class is not thread-safe.  Using classes are
 * expected to provide synchronization.
 */
type NibbleCounters struct {
	counters []uint16
}

func NewNibbleCounters(filterInts uint) *NibbleCounters {
	return &NibbleCounters{
		counters: make([]uint16, filterInts*8),
	}
}

// XXX Unsynchronized
func (nc *NibbleCounters) clear() {
	for i := 0; i < len(nc.counters); i++ {
		nc.counters[i] = 0
	}
}

/**
 * Increment the nibble, ignoring any overflow
 * @param filterWord offset of 32-bit word
 * @param filterBit  offset of bit in that word (so in range 0..31)
 * @return value of nibble after operation
 */
func (nc *NibbleCounters) Inc(filterWord uint, filterBit uint) uint16 {
	counterShort := 8*filterWord + filterBit/4
	counterCell := filterBit % 4
	shiftBy := counterCell * 4
	cellValue := 0xf & (nc.counters[counterShort] >> shiftBy)
	if cellValue < 15 {
		cellValue++
	}
	// mask off the nibble and then OR new value in
	nc.counters[counterShort] &= (NIBBLE_MASK << shiftBy)
	nc.counters[counterShort] |= (cellValue << shiftBy)
	return cellValue
}

/**
 * Decrement the nibble, ignoring any overflow
 * @param filterWord offset of 32-bit word
 * @param filterBit  offset of bit in that word (so in range 0..31)
 * @return value of nibble after operation
 */
func (nc *NibbleCounters) Dec(filterWord uint, filterBit uint) uint16 {
	counterShort := 8*filterWord + filterBit/4
	counterCell := filterBit % 4
	shiftBy := counterCell * 4
	cellValue := 0xf & (nc.counters[counterShort] >> shiftBy)
	if cellValue > 0 {
		cellValue--
	}
	// mask off the nibble and then OR new value in
	nc.counters[counterShort] &= (NIBBLE_MASK << shiftBy)
	nc.counters[counterShort] |= (cellValue << shiftBy)
	return cellValue
}
