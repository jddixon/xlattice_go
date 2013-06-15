package overlay

// xlattice_go/overlay/overlay_test.go

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

	o, err := NewOverlay(name, nil, "tcpip", 0.42)
	c.Assert(err, IsNil)
	c.Assert(o, Not(IsNil))
}
