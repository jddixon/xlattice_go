package overlay

// xlattice_go/overlay/parser_test.go

import (
	"fmt"
	"github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
	"regexp"
	"strings"
)

var _ = fmt.Print

func (s *XLSuite) getAName(rng *rnglib.PRNG) (name string) {
	name = string(rng.NextFileName(8))
	for {
		first := string(name[0])
		if !strings.ContainsAny(name, "-_.") && !strings.ContainsAny(first, "0123456789") {
			break
		}
		name = string(rng.NextFileName(8))
	}
	return
}
func (s *XLSuite) TestParser(c *C) {
	// fmt.Println("TEST_PARSER")
	rng := rnglib.MakeSimpleRNG()
	for i := 0; i < 16; i++ {
		s.doTestParser(c, rng)
	}
}
func (s *XLSuite) doTestParser(c *C, rng *rnglib.PRNG) {

	name := s.getAName(rng)
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
	// XXX ADDR RANGE MISSING
	c.Assert(transport, Equals, o.Transport())
	c.Assert(float32(cost), Equals, o.Cost())

	text := o.String()
	// DEBUG
	// fmt.Printf("serialized overlay is %s\n", text)
	// END
	o2, err := Parse(text)
	c.Assert(err, IsNil)
	c.Assert(text, Equals, o2.String())
}
func (s *XLSuite) TestName(c *C) {
	nameRE, err := regexp.Compile(NAME)
	c.Assert(err, IsNil)
	c.Assert(nameRE.MatchString("foo"), Equals, true)
}
func (s *XLSuite) TestNameGroup(c *C) {
	// fmt.Println("TEST_NAME_GROUP")
	NAME_GROUP := `^\s*(` + NAME + `)\s*$`
	// fmt.Printf("NAME_GROUP: '%s`\n", NAME_GROUP)
	nameGroupRE, err := regexp.Compile(NAME_GROUP)
	c.Assert(err, IsNil)
	fooStr := " foo  "
	c.Assert(nameGroupRE.MatchString(fooStr), Equals, true)
	groups := nameGroupRE.FindStringSubmatch(fooStr)
	c.Assert(groups, Not(IsNil))
	c.Assert(len(groups), Equals, 2)
	c.Assert(groups[1], Equals, "foo")
}
func (s *XLSuite) TestTwoNames(c *C) {
	// fmt.Println("TEST_TWO_NAMES")
	var groups []string
	partialPat := `^overlay:\s*(` + NAME + `),\s*(` + NAME + `),\s*$`
	partialRE, err := regexp.Compile(partialPat)
	c.Assert(err, IsNil)
	const STR = "overlay: foo, bar,"
	groups = partialRE.FindStringSubmatch(STR)
	c.Assert(groups, Not(IsNil))
	c.Assert(groups[0], Equals, STR)
	c.Assert(groups[1], Equals, "foo")
	c.Assert(groups[2], Equals, "bar")
}
func (s *XLSuite) TestAddrRange(c *C) {
	// fmt.Println("TEST_ADDR_RANGE")
	addrRangeRE, err := regexp.Compile(`(` + ADDR_RANGE + `),\s*`)
	c.Assert(err, IsNil)
	cidrBlock := "192.168.0.1/18,  "
	c.Assert(addrRangeRE.MatchString(cidrBlock), Equals, true)
	groups := addrRangeRE.FindStringSubmatch(cidrBlock)
	c.Assert(groups, Not(IsNil))
	c.Assert(len(groups), Equals, 2)
	c.Assert(groups[1], Equals, "192.168.0.1/18")
}
func (s *XLSuite) TestAddrRangePlusCost(c *C) {
	// fmt.Println("TEST_ADDR_RANGE_PLUS_COST")
	addrRangeRE, err := regexp.Compile(`(` + ADDR_RANGE + `),\s*(\d+\.\d*)`)
	c.Assert(err, IsNil)
	cidrBlock := "192.168.0.1/18, 123.456"
	// c.Assert(addrRangeRE.MatchString(cidrBlock), Equals, true)
	groups := addrRangeRE.FindStringSubmatch(cidrBlock)
	c.Assert(groups, Not(IsNil))
	c.Assert(len(groups), Equals, 3)
	// groups[1] is actually 18
	// c.Assert(groups[1], Equals, "192.168.0.1/18")
	c.Assert(groups[2], Equals, "123.456")
} // FOO

func (s *XLSuite) TestCost(c *C) {
	costRE, err := regexp.Compile(`^\d+\.\d*$`)
	c.Assert(err, IsNil)
	c.Assert(costRE.MatchString("237.0000"), Equals, true)
}
func (s *XLSuite) TestPartialRE(c *C) {
	// fmt.Println("TEST_PARTIAL_RE")
	var groups []string

	// IP_OVERLAY = `overlay:\s*(` + NAME + `),\s*(` + NAME + `),\s*(` + ADDR_RANGE + `),\s*(\d\.\d*)`
	// XXX WORKING HERE
	partialPat := `^overlay:\s*(` + NAME + `),\s*(` + NAME + `),\s*(` + ADDR_RANGE + `),\s*$`
	partialRE, err := regexp.Compile(partialPat)
	c.Assert(err, IsNil)

	const STR = "overlay: foo, bar, 192.168.1.10/24,"
	groups = partialRE.FindStringSubmatch(STR)
	c.Assert(groups, Not(IsNil)) // FAILS
}
