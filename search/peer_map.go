package search

import (
	"fmt"
	xn "github.com/jddixon/xlattice_go/node"
)

var _ = fmt.Print

type PeerMap struct {
	lowest *PeerMapCell
}
type PeerMapCell struct {
	byteVal byte
	nextCol *PeerMapCell // points to a cell with same val for this byte
	thisCol *PeerMapCell // points to a cell with higher val for this byte
	peer    *xn.Peer
}

// Add a Peer to the map.  This should be idempotent: adding a Peer
// that is already in the map should have no effect at all.  The cell map
// allows us to efficiently return a reference to a Peer, given its nodeID.

func (m *PeerMap) AddToPeerMap(peer *xn.Peer) (err error) {
	id := peer.GetNodeID().Value()
	byte0 := id[0]

	root := m.lowest
	if root == nil {
		fmt.Printf("empty, adding %s as lowest, byte0 is %d\n", peer.GetName(), byte0)
		m.lowest = &PeerMapCell{byte0, nil, nil, peer}

	} else if byte0 < root.byteVal {
		fmt.Printf("lower, adding %s as lowest, byte0 is %d\n", peer.GetName(), byte0)
		m.lowest = &PeerMapCell{byte0, nil, root, peer}

	} else if byte0 == root.byteVal {
		// THIS DOESN'T WORK AS IT IS:

		// If there is a non-nil pointer in the peer field, we are at
		// the end of the chain.  We remove the pointer and then search
		// forward until we find a differing byte, leaving a trail of
		// cells.  When we find the byte that differs, we want this cell's
		// nextCol to point to the cell for the lower-valued peer and
		// that cell's thisCol to point to the higher-valued peer.
		// We MUST NOT ever run all the way down the chain.
		// XXX STUB XXX

		// Otherwise (the peer field is non-nil) we are not at the
		// end of the chain.
		// XXX STUB XXX

		// OLD CODE, KNOWN TO BE WRONG
		fmt.Printf("adding %s as sibling, byte0 is %d\n", peer.GetName(), byte0)
		root.AddSibling(id, 0, peer)

	} else { // byte0 > root.byteVal
		fmt.Printf("adding %s as higher, byte0 is %d\n", peer.GetName(), byte0)
		root.AddHigher(id, 0, peer)
	}
	return
}

// The nodeID of the peer being added has the same value for this depth.
//
func (p *PeerMapCell) AddSibling(id []byte, depth int, peer *xn.Peer) (
	err error) {

	depth += 1
	nextByte := id[depth]

	if p.nextCol == nil {
		p.nextCol = &PeerMapCell{nextByte, nil, nil, peer}
	} else {
		curSib := p.nextCol // current sibling
		if nextByte < curSib.byteVal {
			curSib = &PeerMapCell{nextByte, nil, curSib, peer}

		} else if nextByte == curSib.byteVal {
			// WON'T WORK, need to look ahead to see which sorts lower
			curSib.AddSibling(id, 0, peer)

		} else { // nextByte > curSib.byteVal
			curSib.AddHigher(id, 0, peer)
		}
	}
	return
}

func (p *PeerMapCell) AddHigher(id []byte, depth int, peer *xn.Peer) (
	err error) {

	nextByte := id[depth]
	if p.nextCol == nil {
		p.nextCol = &PeerMapCell{nextByte, nil, nil, peer}
	} else {
		curHigher := p.thisCol // current higher value
		if nextByte < curHigher.byteVal {
			curHigher = &PeerMapCell{nextByte, nil, curHigher, peer}
		} else if nextByte == curHigher.byteVal {
			curHigher.AddSibling(id, 0, peer)
		} else { // nextByte > curHigher.byteVal
			curHigher.AddHigher(id, 0, peer)
		}
	}
	return
}

func (m *PeerMap) FindPeer(id []byte) (peer *xn.Peer) {
	mapCell := m.lowest
	for depth := 0; depth < len(id); depth++ {
		// continue to check sibling

		myVal := id[depth]
		if myVal < mapCell.byteVal {
			return nil
		} else if myVal == mapCell.byteVal {
			if mapCell.nextCol == nil {
				return mapCell.peer
			} else {
				// we have a sibling - check it
				mapCell = mapCell.nextCol
				continue
			}
		} else {
			// myVal > mapCell.byteVal
			for mapCell = mapCell.thisCol; mapCell != nil; mapCell = mapCell.thisCol {
				if myVal < mapCell.byteVal {
					return
				} else if myVal == mapCell.byteVal {
					mapCell = mapCell.nextCol
					break
				} else {
					continue // down the chain of higher values
				}
			}
			continue // along the chain of matching values
		}
	}
	return
}
