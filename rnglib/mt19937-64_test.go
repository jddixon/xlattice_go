package rnglib

// This is the original Hiroshima test hacked into Go -- and so
// doesn't actually do any useful unit tests.

import (
	"fmt"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestMT64(c *C) {
	var (
		// no size, so it's a slice
		initX  = []uint64{0x12345, 0x23456, 0x34567, 0x45678}
		length = uint64(4)
	)
	m := NewMT64()
	m.init_by_array64(initX, length)
	fmt.Println("1000 outputs of genrand64_int64()")
	for i := 0; i < 1000; i++ {
		// fmt.Printf("%16v ", m.genrand64_int64())
		fmt.Printf("%20v ", m.genrand64_int64())
		if i%5 == 4 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n1000 outputs of genrand64_real2()\n")
	for i := 0; i < 1000; i++ {
		fmt.Printf("%10.8f ", m.genrand64_real2())
		if i%5 == 4 {
			fmt.Printf("\n")
		}
	}
}
