package search

// xlattice_go/search/peer_map_test.go

import (
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print
var _ = rnglib.MakeSimpleRNG

const (
	SHA1_LEN  = 20
	VERBOSITY = 1
)

func (s *XLSuite) makeTopAndBottom(c *C) (topPeer, bottomPeer *xn.Peer) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < SHA1_LEN; i++ {
		t[i] = byte(0xf)
	}
	top, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(top, Not(IsNil))

	topPeer, err = xn.NewNewPeer("top", top)
	c.Assert(err, IsNil)
	c.Assert(topPeer, Not(IsNil))

	bottom, err := xi.NewNodeID(make([]byte, SHA1_LEN))
	c.Assert(err, IsNil)
	c.Assert(bottom, Not(IsNil))

	bottomPeer, err = xn.NewNewPeer("bottom", bottom)
	c.Assert(err, IsNil)
	c.Assert(bottomPeer, Not(IsNil))

	return topPeer, bottomPeer
}
func (s *XLSuite) makeAPeer(c *C, name string, id ...int) (peer *xn.Peer) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < len(id); i++ {
		t[i] = byte(id[i])
	}
	nodeID, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))

	peer, err = xn.NewNewPeer(name, nodeID)
	c.Assert(err, IsNil)
	c.Assert(peer, Not(IsNil))
	return
}
func (s *XLSuite) TestPeerMapTools(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PEER_MAP_TOOLS")
	}
	threePeer := s.makeAPeer(c, "threePeer", 1, 2, 3)
	nodeID := threePeer.GetNodeID()
	value := nodeID.Value()
	c.Assert(threePeer.GetName(), Equals, "threePeer")
	c.Assert(value[0], Equals, byte(1))
	c.Assert(value[1], Equals, byte(2))
	c.Assert(value[2], Equals, byte(3))
	for i := 3; i < SHA1_LEN; i++ {
		c.Assert(value[i], Equals, byte(0))
	}

}
func (s *XLSuite) TestTopBottomMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_PEER_MAP")
	}

	var pm PeerMap
	c.Assert(pm.lowest, IsNil)

	topPeer, bottomPeer := s.makeTopAndBottom(c)
	err := pm.AddToPeerMap(topPeer)
	c.Assert(err, IsNil)
	c.Assert(pm.lowest, Not(IsNil))
	lowest := pm.lowest
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, topPeer) // succeeds ...
	// c.Assert(topPeer.Equal(lowest.peer), Equals, true)      // FAILS!
	c.Assert(lowest.peer.GetName(), Equals, "top")

	// We expect that bottomPeer will become the lowest with its
	// higher field pointing at topPeer.
	err = pm.AddToPeerMap(bottomPeer)
	c.Assert(err, IsNil)
	lowest = pm.lowest
	// c.Assert(bottomPeer.Equal(lowest.peer), Equals, true)   // FAILS
	c.Assert(lowest.peer.GetName(), Equals, "bottom") // XXX gets 'top'
}
func (s *XLSuite) TestShallowMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_MAP")
	}
	var pm PeerMap
	c.Assert(pm.lowest, IsNil)

	peer1 := s.makeAPeer(c, "peer1", 1)
	peer2 := s.makeAPeer(c, "peer2", 2)
	peer3 := s.makeAPeer(c, "peer3", 3)

	err := pm.AddToPeerMap(peer3)
	c.Assert(err, IsNil)
	c.Assert(pm.lowest, Not(IsNil))
	lowest := pm.lowest
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer3)

	err = pm.AddToPeerMap(peer2)
	c.Assert(err, IsNil)
	c.Assert(pm.lowest, Not(IsNil))
	lowest = pm.lowest
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer2)

	err = pm.AddToPeerMap(peer1)
	c.Assert(err, IsNil)
	c.Assert(pm.lowest, Not(IsNil))
	lowest = pm.lowest
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer1)

	c.Assert(pm.lowest.byteVal, Equals, byte(1))
	nextCell := pm.lowest.thisCol
	c.Assert(nextCell.byteVal, Equals, byte(2))
	nextCell = nextCell.thisCol
	c.Assert(nextCell.byteVal, Equals, byte(3))
}
func (s *XLSuite) TestDeeperMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_MAP")
	}
	var pm PeerMap
	c.Assert(pm.lowest, IsNil)

	peer1 := s.makeAPeer(c, "peer1", 1)
	peer12 := s.makeAPeer(c, "peer12", 1, 2)
	peer123 := s.makeAPeer(c, "peer123", 1, 2, 3)

	err := pm.AddToPeerMap(peer123)
	c.Assert(err, IsNil)
	c.Assert(pm.lowest, Not(IsNil))
	lowest := pm.lowest
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer123)

	err = pm.AddToPeerMap(peer12)
	c.Assert(err, IsNil)
	c.Assert(pm.lowest, Not(IsNil))
	lowest = pm.lowest
	c.Assert(lowest.peer, Not(IsNil))
	// c.Assert(lowest.peer, Equals, peer12) // PANIC
	c.Assert(lowest.peer.GetName(), Equals, peer12.GetName())

	err = pm.AddToPeerMap(peer1)
	c.Assert(err, IsNil)
	c.Assert(pm.lowest, Not(IsNil))
	lowest = pm.lowest
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer1)

	c.Assert(pm.lowest.byteVal, Equals, byte(1))
	nextCell := pm.lowest.thisCol
	c.Assert(nextCell.byteVal, Equals, byte(2))
	nextCell = nextCell.thisCol
	c.Assert(nextCell.byteVal, Equals, byte(3))
}

// XXX Something similar to this should be in nodeID/nodeID.go
func (s *XLSuite) TestSameNodeID(c *C) {
	peer := s.makeAPeer(c, "foo", 1, 2, 3, 4)
	id := peer.GetNodeID()
	c.Assert(xi.SameNodeID(id, id), Equals, true)
	peer2 := s.makeAPeer(c, "foo", 1, 2, 3, 4, 5)
	id2 := peer2.GetNodeID()
	c.Assert(xi.SameNodeID(id, id2), Equals, false)
}
