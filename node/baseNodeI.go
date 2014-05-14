package node

// xlattice_go/node/baseNodeI.go

import (
	"crypto/rsa"
	xi "github.com/jddixon/xlNodeID_go"
	xo "github.com/jddixon/xlattice_go/overlay"
)

type BaseNodeI interface {
	GetName() string
	GetNodeID() *xi.NodeID
	GetCommsPublicKey() *rsa.PublicKey
	GetSSHCommsPublicKey() string
	GetSigPublicKey() *rsa.PublicKey
	// overlays
	AddOverlay(o xo.OverlayI) (ndx int, err error)
	SizeOverlays() int
	GetOverlay(n int) xo.OverlayI

	Equal(any interface{}) bool

	Strings() []string
	String() string
}
