package xlattice_go

import (
	"strings"
)

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 */

type EndPointI struct {
	transport *TransportI
	address   *AddressI
}

func NewEndPointI(t *TransportI, a *AddressI) *EndPointI {
	// XXX need some checks
	return &EndPointI{t, a}
}

func (e *EndPointI) getAddress() *AddressI {
	return e.address
}

func (e *EndPointI) getTransport() *TransportI {
	return e.transport
}

// func (e *EndPointI) clone() *EndPointI {
//     NOT IMPLEMENTED
// }

func (e *EndPointI) String() string {
	// e.transport is a pointer to something that satisfies
	//   the Transport interface and similarly for e.address
	var parts = []string{(*e.transport).String(), " ", (*e.address).String()}
	return strings.Join(parts, "")
}
