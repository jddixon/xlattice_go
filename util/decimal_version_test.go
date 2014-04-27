package util

// xlattice_go/util/decimal_version_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestVersionFromBytes(c *C) {
	dv1 := New(1, 2, 3, 4)
	var b []byte = []byte{1, 2, 3, 4}
	dv2, err := VersionFromBytes(b)
	c.Assert(err, IsNil)
	c.Assert(dv1, Equals, dv2)

	dv3, err := VersionFromBytes(b[1:]) // so only 3 bytes long
	c.Assert(err, Equals, WrongLengthForVersion)
	c.Assert(dv3, Equals, DecimalVersion(0))
}
func (s *XLSuite) TestDecimalVersion(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DECIMAL_VERSION")
	}
	rng := xr.MakeSimpleRNG()

	_ = rng

	// always print at least two decimals
	dv := New(1, 0, 0, 0)
	v := dv.String()
	c.Assert(v, Equals, "1.0")
	dv2, err := ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	// don't print more if the values are zero
	dv = New(1, 2, 0, 0)
	v = dv.String()
	c.Assert(v, Equals, "1.2")
	dv2, err = ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	// if the third byte is zero but the fourth isn't, print
	// both
	dv = New(1, 2, 0, 4)
	v = dv.String()
	c.Assert(v, Equals, "1.2.0.4")
	dv2, err = ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	// other cases
	dv = New(1, 2, 3, 0)
	v = dv.String()
	c.Assert(v, Equals, "1.2.3")
	dv2, err = ParseDecimalVersion(v)
	c.Assert(err, IsNil)
	c.Assert(dv2, Equals, dv)

	dv = New(1, 2, 3, 4)
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
