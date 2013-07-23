package overlay

// xlattice_go/overlay/ip_overlay_test.go

import (
	"fmt"
	x "github.com/jddixon/xlattice_go"
	xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
	"net"
	"testing"
)
var _ = fmt.Print

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

func (s *XLSuite) TestCtor(c *C) {
	rng := x.MakeSimpleRNG()
	name := rng.NextFileName(8)

	o, err := NewIPOverlay(name, nil, "tcpip", 0.42)
	c.Assert(err, IsNil)
	c.Assert(o, Not(IsNil))
	c.Assert(name, Equals, o.Name())
	c.Assert("tcpip", Equals, o.Transport())
	c.Assert(float32(0.42), Equals, o.Cost())
}

func (s *XLSuite) TestIsElement(c *C) {
	rng := x.MakeSimpleRNG()
	name := rng.NextFileName(8)
	p10_8 := net.ParseIP("10.0.0.0")[12:]
	a10_8, err := NewAddrRange(p10_8, 8, 32)
	c.Assert(err, IsNil)
	o10_8, err := NewIPOverlay(name, a10_8, "ip", 1.0)
	c.Assert(err, IsNil)

	_ = o10_8				// DEBUG

	// bad transport(s)
	mockE := xt.NewMockEndPoint( "foo", "1234" )
	c.Assert( o10_8.IsElement(mockE), Equals, false)

	// 10/8 ---------------------------------------------------------
	n10_8, n10_8_IPNet, err := net.ParseCIDR("10.0.0.0/8")
	c.Assert( err, IsNil )	
	c.Assert( n10_8, Not(IsNil) )			
	c.Assert( len(n10_8), Equals, 16)			// XXX is 16
	c.Assert( len(n10_8_IPNet.IP), Equals, 4)
	e1 := net.ParseIP("10.11.12.13")[12:]
	c.Assert(e1, Not(IsNil))
	c.Assert(n10_8_IPNet.Contains(e1), Equals, true)
	e2 := net.ParseIP("9.10.11.12")[12:]
	c.Assert(n10_8_IPNet.Contains(e2), Equals, false)

	// 192.168/16 ---------------------------------------------------
	n192_168, n192_168_IPNet, err := net.ParseCIDR("192.168.0.0/16")
	c.Assert( err, IsNil )	
	c.Assert( n192_168, Not(IsNil) )			
	c.Assert( len(n192_168), Equals, 16)			// XXX is 16
	c.Assert( len(n192_168_IPNet.IP), Equals, 4)
	// The first value returned by ParseCIDR is all zeroes.  The
	// IP in the second value, the IPNet, is correct
	e10 := net.ParseIP("192.168.0.0")[12:]
	c.Assert(n192_168_IPNet.Contains(e10), Equals, true)
	e11 := net.ParseIP("192.168.255.255")[12:]
	c.Assert(n192_168_IPNet.Contains(e11), Equals, true)
	e20 := net.ParseIP("192.167.255.255")[12:]
	c.Assert(n192_168_IPNet.Contains(e20), Equals, false)
	e21 := net.ParseIP("192.169.0.0")[12:]
	c.Assert(n192_168_IPNet.Contains(e21), Equals, false)

	
}
