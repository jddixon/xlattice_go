package node

// xlattice_go/overlay/node.go

import (
	xo "github.com/jddixon/xlattice_go/overlay"
)

// This is meant to be a sort of map: the host named can provide
// transport to hosts in the overlays listed.
type Gateway struct {
	peer *Peer // which provides transport
	// XXX not going to work: we need a cost via this peer
	realms []xo.OverlayI // overlays reachable through this peer
}

func NewGateway(p *Peer, o []xo.OverlayI) *Gateway {
	// XXX validations
	sizeNow := len(o)
	r := make([]xo.OverlayI, sizeNow, sizeNow+2)
	copy(r, o)
	return &Gateway{p, r}
}

func (g *Gateway) String() string {
	return "NOT YET IMPLEMENTED"
}
