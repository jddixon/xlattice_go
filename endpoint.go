package xlattice_go

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 */

type EndPoint struct {
	transport *Transport
	address   *Address
}

func NewEndPoint(t *Transport, a *Address) *EndPoint {
	// XXX need some checks
	e := new(EndPoint)
	e.transport = (*Transport)(t)
	e.address = (*Address)(a)
	return e
}

func (e *EndPoint) getAddress() *Address {
	return e.address
}

func (e *EndPoint) getTransport() *Transport {
	return e.transport
}

// func (e *EndPoint) clone() *EndPoint {
//     NOT IMPLEMENTED
// }

func (e *EndPoint) ToString() string {
	// e.transport is a pointer to something that satisfies
	//   the Transport interface and similarly for e.address

	// probably not efficient
	s := (*e.transport).Name() + " " + (*e.address).ToString()
	return s
}
