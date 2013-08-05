package overlay

// xlattice_go/overlay/parser_test.go

import (
	"fmt"
	x "github.com/jddixon/xlattice_go"
	. "launchpad.net/gocheck"
	"regexp"
)

var _ = fmt.Print

func (s *XLSuite) TestParser(c *C) {
	rng := x.MakeSimpleRNG()
	name := rng.NextFileName(8)
	a := rng.Intn(256)
	b := rng.Intn(256)
	_c := rng.Intn(256)
	d := rng.Intn(256)
	bits := rng.Intn(33)
	aRange := fmt.Sprintf("%d.%d.%d.%d/%d", a, b, _c, d, bits)
	transport := "tcp"
	cost := float32(rng.Intn(300)) / 100.0

	ar, err := NewCIDRAddrRange(aRange)
	c.Assert(err, IsNil)

	o, err := NewIPOverlay(name, ar, transport, cost)
	c.Assert(err, IsNil)
	c.Assert(o, Not(IsNil))

	c.Assert(name, Equals, o.Name())
	// ADDR RANGE
	c.Assert(transport, Equals, o.Transport())
	c.Assert(float32(cost), Equals, o.Cost())

	text := o.String()
	o2, err := Parse(text)
	c.Assert(err, IsNil)
	c.Assert(err, Not(IsNil))
	c.Assert(text, Equals, o2.String())
}
func (s *XLSuite) TestRE(c *C) {
	nameRE, err := regexp.Compile(NAME)
	c.Assert(err, IsNil)
	c.Assert(nameRE.MatchString("foo"), Equals, true)

	addrRangeRE, err := regexp.Compile(ADDR_RANGE)
	c.Assert(err, IsNil)
	c.Assert(addrRangeRE.MatchString("127.0.0.1/8"), Equals, true)

	// if you add ^ and $ this test fails
	costRE, err := regexp.Compile(`\d\.\d*`)
	c.Assert(err, IsNil)
	c.Assert(costRE.MatchString("237.0000"), Equals, true)

	// XXX problems from here

	// IP_OVERLAY = `overlay:\s*(` + NAME + `),\s*(` + NAME + `),\s*(` + ADDR_RANGE + `),\s*(\d\.\d*)`
	// partialPat := `overlay:\s*(` + NAME + `),\s*(` + NAME + `),\s*(` + ADDR_RANGE + `),\s*(\d\.\d*)`
	partialPat := `overlay: (` + NAME + `), (` + NAME + `),`
	partialRE, err := regexp.Compile(partialPat)
	c.Assert(err, IsNil)
	const STR = "overlay: foo, bar,"
	groups := partialRE.FindStringSubmatch(STR)
	c.Assert(groups, IsNil)      // succeeds
	c.Assert(groups, Not(IsNil)) // FAILS
	// working here
	_ = groups
	c.Assert(partialRE.MatchString(STR), Equals, true)

}
