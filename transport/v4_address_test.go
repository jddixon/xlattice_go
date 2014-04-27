package transport

// xlattice_go/transport/v4_address_test.go

import (
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	"regexp"
)

func (s *XLSuite) TestV4AddressInterface(c *C) {
	// fmt.Println("TEST_V4_ADDRESS_INTERFACE")
	const myWeb = "127.0.0.1:80"
	w, err := NewV4Address(myWeb)
	c.Assert(err, Equals, nil)
	c.Assert(w, Not(Equals), nil)

	x, err := w.Clone()
	c.Assert(err, Equals, nil)
	c.Assert(x, Not(Equals), nil)
	c.Assert(w.String(), Equals, x.String())
	c.Assert(w.Equal(x), Equals, true)

	foo := AddressI(w)
	_ = foo
}
func (s *XLSuite) TestGoodV4Addrs(c *C) {
	// fmt.Println("TEST_GOOD_V4_ADDRS")
	rng := rnglib.MakeSimpleRNG()
	for i := 0; i < 16; i++ {
		_a := rng.Intn(256)
		_b := rng.Intn(256)
		_c := rng.Intn(256)
		_d := rng.Intn(256)
		_p := rng.Intn(256 * 256)
		var s string
		// half of test cases have ports
		if rng.NextBoolean() {
			s = fmt.Sprintf("%d.%d.%d.%d", _a, _b, _c, _d)
		} else {
			s = fmt.Sprintf("%d.%d.%d.%d:%d", _a, _b, _c, _d, _p)
		}
		a, err := NewV4Address(s)
		c.Assert(err, Equals, nil)
		c.Assert(a, Not(Equals), nil)
		c.Assert(a.String(), Equals, s)
	}
}
func (s *XLSuite) TestQuad(c *C) {
	// fmt.Println("TEST_QUAD")
	MY_PAT := `^(` + QUAD_PAT + `)$`
	myRE, err := regexp.Compile(MY_PAT)
	c.Assert(err, Equals, nil)
	c.Assert(myRE, Not(Equals), nil)

	for i := 0; i < 256; i++ {
		val := fmt.Sprintf("%d", i)
		c.Assert(myRE.MatchString(val), Equals, true)
	}

	c.Assert(myRE.MatchString(""), Equals, false)
	c.Assert(myRE.MatchString("abc"), Equals, false)
	c.Assert(myRE.MatchString("301"), Equals, false)
	c.Assert(myRE.MatchString("256"), Equals, false)
	c.Assert(myRE.MatchString("1a"), Equals, false)
}
func (s *XLSuite) TestDottedQuad(c *C) {
	// fmt.Println("TEST_DOTTED_QUAD")
	rng := rnglib.MakeSimpleRNG()
	// Use of MustCompile makes no difference.
	v4AddrRE, err := regexp.Compile(V4_ADDR_PAT)
	c.Assert(err, Equals, nil)
	c.Assert(v4AddrRE, Not(Equals), nil)

	for i := 0; i < 16; i++ {
		_a := rng.Intn(256)
		_b := rng.Intn(256)
		_c := rng.Intn(256)
		_d := rng.Intn(256)
		val := fmt.Sprintf("%d.%d.%d.%d", _a, _b, _c, _d)
		c.Assert(v4AddrRE.MatchString(val), Equals, true)
	}
	c.Assert(v4AddrRE.MatchString(`abc`), Equals, false)
	// XXX the next four inexplicably return true
	c.Assert(v4AddrRE.MatchString(`1a.2b.3c.4d`), Equals, false)
	c.Assert(v4AddrRE.MatchString(`1.2.3`), Equals, false)
	c.Assert(v4AddrRE.MatchString(`301.2.3`), Equals, false)
	c.Assert(v4AddrRE.MatchString(`1.2.3.4.5`), Equals, false)
}
