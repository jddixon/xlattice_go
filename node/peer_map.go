package node

import (
	"fmt"
	xi "github.com/jddixon/xlattice_go/nodeID"
)

var _ = fmt.Print

type PeerMap struct {
	PeerMapCell
}
type PeerMapCell struct {
	byteVal byte
	pred    *PeerMapCell // predecessor
	nextCol *PeerMapCell // points to a cell with same val for this byte
	thisCol *PeerMapCell // points to a cell with higher val for this col
	peer    *Peer
}

// Add a Peer to the map.  This should be idempotent: adding a Peer
// that is already in the map should have no effect at all.  The cell map
// allows us to efficiently return a reference to a Peer, given its nodeID.

func (m *PeerMap) AddToPeerMap(peer *Peer) (err error) {
	id := peer.GetNodeID().Value()
	// don't make this check on the very first entry
	if m.nextCol != nil && m.FindPeer(id) != nil {
		// it's already present, so ignore
		return
	}
	byte0 := id[0]

	// DEBUG
	zero := make([]byte, xi.SHA1_LEN)
	zeroID, _ := xi.NewNodeID(zero)
	m.peer, _ = NewNewPeer("mapRoot", zeroID)
	// END

	root := m.nextCol
	if root == nil {
		m.nextCol = &PeerMapCell{
			byteVal: byte0,
			pred:    &m.PeerMapCell,
			peer:    peer}
	} else {
		err = root.addAtCell(0, peer, id)
	}
	return
}

// depth is that of cell, with the root cell at column 0, and also the
// index into the id slice.

func (p *PeerMapCell) addAtCell(depth int, peer *Peer, id []byte) (err error) {
	idByte := id[depth]
	// DEBUG
	//fmt.Printf("addAtCell: peer %s, depth %d, idByte %d, p.byteVal %d\n",
	//	peer.GetName(), depth, idByte, p.byteVal)
	// END
	if idByte < p.byteVal {
		// DEBUG
		//fmt.Printf("lower, adding %s as pred, idByte is %d\n",
		//	peer.GetName(), idByte)
		// END
		newCell := &PeerMapCell{
			byteVal: idByte, pred: p.pred, thisCol: p, peer: peer}
		if p.pred.thisCol == nil {
			// pred must be map's base cell
			p.pred.nextCol = newCell
		} else {
			p.pred.thisCol = newCell
		}
		p.pred = newCell

	} else if idByte == p.byteVal {
		// DEBUG
		//fmt.Printf("%s matches at depth %d, idByte is %d\n",
		//	peer.GetName(), depth, idByte)
		// END

		if p.nextCol == nil {
			if p.thisCol != nil {
				fmt.Printf("    thisCol is NOT nil\n")
				// XXX possible error ignored
				p.thisCol.addAtCell(depth, peer, id)
			} else {

				peer2 := p.peer
				var id2 []byte
				if peer2 != nil {
					id2 = peer2.GetNodeID().Value()
				}
				p.peer = nil

				depth++
				nextByte := id[depth]
				nextByte2 := id2[depth]
				curCell := p
				for nextByte == nextByte2 {
					nextCell := &PeerMapCell{byteVal: nextByte, pred: curCell}
					curCell.nextCol = nextCell
					curCell = nextCell
					depth++
					nextByte = id[depth]
					nextByte2 = id2[depth]
				}
				lastCell := &PeerMapCell{byteVal: nextByte, peer: peer}
				lastCell2 := &PeerMapCell{byteVal: nextByte2, peer: peer2}
				if nextByte < nextByte2 {
					curCell.nextCol = lastCell
					lastCell.pred = curCell
					lastCell.thisCol = lastCell2
					lastCell2.pred = lastCell
				} else {
					curCell.nextCol = lastCell2
					lastCell2.pred = curCell
					lastCell2.thisCol = lastCell
					lastCell.pred = lastCell2
				}
			}
		} else {
			// we had a match and we have a nextCol
			lastCell := p
			curCell := p.nextCol
			depth++
			// skip any cells with matching values
			for idByte = id[depth]; idByte == curCell.byteVal; idByte = id[depth] {
				if curCell.nextCol == nil {
					fmt.Printf("    nil nextCol at depth %d, breaking\n",
						depth)
					break
				}
				lastCell = curCell
				curCell = curCell.nextCol
				depth++
				fmt.Printf("    matched on %d; depth becomes %d\n",
					idByte, depth)
			}
			curByte := curCell.byteVal

			// nextCol is nil OR idByte doesn't match ----------------

			if curCell.thisCol != nil {
				// possible error ignored
				curCell.thisCol.addAtCell(depth, peer, id)
			} else {
				// nextCol may NOT be nil but thisCol is nil
				newCell := &PeerMapCell{byteVal: idByte, peer: peer}

				if idByte < curByte {
					// splice newCell in
					//fmt.Printf("    LESS: splicing new cell in at depth %d\n",
					//	depth)

					lastCell.nextCol = newCell
					newCell.pred = lastCell
					newCell.nextCol = nil
					newCell.thisCol = curCell
					curCell.pred = newCell

				} else {
					fmt.Printf("    GREATER: new cell off thisCol\n")
					if curCell.thisCol == nil {
						curCell.thisCol = newCell
						newCell.pred = curCell
					} else {
						// XXX possible error ignored
						curCell.thisCol.addAtCell(depth, peer, id)
					}

				}
			}
		}

	} else { // idByte > p.byteVal
		if p.thisCol == nil {
			p.addThisCol(id, depth, peer)
		} else {
			// XXX possible error ignored
			p.thisCol.addAtCell(depth, peer, id)
		}
	}

	return
}

// The nodeID of the peer being added has the same value for the byte at
// this depth.  id and peer represent the new peer being added, where id
// is the byte slice for its nodeID and peer is a reference to that.
// id2 and peer2 represent any pre-existing value.
//func (p *PeerMapCell) addMatchingToDepth(depth int,
//	id, id2 []byte, peer, peer2 *Peer) (err error) {
//
//	// XXX SHOULD NEVER SEE THIS, but do see it
//	fmt.Printf("ADD_MATCHING_TO_DEPTH, depth %d, peer %s\n",
//		depth, peer.GetName())
//
//	// The byte string id has matched the chain up to this point.
//	// We examine the next byte in id and the byte value for the next
//	// cell in the chain.
//	depth += 1
//	nextByte := id[depth]
//
//	if p.nextCol == nil {
//		if peer2 == nil {
//			p.nextCol = &PeerMapCell{nextByte, p, nil, nil, peer}
//		} else {
//			nextByte2 := id2[depth]
//			if nextByte == nextByte2 {
//				fmt.Printf("Case 1b1, %s\n", peer.GetName())
//				p.nextCol = &PeerMapCell{nextByte, p, nil, nil, nil}
//				p.nextCol.addMatchingToDepth(depth, id, id2, peer, peer2)
//			} else {
//				nextCell := &PeerMapCell{nextByte, nil, nil, nil, peer}
//				nextCell2 := &PeerMapCell{nextByte2, nil, nil, nil, peer2}
//				if nextByte < nextByte2 {
//					fmt.Printf("Case 1b2a, %s\n", peer.GetName())
//					nextCell.thisCol = nextCell2
//					p.nextCol = nextCell
//					nextCell.pred = p
//					nextCell2.pred = nextCell
//				} else {
//					fmt.Printf("Case 1b2b, %s\n", peer.GetName())
//					nextCell2.thisCol = nextCell
//					p.nextCol = nextCell2
//					nextCell2.pred = p
//					nextCell.pred = nextCell2
//				}
//			}
//		}
//	} else {
//		// XXX doesn't handle peer2
//		curCell := p.nextCol
//		if nextByte < curCell.byteVal {
//			// DEBUG
//			var nextPeerStr string
//			if curCell.peer == nil {
//				nextPeerStr = "<nil>"
//			} else {
//				nextPeerStr = curCell.peer.GetName()
//			}
//			fmt.Printf("CASE 2a: %s => %s\n", peer.GetName(), nextPeerStr)
//			// END
//			p.peer = peer
//			p.nextCol = &PeerMapCell{nextByte, p, nil, curCell, peer2}
//
//		} else if nextByte == curCell.byteVal {
//			fmt.Printf("CASE 2b, %s\n", peer.GetName()) // DEBUG
//			peer2 := curCell.peer
//			var id2 []byte
//			if peer2 != nil {
//				id2 = peer2.GetNodeID().Value()
//			}
//			curCell.peer = nil
//			curCell.addMatchingToDepth(depth, id, id2, peer, peer2)
//
//		} else { // nextByte > curCell.byteVal
//			fmt.Printf("CASE 2c, %s\n", peer.GetName()) // DEBUG
//			curCell.addThisCol(id, 0, peer)
//		}
//	}
//	return
//} // GEEP

func (p *PeerMapCell) addThisCol(id []byte, depth int, peer *Peer) (
	err error) {

	nextByte := id[depth]
	//fmt.Printf("addThisCol depth %d, nextByte %d, peer %s\n",
	//	depth, nextByte, peer.GetName())

	if p.thisCol == nil {
		p.thisCol = &PeerMapCell{byteVal: nextByte, pred: p, peer: peer}

		//// DEBUG
		//if p.peer == nil {
		//	fmt.Printf("    %s is sole cell down from <nil>\n", peer.GetName())
		//} else {
		//	fmt.Printf("    %s is sole cell down from %s\n",
		//		peer.GetName(), p.peer.GetName())
		//}
		//// END

	} else {
		// fmt.Println("    thisCol is NOT nil")
		// XXX ignoring possible error
		p.thisCol.addAtCell(depth, peer, id)
	}
	return
}

// At any particular depth, a match is possible only if (a) peer for the
// cell is not nil and (b) we have a byte-wise match

func (m *PeerMap) FindPeer(id []byte) (peer *Peer) {
	curCell := m.nextCol
	if curCell == nil { // no map
		return nil
	}
	// fmt.Printf("FindPeer for %d.%d.%d.%d\n", id[0], id[1], id[2], id[3])

	for depth := 0; depth < len(id); depth++ {
		myVal := id[depth]
		// fmt.Printf("    FindPeer: depth %d, val %d\n", depth, myVal)
		if curCell == nil {
			fmt.Printf("    Internal error: nil curCell at depth %d\n", depth)
			return nil
		}
		if myVal > curCell.byteVal {
			for curCell.thisCol != nil {
				curCell = curCell.thisCol
				if myVal == curCell.byteVal {
					goto maybeEqual
				} else if myVal > curCell.byteVal {
					continue
				}
				break
			}
			//fmt.Printf("    depth %d, %d < %d returning NIL\n",
			//	depth, myVal, curCell.byteVal)
			return nil
		}
	maybeEqual:
		if myVal == curCell.byteVal {
			if curCell.nextCol == nil {
				myNodeID, err := xi.NewNodeID(id)
				if err != nil {
					fmt.Printf("    FindPeer: NewNodeID returns %v", err)
					return nil
				}
				if curCell.peer != nil {
					// fmt.Printf("    peer is %s\n", curCell.peer.GetName())
					if xi.SameNodeID(myNodeID, curCell.peer.GetNodeID()) {
						// fmt.Printf("    *MATCH* on %s\n", curCell.peer.GetName())
						return curCell.peer
					}
				}
			} else {
				// fmt.Printf("    RIGHT, so depth := %d\n", depth+1)
				curCell = curCell.nextCol
				continue
			}

		} else {
			// myVal < curCell.byteVal
			//fmt.Printf("    myval %d > cell's %d\n", myVal, curCell.byteVal)
			return nil
		}
	}
	return
}
