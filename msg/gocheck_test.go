package msg

import (
	"crypto/rand"
	"crypto/rsa"
	xr "github.com/jddixon/rnglib_go"
	xi "github.com/jddixon/xlNodeID_go"
	xu "github.com/jddixon/xlUtil_go"
	. "gopkg.in/check.v1"
	"testing"
)

// IF USING test framework, need a file like this in each package=directory.

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

const (
	VERBOSITY = 1
)

/////////////////////////////////////////////////////////////////
// FROM ../reg; BEING HACKED INTO A TEST OF THIS PACKAGE'S CRYPTO
/////////////////////////////////////////////////////////////////

func (s *XLSuite) makeAnID(c *C, rng *xr.PRNG) (id []byte) {
	id = make([]byte, xu.SHA2_BIN_LEN)
	rng.NextBytes(id)
	return
}
func (s *XLSuite) makeANodeID(c *C, rng *xr.PRNG) (nodeID *xi.NodeID) {
	id := s.makeAnID(c, rng)
	nodeID, err := xi.New(id)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))
	return
}
func (s *XLSuite) makeAnRSAKey(c *C) (key *rsa.PrivateKey) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	c.Assert(err, IsNil)
	c.Assert(key, Not(IsNil))
	return key
}
