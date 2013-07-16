package overlay

// xlattice_go/overlay/ip_overlay.go

import ()

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

type IPOverlay struct {
	addrRange *AddrRange // eg 10/8 in ipv4		XXX THIS IS AN ERROR
	Overlay
}

func NewIPOverlay(name string, addrRange *AddrRange, transport string, cost float32) (*IPOverlay, error) {
	// XXX validate the parameters, please

	overlay, err := New(name, transport, cost)
	if err == nil {
		// XXX validate addrRange
		ipOverlay := IPOverlay{addrRange, *overlay}
		return &ipOverlay, nil
	} else {
		return nil, err
	}
}

func (o *IPOverlay) AddrRange() *AddrRange {
	return o.addrRange
}
