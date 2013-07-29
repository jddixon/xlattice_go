package transport

import (
	"fmt"
	"net"
)

var _ = fmt.Print // DEBUG

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

func (e *TcpEndPoint) Address() AddressI {
	// return a copy
	a, _ := NewV4Address(e.tcpAddr.String())
	return a
}

func (e *TcpEndPoint) Clone() (ep EndPointI, err error) {
	return NewTcpEndPoint(e.Address().String())
}

func (e *TcpEndPoint) Equal(any interface{}) bool {
	if any == nil {
		fmt.Println("Equal: nil other") // DEBUG
		return false
	}
	if any == e {
		return true
	}
	switch v := any.(type) {
	case *TcpEndPoint:
		_ = v
	default:
		fmt.Println("Equal: other not *TcpEndPoint") // DEBUG
		return false
	}
	other := any.(*TcpEndPoint)
	t, ot := e.tcpAddr, other.tcpAddr
	if len(t.IP) != len(ot.IP) {
		fmt.Println("Equal: other has different len(IP)") // DEBUG
		return false
	}
	for i := 0; i < len(t.IP); i++ {
		if t.IP[i] != ot.IP[i] {
			fmt.Println("Equal: other has different IP[i]") // DEBUG
			return false
		}
	}
	return t.Port == ot.Port && t.Zone == ot.Zone
}

func (e *TcpEndPoint) String() string {
	return e.tcpAddr.String()
}

func (e *TcpEndPoint) Transport() string {
	return "tcp"
}

// net.Addr interface ///////////////////////////////////////////////

// This is just an alias for Transport
func (e *TcpEndPoint) Network() string {
	return e.Transport()
}

// Shortcut for Go
func (e *TcpEndPoint) GetTcpAddr() *net.TCPAddr {
	return e.tcpAddr
}
