package reg

// xlattice_go/msg/reg_cluster_test.go

import (
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestClusterAttrs(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_ATTRS")
	}
	c.Assert(CLUSTER_DELETED, Equals, 1)
	c.Assert(FOO, Equals, 2)
	c.Assert(BAR, Equals, 4)
}

func (s *XLSuite) TestClusterMaker(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_MAKER")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster
	maxSize := 2 + rng.Intn(6) // so from 2 to 7
	cl := s.makeACluster(c, rng, maxSize)

	_ = cl // DEBUG

	c.Assert(cl.MaxSize, Equals, maxSize)
	c.Assert(cl.Size(), Equals, maxSize)

	// and that their names are unique
	// XXX STUB ///

	// and that the byName index is correct
	// XXX STUB ///

	// and that the byID index is correct
	// XXX STUB ///

}
func (s *XLSuite) TestClusterSerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_SERIALIZATION")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster
	size := 2 + rng.Intn(6) // so from 2 to 7
	cl := s.makeACluster(c, rng, size)

	// Serialize it
	serialized := cl.String()

	// Reverse the serialization
	deserialized, rest, err := ParseRegCluster(serialized)
	c.Assert(err, IsNil)
	c.Assert(deserialized, Not(IsNil))
	c.Assert(len(rest), Equals, 0)

	// Verify that the deserialized cluster is identical to the original
	c.Assert(deserialized.Equal(cl), Equals, true)

}
