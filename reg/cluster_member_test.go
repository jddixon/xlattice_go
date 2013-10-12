package reg

// xlattice_go/msg/cluster_member_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestCMSerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CM_SERIALIZATION")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster member
	cm := s.makeAClusterMember(c, rng)

	// Serialize it
	serialized := cm.String()

	// Reverse the serialization
	deserialized, rest, err := ParseClusterMember(serialized)
	c.Assert(err, IsNil)
	c.Assert(len(rest), Equals, 0)

	// Verify that the deserialized member is identical to the original
	c.Assert(deserialized.Equal(cm), Equals, true)
}

func (s *XLSuite) TestMembersAndTokens(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MEMBERS_AND_TOKENS")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster member
	cm := s.makeAClusterMember(c, rng)

	token, err := cm.Token()
	c.Assert(err, IsNil)

	cm2, err := NewClusterMemberFromToken(token)
	c.Assert(err, IsNil)
	c.Assert(cm.Equal(cm2), Equals, true)
}
