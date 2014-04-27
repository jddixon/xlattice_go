package reg

// xlattice_go/msg/reg_cluster_test.go

import (
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) TestClusterMaker(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_MAKER")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster
	epCount := uint(1 + rng.Intn(3)) // so from 1 to 3
	maxSize := uint(2 + rng.Intn(6)) // so from 2 to 7
	cl := s.makeACluster(c, rng, epCount, maxSize)

	_ = cl // DEBUG

	c.Assert(cl.MaxSize(), Equals, maxSize)
	c.Assert(cl.Size(), Equals, maxSize) //

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
	epCount := uint(1 + rng.Intn(3)) // so from 1 to 3
	size := uint(2 + rng.Intn(6))    // so from 2 to 7
	cl := s.makeACluster(c, rng, epCount, size)

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
