package overlay

// xlattice_go/overlayI.go

import (
	xt "github.com/jddixon/xlattice_go/transport"
)

/**
 * A Overlay is characterized by an address space, a transport protocol,
 * and possibly a set of rules for navigating the address space using
 * the protocol.
 *
 * A Overlay may either be system-supported, like TCP/IP will
 * normally be, or it may explicitly depend upon an underlying
 * Overlay, in the way that HTTP, for example, is generally
 * implemented over TCP/IP.
 *
 * If the Overlay is system-supported, traffic will be routed and
 * neighbors will be reached by making calls to operating system
 * primitives such as sockets.
 *
 * In some Overlays there is a method which, given an EndPoint x, returns
 * another EndPoint g, a gateway, which can be used to route messages to
 * EndPoint x.
 *
 */

type OverlayI interface {
	Name() string // eg "eu-west-1.compute.amazonaws.com"
	IsElement(xt.EndPointI) bool
	Transport() string // eg "tcp"
	Cost() float32
	Equal(any interface{}) bool
	String() string
}
