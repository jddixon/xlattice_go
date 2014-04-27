package nodeID

// xlattice_go/nodeID/idMapHAMT_test.go

import (
	"bytes"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "gopkg.in/check.v1"
)

var _ = fmt.Print

///////////////////////////////////////////////////
// XXX These tests reply upon code in idMap_test.go
///////////////////////////////////////////////////

// -- utility functions ---------------------------------------------
func (s *XLSuite) addIDToHAMT(c *C, m *IDMapHAMT, baseNode *MockBaseNode) {
	key := baseNode.GetNodeID().Value()
	c.Assert(key, NotNil)
	err := m.Insert(key, baseNode)
	c.Assert(err, IsNil)
}
func (s *XLSuite) findIDInHAMT(c *C, m *IDMapHAMT, baseNode *MockBaseNode) {
	key := baseNode.GetNodeID().Value()
	c.Assert(key, NotNil)
	p, err := m.Find(key)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
	keyBack := p.(*MockBaseNode).GetNodeID().Value()
	c.Assert(bytes.Equal(key, keyBack), Equals, true)
}

// -- tests proper --------------------------------------------------

func (s *XLSuite) TestIDMapHAMTTools(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_ID_MAP_HAMT_TOOLS")
	}
	rng := xr.MakeSimpleRNG()
	threeBaseNode := s.makeABNI(c, rng, "threeBaseNode", 1, 2, 3)
	nodeID := threeBaseNode.GetNodeID()
	value := nodeID.Value()
	c.Assert(threeBaseNode.GetName(), Equals, "threeBaseNode")
	c.Assert(value[0], Equals, byte(1))
	c.Assert(value[1], Equals, byte(2))
	c.Assert(value[2], Equals, byte(3))
	for i := 3; i < SHA1_LEN; i++ {
		c.Assert(value[i], Equals, byte(0))
	}

}
func (s *XLSuite) TestHTopBottomIDMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_TOP_BOTTOM_HAMT_MAP")
	}
	var err error

	m := NewNewIDMapHAMT()

	rng := xr.MakeSimpleRNG()
	topBNI, bottomBNI := s.makeTopAndBottomBNI(c, rng)
	bottomKey := bottomBNI.GetNodeID().Value()
	topKey := topBNI.GetNodeID().Value()

	err = m.Insert(topKey, topBNI)
	c.Assert(err, IsNil)
	err = m.Insert(bottomKey, bottomBNI)
	c.Assert(err, IsNil)
	entryCount, _, _ := m.Size()
	c.Assert(entryCount, Equals, uint(2))

	// insert a duplicate
	err = m.Insert(bottomKey, bottomBNI)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(2))

}

func (s *XLSuite) TestHShallowIDMapHAMT(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_HAMT_MAP")
	}
	var err error
	m := NewNewIDMapHAMT()

	rng := xr.MakeSimpleRNG()
	// 1 or 2 or 3 is first digit of key, guaranteeing shallownes
	baseNode1 := s.makeABNI(c, rng, "baseNode1", 1)
	baseNode2 := s.makeABNI(c, rng, "baseNode2", 2)
	baseNode3 := s.makeABNI(c, rng, "baseNode3", 3)

	key1 := baseNode1.GetNodeID().Value()
	key2 := baseNode2.GetNodeID().Value()
	key3 := baseNode3.GetNodeID().Value()

	// INSERT BNI 3 -------------------------------------------------
	err = m.Insert(key3, baseNode3)
	c.Assert(err, IsNil)

	// INSERT BNI 2 ------------------------------------------------
	err = m.Insert(key2, baseNode2)
	c.Assert(err, IsNil)

	// INSERT BNI 1 ------------------------------------------------
	err = m.Insert(key1, baseNode1)
	c.Assert(err, IsNil)

	c.Assert(err, IsNil)
	entryCount, _, _ := m.Size()
	c.Assert(entryCount, Equals, uint(3))

	// insert a duplicate -------------------------------------------
	err = m.Insert(key1, baseNode1)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(3))
}

func (s *XLSuite) TestHDeeperIDMapHAMT(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DEEPER_HAMT_MAP")
	}
	var err error

	m := NewNewIDMapHAMT()

	rng := xr.MakeSimpleRNG()
	baseNode1 := s.makeABNI(c, rng, "baseNode1", 1)
	baseNode12 := s.makeABNI(c, rng, "baseNode12", 1, 2)
	baseNode123 := s.makeABNI(c, rng, "baseNode123", 1, 2, 3)

	key1 := baseNode1.GetNodeID().Value()
	key12 := baseNode12.GetNodeID().Value()
	key123 := baseNode123.GetNodeID().Value()

	value, err := m.Find(key1)
	c.Assert(err, IsNil)
	c.Assert(value, IsNil)

	value, err = m.Find(key12)
	c.Assert(err, IsNil)
	c.Assert(value, IsNil)

	value, err = m.Find(key123)
	c.Assert(err, IsNil)
	c.Assert(value, IsNil)

	entryCount, _, _ := m.Size()
	c.Assert(entryCount, Equals, uint(0))

	// add baseNode123 ================================================
	err = m.Insert(key123, baseNode123)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(1))

	value, err = m.Find(key123)
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, baseNode123)

	// now add baseNode12 ============================================
	err = m.Insert(key12, baseNode12)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(2))

	value, err = m.Find(key123)
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, baseNode123)

	value, err = m.Find(key12)
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, baseNode12)

	value, err = m.Find(key1)
	c.Assert(err, IsNil)
	c.Assert(value, IsNil)

	// now add baseNode1 =============================================
	err = m.Insert(key1, baseNode1)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(3))

	value, err = m.Find(key1)
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, baseNode1)

	value, err = m.Find(key123)
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, baseNode123)

	value, err = m.Find(key12)
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, baseNode12)

	// insert a duplicate -------------------------------------------
	err = m.Insert(key1, baseNode1)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(3))
}

func (s *XLSuite) TestHFindID(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_ID")
	}
	var err error

	m := NewNewIDMapHAMT()
	c.Assert(err, IsNil)

	rng := xr.MakeSimpleRNG()
	baseNode0123 := s.makeABNI(c, rng, "baseNode0123", 0, 1, 2, 3)
	baseNode1 := s.makeABNI(c, rng, "baseNode1", 1)
	baseNode12 := s.makeABNI(c, rng, "baseNode12", 1, 2)
	baseNode123 := s.makeABNI(c, rng, "baseNode123", 1, 2, 3)
	baseNode4 := s.makeABNI(c, rng, "baseNode4", 4)
	baseNode42 := s.makeABNI(c, rng, "baseNode42", 4, 2)
	baseNode423 := s.makeABNI(c, rng, "baseNode423", 4, 2, 3)
	baseNode5 := s.makeABNI(c, rng, "baseNode5", 5)
	baseNode6 := s.makeABNI(c, rng, "baseNode6", 6)
	baseNode62 := s.makeABNI(c, rng, "baseNode62", 6, 2)
	baseNode623 := s.makeABNI(c, rng, "baseNode623", 6, 2, 3)

	// TODO: randomize order in which baseNodes are added
	s.addIDToHAMT(c, m, baseNode123)
	s.addIDToHAMT(c, m, baseNode12)
	s.addIDToHAMT(c, m, baseNode1)

	s.addIDToHAMT(c, m, baseNode5)

	s.addIDToHAMT(c, m, baseNode4)
	s.addIDToHAMT(c, m, baseNode42)
	s.addIDToHAMT(c, m, baseNode423)

	s.addIDToHAMT(c, m, baseNode6)
	s.addIDToHAMT(c, m, baseNode623)
	s.addIDToHAMT(c, m, baseNode62)

	s.addIDToHAMT(c, m, baseNode0123)

	// adding duplicates should have no effect
	s.addIDToHAMT(c, m, baseNode4)
	s.addIDToHAMT(c, m, baseNode42)
	s.addIDToHAMT(c, m, baseNode423)

	// TODO: randomize order in which finding baseNodes is tested
	s.findIDInHAMT(c, m, baseNode0123)

	s.findIDInHAMT(c, m, baseNode1)
	s.findIDInHAMT(c, m, baseNode12)
	s.findIDInHAMT(c, m, baseNode123)

	s.findIDInHAMT(c, m, baseNode4)
	s.findIDInHAMT(c, m, baseNode42)
	s.findIDInHAMT(c, m, baseNode423)

	s.findIDInHAMT(c, m, baseNode6)
	s.findIDInHAMT(c, m, baseNode62)
	s.findIDInHAMT(c, m, baseNode623)
}
