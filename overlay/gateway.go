package overlay

// xlattice_go/overlay/gateway.go

import (
	x "github.com/jddixon/xlattice_go"
)

// This is meant to be a sort of map: the host named can provide
// transport to hosts in the overlays listed.
type Gateway struct {
	peer *x.Peer // which provides transport
	// XXX not going to work: we need a cost via this peer
	realms []*Overlay // overlays reachable through this peer
}

func NewGateway(p *x.Peer, o []*Overlay) *Gateway {
	// XXX validations
	sizeNow := len(o)
	r := make([]*Overlay, sizeNow, sizeNow+2)
	copy(r, o)
	return &Gateway{p, r}
}

func (g *Gateway) String() string {
	return "NOT YET IMPLEMENTED"
}
