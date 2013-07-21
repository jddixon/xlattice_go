package transport

import (
	"net"
)

/**
 * An EndPoint is specified by a transport and an Address, including
 * the local part.  If the transport is TCP/IP, for example, the
 * Address includes the IP address and the port number.
 *
 */

type TcpEndPoint struct {
	tcpAddr *net.TCPAddr // IP, Port, Zone
}

func NewTcpEndPoint(addr string) (*TcpEndPoint, error) {
	a, err := net.ResolveTCPAddr("tcp", addr)
	if err == nil {
		return &TcpEndPoint{a}, nil
	} else {
		return nil, err
	}
}

func (e *TcpEndPoint) GetAddress() AddressI {
	a, _ := NewV4Address(e.tcpAddr.String())
	return a
}

func (e *TcpEndPoint) GetTransport() string {
	return "tcp"
}

func (e *TcpEndPoint) Clone() (*TcpEndPoint, error) {
	return NewTcpEndPoint(e.GetAddress().String())
}

func (e *TcpEndPoint) String() string {
	return e.tcpAddr.String()
}

// net.Addr interface ///////////////////////////////////////////////

// This is just an alias for GetTransport
func (e *TcpEndPoint) Network() string {
	return e.GetTransport()
}

// Shortcut for Go
func (e *TcpEndPoint) GetTcpAddr() *net.TCPAddr {
	return e.tcpAddr
}
