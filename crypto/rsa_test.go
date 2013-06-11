package crypto

// xlattice_go/crypto/rsa_test.go

import (
	. "launchpad.net/gocheck"
	"testing"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

func (s *XLSuite) TestUnity(c *C) {
	c.Assert(1, Equals, BIG_ONE) // EXPECTED to fail
}
