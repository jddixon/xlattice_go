package overlay

// xlattice_go/overlay/ip_overlay.go

import (
	"fmt"
	xt "github.com/jddixon/xlattice_go/transport"
	"net"
	"strings"
)

var _ = fmt.Print // DEBUG

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
	addrRange *AddrRange // eg 10/8 in ipv4
	BaseOverlay
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
	// XXX should clone
	return o.addrRange
}

// XXX This belongs in ../transport
func CompatibleTransports(overlayT, endPointT string) bool {
	if overlayT == endPointT {
		return true
	}
	// more elaborate structure needed here
	if overlayT == "ip" && (endPointT == "tcp" || endPointT == "udp") {
		return true
	}
	return false
}
func (o *IPOverlay) IsElement(e xt.EndPointI) bool {
	oT := o.Transport()
	eT := e.Transport()
	if !CompatibleTransports(oT, eT) {
		return false
	}

	eA := e.Address().String()
	parts := strings.Split(eA, ":")

	bs := net.ParseIP(parts[0]) // returns an IP, a []byte
	if bs == nil {
		fmt.Printf("could not parse '%s'\n", eA)
		return false
	}
	return o.addrRange.ipNet.Contains(bs)
}
func (o *IPOverlay) String() string {
	return fmt.Sprintf("overlay: %s, %s, %s, %f", o.name,
		o.transport, o.addrRange.String(), o.cost)
}

func (o *IPOverlay) Equal(any interface{}) bool {
	if any == o {
		return true
	}
	if any == nil {
		return false
	}
	switch v := any.(type) {
	case *IPOverlay:
		_ = v
	default:
		return false
	}
	other := any.(*IPOverlay)
	if o.addrRange.String() != other.addrRange.String() {
		return false
	}
	return true
}
