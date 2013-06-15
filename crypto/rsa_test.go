package crypto

// xlattice_go/crypto/rsa_test.go

import (
	. "launchpad.net/gocheck"
	"math/big"
	"testing"
)

// gocheck tie-in /////////////////////
func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

// end gocheck setup //////////////////

// Fiddling around to see whether gocheck could compare bigInts (answer: no).
func (s *XLSuite) TestUnity(c *C) {
	c.Assert(big.NewInt(1).Int64(), Equals, (*BIG_ONE).Int64())
}
