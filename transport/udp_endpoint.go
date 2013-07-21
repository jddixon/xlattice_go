package transport

import (
	"net"
)

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is UDP, for example, the
 * Address includes the IP address and the port number.
 *
 */

type UdpEndPoint struct {
	udpAddr *net.UDPAddr // IP, Port, Zone
}

func NewUdpEndPoint(addr string) (*UdpEndPoint, error) {
	a, err := net.ResolveUDPAddr("udp", addr)
	if err == nil {
		return &UdpEndPoint{a}, nil
	} else {
		return nil, err
	}
}

func (e *UdpEndPoint) GetAddress() AddressI {
	a, _ := NewV4Address(e.udpAddr.String())
	return a
}

func (e *UdpEndPoint) GetTransport() string {
	return "udp"
}

func (e *UdpEndPoint) Clone() (*UdpEndPoint, error) {
	return NewUdpEndPoint(e.GetAddress().String())
}

func (e *UdpEndPoint) String() string {
	return e.udpAddr.String()
}

// net.Addr interface ///////////////////////////////////////////////

// This is just an alias for GetTransport
func (e *UdpEndPoint) Network() string {
	return e.GetTransport()
}

// Shortcut for Go
func (e *UdpEndPoint) GetUdpAddr() *net.UDPAddr {
	return e.udpAddr
}
