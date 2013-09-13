package node

import (
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
)

var _ = fmt.Print

///////////////////////
// PROBABLY JUST A HACK
///////////////////////

// This was PeerMap until I recognized that we were only using BaseNode
// attributes.  So I crudely renamed Peer to BaseNode, peer to baseNode,
// and so forth throughout.  It does pass its tests.

type BaseNodeMap struct {
	BaseNodeMapCell
}
type BaseNodeMapCell struct {
	ByteVal  byte
	Pred     *BaseNodeMapCell // predecessor
	NextCol  *BaseNodeMapCell // points to a cell with same val for this byte
	ThisCol  *BaseNodeMapCell // points to a cell with higher val for this col
	CellNode *BaseNode
}

// Add a BaseNode to the map.  This should be idempotent: adding a BaseNode
// that is already in the map should have no effect at all.  The cell map
// allows us to efficiently return a reference to a BaseNode, given its nodeID.

func (m *BaseNodeMap) AddToBaseNodeMap(baseNode *BaseNode) (err error) {
	id := baseNode.GetNodeID().Value()
	// don't make this check on the very first entry
	if m.NextCol != nil && m.FindBaseNode(id) != nil {
		// it's already present, so ignore
		return
	}
	byte0 := id[0]

	// DEBUG
	zero := make([]byte, xi.SHA1_LEN)
	zeroID, _ := xi.NewNodeID(zero)
	m.CellNode, _ = NewNewBaseNode("mapRoot", zeroID)
	// END

	root := m.NextCol
	if root == nil {
		m.NextCol = &BaseNodeMapCell{
			ByteVal:  byte0,
			Pred:     &m.BaseNodeMapCell,
			CellNode: baseNode}
	} else {
		err = root.addAtCell(0, baseNode, id)
	}
	return
}

// depth is that of cell, with the root cell at column 0, and also the
// index into the id slice.

func (p *BaseNodeMapCell) addAtCell(depth int, baseNode *BaseNode, id []byte) (err error) {
	idByte := id[depth]
	// DEBUG
	//fmt.Printf("addAtCell: baseNode %s, depth %d, idByte %d, p.ByteVal %d\n",
	//	CellNode.GetName(), depth, idByte, p.ByteVal)
	// END
	if idByte < p.ByteVal {
		// DEBUG
		//fmt.Printf("lower, adding %s as pred, idByte is %d\n",
		//	CellNode.GetName(), idByte)
		// END
		newCell := &BaseNodeMapCell{
			ByteVal: idByte, Pred: p.Pred, ThisCol: p, CellNode: baseNode}
		if p.Pred.ThisCol == nil {
			// pred must be map's base cell
			p.Pred.NextCol = newCell
		} else {
			p.Pred.ThisCol = newCell
		}
		p.Pred = newCell

	} else if idByte == p.ByteVal {
		// DEBUG
		//fmt.Printf("%s matches at depth %d, idByte is %d\n",
		//	CellNode.GetName(), depth, idByte)
		// END

		if p.NextCol == nil {
			if p.ThisCol != nil {
				fmt.Printf("    ThisCol is NOT nil\n")
				// XXX possible error ignored
				p.ThisCol.addAtCell(depth, baseNode, id)
			} else {

				baseNode2 := p.CellNode
				var id2 []byte
				if baseNode2 != nil {
					id2 = baseNode2.GetNodeID().Value()
				}
				p.CellNode = nil

				depth++
				nextByte := id[depth]
				nextByte2 := id2[depth]
				curCell := p
				for nextByte == nextByte2 {
					nextCell := &BaseNodeMapCell{ByteVal: nextByte, Pred: curCell}
					curCell.NextCol = nextCell
					curCell = nextCell
					depth++
					nextByte = id[depth]
					nextByte2 = id2[depth]
				}
				lastCell := &BaseNodeMapCell{ByteVal: nextByte, CellNode: baseNode}
				lastCell2 := &BaseNodeMapCell{ByteVal: nextByte2, CellNode: baseNode2}
				if nextByte < nextByte2 {
					curCell.NextCol = lastCell
					lastCell.Pred = curCell
					lastCell.ThisCol = lastCell2
					lastCell2.Pred = lastCell
				} else {
					curCell.NextCol = lastCell2
					lastCell2.Pred = curCell
					lastCell2.ThisCol = lastCell
					lastCell.Pred = lastCell2
				}
			}
		} else {
			// we had a match and we have a NextCol
			lastCell := p
			curCell := p.NextCol
			depth++
			// skip any cells with matching values
			for idByte = id[depth]; idByte == curCell.ByteVal; idByte = id[depth] {
				if curCell.NextCol == nil {
					fmt.Printf("    nil NextCol at depth %d, breaking\n",
						depth)
					break
				}
				lastCell = curCell
				curCell = curCell.NextCol
				depth++
				fmt.Printf("    matched on %d; depth becomes %d\n",
					idByte, depth)
			}
			curByte := curCell.ByteVal

			// NextCol is nil OR idByte doesn't match ----------------

			if curCell.ThisCol != nil {
				// possible error ignored
				curCell.ThisCol.addAtCell(depth, baseNode, id)
			} else {
				// NextCol may NOT be nil but ThisCol is nil
				newCell := &BaseNodeMapCell{ByteVal: idByte, CellNode: baseNode}

				if idByte < curByte {
					// splice newCell in
					//fmt.Printf("    LESS: splicing new cell in at depth %d\n",
					//	depth)

					lastCell.NextCol = newCell
					newCell.Pred = lastCell
					newCell.NextCol = nil
					newCell.ThisCol = curCell
					curCell.Pred = newCell

				} else {
					fmt.Printf("    GREATER: new cell off ThisCol\n")
					if curCell.ThisCol == nil {
						curCell.ThisCol = newCell
						newCell.Pred = curCell
					} else {
						// XXX possible error ignored
						curCell.ThisCol.addAtCell(depth, baseNode, id)
					}

				}
			}
		}

	} else { // idByte > p.ByteVal
		if p.ThisCol == nil {
			p.addThisCol(id, depth, baseNode)
		} else {
			// XXX possible error ignored
			p.ThisCol.addAtCell(depth, baseNode, id)
		}
	}

	return
}

func (p *BaseNodeMapCell) addThisCol(id []byte, depth int, baseNode *BaseNode) (
	err error) {

	nextByte := id[depth]
	//fmt.Printf("addThisCol depth %d, nextByte %d, baseNode %s\n",
	//	depth, nextByte, CellNode.GetName())

	if p.ThisCol == nil {
		p.ThisCol = &BaseNodeMapCell{ByteVal: nextByte, Pred: p, CellNode: baseNode}

		//// DEBUG
		//if p.CellNode == nil {
		//	fmt.Printf("    %s is sole cell down from <nil>\n", CellNode.GetName())
		//} else {
		//	fmt.Printf("    %s is sole cell down from %s\n",
		//		CellNode.GetName(), p.CellNode.GetName())
		//}
		//// END

	} else {
		// fmt.Println("    ThisCol is NOT nil")
		// XXX ignoring possible error
		p.ThisCol.addAtCell(depth, baseNode, id)
	}
	return
}

// At any particular depth, a match is possible only if (a) baseNode for the
// cell is not nil and (b) we have a byte-wise match

func (m *BaseNodeMap) FindBaseNode(id []byte) (baseNode *BaseNode) {
	curCell := m.NextCol
	if curCell == nil { // no map
		return nil
	}
	// fmt.Printf("FindBaseNode for %d.%d.%d.%d\n", id[0], id[1], id[2], id[3])

	for depth := 0; depth < len(id); depth++ {
		myVal := id[depth]
		// fmt.Printf("    FindBaseNode: depth %d, val %d\n", depth, myVal)
		if curCell == nil {
			fmt.Printf("    Internal error: nil curCell at depth %d\n", depth)
			return nil
		}
		if myVal > curCell.ByteVal {
			for curCell.ThisCol != nil {
				curCell = curCell.ThisCol
				if myVal == curCell.ByteVal {
					goto maybeEqual
				} else if myVal > curCell.ByteVal {
					continue
				}
				break
			}
			//fmt.Printf("    depth %d, %d < %d returning NIL\n",
			//	depth, myVal, curCell.ByteVal)
			return nil
		}
	maybeEqual:
		if myVal == curCell.ByteVal {
			if curCell.NextCol == nil {
				myNodeID, err := xi.NewNodeID(id)
				if err != nil {
					fmt.Printf("    FindBaseNode: NewNodeID returns %v", err)
					return nil
				}
				if curCell.CellNode != nil {
					// fmt.Printf("    baseNode is %s\n", curCell.CellNode.GetName())
					if xi.SameNodeID(myNodeID, curCell.CellNode.GetNodeID()) {
						// fmt.Printf("    *MATCH* on %s\n", curCell.CellNode.GetName())
						return curCell.CellNode
					}
				}
			} else {
				// fmt.Printf("    RIGHT, so depth := %d\n", depth+1)
				curCell = curCell.NextCol
				continue
			}

		} else {
			// myVal < curCell.ByteVal
			//fmt.Printf("    myval %d > cell's %d\n", myVal, curCell.ByteVal)
			return nil
		}
	}
	return
}
