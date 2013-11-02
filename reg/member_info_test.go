package reg

// xlattice_go/msg/member_info_test.go

import (
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestMISerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MI_SERIALIZATION")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster member
	cm := s.makeAMemberInfo(c, rng)

	// Serialize it
	serialized := cm.String()

	// Reverse the serialization
	deserialized, rest, err := ParseMemberInfo(serialized)
	c.Assert(err, IsNil)
	c.Assert(len(rest), Equals, 0)

	// Verify that the deserialized member is identical to the original
	c.Assert(deserialized.Equal(cm), Equals, true)
}

func (s *XLSuite) TestMemberInfoAndTokens(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_MEMBER_INFO_AND_TOKENS")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster member
	cm := s.makeAMemberInfo(c, rng)

	token, err := cm.Token()
	c.Assert(err, IsNil)

	cm2, err := NewMemberInfoFromToken(token)
	c.Assert(err, IsNil)
	c.Assert(cm.Equal(cm2), Equals, true)
}
