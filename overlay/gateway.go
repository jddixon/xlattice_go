package overlay

// xlattice_go/overlay/gateway.go

import (
	x "github.com/jddixon/xlattice_go"
)

// This is meant to be a sort of map: the host named can provide
// transport to hosts in the overlays listed.
type Gateway struct {
	peer *x.Peer // who provides transport
	// XXX not going to work: we need a cost via this peer
	realm []*Overlay // overlays reachable through this peer
}
