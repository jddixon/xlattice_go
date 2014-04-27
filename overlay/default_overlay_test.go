package overlay

// xlattice_go/overlay/default_overlay_test.go

import (
	"github.com/jddixon/xlattice_go/rnglib"
	xt "github.com/jddixon/xlattice_go/transport"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) shouldGetDefault(c *C, addr string) OverlayI {
	e, err := xt.NewTcpEndPoint(addr)
	c.Assert(err, Equals, nil)
	c.Assert(e, Not(Equals), nil)

	o, err := DefaultOverlay(e)
	c.Assert(err, Equals, nil)
	c.Assert(o.Name(), Not(Equals), "") // NPE?
	c.Assert(o.Transport(), Equals, "ip")
	c.Assert(o.Cost(), Equals, float32(1.0))
	return o
}
func (s *XLSuite) TestDefaultOverlay(c *C) {
	rng := rnglib.MakeSimpleRNG()
	_ = rng

	o := s.shouldGetDefault(c, "127.0.0.1:27")
	c.Assert(o.Name(), Equals, "localhost")
	aRange := o.(*IPOverlay).AddrRange()
	expectedAR, err := NewCIDRAddrRange("127.0.0.0/8")
	c.Assert(err, Equals, nil)
	c.Assert(expectedAR.Equal(aRange), Equals, true)

	o = s.shouldGetDefault(c, "10.0.29.1:52")
	c.Assert(o.Name(), Equals, "privateA")
	aRange = o.(*IPOverlay).AddrRange()
	expectedAR, err = NewCIDRAddrRange("10.0.0.0/8")
	c.Assert(err, Equals, nil)
	c.Assert(expectedAR.Equal(aRange), Equals, true)

	o = s.shouldGetDefault(c, "172.17.9.1:5")
	c.Assert(o.Name(), Equals, "privateB")
	aRange = o.(*IPOverlay).AddrRange()
	expectedAR, err = NewCIDRAddrRange("172.16.0.0/12")
	c.Assert(err, Equals, nil)
	c.Assert(expectedAR.Equal(aRange), Equals, true)

	o = s.shouldGetDefault(c, "192.168.136.254:95")
	c.Assert(o.Name(), Equals, "privateC")
	aRange = o.(*IPOverlay).AddrRange()
	expectedAR, err = NewCIDRAddrRange("192.168.0.0/16")
	c.Assert(err, Equals, nil)
	c.Assert(expectedAR.Equal(aRange), Equals, true)

	o = s.shouldGetDefault(c, "92.168.136.254:95")
	c.Assert(o.Name(), Equals, "globalV4")
	aRange = o.(*IPOverlay).AddrRange()
	expectedAR, err = NewCIDRAddrRange("0.0.0.0/0")
	c.Assert(err, Equals, nil)
	c.Assert(expectedAR.Equal(aRange), Equals, true)

}
