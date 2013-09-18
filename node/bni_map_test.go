package node

// xlattice_go/node/bni_map_test.go

import (
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
	"github.com/jddixon/xlattice_go/rnglib"
	. "launchpad.net/gocheck"
)

var _ = fmt.Print
var _ = rnglib.MakeSimpleRNG

func (s *XLSuite) makeTopAndBottomBN(c *C) (topBaseNode, bottomBaseNode *BaseNode) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < SHA1_LEN; i++ {
		t[i] = byte(0xf)
	}
	top, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(top, Not(IsNil))

	topBaseNode, err = NewNewBaseNode("top", top)
	c.Assert(err, IsNil)
	c.Assert(topBaseNode, Not(IsNil))

	bottom, err := xi.NewNodeID(make([]byte, SHA1_LEN))
	c.Assert(err, IsNil)
	c.Assert(bottom, Not(IsNil))

	bottomBaseNode, err = NewNewBaseNode("bottom", bottom)
	c.Assert(err, IsNil)
	c.Assert(bottomBaseNode, Not(IsNil))

	return topBaseNode, bottomBaseNode
}
func (s *XLSuite) makeABaseNode(c *C, name string, id ...int) (baseNode *BaseNode) {
	t := make([]byte, SHA1_LEN)
	for i := 0; i < len(id); i++ {
		t[i] = byte(id[i])
	}
	nodeID, err := xi.NewNodeID(t)
	c.Assert(err, IsNil)
	c.Assert(nodeID, Not(IsNil))

	baseNode, err = NewNewBaseNode(name, nodeID)
	c.Assert(err, IsNil)
	c.Assert(baseNode, Not(IsNil))
	return
}
func (s *XLSuite) TestBNIMapTools(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_BASE_NODE_MAP_TOOLS")
	}
	threeBaseNode := s.makeABaseNode(c, "threeBaseNode", 1, 2, 3)
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
func (s *XLSuite) TestTopBottomBNMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_TOP_BOTTOM_MAP")
	}

	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	topBaseNode, bottomBaseNode := s.makeTopAndBottomBN(c)
	err := pm.AddToBNIMap(topBaseNode)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	lowest := pm.NextCol
	c.Assert(lowest.CellNode, Not(IsNil))
	// THESE THREE TESTS ARE LOGICALLY EQUIVALENT ----------------------
	c.Assert(lowest.CellNode, Equals, topBaseNode) // succeeds ...
	c.Assert(xi.SameNodeID(lowest.CellNode.GetNodeID(), topBaseNode.GetNodeID()),
		Equals, true)
	// XXX This fails, but it's a bug in BaseNode.Equal()
	// c.Assert(topBaseNode.Equal(lowest.CellNode), Equals, true)
	// END LOGICALLY EQUIVALENT -----------------------------------------
	c.Assert(lowest.CellNode.GetName(), Equals, "top")

	// We expect that bottomBaseNode will become the lowest with its
	// higher field pointing at topBaseNode.
	err = pm.AddToBNIMap(bottomBaseNode)
	c.Assert(err, IsNil)
	lowest = pm.NextCol
	// c.Assert(bottomBaseNode.Equal(lowest.CellNode), Equals, true)   // FAILS
	c.Assert(lowest.CellNode.GetName(), Equals, "bottom") // XXX gets 'top'
}
func (s *XLSuite) TestShallowBNMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_SHALLOW_MAP")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	baseNode1 := s.makeABaseNode(c, "baseNode1", 1)
	baseNode2 := s.makeABaseNode(c, "baseNode2", 2)
	baseNode3 := s.makeABaseNode(c, "baseNode3", 3)

	// ADD BASE_NODE 3 ---------------------------------------------------
	err := pm.AddToBNIMap(baseNode3)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	cell3 := pm.NextCol
	c.Assert(cell3.ByteVal, Equals, byte(3))
	c.Assert(cell3.CellNode, Not(IsNil))
	c.Assert(cell3.CellNode.GetName(), Equals, baseNode3.GetName())

	// INSERT BASE_NODE 2 ------------------------------------------------
	err = pm.AddToBNIMap(baseNode2)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	cell2 := pm.NextCol
	c.Assert(cell2.ByteVal, Equals, byte(2)) // FAILS, is 3
	c.Assert(cell2.ThisCol.ByteVal, Equals, byte(3))
	c.Assert(cell2.CellNode, Not(IsNil))
	c.Assert(cell2.CellNode.GetName(), Equals, baseNode2.GetName()) // FAILS

	// DumpBNIMap(&pm, "dump of shallow map, baseNodes 3 and 2")

	// INSERT BASE_NODE 1 ------------------------------------------------
	err = pm.AddToBNIMap(baseNode1)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	cell1 := pm.NextCol
	c.Assert(cell1.ByteVal, Equals, byte(1))
	c.Assert(cell1.CellNode, Not(IsNil))
	c.Assert(cell1.CellNode.GetName(), Equals, baseNode1.GetName())

	// DumpBNIMap(&pm, "dump of shallow map, baseNodes 3,2,1")

	rootCell := pm.NextCol
	c.Assert(rootCell.ByteVal, Equals, byte(1))
	c.Assert(rootCell.CellNode.GetName(), Equals, "baseNode1")
	nextCell := rootCell.ThisCol
	c.Assert(nextCell, Not(IsNil))
	c.Assert(nextCell.ByteVal, Equals, byte(2))
	nextCell = nextCell.ThisCol
	c.Assert(nextCell.ByteVal, Equals, byte(3))
}

func (s *XLSuite) TestDeeperBNMap(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_DEEPER_MAP")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	baseNode1 := s.makeABaseNode(c, "baseNode1", 1)
	baseNode12 := s.makeABaseNode(c, "baseNode12", 1, 2)
	baseNode123 := s.makeABaseNode(c, "baseNode123", 1, 2, 3)

	// add baseNode123 ================================================
	err := pm.AddToBNIMap(baseNode123)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	lowest := pm.NextCol
	c.Assert(lowest.CellNode, Not(IsNil))
	c.Assert(lowest.CellNode, Equals, baseNode123)

	// now add baseNode12 ============================================
	err = pm.AddToBNIMap(baseNode12)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	col0 := pm.NextCol

	// DumpBNIMap(&pm, "after baseNode123 then baseNode12 added")

	// column 0 check - expect an empty cell
	c.Assert(col0.ThisCol, IsNil)
	c.Assert(col0.CellNode, IsNil)

	// column 1 check - another empty cell
	col1 := col0.NextCol
	c.Assert(col1, Not(IsNil))
	c.Assert(col1.ThisCol, IsNil)
	c.Assert(col1.CellNode, IsNil)

	// column 2a checks - baseNode12 with baseNode123 on the NextCol chain
	col2a := col1.NextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.NextCol, IsNil)
	c.Assert(col2a.CellNode, Not(IsNil))
	c.Assert(col2a.CellNode.GetName(), Equals, "baseNode12")

	// column 2b checks
	col2b := col2a.ThisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.NextCol, IsNil)
	c.Assert(col2b.ThisCol, IsNil)
	c.Assert(col2b.CellNode, Not(IsNil))
	c.Assert(col2b.CellNode.GetName(), Equals, "baseNode123")

	// now add baseNode1 =============================================
	err = pm.AddToBNIMap(baseNode1)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
	col0 = pm.NextCol

	// DumpBNIMap(&pm, "after baseNode123, baseNode12, then baseNode1 added")

	// column 0 checks - an empty cell
	c.Assert(col0.CellNode, IsNil)
	c.Assert(col0.ThisCol, IsNil)

	// column 1a check -
	col1a := col0.NextCol
	c.Assert(col1a, Not(IsNil))
	c.Assert(col1a.NextCol, IsNil)
	c.Assert(col1a.ThisCol, Not(IsNil))
	c.Assert(col1a.CellNode, Not(IsNil))
	c.Assert(col1a.CellNode, Equals, baseNode1)
	c.Assert(col1a.CellNode.GetName(), Equals, "baseNode1")

	// column 1b checks - another empty cell
	col1b := col1a.ThisCol
	c.Assert(col1b.CellNode, IsNil)
	c.Assert(col1b.ThisCol, IsNil)

	// column 2a checks - baseNode12 with baseNode123 on the NextCol chain
	col2a = col1b.NextCol
	c.Assert(col2a, Not(IsNil))
	c.Assert(col2a.NextCol, IsNil)
	c.Assert(col2a.CellNode, Not(IsNil))
	c.Assert(col2a.CellNode.GetName(), Equals, "baseNode12")

	// column 2b checks
	col2b = col2a.ThisCol
	c.Assert(col2b, Not(IsNil))
	c.Assert(col2b.NextCol, IsNil)
	c.Assert(col2b.ThisCol, IsNil)
	c.Assert(col2b.CellNode, Not(IsNil))
	c.Assert(col2b.CellNode.GetName(), Equals, "baseNode123")

	c.Assert(col0.ByteVal, Equals, byte(1))
	c.Assert(col1a.ByteVal, Equals, byte(0))
	c.Assert(col1b.ByteVal, Equals, byte(2))
	c.Assert(col2a.ByteVal, Equals, byte(0))
	c.Assert(col2b.ByteVal, Equals, byte(3))

	// add 123, then 1, then 12 ----------------------------------

	// XXX STUB XXX

}

func (s *XLSuite) addABaseNode(c *C, pm *BNIMap, baseNode *BaseNode) {
	err := pm.AddToBNIMap(baseNode)
	c.Assert(err, IsNil)
	c.Assert(pm.NextCol, Not(IsNil))
}
func (s *XLSuite) findABaseNode(c *C, pm *BNIMap, baseNode *BaseNode) {
	nodeID := baseNode.GetNodeID()
	d := nodeID.Value()
	c.Assert(d, Not(IsNil))
	p := pm.FindBaseNode(d)
	// DEBUG
	if p == nil {
		fmt.Printf("can't find a match for %d.%d.%d.%d\n", d[0], d[1], d[2], d[3])
	}
	// END
	c.Assert(p, Not(IsNil))
	nodeIDBack := p.GetNodeID()
	c.Assert(xi.SameNodeID(nodeID, nodeIDBack), Equals, true)

}
func (s *XLSuite) TestFindFlatBaseNodes(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_FLAT_BASE_NODES")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	baseNode1 := s.makeABaseNode(c, "baseNode1", 1)
	baseNode2 := s.makeABaseNode(c, "baseNode2", 2)
	baseNode4 := s.makeABaseNode(c, "baseNode4", 4)
	baseNode5 := s.makeABaseNode(c, "baseNode5", 5)
	baseNode6 := s.makeABaseNode(c, "baseNode6", 6)

	// TODO: randomize order in which baseNodes are added

	// ADD 1 AND THEN 5 ---------------------------------------------
	s.addABaseNode(c, &pm, baseNode1)
	s.addABaseNode(c, &pm, baseNode5)

	cell1 := pm.NextCol
	c.Assert(cell1.Pred, Equals, &pm.BNIMapCell)
	c.Assert(cell1.NextCol, IsNil)

	cell5 := cell1.ThisCol
	c.Assert(cell5, Not(IsNil)) // FAILS
	c.Assert(cell5.ByteVal, Equals, byte(5))
	c.Assert(cell5.Pred, Equals, cell1)
	c.Assert(cell5.NextCol, IsNil)
	c.Assert(cell5.ThisCol, IsNil)

	// INSERT 4 -----------------------------------------------------
	s.addABaseNode(c, &pm, baseNode4)

	cell4 := cell1.ThisCol
	c.Assert(cell4.ByteVal, Equals, byte(4))
	c.Assert(cell4.Pred, Equals, cell1)
	c.Assert(cell4.NextCol, IsNil)
	c.Assert(cell4.ThisCol, Equals, cell5)
	c.Assert(cell5.Pred, Equals, cell4)

	// ADD 6 --------------------------------------------------------
	s.addABaseNode(c, &pm, baseNode6)

	cell6 := cell5.ThisCol
	c.Assert(cell6.ByteVal, Equals, byte(6))
	c.Assert(cell6.Pred, Equals, cell5)
	c.Assert(cell6.NextCol, IsNil)
	c.Assert(cell6.ThisCol, IsNil)

	// INSERT 2 -----------------------------------------------------
	s.addABaseNode(c, &pm, baseNode2)

	cell2 := cell1.ThisCol
	c.Assert(cell2.ByteVal, Equals, byte(2))
	c.Assert(cell2.Pred, Equals, cell1)
	c.Assert(cell2.NextCol, IsNil)
	c.Assert(cell2.ThisCol, Equals, cell4)
	c.Assert(cell4.Pred, Equals, cell2)

	// DumpBNIMap(&pm, "after adding baseNode2")

	// TODO: randomize order in which finding baseNodes is tested
	s.findABaseNode(c, &pm, baseNode1)
	s.findABaseNode(c, &pm, baseNode2)
	s.findABaseNode(c, &pm, baseNode4)
	s.findABaseNode(c, &pm, baseNode5)
	s.findABaseNode(c, &pm, baseNode6)
}
func (s *XLSuite) TestFindBaseNode(c *C) {
	if VERBOSITY > 0 {
		fmt.Println("TEST_FIND_BASE_NODE")
	}
	var pm BNIMap
	c.Assert(pm.NextCol, IsNil)

	baseNode0123 := s.makeABaseNode(c, "baseNode0123", 0, 1, 2, 3)
	baseNode1 := s.makeABaseNode(c, "baseNode1", 1)
	baseNode12 := s.makeABaseNode(c, "baseNode12", 1, 2)
	baseNode123 := s.makeABaseNode(c, "baseNode123", 1, 2, 3)
	baseNode4 := s.makeABaseNode(c, "baseNode4", 4)
	baseNode42 := s.makeABaseNode(c, "baseNode42", 4, 2)
	baseNode423 := s.makeABaseNode(c, "baseNode423", 4, 2, 3)
	// baseNode5 := s.makeABaseNode(c, "baseNode5", 5)
	baseNode6 := s.makeABaseNode(c, "baseNode6", 6)
	baseNode62 := s.makeABaseNode(c, "baseNode62", 6, 2)
	baseNode623 := s.makeABaseNode(c, "baseNode623", 6, 2, 3)

	// TODO: randomize order in which baseNodes are added
	s.addABaseNode(c, &pm, baseNode123)
	s.addABaseNode(c, &pm, baseNode12)
	s.addABaseNode(c, &pm, baseNode1)
	//DumpBNIMap(&pm, "after adding baseNode1, baseNode12, baseNode123, before baseNode4")

	// s.addABaseNode(c, &pm, baseNode5)
	// DumpBNIMap(&pm, "after adding baseNode5")

	s.addABaseNode(c, &pm, baseNode4)
	s.addABaseNode(c, &pm, baseNode42)
	s.addABaseNode(c, &pm, baseNode423)
	// DumpBNIMap(&pm, "after adding baseNode4, baseNode42, baseNode423")

	s.addABaseNode(c, &pm, baseNode6)
	// DumpBNIMap(&pm, "after adding baseNode6")
	s.addABaseNode(c, &pm, baseNode623)
	//DumpBNIMap(&pm, "after adding baseNode623")
	s.addABaseNode(c, &pm, baseNode62)
	//DumpBNIMap(&pm, "after adding baseNode62")

	s.addABaseNode(c, &pm, baseNode0123)
	//DumpBNIMap(&pm, "after adding baseNode0123")

	// adding duplicates should have no effect
	s.addABaseNode(c, &pm, baseNode4)
	s.addABaseNode(c, &pm, baseNode42)
	s.addABaseNode(c, &pm, baseNode423)

	// TODO: randomize order in which finding baseNodes is tested
	s.findABaseNode(c, &pm, baseNode0123) // XXX

	s.findABaseNode(c, &pm, baseNode1)
	s.findABaseNode(c, &pm, baseNode12)
	s.findABaseNode(c, &pm, baseNode123)

	s.findABaseNode(c, &pm, baseNode4)
	s.findABaseNode(c, &pm, baseNode42)
	s.findABaseNode(c, &pm, baseNode423)

	s.findABaseNode(c, &pm, baseNode6)
	s.findABaseNode(c, &pm, baseNode62)
	s.findABaseNode(c, &pm, baseNode623)
}
