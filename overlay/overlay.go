package overlay

import (
	x "github.com/jddixon/xlattice_go"
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

type Overlay struct {
	name      string       // eg "eu-west-1.compute.amazonaws.com"
	addrRange *x.AddrRange // eg 10/8 in ipv4
	transport string       // eg "tcpip"
	cost      float32
}

func NewOverlay(name string, addrRange *x.AddrRange, transport string, cost float32) (*Overlay, error) {
	// XXX validate the parameters, please

	return &Overlay{name, addrRange, transport, cost}, nil
}

func (o *Overlay) Name() string {
	return o.name
}

func (o *Overlay) Equal(any interface{}) bool {
	// XXX STUB XXX
	return false
}
