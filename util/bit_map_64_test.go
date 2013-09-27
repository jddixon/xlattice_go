package util

// xlattice_go/util/bit_map_64_test.go

import (
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestBitMap64(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BIT_MAP_64")
	}
	rng := rnglib.MakeSimpleRNG()

	_ = rng

	c.Assert((*LowNMap(0)).Bits, Equals, uint64(0))

	b := uint64(1)
	for n := uint(1); n <= 64; n++ {
		c.Assert((*LowNMap(n)).Bits, Equals, b)
		b = (b << 1) | 1
	}
}
func (s *XLSuite) TestCounts(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_COUNTS")
	}
	rng := rnglib.MakeSimpleRNG()
	
	for i := 0; i < 8 ; i++ {
		x := uint64(rng.Int63())
		
		// these operate on uint64s, not BitMap64s
		count3 := popCount3(x)
		c.Assert( count3, Equals, popCount4(x))

		// test the BitMap64 operation
		b := NewBitMap64(x)
		c.Assert(b.Count(), Equals, count3)
	}
}

func (s *XLSuite) TestOtherFuncs(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_OTHER_FUNCS")
	}
	rng := rnglib.MakeSimpleRNG()

	for i := 0; i < 8 ; i++ {
		x := uint64(rng.Int63())
		y := uint64(rng.Int63())
		for x == y {
			y = uint64(rng.Int63())
		}
		xx := NewBitMap64(x)
		yy := NewBitMap64(y)
		c.Assert(xx.Equal(xx), Equals, true)
		c.Assert(xx.Equal(x), Equals, false)
		c.Assert(xx.Equal(yy), Equals, false)
		zz := xx.Clone()
		c.Assert(xx.Equal(zz), Equals, true)

		zero := NewBitMap64(0)
		one  := NewBitMap64(0xffffffffffffffff)

		c.Assert(one.Any(), Equals, true)
		c.Assert(one.None(), Equals, false)
		c.Assert(zero.Any(), Equals, false)
		c.Assert(zero.None(), Equals, true)
		
		c.Assert(zero.Union(xx).Equal(xx), Equals,true)	
		c.Assert(one.Intersection(xx).Equal(xx), Equals, true)

		c.Assert(xx.Union(xx.Complement()).Equal(one), Equals, true)
		c.Assert(xx.Intersection(xx.Complement()).Equal(zero), Equals, true)

		count := xx.Count()
		n := uint(rng.Intn(64))	// so values range from 0 to 63
		if xx.Test(n) {
			// it's set
			rr := xx.Flip(n)
			c.Assert( rr.Count(), Equals, count - 1 )
			c.Assert( xx.Difference(rr).Count(), Equals, 1)


		} else {
			// the bit is not set
			rr := xx.Flip(n)
			c.Assert( rr.Count(), Equals, count + 1 )
			c.Assert( rr.Difference(xx).Count(), Equals, 1)
		}
	}
}
