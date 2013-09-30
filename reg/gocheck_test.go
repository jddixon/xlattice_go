package reg

import (
	"crypto/rand"
	"crypto/rsa"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
	"testing"
)

// IF USING gocheck, need a file like this in each package=directory.

func Test(t *testing.T) { TestingT(t) }

type XLSuite struct{}

var _ = Suite(&XLSuite{})

const (
	VERBOSITY = 1
)

func (s *XLSuite) makeAnID(c *C, rng *xr.PRNG) (id []byte) {
	id = make([]byte, SHA3_LEN)
	rng.NextBytes(&id)
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

// Using functions must check to ensure members have unique names

func (s *XLSuite) makeAClusterMember(c *C, rng *xr.PRNG) *ClusterMember {
	attrs := uint64(rng.Int63())
	bn, err := xn.NewBaseNode(
		rng.NextFileName(8),
		s.makeANodeID(c, rng),
		&s.makeAnRSAKey(c).PublicKey,
		&s.makeAnRSAKey(c).PublicKey,
		nil) // overlays
	c.Assert(err, IsNil)
	return &ClusterMember{
		attrs:    attrs,
		BaseNode: *bn,
	}
}

// Make a RegCluster for test purposes.  Cluster member names are guaranteed
// to be unique but the name of the cluster itself may not be.

func (s *XLSuite) makeACluster(c *C, rng *xr.PRNG, size int) (rc *RegCluster) {

	var err error
	c.Assert(1 < size && size <= 64, Equals, true)

	attrs := uint64(rng.Int63())
	name := rng.NextFileName(8) // no guarantee of uniqueness
	id := s.makeANodeID(c, rng)

	rc, err = NewRegCluster(attrs, name, id, size)
	c.Assert(err, IsNil)

	for count := 0; count < size; count++ {
		cm := s.makeAClusterMember(c, rng)
		for {
			if _, ok := rc.MembersByName[cm.GetName()]; ok {
				// name is in use, so try again
				cm = s.makeAClusterMember(c, rng)
			} else {
				err = rc.AddMember(cm)
				c.Assert(err, IsNil)
				break
			}
		}
	}
	return
}
