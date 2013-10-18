package util

// xlattice_go/util/decimal_version_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) makeVersion(a, b, c, d uint) (dv DecimalVersion) {
	return DecimalVersion(uint((0xff & a) |
		((0xff & b) << 8) |
		((0xff & c) << 16) |
		((0xff & d) << 24)))
}
func (s *XLSuite) TestDecimalVersion(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DECIMAL_VERSION")
	}
	rng := xr.MakeSimpleRNG()

	_ = rng

	// always print at least two decimals
	dv := s.makeVersion(1, 0, 0, 0)
	v := dv.String()
	c.Assert(v, Equals, "1.0")
	dv2, err := ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	// don't print more if the values are zero
	dv = s.makeVersion(1, 2, 0, 0)
	v = dv.String()
	c.Assert(v, Equals, "1.2")
	dv2, err = ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	// if the third byte is zero but the fourth isn't, print
	// both
	dv = s.makeVersion(1, 2, 0, 4)
	v = dv.String()
	c.Assert(v, Equals, "1.2.0.4")
	dv2, err = ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	// other cases
	dv = s.makeVersion(1, 2, 3, 0)
	v = dv.String()
	c.Assert(v, Equals, "1.2.3")
	dv2, err = ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	dv = s.makeVersion(1, 2, 3, 4)
	v = dv.String()
	c.Assert(v, Equals, "1.2.3.4")
	dv2, err = ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	for i := 0; i < 8; i++ {
		n := rng.Uint32()
		dv := DecimalVersion(n)
		v = dv.String()
		dv2, err := ParseDecimalVersion(v)
		c.Assert(err, IsNil)
		c.Assert(dv2, Equals, dv)
	}
}
