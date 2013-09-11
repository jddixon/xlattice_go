package node

// xlattice_go/search/peer_map_test.go

import (
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print
var _ = rnglib.MakeSimpleRNG

const (
	SHA1_LEN = 20
)

func (s *XLSuite) makeTopAndBottom(c *C) (topPeer, bottomPeer *Peer) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < SHA1_LEN; i++ {
		t[i] = byte(0xf)
	}
	top, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(top, Not(IsNil))

	topPeer, err = NewNewPeer("top", top)
	c.Assert(err, IsNil)
	c.Assert(topPeer, Not(IsNil))

	bottom, err := xi.NewNodeID(make([]byte, SHA1_LEN))
	c.Assert(err, IsNil)
	c.Assert(bottom, Not(IsNil))

	bottomPeer, err = NewNewPeer("bottom", bottom)
	c.Assert(err, IsNil)
	c.Assert(bottomPeer, Not(IsNil))

	return topPeer, bottomPeer
}
func (s *XLSuite) makeAPeer(c *C, name string, id ...int) (peer *Peer) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < len(id); i++ {
		t[i] = byte(id[i])
	}
	nodeID, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))

	peer, err = NewNewPeer(name, nodeID)
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
		fmt.Println("TEST_TOP_BOTTOM_MAP")
	}

	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	topPeer, bottomPeer := s.makeTopAndBottom(c)
	err := pm.AddToPeerMap(topPeer)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	lowest := pm.nextCol
	c.Assert(lowest.peer, Not(IsNil))
	// THESE THREE TESTS ARE LOGICALLY EQUIVALENT ----------------------
	c.Assert(lowest.peer, Equals, topPeer) // succeeds ...
	c.Assert(xi.SameNodeID(lowest.peer.GetNodeID(), topPeer.GetNodeID()),
		Equals, true)
	// XXX This fails, but it's a bug in Peer.Equal()
	// c.Assert(topPeer.Equal(lowest.peer), Equals, true)
	// END LOGICALLY EQUIVALENT -----------------------------------------
	c.Assert(lowest.peer.GetName(), Equals, "top")

	// We expect that bottomPeer will become the lowest with its
	// higher field pointing at topPeer.
	err = pm.AddToPeerMap(bottomPeer)
	c.Assert(err, IsNil)
	lowest = pm.nextCol
	// c.Assert(bottomPeer.Equal(lowest.peer), Equals, true)   // FAILS
	c.Assert(lowest.peer.GetName(), Equals, "bottom") // XXX gets 'top'
}
func (s *XLSuite) TestShallowMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_MAP")
	}
	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	peer1 := s.makeAPeer(c, "peer1", 1)
	peer2 := s.makeAPeer(c, "peer2", 2)
	peer3 := s.makeAPeer(c, "peer3", 3)

	// ADD PEER 3 ---------------------------------------------------
	err := pm.AddToPeerMap(peer3)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	cell3 := pm.nextCol
	c.Assert(cell3.byteVal, Equals, byte(3))
	c.Assert(cell3.peer, Not(IsNil))
	c.Assert(cell3.peer.GetName(), Equals, peer3.GetName())

	// INSERT PEER 2 ------------------------------------------------
	err = pm.AddToPeerMap(peer2)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	cell2 := pm.nextCol
	c.Assert(cell2.byteVal, Equals, byte(2)) // FAILS, is 3
	c.Assert(cell2.thisCol.byteVal, Equals, byte(3))
	c.Assert(cell2.peer, Not(IsNil))
	c.Assert(cell2.peer.GetName(), Equals, peer2.GetName()) // FAILS

	// DumpPeerMap(&pm, "dump of shallow map, peers 3 and 2")

	// INSERT PEER 1 ------------------------------------------------
	err = pm.AddToPeerMap(peer1)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	cell1 := pm.nextCol
	c.Assert(cell1.byteVal, Equals, byte(1))
	c.Assert(cell1.peer, Not(IsNil))
	c.Assert(cell1.peer.GetName(), Equals, peer1.GetName())

	// DumpPeerMap(&pm, "dump of shallow map, peers 3,2,1")

	rootCell := pm.nextCol
	c.Assert(rootCell.byteVal, Equals, byte(1))
	c.Assert(rootCell.peer.GetName(), Equals, "peer1")
	nextCell := rootCell.thisCol
	c.Assert(nextCell, Not(IsNil))
	c.Assert(nextCell.byteVal, Equals, byte(2))
	nextCell = nextCell.thisCol
	c.Assert(nextCell.byteVal, Equals, byte(3))
}

func (s *XLSuite) TestDeeperMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DEEPER_MAP")
	}
	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	peer1 := s.makeAPeer(c, "peer1", 1)
	peer12 := s.makeAPeer(c, "peer12", 1, 2)
	peer123 := s.makeAPeer(c, "peer123", 1, 2, 3)

	// add peer123 ================================================
	err := pm.AddToPeerMap(peer123)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	lowest := pm.nextCol
	c.Assert(lowest.peer, Not(IsNil))
	c.Assert(lowest.peer, Equals, peer123)

	// now add peer12 ============================================
	err = pm.AddToPeerMap(peer12)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	col0 := pm.nextCol

	// DumpPeerMap(&pm, "after peer123 then peer12 added")

	// column 0 check - expect an empty cell
	c.Assert(col0.thisCol, IsNil)
	c.Assert(col0.peer, IsNil)

	// column 1 check - another empty cell
	col1 := col0.nextCol
	c.Assert(col1, Not(IsNil))
	c.Assert(col1.thisCol, IsNil)
	c.Assert(col1.peer, IsNil)

	// column 2a checks - peer12 with peer123 on the nextCol chain
	col2a := col1.nextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.nextCol, IsNil)
	c.Assert(col2a.peer, Not(IsNil))
	c.Assert(col2a.peer.GetName(), Equals, "peer12")

	// column 2b checks
	col2b := col2a.thisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.nextCol, IsNil)
	c.Assert(col2b.thisCol, IsNil)
	c.Assert(col2b.peer, Not(IsNil))
	c.Assert(col2b.peer.GetName(), Equals, "peer123")

	// now add peer1 =============================================
	err = pm.AddToPeerMap(peer1)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
	col0 = pm.nextCol

	// DumpPeerMap(&pm, "after peer123, peer12, then peer1 added")

	// column 0 checks - an empty cell
	c.Assert(col0.peer, IsNil)
	c.Assert(col0.thisCol, IsNil)

	// column 1a check -
	col1a := col0.nextCol
	c.Assert(col1a, Not(IsNil))
	c.Assert(col1a.nextCol, IsNil)
	c.Assert(col1a.thisCol, Not(IsNil))
	c.Assert(col1a.peer, Not(IsNil))
	c.Assert(col1a.peer, Equals, peer1)
	c.Assert(col1a.peer.GetName(), Equals, "peer1")

	// column 1b checks - another empty cell
	col1b := col1a.thisCol
	c.Assert(col1b.peer, IsNil)
	c.Assert(col1b.thisCol, IsNil)

	// column 2a checks - peer12 with peer123 on the nextCol chain
	col2a = col1b.nextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.nextCol, IsNil)
	c.Assert(col2a.peer, Not(IsNil))
	c.Assert(col2a.peer.GetName(), Equals, "peer12")

	// column 2b checks
	col2b = col2a.thisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.nextCol, IsNil)
	c.Assert(col2b.thisCol, IsNil)
	c.Assert(col2b.peer, Not(IsNil))
	c.Assert(col2b.peer.GetName(), Equals, "peer123")

	c.Assert(col0.byteVal, Equals, byte(1))
	c.Assert(col1a.byteVal, Equals, byte(0))
	c.Assert(col1b.byteVal, Equals, byte(2))
	c.Assert(col2a.byteVal, Equals, byte(0))
	c.Assert(col2b.byteVal, Equals, byte(3))

	// add 123, then 1, then 12 ----------------------------------

	// XXX STUB XXX

}

func (s *XLSuite) addAPeer(c *C, pm *PeerMap, peer *Peer) {
	err := pm.AddToPeerMap(peer)
	c.Assert(err, IsNil)
	c.Assert(pm.nextCol, Not(IsNil))
}
func (s *XLSuite) findAPeer(c *C, pm *PeerMap, peer *Peer) {
	nodeID := peer.GetNodeID()
	d := nodeID.Value()
	c.Assert(d, Not(IsNil))
	p := pm.FindPeer(d)
	// DEBUG
	if p == nil {
		fmt.Printf("can't find a match for %d.%d.%d.%d\n", d[0], d[1], d[2], d[3])
	}
	// END
	c.Assert(p, Not(IsNil))
	nodeIDBack := p.GetNodeID()
	c.Assert(xi.SameNodeID(nodeID, nodeIDBack), Equals, true)

}
func (s *XLSuite) TestFindFlatPeers(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_FLAT_PEERS")
	}
	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	peer1 := s.makeAPeer(c, "peer1", 1)
	peer2 := s.makeAPeer(c, "peer2", 2)
	peer4 := s.makeAPeer(c, "peer4", 4)
	peer5 := s.makeAPeer(c, "peer5", 5)
	peer6 := s.makeAPeer(c, "peer6", 6)

	// TODO: randomize order in which peers are added

	// ADD 1 AND THEN 5 ---------------------------------------------
	s.addAPeer(c, &pm, peer1)
	s.addAPeer(c, &pm, peer5)

	cell1 := pm.nextCol
	c.Assert(cell1.pred, Equals, &pm.PeerMapCell)
	c.Assert(cell1.nextCol, IsNil)

	cell5 := cell1.thisCol
	c.Assert(cell5, Not(IsNil)) // FAILS
	c.Assert(cell5.byteVal, Equals, byte(5))
	c.Assert(cell5.pred, Equals, cell1)
	c.Assert(cell5.nextCol, IsNil)
	c.Assert(cell5.thisCol, IsNil)

	// INSERT 4 -----------------------------------------------------
	s.addAPeer(c, &pm, peer4)

	cell4 := cell1.thisCol
	c.Assert(cell4.byteVal, Equals, byte(4))
	c.Assert(cell4.pred, Equals, cell1)
	c.Assert(cell4.nextCol, IsNil)
	c.Assert(cell4.thisCol, Equals, cell5)
	c.Assert(cell5.pred, Equals, cell4)

	// ADD 6 --------------------------------------------------------
	s.addAPeer(c, &pm, peer6)

	cell6 := cell5.thisCol
	c.Assert(cell6.byteVal, Equals, byte(6))
	c.Assert(cell6.pred, Equals, cell5)
	c.Assert(cell6.nextCol, IsNil)
	c.Assert(cell6.thisCol, IsNil)

	// INSERT 2 -----------------------------------------------------
	s.addAPeer(c, &pm, peer2)

	cell2 := cell1.thisCol
	c.Assert(cell2.byteVal, Equals, byte(2))
	c.Assert(cell2.pred, Equals, cell1)
	c.Assert(cell2.nextCol, IsNil)
	c.Assert(cell2.thisCol, Equals, cell4)
	c.Assert(cell4.pred, Equals, cell2)

	// DumpPeerMap(&pm, "after adding peer2")

	// TODO: randomize order in which finding peers is tested
	s.findAPeer(c, &pm, peer1)
	s.findAPeer(c, &pm, peer2)
	s.findAPeer(c, &pm, peer4)
	s.findAPeer(c, &pm, peer5)
	s.findAPeer(c, &pm, peer6)
}
func (s *XLSuite) TestFindPeer(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_PEER")
	}
	var pm PeerMap
	c.Assert(pm.nextCol, IsNil)

	peer0123 := s.makeAPeer(c, "peer0123", 0, 1, 2, 3)
	peer1 := s.makeAPeer(c, "peer1", 1)
	peer12 := s.makeAPeer(c, "peer12", 1, 2)
	peer123 := s.makeAPeer(c, "peer123", 1, 2, 3)
	peer4 := s.makeAPeer(c, "peer4", 4)
	peer42 := s.makeAPeer(c, "peer42", 4, 2)
	peer423 := s.makeAPeer(c, "peer423", 4, 2, 3)
	// peer5 := s.makeAPeer(c, "peer5", 5)
	peer6 := s.makeAPeer(c, "peer6", 6)
	peer62 := s.makeAPeer(c, "peer62", 6, 2)
	peer623 := s.makeAPeer(c, "peer623", 6, 2, 3)

	// TODO: randomize order in which peers are added
	s.addAPeer(c, &pm, peer123)
	s.addAPeer(c, &pm, peer12)
	s.addAPeer(c, &pm, peer1)
	//DumpPeerMap(&pm, "after adding peer1, peer12, peer123, before peer4")

	// s.addAPeer(c, &pm, peer5)
	// DumpPeerMap(&pm, "after adding peer5")

	s.addAPeer(c, &pm, peer4)
	s.addAPeer(c, &pm, peer42)
	s.addAPeer(c, &pm, peer423)
	// DumpPeerMap(&pm, "after adding peer4, peer42, peer423")

	s.addAPeer(c, &pm, peer6)
	// DumpPeerMap(&pm, "after adding peer6")
	s.addAPeer(c, &pm, peer623)
	//DumpPeerMap(&pm, "after adding peer623")
	s.addAPeer(c, &pm, peer62)
	//DumpPeerMap(&pm, "after adding peer62")

	s.addAPeer(c, &pm, peer0123)
	//DumpPeerMap(&pm, "after adding peer0123")

	// adding duplicates should have no effect
	s.addAPeer(c, &pm, peer4)
	s.addAPeer(c, &pm, peer42)
	s.addAPeer(c, &pm, peer423)

	// TODO: randomize order in which finding peers is tested
	s.findAPeer(c, &pm, peer0123) // XXX

	s.findAPeer(c, &pm, peer1)
	s.findAPeer(c, &pm, peer12)
	s.findAPeer(c, &pm, peer123)

	s.findAPeer(c, &pm, peer4)
	s.findAPeer(c, &pm, peer42)
	s.findAPeer(c, &pm, peer423)

	s.findAPeer(c, &pm, peer6)
	s.findAPeer(c, &pm, peer62)
	s.findAPeer(c, &pm, peer623)
}
