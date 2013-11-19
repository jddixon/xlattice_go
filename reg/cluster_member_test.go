package reg

// xlattice_go/reg/cluster_member_test.go

import (
	"fmt"
	//xc "github.com/jddixon/xlattice_go/crypto"
	xi "github.com/jddixon/xlattice_go/nodeID"
	xr "github.com/jddixon/xlattice_go/rnglib"
	//xt "github.com/jddixon/xlattice_go/transport"
	. "launchpad.net/gocheck"
)

func (s *XLSuite) TestClusterMemberSerialization(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_CLUSTER_MEMBER_SERIALIZATION")
	}
	rng := xr.MakeSimpleRNG()

	// Generate a random cluster
	epCount := uint(1 + rng.Intn(3)) // so from 1 to 3
	size := uint(2 + rng.Intn(6))    // so from 2 to 7
	cl := s.makeACluster(c, rng, epCount, size)

	// We are going to overwrite cluster member zero's attributes
	// with those of the new cluster member.
	myNode, myCkPriv, mySkPriv := s.makeHostAndKeys(c, rng)
	myAttrs := cl.Members[0].Attrs

	var myEnds []string
	for i := uint(0); i < epCount; i++ {
		myEnds = append(myEnds, "127.0.0.1:0")
	}
	myMemberInfo, err := NewMemberInfo(myNode.GetName(), myNode.GetNodeID(),
		&myCkPriv.PublicKey, &mySkPriv.PublicKey, myAttrs, myEnds)
	c.Assert(err, IsNil)
	// overwrite member 0
	cl.Members[0] = myMemberInfo

	myClusterID, err := xi.New(cl.ID)
	c.Assert(err, IsNil)

	cm := &ClusterMember{
		Attrs:        myAttrs,
		ClusterName:  cl.Name,
		ClusterAttrs: cl.Attrs,
		ClusterID:    myClusterID,
		ClusterSize:  uint32(cl.MaxSize()),
		SelfIndex:    uint32(0),
		Members:      cl.Members, // []*MemberInfo
		EpCount:      uint32(epCount),
		Node:         *myNode,
	}

	// simplest test of Equal()
	c.Assert(cm.Equal(cm), Equals, true)

	// Serialize it
	serialized := cm.String()

	// Reverse the serialization
	deserialized, rest, err := ParseClusterMember(serialized)
	c.Assert(err, IsNil)
	c.Assert(deserialized, Not(IsNil))
	c.Assert(len(rest), Equals, 0)

	// Verify that the deserialized cluster is identical to the original
	// First version:
	c.Assert(deserialized.Equal(cm), Equals, true)

	// Second version of identity test:
	serialized2 := deserialized.String()
	c.Assert(serialized2, Equals, serialized)
}
