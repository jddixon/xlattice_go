package nodeID

// xlattice_go/nodeID/idMap_test.go

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	xr "github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print

const (
	MY_MAX_DEPTH = uint(16)
)

// -- MockBaseNode --------------------------------------------------
type MockBaseNode struct {
	id   *NodeID
	name string
}

func NewMockBaseNode(id *NodeID, name string) *MockBaseNode {
	return &MockBaseNode{
		id:   id,
		name: name,
	}
}
func (mbn *MockBaseNode) GetNodeID() *NodeID {
	return mbn.id
}

func (mbn *MockBaseNode) GetName() string {
	return mbn.name
}

// -- utility functions ---------------------------------------------
func (s *XLSuite) makePubKey(c *C) *rsa.PublicKey {
	key, err := rsa.GenerateKey(rand.Reader, 512) // 512 because cheaper
	c.Assert(err, IsNil)
	return &key.PublicKey
}

func (s *XLSuite) makePeerGivenID(c *C, rng *xr.PRNG, name string,
	id []byte) (member *MockBaseNode) {

	nodeID, err := New(id)
	c.Assert(err, IsNil)

	member = NewMockBaseNode(nodeID, name)
	c.Assert(err, IsNil)
	return
}
func (s *XLSuite) makeTopAndBottomBNI(c *C, rng *xr.PRNG) (
	topBNI, bottomBNI *MockBaseNode) {

	// top contains  a slice of SHA1_LEN 0xFF as NodeID
	t := make([]byte, SHA1_LEN)
	for i := 0; i < SHA1_LEN; i++ {
		t[i] = byte(0xff)
	}
	topBNI = s.makePeerGivenID(c, rng, "top", t)

	// bottom contains  a slice of SHA1_LEN zeroes as NodeID
	b := make([]byte, SHA1_LEN)
	bottomBNI = s.makePeerGivenID(c, rng, "bottom", b)

	return topBNI, bottomBNI
}

// Create a randomish Peer for use as a BNI, assigning it a
// nodeID based upon the variable-length list of ints passed

func (s *XLSuite) makeABNI(c *C, rng *xr.PRNG, name string, id ...int) (
	bni *MockBaseNode) {

	t := make([]byte, SHA1_LEN)
	for i := 0; i < len(id); i++ {
		t[i] = byte(id[i])
	}
	bni = s.makePeerGivenID(c, rng, name, t)
	c.Assert(bni, Not(IsNil))
	return
}

// -- tests proper --------------------------------------------------

func (s *XLSuite) TestIDMapTools(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_ID_MAP_TOOLS")
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
func (s *XLSuite) TestTopBottomIDMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_TOP_BOTTOM_MAP")
	}
	m, err := NewIDMap(MY_MAX_DEPTH)
	c.Assert(err, IsNil)
	c.Assert(m.MaxDepth, Equals, MY_MAX_DEPTH)

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

	for i := 0; i < 256; i++ {
		cell := m.Cells[i]
		c.Assert(cell.Next, IsNil)
		cellKey := cell.Key // a pointer
		if i == 0 {
			c.Assert(cellKey, NotNil) // XXX FAILS
			c.Assert(bytes.Equal(*cellKey, bottomKey), Equals, true)
		} else if i == 255 {
			c.Assert(cellKey, NotNil)
			c.Assert(bytes.Equal(*cellKey, topKey), Equals, true)
		} else {
			c.Assert(cellKey, IsNil)
		}
	}

	lowest := m.Cells[0]
	highest := m.Cells[255]

	c.Assert(lowest.Value, Equals, bottomBNI)
	c.Assert(highest.Value, Equals, topBNI)

	c.Assert(SameNodeID(
		lowest.Value.(*MockBaseNode).GetNodeID(), bottomBNI.GetNodeID()),
		Equals, true)

	c.Assert(lowest.Value.(*MockBaseNode).GetName(), Equals, "bottom")
	c.Assert(highest.Value.(*MockBaseNode).GetName(), Equals, "top")
}

func (s *XLSuite) TestShallowIDMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_MAP")
	}
	m, err := NewIDMap(MY_MAX_DEPTH)
	c.Assert(err, IsNil)
	c.Assert(m.MaxDepth, Equals, MY_MAX_DEPTH)

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
	cell3 := &m.Cells[3]
	c.Assert(cell3.Value, NotNil)
	c.Assert(cell3.Value.(*MockBaseNode).GetName(), Equals, baseNode3.GetName())

	// INSERT BNI 2 ------------------------------------------------
	err = m.Insert(key2, baseNode2)
	c.Assert(err, IsNil)
	cell2 := &m.Cells[2]
	c.Assert(cell2.Value, NotNil)
	c.Assert(cell2.Value.(*MockBaseNode).GetName(), Equals, baseNode2.GetName())

	// INSERT BNI 1 ------------------------------------------------
	err = m.Insert(key1, baseNode1)
	c.Assert(err, IsNil)
	cell1 := &m.Cells[1]
	c.Assert(cell1.Value, NotNil)
	c.Assert(cell1.Value.(*MockBaseNode).GetName(), Equals, baseNode1.GetName())

	c.Assert(err, IsNil)
	entryCount, _, _ := m.Size()
	c.Assert(entryCount, Equals, uint(3))

	// insert a duplicate -------------------------------------------
	err = m.Insert(key1, baseNode1)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(3))

	for i := 0; i < 256; i++ {
		cell := &m.Cells[i]
		c.Assert(cell.Next, IsNil)
		if i < 1 || i > 3 {
			c.Assert(cell.Key, IsNil)
			c.Assert(cell.Value, IsNil)
		}
	}
}

func (s *XLSuite) TestDeeperIDMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DEEPER_MAP")
	}
	m, err := NewIDMap(MY_MAX_DEPTH)
	c.Assert(err, IsNil)
	c.Assert(m.MaxDepth, Equals, MY_MAX_DEPTH)

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

	// add baseNode123 ================================================
	err = m.Insert(key123, baseNode123)
	c.Assert(err, IsNil)
	cell1 := &m.Cells[1]
	c.Assert(cell1.Next, IsNil)
	c.Assert(cell1.Key, NotNil)
	c.Assert(cell1.Value, NotNil)

	value, err = m.Find(key123)
	c.Assert(err, IsNil)
	c.Assert(value, NotNil)
	c.Assert(value, Equals, baseNode123)

	// now add baseNode12 ============================================
	// This should clear cell1 and create cell120 and cell123
	err = m.Insert(key12, baseNode12)
	c.Assert(err, IsNil)
	m1 := cell1.Next
	c.Assert(m1, NotNil)

	cell12 := &m1.Cells[2]
	c.Assert(cell12, NotNil)
	c.Assert(cell12.Next, NotNil)
	c.Assert(cell12.Key, IsNil)
	c.Assert(cell12.Value, IsNil)
	m12 := cell12.Next

	cell120 := &m12.Cells[0]
	c.Assert(cell120, NotNil)
	c.Assert(cell120.Next, IsNil)
	c.Assert(cell120.Key, NotNil)
	c.Assert(cell120.Value, NotNil)

	cell123 := &m12.Cells[3]
	c.Assert(cell123, NotNil)
	c.Assert(cell123.Next, IsNil)
	c.Assert(cell123.Key, NotNil)
	c.Assert(cell123.Value, NotNil)

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

	c.Assert(err, IsNil)
	entryCount, _, _ := m.Size()
	c.Assert(entryCount, Equals, uint(3))

	// insert a duplicate -------------------------------------------
	err = m.Insert(key1, baseNode1)
	c.Assert(err, IsNil)
	entryCount, _, _ = m.Size()
	c.Assert(entryCount, Equals, uint(3))

}

func (s *XLSuite) addAnID(c *C, m *IDMap, baseNode *MockBaseNode) {
	key := baseNode.GetNodeID().Value()
	c.Assert(key, NotNil)
	err := m.Insert(key, baseNode)
	c.Assert(err, IsNil)
}
func (s *XLSuite) findAnID(c *C, m *IDMap, baseNode *MockBaseNode) {
	key := baseNode.GetNodeID().Value()
	c.Assert(key, NotNil)
	p, err := m.Find(key)
	c.Assert(err, IsNil)
	c.Assert(p, NotNil)
	keyBack := p.(*MockBaseNode).GetNodeID().Value()
	c.Assert(bytes.Equal(key, keyBack), Equals, true)
}

func (s *XLSuite) TestFindID(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_ID")
	}
	m, err := NewIDMap(MY_MAX_DEPTH)
	c.Assert(err, IsNil)
	c.Assert(m.MaxDepth, Equals, MY_MAX_DEPTH)

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
	s.addAnID(c, m, baseNode123)
	s.addAnID(c, m, baseNode12)
	s.addAnID(c, m, baseNode1)

	s.addAnID(c, m, baseNode5)

	s.addAnID(c, m, baseNode4)
	s.addAnID(c, m, baseNode42)
	s.addAnID(c, m, baseNode423)

	s.addAnID(c, m, baseNode6)
	s.addAnID(c, m, baseNode623)
	s.addAnID(c, m, baseNode62)

	s.addAnID(c, m, baseNode0123)

	// adding duplicates should have no effect
	s.addAnID(c, m, baseNode4)
	s.addAnID(c, m, baseNode42)
	s.addAnID(c, m, baseNode423)

	// TODO: randomize order in which finding baseNodes is tested
	s.findAnID(c, m, baseNode0123)

	s.findAnID(c, m, baseNode1)
	s.findAnID(c, m, baseNode12)
	s.findAnID(c, m, baseNode123)

	s.findAnID(c, m, baseNode4)
	s.findAnID(c, m, baseNode42)
	s.findAnID(c, m, baseNode423)

	s.findAnID(c, m, baseNode6)
	s.findAnID(c, m, baseNode62)
	s.findAnID(c, m, baseNode623)
}
