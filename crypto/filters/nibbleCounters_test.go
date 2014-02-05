package filters

// xlattice_go/crypto/filters/nibbleCounters_test.go

import (
	"fmt" // DEBUG
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

/**
 * Tests the counters associated with Bloom filters for sets whose members
 * are 20-byte SHA1 digests.
 */
func (s *XLSuite) doTestBit(c *C, nibCount *NibbleCounters,
	filterWord uint, filterBit uint) {

	for i := uint16(0); i < 18; i++ {
		value := nibCount.Inc(filterWord, filterBit)
		if i < 15 {
			c.Check(value, Equals, i+1) // XXX SHOULD BE Assert
		} else {
			c.Assert(value, Equals, uint16(15))
		}
	}
	for i := uint16(0); i < 18; i++ {
		value := nibCount.Dec(filterWord, filterBit)
		if i < 15 {
			c.Assert(value, Equals, 14-i)
		} else {
			c.Assert(value, Equals, uint16(0))
		}
	}
}
func (s *XLSuite) doTestWord(c *C, nibCount *NibbleCounters, filterWord uint) {
	// test the low order bit, the high order bit, and a few from
	// the middle
	s.doTestBit(c, nibCount, filterWord, 0)
	s.doTestBit(c, nibCount, filterWord, 1)
	s.doTestBit(c, nibCount, filterWord, 14)
	s.doTestBit(c, nibCount, filterWord, 15)
	s.doTestBit(c, nibCount, filterWord, 16)
	s.doTestBit(c, nibCount, filterWord, 31)

}
func (s *XLSuite) doTest(c *C, m uint) {

	filterInts := uint(1<<m) / 32 // 2^m bits fit into this many ints
	// test the low order word, the high order word, and a few from
	// the middle
	nibCount := NewNibbleCounters(filterInts)
	s.doTestWord(c, nibCount, 0)
	s.doTestWord(c, nibCount, (filterInts/2)-1)
	s.doTestWord(c, nibCount, (filterInts / 2))
	s.doTestWord(c, nibCount, (filterInts/2)+1)
	s.doTestWord(c, nibCount, filterInts-1)
}
func (s *XLSuite) TestNibs(c *C) {
	s.doTest(c, 12)
	s.doTest(c, 13)
	s.doTest(c, 14)
	s.doTest(c, 20)
}
