package node

import ()

type PeerMap struct {
	lowest *PeerMapCell
}
type PeerMapCell struct {
	byteVal byte
	sameVal *PeerMapCell // points to a cell with same val for this byte
	higher  *PeerMapCell // points to a cell with higher val for this byte
	peer    *Peer
}

// Add a Peer to the map.  This should be idempotent: adding a Peer
// that is already in the map should have no effect at all.  The cell map
// allows us to very efficiently return a reference to a Peer, given its
// nodeID.

func (m *PeerMap) AddToPeerMap(peer *Peer) (err error) {
	id := peer.GetNodeID().Value()
	byte0 := id[0]

	root := m.lowest
	if root == nil || byte0 < root.byteVal {
		root = &PeerMapCell{byte0, nil, root, peer}
	} else if byte0 == root.byteVal {
		root.AddSibling(id, 0, peer)
	} else { // byte0 > root.byteVal
		root.AddHigher(id, 0, peer)
	}
	return
}

// The nodeID of the peer being added has the same value for this depth.
//
func (p *PeerMapCell) AddSibling(id []byte, depth int, peer *Peer) (err error) {
	depth += 1
	nextByte := id[depth]

	if p.sameVal == nil {
		p.sameVal = &PeerMapCell{nextByte, nil, nil, peer}
	} else {
		curSib := p.sameVal // current sibling
		if nextByte < curSib.byteVal {
			curSib = &PeerMapCell{nextByte, nil, curSib, peer}
		} else if nextByte == curSib.byteVal {
			curSib.AddSibling(id, 0, peer)
		} else { // nextByte > curSib.byteVal
			curSib.AddHigher(id, 0, peer)
		}
	}
	return
}

func (p *PeerMapCell) AddHigher(id []byte, depth int, peer *Peer) (err error) {
	nextByte := id[depth]

	if p.sameVal == nil {
		p.sameVal = &PeerMapCell{nextByte, nil, nil, peer}
	} else {
		curHigher := p.higher // current higher value
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

func (m *PeerMap) FindPeer(id []byte) (peer *Peer) {
	curCell := m.lowest
	for depth := 0; depth < len(id); depth++ {
		// WORKING HERE
		_ = curCell
	}
	return
}
