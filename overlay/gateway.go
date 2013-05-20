package overlay

// xlattice_go/overlay/gateway.go

import (
	x "github.com/jddixon/xlattice_go"
)

// This is meant to be a sort of map: the host named can provide
// transport to hosts in the realm = the destination overlay.
type Gateway struct {
	host  *x.NodeID
	realm *x.Overlay
}
