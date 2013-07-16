package overlay

// xlattice_go/overlay/ip_overlay_test.go

import (
	x "github.com/jddixon/xlattice_go"
	. "launchpad.net/gocheck"
	"testing"
)

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
