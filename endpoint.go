package xlattice_go

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 */

type EndPoint struct {
	transport string // should be pointer to Transport
	address   *Address
}

func NewEndPoint(t string, a *Address) *EndPoint {
	// XXX need some checks
	return &EndPoint{t, a}
}

func (e *EndPoint) getAddress() *Address {
	return e.address
}

//func (e *EndPoint) getTransport() *Transport {
func (e *EndPoint) getTransport() string {
	return e.transport
}

// func (e *EndPoint) clone() *EndPoint {
//     NOT IMPLEMENTED
// }

func (e *EndPoint) ToString() string {
	// e.transport is a pointer to something that satisfies
	//   the Transport interface and similarly for e.address

	// probably not efficient
	s := e.transport + " " + (*e.address).ToString()
	return s
}
